package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestOptimizedSuffixArraysTimeOnBigString(t *testing.T) {
	file, err := os.OpenFile("TimeOptimizedSAISOnBigString.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	for i := 100000; i <= 10000000; i += 1000000 {
		info := new(Info)
		generateRandomNucleotide(i, info)

		//Create alphabet
		generateAlphabet(info)

		//Generate C table

		generateCTable(info)
		fmt.Println("SAIS")

		//SAIS
		start := time.Now()
		info.SA = SAIS(info, info.input)
		timeSAIS := time.Since(start).Seconds()

		fmt.Println("SAIS was created in", timeSAIS)

		//Printing to file with the result of the two test
		s1 := "SAIS " + strconv.Itoa(int(timeSAIS))
		s3 := " length " + strconv.Itoa(int(i))
		n := s1 + s3 + "\n"
		_, err = file.WriteString(n)
		check(err)
	}
}

func TestTimeEverything(t *testing.T) {
	file1, err := os.OpenFile("TimeCTable.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	file2, err := os.OpenFile("TimeSAIS.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	file3, err := os.OpenFile("TimeReverseSAIS.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	file4, err := os.OpenFile("TimeOTable.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	file5, err := os.OpenFile("TimeExactMatch.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	file6, err := os.OpenFile("TimeRecApproxMatch.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	for i := 1000000000; i <= 1000000000; i += 1000000 {
		info := new(Info)
		fmt.Println("Generating nucleotide")
		generateRandomNucleotide(i, info)
		fmt.Println("Generating alphabet")
		generateAlphabet(info)
		info.threshHold = 1
		info.key = "AAT"

		//C Table
		fmt.Println("Generating C Table")
		start := time.Now()
		generateCTable(info)
		endTimeCTable := time.Since(start).Milliseconds()
		fmt.Println("C Table is created for ", i)

		//SAIS
		start = time.Now()
		info.SA = SAIS(info, info.input)
		endTimeSAIS := time.Since(start).Seconds()
		fmt.Println("SAIS is created for ", i)

		//Reverse SAIS and input
		start = time.Now()
		info.reverseInput = Reverse(info.input[0:len(info.input)-1]) + "$"
		info.reverseSA = SAIS(info, Reverse(info.input[0:len(info.input)-1])+"$")
		endTimeReverseSAIS := time.Since(start).Seconds()
		fmt.Println("Reverse SAIS and Reverse input is created for ", i)

		//O Table
		start = time.Now()
		generateOTable(info)
		endTimeOTable := time.Since(start).Milliseconds()
		fmt.Println("O Table is created for ", i)

		//Exact Match
		start = time.Now()
		initBwtSearch(info)
		exactMatch := indexBwtSearch(info)
		endTimeExactMatch := time.Since(start).Microseconds()
		fmt.Println("Exact Match is created for ", i, " and there were ", len(exactMatch), " matches")

		//Rec Approx Match
		start = time.Now()
		bwtApprox := new(bwtApprox)
		initBwtApproxIter(info.threshHold, info, bwtApprox)
		endTimeRecApprox := time.Since(start).Microseconds()
		fmt.Println("Rec Approx is created for ", i, " and there were ", len(bwtApprox.cigar), " matches")

		//Write to file C table
		reformatedStringCTable := "CTable " + strconv.Itoa(int(endTimeCTable)) + " Size " + strconv.Itoa(i) + "\n"
		_, err = file1.WriteString(reformatedStringCTable)
		check(err)

		//Write to file SAIS
		reformatedStringSAIS := "SAIS " + strconv.Itoa(int(endTimeSAIS)) + " Size " + strconv.Itoa(i) + "\n"
		_, err = file2.WriteString(reformatedStringSAIS)
		check(err)

		//Write to file Reverse SAIS
		reformatedStringReverseSAIS := "ReverseSAIS " + strconv.Itoa(int(endTimeReverseSAIS)) + " Size " + strconv.Itoa(i) + "\n"
		_, err = file3.WriteString(reformatedStringReverseSAIS)
		check(err)

		//Write to file O table
		reformatedStringOTable := "OTable " + strconv.Itoa(int(endTimeOTable)) + " Size " + strconv.Itoa(i) + "\n"
		_, err = file4.WriteString(reformatedStringOTable)
		check(err)

		//Write to file Exact Match
		reformatedStringExactMatch := "ExactMatch " + strconv.Itoa(int(endTimeExactMatch)) + " Size " + strconv.Itoa(i) + "\n"
		_, err = file5.WriteString(reformatedStringExactMatch)
		check(err)

		//Write to file Rec Approx
		reformatedStringRecApprox := "RecApproxMatch " + strconv.Itoa(int(endTimeRecApprox)) + " Size " + strconv.Itoa(i) + "\n"
		_, err = file6.WriteString(reformatedStringRecApprox)
		check(err)
	}
}

//Has already been ran
func TestTimeCTable(t *testing.T) {
	file, err := os.OpenFile("TimeCTable.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	for i := 0; i < 1000000; i += 1000 {
		info := new(Info)
		generateRandomNucleotide(i, info)
		generateAlphabet(info)
		fmt.Println("C Table is created for ", i)
		start := time.Now()
		generateCTable(info)
		endTime := time.Since(start).Milliseconds()

		reformatedString := "CTable " + strconv.Itoa(int(endTime)) + " Size " + strconv.Itoa(i) + "\n"
		_, err = file.WriteString(reformatedString)
		check(err)
	}
}

//Has already been ran
func TestOptimizedSuffixArraysTime(t *testing.T) {
	//File handling
	//err := os.Remove("TimeOptimizedSAIS.txt")
	//check(err)
	file, err := os.OpenFile("TimeOptimizedSAIS.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	//114000 - 118000
	//148000 - 160000

	for i := 114000; i <= 118000; i += 1000 {
		info := new(Info)
		generateRandomNucleotide(i, info)

		//Create alphabet
		generateAlphabet(info)

		//Generate C table

		generateCTable(info)
		fmt.Println("SAIS")
		//SAIS
		start := time.Now()
		info.SA = SAIS(info, info.input)
		timeSAIS := time.Since(start).Milliseconds()

		fmt.Println("Reverse SAIS, time for SAIS: ", timeSAIS)
		//Reverse SAIS
		start = time.Now()
		info.reverseInput = Reverse(info.input[0:len(info.input)-1]) + "$"
		info.reverseSA = SAIS(info, Reverse(info.input[0:len(info.input)-1])+"$")
		timeReverseSAIS := time.Since(start).Milliseconds()

		//Printing to file with the result of the two test
		s1 := "SAIS " + strconv.Itoa(int(timeSAIS))
		s2 := " ReverseSAIS " + strconv.Itoa(int(timeReverseSAIS))
		s3 := " length " + strconv.Itoa(int(i))
		n := s1 + s2 + s3 + "\n"
		_, err = file.WriteString(n)
		check(err)
		fmt.Printf("SAIS %v ReverseSAIS %v \n", timeSAIS, timeReverseSAIS)
	}
}

//Has already been ran
func TestNaiveSAAndSAIS(t *testing.T) {
	//File handling
	//err := os.Remove("TimeNaiveAndSAIS.txt")
	//check(err)
	file, err := os.OpenFile("TimeNaiveAndSAIS.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	for i := 1000; i <= 100000; i += 1000 {
		fmt.Println("initializing resources for SAIS")
		info := new(Info)
		generateRandomNucleotide(i, info)

		//Create alphabet
		generateAlphabet(info)

		//Generate C table
		generateCTable(info)

		//SAIS
		start := time.Now()
		info.SA = SAIS(info, info.input)
		timeSAIS := time.Since(start).Milliseconds()
		fmt.Println("SAIS is done in", timeSAIS, "ms")

		//Naive SA
		fmt.Println("initializing resources for Naive SA")
		info2 := new(Info)
		info2.input = info.input
		info2.alphabet = info.alphabet
		info2.cTable = info.cTable

		start = time.Now()
		createSuffixArrayNaive(info2)
		sortSuffixArrayNaive(info2)
		timeNaiveSA := time.Since(start).Milliseconds()
		fmt.Println("Naive SA is done in", timeNaiveSA, "ms")

		//Printing to file with the result of the two test
		s1 := "SAIS " + strconv.Itoa(int(timeSAIS))
		s2 := " NaiveSA " + strconv.Itoa(int(timeNaiveSA))
		s3 := " length " + strconv.Itoa(int(i))
		n := s1 + s2 + s3 + "\n"
		_, err = file.WriteString(n)
		check(err)
		fmt.Printf("SAIS %v NaiveSA %v \n", timeSAIS, timeNaiveSA)
	}
}

//TODO ready to run
func TestOTableWithMemoryPrint(t *testing.T) {
	//get input from txt file into an int array
	dataSA, err := ioutil.ReadFile("SuffixArray1000000SA.txt")
	check(err)

	var sa []int
	for s := 0; s < len(dataSA); s++ {
		sa = append(sa, int(dataSA[s]))
	}
	fmt.Println("Sa is created")

	//Input
	dataInput, err := ioutil.ReadFile("SuffixArray1000000Input.txt")
	check(err)

	var input string
	input = string(dataInput)
	fmt.Println("input is created", len(input))

	//Alphaet
	dataAlphabet, err := ioutil.ReadFile("SuffixArray1000000Alphabet.txt")
	check(err)

	var alphabet []string
	for a := 0; a < len(dataAlphabet); a++ {
		alphabet = append(alphabet, string(dataAlphabet[a]))
	}
	fmt.Println("alphabet is created")

	info := new(Info)
	info.input = input
	info.SA = sa
	info.alphabet = alphabet

	//Generate C table
	generateCTable(info)

	//Generate O Table
	runtime.GC()
	MemUsage("OtableMemory.txt")
	generateOTable(info)
	MemUsage("OtableMemory.txt")
	runtime.GC()

}

func TestExactMatch(t *testing.T) {
	info := new(Info)
	generateRandomNucleotide(100000, info)

	info.key = "AAA"

	info.threshHold = 0
	generateAlphabet(info)

	//Generate C table
	generateCTable(info)

	info.SA = SAIS(info, info.input)
	//Otable
	generateOTable(info)

	fmt.Println(naiveExactSearch(info))
	initBwtSearch(info)
	fmt.Println(len(indexBwtSearch(info)))
	fmt.Println(info.input)
}

//function to read/write to a file
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func MemUsage(fileName string) {
	fileSA, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v MiB", bytesToMegabyte(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bytesToMegabyte(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bytesToMegabyte(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)

	reformatedString := "Alloc " + strconv.Itoa(int(bytesToMegabyte(m.Alloc))) +
		" TotalAlloc " + strconv.Itoa(int(bytesToMegabyte(m.TotalAlloc))) +
		" Sys " + strconv.Itoa(int(bytesToMegabyte(m.Sys))) +
		" NumGC " + strconv.Itoa(int(m.NumGC)) + "\n"

	_, err = fileSA.WriteString(reformatedString)
	check(err)
}

func bytesToMegabyte(b uint64) uint64 {
	return b / 1024 / 1024
}
