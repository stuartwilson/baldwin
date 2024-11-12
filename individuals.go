package main

type Basic struct {
	genome       *[]int
	fitness      float64
	modules      int
	moduleLength int
	trials       int
}

func NewBasic(n int, p *ProbabilitySelector, modules int, trials int) *Basic {
	return &Basic{
		genome:       NewGenome(n, p),
		modules:      modules,
		moduleLength: n / modules,
		trials:       trials,
	}
}

func (ind *Basic) GetGenome() *[]int {
	return ind.genome
}

func (ind *Basic) ComputeFitness(target []int) {

	q := 0.0
	k := 0
	for i := 0; i < ind.modules; i++ {
		module := make([]int, 0)
		targ := make([]int, 0)
		for j := 0; j < ind.moduleLength; j++ {
			module = append(module, (*ind.genome)[k])
			targ = append(targ, target[k])
			k++
		}
		if getFit(module, targ, ind.trials/ind.modules) {
			q++
		}
	}
	q /= float64(ind.modules)
	ind.fitness = q
}

func (ind *Basic) GetFitness() float64 {
	return ind.fitness
}

type Hoppy struct {
	genome  *[]int
	fitness float64
	trials  int
	H       *Hopfield
}

func NewHoppy(n int, p *ProbabilitySelector, trials int) *Hoppy {

	nG := (n*n - n) / 2
	return &Hoppy{
		genome: NewGenome(nG, p),
		trials: trials,
		H:      NewHopfield(n),
	}
}

func (ind *Hoppy) GetGenome() *[]int {
	return ind.genome
}

func (ind *Hoppy) ComputeFitness(target []int) {

	targ := make([]bool, len(target))
	for i := 0; i < len(target); i++ {
		targ[i] = target[i] > 0
	}

	f, _, _ := ind.H.Evaluate(*ind.genome, targ, ind.trials)

	if f {
		ind.fitness = 1.0
	} else {
		ind.fitness = 0.0
	}
}

func (ind *Hoppy) GetFitness() float64 {
	return ind.fitness
}
