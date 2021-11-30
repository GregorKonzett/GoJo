import subprocess

amount = 5

producerAmounts = 4
vals = 100


print("Junctions:")
for i in range(producerAmounts):
    prodConAmount = 10 ** i
    print("Producer/Consumer Amount: " + str(prodConAmount))
    for j in range(amount):
        subprocess.call(["./go_build_queuebenchmarking_go.exe","junction",str(prodConAmount), str(prodConAmount), str(vals)])

print("Mutex")
for i in range(producerAmounts):
    prodConAmount = 10 ** i
    print("Producer/Consumer Amount: " + str(prodConAmount))
    for j in range(amount):
        subprocess.call(["./go_build_queuebenchmarking_go.exe","mutex",str(prodConAmount), str(prodConAmount), str(vals)])