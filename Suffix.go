package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main() {
	key := "AT"
	inputString := generateRandomNucleotide(5000) + "$"
	threshHold := 1

	//Create alfabet
	alphabet := generateAlfabet(inputString)

	//Creat the SA
	var SA []string
	SA = createSuffixArray(inputString)

	//sort the SA
	sortedSA := sortSuffixArray(SA)

	//Generate C table
	timeExact := []time.Duration{}
	start := time.Now()
	cTable := generateCTable(inputString, alphabet)

	//Generate O Table
	oTable := generateOTable(inputString, alphabet, sortedSA)

	//Init BWT search

	L, R := init_BWT_search(key, alphabet, sortedSA, cTable, oTable)
	timeExact = append(timeExact, time.Now().Sub(start))

	start = time.Now()
	r := naiveExactSearch(inputString, key)
	timeExact = append(timeExact, time.Now().Sub(start))

	match := index_BWT_search(sortedSA, L, R)
	fmt.Println(match)
	approxMatch := naiveApproxSearch(inputString, key, threshHold)
	fmt.Println(approxMatch)

	finePrint(inputString, alphabet, SA, sortedSA, cTable, oTable, r, L, R, timeExact, match, approxMatch)

	//Create Reverse O Table
	reverseInput := reverse(inputString)
	rsa := createSuffixArray(reverseInput)
	irsa := sortSuffixArray(rsa)
	ROTable := generateOTable(reverseInput, alphabet, irsa)

	dTable := generateDTable(inputString, alphabet, key, irsa, cTable, ROTable)

}

func generateDTable(inputString string, alphabet []string, key string, irsa []int, cTable []int, roTable [][]int) []int {
	m := len(key)
	DTable := []int{}
	minEdit := 0
	L := 0
	R := len(irsa)
	for i := 0; i < m; i++ {
		var a int
		for j := range alphabet {
			if string(key[i]) == alphabet[j] {
				a = j
			}
		}
		L = cTable[a] + roTable[a][L]
		R = cTable[a] + roTable[a][R]

		if L >= R {
			minEdit += 1
			L = 0
			R = len(irsa)
		}

		DTable = append(DTable, minEdit)
	}
	return DTable
}

func reverse(inputString string) string {
	chars := []rune(inputString)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

func naiveApproxSearch(inputString string, key string, hold int) []int {
	match := []int{}
	for i := 0; i < len(inputString)-len(key); i++ {
		hammingDistance := 0
		for j := i; j < i+len(key); j++ {
			if inputString[j] != key[j-i] {
				hammingDistance += 1
				if hammingDistance > hold {
					break
				}
			}
			if j == (i + len(key) - 1) {
				match = append(match, i)
			}
		}
	}
	return match
}

func index_BWT_search(SA []int, l int, r int) []int {
	match := []int{}

	for i := 0; i < (r - l); i++ {
		match = append(match, SA[l+i])
	}

	return match
}

func init_BWT_search(key string, alf []string, sa []int, cTable []int, oTable [][]int) (int, int) {
	n := len(sa)
	m := len(key)

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
		for j := range alf {
			if string(key[i]) == alf[j] {
				a = j
			}
		}

		L = cTable[a] + oTable[a][L]
		R = cTable[a] + oTable[a][R]
		i -= 1
	}

	return L, R
}

func generateAlfabet(inputString string) []string {
	alfabet := []string{}

	for s := range inputString {
		found := false
		for i := range alfabet {
			if string(inputString[s]) == alfabet[i] {
				found = true
				break
			}
		}
		if !found {
			alfabet = append(alfabet, string(inputString[s]))
		}
	}
	sort.Strings(alfabet)
	return alfabet
}

func bwt(x string, SA []int, i int) string {
	x_index := SA[i]
	if x_index == 0 {
		return string(x[len(x)-1])
	} else {
		return string(x[x_index-1])
	}
}

func generateOTable(x string, alfabet []string, sa []int) [][]int {
	o_Table := [][]int{}

	for range alfabet {
		o_Table = append(o_Table, []int{0})
	}

	for i := range sa {
		for j := range alfabet {
			if bwt(x, sa, i) == alfabet[j] {
				o_Table[j] = append(o_Table[j], o_Table[j][i]+1)
			} else {
				o_Table[j] = append(o_Table[j], o_Table[j][i])
			}
		}
	}
	return o_Table
}

func generateCTable(n string, alf []string) []int {
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

	return cTable
}

func sortSuffixArray(array []string) []int {

	var index_SA = []int{}
	var oldArray = make([]string, len(array))
	copy(oldArray, array)

	sort.Strings(array)
	for s := range array {
		index_SA = append(index_SA, indexOf(array[s], oldArray))
	}
	return index_SA
}

func generateRandomNucleotide(size int) string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("ATCG")

	nucleotide := make([]rune, size)

	for i := range nucleotide {
		nucleotide[i] = letters[rand.Intn(len(letters))]
	}
	return string(nucleotide)
}

func createSuffixArray(input string) []string {
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

	return suffixArray
}

func indexOf(lookingFor string, lookingIn []string) int {
	for i := range lookingFor {
		if lookingIn[i] == lookingFor {
			return i
		}
	}
	return -1
}

func finePrint(input string, alphabet []string, sa []string, sortedSA []int, ctable []int, otable [][]int, r int, l int, R int, exact []time.Duration, match []int, approxMatch []int) {
	//Input string
	fmt.Println("\nInput String:")
	fmt.Println(input)
	fmt.Println()

	//Alphabet
	fmt.Println("\nAlphabet over input string:")
	fmt.Println(alphabet)
	fmt.Println()

	//Print sorted array in Strings
	fmt.Println("\nSuffix Array with sort in strings:")
	for i := range sa {
		fmt.Println(i, sa[i])
	}
	fmt.Println()

	// Print sorted array in integers
	fmt.Println("\nSuffix Array with sort:")
	fmt.Println(sortedSA)
	fmt.Println()

	//C Table print
	fmt.Println("C Table:")
	fmt.Println(alphabet)
	fmt.Println(ctable)
	fmt.Println()

	//O Table Print
	fmt.Println("Otable:")
	printbwt := "     "
	for i := range sa {
		printbwt += bwt(input, sortedSA, i) + " "
	}
	fmt.Println(printbwt)
	for i := range otable {
		fmt.Println(alphabet[i], otable[i])
	}
	fmt.Println()

	//Complexity
	fmt.Println("Time taken for exact match:")
	fmt.Println("Naive match: ", exact[1])
	fmt.Println("BWT search match: ", exact[0])
	fmt.Println()

	//match
	if r == (R - l) {
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

func naiveExactSearch(n string, k string) int {
	counter := 0
	indices := []int{}

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
