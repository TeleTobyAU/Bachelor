package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInput(t *testing.T) {
	info := new(Info)
	info.input = "mississippi"
	if info.input != "mississippi" {
		t.Errorf("Input was %s, instead of mississippi", info.input)
	}
}

func TestAlphabet(t *testing.T) {
	info := new(Info)
	info.input = "mississippi$"

	info.alphabet = generateAlphabet(info.input)
	a := []string{"$", "i", "m", "p", "s"}

	if !reflect.DeepEqual(info.alphabet, a) {
		t.Errorf("Alphabet was %s, instead of %s", info.alphabet, a)
	}
	for i := range a {
		if a[i] != info.alphabet[i] {
			t.Errorf("the order of the alphabet %s is not the same as %s. %s is not the same as %s.", info.alphabet, a, info.alphabet[i], a[i])
		}
	}

}

func TestReverse(t *testing.T) {
	info := new(Info)
	info.input = "mississippi"

	info.reverseInput = Reverse(info.input) + "$"

	if info.reverseInput != "ippississim$" {
		t.Errorf("Reverse input was %s, instead of ippississim$", info.reverseInput)
	}
}

func TestSuffixAndReverseSuffix(t *testing.T) {
	info := new(Info)
	info.key = "a"
	info.input = "abcab$"
	//Sets a thresh hold
	info.threshHold = 1

	//Create alphabet
	info.alphabet = generateAlphabet(info.input)

	//Generate C table
	generateCTable(info)

	info.reverseInput = Reverse(info.input[0:len(info.input)-1]) + "$"
	info.SA = SAIS(info.input)
	info.reverseSA = SAIS(info.reverseInput)

	reversedSA := []int{5, 4, 1, 3, 0, 2}
	sa := []int{5, 3, 0, 4, 1, 2}

	if !reflect.DeepEqual(info.reverseSA, reversedSA) {
		t.Errorf("Reversed suffix array %v, is not %v", info.reverseSA, reversedSA)
	}

	if !reflect.DeepEqual(info.SA, sa) {
		t.Errorf("Suffix array %v, is not %v", info.SA, sa)
	}
}

func TestCigar(t *testing.T) {
	info := new(Info)
	info.key = "iis"
	info.input = "mmiissiissiippii$"
	//Sets a thresh hold
	info.threshHold = 1

	//Create alphabet
	info.alphabet = generateAlphabet(info.input)

	//Generate C table
	generateCTable(info)

	info.reverseInput = Reverse(info.input[0:len(info.input)-1]) + "$"
	info.SA = SAIS(info.input)
	info.reverseSA = SAIS(info.reverseInput)

	//Generate O Table
	generateOTable(info)

	//Init BWT search
	bwtApprox := new(bwtApprox)
	initBwtApproxIter(info.threshHold, info, bwtApprox)

	cigar := []string{"2=1X", "3=", "1I2=", "1=1X1=", "1=1I1=", "2=1D1=", "2=1I"}

	if !reflect.DeepEqual(bwtApprox.cigar, cigar) {
		t.Errorf("CIGAR for mmiissiissiippii$ was %s, but should have been %s", bwtApprox.cigar, cigar)
	}
}

func TestOtable(t *testing.T) {
	info := new(Info)
	info.input = "mmiissiissiippii$"
	info.key = "iss"
	info.alphabet = generateAlphabet(info.input)

	generateCTable(info)

	info.SA = SAIS(info.input)

	generateOTable(info)

	otable := [][]int{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1},
		{0, 1, 2, 2, 2, 2, 2, 3, 4, 5, 5, 5, 5, 6, 6, 6, 7, 8},
		{0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2},
		{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2},
		{0, 0, 0, 0, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 3, 4, 4, 4}}

	printOtable := false

	if !reflect.DeepEqual(info.oTable, otable) {
		t.Errorf("O table is incorrect")
		printOtable = true
	}

	if printOtable {
		//O Table Print
		fmt.Println("Failed Otable calculated:")
		printbwt := "     "
		for i := range info.SA {
			printbwt += bwt(info.input, info.SA, i) + " "
		}
		fmt.Println(printbwt)
		for i := range info.oTable {
			fmt.Println(info.alphabet[i], info.oTable[i])
		}
		fmt.Println()

		//Correct O Table Print
		fmt.Println("Correct Otable:")
		printbwt = "     "
		for i := range info.SA {
			printbwt += bwt(info.input, info.SA, i) + " "
		}
		fmt.Println(printbwt)
		for i := range otable {
			fmt.Println(info.alphabet[i], otable[i])
		}
		fmt.Println()
	}
}
