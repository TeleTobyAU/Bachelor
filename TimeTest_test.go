package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"
)

/*
func TestOTableTime(t *testing.T) {
	//File handling
	err := os.Remove("TimeDataOTable.txt")
	check(err)
	file, err := os.OpenFile("TimeDataOTable.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	for i := 10000; i <= 100000; i += 1000 {
		info := new(Info)
		generateRandomNucleotide(i, info)
		//info.input = "mmiissiissiippii$"

		//Create alphabet
		generateAlphabet(info)

		//Generate C table

		generateCTable(info)
		info.reverseInput = Reverse(info.input[0:len(info.input)-1]) + "$"
		info.SA = SAIS(info, info.input)
		info.reverseSA = SAIS(info, Reverse(info.input[0:len(info.input)-1])+"$")

		//Generate O Table
		start := time.Now()
		MemUsage()
		generateOTable(info)
		MemUsage()
		timeOTable := time.Since(start).Milliseconds()

		s1 :=  "OTable " + strconv.Itoa(int(timeOTable))
		s2 := " length " + strconv.Itoa(int(i))
		n := s1 + s2 + "\n"
		_, err = file.WriteString(n)
		check(err)
	}


}
*/

func TestOptimizedSuffixArraysTime(t *testing.T) {
	//File handling
	err := os.Remove("TimeOptimizedSAIS.txt")
	check(err)
	file, err := os.OpenFile("TimeOptimizedSAIS.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	check(err)

	for i := 1000; i <= 1000000; i += 1000 {
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

//functions to write to a file
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func MemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v MiB", bytesToMegabyte(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bytesToMegabyte(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bytesToMegabyte(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bytesToMegabyte(b uint64) uint64 {
	return b / 1024 / 1024
}

/*
	 //Init info 2
			info2 := new(Info)
			info2.input = info.input
			generateAlphabet(info2)

			start = time.Now()
			//SA
			createSuffixArray(info2)
			sortSuffixArray(info2)
			timeSA := time.Since(start).Milliseconds()

			if !reflect.DeepEqual(info.SA, info2.SA) {
				fmt.Println("Sufiix arrays isn't equal")
			}


			s1 := "SAIS " + strconv.Itoa(int(timeSAIS))
			s2 := " SA " + strconv.Itoa(int(timeSA))
			s3 := " length " + strconv.Itoa(int(i))
			n := s1 + s2 + s3 + "\n"
			_, err = file.WriteString(n)
			check(err)
			fmt.Printf("SAIS %v SA %v \n", timeSAIS, timeSA)
		}
*/
