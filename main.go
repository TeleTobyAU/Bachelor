package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

//go run main.go Structs.go HelperFunctions.go Search.go Suffixv1.go

func main() {
	info := new(Info)
	info.Alphabet = []string{"$", "A", "C", "G", "T"}
	bwtApprox := new(BwtApprox)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Testing BWA and SAIS")
	fmt.Println("Commands are: quit, test, and custom")
	fmt.Println("---------------------")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		text = strings.TrimSpace(text)
		text = strings.ToLower(text)

		if strings.Compare("quit", text) == 0 {
			fmt.Println("Okay, bye bye!")
			os.Exit(1)
		} else if strings.Compare("test", text) == 0 {
			info.Input = GenerateRandomNucleotide(100)

			info.SA = SAISv1(info.Input)
			info.ReverseInput = Reverse(info.Input[0:len(info.Input)-1]) + "$"
			info.ReverseSA = SAISv1(info.ReverseInput)

			info.CTable = GenerateCTableOptimized(info.Input, info.Alphabet)
			GenerateOTable(info)
			GenerateOTableReverse(info)

			keys := []string{}
			for i := 0; i < 5; i++ {
				keySize := 5
				keys = append(keys, GenerateRandomNucleotide(keySize)[:keySize])
				time.Sleep(100)
			}

			fmt.Println("----------------------------------")
			fmt.Println("The input string is:", info.Input)
			fmt.Println("The keys are:", keys)
			fmt.Println("----------------------------------")

			cigars := [][]string{}
			positionL := [][]int{}
			positionR := [][]int{}
			for i := 0; i < 5; i++ {
				bwtApprox.ThreshHold = 1
				bwtApprox.Key = keys[i]
				InitBwtApproxIter(bwtApprox.ThreshHold, info, bwtApprox)

				cigars = append(cigars, bwtApprox.Cigar)
				bwtApprox.Cigar = []string{}

				positionL = append(positionL, bwtApprox.Ls)
				positionR = append(positionR, bwtApprox.Rs)
				bwtApprox.Rs = []int{}
				bwtApprox.Ls = []int{}
			}

			fmt.Println("CIGARs and positions for each key")
			for i := 0; i < 5; i++ {
				fmt.Println("\n", keys[i])
				if len(cigars[i]) == 0 {
					fmt.Println("There where no matches for this key")
				}
				for j := 0; j < len(cigars[i]); j++ {
					indices := []int{}
					for k := 0; k < len(positionL[i]); k++ {
						indices = append(indices, info.SA[positionL[i][k]])
					}
					fmt.Print("At positions: ", indices, " in the input, we have a match: ")
					fmt.Println(cigars[i][j])
				}
			}
		} else if strings.Compare("custom", text) == 0 {
			fmt.Print("Please provide an input string, must be string of A, C, G and T: ")
			input, _ := reader.ReadString('\n')
			// convert CRLF to LF
			input = strings.Replace(input, "\n", "", -1)
			input = strings.TrimSpace(input)

			info.Input = input + "$"

			info.SA = SAISv1(info.Input)
			info.ReverseInput = Reverse(info.Input[0:len(info.Input)-1]) + "$"
			info.ReverseSA = SAISv1(info.ReverseInput)

			info.CTable = GenerateCTableOptimized(info.Input, info.Alphabet)
			GenerateOTable(info)
			GenerateOTableReverse(info)

			fmt.Print("Please provide key, must be string of A, C, G and T: ")
			key, _ := reader.ReadString('\n')
			key = strings.Replace(key, "\n", "", -1)
			key = strings.TrimSpace(key)

			fmt.Print("Please provide the maximum allowed variations: ")
			max, _ := reader.ReadString('\n')
			max = strings.Replace(max, "\n", "", -1)
			max = strings.TrimSpace(max)
			maxInt, _ := strconv.Atoi(max)

			bwtApprox.ThreshHold = maxInt
			bwtApprox.Key = key
			InitBwtApproxIter(bwtApprox.ThreshHold, info, bwtApprox)

			fmt.Println("CIGARs and positions for each key")
			if len(bwtApprox.Cigar) == 0 {
				fmt.Println("There where no matches for this key")
			} else {
				indices := []int{}
				for k := 0; k < len(bwtApprox.Ls); k++ {
					indices = append(indices, info.SA[bwtApprox.Ls[k]])
				}
				fmt.Print("At positions: ", indices, " in the input, we have a match: ")
				fmt.Println(bwtApprox.Cigar)
			}
		}
	}
}
