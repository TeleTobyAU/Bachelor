package main

import (
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

	generateAlphabet(info)
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
	generateAlphabet(info)

	//Generate C table
	generateCTable(info)

	info.reverseInput = Reverse(info.input[0:len(info.input)-1]) + "$"
	info.SA = SAIS(info, info.input)
	info.RSA = SAIS(info, Reverse(info.input[0:len(info.input)-1])+"$")

	reversedSA := []int{5, 4, 1, 3, 0, 2}
	sa := []int{5, 3, 0, 4, 1, 2}

	if !reflect.DeepEqual(info.RSA, reversedSA) {
		t.Errorf("Reversed suffix array %v, is not %v", info.RSA, reversedSA)
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
	generateAlphabet(info)

	//Generate C table
	generateCTable(info)

	info.reverseInput = Reverse(info.input[0:len(info.input)-1]) + "$"
	info.SA = SAIS(info, info.input)
	info.RSA = SAIS(info, Reverse(info.input[0:len(info.input)-1])+"$")

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
