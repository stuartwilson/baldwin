import numpy as np
import pylab as pl
import json
import sys
import os

#files = ['data/test1.json','data/test2.json','data/test3.json']

folder = sys.argv[1]
files = os.listdir(folder)

jsonFiles = []
for file in files:
    if file.split('.')[1]=="json":
        jsonFiles.append(file)


F = pl.figure(figsize=(12,4))
f1 = F.add_subplot(121)
f2 = F.add_subplot(122)

legend = []

for jsonFile in jsonFiles:
    with open(folder+jsonFile,'r') as file:
	    data = json.load(file)

    plastic = data['plastic']
    fitness = data['fitness']

    legend.append(jsonFile)

    #legend.append("genes="+str(data['genes'])+", unique="+str(data['numUnique']))

    f1.plot(fitness)
    f2.plot(plastic)



f1.set_xlabel('generations')
f1.set_ylabel('fitness')
f1.legend(legend)

f2.set_xlabel('generations')
f2.set_ylabel('plastic')

pl.show()
