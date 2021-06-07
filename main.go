package main

import (
	"fmt"
)

func main() {
	info := new(Info)
	info.Key = "iiss"
	//generateRandomNucleotide(500, info)//
	info.Input = "mmiissiissiippii$"

	//Sets a thresh hold
	info.ThreshHold = 1

	//Create alphabet
	info.Alphabet = GenerateAlphabet(info.Input)

	//Generate C table
	GenerateCTableOptimized(info)

	//Generating SAIS
	info.SA = SAIS(info.Input)

	//Reverse the SA string and input string
	info.ReverseInput = Reverse(info.Input[0:len(info.Input)-1]) + "$"
	info.ReverseSA = SAIS(info.ReverseInput) //Making sure the sentinel remains at the end after versing

	//Generate O Table
	GenerateOTable(info)

	//Init BTW search
	ExactMatch(info)

	//Init BWT rec search
	bwtApprox := new(BwtApprox)
	InitBwtApproxIter(info.ThreshHold, info, bwtApprox)

	//Print Cigar
	fmt.Println("CIGAR")
	fmt.Println(bwtApprox.Cigar)
	fmt.Println(info.Alphabet)
	fmt.Println(info.CTable)
	fmt.Println(bwtApprox.DTable)
}
