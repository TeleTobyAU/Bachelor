package main

import (
	"fmt"
	"math/rand"
	"sort"
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

type bwt_Approx struct {
	bwt_table           *Info
	key                 string
	L, R, next_interval int
	Ls                  []int
	Rs                  []int
	cigar               []string
	m                   int
	edit_buff           []rune
	dTable              []int
	match_lengths       []int
}

func main() {
	info := new(Info)
	info.key = "ii"
	//generateRandomNucleotide(10000, info)//
	info.input = "mmiissiissiippii$"

	//Reverse the input string
	reverse(info)

	//Sets a thresh hold
	info.threshHold = 1

	//Create alphabet
	generateAlfabet(info)

	//Creat SA and reversed SA
	createSuffixArray(info)

	//sorting SA and reversed SA
	sortSuffixArray(info)

	//Generate C table
	var timeExact []time.Duration
	start := time.Now()
	generateCTable(info)

	//Generate O Table
	generateOTable(info)

	//Init BWT search
	fmt.Println("INPUT", info.input, info.key)
	init_BWT_search(info)
	timeExact = append(timeExact, time.Now().Sub(start))

	start = time.Now()
	//r := naiveExactSearch(info)
	timeExact = append(timeExact, time.Now().Sub(start))

	//match := index_BWT_search(info)
	//fmt.Println(match)
	approxMatch := naiveApproxSearch(info)
	fmt.Println("approx match", approxMatch)

	//finePrint(info.StringSA, r, info, timeExact, match, approxMatch)

	//Create D Table
	fmt.Println(info.dTable)

	bwt_approx := new(bwt_Approx)
	init_bwt_approx_iter(1, info, bwt_approx)

	fmt.Println("\n HAHAHA", bwt_approx.match_lengths, bwt_approx.Ls, bwt_approx.Rs)
	for i := 0; i < len(bwt_approx.bwt_table.StringSA); i++ {
		fmt.Println(i, bwt_approx.bwt_table.StringSA[i])
	}

}

func init_bwt_approx_iter(max_edit int, info *Info, approx *bwt_Approx) {
	//Init struct bwt_Approx
	approx.bwt_table = info
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
		//Lookup method
		a := indexOf(string(info.key[i]), info.alphabet)
		fmt.Println(string(info.key[i]), info.alphabet)
		fmt.Println(a)
		L = info.cTable[a] + info.roTable[a][L]
		R = info.cTable[a] + info.roTable[a][R]

		if L >= R {
			minEdit++
			L = 0
			R = len(info.SA)
		}

		approx.dTable = append(approx.dTable, minEdit)
	}
	fmt.Println("D table = ", approx.dTable)

	//Set up edits buffer.
	fmt.Println("Set up edits buffer - init")
	m = len(info.key)
	approx.m = m
	//approx.edit_buff = append(approx.edit_buff, '\000')

	//Start searching
	fmt.Println("Start searching - init")
	L = 0
	R = len(info.SA)
	i := len(info.key) - 1
	edits := approx.edit_buff

	//M-Operations
	fmt.Println("M-operation, edits =", edits)
	a_match := indexOf(string(info.key[i]), info.alphabet)

	for a := 1; a < len(info.alphabet); a++ {
		new_L := info.cTable[a] + info.oTable[a][L]
		new_R := info.cTable[a] + info.oTable[a][R]

		edit_cost := 1
		fmt.Println(a, a_match)
		if a == a_match {
			edit_cost = 0
		}
		if max_edit-edit_cost < 0 {
			continue
		}
		if new_L >= new_R {
			continue
		}

		edits = append(edits, 'M')

		rec_approx_matching(info, approx, new_L, new_R, i-1, 1, max_edit-edit_cost, edits)
	}

	// I-operation
	fmt.Println("I-operation - init")
	edits = append(edits, 'I')

	rec_approx_matching(info, approx, L, R, i-1, 0, max_edit-1, edits)

	// Make sure we start at the first interval.
	info.L = m
	info.R = 0 // TODO meaning
	approx.next_interval = 0

}

func rec_approx_matching(info *Info, approx *bwt_Approx, L int, R int, i int, match_length int, leftEdit int, edit []rune) {
	//TODO struct
	approx.bwt_table = info
	lowerLimit := 0
	if i >= 0 {
		lowerLimit = approx.dTable[i]
	}

	if leftEdit < lowerLimit {
		return // We can never get a match from here.
	}
	if i < 0 { // We have a match
		approx.Ls = append(approx.Ls, L)
		approx.Rs = append(approx.Rs, R)
		approx.match_lengths = append(approx.match_lengths, match_length)

		// Extract the edits and reverse them.
		// := make([]rune, len(approx.edit_buff))
		rev_edits := []rune{}
		rev_edits = append(rev_edits, edit...)

		for i, j := 0, len(rev_edits)-1; i < j; i, j = i+1, j-1 {
			rev_edits[i], rev_edits[j] = rev_edits[j], rev_edits[i] //TODO
		}

		//Building cigar from edits
		fmt.Println("edits = ", string(edit), "reverse edit = ", string(rev_edits)) //TODO edits_to_cigar
		//cigar := new(rune)
		//edits_to_cigar(cigar, rev_edits)
		fmt.Println("PSFPSDPFSDPFOSFDODSP", edit)
		return
	}

	//M-operation
	fmt.Println("M-operation - rec")
	a_match := indexOf(string(info.key[i]), info.alphabet)

	for a := 1; a < len(info.alphabet); a++ {
		new_L := info.cTable[a] + info.oTable[a][L]
		new_R := info.cTable[a] + info.oTable[a][R]

		edit_cost := 1
		if a == a_match {
			edit_cost = 0
		}

		if leftEdit-edit_cost < 0 {
			continue
		}
		if new_L >= new_R {
			continue
		}

		edit = append(edit, 'M')

		rec_approx_matching(info, approx, new_L, new_R, i-1, match_length+1, leftEdit-edit_cost, edit)
	}
	//I operation
	fmt.Println("I-operation - rec")
	edit = append(edit, 'I')
	rec_approx_matching(info, approx, L, R, i-1, match_length, leftEdit-1, edit)

	// D operations
	fmt.Println("D-operation - rec")
	edit = append(edit, 'D')

	for a := 1; a < len(info.alphabet); a++ {
		new_L := info.cTable[a] + info.oTable[a][L]
		new_R := info.cTable[a] + info.oTable[a][R]

		if new_L >= new_R {
			continue
		}
		rec_approx_matching(info, approx, new_L, new_R, i, match_length+1, leftEdit-1, edit)
	}
	fmt.Println("Edits", string(edit))
}

func edits_to_cigar(cigar *rune, edits []rune) {
	for i := 0; i < len(edits); i++ {
		next := scan(edits)
		println("huuuhuuuhuu", string(next))
	}

}

func scan(edits []rune) []rune {
	p := edits
	for i := 0; p[i] == edits[i]; i++ {
		println("HEHEHEHEHEH")
		p[i] = p[i]
		if i < len(edits) {
			println("hohohoho")
			break
		}
	}
	return p
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

func index_BWT_search(info *Info) []int {
	match := []int{}

	for i := 0; i < (info.R - info.L); i++ {
		match = append(match, info.SA[info.L+i])
	}

	return match
}

func init_BWT_search(info *Info) {
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

func generateAlfabet(info *Info) {
	alphabet := []string{}
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
		o_Table := [][]int{}
		alphabet := info.alphabet
		sa := info.SA
		x := info.input
		if k == 1 {
			sa = info.reverseSA
			x = info.reverseInput
		}

		for range alphabet {
			o_Table = append(o_Table, []int{0})
		}

		for i := range sa {
			for j := range alphabet {
				if bwt(x, sa, i) == alphabet[j] {
					o_Table[j] = append(o_Table[j], o_Table[j][i]+1)
				} else {
					o_Table[j] = append(o_Table[j], o_Table[j][i])
				}
			}
		}
		if k == 1 {
			info.roTable = o_Table
			break
		}
		info.oTable = o_Table
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
		var index_SA = []int{}
		var oldArray = make([]string, len(SA))
		copy(oldArray, SA)

		sort.Strings(SA)
		for s := range SA {
			index_SA = append(index_SA, indexOf(SA[s], oldArray))
		}
		if i == 1 {
			info.reverseSA = index_SA
			break
		}
		info.SA = index_SA
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

func indexOf(lookingFor string, lookingIn []string) int {
	for i := range lookingIn {
		if lookingIn[i] == lookingFor {
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
