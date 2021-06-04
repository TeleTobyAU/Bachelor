import RandomRNAStringGenerator
import time


#The BWT table struct holds all of our important information
class BWTTable:
    n = None
    alphabet = None
    sa = None
    bwt = None
    rsa = None
    cTable = None
    oIndex = None
    oTable = None

    def __init__(self, n):
        if n == None:
            raise Exception("Input string is required!")
        if type(n) != type(str()):
            raise Exception("Input has to be a string!")

        self.n = "mississippi$"

        self.findAlphabet()

        self.genCTable()

        self.generateSuffixArray()
        self.printStuff()


    #Naive solutions ---------------------------------------------------------------------------------------------------
    def printStuff(self):
        print(self.n, "\n--------------------")
        suffixes = []
        for i in range(len(self.n)):
            suffixes.append(self.n[i:] + self.n[:i])
            if i < 10:
                print(i, " " + suffixes[i])
            else:
                print(i, suffixes[i])

        print("--------------------")
        suffixArray = sorted(suffixes)
        for i in range(len(suffixArray)):
            if self.sa[i] < 10:
                print(self.sa[i], " " + suffixArray[i])
            else:
                print(self.sa[i], suffixArray[i])

        print("--------------------")
        suffixArray = sorted(suffixes)
        for i in range(len(suffixArray)):
            if self.sa[i] < 10:
                print(self.sa[i], " " + suffixArray[i][:-1]+ " " + suffixArray[i][-1])
            else:
                print(self.sa[i], suffixArray[i][:-1] + " " + suffixArray[i][-1])
        print("--------------------")
        print("BWT:", self.generateBWT(suffixes, self.sa))
        print("--------------------")

        print(self.n + "$")
        n = self.n + "$"
        inString = n[::-1]
        LSTypes = "S"
        for i in range(1, len(n)):
            if inString[i - 1] == inString[i]:
                LSTypes += LSTypes[i - 1]
            else:
                if inString[i - 1] < inString[i]:
                    LSTypes += "L"
                else:
                    LSTypes += "S"
        LSTypes = LSTypes[::-1]
        print(LSTypes)

        LMSIndices = []
        if LSTypes[0] == "S": LMSIndices.append(0)
        for i in range(len(LSTypes)):
            if LSTypes[i] == "S" and LSTypes[i - 1] != "S":
                LMSIndices.append(i)

        printer = ""
        for i in range(len(LSTypes)):
            if i in LMSIndices:
                printer += "*"
            else:
                printer += " "
        print(printer)

        printer = ""
        for i in range(len(LSTypes)):
            if i in LMSIndices:
                printer += "|"
            else:
                printer += "-"
        print(printer)

        print("--------------------\n")

        for i in range(len(self.alphabet)):
            print(self.alphabet[i], " ", self.cTable[i])

        print("-------------------------------------------------")

        self.genOTable()
        bwt = self.generateBWT(suffixes, self.sa)
        printer = "      "
        for i in bwt:
            printer += i + "  "
        print(printer)

        for i in range(len(self.alphabet)):
            print(self.alphabet[i], self.oTable[i])

        print("-------------------------------------------------------")

        self.generateSuffixArray(True)
        self.genOTable(True)
        bwt = self.generateBWT(suffixes, self.rsa)
        printer = "      "
        for i in bwt:
            printer += i + "  "
        print(printer)

        for i in range(len(self.alphabet)):
            print(self.alphabet[i], self.roTable[i])

        print("-------------------------------------------------------")



    def generateSuffixArray(self, r=False):
        if r:
            n = self.n[::-1]
        else:
            n = self.n

        suffixes = []
        for i in range(len(n)):
            suffixes.append(n[i: len(n)] + n[0: i])

        suffixArray = []
        for s in sorted(suffixes):
            suffixArray.append(suffixes.index(s))

        if r:
            self.rsa = suffixArray
        else:
            self.sa = suffixArray


    def generateBWT(self, suffixes, suffixArray):
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


    def findAlphabet(self):
        alphabet = []
        for i in self.n:
            if i not in alphabet: alphabet.append(i)
        alphabet.sort()
        self.alphabet = alphabet


    def BWTLookup(self, i, r=False):
        if r:
            n = self.n[::-1]
            nIndex = self.rsa[i]
        else:
            n = self.n
            nIndex = self.sa[i]

        if nIndex == 0:
            return n[-1]
        else:
            return n[nIndex - 1]



    #Burrows Wheeler transformation search -----------------------------------------------------------------------------
    def genCTable (self):
        cTable = []
        for i in self.alphabet:
            cTable.append(0)
            for j in self.n:
                if i > j: cTable[self.alphabet.index(i)] += 1

        self.cTable = cTable


    def genOTable(self, r=False):
        if r:
            n = self.n[::-1]
        else:
            n = self.n

        oIndices = self.alphabet.copy()
        oTable = []

        for i in range(len(oIndices)):
            oTable.append([0])

        for i in range(len(n)):
            for j in range(len(oIndices)):
                if self.BWTLookup(i, r) == self.alphabet[j]:
                    oTable[j].append(oTable[j][i] + 1)
                else:
                    oTable[j].append(oTable[j][i])

        if r:
            self.roIndices = oIndices
            self.roTable = oTable
        else:
            self.oIndices = oIndices
            self.oTable = oTable


    def initBwtSearchIter(self, k):
        L = 0
        R = len(self.n)

        #Check if the key's characters are in our alphabet
        for i in k:
            if i not in self.alphabet:
                print(i, "from the key is not in the alphabet!")
                exit(-1)

        #If our key is longer than our string there will be no match
        if len(k) > len(self.n):
            R = 0
            L = 1

        i = len(k) - 1

        while i >= 0 and L < R:
            a = self.oIndices[0:].index(k[i])
            L = self.cTable[a] + self.oTable[a][L]
            R = self.cTable[a] + self.oTable[a][R]
            i -= 1

        i = L
        LRI = (L, R, i)
        self.getMatchIndices(LRI)


    def getMatchIndices(self, LRI):
        matches = []
        L = LRI[0]
        R = LRI[1]

        for i in range(R - L):
            matches.append(self.sa[L + i])
        matches.sort()

        print("Matches:", matches)
        print()

    #Approx matching stuff
    def initBWTApproxSearch(self, k, maxEdits):
        #Generate the D table
        dTable = []
        minEdits = 0
        L = 0
        R = len(self.n)

        for i in range(len(k)):
            a = self.roIndices[0:].index(k[i])
            L = self.cTable[a] + self.roTable[a][L]
            R = self.cTable[a] + self.roTable[a][R]

            if L >= R:
                minEdits += 1
                L = 0
                R = len(self.n)
            dTable.append(minEdits)
        self.dTable = dTable

        #Start searching
        L = 0
        R = len(self.n)
        i = len(k) - 1
        matchA = self.alphabet.index(k[i])
        edits = []

        for a in range(1, len(self.alphabet)):
            newL = self.cTable[a] + self.oTable[a][L]
            newR = self.cTable[a] + self.oTable[a][R]

            if a == matchA:
                editCost = 0
            else:
                editCost = 1

            edits.append("M")
            self.recursiveApproxMatch(newL, newR, i - 1, maxEdits - editCost, edits)

            edits.append("I")
            self.recursiveApproxMatch(L, R, i - 1, maxEdits - 1, edits)

            L = len(k)
            R = 0
            nextInterval = 0

        return L, R


    def recursiveApproxMatch(self, L, R, i, editsLeft, edits):
        if i >= 0:
            lowerLim = self.dTable[i]
        else:
            lowerLim = 0

        if editsLeft < lowerLim:
            return

        iva = [[], []]
        if i < 0:
            iva[0].append(L)
            iva[1].append(R)

        revEdits = reversed(edits)
        #editsToCigar(cigar, revEdits) #todo very unsure about this

        return


    #Linear suffix array construction by almost pure induced sorting ---------------------------------------------------
    def SAIS(self, reverse=False):
        if reverse:
            n = self.n[:-1]
            n = n[::-1] + "$"
        else:
            n = self.n

        def classifyLAndSTypes(n):
            inString = n[::-1]
            LSTypes = "S"
            for i in range(1, len(n)):
                 if inString[i - 1] == inString[i]:
                    LSTypes += LSTypes[i - 1]
                 else:
                    if inString[i - 1] < inString[i]:
                        LSTypes += "L"
                    else:
                        LSTypes += "S"
            LSTypes = LSTypes[::-1]

            return LSTypes

        def findLMSIndices(LSTypes):
            LMSIndices = []
            if LSTypes[0] == "S": LMSIndices.append(0)
            for i in range(len(LSTypes)):
                if LSTypes[i] == "S" and LSTypes[i - 1] != "S":
                    LMSIndices.append(i)
            return LMSIndices

        def findBucketBeginnings(self):
            beginnings = [0] * len(self.alphabet)
            for i in range(1, len(self.alphabet)):
                beginnings[i] = self.cTable[i] #Kinda hacky to use the C table already, I know
            return beginnings

        def findBucketEnds(self, n):
            ends = [-1] * len(self.alphabet) #TODO This -1 is weird
            for i in range(len(n) - 1, -1, -1):
                j = self.alphabet.index(n[i]) #The numerical/remapped letter
                for k in range(j, len(self.alphabet)):
                    ends[k] += 1
            return ends

        def getBuckets(self, n):
            beginnings = findBucketBeginnings(self)
            ends = findBucketEnds(self, n)

            buckets = []
            for i in range(len(beginnings)):
                buckets.append([beginnings[i], ends[i]])
            return buckets

        def generateCompressedN(self, n, LMSSubstrings):
            s = [LMSSubstrings.index(n[0: LMSIndices[0] + 1])]
            for i in range(1, len(LMSIndices)):
                s.append(LMSSubstrings.index(n[LMSIndices[i - 1]: LMSIndices[i] + 1]))
            return s

        LSTypes = classifyLAndSTypes(n)
        LMSIndices = findLMSIndices(LSTypes)


        SA = [0] * len(n)
        buckets = getBuckets(self, n)

        #SA-IS Step 1, placing LMS indices in the SA
        for i in range(len(n)):
            if i in LMSIndices:
                remappedIndex = self.alphabet.index(n[i])
                SA[buckets[remappedIndex][1]] = i
                buckets[remappedIndex][1] -= 1

        # SA-IS Step 2, placing L types in the SA
        for i in range(len(n) + 1):
            if i >= len(n): break

            if SA[i] == 0: continue

            j = SA[i] - 1
            if LSTypes[j] == "L":
                remappedIndex = self.alphabet.index(n[j])
                SA[buckets[remappedIndex][0]] = j
                buckets[remappedIndex][0] += 1

        # SA-IS Step 3, placing remaining S types in the SA
        buckets = getBuckets(self, n)
        for i in range(len(n), 0, -1):

            if SA[i - 1] == 0: continue

            j = SA[i - 1] - 1
            if (LSTypes[j] == "S"):
                remappedIndex = self.alphabet.index(n[j])
                SA[buckets[remappedIndex][1]] = j
                buckets[remappedIndex][1] -= 1

        if reverse:
            self.rsa = SA
        else:
            self.sa = SA


    #Pretty printing -------------------------------------------------------------------------------------------------------
    def prettyPrint(self):
        print("\nString we are working on:", self.n)
        print("In the alphabet:", self.alphabet)
        print()

        print("Suffix array:", self.sa)
        print()

        print("C table:")
        for i in range(len(self.cTable)): print(self.alphabet[i], ":", self.cTable[i])
        print()

        print("Otable:")
        printBwt = "      "
        for i in range(len(self.n)):
            printBwt += self.BWTLookup(i) + "  "
        print(printBwt)
        for i in range(len(self.oTable)):
            print(self.oIndices[i], self.oTable[i])
        print("Elements in O table:", len(self.oTable) * len(self.oTable[0]))
        print()

        print("rOtable:")
        printBwt = "      "
        for i in range(len(self.n)):
            printBwt += self.BWTLookup(i, True) + "  "
        print(printBwt)
        for i in range(len(self.roTable)):
            print(self.roIndices[i], self.roTable[i])
        print("Elements in O table:", len(self.roTable) * len(self.roTable[0]))
        print()

        print("Searching")
        self.initBwtSearchIter("AAA")


#Tests -----------------------------------------------------------------------------------------------------------------
def testMississippi():
    n = "mmiissiissiippii$"
    bwtTable = BWTTable(n)
    k = "iis"
    #bwtTable.prettyPrint()

def testSATime():
    n = RandomRNAStringGenerator.generateString(1000000) + "$"
    start = time.time_ns()
    bwt = BWTTable(n)
    bwt.generateSuffixArray()


def testNuc():
    print("Starting at:", time.strftime("%H:%M:%S", time.localtime()))
    n = RandomRNAStringGenerator.generateString(100) + "$"
    bwt = BWTTable(n)
    print("100 done at:", time.strftime("%H:%M:%S", time.localtime()))
    n = RandomRNAStringGenerator.generateString(1000) + "$"
    bwt = BWTTable(n)
    print("1.000 done at:", time.strftime("%H:%M:%S", time.localtime()))
    n = RandomRNAStringGenerator.generateString(10000) + "$"
    bwt = BWTTable(n)
    print("10.000 done at:", time.strftime("%H:%M:%S", time.localtime()))
    n = RandomRNAStringGenerator.generateString(100000) + "$"
    bwt = BWTTable(n)
    print("100.000 done at:", time.strftime("%H:%M:%S", time.localtime()))
    n = RandomRNAStringGenerator.generateString(1000000) + "$"
    bwt = BWTTable(n)
    print("1.000.000 done at:", time.strftime("%H:%M:%S", time.localtime()))
    n = RandomRNAStringGenerator.generateString(10000000) + "$"
    bwt = BWTTable(n)
    print("10.000.000 done at:", time.strftime("%H:%M:%S", time.localtime()))
    n = RandomRNAStringGenerator.generateString(100000000) + "$"
    bwt = BWTTable(n)
    print("100.000.000 done at:", time.strftime("%H:%M:%S", time.localtime()))
    n = RandomRNAStringGenerator.generateString(1000000000) + "$"
    bwt = BWTTable(n)
    print("1.000.000.000 done at:", time.strftime("%H:%M:%S", time.localtime()))

def testGoogol():
    n = "GTCGGTATCGGTGGGCGTGCGCCAACCTGGGCAGAGTTGATTCTTGCTTTCCCGCTCATACTACATCCGGAAGCAGATCCAGGCGACCGGAACCGAGCGC$"
    bwtTable = BWTTable(n)
    print("sais   ", bwtTable.sa)
    print("nielses: 100, 90, 23, 70, 62, 91, 85, 24, 59, 32, 74, 71, 95, 80, 34, 57, 76, 64, 6, 39, 99, 22, 31, 73, 79, 56, 63, 21, 78, 50, 92, 51, 86, 66, 25, 83, 93, 97, 19, 52, 87, 67, 2, 8, 15, 60, 54, 26, 42, 46, 89, 69, 84, 94, 33, 75, 38, 98, 30, 72, 20, 82, 96, 18, 14, 53, 45, 88, 68, 29, 81, 13, 28, 12, 3, 9, 4, 0, 16, 10, 35, 61, 58, 5, 55, 77, 49, 65, 1, 7, 41, 37, 17, 44, 27, 11, 48, 40, 36, 43, 47")
    bwtTable.generateSuffixArray(bwtTable)
    print("naiv   ", bwtTable.sa)
    k = ""

def testABBCABA():
    n = "ABBCABA$"
    k = "AB"

#testGoogol()
#testSATime()
testMississippi()
#testABBCABA()
#testNuc()