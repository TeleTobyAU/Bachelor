package main

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

type Info struct {
	input        string
	reverseInput string
	alphabet     []string
	threshHold   int
	SA           []int
	RSA          []int
	cTable       []int
	oTable       [][]int
	roTable      [][]int
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
	key := "na"
	//generateRandomNucleotide(10000, info)//
	info.input = "banana$"

	//Allowed mismatch degree for searching
	info.threshHold = 1

	info.alphabet = generateAlphabet(info.input)

	generateCTable(info)

	info.SA = SAIS(info, info.input)

	info.reverseInput = Reverse(info.input[0:len(info.input)-1]) + "$"
	info.RSA = SAIS(info, Reverse(info.input[0:len(info.input)-1])+"$") //Making sure the sentinel remains at the end after versing

	generateOTable(info)

	bwtApprox := new(bwtApprox)
	initBwtApproxIter(key, info.threshHold, info, bwtApprox)
	fmt.Println("D Table- ", bwtApprox.dTable)
	for i := 0; i < len(info.SA); i++ {
		fmt.Println(info.reverseInput[info.RSA[i]:])
	}

	for i := 0; i < len(bwtApprox.Ls); i++ {
		fmt.Println("From index", bwtApprox.Ls[i], "to", bwtApprox.Rs[i], "in SA")
		for j := bwtApprox.Ls[i]; j < bwtApprox.Rs[i]; j++ {
			fmt.Println(info.SA[j], bwtApprox.cigar[i], info.input[info.SA[j]:])
		}
		fmt.Println()
	}
	fmt.Println("Input -", info.input)
	fmt.Println("Key -", key)
	fmt.Println("CIGAR -", bwtApprox.cigar)
	fmt.Println("Ls -", bwtApprox.Ls)
	fmt.Println("Rs -", bwtApprox.Rs)
	fmt.Println("Match length -", bwtApprox.matchLengths)
}

/**
Linear suffix array construction by almost pure induced sorting.
Algorithm derived from Zhang and Chan 2009
*/
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
		if LSTypes[j] == 'S' {
			remappedi := IndexOf(string(n[j]), info.alphabet)
			SA[buckets[remappedi][1]] = j
			buckets[remappedi][1] -= 1
		}
	}

	return SA
}

/**
Helper method used in the SA-IS algorithm
Calculates the intervals in the suffix array for each letter
*/
func getBuckets(info *Info, n string) [][]int {
	//Find bucket beginnings, this will be the same as the C table, so I am reusing it
	beginnings := info.cTable

	//Find bucket ends
	ends := []int{}
	for i := 0; i < len(info.alphabet); i++ {
		ends = append(ends, -1)
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

func initBwtApproxIter(key string, maxEdit int, info *Info, approx *bwtApprox) {
	//Init struct bwt_Approx
	approx.bwtTable = info
	approx.key = key

	//Set up edits buffer.
	keyLength := len(approx.key)
	approx.keyLength = keyLength

	generateDTable(approx, info)

	//Start searching
	L := 0
	R := len(info.SA)
	i := keyLength - 1
	edits := &approx.editBuff

	//M-Operations
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
		recApproxMatching(info, approx, newL, newR, i-1, 1, maxEdit-editCost, edits)
		*edits = (*edits)[:len(*edits)-1]

	}

	// I-operation
	*edits = append(*edits, 'I')
	recApproxMatching(info, approx, L, R, i-1, 0, maxEdit-1, edits)
	*edits = (*edits)[:len(*edits)-1]

	approx.L = keyLength
	approx.R = 0
	approx.nextInterval = 0
}

func recApproxMatching(info *Info, approx *bwtApprox, L int, R int, i int, matchLength int, editLeft int, edits *[]rune) {
	approx.bwtTable = info
	var lowerLimit int
	if i >= 0 {
		lowerLimit = approx.dTable[i]
	} else {
		lowerLimit = 0
	}
	// We can never get a match from here.
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
		var revEdits []rune
		revEdits = append(revEdits, *edits...)
		for i, j := 0, len(revEdits)-1; i < j; i, j = i+1, j-1 {
			revEdits[i], revEdits[j] = revEdits[j], revEdits[i]
		}

		//Building cigar from edits
		approx.cigar = append(approx.cigar, editsToCigar(revEdits))
		return
	}

	//M-operation
	aMatch := IndexOf(string(approx.key[i]), info.alphabet)

	for a := 1; a < len(info.alphabet); a++ {

		newL := info.cTable[a] + info.oTable[a][L]
		newR := info.cTable[a] + info.oTable[a][R]

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

		recApproxMatching(info, approx, newL, newR, i-1, matchLength+1, editLeft-editCost, edits)
		*edits = (*edits)[:len(*edits)-1]
	}

	//I operation
	*edits = append(*edits, 'I')
	recApproxMatching(info, approx, L, R, i-1, matchLength, editLeft-1, edits)
	*edits = (*edits)[:len(*edits)-1]

	// D operations
	*edits = append(*edits, 'D')

	for a := 1; a < len(info.alphabet); a++ {
		newL := info.cTable[a] + info.oTable[a][L]
		newR := info.cTable[a] + info.oTable[a][R]

		if newL >= newR {
			continue
		}
		recApproxMatching(info, approx, newL, newR, i, matchLength+1, editLeft-1, edits)
	}
	*edits = (*edits)[:len(*edits)-1]
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

/**

 */
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

/**
Converts the edits recorded in the approximate search methods to the CIGAR format
Insertions as I
Deletions as D
Matches as = for exact matches, and X for mismatches
*/
func editsToCigar(edits []rune) string {
	var curr rune
	var counter int
	var cigar string
	for i := 0; i < len(edits); i++ {
		if edits[i] != curr {
			strCounter := strconv.FormatInt(int64(counter), 10)
			cigar += strCounter + string(curr)
			curr = edits[i]
			counter = 1
		} else {
			counter += 1
		}
	}
	strCounter := strconv.FormatInt(int64(counter), 10)
	cigar += strCounter + string(curr)
	return cigar[1:]
}

/**
Looks up the character in the BWT at index i
*/
func bwt(x string, SA []int, i int) string {
	xIndex := SA[i]
	if xIndex == 0 {
		return string(x[len(x)-1])
	} else {
		return string(x[xIndex-1])
	}
}

/**
Generates an alphabet for a string
*/
func generateAlphabet(input string) []string {
	var alphabet []string

	for s := range input {
		found := false
		for i := range alphabet {
			if string(input[s]) == alphabet[i] {
				found = true
				break
			}
		}
		if !found {
			alphabet = append(alphabet, string(input[s]))
		}
	}
	sort.Strings(alphabet)
	return alphabet
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

//https://stackoverflow.com/questions/1752414/how-to-reverse-a-string-in-go
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
