import matplotlib.pyplot as plt

print("Input format: SAIS <Time ms>, Reverse SAIS <Time ms>, Number of charters <Size>")
file = open("../DATA/TimeOptimizedSAIS.txt")
data = []
SAIS1 = []
reverseSAIS1 = []
size = []
data = file.readlines()

for line in data:
    SAIS1.append(int(line.split()[1]))
    reverseSAIS1.append(int(line.split()[3]))
    size.append(int(line.split()[5]))
SAIS = [x / 1000 for x in SAIS1]
reverseSAIS = [x / 1000 for x in reverseSAIS1]

plt.plot(size, reverseSAIS, label='Reverse SAIS')
plt.plot(size, SAIS, label='SAIS')
print(SAIS)

plt.ylabel('Time in second')
plt.xlabel('Size of input')
plt.title('Plot for SAIS and Reverse SAIS')

plt.legend()
plt.show()

print(data)
