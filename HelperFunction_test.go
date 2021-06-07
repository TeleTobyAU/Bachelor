package main

import (
	"reflect"

	"testing"
)

func TestAlphabet(t *testing.T) {
	input := "mississippi$"
	alphabet := GenerateAlphabet(input)
	a := []string{"$", "i", "m", "p", "s"}

	if !reflect.DeepEqual(alphabet, a) {
		t.Errorf("Alphabet was %s, instead of %s", alphabet, a)
	}
	for i := range a {
		if a[i] != alphabet[i] {
			t.Errorf("the order of the alphabet %s is not the same as %s. %s is not the same as %s.", alphabet, a, alphabet[i], a[i])
		}
	}

}

func TestReverse(t *testing.T) {
	input := "mississippi"

	reverseInput := Reverse(input) + "$"

	if reverseInput != "ippississim$" {
		t.Errorf("Reverse input was %s, instead of ippississim$", reverseInput)
	}
}

func TestRandomStringLen(t *testing.T) {
	//Generates a string of len size + 1 sentinel
	randomString := GenerateRandomNucleotide(1000000)
	if len(randomString) != 1000001 {
		t.Errorf("Length of reverse input was %v, instead of 1000001", len(randomString))
	}
}

func TestIndexOf(t *testing.T) {
	a := []int{4, 2, 7, 6, 8, 1, 5}
	output := IndexOf(7, a)

	if output != 2 {
		t.Errorf("IndexOf found index %v but should have been 2", output)
	}
}
