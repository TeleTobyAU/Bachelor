package main

import (
	"fmt"
	"sort"
	"strconv"
)

func ExactMatch(exact *BwtExact) {
	InitBwtSearch(exact)
	exactMatch := IndexBwtSearch(exact)
	sort.Ints(exactMatch)

	j := 0
	for i := 0; i < len(exact.bwtTable.Input); i++ {
		if i >= exactMatch[j] && i < (exactMatch[j]+len(exact.Key)) {
			for j := 0; j < len(exact.Key); j++ {
				fmt.Print("\033[33m", string(exact.bwtTable.Input[i]))
				if len(exact.Key) != 1 {
					i++
				}
			}
			if j < len(exactMatch)-1 {
				j++
			}
			continue
		}
		fmt.Print("\033[0m", string(exact.bwtTable.Input[i]))
	}
	fmt.Println("\033[0m")
}

func InitBwtApproxIter(maxEdit int, info *Info, approx *BwtApprox) {
	//Init struct bwt_Approx
	approx.bwtTable = info

	//Set up edits buffer.
	keyLength := len(approx.Key)
	approx.keyLength = keyLength

	//Building D table
	generateDTable(approx, info)

	//Start searching
	L := 0
	R := len(info.SA)
	i := keyLength - 1
	edits := &approx.editBuff

	//X- and =-operation
	aMatch := IndexOf(string(approx.Key[i]), info.Alphabet)

	for a := 1; a < len(info.Alphabet); a++ {
		newL := info.CTable[a] + info.OTable[a][L]
		newR := info.CTable[a] + info.OTable[a][R]

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
			//A mismatch is described as a X
			*edits = append(*edits, 'X')
		} else {
			//A match is described as =
			*edits = append(*edits, '=')
		}
		recApproxMatching(approx, newL, newR, i-1, 1, maxEdit-editCost, edits)

		//Remove the charter that just was appended.
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

func recApproxMatching(approx *BwtApprox, L int, R int, i int, matchLength int, editLeft int, edits *[]rune) {
	//initializing variables for rec approx
	C := approx.bwtTable.CTable
	O := approx.bwtTable.OTable
	alphabet := approx.bwtTable.Alphabet
	var lowerLimit int
	var revEdits []rune

	if i >= 0 {
		lowerLimit = approx.DTable[i]
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
		approx.Cigar = append(approx.Cigar, editsToCigar(revEdits))
		return
	}

	//X- and =-operation
	aMatch := IndexOf(string(approx.Key[i]), alphabet)

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

func IndexBwtSearch(exact *BwtExact) []int {
	var match []int

	for i := 0; i < (exact.R - exact.L); i++ {
		match = append(match, exact.bwtTable.SA[exact.L+i])
	}

	return match
}

func InitBwtSearch(exact *BwtExact) {
	n := len(exact.bwtTable.SA)
	m := len(exact.Key)
	key := exact.Key
	alphabet := exact.bwtTable.Alphabet
	CTable := exact.bwtTable.CTable
	OTable := exact.bwtTable.OTable

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
		for j := range alphabet {
			if string(key[i]) == alphabet[j] {
				a = j
			}
		}

		L = CTable[a] + OTable[a][L]
		R = CTable[a] + OTable[a][R]
		i -= 1
	}

	exact.L = L
	exact.R = R
}

func Bwt(x string, SA []int, i int) string {
	xIndex := SA[i]
	if xIndex == 0 {
		return string(x[len(x)-1])
	} else {
		return string(x[xIndex-1])
	}
}

func generateDTable(approx *BwtApprox, info *Info) {
	minEdit := 0
	L := 0
	R := len(info.SA)
	for i := 0; i < approx.keyLength; i++ {
		a := IndexOf(string(approx.Key[i]), info.Alphabet)

		L = info.CTable[a] + info.roTable[a][L]
		R = info.CTable[a] + info.roTable[a][R]

		if L >= R {
			minEdit++
			L = 0
			R = len(info.SA)
		}

		if len(info.roTable) != 0 {
			approx.DTable = append(approx.DTable, minEdit)
		}

	}
}

func GenerateOTableReverse(info *Info) {
	reverseOTable := [][]int{}
	reverseAlphabet := info.Alphabet
	reverseSA := info.ReverseSA
	reverseInput := Reverse(info.Input[0:len(info.Input)-1]) + "$"

	for range reverseAlphabet {
		reverseOTable = append(reverseOTable, []int{0})
	}

	for i := range reverseSA {
		for j := range reverseAlphabet {
			if Bwt(reverseInput, reverseSA, i) == reverseAlphabet[j] {
				reverseOTable[j] = append(reverseOTable[j], reverseOTable[j][i]+1)
			} else {
				reverseOTable[j] = append(reverseOTable[j], reverseOTable[j][i])
			}
		}
	}
	info.roTable = reverseOTable
}

func GenerateOTable(info *Info) {
	oTable := [][]int{}
	alphabet := info.Alphabet
	sa := info.SA
	x := info.Input

	for range alphabet {
		oTable = append(oTable, []int{0})
	}
	for i := range sa {
		for j := range alphabet {
			if Bwt(x, sa, i) == alphabet[j] {
				oTable[j] = append(oTable[j], oTable[j][i]+1)
			} else {
				oTable[j] = append(oTable[j], oTable[j][i])
			}
		}
	}
	info.OTable = oTable
}

// GenerateCTableOptimized
//C table, is number of lexicographically smaller charter than alphabet i in string x.
///**
func GenerateCTableOptimized32(input string, alphabet []string, type32 bool) ([]int, []int32) {

	counter := make([]int, len(alphabet))
	for i := range input {
		switch input[i] {
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

	cTable := make([]int, len(alphabet))
	for i := 0; i < len(counter); i++ {
		for j := i - 1; j >= 0; j-- {
			cTable[i] += counter[j]
		}
	}

	if type32 {
		var cTableInt32 []int32

		for x := range cTable {
			cTableInt32 = append(cTableInt32, int32(cTable[x]))
		}
		return nil, cTableInt32
	}
	return cTable, nil
}

func GenerateCTableOptimized(input string, alphabet []string) []int {

	counter := make([]int, len(alphabet))
	for i := range input {
		switch input[i] {
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

	cTable := make([]int, len(alphabet))
	for i := 0; i < len(counter); i++ {
		for j := i - 1; j >= 0; j-- {
			cTable[i] += counter[j]
		}
	}
	return cTable
}

func GenerateCTable(info *Info) {
	var cTable []int
	for i := range info.Alphabet {
		cTable = append(cTable, 0)
		for j := range info.Input {
			if info.Alphabet[i] > string(info.Input[j]) {
				cTable[i] += 1
			}
		}
	}
	info.CTable = cTable
}

func GenerateCTable32(info *InfoInt32) {
	var cTable []int32
	for i := range info.Alphabet {
		cTable = append(cTable, 0)
		for j := range info.Input {
			if info.Alphabet[i] > string(info.Input[j]) {
				cTable[i] += 1
			}
		}
	}
	info.CTable = cTable
}
