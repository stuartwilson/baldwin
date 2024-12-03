import numpy as np
import pylab as pl
import json

with open('test.json','r') as file:
	data = json.load(file)

plastic = data['plastic']
fitness = data['fitness']

F = pl.figure()
f = F.add_subplot(111)
f.plot(plastic)

pl.show()
