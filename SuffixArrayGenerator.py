# Takes a string input from terminal and generates corresponding suffix array
# Set printResult to true if you want to print
printResult = True
# Sentinel is $

import RandomRNAStringGenerator

def generateSuffixArray(inputString):
    workingString = inputString + '$'

    suffixArray = []
    for i in range (len(workingString)):
        suffixArray.append(workingString[i : len(workingString)] + workingString[0 : i])

    if printResult:
        print("\nPivots:")
        for s in suffixArray:
            print(s)

    return suffixArray

def generateBWT(suffixArray):
    sortedSuffixArray = suffixArray.copy()
    sortedSuffixArray.sort()

    bwt = ""
    for i in range(len(sortedSuffixArray)):
        item = sortedSuffixArray[i]
        bwt += item[-1]

    if printResult:
        print("\nSorted suffix array:")
        for s in sortedSuffixArray:
            print(s)
        print("\nBurrows Wheeler transformation string:")
        print(bwt)

def exactSearch(BWTString, suffixArray, stringsToBeMatched):
    m = len(stringsToBeMatched)
    matches = 0
    for i in stringsToBeMatched:
        for j in suffixArray:
            if i in j:
                matches += 1
                break
    print("\nNumber of matches:", matches)
    print("Match percentage:", (matches/m)*100,"%")

def genCTable (input):
    #Count the number of occurences of each nucleotide
    output = [0,0,0,0,0]
    for i in input:
        if i == '$':
            output[0] += 1
        if i == 'A':
            output[1] += 1
        if i == 'C':
            output[2] += 1
        if i == 'G':
            output[3] += 1
        if i == 'T':
            output[4] += 1

    if printResult:
        print("\n[$, A, C, G, T]")
        print(output)

    return output

#inp = input("Type string you wish to create suffix array from: \n")
#inp = "mississipi"
inp = RandomRNAStringGenerator.generateString(1000)
suffixArray = generateSuffixArray(inp)
bwt = generateBWT(suffixArray)
genCTable(inp)

matchItems = [
"ACCGT",
"AAATA",
"AAGTC",
"TCGGG",
"TCTCT",
"AGGAG",
"TGAGT",
"TAATT",
"TGGAT",
"TATCT",
"CATAT",
"ACATG",
"CTGCT",
"AGCCT",
"AAGCG",
"AAGCG",
"TTCGT",
"CCATT",
"ACGAT",
"AGAGT",
"AGTAG",
"TAAAG",
"TTTCT",
"TGTTA",
"GCAAC",
"CTGCT",
"ATCTG",
"ACCAC",
"AGTCT",
"TCGAC",
"ACTTA",
"CCAAC",
"CGGAA",
"GGATC",
"CTGCT",
"ACTCT",
"CGCCA",
"AGATG",
"ATGCG",
"GTTAC",
"GCCCT",
"AATTG",
"ATGGC",
"GGGAT",
"AGGTG",
"TAGAA",
"CTCTT",
"AACTG",
"ATGGC",
"CGTCA",
"GGTAA",
"CGAGA",
"ATTTG",
"TAACA",
"TTGCA",
"GTTAC",
"TAGTG",
"AACTC",
"GTCCC",
"TAGTG",
"CTGCA",
"GGATA",
"GTTCG",
"AATGC",
"GCTAA",
"GCTTC",
"AACGT",
"TTTCA",
"TTATT",
"TCATG",
"CGAGC",
"GTTTC",
"CATGA",
"TTGCT",
"AACGG",
"CATTG",
"GGACT",
"TCCAC",
"TCGGA",
"AAAAG",
"CTGCG",
"GATAG",
"TTCAC",
"CGTCG",
"GACGT",
"GGCCT",
"TGCCA",
"CCGGC",
"CGTAA",
"AGGAC",
"ACAGG",
"ACATT",
"TCGTC",
"CGGAC",
"AAGAC",
"TACAC",
"CTGGA",
"GAGCG",
"TCATA",
"GTGTG"
]

exactSearch(bwt, suffixArray, matchItems)
