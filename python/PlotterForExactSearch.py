import matplotlib.pyplot as plt

print("Input format: Exact Search <Time ms>, Number of charters <Size>")
file = open("../DATA/TimeExactMatch.txt")
data = []
timeExact = []
size = []
data = file.readlines()

for line in data:
    timeExact.append(float(line.split()[1]))
    size.append(float(line.split()[3]))

plt.plot(size, timeExact, label='Exact Search')

plt.ylabel('Time in Millisecond')
plt.xlabel('Size of key')
plt.title('Run time to find Exact Match')

plt.legend()
plt.show()

print(data)
