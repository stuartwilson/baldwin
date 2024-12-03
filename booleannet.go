package main

import (
	"math/rand"
)

type GRN struct {
	Nnodes        int
	Ncombinations int
	Nbits         int
	State         []bool
	Genome        []int
	fitness       float64
	trials        int
}

func NewGRN(n int, p *ProbabilitySelector, trials int) *GRN {
	N := int(pow(2, n))
	grn := GRN{
		Nnodes:        n,
		Ncombinations: N,
		Nbits:         n * N,
		State:         make([]bool, n),
		Genome:        NewGenome(n*N, p),
		trials:        trials,
	}
	return &grn
}

func (grn *GRN) GetCycle(initialState []bool) [][]bool {
	genome := make([]bool, grn.Nbits)
	for i := 0; i < grn.Nbits; i++ {
		if grn.Genome[i] > 0 {
			genome[i] = true
			if grn.Genome[i] == 2 && rand.Float64() < 0.5 {
				genome[i] = false
			}
		}
	}
	grn.State = initialState
	visited := make([][]bool, 0)
	for t := 0; t < grn.Ncombinations; t++ {
		visited = append(visited, grn.State)
		grn.Step(genome)
		for v := 0; v < len(visited); v++ {
			if match(grn.State, visited[v]) { // limit cycle detected
				return visited[v:]
			}
		}
	}
	return nil
}

func (grn *GRN) Step(genome []bool) {
	stateIndex := 0
	power := 1
	for i := 0; i < grn.Nnodes; i++ {
		if grn.State[i] {
			stateIndex += power
		}
		power *= 2
	}
	tableIndex := 0
	for i := 0; i < grn.Nnodes; i++ {
		grn.State[i] = genome[tableIndex+stateIndex]
		tableIndex += grn.Ncombinations
	}
}

func (grn *GRN) GetGenome() []int {
	return grn.Genome
}

func (grn *GRN) SetGenome(genome []int) {
	grn.Genome = genome
}

func (grn *GRN) GetFitness() float64 {
	return grn.fitness
}

func (grn *GRN) ComputeFitness(target []int) {
	for t := 0; t < grn.trials; t++ {
		cycle := grn.GetCycle(make([]bool, grn.Nnodes))
		if len(cycle) == 1 {
			targ := make([]bool, len(target))
			for i := 0; i < len(target); i++ {
				targ[i] = target[i] > 0
			}
			if match(cycle[0], targ) {
				grn.fitness = 1.0
				return
			}
		}
	}
	grn.fitness = 0.0
}
