package Bachelor

import (
	"reflect"
	"testing"
)

func TestNaiveSA(t *testing.T) {
	info := new(Info)
	info.input = "GGCAATATCTGTAAGCTTAGTGTGCGTGCTTTGTCTGCACCTCTAGGTACGCTGATCGTACAGTTGGCGTAGGCTCCTATACCGGGAACCCTCTGTGAAA$"
	generateAlphabet(info)

	correctSA := []int{100, 99, 98, 97, 86, 12, 3, 59, 87, 80, 38, 48, 13, 70, 44,
		18, 61, 78, 4, 54, 6, 2, 37, 60, 88, 81, 75, 39, 89, 49, 82, 56, 67, 24, 42,
		76, 73, 40, 90, 51, 34, 8, 92, 15, 28, 96, 85, 53, 1, 36, 66, 23, 72, 50, 14,
		27, 84, 0, 65, 71, 83, 45, 10, 57, 46, 68, 32, 94, 21, 25, 19, 62, 11, 58, 79,
		47, 69, 43, 17, 77, 5, 74, 55, 41, 33, 7, 91, 95, 52, 35, 22, 26, 64, 9, 31, 93, 20, 16, 63, 30, 29}
	//SA
	createSuffixArrayNaive(info)
	sortSuffixArrayNaive(info)

	if !reflect.DeepEqual(correctSA, info.SA) {
		t.Errorf("SA is wrong %v and should have been %v", info.SA, correctSA)
	}
}

func TestReverseNaiveSA(t *testing.T) {
	info := new(Info)
	info.reverseInput = Reverse("GGCAATATCTGTAAGCTTAGTGTGCGTGCTTTGTCTGCACCTCTAGGTACGCTGATCGTACAGTTGGCGTAGGCTCCTATACCGGGAACCCTCTGTGAAA$")
	generateAlphabet(info)

	correctReverseSA := []int{0, 1, 96, 13, 2, 87, 39, 97, 62, 14, 46, 3, 94, 20, 22, 56, 30,
		41, 52, 88, 82, 12, 61, 19, 40, 51, 11, 60, 18, 10, 24, 85, 49, 98, 27, 33, 63, 72, 76,
		44, 92, 58, 8, 25, 66, 100, 86, 38, 55, 29, 81, 50, 17, 32, 75, 43, 99, 54, 28, 16, 15, 34,
		47, 90, 6, 64, 79, 73, 4, 77, 35, 68, 95, 45, 93, 21, 59, 9, 23, 84, 48, 26, 71, 91, 57, 7, 65,
		37, 80, 31, 74, 42, 53, 89, 5, 78, 67, 83, 70, 36, 69}
	//Reverse SA
	createReverseSuffixArrayNaive(info)
	sortReverseSuffixArrayNaive(info)

	if !reflect.DeepEqual(correctReverseSA, info.reverseSA) {
		t.Errorf("Reverse SA is wrong %v and should have been %v", info.reverseSA, correctReverseSA)
	}
}
