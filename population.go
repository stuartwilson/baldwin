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
	n     int
}

func NewProbabilitySelector(x []float64) *ProbabilitySelector {
	n := len(x)
	lower := make([]float64, n)
	upper := make([]float64, n)
	cumSum := 0.0
	for i := 0; i < n; i++ {
		lower[i] = cumSum
		cumSum += x[i]
		upper[i] = cumSum
	}
	for i := 0; i < n; i++ {
		lower[i] /= cumSum
		upper[i] /= cumSum
	}
	return &ProbabilitySelector{lower, upper, n}
}

func (s *ProbabilitySelector) Select() int {
	r := rand.Float64()
	for i := 0; i < s.n; i++ {
		if s.Lower[i] <= r && r < s.Upper[i] {
			return i
		}
	}
	return s.n - 1
}

type IndividualI interface {
	GetGenome() []int
	SetGenome([]int)
	ComputeFitness([]int) //float64
	GetFitness() float64
}

func NewGenome(n int, p *ProbabilitySelector) []int {
	x := make([]int, n)
	for i := 0; i < n; i++ {
		x[i] = p.Select()
	}
	return x
}

func Combine(mum, dad []int, crossover int) []int {
	return append(mum[:crossover], dad[crossover:]...)
}

type Population []IndividualI

func Evolve(initial Population, generations int, target []int, nUnstable int) ([]float64, []float64) {

	// store the number of units of each type
	countPlastic := make([]float64, generations)
	meanFitness := make([]float64, generations)

	P := len(initial)
	N := len(initial[0].GetGenome())
	minFitness := 1.0 / float64(N)
	unstable := rand.Perm(N)[:nUnstable]

	current := initial
	for g := 0; g < generations; g++ {

		// evaluate fitness of each individual
		F := make([]float64, P)
		for i := 0; i < P; i++ {
			go current[i].ComputeFitness(target)
		}
		for i := 0; i < P; i++ {
			F[i] = minFitness + (1-minFitness)*current[i].GetFitness()
		}
		// compute next generation
		selector := NewProbabilitySelector(F)
		next := make([][]int, P)
		for i := 0; i < P; i++ {
			next[i] = Combine(current[selector.Select()].GetGenome(), current[selector.Select()].GetGenome(), 1+int(rand.Float64()*float64(N-2)))
		}
		// switch populations
		for i := 0; i < P; i++ {
			current[i].SetGenome(next[i])
		}
		// randomly assign target state of unstable units
		for _, i := range unstable {
			if rand.Float64() < 0.5 {
				target[i] = 0
			} else {
				target[i] = 1
			}
		}

		// store results for analysis
		mf := 0.0
		for i := 0; i < P; i++ {
			x := current[i].GetGenome()
			for j := 0; j < len(x); j++ {
				if x[j] == 2 {
					countPlastic[g]++
				}
			}
			mf += current[i].GetFitness()
		}
		meanFitness[g] = mf / float64(P)
		countPlastic[g] /= float64(P) * float64(N)
		fmt.Println(g, meanFitness[g], countPlastic[g])
	}

	return countPlastic, meanFitness
}
