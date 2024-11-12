package main

import (
	"fmt"
	"math/rand"
)

type Hopfield struct {
	N int
	W [][]float64
	X []bool
	M [][]bool
}

func NewHopfield(n int) *Hopfield {
	w := make([][]float64, n)
	for i := 0; i < n; i++ {
		w[i] = make([]float64, n)
	}
	return &Hopfield{
		N: n,
		W: w,
		X: make([]bool, n),
		M: make([][]bool, 0),
	}
}

func (h *Hopfield) SetState(x []bool) {
	h.X = x
}

func (h *Hopfield) Relax() (int, int) {
	repeated := false
	for !repeated {
		repeated = h.Step()
	}
	comparison := h.Compare()
	maxVal := 0
	maxInd := 0
	for i := 0; i < len(comparison); i++ {
		if comparison[i] > maxVal {
			maxVal = comparison[i]
			maxInd = i
		}
	}
	return maxVal, maxInd
}

func (h *Hopfield) Step() bool {
	next := make([]bool, h.N)
	allSame := true
	for i := 0; i < h.N; i++ {
		sum := 0.0
		for j := 0; j < h.N; j++ {
			if i != j {
				if h.X[j] {
					sum += h.W[i][j]
				} else {
					sum -= h.W[i][j]
				}
			}
		}
		next[i] = sum > 0
		if next[i] != h.X[i] {
			allSame = false
		}
	}
	h.X = next
	return allSame
}

func (h *Hopfield) Remember(m []bool) {
	h.M = append(h.M, m)
	for i := 0; i < h.N; i++ {
		for j := 0; j < h.N; j++ {
			if i != j {
				if m[i] == m[j] {
					h.W[i][j]++
				} else {
					h.W[i][j]--
				}
			}
		}
	}
}

func (h *Hopfield) Compare() []int {
	same := make([]int, len(h.M))
	for i := 0; i < len(h.M); i++ {
		for j := 0; j < h.N; j++ {
			if h.X[j] == h.M[i][j] {
				same[i]++
			}
		}
	}
	return same
}

func (h *Hopfield) Evaluate(genome []int, target []bool, trials int) (bool, int, error) {
	if len(genome) != (h.N*h.N-h.N)/2 {
		return false, 0, fmt.Errorf("wrong genome length")
	}

	for t := 0; t < trials; t++ {

		// initial random state
		for i := 0; i < h.N; i++ {
			h.X[i] = rand.Float64() < 0.5
		}

		// set Hopfield weights from genome
		k := 0
		for i := 0; i < h.N; i++ {
			for j := i + 1; j < h.N; j++ {
				w := float64(genome[k])
				// randomise initial
				if genome[k] == 2 {
					if rand.Float64() < 0.5 {
						w = 0
					} else {
						w = 1
					}
				}
				h.W[i][j] = w
				h.W[j][i] = w
				k++
			}
		}

		// relax the dynamics
		repeated := false
		for !repeated {
			repeated = h.Step()
		}

		match := true
		for i := 0; i < h.N; i++ {
			if h.X[i] != target[i] {
				match = false
				break
			}
		}
		if match {
			return true, t, nil
		}

	}
	return false, trials, nil
}
