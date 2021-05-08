from Main.python import RandomRNAStringGenerator
import time


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
    matches = []
    for i in range(len(n)):
        if n[i : i + len(k)] == k:
            matches.append(i)

    return matches


def approximateSearch(n, k, threshHold=0):
    matches = []
    for i in range(len(n) - len(k)):
        hammingDistance = 0
        for j in range(i, i + len(k)):
            if n[j] != k[j - i]:
                hammingDistance += 1
                if hammingDistance > threshHold:
                    break
            if j == i + len(k) - 1:
                matches.append(i)

    return matches


def findAlphabet(n):
    alphabet = []
    for i in n:
        if i not in alphabet: alphabet.append(i)
    alphabet.sort()
    return alphabet


#Burrows Wheeler transofmration search ---------------------------------------------------------------------------------
def genCTable (n, alphabet):
    #Define alphabet if none given
    if alphabet == None:
        alphabet = findAlphabet(n)

    output = []
    for i in alphabet:
        output.append(0)
        for j in n:
            if i > j: output[alphabet.index(i)] += 1

    return output


def genOTable(BWT, alphabet):
    oIndices = alphabet.copy()
    oTable = []

    for i in range(len(alphabet)):
        oTable.append([0])

    for i in range(len(BWT)):
        for j in range(len(alphabet)):
            if BWT[i] == alphabet[j]:
                oTable[j].append(oTable[j][i] + 1)
            else:
                oTable[j].append(oTable[j][i])

    return oIndices, oTable


def initBwtSearchIter(n, k, cTable, oIndex, oTable):
    L = 0
    R = len(n)

    #If our key is longer than our string there will be no match
    if len(k) > len(n):
        R = 0
        L = 1

    i = len(k) - 1

    while i >= 0 and L < R:
        a = oIndex[0:].index(k[i])
        L = cTable[a] + oTable[a][L]
        R = cTable[a] + oTable[a][R]
        i -= 1

    i = L
    LRI = (L, R, i)
    return LRI


def getMatchIndices(sa, LRI):
    matches = []
    L = LRI[0]
    R = LRI[1]

    for i in range(R - L):
        matches.append(sa[L + i])
    matches.sort()

    return matches


def genDTable(n, alp, k, rsa, cTable, roIndex, roTable):
    dTable = []
    minEdits = 0
    L = 0
    R = len(n)

    for i in range(len(k)):
        a = roIndex[0:].index(k[i])
        L = cTable[a] + roTable[a][L]
        R = cTable[a] + roTable[a][R]

        if L >= R:
            minEdits += 1
            L = 0
            R = len(n)
        dTable.append(minEdits)

    return dTable


def initBWTApproxSearch(n, alp, k, cTable, oIndex, oTable, maxEdits):
    L = 0
    R = len(n)
    i = len(k) - 1
    matchA = k[i]

    for i in range(1, len(alp) - 1):
        newL = cTable[a] + oTable[a][L]
        newR = cTable[a] + oTable[a][R]

        if a == matchA:
            editCost = 0
        else:
            editCost = 1

        if not maxEdits - editCost < 0:
            break
        if not newL >= newR:
            break

        edits = "M"
        recursiveApproxMatch(L, R, i - 1, 0, maxEdits - 1, edits + 1)

    return -1


def recursiveApproxMatch(L, R, i , editsLeft, maxEdits , edits):


    return -1


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


def findLMSSuffixes(n, LSTypesString):
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

    LMSSubStr = []
    for i in LMSSubstringIndices:
        LMSSubStr.append(n[i[0] : i[1] + 1])

    return LMSSubStr


def genBuckets(n, alphabet, suffixes, LMSIndices, LMSsubstrings):
    buckets = []
    for i in alphabet:
        buckets.append([i])
    print(buckets)

    for i in suffixes:
        for j in range(len(buckets)):
            if i[0] == buckets[j][0]:
                buckets[j].append(i)
                break

    for i in buckets:
        i.sort()

    return buckets



#Pretty printing -------------------------------------------------------------------------------------------------------
def prettyPrint(n, k, alphabet, S=None, SA=None, buckets=None, BWT=None, LSChars=None, LMSIndices=None, LMSSubStr=None, cTable=None, oIndices=None, oTable=None, LRI=None, searchResults=None, approxResults=None):
    print("\nString we are working on:", n)
    print("In the alphabet:", alphabet)
    print()

    if S != None:
        print("Suffixes:")
        for i in S:
            print(i)
        print()
        for i in sorted(S):
            print(sorted(S).index(i), i)
        print()

    if SA != None:
        print("Suffix array:", SA)
        print()

    if buckets != None:
        print("Buckets:")
        for i in buckets:
            print(i)
        print()

    if BWT != None:
        print("Burrows Wheeler transformation string:\n", BWT)
        print()

    if cTable != None:
        print("C table:")
        for i in range(len(cTable)): print(alphabet[i], ":", cTable[i])
        print()
        print(cTable)
        print()

    if oTable != None and oIndices != None:
        print("Otable:")
        printBwt = "      "
        for i in BWT:
            printBwt += i + "  "
        print(printBwt)
        for i in range(len(oTable)):
            print(oIndices[i], oTable[i])
        print()

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
        print()

    if k != None:
        print("Looking for", k, "in", n)

    if LRI != None:
        print("L:", LRI[0])
        print("R:", LRI[1])
        print("Number of matches:", LRI[1] - LRI[0])
        print()

    if searchResults != None:
        width = 300
        widthCheck = width
        print("Found results at indices:", searchResults)
        printString = ""
        i = 0
        while i < len(n):
            if i not in searchResults:
                printString += " "
                i += 1
            else:
                printString += "|"
                if len(k) > 1:
                    for j in range(len(k) - 2):
                        printString += "-"
                    printString += "|"
                i += len(k)
            if i >= widthCheck:
                print(n[i - width:i])
                print(printString)
                printString = ""
                widthCheck += width
        print(n[widthCheck - width : len(n) - 1])
        print(printString)
        print()

    if approxResults != None:
        print("Approximative search results:")
        print("Found", len(approxResults), "approximate hits")
        print("their starting indices listed below:")
        print(approxResults)

    print("--------------------------------------------------------------------------")
    print()

#Tests -----------------------------------------------------------------------------------------------------------------
def testMississippi():
    n = "mmiissiissiippii$"
    k = "iss"
    runTest(n, k, 1)

def testRandomNucleotideString(nLen, kLen):
    n = RandomRNAStringGenerator.generateString(nLen) + "$"
    k = RandomRNAStringGenerator.generateString(kLen)
    #timeTest(n, k)
    runTest(n, k, 1)

def testRandomNucleotideStringMoreK(nLen, kLen, kNum, a=0):
    n = RandomRNAStringGenerator.generateString(nLen) + "$"
    k = []
    for i in range(kNum):
        k.append(RandomRNAStringGenerator.generateString(kLen))
    timeTestSeveralK(n, k)

def testGoogol():
    n = "googol$"
    k = ""
    runTest(n, k, 0)

def testABBCABA():
    n = "ABBCABA$"
    k = "AB"
    runTest(n, k, 0)


def runTest(n, k, approxSearchDistance):
    s = generateSuffixes(n)
    sa = generateSuffixArray(s)
    bwt = generateBWT(s, sa)
    alp = findAlphabet(n)
    LS = LSTypes(n)
    LMS = findLMSSuffixes(n, LS)
    LMSS = findLMSSubstring(n, LMS)
    buckets = genBuckets(n, alp, s, LMS, LMSS)

    cTable = genCTable(n, alp)
    oInd, oTable = genOTable(bwt, alp)
    lri = initBwtSearchIter(n, k, cTable, oInd, oTable)
    matchIndexes = getMatchIndices(sa, lri)

    approx = approximateSearch(n, k, approxSearchDistance)

    prettyPrint(n, k, alp, s, sa, buckets, bwt, LS, LMS, LMSS, cTable, oInd, oTable, lri, matchIndexes, approx)

    n = n[::-1]
    s = generateSuffixes(n)
    sa = generateSuffixArray(s)
    bwt = generateBWT(s, sa)
    cTable = genCTable(n, alp)
    oInd, oTable = genOTable(bwt, alp)
    dTable = genDTable(n, alp, k, sa, cTable, oInd, oTable)
    print("dtable", dTable)
    lri = initBwtSearchIter(n, k, cTable, oInd, oTable)
    matchIndexes = getMatchIndices(sa, lri)

    prettyPrint(n, k, alp, s, sa, buckets, bwt, LS, LMS, LMSS, cTable, oInd, oTable, lri, matchIndexes, approx)

def timeTest(n, k):
    start = time.time_ns()
    s = generateSuffixes(n)
    print("Suffixes done", (time.time_ns() - start) / 1000, "micro seconds")
    start = time.time_ns()
    sa = generateSuffixArray(s)
    print("Sorting done", (time.time_ns() - start) / 1000, "micro seconds")
    print()

    start = time.time_ns()
    bwt = generateBWT(s, sa)
    print("Generated bwt", (time.time_ns() - start) / 1000, "micro seconds")
    start = time.time_ns()
    alp = findAlphabet(n)
    print("Generated alphabet", (time.time_ns() - start) / 1000, "micro seconds")
    print()

    start = time.time_ns()
    LS = LSTypes(n)
    print("Found LS types", (time.time_ns() - start) / 1000, "micro seconds")
    start = time.time_ns()
    LMS = findLMSSuffixes(n, LS)
    print("Found LMS", (time.time_ns() - start) / 1000, "micro seconds")
    start = time.time_ns()
    LMSS = findLMSSubstring(n, LMS)
    print("Found LMS substrings", (time.time_ns() - start) / 1000, "micro seconds")
    print()

    start = time.time_ns()
    exactSearch(n, k)
    print("Searching naive:", (time.time_ns() - start) / 1000, "micro seconds")
    print()

    start = time.time_ns()
    cTable = genCTable(n, alp)
    print("C table generation", (time.time_ns() - start) / 1000, "micro seconds")
    oTableTime = time.time_ns()
    oInd, oTable = genOTable(bwt, alp)
    print("O table generation", (time.time_ns() - oTableTime) / 1000, "micro seconds")
    searchTime = time.time_ns()
    lri = initBwtSearchIter(n, k, cTable, oInd, oTable)
    print("Searching", (time.time_ns() - searchTime) / 1000, "micro seconds")
    print("Total search time", (time.time_ns() - start) / 1000, "micro seconds")
    print()

    start = time.time_ns()
    matches = getMatchIndices(sa, lri)
    print("matching", (time.time_ns() - start) / 1000, "micro seconds")

def timeTestSeveralK(n, k):
    start = time.time_ns()
    s = generateSuffixes(n)
    print("Suffixes done", (time.time_ns() - start) / 1000, "micro seconds")
    start = time.time_ns()
    sa = generateSuffixArray(s)
    print("Sorting done", (time.time_ns() - start) / 1000, "micro seconds")
    print()

    start = time.time_ns()
    bwt = generateBWT(s, sa)
    print("Generated bwt", (time.time_ns() - start) / 1000, "micro seconds")
    start = time.time_ns()
    alp = findAlphabet(n)
    print("Generated alphabet", (time.time_ns() - start) / 1000, "micro seconds")
    print()

    start = time.time_ns()
    LS = LSTypes(n)
    print("Found LS types", (time.time_ns() - start) / 1000, "micro seconds")
    start = time.time_ns()
    LMS = findLMSSuffixes(n, LS)
    print("Found LMS", (time.time_ns() - start) / 1000, "micro seconds")
    start = time.time_ns()
    LMSS = findLMSSubstring(n, LMS)
    print("Found LMS substrings", (time.time_ns() - start) / 1000, "micro seconds")
    print()

    start = time.time_ns()
    for i in k:
        exactSearch(n, k)
    print("Searching naive:", (time.time_ns() - start) / 1000, "micro seconds")
    print()

    start = time.time_ns()
    cTable = genCTable(n, alp)
    print("C table generation", (time.time_ns() - start) / 1000, "micro seconds")
    oTableTime = time.time_ns()
    oInd, oTable = genOTable(bwt, alp)
    print("O table generation", (time.time_ns() - oTableTime) / 1000, "micro seconds")
    searchTime = time.time_ns()
    for i in k:
        lri = initBwtSearchIter(n, i, cTable, oInd, oTable)
    print("Searching", (time.time_ns() - searchTime) / 1000, "micro seconds")
    print("Total search time", (time.time_ns() - start) / 1000, "micro seconds")
    print()

    start = time.time_ns()
    matches = getMatchIndices(sa, lri)
    print("matching", (time.time_ns() - start) / 1000, "micro seconds")

#testGoogol()
testMississippi()
#testRandomNucleotideString(100, 10)
#testRandomNucleotideStringMoreK(10000, 5, 100)
#testABBCABA()