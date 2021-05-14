package main

import "fmt"

func main() {
	info := new(Info)
	info.key = "iiss"
	//generateRandomNucleotide(500, info)//
	info.input = "mmiissiissiippii$"

	//Sets a thresh hold
	info.threshHold = 1

	//Create alphabet
	generateAlphabet(info)

	//Generate C table
	generateCTable(info)

	//Generating SAIS
	info.SA = SAIS(info, info.input)

	//Reverse the SA string and input string
	info.reverseInput = Reverse(info.input[0:len(info.input)-1]) + "$"
	info.reverseSA = SAIS(info, Reverse(info.input[0:len(info.input)-1])+"$") //Making sure the sentinel remains at the end after versing

	//Generate O Table
	generateOTable(info)

	//Init BTW search
	exactMatch(info)

	//Init BWT rec search
	bwtApprox := new(bwtApprox)
	initBwtApproxIter(info.threshHold, info, bwtApprox)

	//Print Cigar
	fmt.Println("CIGAR")
	fmt.Println(bwtApprox.cigar)
	fmt.Println(info.alphabet)
	fmt.Println(info.cTable)
	fmt.Println(bwtApprox.dTable)
}
