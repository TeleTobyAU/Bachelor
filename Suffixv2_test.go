package main

import (
	"reflect"

	"testing"
)

func TestSuffixAndReverseSuffix(t *testing.T) {
	input := "abcab$"

	SA := SAISv1(input)

	sa := []int{5, 3, 0, 4, 1, 2}

	if !reflect.DeepEqual(SA, sa) {
		t.Errorf("Suffix array %v, is not %v", SA, sa)
	}
}

func TestReverseSuffix(t *testing.T) {
	input := "abcab$"

	reverseInput := Reverse(input[0:len(input)-1]) + "$"

	reverseSA := SAISv1(reverseInput)
	reversedSA := []int{5, 4, 1, 3, 0, 2}
	if !reflect.DeepEqual(reverseSA, reversedSA) {
		t.Errorf("Reversed suffix array %v, is not %v", reverseSA, reversedSA)
	}
}

func TestSuffixWithRecursive(t *testing.T) {
	input := "mmiissiissiippii$"

	output := SAISv1(input)

	sa := []int{16, 15, 14, 10, 6, 2, 11, 7, 3, 1, 0, 13, 12, 9, 5, 8, 4}

	if !reflect.DeepEqual(output, sa) {
		t.Errorf("Suffix array %v, is not %v", output, sa)
	}
}

func TestSuffixWith1000charactersLong(t *testing.T) {
	info := new(NaiveStruct)
	info.Input = GenerateRandomNucleotide(10)

	CreateSuffixArrayNaive(info)
	SortSuffixArrayNaive(info)
	naiveSa := info.SA

	sais := SAISv1(info.Input)

	//Converts to int32
	var naiveSaInt32 []int32
	for x := range naiveSa {
		naiveSaInt32 = append(naiveSaInt32, int32(naiveSa[x]))
	}

	if !reflect.DeepEqual(sais, naiveSa) {
		t.Errorf("Suffix array %v, is not %v", sais, naiveSa)
	}
}
