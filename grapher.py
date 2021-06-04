import matplotlib.pyplot as plt

#input format: SAIS tid SA tid length size

print("Input format should be: SAIS <time> SA <time> length <size>")
print("Press <enter> twice when done")

lines = []
while True:
    line = input()
    if line:
        lines.append(line)
    else:
        break

sais = []
sa = []
length = []

for i in lines:
    inp = i.split()
    sais.append(int(inp[1]))
    sa.append(int(inp[3]))
    length.append(int(inp[5]))

plt.plot(length, sa, label="Naive SA")
plt.plot(length, sais, label="SA-IS")
plt.xlabel("Sequence length")
plt.ylabel("Time in ms")
plt.title("SAIS running times")
plt.legend()
plt.show()