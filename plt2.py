import numpy as np
import pylab as pl
import json

with open('test.json','r') as file:
	data = json.load(file)

plastic = data['plastic']
fitness = data['fitness']

F = pl.figure(figsize=(12,4))

f = F.add_subplot(121)
f.plot(fitness)
f.set_xlabel('generations')
f.set_ylabel('fitness')

f = F.add_subplot(122)
f.plot(plastic)
f.set_xlabel('generations')
f.set_ylabel('plastic')

pl.show()
