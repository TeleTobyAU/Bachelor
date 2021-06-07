import matplotlib.pyplot as plt

print("Input format: CTable <Time ms>, Number of charters <Size>")
file = open("../DATA/TimeCTable.txt")
data = []
timeCTable = []
timeSecCTable = []
size = []
data = file.readlines()

for line in data:
    timeCTable.append(float(line.split()[1]))
    size.append(float(line.split()[3]))

for time in timeCTable:
    timeSecCTable.append(float(time) / 1000.0)


plt.plot(size, timeSecCTable, label='C Table')

plt.ylabel('Time in Second')
plt.xlabel('Size of input x in billions')
plt.title('Time to generate C table')

plt.legend()
plt.show()

print(data)
