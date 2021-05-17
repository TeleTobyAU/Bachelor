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

const UNDEFINED = int(^uint(0) >> 1)

func main() {
	info := new(Info)

	info.input = "mmiissiissiippii$" //"GGCTTTCCGTTGGCCATAAGGGTCTCTGGAGACGTATATCGGGTCCTAAGTGCATACGACCAATTAAAGCGACGACGCTGAGATCGCAAATAGGATAGTC$"
	generateAlphabet(info)

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
func SAIS(info *Info, x string) []int {
	start := time.Now()
	n, alphSize := str2int(x)
	SA := make([]int, len(n))
	names := make([]int, len(n))
	sumString := make([]int, len(n))
	sumOffset := make([]int, len(n))
	LSTypes := make([]bool, len(n))
	maxAlph := len(n) + 1
	if alphSize > len(n) {
		maxAlph = alphSize
	}
	buckets := make([]int, maxAlph)
	bucketEnd := make([]int, maxAlph)

	sortSA(n, &SA, &names, &sumString, &sumOffset, &buckets, &bucketEnd, &LSTypes, alphSize)
	fmt.Println("total", time.Since(start))

	return SA
}

func str2int(x string) ([]int, int) {
	alpha := map[byte]int{}
	for i := range x {
		alpha[x[i]] = 1
	}
	tempAlph := make([]byte, len(alpha))
	i := 0
	for c := range alpha {
		tempAlph[i] = c
		i++
	}

	sort.Slice(tempAlph, func(i, j int) bool {
		return tempAlph[i] < tempAlph[j]
	})

	for i, c := range tempAlph {
		alpha[c] = i
	}

	fmt.Println(alpha)
	out := make([]int, len(x))
	for i := range x {
		out[i] = alpha[x[i]]
	}
	fmt.Println(out)

	return out, len(alpha)
}

func classifyLS(n []int, LSTypes *[]bool) {
	(*LSTypes)[len(n)-1] = true

	for i := len(n) - 2; i >= 0; i-- {
		if n[i] == n[i+1] {
			(*LSTypes)[i] = (*LSTypes)[i+1]
		} else {
			if n[i] > n[i+1] {
				(*LSTypes)[i] = false
			} else {
				(*LSTypes)[i] = true
			}
		}
	}
}

func isLMSIndex(LSString []bool, i int) bool {
	if i == 0 {
		return false
	} else {
		return LSString[i] && !LSString[i-1]
	}
}

func placeLMS(n []int, alphSize int, SA *[]int, LSTypes *[]bool, buckets *[]int, bucketEnds *[]int) {
	for i := range *SA {
		(*SA)[i] = UNDEFINED
	}

	findBucketEnds(alphSize, buckets, bucketEnds)

	//SA-IS step 1, placing LMS substrings in saisStruct
	for i := 0; i < len(n); i++ {
		if isLMSIndex(*LSTypes, i) {
			(*bucketEnds)[n[i]]--
			(*SA)[(*bucketEnds)[n[i]]] = i
		}
	}

}

func induceL(n []int, alphSize int, SA *[]int, LSTypes *[]bool, buckets *[]int, bucketStarts *[]int) {
	bucketBeginnings(alphSize, buckets, bucketStarts)
	for i := 0; i < len(n); i++ {
		if (*SA)[i] == UNDEFINED {
			continue
		}

		if (*SA)[i] == 0 {
			continue
		}

		j := (*SA)[i] - 1
		if !(*LSTypes)[j] {
			(*SA)[(*bucketStarts)[n[j]]] = j
			(*bucketStarts)[n[j]]++
		}
	}
}

func induceS(n []int, alphSize int, SA *[]int, LSTypes *[]bool, buckets *[]int, bucketEnds *[]int) {
	findBucketEnds(alphSize, buckets, bucketEnds)
	for i := len(*SA); i > 0; i-- {
		if (*SA)[i-1] == 0 {
			continue
		}

		j := (*SA)[i-1] - 1
		if (*LSTypes)[j] {
			(*bucketEnds)[n[j]]--
			(*SA)[(*bucketEnds)[n[j]]] = j
		}

	}
}

func computeBuckets(n []int, buckets *[]int) {
	for i := 0; i < len(n); i++ {
		if n[i] != -1 {
			(*buckets)[n[i]]++
		}
	}
}

func findBucketEnds(alphSize int, buckets *[]int, bucketEnds *[]int) {
	(*bucketEnds)[0] = (*buckets)[0]
	for i := 1; i < alphSize; i++ {
		(*bucketEnds)[i] = (*bucketEnds)[i-1] + (*buckets)[i]
	}
}

func bucketBeginnings(alphSize int, buckets *[]int, bucketStarts *[]int) {
	(*bucketStarts)[0] = 0
	for i := 1; i < alphSize; i++ {
		(*bucketStarts)[i] = (*bucketStarts)[i-1] + (*buckets)[i-1]
	}
}

func equalLMS(n []int, LSTypes *[]bool, i int, j int) bool {
	if i == len(n) || j == len(n) {
		return false
	}
	k := 0
	for {
		iLMS := isLMSIndex(*LSTypes, i+k)
		jLMS := isLMSIndex(*LSTypes, j+k)
		if k > 0 && iLMS && jLMS {
			return true
		}
		if iLMS != jLMS || n[i+k] != n[j+k] || ((*LSTypes)[i+k]) != ((*LSTypes)[j+k]) {
			return false
		}
		k++
	}
}

func reduceSA(n []int, SA *[]int, names *[]int, LSTypes *[]bool, newAlphSize *int, sumString *[]int, sumOffset *[]int, newStrLen *int) {
	name := 0

	for i := range *names {
		(*names)[i] = UNDEFINED
	}
	(*names)[(*SA)[0]] = name

	lastS := (*SA)[0]

	for i := 1; i < len(n); i++ {
		j := (*SA)[i]
		if !isLMSIndex(*LSTypes, j) {
			continue
		}
		if !equalLMS(n, LSTypes, lastS, j) {
			name++
		}
		lastS = j
		(*names)[j] = name
	}
	*newAlphSize = name + 1

	j := 0
	for i := 0; i < len(n); i++ {
		name = (*names)[i]
		if name == UNDEFINED {
			continue
		}

		(*sumOffset)[j] = i
		(*sumString)[j] = name
		j++
	}

	var temp []int
	for i := 0; i < len(*sumString); i++ {
		if (*sumString)[i] != 0 {
			temp = append(temp, (*sumString)[i])
		}
	}
	*sumString = append(temp, 0)

	*newStrLen = j - 1
}

func recursiveSorting(n []int, SA *[]int, names *[]int, LSTypes *[]bool, buckets *[]int, bucketEnds *[]int, reducedString *[]int, reducedOffset *[]int, alphSize int) {
	classifyLS(n, LSTypes)

	computeBuckets(n, buckets)

	placeLMS(n, alphSize, SA, LSTypes, buckets, bucketEnds)
	fmt.Println("n", n)
	fmt.Println(LSTypes)
	fmt.Println(SA)

	induceL(n, alphSize, SA, LSTypes, buckets, bucketEnds)
	fmt.Println(n)
	fmt.Println(LSTypes)
	fmt.Println(SA)

	induceS(n, alphSize, SA, LSTypes, buckets, bucketEnds)
	fmt.Println(n)
	fmt.Println(LSTypes)
	fmt.Println(SA)
	fmt.Println("---------------------")

	var newAlphSize int
	var newStrLen int
	fmt.Println(newStrLen, len(n))
	reduceSA(n, SA, names, LSTypes, &newAlphSize, reducedString, reducedOffset, &newStrLen)
	fmt.Println(reducedString)

	newSa := make([]int, len(n))
	newNames := make([]int, len(n))
	newLSTypes := make([]bool, len(n))
	newSumStr := make([]int, len(n))
	newSumOffset := make([]int, len(n))
	newBuckets := make([]int, newAlphSize)
	newBucketEnds := make([]int, newAlphSize)
	sortSA(*reducedString, &newSa, &newNames, &newSumStr, &newSumOffset, &newBuckets, &newBucketEnds, &newLSTypes, newAlphSize)

	remapLMS(n, buckets, bucketEnds, alphSize, newStrLen, &newSa, reducedOffset, SA)

	induceL(n, alphSize, SA, LSTypes, buckets, bucketEnds)

	induceS(n, alphSize, SA, LSTypes, buckets, bucketEnds)
}

func sortSA(n []int, SA *[]int, names *[]int, sumString *[]int, sumOffset *[]int, buckets *[]int, bucketEnds *[]int, LSTypes *[]bool, alphSize int) {
	if len(n) == 0 {
		(*SA)[0] = 0
		return
	}

	if alphSize == len(n)+1 {
		(*SA)[0] = len(n)
		for i := 0; i < len(n)-1; i++ {
			j := n[i]
			(*SA)[j] = i
		}
	} else {
		recursiveSorting(n, SA, names, LSTypes, buckets, bucketEnds, sumString, sumOffset, alphSize)
		fmt.Println("hej")
	}
}

func remapLMS(n []int, buckets *[]int, bucketEnds *[]int, alphSize int, reducedLength int, reducedSA *[]int, reducedOffset *[]int, SA *[]int) {
	findBucketEnds(alphSize, buckets, bucketEnds)
	for i := reducedLength + 1; i > 0; i-- {
		idx := (*reducedOffset)[(*reducedSA)[i-1]]
		(*bucketEnds)[n[idx]]--
		(*SA)[(*bucketEnds)[n[idx]]] = idx
	}
	(*SA)[0] = len(n) - 1
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
	xIndex := SA[i]
	if xIndex == 0 {
		return string(x[len(x)-1])
	} else {
		return string(x[xIndex-1])
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
	counter := make([]int, len(info.alphabet))

	for i := range info.input {
		switch info.input[i] {
		case '$':
			counter[0]++
		case 'A':
			counter[1]++
		case 'C':
			counter[2]++
		case 'G':
			counter[3]++
		case 'T':
			counter[4]++
		}
	}

	cTable := make([]int, len(info.alphabet))
	for i := 0; i < len(counter); i++ {
		for j := i - 1; j >= 0; j-- {
			cTable[i] += counter[j]
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
