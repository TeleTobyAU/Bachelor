import matplotlib.pyplot as plt

print("Input format: OTable <Time ms>, Number of charters <Size>")
file = open("../DATA/TimeSAIS.txt")
data = []
sais = []
size = []
data = file.readlines()

for line in data:
    sais.append(int(line.split()[1]))
    size.append(int(line.split()[3]))

plt.plot(size, sais, label='SAIS')

plt.ylabel('Times in seconds')
plt.xlabel('Size of input string x 1000')
plt.title('Plot for SAIS')

plt.legend()
plt.show()
