package baldwin

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Hoppy struct {
	genome     []int
	fitness    float64
	trials     int
	H          *Hopfield
	switchCase int
}

func NewHoppy(n int, p *ProbabilitySelector, trials int, extra []string) *Hoppy {

	var err error
	switchCase := 0 // default switch case is 1
	if len(extra) > 0 {
		switchCase, err = strconv.Atoi(extra[0])
		if err != nil {
			fmt.Println("Switch case must be an integer")
			return nil
		}
	}

	nG := (n*n - n) / 2
	return &Hoppy{
		genome:     NewGenome(nG, p),
		trials:     trials,
		H:          NewHopfield(n),
		switchCase: switchCase,
	}
}

func (ind *Hoppy) GetGenome() []int {
	return ind.genome
}

func (ind *Hoppy) SetGenome(genome []int) {
	copy(ind.genome, genome)
}

func (ind *Hoppy) ComputeFitness(target []int) {

	targ := make([]bool, len(target))
	for i := 0; i < len(target); i++ {
		targ[i] = target[i] > 0
	}

	f, _, _ := ind.H.Evaluate(ind.genome, targ, ind.trials, ind.switchCase)

	if f {
		ind.fitness = 1.0
	} else {
		ind.fitness = 0.0
	}
}

func (ind *Hoppy) GetFitness() float64 {
	return ind.fitness
}

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
	copy(h.X, next)
	return allSame
}

func (h *Hopfield) Evaluate(genome []int, target []bool, trials int, switchCase int) (bool, int, error) {

	if len(genome) != (h.N*h.N-h.N)/2 {
		return false, 0, fmt.Errorf("wrong genome length")
	}

	stateSpaceSize := int(pow(2, h.N))

	// set Hopfield weights from genome
	k := 0
	for i := 0; i < h.N; i++ {
		for j := i + 1; j < h.N; j++ {
			if genome[k] != 2 {
				w := float64(genome[k])*2 - 1
				h.W[i][j] = w
				h.W[j][i] = w
				k++
			}
		}
	}

	for t := 0; t < trials; t++ {

		switch switchCase {
		case 0:
			// initial random state
			for i := 0; i < h.N; i++ {
				h.X[i] = rand.Float64() < 0.5
			}
		default:
			fmt.Println("invalid case")
			return false, trials, fmt.Errorf("invalid case")
		}

		// set Hopfield weights from genome
		k := 0
		for i := 0; i < h.N; i++ {
			for j := i + 1; j < h.N; j++ {
				// randomise initial
				var w float64
				if genome[k] == 2 {
					if rand.Float64() < 0.5 {
						w = -1.0
					} else {
						w = 1.0
					}
					h.W[i][j] = w
					h.W[j][i] = w
				}
				k++
			}
		}

		// relax the dynamics

		//repeated := false
		//fmt.Println(stateSpaceSize)
		for i := 0; i < stateSpaceSize/2; i++ {
			if h.Step() {
				break
			}
		}

		if match(h.X, target) {
			return true, t, nil
		}

	}
	return false, trials, nil
}

/* //TODO: ALL GOOD BELOW BUT COMMENTED AS NOT STRICTLY NECESSARY

func (h *Hopfield) SetState(x []bool) {
	h.X = x
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

*/
