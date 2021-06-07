import matplotlib.pyplot as plt

print("Input format: OTable <Time ms>, Number of charters <Size>")
file = open("../DATA/TimeRecApprox2.txt")
data = file.readlines()
size = []
edits0 = []
edits1 = []
edits2 = []
edits3 = []

i = 0
for line in data:
    i += 1
    if 2 < i < 23:
        edits0.append(int(line.split()[1]))

    if 24 < i < 45:
        edits1.append(int(line.split()[1]))
    if 46 < i < 67:
        edits2.append(int(line.split()[1]))
    if 68 < i:
        edits3.append(int(line.split()[1]))
j = 10
for i in range(20):
    size.append(j)
    j += 10


plt.plot(size, edits3, label='Edits: 3')
plt.plot(size, edits2, label='Edits: 2')
plt.plot(size, edits1, label='Edits: 1')
plt.plot(size, edits0, label='Edits: 0')




plt.ylabel('Time in milliseconds')
plt.xlabel('Key length')
plt.title('Complexity of key length rec approx')

plt.legend()
plt.show()
