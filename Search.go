package main

import (
	"fmt"
	"sort"
	"strconv"
)

func exactMatch(info *Info) {
	initBwtSearch(info)
	exactMatch := indexBwtSearch(info)
	sort.Ints(exactMatch)

	fmt.Println("Exact match result.\nYellow indicate a match", len(exactMatch))
	j := 0
	for i := 0; i < len(info.input); i++ {
		if i >= exactMatch[j] && i < (exactMatch[j]+len(info.key)) {
			for j := 0; j < len(info.key); j++ {
				fmt.Print("\033[33m", string(info.input[i]))
				if len(info.key) != 1 {
					i++
				}
			}
			if j < len(exactMatch)-1 {
				j++
			}
			continue
		}
		fmt.Print("\033[0m", string(info.input[i]))
	}
	fmt.Println("\033[0m")
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
	alphabet := info.alphabet

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

		L = info.cTable[a] + info.oTable[a][L]
		R = info.cTable[a] + info.oTable[a][R]
		i -= 1
	}

	info.L = L
	info.R = R
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

/**
C table, is number of lexicographically smaller charter than alphabet i in string x.
*/
func generateCTable(info *Info) {
	counter := make([]int, len(info.alphabet))
	fmt.Println(counter)
	//TODO Vi kan ikke udregne for stringen MMIISSIISSIIPPII f.eks. Vi har hard coded det til kun at have alphabetet fra et nucleotide.
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
