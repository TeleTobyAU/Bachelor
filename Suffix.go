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
	input           string
	reverseInput    string
	alphabet        []string
	threshHold      int
	key             string
	SA              []int
	StringSA        []string
	reverseSA       []int
	stringReverseSA []string
	cTable          []int
	oTable          [][]int
	roTable         [][]int
	dTable          []int
	L               int
	Ls              []int
	R               int
	Rs              []int
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

func main() {
	info := new(Info)
	info.key = "iis"
	//generateRandomNucleotide(500, info)//
	info.input = "mmiissiissiippii$"

	//Sets a thresh hold
	info.threshHold = 1

	//Create alphabet
	generateAlphabet(info)

	//Generate C table
	generateCTable(info)

	//Generating SAIS
	info.SA = SAIS(info, info.input)

	//Reverse the SA string and input string
	//fmt.Println("Reversed string", Reverse(info.input))
	info.reverseInput = Reverse(info.input[0:len(info.input)-1]) + "$"
	info.reverseSA = SAIS(info, Reverse(info.input[0:len(info.input)-1])+"$") //Making sure the sentinel remains at the end after versing

	//Generate O Table
	generateOTable(info)

	//Init BWT search
	bwtApprox := new(bwtApprox)
	initBwtApproxIter(info.threshHold, info, bwtApprox)

	//Print Cigar
	fmt.Println("CIGAR -", bwtApprox.cigar)
}

func SAIS(info *Info, n string) []int {
	//Classify L and S types
	LSTypes := "S"
	reversedN := Reverse(n)
	for i := 1; i < len(n); i++ {
		if reversedN[i-1] == reversedN[i] {
			LSTypes += string(LSTypes[i-1])
		} else {
			if reversedN[i-1] < reversedN[i] {
				LSTypes += "L"
			} else {
				LSTypes += "S"
			}

		}
	}
	LSTypes = Reverse(LSTypes)

	//Find LMS Indices
	LMSIndices := []int{}
	if LSTypes[0] == 'S' {
		LMSIndices = append(LMSIndices, 0)
	}

	for i := 1; i < len(LSTypes); i++ {
		if LSTypes[i] == 'S' && LSTypes[i-1] != 'S' {
			LMSIndices = append(LMSIndices, i)
		}
	}

	buckets := getBuckets(info, n)

	//Initializing SA to -1 at all positions
	SA := []int{}
	for i := 0; i < len(n); i++ {
		SA = append(SA, -1)
	}

	//SA-IS step 1, placing LMS substrings in SA
	for i := 0; i < len(n); i++ {
		if IndexOf(i, LMSIndices) != -1 {
			remappedi := IndexOf(string(n[i]), info.alphabet)
			SA[buckets[remappedi][1]] = i
			buckets[remappedi][1] -= 1
		}
	}

	//SA-IS step 2, placing L types in SA
	for i := 0; i < len(n)+1; i++ {
		if i >= len(n) {
			break
		}

		if SA[i] == -1 {
			continue
		}

		j := SA[i] - 1
		if j >= 0 {
			if LSTypes[j] == 'L' {
				remappedi := IndexOf(string(n[j]), info.alphabet)
				SA[buckets[remappedi][0]] = j
				buckets[remappedi][0] += 1
			}
		}
	}

	//SA-IS step 3, placing and sorting S types
	buckets = getBuckets(info, n) //We have to reset the buckets since our ends where modified with the insertion of LMS indices earlier
	for i := len(n); i > 0; i-- {
		if SA[i-1] == 0 {
			continue
		}

		j := SA[i-1] - 1
		if j < 0 {
			j = len(n) + j
		}
		if LSTypes[j] == 'S' { //Something is fucky when doing this with reversed n
			remappedi := IndexOf(string(n[j]), info.alphabet)
			SA[buckets[remappedi][1]] = j
			buckets[remappedi][1] -= 1
		}
	}

	return SA
}

func getBuckets(info *Info, n string) [][]int {
	//Find bucket beginnings, this will be the same as the C table, so I am reusing it
	beginnings := info.cTable

	//Find bucket ends
	ends := []int{}
	for i := 0; i < len(info.alphabet); i++ {
		ends = append(ends, -1) //Kinda wish this was 0, but get an error when it is, something must be wrong somewhere TODO
	}
	for i := len(n) - 1; i > -1; i-- {
		j := IndexOf(string(n[i]), info.alphabet)
		for k := j; k < len(info.alphabet); k++ {
			ends[k] = ends[k] + 1
		}
	}

	//Build buckets from beginnings and ends
	buckets := [][]int{}
	for i := 0; i < len(beginnings); i++ {
		buckets = append(buckets, []int{beginnings[i], ends[i]})
	}
	return buckets
}

//https://stackoverflow.com/questions/1752414/how-to-reverse-a-string-in-go
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func initBwtApproxIter(maxEdit int, info *Info, approx *bwtApprox) {
	//Init struct bwt_Approx
	approx.bwtTable = info
	approx.key = info.key

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

	//X- and =-operation
	aMatch := IndexOf(string(info.key[i]), info.alphabet)

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
	info.L = keyLength
	info.R = 0
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

func naiveApproxSearch(info *Info) []int {
	match := []int{}
	for i := 0; i < len(info.input)-len(info.key); i++ {
		hammingDistance := 0
		for j := i; j < i+len(info.key); j++ {
			if info.input[j] != info.key[j-i] {
				hammingDistance += 1
				if hammingDistance > info.threshHold {
					break
				}
			}
			if j == (i + len(info.key) - 1) {
				match = append(match, i)
			}
		}
	}
	return match
}

func indexBwtSearch(info *Info) []int {
	match := []int{}

	for i := 0; i < (info.R - info.L); i++ {
		match = append(match, info.SA[info.L+i])
	}

	return match
}

func initBwtSearch(info *Info) {
	n := len(info.SA)
	m := len(info.key)
	key := info.key
	alph := info.alphabet

	L := 0
	R := n

	if m > n {
		R = 0
		L = 1
	}
	i := m - 1
	for i >= 0 && L < R {

		//Find Index of key[i] in O table
		var a int
		for j := range alph {
			if string(key[i]) == alph[j] {
				a = j
			}
		}

		L = info.cTable[a] + info.oTable[a][L]
		R = info.cTable[a] + info.oTable[a][R]
		i -= 1
	}

	info.L = L
	info.R = R
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
			sa = info.reverseSA
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

// C table, is number of lexicographically smaller charter than alphabet i in string x.
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

func generateRandomNucleotide(size int, info *Info) {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("ATCG")

	nucleotide := make([]rune, size)

	for i := range nucleotide {
		nucleotide[i] = letters[rand.Intn(len(letters))]
	}
	info.input = string(nucleotide) + "$"
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

//No longer in use
func sortSuffixArray(info *Info) {
	for i := 0; i < 2; i++ {
		SA := info.StringSA
		if i == 1 {
			SA = info.stringReverseSA
		}
		var indexSa = []int{}
		var oldArray = make([]string, len(SA))
		copy(oldArray, SA)

		sort.Strings(SA)
		for s := range SA {
			indexSa = append(indexSa, IndexOf(SA[s], oldArray))
		}
		if i == 1 {
			info.reverseSA = indexSa
			break
		}
		info.SA = indexSa
	}

}

func createSuffixArray(info *Info) {
	for j := 0; j < 2; j++ {
		input := info.input
		if j == 1 {
			input = info.reverseInput
		}
		length := len(input)
		suffixArray := []string{}
		suffix := ""

		for i := 0; i < length; i++ {

			if i != 0 {
				suffix = suffix + string(input[i-1])
			}

			slicePiece := input[i:length] + suffix

			suffixArray = append(suffixArray, slicePiece)

		}
		if j == 1 {
			info.stringReverseSA = suffixArray
			break
		}
		info.StringSA = suffixArray
	}
}

func findBWT(array []string) []string {
	length := len(array)
	bwt := []string{}
	for _, s := range array {
		bwt = append(bwt, string(s[length-1]))
	}

	return bwt
}

func naiveExactSearch(info *Info) int {
	counter := 0
	indices := []int{}
	k := info.key
	n := info.input

	for i := range n {
		if n[i] == k[0] {
			for j := range k {
				if k[j] == n[i+j] && len(k)+i < len(n) {
					if j+1 == len(k) {
						counter += 1
						indices = append(indices, i)
					}
				} else {
					break
				}

			}
		}

	}
	return counter
}
