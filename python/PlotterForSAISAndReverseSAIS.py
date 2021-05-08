import matplotlib.pyplot as plt

print("Input format: SAIS <Time ms>, Reverse SAIS <Time ms>, Number of charters <Size>")
file = open("../TimeOptimizedSAIS.txt")
data = []
SAIS = []
reverseSAIS = []
size = []
data = file.readlines()

for line in data:
    SAIS.append(int(line.split()[1]))
    reverseSAIS.append(int(line.split()[3]))
    size.append(int(line.split()[5]))

plt.plot(size, reverseSAIS, label='Reverse SAIS')
plt.plot(size, SAIS, label='SAIS')
print(SAIS)

plt.ylabel('Time in ms')
plt.xlabel('Size of input')
plt.title('Plot for SAIS and Reverse SAIS')

plt.legend()
plt.show()

print(data)
