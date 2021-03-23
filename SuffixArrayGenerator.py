# Takes a string input from terminal and generates corresponding suffix array
# Set printResult to true if you want to print
printResult = True
# Sentinel is $

import RandomRNAStringGenerator


#Naive solutions -------------------------------------------------------------------------------------------------------
def generateSuffixes(inputString):
    suffixes = []
    for i in range(len(inputString)):
        suffixes.append(inputString[i: len(inputString)] + inputString[0: i])
    return suffixes


def generateSuffixArray(suffixes):
    suffixArray = []
    for s in sorted(suffixes):
        suffixArray.append(suffixes.index(s))
    return suffixArray


def generateBWT(suffixes, suffixArray):
    bwt = ""
    for i in range(len(suffixArray)):
        bwt += suffixes[suffixArray[i]][-1]
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


def findAlphabet(n):
    alphabet = []
    for i in n:
        if i not in alphabet: alphabet.append(i)
    alphabet.sort()
    return alphabet


#Burrows Wheeler transofmration search ---------------------------------------------------------------------------------
def genCTable (n, alphabet):
    if printResult: print("\nThe C table:")

    #Define alphabet if none given
    if alphabet == None:
        alphabet = findAlphabet(n)

    output = []
    for i in alphabet:
        output.append(0)
        for j in n:
            if i > j: output[alphabet.index(i)] += 1

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



#Linear suffix array construction by almost pure induced sorting -------------------------------------------------------
def LSTypes(n):
    outString = "S"
    inString = n[::-1]
    for i in range(1, len(n)):
        if inString[i - 1] == inString[i]:
            outString += outString[i - 1]
        else:
            if inString[i - 1] < inString[i]:
                outString += "L"
            else:
                outString += "S"
    outString = outString[::-1]

    return outString


def findLMSCSuffixes(n, LSTypesString):
    LMSIndices = []
    if LSTypesString[0] == "S": LMSIndices.append(0)
    for i in range(len(LSTypesString)):
        if LSTypesString[i] == "S" and LSTypesString[i - 1] != "S":
            LMSIndices.append(i)

    return LMSIndices


def findLMSSubstring(n, LMSIndexes):
    LMSSubstringIndices = []
    prev = -1
    for i in range(len(n)):
        if i in LMSIndexes:
            if prev != -1:
                LMSSubstringIndices.append((prev, i))
                prev = i
            else:
                prev = i

    LMSSubStr = []
    for i in LMSSubstringIndices:
        LMSSubStr.append(n[i[0] : i[1] + 1])

    return LMSSubStr



#Pretty printing -------------------------------------------------------------------------------------------------------
def prettyPrint(n, alphabet, S=None, SA=None, BWT=None, LSChars=None, LMSIndices=None, LMSSubStr=None, cTable=None, oTable=None):
    print("String we are working on:", n)
    print()

    if S != None:
        print("Suffixes:")
        for i in S:
            print(i)
        print()

    if SA != None:
        print("Suffix array:", SA)
        print()

    if BWT != None:
        print("Burrows Wheeler transformation string:\n", BWT)
        print()

    if cTable != None:
        print("C table:")
        for i in range(len(cTable)): print(alphabet[i], ":", cTable[i])
        print()

    if oTable != None:
        print(oTable)

    if LSChars != None:
        print("L and S types")
        print(n, "<- Input String")
        print(LSChars, "<- L and S types for each character")

        if LMSIndices != None:
            printString = ""
            for i in range(len(n)):
                if i in LMSIndices:
                    printString += "*"
                else:
                    printString += " "
            print(printString, "<- LMS characters marked with *")

            if LMSSubStr != None:
                printString = ""
                for i in range(len(n)):
                    if i in LMSIndices:
                        printString += "|"
                    else:
                        printString += "-"
                print(printString, "<- LMS substrings visualized")
                print("\nLMS substrings:")
                for i in LMSSubStr: print(i)


#Tests -----------------------------------------------------------------------------------------------------------------
def testMississippi():
    n = "mmiissiissiippii$"
    matchItem = "ss"

    s = generateSuffixes(n)
    sa = generateSuffixArray(s)
    bwt = generateBWT(s, sa)
    alp = findAlphabet(n)
    LS = LSTypes(n)
    LMS = findLMSCSuffixes(n, LS)
    LMSS = findLMSSubstring(n, LMS)
    prettyPrint(n, alp, s, sa, bwt, LS, LMS, LMSS)
    exactSearch(n, matchItem)


def testRandomNucleotideString(nLen, kLen):
    n = RandomRNAStringGenerator.generateString(nLen) + "$"
    matchItem = RandomRNAStringGenerator.generateString(kLen)

    s = generateSuffixes(n)
    sa = generateSuffixArray(s)
    bwt = generateBWT(s, sa)
    alp = findAlphabet(n)
    LS = LSTypes(n)
    LMS = findLMSCSuffixes(n, LS)
    LMSS = findLMSSubstring(n, LMS)
    cTable = genCTable(n, alp)
    prettyPrint(n, alp, s, sa, bwt, LS, LMS, LMSS, cTable)
    exactSearch(n, matchItem)


def testGoogol():
    n = "googol$"
    s = generateSuffixes(n)
    sa = generateSuffixArray(s)
    bwt = generateBWT(s, sa)
    alp = findAlphabet(n)
    LS = LSTypes(n)
    LMS = findLMSCSuffixes(n, LS)
    LMSS = findLMSSubstring(n, LMS)
    prettyPrint(n, alp, s, sa, bwt, LS, LMS, LMSS)


def testABBCABA():
    n = "ABBCABA$"

    s = generateSuffixes(n)
    sa = generateSuffixArray(s)
    bwt = generateBWT(s, sa)
    alp = findAlphabet(n)
    LS = LSTypes(n)
    LMS = findLMSCSuffixes(n, LS)
    LMSS = findLMSSubstring(n, LMS)
    prettyPrint(n, alp, s, sa, bwt, LS, LMS, LMSS)


#testGoogol()
#testMississippi()
testRandomNucleotideString(50, 2)
#testABBCABA()
