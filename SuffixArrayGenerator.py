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

def exactSearch(n, k):
    print("\nWe are matching", k, "to", n + ":")

    matches = 0
    for i in range(len(n)):
        if n[i : i + len(k)] == k:
            matches += 1

            if printResult:
                print(k, "matched at index", i)
    if printResult:
        print("Matches: ", matches)


def genCTable (input, alphabet):
    if printResult:
        print("\nThe C table:")

    #Define alphabet if none given
    input += '$'
    if alphabet == None:
        alphabet = []
        for i in input:
            if i not in alphabet:
                alphabet.append(i)
        alphabet.sort()

    output = []
    for i in alphabet:
        output.append(0)
        for j in input:
            if i > j:
                output[alphabet.index(i)] += 1

    if printResult:
        for i in range(len(output)):
            print(alphabet[i], ":", output[i])

    return output

#inp = input("Type string you wish to create suffix array from: \n")
#inp = "mississippi"
inp = RandomRNAStringGenerator.generateString(100)
suffixArray = generateSuffixArray(inp)
bwt = generateBWT(suffixArray)
genCTable(inp, None)

matchItem = "GCCGT"

exactSearch(inp, matchItem)
