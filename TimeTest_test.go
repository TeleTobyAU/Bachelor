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

func TestRun(t *testing.T) {
	//OtableCreation(t)

	fmt.Println("Input string = HELLO")
	fmt.Println("Key = ELO and threshold = 1")
	fmt.Println("X and =")
	fmt.Println("1=1X2=")
	fmt.Println("Insertion")
	fmt.Println("2=1I, 1I2=")
	fmt.Println("Deletion")
	fmt.Println("1=1D2=")

	OtableCreation(t)

}

func OtableCreation(t *testing.T) {
	info := new(Info)
	info.Input = "mississippi$"

	info.Alphabet = GenerateAlphabet(info.Input)

	GenerateCTable(info)

	info.SA = SAISv1(info.Input)

	GenerateOTable(info)
	fmt.Println(info.SA)

	fmt.Println("F          L")

	for i := range info.SA {
		fmt.Println(info.Input[info.SA[i]:] + info.Input[:info.SA[i]])
	}

	fmt.Println("F   L")
	for i := range info.SA {
		fmt.Println(info.Input[info.SA[i]:][0:1] + "   " + (info.Input[info.SA[i]:] + info.Input[:info.SA[i]])[11:12])
	}

	fmt.Println(" $ i m p s ")
	fmt.Println(info.CTable)

	fmt.Println()

	printbwt := "  "
	for i := range info.Alphabet {
		printbwt += info.Alphabet[i] + " "
	}
	fmt.Println(len(info.OTable))
	fmt.Println(printbwt)

	j := 0
	for i := range info.OTable[j] {
		if i == 0 {
			fmt.Println("L", info.OTable[j][i], info.OTable[j+1][i], info.OTable[j+2][i], info.OTable[j+3][i], info.OTable[j+4][i])
			continue
		}
		fmt.Println(Bwt(info.Input, info.SA, i-1), info.OTable[j][i], info.OTable[j+1][i], info.OTable[j+2][i], info.OTable[j+3][i], info.OTable[j+4][i])
	}

	fmt.Println()
	fmt.Println(printbwt)
	j = 0
	for i := range info.OTable[j] {
		if i == 0 {
			fmt.Println("L", info.OTable[j][i], info.OTable[j+1][i], info.OTable[j+2][i], info.OTable[j+3][i], info.OTable[j+4][i])
			continue
		}
		if i == 5 || i == 10 {
			fmt.Println(Bwt(info.Input, info.SA, i-1), info.OTable[j][i], info.OTable[j+1][i], info.OTable[j+2][i], info.OTable[j+3][i], info.OTable[j+4][i])
		} else {
			fmt.Println(Bwt(info.Input, info.SA, i-1), "-", "-", "-", "-", "-")
		}
	}

	fmt.Println()
	for i := range info.SA {
		if i == 6 {
			fmt.Println("CTable[3] + OTable[3, 0]  ->", info.Input[info.SA[i]:]+info.Input[:info.SA[i]])
		} else if i == 8 {
			fmt.Println("CTable[3] + OTable[3, 12] ->", info.Input[info.SA[i]:]+info.Input[:info.SA[i]])
		} else {
			fmt.Println("                            ", info.Input[info.SA[i]:]+info.Input[:info.SA[i]])
		}
	}

	fmt.Println()
	for i := range info.SA {
		if i == 1 {
			fmt.Println("               CTable[1]  ->", info.Input[info.SA[i]:]+info.Input[:info.SA[i]])
		} else if i == 2 {
			fmt.Println("CTable[1] + OTable[1, 6]  ->", info.Input[info.SA[i]:]+info.Input[:info.SA[i]])
		} else if i == 3 {
			fmt.Println("CTable[1] + OTable[1, 8]  ->", info.Input[info.SA[i]:]+info.Input[:info.SA[i]])
		} else {
			fmt.Println("                            ", info.Input[info.SA[i]:]+info.Input[:info.SA[i]])
		}
	}
}

//Has already been ran
func OptimizedSuffixArraysTime(t *testing.T) {
	//File handling
	//err := os.Remove("TimeOptimizedSAIS.txt")
	//check(err)
	file, err := os.OpenFile("DATA/TimeOptimizedSAISv2.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	for i := 100000000; i <= 1000000000; i += 100000000 {
		//Generating nucleotide
		fmt.Println("Generating input")
		input := GenerateRandomNucleotide(i)

		//SAIS
		fmt.Println("Creating SA-IS")
		start := time.Now()
		SA := SAIS(input)
		timeSAIS := time.Since(start).Seconds()
		fmt.Println("SA-IS created in", timeSAIS, "with len", len(SA))

		SA = nil

		//Printing to file with the result
		s1 := "SAIS " + strconv.Itoa(int(timeSAIS))
		s3 := " length " + strconv.Itoa(int(i))
		n := s1 + s3 + "\n"
		_, err = file.WriteString(n)
		check(err)
	}
}

//Has already been ran
func NaiveSAAndSAIS(t *testing.T) {
	//File handling
	//err := os.Remove("TimeNaiveAndSAIS.txt")
	//check(err)
	file, err := os.OpenFile("DATA/TimeNaiveAndSAIS.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	for i := 1000; i <= 100000; i += 1000 {
		fmt.Println("initializing resources for SAIS")
		info := new(InfoInt32)
		info.Input = GenerateRandomNucleotide(i)

		//Create alphabet
		info.Alphabet = GenerateAlphabet(info.Input)

		//Generate C table
		_, info.CTable = GenerateCTableOptimized(info.Input, info.Alphabet, true)

		//SAIS
		start := time.Now()
		info.SA = SAIS(info.Input)
		timeSAIS := time.Since(start).Milliseconds()
		fmt.Println("SAIS is done in", timeSAIS, "ms")

		/*//Naive SA
		fmt.Println("initializing resources for Naive SA")
		info2 := new(NaiveStruct)
		info2.Input = info.Input
		info2.Alphabet = info.Alphabet

		start = time.Now()
		CreateSuffixArrayNaive(info2)
		SortSuffixArrayNaive(info2)
		timeNaiveSA := time.Since(start).Milliseconds()
		fmt.Println("Naive SA is done in", timeNaiveSA, "ms")
		*/
		//Printing to file with the result of the two test
		s1 := "SAIS " + strconv.Itoa(int(timeSAIS))
		s3 := " length " + strconv.Itoa(int(i))
		n := s1 + s3 + "\n"
		_, err = file.WriteString(n)
		check(err)
		fmt.Printf("SAIS %v NaiveSA %v \n", timeSAIS, 0)

	}
}

//Has already been ran
func TimeCTable(t *testing.T) {
	file, err := os.OpenFile("DATA/TimeCTable.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	for i := 100000000; i <= 3500000000; i += 100000000 {
		info := new(Info)
		info.Input = GenerateRandomNucleotide(i)
		info.Alphabet = GenerateAlphabet(info.Input)
		fmt.Println("C Table is created for ", i)
		start := time.Now()
		info.CTable, _ = GenerateCTableOptimized(info.Input, info.Alphabet, false)
		endTime := time.Since(start).Milliseconds()

		reformatedString := "CTable " + strconv.Itoa(int(endTime)) + " Size " + strconv.Itoa(i) + "\n"
		_, err = file.WriteString(reformatedString)
		check(err)
	}
}

func TimeOTable(t *testing.T) {
	file, err := os.OpenFile("DATA/TimeOTable.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	file2, err2 := os.OpenFile("DATA/TimeExactMatch.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err2)
	for j := 11; j <= 41; j += 10 {
		a := GenerateRandomNucleotide(j)
		for i := 100000000; i <= 1500000000; i += 100000000 {
			info := new(InfoInt32)
			info.Input = GenerateRandomNucleotide(i)
			fmt.Println("Len of input:", len(info.Input))

			info.Key = a[:j-1]
			fmt.Println("Key len:", len(info.Key))

			info.Alphabet = GenerateAlphabet(info.Input)
			fmt.Println("Alphabet len:", len(info.Alphabet))

			GenerateCTable32(info)
			fmt.Println("C Table:", info.CTable)

			info.SA = SAIS(info.Input)
			fmt.Println("Len of SA:", len(info.SA))

			start := time.Now()
			GenerateOTable32(info)
			endTimeOTable := time.Since(start).Milliseconds()
			fmt.Println("Otable is created")

			start = time.Now()
			InitBwtSearch32(info)
			a := IndexBwtSearch32(info)
			endTimeExact := time.Since(start).Milliseconds()
			fmt.Println("Len of matches:", len(a))

			reformatedStringO := "OTable " + strconv.Itoa(int(endTimeOTable)) + " Size " + strconv.Itoa(i) + "\n"
			_, err = file.WriteString(reformatedStringO)
			check(err)
			reformatedStringE := "Exact " + strconv.Itoa(int(endTimeExact)) + " Size " + strconv.Itoa(i) + "\n"
			_, err2 = file2.WriteString(reformatedStringE)
			check(err2)
		}
		reformatedStringO := "-----------------------------------------------------------\n"
		_, err = file.WriteString(reformatedStringO)

		reformatedStringE := "-----------------------------------------------------------\n"
		_, err2 = file2.WriteString(reformatedStringE)
		check(err2)
	}
}

func TimeExactSearch(t *testing.T) {
	file, err := os.OpenFile("DATA/TimeExactMatch.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	info := new(InfoInt32)
	info.Input = GenerateRandomNucleotide(150000000)
	fmt.Println("Len of input:", len(info.Input))

	fmt.Println("Key len:", len(info.Key))

	info.Alphabet = GenerateAlphabet(info.Input)
	fmt.Println("Alphabet len:", len(info.Alphabet))

	GenerateCTable32(info)
	fmt.Println("C Table:", info.CTable)

	info.SA = SAIS(info.Input)
	fmt.Println("Len of SA:", len(info.SA))

	GenerateOTable32(info)
	fmt.Println("Otable is created")

	for i := 1; i <= 100; i++ {
		info.Key = GenerateRandomNucleotide(i)[:i-1]
		start := time.Now()
		InitBwtSearch32(info)
		a := IndexBwtSearch32(info)
		endTimeExact := time.Since(start).Milliseconds()
		fmt.Println("Len of matches:", len(a))

		reformatedStringE := "Exact " + strconv.Itoa(int(endTimeExact)) + " Size " + strconv.Itoa(i) + "\n"
		_, err = file.WriteString(reformatedStringE)
		check(err)
	}
}

//Has already been ran
func TimeTestRecApprox(t *testing.T) {
	file, err := os.OpenFile("DATA/TimeRecApprox1.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	info := new(Info)

	//Generating nucleotide
	fmt.Println("Generating input")
	info.Input = GenerateRandomNucleotide(250000000)

	info.Alphabet = GenerateAlphabet(info.Input)

	//SAIS
	fmt.Println("Creating SA-IS")
	info.SA = SAISv1(info.Input)

	fmt.Println("Creating reverse SAIS")
	info.ReverseInput = Reverse(info.Input[0:len(info.Input)-1]) + "$"
	info.ReverseSA = SAISv1(info.ReverseInput)

	//C Table
	fmt.Println("Generating C table")
	GenerateCTable(info)

	//O table
	fmt.Println("Generating O table")
	GenerateOTable(info)

	//Reverse O table
	fmt.Println("Generating reverse O table")
	GenerateOTableReverse(info)

	//Creating Rec Approx struct
	fmt.Println("Creating and init Approx struct")
	recApprox := new(BwtApprox)
	recApprox.bwtTable = info
	fmt.Println("Len of input ", len(info.SA))

	for j := 50; j <= 100; j += 50 {
		a := GenerateRandomNucleotide(j)
		info.Key = a[:j-1]
		fmt.Println(info.Key)
		for i := 0; i <= 8; i += 1 {
			info.ThreshHold = i
			startTime := time.Now()
			InitBwtApproxIter(info.ThreshHold, info, recApprox)
			endTime := time.Since(startTime)
			fmt.Println("endTime for ", i, ": ", endTime)
			reformattedStringRecApprox := "RecApproxMatch " + endTime.String() + "\n"
			_, err = file.WriteString(reformattedStringRecApprox)
			check2(err, reformattedStringRecApprox)
		}
		_, err = file.WriteString("-----------------------------------------------------------\n")
		check(err)
	}
}

func TimeTestRecApprox2(t *testing.T) {
	file, err := os.OpenFile("DATA/TimeRecApprox2.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	info := new(Info)

	//Generating nucleotide
	fmt.Println("Generating input")
	info.Input = GenerateRandomNucleotide(500000000)

	info.Alphabet = GenerateAlphabet(info.Input)

	//SAIS
	fmt.Println("Creating SA-IS")
	info.SA = SAISv1(info.Input)

	fmt.Println("Creating reverse SAIS")
	info.ReverseInput = Reverse(info.Input[0:len(info.Input)-1]) + "$"
	info.ReverseSA = SAISv1(info.ReverseInput)

	//C Table
	fmt.Println("Generating C table")
	GenerateCTable(info)

	//O table
	fmt.Println("Generating O table")
	GenerateOTable(info)

	//Reverse O table
	fmt.Println("Generating reverse O table")
	GenerateOTableReverse(info)

	//Creating Rec Approx struct
	fmt.Println("Creating and init Approx struct")
	recApprox := new(BwtApprox)
	recApprox.bwtTable = info
	fmt.Println("Len of input ", len(info.SA))
	timeREC := int64(0)
	for k := 0; k <= 5; k++ {
		reformattedStringRecApprox := "Edits " + strconv.Itoa(k) + ":\n-----------------------------------------------------------\n"
		_, err = file.WriteString(reformattedStringRecApprox)
		check2(err, reformattedStringRecApprox)
		for j := 10; j <= 200; j += 10 {
			for i := 0; i < 100; i++ {
				a := GenerateRandomNucleotide(j)
				info.Key = a[:j-1]
				info.ThreshHold = k
				startTime := time.Now()
				InitBwtApproxIter(info.ThreshHold, info, recApprox)
				endTime := time.Since(startTime).Milliseconds()
				//fmt.Println("endTime for ", j, ": ", endTime)
				timeREC += endTime
			}
			reformattedStringRecApprox := "RecApproxMatch " + strconv.Itoa(int(timeREC/100)) + " Size " + strconv.Itoa(j) + "\n"
			_, err = file.WriteString(reformattedStringRecApprox)
			check2(err, reformattedStringRecApprox)
			timeREC = 0
		}
	}

}

func TimeEverything(t *testing.T) {
	file1, err := os.OpenFile("DATA/TimeCTable.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	file2, err := os.OpenFile("DATA/TimeSAIS.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	file3, err := os.OpenFile("DATA/TimeReverseSAIS.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	file4, err := os.OpenFile("DATA/TimeOTable.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	file5, err := os.OpenFile("DATA/TimeExactMatch.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	file6, err := os.OpenFile("DATA/TimeRecApproxMatch.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	for i := 1000000000; i <= 4000000000; i += 1000000000 {
		info := new(Info)
		fmt.Println("Generating nucleotide")
		info.Input = GenerateRandomNucleotide(i)
		fmt.Println("Generating alphabet")
		info.Alphabet = GenerateAlphabet(info.Input)
		info.ThreshHold = 1
		info.Key = "AAT"

		//C Table
		fmt.Println("Generating C Table")
		start := time.Now()
		info.CTable, _ = GenerateCTableOptimized(info.Input, info.Alphabet, false)
		endTimeCTable := time.Since(start).Milliseconds()
		fmt.Println("C Table is created for ", i)

		//SAIS
		fmt.Println("Starting on SAIS")
		start = time.Now()
		info.SA = SAISv1(info.Input)
		endTimeSAIS := time.Since(start).Seconds()
		fmt.Println("SAIS is created for ", i)

		//Reverse SAIS and input
		fmt.Println("Starting on reverse input")
		info.ReverseInput = Reverse(info.Input[0:len(info.Input)-1]) + "$"
		fmt.Println("Reverse input created")
		fmt.Println("Starting on reverse SAIS")

		start = time.Now()
		info.ReverseSA = SAISv1(info.ReverseInput)
		endTimeReverseSAIS := time.Since(start).Seconds()
		fmt.Println("Reverse SAIS and Reverse input is created for ", i)

		fmt.Println("Starting on O table")
		//O Table
		start = time.Now()
		GenerateOTable(info)
		endTimeOTable := time.Since(start).Milliseconds()
		fmt.Println("O Table is created for ", i)

		//Exact Match
		start = time.Now()
		InitBwtSearch(info)
		exactMatch := IndexBwtSearch(info)
		endTimeExactMatch := time.Since(start).Microseconds()
		fmt.Println("Exact Match is created for ", i, " and there were ", len(exactMatch), " matches")

		//Rec Approx Match
		start = time.Now()
		bwtApprox := new(BwtApprox)
		InitBwtApproxIter(info.ThreshHold, info, bwtApprox)
		endTimeRecApprox := time.Since(start).Microseconds()
		fmt.Println("Rec Approx is created for ", i, " and there were ", len(bwtApprox.Cigar), " matches")

		//Write to file C table
		reformatedStringCTable := "CTable " + strconv.Itoa(int(endTimeCTable)) + " Size " + strconv.Itoa(i) + "\n"
		_, err = file1.WriteString(reformatedStringCTable)
		check2(err, reformatedStringCTable)

		//Write to file SAIS
		reformatedStringSAIS := "SAIS " + strconv.Itoa(int(endTimeSAIS)) + " Size " + strconv.Itoa(i) + "\n"
		_, err = file2.WriteString(reformatedStringSAIS)
		check2(err, reformatedStringSAIS)

		//Write to file Reverse SAIS
		reformatedStringReverseSAIS := "ReverseSAIS " + strconv.Itoa(int(endTimeReverseSAIS)) + " Size " + strconv.Itoa(i) + "\n"
		_, err = file3.WriteString(reformatedStringReverseSAIS)
		check2(err, reformatedStringReverseSAIS)

		//Write to file O table
		reformatedStringOTable := "OTable " + strconv.Itoa(int(endTimeOTable)) + " Size " + strconv.Itoa(i) + "\n"
		_, err = file4.WriteString(reformatedStringOTable)
		check2(err, reformatedStringOTable)

		//Write to file Exact Match
		reformatedStringExactMatch := "ExactMatch " + strconv.Itoa(int(endTimeExactMatch)) + " Size " + strconv.Itoa(i) + "\n"
		_, err = file5.WriteString(reformatedStringExactMatch)
		check2(err, reformatedStringExactMatch)

		//Write to file Rec Approx
		reformatedStringRecApprox := "RecApproxMatch " + strconv.Itoa(int(endTimeRecApprox)) + " Size " + strconv.Itoa(i) + "\n"
		_, err = file6.WriteString(reformatedStringRecApprox)
		check2(err, reformatedStringRecApprox)
	}
}

func OTableWithMemoryPrint(t *testing.T) {
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
	info.Input = input
	info.SA = sa
	info.Alphabet = alphabet

	//Generate C table
	info.CTable, _ = GenerateCTableOptimized(info.Input, info.Alphabet, false)

	//Generate O Table
	runtime.GC()
	MemUsage("OtableMemory.txt")
	GenerateOTable(info)
	MemUsage("OtableMemory.txt")
	runtime.GC()

}

//function to read/write to a file
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func check2(err error, failedString string) {
	if err != nil {
		//panic(err)
		fmt.Println("Fail with ", err)
		fmt.Println("Printing string to terminal instead")
		fmt.Println()
		fmt.Println(failedString)
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
