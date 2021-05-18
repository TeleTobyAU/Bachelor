package main

type NaiveStruct struct {
	input, reverseInput, key  string
	alphabet                  []string
	thresHold                 int
	SA, reverseSA             []int
	stringSA, stringReverseSA []string
}

type Info struct {
	input        string
	reverseInput string
	alphabet     []string
	threshHold   int
	key          string
	SA           []int
	reverseSA    []int
	cTable       []int
	oTable       [][]int
	roTable      [][]int
	L            int
	R            int
}

type bwtApprox struct {
	bwtTable           *Info
	key                string
	L, R, nextInterval int
	Ls                 []int
	Rs                 []int
	cigar              []string
	keyLength          int
	editBuff           []rune
	dTable             []int
	matchLengths       []int
}
