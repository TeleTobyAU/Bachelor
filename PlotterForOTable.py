import matplotlib.pyplot as plt

print("Input format: OTable <Time ms>, Number of charters <Size>")
file = open("TimeDataOTable.txt")
data = []
OTable = []
size = []
data = file.readlines()

for line in data:
    OTable.append(int(line.split()[1]))
    size.append(int(line.split()[3]))

plt.plot(size, OTable, label='OTable')

plt.ylabel('Time in ms')
plt.xlabel('Size of input')
plt.title('Plot for OTable')

plt.legend()
plt.show()
