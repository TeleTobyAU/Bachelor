import matplotlib.pyplot as plt

print("Input format: ReverseSAIS <Time ms>, Number of charters <Size>")
file = open("../DATA/TimeReverseSAIS.txt")
data = []
reverseSais = []
size = []
data = file.readlines()

for line in data:
    reverseSais.append(int(line.split()[1]))
    size.append(int(line.split()[3]))

plt.plot(size, reverseSais, label='Reverse SAIS')

plt.ylabel('Times in seconds')
plt.xlabel('Size of input string x 1000')
plt.title('Plot for Reverse SAIS')

plt.legend()
plt.show()
