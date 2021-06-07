package main

import (
	"fmt"
)

func main() {
	info := new(Info)
	//generateRandomNucleotide(500, info)//
	info.Input = "mmiissiissiippii$"

	//Create alphabet
	info.Alphabet = GenerateAlphabet(info.Input)

	//Generate C table
	info.CTable = GenerateCTableOptimized(info.Input, info.Alphabet)

	//Generating SAIS
	info.SA = SAISv1(info.Input)

	//Reverse the SA string and input string
	info.ReverseInput = Reverse(info.Input[0:len(info.Input)-1]) + "$"
	info.ReverseSA = SAISv1(info.ReverseInput) //Making sure the sentinel remains at the end after versing

	//Generate O Table
	GenerateOTable(info)

	//Init BTW search
	bwtE := new(BwtExact)
	bwtE.Key = "iiss"
	bwtE.bwtTable = info

	ExactMatch(bwtE)

	//Init BWT rec search
	bwtApprox := new(BwtApprox)
	bwtApprox.ThreshHold = 1
	InitBwtApproxIter(bwtApprox.ThreshHold, info, bwtApprox)

	//Print Cigar
	fmt.Println("CIGAR")
	fmt.Println(bwtApprox.Cigar)
	fmt.Println(info.Alphabet)
	fmt.Println(info.CTable)
	fmt.Println(bwtApprox.DTable)
}
