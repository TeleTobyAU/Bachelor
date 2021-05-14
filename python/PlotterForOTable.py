import matplotlib.pyplot as plt

print("Input format: OTable <Time ms>, Number of charters <Size>")
file = open("../DATA/TimeOTable.txt")
data = []
timeOTable = []
size = []
data = file.readlines()

for line in data:
    timeOTable.append(int(line.split()[1]))
    size.append(int(line.split()[3]))


plt.plot(size, timeOTable, label='O Table')

plt.ylabel('Time in Millisecond')
plt.xlabel('Size of input x 1000')
plt.title('Time to create O table')

plt.legend()
plt.show()
