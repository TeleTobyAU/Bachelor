package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGenerateCTable(t *testing.T) {
	info := new(Info)
	info.Input = "AAATTTTTCCCCGGGG$"

	info.Alphabet = GenerateAlphabet(info.Input)

	info.CTable, _ = GenerateCTableOptimized(info.Input, info.Alphabet, false)

	ctable := []int{0, 1, 4, 8, 12}

	if !reflect.DeepEqual(ctable, info.CTable) {
		t.Errorf("Failed to create C table %v, should have been %v", info.CTable, ctable)
	}
}

func TestOtable(t *testing.T) {
	info := new(Info)
	info.Input = "mmiissiissiippii$"
	info.Key = "iss"
	info.Alphabet = GenerateAlphabet(info.Input)

	info.CTable, _ = GenerateCTableOptimized(info.Input, info.Alphabet, false)

	info.SA = SAISv1(info.Input)

	GenerateOTable(info)

	otable := [][]int{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1},
		{0, 1, 2, 2, 2, 2, 2, 3, 4, 5, 5, 5, 5, 6, 6, 6, 7, 8},
		{0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2},
		{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2},
		{0, 0, 0, 0, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 3, 4, 4, 4}}

	printOtable := false

	if !reflect.DeepEqual(info.OTable, otable) {
		t.Errorf("O table is incorrect")
		printOtable = true
	}

	if printOtable {
		//O Table Print
		fmt.Println("Failed Otable calculated:")
		printbwt := "     "
		for i := range info.SA {
			printbwt += Bwt(info.Input, info.SA, i) + " "
		}
		fmt.Println(printbwt)
		for i := range info.OTable {
			fmt.Println(info.Alphabet[i], info.OTable[i])
		}
		fmt.Println()

		//Correct O Table Print
		fmt.Println("Correct Otable:")
		printbwt = "     "
		for i := range info.SA {
			printbwt += Bwt(info.Input, info.SA, i) + " "
		}
		fmt.Println(printbwt)
		for i := range otable {
			fmt.Println(info.Alphabet[i], otable[i])
		}
		fmt.Println()
	}
}

func TestExactMatch(t *testing.T) {
	info := new(Info)
	info.Key = "ATCG"
	info.Input = GenerateRandomNucleotide(100000)
	//Sets a thresh hold

	//Create alphabet
	info.Alphabet = GenerateAlphabet(info.Input)

	//Generate SA-IS
	info.SA = SAISv1(info.Input)

	//Generate C table
	GenerateCTable(info)

	//Generate O Table
	GenerateOTable(info)

	naiveSolutions := NaiveExactSearch(info.Key, info.Input)

	InitBwtSearch(info)
	output := IndexBwtSearch(info)

	if naiveSolutions != len(output) {
		t.Errorf("Exact Match failed %v, but should have been %v", len(output), naiveSolutions)
	}
}

func TestCigar(t *testing.T) {
	info := new(Info)
	info.Key = "iissii"
	info.Input = "mmiissiissiippii$"
	//Sets a thresh hold
	info.ThreshHold = 1

	//Create alphabet
	info.Alphabet = GenerateAlphabet(info.Input)

	//Generate C table
	GenerateCTable(info)

	info.ReverseInput = Reverse(info.Input[0:len(info.Input)-1]) + "$"
	info.SA = SAISv1(info.Input)
	info.ReverseSA = SAISv1(info.ReverseInput)

	//Generate O Table
	GenerateOTable(info)

	//Init BWT search
	bwtApprox := new(BwtApprox)
	InitBwtApproxIter(info.ThreshHold, info, bwtApprox)

	//Old cigar [2=1X 3= 1I2= 1=1X1= 1=1I1= 2=1D1= 2=1I]
	cigar := []string{"2=1X", "3=", "1I2=", "1=1X1=", "1=1I1=", "2=1D1=", "2=1I"}

	if !reflect.DeepEqual(bwtApprox.Cigar, cigar) {
		t.Errorf("CIGAR for mmiissiissiippii$ was %s, but should have been %s", bwtApprox.Cigar, cigar)
	}
}
