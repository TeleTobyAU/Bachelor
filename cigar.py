s = "MMIIMMM"

curr = "-1"
counter = 0
prString = ""
for i in s:
    if i != curr:
        prString += curr + str(counter)
        curr = i
        counter = 1
    else:
        counter += 1

prString += curr + str(counter)

print("Input:", s)
print("Cigar:", prString[3:])