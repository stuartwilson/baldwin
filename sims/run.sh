#!/bin/bash

dir="data/funky/"
mkdir $dir
rm ${dir}*
cp run.sh ${dir}run_copy.sh

type="Nowlan"
n=5
pop=1000
gens=500
#trials=500
unstable=0
plastic=0.5

i=2
for trials in 2
do
	go run main.go ${dir}sim_$i.json $type $n $pop $gens $trials $unstable $plastic $i 0
	i=$(($i+1))
done

python plt.py $dir
