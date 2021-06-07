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
	SA           []int
	ReverseSA    []int
	CTable       []int
	OTable       [][]int
	roTable      [][]int
}

type InfoInt32 struct {
	Input        string
	ReverseInput string
	Alphabet     []string
	SA           []int32
	ReverseSA    []int32
	CTable       []int32
	OTable       [][]int32
	roTable      [][]int32
}

type BwtExact struct {
	bwtTable *Info
	Key      string
	L        int
	R        int
}

type BwtExact32 struct {
	bwtTable *InfoInt32
	Key      string
	L        int32
	R        int32
}

type BwtApprox struct {
	bwtTable           *Info
	Key                string
	ThreshHold         int
	L, R, nextInterval int
	Ls                 []int
	Rs                 []int
	Cigar              []string
	keyLength          int
	editBuff           []rune
	DTable             []int
	matchLengths       []int
}
