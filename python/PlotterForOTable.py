import matplotlib.pyplot as plt

print("Input format: OTable <Time ms>, Number of charters <Size>")
file = open("../DATA/TimeOTable.txt")
file2 = open("../DATA/TimeCTable.txt")
timeOTable = []
timeSecOTable = []
timeCTable = []
timeSecCTable = []
size = []
data = file.readlines()
data2 = file2.readlines()


for line in data:
    timeOTable.append(float(line.split()[1]))
    size.append(float(line.split()[3]))

i = 0
for line in data2:

    if i >14:
        break
    timeCTable.append(float(line.split()[1]))
    i += 1

for time in timeOTable:
    timeSecOTable.append(time / 1000)

for time in timeCTable:
    timeSecCTable.append(time / 1000)

print(timeSecCTable)
#print(timeSecOTable)
plt.plot(size, timeSecOTable, label='O Table')
plt.plot(size, timeSecCTable, label='C Table')


plt.ylabel('Time in Second')
plt.xlabel('Size of input x Billions')
plt.title('Time to generate C and O tables')

plt.legend()
plt.show()
