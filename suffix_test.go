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
	info.input = "mississippi$"

	reverse(info)

	if info.reverseInput != "$ippississim" {
		t.Errorf("Reverse input was %s, instead of $ippississim", info.reverseInput)
	}
}

func TestSuffixAndReverseSuffix(t *testing.T) {
	info := new(Info)
	info.input = "abcab$"

	reverse(info)
	generateAlphabet(info)
	createSuffixArray(info)

	a := []string{"abcab$", "bcab$a", "cab$ab", "ab$abc", "b$abca", "$abcab"}
	b := []string{"$bacba", "bacba$", "acba$b", "cba$ba", "ba$bac", "a$bacb"}

	if !reflect.DeepEqual(info.stringReverseSA, b) {
		t.Errorf("suffix was %s, instead of %s", info.stringReverseSA, b)
	}

	if !reflect.DeepEqual(info.StringSA, a) {
		t.Errorf("suffix was %s, instead of %s", info.StringSA, a)
	}

	for i := range a {
		if a[i] != info.StringSA[i] {
			t.Errorf("the order of the suffix %s is not the same as %s. %s is not the same as %s.", info.StringSA, a, info.StringSA[i], a[i])
		}
	}
}

func TestSortSuffixAndSortReverseSuffix(t *testing.T) {
	info := new(Info)
	info.input = "abcab$"

	reverse(info)
	generateAlphabet(info)
	createSuffixArray(info)
	sortSuffixArray(info)

	a := []int{5, 3, 0, 4, 1, 2}
	b := []int{0, 5, 2, 4, 1, 3}

	if !reflect.DeepEqual(info.SA, a) {
		t.Errorf("Sorted suffix was %v, instead of %v", info.SA, b)
	}

	if !reflect.DeepEqual(info.reverseSA, b) {
		t.Errorf("Sorted reverse suffix was %v, instead of %v", info.reverseSA, b)
	}

	for i := range a {
		if a[i] != info.SA[i] {
			t.Errorf("the order of the sorted suffix %v is not the same as %v. %v is not the same as %v.", info.SA, a, info.SA[i], a[i])
		}
	}
}
