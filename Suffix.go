package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	array := createSuffixArray("AATAAACCTTACCTAGCACTCCATCATGTCTTATGGCGCGTGATTTGCCCCGGACTCAGG$")
	fmt.Println("\nSuffix Array without sort: \n", array)
	sortedArray := sortSuffixArray(array)
	fmt.Println("\nSuffix Array with sort: \n", sortedArray)

	BWT := findBWT(sortedArray)
	BWTString := strings.Join(BWT, "")
	fmt.Println("\nBurrows–Wheeler transform: \n", BWTString)

}

func findBWT(array []string) []string {
	length := len(array)
	bwt := []string{}
	for _, s := range array {
		bwt = append(bwt, string(s[length-1]))
	}

	return bwt
}

func sortSuffixArray(array []string) []string {

	sort.Strings(array)

	return array
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

/*
   input := "Hello"
   	index := suffixarray.New([]byte(input))
   	offsets := index.Lookup([]byte("$"), -1)

   	//Print out the index of
   	fmt.Printf("offsets %v \n", offsets)
*/
