package main

import (
	"fmt"
	"math/rand"
)

func pow(x float64, exp int) float64 {
	y := 1.0
	for i := 0; i < exp; i++ {
		y *= x
	}
	return y
}

type ProbabilitySelector struct {
	Lower []float64
	Upper []float64
}

func NewProbabilitySelector(x []float64) *ProbabilitySelector {
	lower := make([]float64, len(x))
	upper := make([]float64, len(x))
	cumSum := 0.0
	for i := 0; i < len(x); i++ {
		lower[i] = cumSum
		cumSum += x[i]
		upper[i] = cumSum
	}
	for i := 0; i < len(x); i++ {
		lower[i] /= cumSum
		upper[i] /= cumSum
	}
	return &ProbabilitySelector{lower, upper}
}

func (s *ProbabilitySelector) Select() int {
	r := rand.Float64()
	for i := 0; i < len(s.Lower); i++ {
		if s.Lower[i] <= r && r < s.Upper[i] {
			return i
		}
	}
	return len(s.Lower) - 1
}

// define fitness function
func getFit(x, target []int, trials int) bool {
	a := 0
	for i := 0; i < len(x); i++ {
		if x[i] == 2 {
			a++
		} else {
			if x[i] != target[i] {
				return false
			}
		}
	}
	return rand.Float64() < 1.0-pow((1.0-pow(0.5, a)), trials)
}

type IndividualI interface {
	GetGenome() *[]int
	ComputeFitness([]int)
	GetFitness() float64
}

func NewGenome(n int, p *ProbabilitySelector) *[]int {
	x := make([]int, n)
	for i := 0; i < n; i++ {
		x[i] = p.Select()
	}
	return &x
}

func Combine(mum, dad IndividualI, crossover int) IndividualI {
	m := mum.GetGenome()
	d := dad.GetGenome()
	child := append((*m)[:crossover], (*d)[crossover:]...)
	m = &child
	return mum
}

type Population []IndividualI

func Evolve(initial Population, generations int, target []int, nUnstable int) ([]float64, []float64) {

	P := len(initial)
	N := len(*initial[0].GetGenome())
	minFitness := 1.0 / float64(N)

	// store the number of units of each type
	countPlastic := make([]float64, generations)
	meanFitness := make([]float64, generations)
	//countCorrect := make([]int, generations)

	unstable := rand.Perm(N)[:nUnstable]

	// evolution loop
	current := initial
	for g := 0; g < generations; g++ {
		fmt.Println(g)

		/*
			//count units(for plotting)
			for _, i := range current {
				x := *i.GetGenome()
				for j := 0; j < len(x); j++ {
					if x[j] == 2 {
						countPlastic[g]++
					} else {
						if x[j] == target[j] {
							countCorrect[g]++
						}
					}
				}
			}
			countPlastic[g] /= P
			countCorrect[g] /= P
		*/

		mf := 0.0
		for _, i := range current {
			x := *i.GetGenome()
			for j := 0; j < len(x); j++ {
				if x[j] == 2 {
					countPlastic[g]++
				}
			}
			mf += i.GetFitness()
		}
		meanFitness[g] = mf / float64(P)
		countPlastic[g] /= float64(P)

		// evaluate fitness of each individual
		F := make([]float64, P)
		for i := 0; i < P; i++ {
			current[i].ComputeFitness(target)
			F[i] = minFitness + (1-minFitness)*current[i].GetFitness()
		}
		selector := NewProbabilitySelector(F)

		next := make(Population, P)
		for i := 0; i < P; i++ {
			mum := current[selector.Select()]
			dad := current[selector.Select()]
			crossover := 1 + int(rand.Float64()*float64(N-2))
			next[i] = Combine(mum, dad, crossover)
		}
		current = next

		// randomly assign target state of unstable units
		for _, i := range unstable {
			if rand.Float64() < 0.5 {
				target[i] = 0
			} else {
				target[i] = 1
			}
		}
	}
	return countPlastic, meanFitness
}
