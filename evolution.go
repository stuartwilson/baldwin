package main

import "math/rand"

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
	GetFitness([]int) float64
}

func NewGenome(n int, p *ProbabilitySelector) *[]int {
	x := make([]int, n)
	for i := 0; i < n; i++ {
		x[i] = p.Select()
	}
	return &x
}

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

func (ind *Basic) GetFitness(target []int) float64 {

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
	return ind.fitness
}

func Combine(mum, dad IndividualI, crossover int) IndividualI {
	m := mum.GetGenome()
	d := dad.GetGenome()
	child := append((*m)[:crossover], (*d)[crossover:]...)
	m = &child
	return mum
}

type Population struct {
	N            int // number of heritable units
	M            int // number of modules (note that n/m should be an integer)
	U            int // number of unstable units
	P            int // number in population
	moduleLength int
	minFitness   float64
	Pop          []IndividualI
}

func NewPopulation(p, n, m, trials int, ps *ProbabilitySelector) *Population {
	pop := Population{
		N:          n,
		P:          p,
		minFitness: 1.0 / float64(n),
		Pop:        make([]IndividualI, p),
	}

	for i := 0; i < p; i++ {
		pop.Pop[i] = NewBasic(n, ps, m, trials)
	}

	return &pop
}

func (p *Population) Evolve(generations int, target []int, nUnstable int) ([]int, []int) {

	// store the number of units of each type
	countPlastic := make([]int, generations)
	countCorrect := make([]int, generations)

	unstable := rand.Perm(p.N)[:nUnstable]

	// evolution loop
	for g := 0; g < generations; g++ {

		//count units(for plotting)
		for _, i := range p.Pop {
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
		countPlastic[g] /= p.P
		countCorrect[g] /= p.P

		// evaluate fitness of each individual
		F := make([]float64, p.P)
		for i := 0; i < p.P; i++ {
			F[i] = p.minFitness + (1-p.minFitness)*p.Pop[i].GetFitness(target)
		}
		selector := NewProbabilitySelector(F)

		next := make([]IndividualI, p.P)
		for i := 0; i < p.P; i++ {
			mum := p.Pop[selector.Select()]
			dad := p.Pop[selector.Select()]
			crossover := 1 + int(rand.Float64()*float64(p.N-2))
			next[i] = Combine(mum, dad, crossover)
		}
		p.Pop = next

		// randomly assign target state of unstable units
		for _, i := range unstable {
			if rand.Float64() < 0.5 {
				target[i] = 0
			} else {
				target[i] = 1
			}
		}
	}
	return countCorrect, countPlastic
}

// SIMULATION PARAMETERS

//n=20
//m=1
//u=5
//t=1000
//P=1000
//G=100

//# plotting
//F = pl.figure()
//f = F.add_subplot(111)
//f.plot(countCorrect,linestyle='dashed',color='k')
//f.plot(countPlastic,'-',color='k')
//f.plot(n-countCorrect-countPlastic,linestyle='dotted',color='k')
//f.legend(['correct','plastic','incorrect'],frameon=False)
//f.axis([0,G-1,0,n])
//f.set_xlabel('generations')
//f.set_ylabel('units')
//pl.show()

/*
import numpy as np
import pylab as pl
import sys

'''
    This script implements the Population model of the evolution of phenotypic plasticty described in the paper 'Tipping the scales between genetic assimilation and phenotypic plasticity'. Set m=1, m=2, or m=4 to re-produce plots of the form shown in Fig. 3A, 3B, and 3C, respectively.
'''


# define fitness function
def getFit(x,target,t):
    if(np.any((x<2)*(x!=target))): return False
    a = np.sum(x==2)*1.
    adaptation_success = 1.0-pow((1.0-pow(0.5,a)),t)
    return (np.random.rand() < adaptation_success)

# (optionally) supply seed for random number generator at the command line
if(len(sys.argv)>1): np.random.seed(int(sys.argv[1]))

# SIMULATION PARAMETERS

n=20              # number of heritable units
m=1               # number of modules (note that n/m should be an integer)
u=5               # number of unstable units
t=1000            # number of within-lifetime adaptation trials
P=1000            # number in population
G=100             # number of generations

moduleLength=int(n/m)
s=1/n

# initial population (0 and 1 represent fixed states, 2 represents plastic)
X = np.fmax(1*(np.random.rand(P,n)<0.5),2*(np.random.rand(P,n)<0.5))

# store the number of units of each type
countPlastic = np.zeros(G)
countCorrect = np.zeros(G)

# define target phenotype
target = np.ones(n)
unstable_units = np.random.permutation(u)

# evolution loop
for g in range(G):

    # count units (for plotting)
    countPlastic[g] = np.sum(X==2)/P
    countCorrect[g] = np.sum(X==np.tile(target,[P,1]))/P

    # evaluate fitness of each individual
    f = np.zeros(P)
    for i in range(P):
        q = 0.0
        for j in range(m):
            module = range(j*moduleLength,(j+1)*moduleLength)
            q += getFit(X[i,module],target[module],t/m)
        f[i] = q/m

    # PDF for selecting parents
    S = s+(1-s)*f
    S = S/np.sum(S)
    boundB = np.cumsum(S)
    boundA = np.hstack([0,boundB[:-1]])

    # Populate next generation

    Xnext = np.zeros([P,n])

    for i in range(P):

        # identify parents
        r = np.random.rand(2)
        parent1 = np.where((boundA<=r[0])*(r[0]<boundB))
        parent2 = np.where((boundA<=r[1])*(r[1]<boundB))

        # create offspring
        crossover = 1+int(np.random.rand()*(n-2))
        Xnext[i] = X[parent1]
        Xnext[i,crossover:] = X[parent2,crossover:]

    X = Xnext

    # randomly assign target state of unstable units
    target[unstable_units] = 1*(np.random.rand(u)<0.5)


# plotting
F = pl.figure()
f = F.add_subplot(111)
f.plot(countCorrect,linestyle='dashed',color='k')
f.plot(countPlastic,'-',color='k')
f.plot(n-countCorrect-countPlastic,linestyle='dotted',color='k')
f.legend(['correct','plastic','incorrect'],frameon=False)
f.axis([0,G-1,0,n])
f.set_xlabel('generations')
f.set_ylabel('units')
pl.show()

*/
