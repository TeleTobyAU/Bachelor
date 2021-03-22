# Takes a string input from terminal and generates corresponding suffix array
# Set printResult to true if you want to print
printResult = True
# Sentinel is $

import RandomRNAStringGenerator



def generateSuffixes(inputString):
    suffixes = []
    for i in range(len(inputString)):
        suffixes.append(inputString[i: len(inputString)] + inputString[0: i])

    if printResult:
        print("\nSuffixes:")
        for s in suffixes:
            print(s)

        print("\nSorted suffixes:")
        for s in sorted(suffixes):
            print(s)

    return suffixes


def generateSuffixArray(suffixes):
    suffixArray = []
    for s in sorted(suffixes):
        suffixArray.append(suffixes.index(s))

    if printResult:
        print("\nSuffixarray:")
        print(suffixArray)

    return suffixArray


def generateBWT(suffixes, suffixArray):
    bwt = ""
    for i in range(len(suffixArray)):
        bwt += suffixes[suffixArray[i]][-1]

    if printResult:
        print("\nBurrows Wheeler transformation string:")
        print(bwt)

    return bwt


def exactSearch(n, k):
    if printResult: print("\nWe are matching", k, "to", n + ":")

    matches = []
    for i in range(len(n)):
        if n[i : i + len(k)] == k:
            matches.append(i)
            if printResult: print(k, "matched at index", i)

    if printResult: print("Number of matches: ", len(matches))

    return matches


def genCTable (n, alphabet):
    if printResult: print("\nThe C table:")

    #Define alphabet if none given
    if alphabet == None:
        alphabet = []
        for i in n:
            if i not in alphabet: alphabet.append(i)
        alphabet.sort()

    output = []
    for i in alphabet:
        output.append(0)
        for j in n:
            if i > j: output[alphabet.index(i)] += 1

    if printResult:
        for i in range(len(output)): print(alphabet[i], ":", output[i])

    return output


def genOTable(n, sa, BWT, alphabet):
    if printResult: print("\nThe O table:")

    #Define alphabet if none given
    if alphabet == None:
        alphabet = []
        for i in n:
            if i not in alphabet: alphabet.append(i)
        alphabet.sort()

    oTableSize = len(alphabet) #* len(sa) + 1 #* len(oTableSize)
    oIndicesSize = len(sa) + 1 #* len(oIndicesSize)
    if printResult: print("O table size:", oTableSize, "\nO indices size:", oIndicesSize)

    oTable = []
    oIndices = []

    #First weird loop
    #put lists of correct size into the o indices table
    for i in range(oTableSize):
        oIndices.append([])
        for j in range(oIndicesSize):
            oIndices[i].append(None)

    for i in range(1, oTableSize):
        oIndices[i][0] = 0

    if printResult:
        for i in oIndices:
            print(i)



def testMississippi():
    inp = "mississippi$"
    matchItem = "ss"

    s = generateSuffixes(inp)
    suffixArray = generateSuffixArray(s)
    bwt = generateBWT(s, suffixArray)
    genCTable(inp, None)
    genOTable(inp, suffixArray, bwt, None)
    exactSearch(inp, matchItem)


def testRandomNucleotideString(nLen, kLen):
    inp = RandomRNAStringGenerator.generateString(nLen) + "$"
    matchItem = RandomRNAStringGenerator.generateString(kLen)

    s = generateSuffixes(inp)
    suffixArray = generateSuffixArray(s)
    bwt = generateBWT(s, suffixArray)
    genCTable(inp, None)
    genOTable(inp, suffixArray, bwt, None)
    exactSearch(inp, matchItem)


def testGoogol():
    n = "googol$"
    s = generateSuffixes(n)
    sa = generateSuffixArray(s)
    bwt = generateBWT(s, sa)



#testGoogol()
testMississippi()
#testRandomNucleotideString(100, 2)
