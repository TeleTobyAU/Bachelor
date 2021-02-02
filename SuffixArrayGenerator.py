# Takes a string input from terminal and generates corresponding suffix array
# Set printResult to true if you want to print
printResult = True
# Sentinel is $

def generateSuffixArray(inputString):
    workingString = inputString + '$'

    suffixArray = []
    for i in range (len(workingString)):
        suffixArray.append(workingString[i : len(workingString)] + workingString[0 : i])

    sortedSuffixArray = suffixArray.copy()
    sortedSuffixArray.sort()

    bwt = ""
    for i in range (len(sortedSuffixArray)):
        item = sortedSuffixArray[i]
        bwt += item[-1]

    if printResult:
        print("Suffix array:")
        print(suffixArray)
        print("Sorted suffix array:")
        print(sortedSuffixArray)
        print("Burrows Wheeler transformation string:")
        print(bwt)

    return suffixArray

inp = input("Type string you wish to create suffix array from: \n")
generateSuffixArray(inp)
