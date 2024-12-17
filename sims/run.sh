#!/bin/bash

dir="data/grn1/"
mkdir $dir
rm ${dir}*
cp run.sh ${dir}run_copy.sh

type="GRN"
n=5
pop=1000
gens=500
#trials=500
unstable=0
plastic=0.25

i=1
for trials in 50 100 200 500
do
	go run main.go ${dir}sim_$i.json $type $n $pop $gens $trials $unstable $plastic $i 0
	i=$(($i+1))
done

python plt.py $dir
