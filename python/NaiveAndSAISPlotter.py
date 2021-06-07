import matplotlib.pyplot as plt

print("Input format: SAIS <Time ms>, Reverse SAIS <Time ms>, Number of charters <Size>")
file = open("../DATA/TimeNaiveAndSAIS.txt")
data = []
SAIS = []
naiveSA = []
size = []
data = file.readlines()

i = 0
for line in data:
    if i < 100:
        naiveSA.append(int(line.split()[3]))
        size.append(int(line.split()[5]))
    else:
        SAIS.append(int(line.split()[1]))
    i += 1


newSAIS = []
newNaiveSA = []
for x in range(len(SAIS)):
    newSAIS.append(SAIS[x] / 1000)
    newNaiveSA.append(naiveSA[x] / 1000)

plt.plot(size, newNaiveSA, label='Naive SA')
plt.plot(size, newSAIS, label='SAIS')
print(SAIS)

plt.ylabel('Time in seconds')
plt.xlabel('Size of input')
plt.title('Plot for SAIS and Naive SA')

plt.legend()
plt.show()

print(data)