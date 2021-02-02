# Takes a string input from terminal and generates corresponding suffix array
# Set printResult to true if you want to print
printResult = True
# Sentinel is $

def generateSuffixArray(inputString):
    workingString = inputString + '$'

    suffixArray = []
    for i in range (len(workingString)):
        suffixArray.append(workingString[i : len(workingString)] + workingString[0 : i])

    if printResult:
        print("\nSuffix array:")
        print(suffixArray)

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
        print(sortedSuffixArray)
        print("\nBurrows Wheeler transformation string:")
        print(bwt)

#inp = input("Type string you wish to create suffix array from: \n")
#suffixArray = generateSuffixArray(inp)
#generateBWT(suffixArray)
