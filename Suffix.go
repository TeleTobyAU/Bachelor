package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"time"
)

type Info struct {
	input      string
	alphabet   []string
	threshHold int
	SA         []int
	RSA        []int
	cTable     []int
	oTable     [][]int
	roTable    [][]int
}

type bwtApprox struct {
	bwtTable           *Info
	key                string
	L, R, nextInterval int
	Ls                 []int
	Rs                 []int
	cigar              []string
	keyLength          int
	editBuff           []rune
	dTable             []int
	matchLengths       []int
}

type SaisStruct struct {
	SA             []int
	LSTypes        string
	buckets        []int
	beginnings     []int
	ends           []int
	LMSAlphabet    []string
	summaryOffsets []int
	summaryString  []int
	newStrLen      int
	newAlphSize    int
}

const UNDEFINED = int(^uint(0) >> 1)

func main() {
	info := new(Info)
	generateRandomNucleotide(1000000000, info)
	fmt.Println(len(info.input))

	//info.input = "GTCGGTATCGGTGGGCGTGCGCCAACCTGGGCAGAGTTGATTCTTGCTTTCCCGCTCATACTACATCCGGAAGCAGATCCAGGCGACCGGAACCGAGCGC$"
	//info.input = "mmiissiissiippii$"

	info.threshHold = 1

	generateAlphabet(info)

	generateCTable(info)
	info.SA = SAIS(info, info.input)

}

func generateRandomNucleotide(size int, info *Info) {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("ATCG")

	nucleotide := make([]rune, size)

	for i := range nucleotide {
		nucleotide[i] = letters[rand.Intn(len(letters))]
	}
	info.input = string(nucleotide) + "$"
}

/**
Linear suffix array construction by almost pure induced sorting.
Algorithm derived from Zhang and Chan 2009
*/
func SAIS(info *Info, n string) []int {
	start := time.Now()
	saisStruct := new(SaisStruct)

	sortSA(info, saisStruct, n)
	fmt.Println("total", time.Since(start))

	return saisStruct.SA
}

func classifyLS(n string) string {
	LSTypes := make([]rune, len(n))
	LSTypes[len(n)-1] = 'S'

	for i := len(n) - 2; i >= 0; i-- {
		if n[i] == n[i+1] {
			LSTypes[i] = LSTypes[i+1]
		} else {
			if n[i] > n[i+1] {
				LSTypes[i] = 'L'
			} else {
				LSTypes[i] = 'S'
			}
		}
	}
	return string(LSTypes)
}

func isLMSIndex(LSString string, i int) bool {
	if i == 0 {
		return false
	} else {
		return LSString[i] == 'S' && LSString[i-1] != 'S'
	}
}

func placeLMS(saisStruct *SaisStruct, info *Info, n string) []int {
	//Initialize SA to UNDEFINED
	SA := make([]int, len(n))
	for i := range SA {
		SA[i] = UNDEFINED
	}

	//SA-IS step 1, placing LMS substrings in saisStruct
	for i := 0; i < len(n); i++ {
		if isLMSIndex(saisStruct.LSTypes, i) {
			remappedi := IndexOf(string(n[i]), info.alphabet)
			saisStruct.ends[remappedi]--
			SA[saisStruct.ends[remappedi]] = i
		}
	}
	return SA
}

func induceL(saisStruct *SaisStruct, info *Info, n string) []int {
	SA := saisStruct.SA
	for i := 0; i < len(n); i++ {
		if SA[i] == UNDEFINED {
			continue
		}

		if SA[i] == 0 {
			continue
		}

		j := SA[i] - 1
		if saisStruct.LSTypes[j] == 'L' {
			remappedi := IndexOf(string(n[j]), info.alphabet)
			SA[saisStruct.beginnings[remappedi]-1] = j
			saisStruct.beginnings[remappedi]++
		}
	}
	return SA
}

func induceS(saisStruct *SaisStruct, info *Info, n string) []int {
	//SA-IS step 3, placing and sorting S types
	SA := saisStruct.SA
	saisStruct.ends = bucketEnds(info, saisStruct)
	for i := len(n); i > 0; i-- {
		if saisStruct.SA[i-1] == 0 {
			continue
		}

		j := SA[i-1] - 1

		if saisStruct.LSTypes[j] == 'S' {
			remappedi := IndexOf(string(n[j]), info.alphabet)
			saisStruct.ends[remappedi]--
			SA[saisStruct.ends[remappedi]] = j
		}

	}
	return SA
}

func computeBuckets(info *Info) []int {
	buckets := make([]int, len(info.alphabet))
	for i := 0; i < len(info.input); i++ {
		remappedi := IndexOf(string(info.input[i]), info.alphabet)
		if remappedi != -1 {
			buckets[remappedi]++
		}
	}
	return buckets
}

func bucketEnds(info *Info, saisStruct *SaisStruct) []int {
	ends := make([]int, len(info.alphabet))
	ends[0] = saisStruct.buckets[0]
	for i := 1; i < len(info.alphabet); i++ {
		ends[i] = ends[i-1] + saisStruct.buckets[i]
	}
	return ends
}

func bucketBeginnings(info *Info, saisStruct *SaisStruct) []int {
	beginnings := make([]int, len(info.alphabet))
	beginnings[0] = saisStruct.buckets[0]
	for i := 1; i < len(info.alphabet); i++ {
		beginnings[i] = beginnings[i-1] + saisStruct.buckets[i-1]
	}
	return beginnings
}

func equalLMS(saisStruct *SaisStruct, n string, i int, j int) bool {
	if i == len(n) || j == len(n) {
		return false
	}
	k := 0
	for {
		iLMS := isLMSIndex(saisStruct.LSTypes, i+k)
		jLMS := isLMSIndex(saisStruct.LSTypes, j+k)
		if k > 0 && iLMS && jLMS {
			return true
		}
		if iLMS != jLMS || n[i+k] != n[j+k] || (saisStruct.LSTypes[i+k] == 'S') != (saisStruct.LSTypes[j+k] == 'S') {
			return false
		}
		k++
	}
}

func reduceSA(saisStruct *SaisStruct, n string) {
	saisStruct.summaryOffsets = make([]int, len(n))
	name := 0
	names := make([]int, len(n)+1)
	for i := range names {
		names[i] = UNDEFINED
	}
	names[saisStruct.SA[0]] = name
	lastS := saisStruct.SA[0]

	for i := 1; i < len(n); i++ {
		j := saisStruct.SA[i]
		if !isLMSIndex(saisStruct.LSTypes, j) {
			continue
		}
		if !equalLMS(saisStruct, n, lastS, j) {
			name++
		}
		lastS = j
		names[j] = name
	}
	saisStruct.newAlphSize = name + 1

	j := 0
	for i := 0; i < len(n); i++ {
		name = names[i]
		if name == UNDEFINED {
			continue
		}
		saisStruct.summaryOffsets[j] = i
		saisStruct.summaryString = append(saisStruct.summaryString, name)
		j++
	}

	saisStruct.newStrLen = j - 1
}

func recursiveSorting(info *Info, saisStruct *SaisStruct, n string) []int {
	start := time.Now()
	saisStruct.LSTypes = classifyLS(n)
	fmt.Println("LS")
	fmt.Println(time.Since(start))
	start = time.Now()

	saisStruct.buckets = computeBuckets(info)
	fmt.Println("buckets")
	fmt.Println(time.Since(start))
	start = time.Now()

	saisStruct.beginnings = bucketBeginnings(info, saisStruct)
	fmt.Println("beginning")
	fmt.Println(time.Since(start))
	start = time.Now()

	saisStruct.ends = bucketEnds(info, saisStruct)
	fmt.Println("ends")
	fmt.Println(time.Since(start))
	start = time.Now()

	saisStruct.SA = placeLMS(saisStruct, info, n)
	fmt.Println("placelms")
	fmt.Println(time.Since(start))
	start = time.Now()

	saisStruct.SA = induceL(saisStruct, info, n)
	fmt.Println("inducel")
	fmt.Println(time.Since(start))
	start = time.Now()

	saisStruct.beginnings = bucketBeginnings(info, saisStruct)
	fmt.Println("beginnings2")
	fmt.Println(time.Since(start))
	start = time.Now()

	saisStruct.SA = induceS(saisStruct, info, n)
	fmt.Println("induces")
	fmt.Println(time.Since(start))
	start = time.Now()

	reduceSA(saisStruct, n)
	fmt.Println("reduce")
	fmt.Println(time.Since(start))
	start = time.Now()
	SA := saisStruct.SA

	//SA = sortSA(info, saisStruct, n)

	fmt.Println(time.Since(start))

	return SA
}

func sortSA(info *Info, saisStruct *SaisStruct, n string) []int {
	SA := saisStruct.SA
	if len(n) == 0 {
		SA[0] = 0
		return SA
	}

	if len(info.alphabet) == len(n)+1 {
		SA[0] = len(n)
		for i := 0; i < len(n); i++ {
			j := IndexOf(n[i], info.alphabet)
			SA[j] = i
		}
	} else {
		recursiveSorting(info, saisStruct, n)
	}

	return SA
}

func remapLMS(info *Info, saisStruct *SaisStruct, n string) []int {
	SA := saisStruct.SA
	saisStruct.ends = bucketEnds(info, saisStruct)
	for i := saisStruct.newStrLen + 1; i > 0; i-- {
		idx := saisStruct.summaryOffsets[saisStruct.summaryString[i-1]]
		bucketIdx := IndexOf(string(n[idx]), info.alphabet)
		SA[saisStruct.ends[bucketIdx]] = idx
	}
	SA[0] = len(n) - 1

	return SA
}

//https://stackoverflow.com/questions/1752414/how-to-reverse-a-string-in-go
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func initBwtApproxIter(key string, maxEdit int, info *Info, approx *bwtApprox) {
	//Init struct bwt_Approx
	approx.bwtTable = info
	approx.key = key

	//Set up edits buffer.
	keyLength := len(approx.key)
	approx.keyLength = keyLength

	//Building D table
	generateDTable(approx, info)

	//Start searching
	L := 0
	R := len(info.SA)
	i := keyLength - 1
	edits := &approx.editBuff

	//X- and = -operation
	aMatch := IndexOf(string(key[i]), info.alphabet)

	for a := 1; a < len(info.alphabet); a++ {
		newL := info.cTable[a] + info.oTable[a][L]
		newR := info.cTable[a] + info.oTable[a][R]

		var editCost int
		if a == aMatch {
			editCost = 0
		} else {
			editCost = 1
		}

		if maxEdit-editCost < 0 {
			continue
		}
		if newL >= newR {
			continue
		}

		if editCost == 1 {
			*edits = append(*edits, 'X')
		} else {
			*edits = append(*edits, '=')
		}
		recApproxMatching(approx, newL, newR, i-1, 1, maxEdit-editCost, edits)
		*edits = (*edits)[:len(*edits)-1]

	}

	// I-operation
	*edits = append(*edits, 'I')
	recApproxMatching(approx, L, R, i-1, 0, maxEdit-1, edits)
	*edits = (*edits)[:len(*edits)-1]

	// Make sure we start at the first interval.
	approx.L = keyLength
	approx.R = 0
	approx.nextInterval = 0
}

func recApproxMatching(approx *bwtApprox, L int, R int, i int, matchLength int, editLeft int, edits *[]rune) {
	//initializing variables for rec approx
	C := approx.bwtTable.cTable
	O := approx.bwtTable.oTable
	alphabet := approx.bwtTable.alphabet
	var lowerLimit int
	var revEdits []rune

	if i >= 0 {
		lowerLimit = approx.dTable[i]
	} else {
		lowerLimit = 0
	}

	//We can never get a match from here.
	//If lowerLimit is greater than edits left it's not possible to continue.
	if editLeft < lowerLimit {
		return
	}

	if !(L < R) {
		return
	}

	// We have a match
	if i < 0 {
		approx.Ls = append(approx.Ls, L)
		approx.Rs = append(approx.Rs, R)
		approx.matchLengths = append(approx.matchLengths, matchLength)

		// Extract the edits and reverse them.
		revEdits = append(revEdits, *edits...)

		for i, j := 0, len(revEdits)-1; i < j; i, j = i+1, j-1 {
			revEdits[i], revEdits[j] = revEdits[j], revEdits[i]
		}
		//Building cigar from edits
		approx.cigar = append(approx.cigar, editsToCigar(revEdits))
		return
	}

	//X- and =-operation
	aMatch := IndexOf(string(approx.key[i]), alphabet)

	for a := 1; a < len(alphabet); a++ {
		newL := C[a] + O[a][L]
		newR := C[a] + O[a][R]

		var editCost int
		if a == aMatch {
			editCost = 0
		} else {
			editCost = 1
		}

		if editLeft-editCost < 0 {
			continue
		}
		if newL >= newR {
			continue
		}

		if editCost == 1 {
			*edits = append(*edits, 'X')
		} else {
			*edits = append(*edits, '=')
		}
		recApproxMatching(approx, newL, newR, i-1, matchLength+1, editLeft-editCost, edits)
		*edits = (*edits)[:len(*edits)-1]
	}

	//I operation
	*edits = append(*edits, 'I')
	recApproxMatching(approx, L, R, i-1, matchLength, editLeft-1, edits)
	*edits = (*edits)[:len(*edits)-1]

	// D operations
	*edits = append(*edits, 'D')

	for a := 1; a < len(alphabet); a++ {
		newL := C[a] + O[a][L]
		newR := C[a] + O[a][R]

		if newL >= newR {
			continue
		}
		recApproxMatching(approx, newL, newR, i, matchLength+1, editLeft-1, edits)
	}
	*edits = (*edits)[:len(*edits)-1]
}

/**
Converts the edits recorded in the approximate search methods to the CIGAR format
Insertions as I
Deletions as D
Matches as = for exact matches, and X for mismatches
*/
func editsToCigar(edits []rune) string {
	var cigar string
	curr := edits[0]
	counter := 0

	for i := 0; i < len(edits); i++ {
		if edits[i] == curr {
			counter++
		} else {
			strCounter := strconv.FormatInt(int64(counter), 10)
			cigar += strCounter + string(curr)
			curr = edits[i]
			counter = 1
		}
	}
	strCounter := strconv.FormatInt(int64(counter), 10)
	cigar += strCounter + string(curr)
	return cigar
}

func generateAlphabet(info *Info) {
	var alphabet []string
	inputString := info.input

	for s := range inputString {
		found := false
		for i := range alphabet {
			if string(inputString[s]) == alphabet[i] {
				found = true
				break
			}
		}
		if !found {
			alphabet = append(alphabet, string(inputString[s]))
		}
	}
	sort.Strings(alphabet)
	info.alphabet = alphabet
}

func bwt(x string, SA []int, i int) string {
	x_index := SA[i]
	if x_index == 0 {
		return string(x[len(x)-1])
	} else {
		return string(x[x_index-1])
	}
}

func generateDTable(approx *bwtApprox, info *Info) {
	minEdit := 0
	L := 0
	R := len(info.SA)
	for i := 0; i < approx.keyLength; i++ {
		a := IndexOf(string(approx.key[i]), info.alphabet)

		L = info.cTable[a] + info.roTable[a][L]
		R = info.cTable[a] + info.roTable[a][R]

		if L >= R {
			minEdit++
			L = 0
			R = len(info.SA)
		}

		if len(info.roTable) != 0 {
			approx.dTable = append(approx.dTable, minEdit)
		}

	}
}

func generateOTable(info *Info) {
	for k := 0; k < 2; k++ {
		oTable := [][]int{}
		alphabet := info.alphabet
		sa := info.SA
		x := info.input

		if k == 1 {
			sa = info.RSA
			x = Reverse(info.input[0:len(info.input)-1]) + "$"
		}
		for range alphabet {
			oTable = append(oTable, []int{0})
		}

		for i := range sa {
			for j := range alphabet {
				if bwt(x, sa, i) == alphabet[j] {
					oTable[j] = append(oTable[j], oTable[j][i]+1)
				} else {
					oTable[j] = append(oTable[j], oTable[j][i])
				}
			}
		}
		if k == 1 {
			info.roTable = oTable
			break
		}
		info.oTable = oTable
	}
}

/**
Cumulative sum, calculates the number of occurrences of characters
in the input string smaller than each character of the alphabet.
Used in the BWT search and bucket generation
*/
func generateCTable(info *Info) {
	cTable := []int{}
	for i := range info.alphabet {
		cTable = append(cTable, 0)
		for j := range info.input {
			if info.alphabet[i] > string(info.input[j]) {
				cTable[i] += 1
			}
		}
	}
	info.cTable = cTable
}

//https://github.com/heapwolf/go-indexof/blob/master/indexof.go
func IndexOf(params ...interface{}) int {
	v := reflect.ValueOf(params[0])
	arr := reflect.ValueOf(params[1])

	var t = reflect.TypeOf(params[1]).Kind()

	if t != reflect.Slice && t != reflect.Array {
		panic("Type Error! Second argument must be an array or a slice.")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == v.Interface() {
			return i
		}
	}
	return -1
}
