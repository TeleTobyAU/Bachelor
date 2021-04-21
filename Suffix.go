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
	m                  int
	editBuff           []rune
	dTable             []int
	matchLengths       []int
}

func main() {
	info := new(Info)
	info.key = "iss"
	//generateRandomNucleotide(10000, info)//
	info.input = "mmiissiissiippii$"

	//Reverse the input string
	reverse(info)

	//Sets a thresh hold
	info.threshHold = 1

	//Create alphabet
	generateAlphabet(info)

	//Generate C table
	generateCTable(info)

	info.SA = SAIS(info, info.input)
	info.reverseSA = SAIS(info, Reverse(info.input[0:len(info.input)-1])+"$") //Making sure the sentinel remains at the end after versing

	//Generate O Table
	generateOTable(info)

	//Init BWT search
	bwtApprox := new(bwtApprox)
	initBwtApproxIter(info.threshHold, info, bwtApprox)
	fmt.Println(info.StringSA)

	for i := 0; i < len(bwtApprox.Ls); i++ {
		fmt.Println("From index", bwtApprox.Ls[i], "to", bwtApprox.Rs[i], "in SA")
		for j := bwtApprox.Ls[i]; j < bwtApprox.Rs[i]; j++ {
			fmt.Println(info.SA[j], bwtApprox.cigar[i], info.input[j:])
		}
		fmt.Println()
	}
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
	fmt.Println(buckets)

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
	fmt.Println("SA after SAIS step 1:", SA)

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
	fmt.Println("SA after SAIS step 2:", SA)

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
	fmt.Println("SA after SAIS step 3:", SA, "\n")

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
	approx.Ls = []int{}
	approx.Rs = []int{}
	approx.cigar = []string{}

	//Building D table
	m := len(info.key)
	minEdit := 0
	L := 0
	R := len(info.SA)
	for i := 0; i < m; i++ {
		a := IndexOf(string(info.key[i]), info.alphabet)
		L = info.cTable[a] + info.roTable[a][L]
		R = info.cTable[a] + info.roTable[a][R]

		if L >= R {
			minEdit++
			L = 0
			R = len(info.SA)
		}

		approx.dTable = append(approx.dTable, minEdit)
	}

	//Set up edits buffer.
	m = len(info.key)
	approx.m = m
	//approx.edit_buff = append(approx.edit_buff, '\000')

	//Start searching
	L = 0
	R = len(info.SA)
	i := m - 1
	edits := approx.editBuff //TODO maybe pointer

	//M-Operations
	aMatch := IndexOf(string(info.key[i]), info.alphabet) //TODO look at this later

	for a := 1; a < len(info.alphabet); a++ {
		newL := info.cTable[a] + info.oTable[a][L]
		newR := info.cTable[a] + info.oTable[a][R]

		editCost := 1
		if a == aMatch {
			editCost = 0
		}
		if maxEdit-editCost < 0 {
			continue
		}
		if newL >= newR {
			continue
		}

		edits = append(edits, 'M')
		recApproxMatching(info, approx, newL, newR, i-1, 1, maxEdit-editCost, edits)
	}

	// I-operation
	edits = append(edits, 'I')
	recApproxMatching(info, approx, L, R, i-1, 0, maxEdit-1, edits)

	// Make sure we start at the first interval.
	info.L = m
	info.R = 0 // TODO meaning
	approx.nextInterval = 0
}

func recApproxMatching(info *Info, approx *bwtApprox, L int, R int, i int, matchLength int, leftEdit int, edits []rune) {
	//TODO struct
	approx.bwtTable = info
	lowerLimit := 0
	if i >= 0 {
		lowerLimit = approx.dTable[i]
	}

	if leftEdit < lowerLimit {
		return // We can never get a match from here.
	}

	if L >= R {
		return
	}

	if i < 0 { // We have a match
		approx.Ls = append(approx.Ls, L)
		approx.Rs = append(approx.Rs, R)
		approx.matchLengths = append(approx.matchLengths, matchLength)

		// Extract the edits and reverse them.
		var revEdits []rune
		revEdits = append(revEdits, edits...)

		for i, j := 0, len(revEdits)-1; i < j; i, j = i+1, j-1 {
			revEdits[i], revEdits[j] = revEdits[j], revEdits[i] //TODO
		}

		//Building cigar from edits
		approx.cigar = append(approx.cigar, editsToCigar(revEdits))
		return
	}

	//M-operation
	aMatch := IndexOf(string(info.key[i]), info.alphabet)

	for a := 1; a < len(info.alphabet); a++ {

		newL := info.cTable[a] + info.oTable[a][L]
		newR := info.cTable[a] + info.oTable[a][R]

		editCost := 1
		if a == aMatch {
			editCost = 0
		}

		if leftEdit-editCost < 0 {
			continue
		}
		if newL >= newR {
			continue
		}

		edits = append(edits, 'M')

		recApproxMatching(info, approx, newL, newR, i-1, matchLength+1, leftEdit-editCost, edits)
	}

	//I operation
	edits = append(edits, 'I')
	recApproxMatching(info, approx, L, R, i-1, matchLength, leftEdit-1, edits)

	// D operations
	edits = append(edits, 'D')

	for a := 1; a < len(info.alphabet); a++ {
		newL := info.cTable[a] + info.oTable[a][L]
		newR := info.cTable[a] + info.oTable[a][R]

		if newL >= newR {
			continue
		}
		recApproxMatching(info, approx, newL, newR, i, matchLength+1, leftEdit-1, edits)
	}
}

func editsToCigar(edits []rune) string {
	curr := '\000'
	counter := 0
	cigar := ""
	for i := 0; i < len(edits); i++ {
		if edits[i] != curr {
			strCounter := strconv.FormatInt(int64(counter), 10)
			cigar += string(curr) + strCounter
			curr = edits[i]
			counter = 1
		} else {
			counter += 1
		}
	}
	strCounter := strconv.FormatInt(int64(counter), 10)
	cigar += string(curr) + strCounter
	return cigar[2:]
}

func reverse(info *Info) {
	chars := []rune(info.input)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	info.reverseInput = string(chars)
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

func generateOTable(info *Info) {
	for k := 0; k < 2; k++ {
		oTable := [][]int{}
		alphabet := info.alphabet
		sa := info.SA
		x := info.input
		if k == 1 {
			sa = info.reverseSA
			x = info.reverseInput
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

func generateCTable(info *Info) {
	alf := info.alphabet
	n := info.input

	sort.Strings(alf)
	cTable := []int{}
	for i := range alf {
		cTable = append(cTable, 0)
		for j := range n {
			if alf[i] > string(n[j]) {
				cTable[i] += 1
			}

		}
	}

	info.cTable = cTable
}

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

func generateRandomNucleotide(size int, info *Info) {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("ATCG")

	nucleotide := make([]rune, size)

	for i := range nucleotide {
		nucleotide[i] = letters[rand.Intn(len(letters))]
	}
	info.input = string(nucleotide) + "$"
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

func finePrint(SA []string, r int, info *Info, exact []time.Duration, match []int, approxMatch []int) {
	//Input string
	fmt.Println("\nInput String:")
	fmt.Println(info.input)
	fmt.Println()

	//Alphabet
	fmt.Println("\nAlphabet over input string:")
	fmt.Println(info.alphabet)
	fmt.Println()

	//Print sorted array in Strings
	fmt.Println("\nSuffix Array with sort in strings:")
	for i := range SA {
		fmt.Println(i, SA[i])
	}
	fmt.Println()

	// Print sorted array in integers
	fmt.Println("\nSuffix Array with sort:")
	fmt.Println(info.SA)
	fmt.Println()

	//C Table print
	fmt.Println("C Table:")
	fmt.Println(info.alphabet)
	fmt.Println(info.cTable)
	fmt.Println()

	//O Table Print
	fmt.Println("Otable:")
	printbwt := "     "
	for i := range SA {
		printbwt += bwt(info.input, info.SA, i) + " "
	}
	fmt.Println(printbwt)
	for i := range info.oTable {
		fmt.Println(info.alphabet[i], info.oTable[i])
	}
	fmt.Println()

	//Complexity
	fmt.Println("Time taken for exact match:")
	fmt.Println("Naive match: ", exact[1])
	fmt.Println("BWT search match: ", exact[0])
	fmt.Println()

	//match
	if r == (info.R - info.L) {
		fmt.Println("Matches found: ", r)
	}
	fmt.Println()

	//Index for matches in string
	fmt.Println("Index for matches in string")
	sort.Ints(match)
	for i := range match {
		fmt.Println("Match number", i+1, "is at index", match[i])
	}
	fmt.Println()

	//Naive Approx search
	fmt.Println("Index for matches for approx")
	fmt.Println(approxMatch)
	fmt.Println()

}

//No longer in use

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
