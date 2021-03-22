import random

# Generates a random string of A, C, G, T

def generateString(length):
    DNAString = ""
    for i in range (length):
        r = random.randint(0, 3)
        if r == 0:
            DNAString += 'A'
        if r == 1:
            DNAString += 'C'
        if r == 2:
            DNAString += 'G'
        if r == 3:
            DNAString += 'T'

    return DNAString
