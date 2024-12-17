package baldwin

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Basic struct {
	genome       []int
	fitness      float64
	modules      int
	moduleLength int
	trials       int
}

func NewBasic(n int, p *ProbabilitySelector, trials int, extra []string) *Basic {

	//
	var err error
	modules := 1 // default number of modules is 1
	if len(extra) > 0 {
		modules, err = strconv.Atoi(extra[0])
		if err != nil {
			fmt.Println("Number of modules must be an integer")
			return nil
		}
	}

	return &Basic{
		genome:       NewGenome(n, p),
		modules:      modules,
		moduleLength: n / modules,
		trials:       trials,
	}
}

func (ind *Basic) GetGenome() []int {
	return ind.genome
}

func (ind *Basic) SetGenome(genome []int) {
	copy(ind.genome, genome)
}

func (ind *Basic) getFit(x, target []int, trials int) (bool, float64) {
	a := 0
	for i := 0; i < len(x); i++ {
		if x[i] == 2 {
			a++
		} else {
			if x[i] != target[i] {
				return false, 0.0
			}
		}
	}
	prob := 1.0 - pow((1.0-pow(0.5, a)), trials)
	return rand.Float64() < prob, prob
}

func (ind *Basic) ComputeFitness(target []int) {

	q := 0.0
	k := 0
	for i := 0; i < ind.modules; i++ {
		module := make([]int, 0)
		targ := make([]int, 0)
		for j := 0; j < ind.moduleLength; j++ {
			module = append(module, ind.genome[k])
			targ = append(targ, target[k])
			k++
		}
		fit, _ := ind.getFit(module, targ, ind.trials/ind.modules)
		if fit {
			q++
		}
	}
	q /= float64(ind.modules)

	ind.fitness = q
}

func (ind *Basic) GetFitness() float64 {
	return ind.fitness
}
