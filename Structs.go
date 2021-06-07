package main

type NaiveStruct struct {
	Input, ReverseInput, key  string
	Alphabet                  []string
	thresHold                 int
	SA, ReverseSA             []int
	stringSA, stringReverseSA []string
}

type Info struct {
	Input        string
	ReverseInput string
	Alphabet     []string
	ThreshHold   int
	Key          string
	SA           []int
	ReverseSA    []int
	CTable       []int
	OTable       [][]int
	roTable      [][]int
	L            int
	R            int
}

type InfoInt32 struct {
	Input        string
	ReverseInput string
	Alphabet     []string
	ThreshHold   int
	Key          string
	SA           []int32
	ReverseSA    []int32
	CTable       []int32
	OTable       [][]int32
	roTable      [][]int32
	L            int32
	R            int32
}

type BwtApprox struct {
	bwtTable           *Info
	key                string
	L, R, nextInterval int
	Ls                 []int
	Rs                 []int
	Cigar              []string
	keyLength          int
	editBuff           []rune
	DTable             []int
	matchLengths       []int
}
