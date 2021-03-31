package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main() {
	inputString := generateRandomNucleotide(10) + "$"

	//Create alfabet
	alphabet := generateAlfabet(inputString)

	//Creat the SA
	var SA []string
	SA = createSuffixArray(inputString)

	//sort the SA
	sortedSA := sortSuffixArray(SA)

	//Generate C table
	cTable := generateCTable(inputString, alphabet)

	//Generate O Table
	oTable := generateOTable(inputString, alphabet, sortedSA)

	finePrint(inputString, alphabet, SA, sortedSA, cTable, oTable)

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

//TODO rewrite to own solution
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

func finePrint(input string, alphabet []string, sa []string, sortedSA []int, ctable []int, otable [][]int) {
	//Input string
	fmt.Println("\nInput String:")
	fmt.Println(input)
	fmt.Println()

	//Alphabet
	fmt.Println("\nAlphabet over input string:")
	fmt.Println(alphabet)
	fmt.Println()

	//suffix array without sort
	fmt.Println("\nSuffix:")
	fmt.Println(sa)
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

func naiveExactSearch(n string, k string) {
	counter := 0
	indices := []int{}
	fmt.Println(n)
	fmt.Println("This is the string we are searching for " + k)

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

	fmt.Println("number of exact match: ", counter)
	fmt.Println("indices of the exact match ", indices)
}
