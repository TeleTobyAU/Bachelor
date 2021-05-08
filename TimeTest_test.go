package Bachelor

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

//TODO ready to run
func TestOptimizedSuffixArraysTime(t *testing.T) {
	//File handling
	//err := os.Remove("TimeOptimizedSAIS.txt")
	//check(err)
	file, err := os.OpenFile("TimeOptimizedSAIS.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	for i := 232000; i <= 300000; i += 1000 {
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

//TODO ready to run
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

	for i := 0; i < 10; i++ {
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

}

//Do not run this, as this has already been ran, takes about 3 hours to run.
func TestCreateBigSuffixArray(t *testing.T) {
	fileSA, err := os.OpenFile("SuffixArray1000000SA.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	fileInput, err := os.OpenFile("SuffixArray1000000Input.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	fileAlphabet, err := os.OpenFile("SuffixArray1000000Alphabet.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	info := new(Info)
	generateRandomNucleotide(1000000, info)

	//Create alphabet
	generateAlphabet(info)

	//Generate C table

	generateCTable(info)

	//Generating SAIS
	info.SA = SAIS(info, info.input)

	var reformatedSA string

	for i := 0; i < len(info.SA); i++ {
		if i == 0 {
			reformatedSA = strconv.Itoa(info.SA[i])
			continue
		}
		reformatedSA = reformatedSA + ", " + strconv.Itoa(info.SA[i])
	}

	_, err = fileSA.WriteString(reformatedSA)
	check(err)
	_, err = fileInput.WriteString(info.input)
	check(err)
	_, err = fileAlphabet.WriteString(strings.Join(info.alphabet, ""))
	check(err)
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
