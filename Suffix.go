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
	R               int
}

func main() {
	info := new(Info)
	info.key = "iss"
	info.input = "wfsdfsdf"
	info.input = "mmiissiissiippii" //generateRandomNucleotide(5000, &info)

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
	r := naiveExactSearch(info)
	timeExact = append(timeExact, time.Now().Sub(start))

	match := index_BWT_search(info)
	fmt.Println(match)
	approxMatch := naiveApproxSearch(info)
	fmt.Println(approxMatch)

	finePrint(info.StringSA, r, info, timeExact, match, approxMatch)

	//Create D Table
	generateDTable(info)
	fmt.Println(info.dTable)
	/*
		approxMatchOPT := init_bwt_approx_iter(inputString, key, alphabet, reverseInput, rsa, cTable, oTable, ROTable, dTable, 10)
	*/
}

/*
func init_bwt_approx_iter(inputString string, key string, alphabet []string, reverseInput string, rsa []string, ctable []int, rotable [][]int, otable [][]int, dtable []int, max_edit int) []int {
	approxMatch := []int{}
	L := 0
	R := len(rsa)
	i := len(key) - 1

	edits := 'r'
	fmt.Println(edits)
	a_match := key[i]

	for j := 1; j < len(alphabet); j++ {
		new_L := ctable[j] + otable[j][L]
		new_R := ctable[j] + otable[j][R]

		edit_cost := 1
		if j == int(a_match) {
			edit_cost = 0
		}

		if max_edit-edit_cost > 0 {
			break
		}
		if new_L < new_R {
			break
		}

		edits = 'M'
		rec_approx_matching(new_L, new_R, i-1, 1, max_edit-edit_cost, edits+1)
	}

	edits = 'I'
	rec_approx_matching(L, R, i-1, 0, max_edit-1, edits+1)

	L = len(key)
	R = 0
	//TODO next interval
	return approxMatch
}

func rec_approx_matching(L int, R int, i int, i2 int, i3 int, i4 int32) {
	lowerLimit := 0
	if i >= 0 {
		lowerLimit =
	}
}
*/
func generateDTable(info *Info) {
	m := len(info.key)
	DTable := []int{}
	minEdit := 0
	L := 0
	R := len(info.reverseSA)
	for i := 0; i < m; i++ {
		var a int
		for j := range info.alphabet {
			if string(info.key[i]) == info.alphabet[j] {
				a = j
			}
		}
		L = info.cTable[a] + info.roTable[a][L]
		R = info.cTable[a] + info.roTable[a][R]

		if L >= R {
			minEdit += 1
			L = 0
			R = len(info.reverseSA)
		}

		DTable = append(DTable, minEdit)
	}
	info.dTable = DTable
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
	for i := range lookingFor {
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
