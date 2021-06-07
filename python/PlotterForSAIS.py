import matplotlib.pyplot as plt

print("Input format: SAIS <Time s>, Number of charters <Size>")
file = open("../DATA/TimeOptimizedSAIS.txt")
data = file.readlines()
file2 = open("../DATA/TimeOptimizedSAISv2.txt")
data2 = file2.readlines()
sais = []
saisV2 = []
size = []
sizeV2 = []

for line in data:
    sais.append(int(line.split()[1]))
    size.append(float(line.split()[3]))

for line in data2:
    saisV2.append(int(line.split()[1]))
    sizeV2.append(float(line.split()[3]))

for x in range(len(size)):
    sais[x] = sais[x] / 60

for x in range(len(sizeV2)):
    saisV2[x] = saisV2[x] / 60
print(sizeV2)

plt.plot(size, sais, label='SA-IS V1')
plt.plot(sizeV2, saisV2, label='SA-IS V2')

plt.ylabel('Times in minutes')
plt.xlabel('Size of input string billions')
plt.title('Plot for SA-IS')

plt.legend()
plt.show()
