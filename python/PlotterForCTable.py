import matplotlib.pyplot as plt

print("Input format: CTable <Time ms>, Number of charters <Size>")
file = open("../DATA/TimeCTable.txt")
data = []
timeCTable = []
size = []
data = file.readlines()

for line in data:
    timeCTable.append(int(line.split()[1]))
    size.append(int(line.split()[3]))


plt.plot(size, timeCTable, label='C Table')

plt.ylabel('Time in Millisecond')
plt.xlabel('Size of input x 1000')
plt.title('Time to create C table')

plt.legend()
plt.show()

print(data)
