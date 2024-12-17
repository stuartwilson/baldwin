package baldwin

import (
	"fmt"
	"math/rand"
	"strconv"
)

type GRN struct {
	Nnodes        int
	Ncombinations int
	Nbits         int
	State         []bool
	Genome        []int
	fitness       float64
	trials        int
	switchCase    int
}

func NewGRN(n int, p *ProbabilitySelector, trials int, extra []string) *GRN {

	var err error
	switchCase := 0 // default switch case is 1
	if len(extra) > 0 {
		switchCase, err = strconv.Atoi(extra[0])
		if err != nil {
			fmt.Println("Switch case must be an integer")
			return nil
		}
	}

	N := int(pow(2, n))
	grn := GRN{
		Nnodes:        n,
		Ncombinations: N,
		Nbits:         n * N,
		State:         make([]bool, n),
		Genome:        NewGenome(n*N, p),
		trials:        trials,
		switchCase:    switchCase,
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

/*
func (grn *GRN) ComputeFitness(target []int) {

	var initial []bool
	switch grn.switchCase {
	case 0:
		initial = make([]bool, grn.Nnodes)
	default:
		fmt.Println("invalid switch case")
		return
	}

	for t := 0; t < grn.trials; t++ {
		cycle := grn.GetCycle(initial)
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
*/

// HARD-CODING THIS ONE TO EVALUATE TWO
func (grn *GRN) ComputeFitness(target []int) {

	initial1 := make([]bool, grn.Nnodes)
	initial2 := make([]bool, grn.Nnodes)
	initial2[0] = true
	target1 := make([]bool, grn.Nnodes)
	target2 := make([]bool, grn.Nnodes)
	target2[0] = true
	for i := 1; i < grn.Nnodes; i++ {
		target1[i] = !target1[i-1]
		target2[i] = !target2[i-1]
	}

	for t := 0; t < grn.trials; t++ {

		cycle1 := grn.GetCycle(initial1)
		if len(cycle1) == 1 {
			if match(cycle1[0], target1) {
				cycle2 := grn.GetCycle(initial2)
				if len(cycle2) == 1 {
					if match(cycle2[0], target2) {
						grn.fitness = 1.0
						//grn.fitness = float64(grn.trials-t-1) / float64(grn.trials-1)
						return
					}
				}
			}
		}

	}
	grn.fitness = 0.0
}
