package main

import (
	"math/rand"
	"reflect"
	"sort"
	"time"
)

//https://stackoverflow.com/questions/1752414/how-to-reverse-a-string-in-go
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

//https://github.com/heapwolf/go-indexof/blob/master/indexof.go
func IndexOf(params ...interface{}) int {
	v := reflect.ValueOf(params[0])
	arr := reflect.ValueOf(params[1])

	var t = reflect.TypeOf(params[1]).Kind()

	if t != reflect.Slice && t != reflect.Array {
		panic("Type Error! Second argument must be an array or a slice.")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == v.Interface() {
			return i
		}
	}
	return -1
}

func generateAlphabet(info *Info) {
	var alphabet []string
	inputString := info.input

	for s := range inputString {
		found := false
		for i := range alphabet {
			if string(inputString[s]) == alphabet[i] {
				found = true
				break
			}
		}
		if !found {
			alphabet = append(alphabet, string(inputString[s]))
		}
	}
	sort.Strings(alphabet)
	info.alphabet = alphabet
}

func generateRandomNucleotide(size int, info *Info) {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("ATCG")

	nucleotide := make([]rune, size)

	for i := range nucleotide {
		nucleotide[i] = letters[rand.Intn(len(letters))]
	}
	info.input = string(nucleotide) + "$"
}
