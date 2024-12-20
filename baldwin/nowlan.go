package baldwin

import (
	"math/rand"
)

// HOPFIELD NET

type Gene struct {
	weight  float64
	plastic bool
}

func (g *Gene) SetFromGenotype(geneotype int) {
	g.plastic = geneotype == 2
	if g.plastic {
		if rand.Float64() < 0.5 {
			g.weight = 1
		} else {
			g.weight = -1
		}
	} else {
		g.weight = float64(geneotype)*2 - 1
	}
}

type Hop struct {
	N             int
	G             int
	W             []Gene
	X             []bool
	target        []bool
	genome        []int
	IsSetToCorrel bool
	fitness       float64
}

func NewHop(n int, p *ProbabilitySelector) *Hop {
	g := (n*n - n) / 2
	h := Hop{
		N:             n,
		G:             g,
		X:             make([]bool, n),
		W:             make([]Gene, 0),
		target:        make([]bool, n),
		genome:        NewGenome(g, p),
		IsSetToCorrel: false,
	}
	for i := 0; i < g; i++ {
		h.W = append(h.W, Gene{})
	}
	return &h
}

func (h *Hop) SetToCorrel(patterns []bool) {
	weightIndex := 0
	for i := 0; i < h.N; i++ {
		for j := i + 1; j < h.N; j++ {
			if h.W[weightIndex].plastic {
				correl := patterns[i] == patterns[j]
				if correl {
					h.W[weightIndex].weight = 1.0
				} else {
					h.W[weightIndex].weight = -1.0
				}
			}
			weightIndex++
		}
	}
}

func (h *Hop) Step() bool {
	if true { // set
		if !h.IsSetToCorrel {
			if match(h.X, h.target) {
				h.SetToCorrel(h.target)
				h.IsSetToCorrel = true
			}
		}
	}

	sums := make([]float64, h.N)
	weightIndex := 0
	for i := 0; i < h.N; i++ {
		for j := i + 1; j < h.N; j++ {
			weight := h.W[weightIndex].weight

			if false { // randomize plastic on every step
				if !h.IsSetToCorrel && h.genome[weightIndex] == 2 {
					if rand.Float64() < 0.5 {
						weight = 1
					} else {
						weight = -1.0
					}
				}
			}

			if h.X[j] {
				sums[i] += weight
			} else {
				sums[i] -= weight
			}
			if h.X[i] {
				sums[j] += weight
			} else {
				sums[j] -= weight
			}
			weightIndex++
		}
	}

	allSame := true
	for i := 0; i < h.N; i++ {
		thresh := sums[i] > 0
		if thresh != h.X[i] {
			allSame = false
		}
		h.X[i] = thresh
	}
	return allSame
}

// FITNESS

func (h *Hop) ComputeFitness(target []int) {

	// INITIALIZE NETWORK STATE
	for i := 0; i < h.N; i++ {
		h.X[i] = false
	}

	// INITIALIZE TARGET
	for i := 0; i < h.N; i++ {
		if target[i] > 0 {
			h.target[i] = true
		}
	}

	// INITIALIZE FROM GENOME
	for i := 0; i < h.G; i++ {
		h.W[i].SetFromGenotype(h.genome[i])
	}

	// EVALUATE
	stateSpaceSize := int(pow(2, h.N))
	for i := 0; i < stateSpaceSize; i++ {
		if h.Step() {
			if match(h.X, h.target) {
				h.fitness = 1
			} else {
				h.fitness = 0
			}
			return
		}
	}
	h.fitness = 0

}

func (h *Hop) GetGenome() []int {
	return h.genome
}

func (h *Hop) SetGenome(genome []int) {
	copy(h.genome, genome)
	// INITIALIZE FROM GENOME
	for i := 0; i < h.G; i++ {
		h.W[i].SetFromGenotype(h.genome[i])
	}
	h.IsSetToCorrel = false
}

func (h *Hop) GetFitness() float64 {
	return h.fitness
}

//func Correl(m []bool) {
//weights := make([]int, len(m))
//}

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

*/

/*
func (g Gene) SetPlastic(correlated bool) {
	g.correl *= g.nSet
	if correlated {
		g.correl++
	} else {
		g.correl--
	}
	g.nSet++
	g.correl /= g.nSet
}

func (g Gene) SetProbability(rate float64) {
	g.prob = 1.0 / (1 + math.Exp(-g.correl*rate))
}

func (g Gene) Get() float64 {
	if g.plastic {
		if rand.Float64() < g.prob {
			return 1.0
		} else {
			return -1.0
		}
	} else {
		return g.weight
	}
}

func (g Gene) SetWeight() {
	g.weight = g.Get()
}

*/

/*
func (h *Hop) Set(genome []int, patterns []bool, probCorrect float64) {
	//if len(genome) != h.G {
	//	fmt.Println("Invalid genome")
	//	return
	//}

	// get correlations
	correl := make([]bool, h.G)
	weightIndex := 0
	for i := 0; i < h.N; i++ {
		for j := i + 1; j < h.N; j++ {
			correl[weightIndex] = patterns[i] == patterns[j]
			weightIndex++
		}
	}

	for i := 0; i < len(genome); i++ {
		//h.W[i].SetState(genome[i])
		if h.W[i].plastic {
			if (rand.Float64() < probCorrect) == correl[i] {
				h.W[i].weight = 1.0
			} else {
				h.W[i].weight = -1.0
			}
		}

		//h.W[i].SetPlastic(correl[i])
		//h.W[i].SetProbability(rate)
		//h.W[i].SetWeight() // THIS LOCKS THE VALUE OF THE PLASTIC WEIGHTS ONCE
	}
}
*/

// NOWLAN

//type Nowlan struct {
//	genome  []int
//	fitness float64
//	trials  int
//	H       *Hop
//	//switchCase int
//}
//
//func NewNowlan(n int, p *ProbabilitySelector, trials int, extra []string) *Nowlan {
//
//	/*
//		var err error
//		switchCase := 0 // default switch case is 1
//		if len(extra) > 0 {
//			switchCase, err = strconv.Atoi(extra[0])
//			if err != nil {
//				fmt.Println("Switch case must be an integer")
//				return nil
//			}
//		}
//
//	*/
//
//	h := NewHop(n)
//	genome := NewGenome(h.G, p)
//	h.genome = genome
//	return &Nowlan{
//		genome: genome,
//		trials: trials,
//		H:      h,
//		//switchCase: switchCase,
//	}
//
//}

//func (h *Hop) ComputeFitness(target []int) {
//	targ := make([]bool, len(target))
//	for i := 0; i < len(target); i++ {
//		targ[i] = target[i] > 0
//	}
//f := h.Evaluate(h.genome, targ)
//if f {
//	h.fitness = 1.0
//} else {
//	h.fitness = 0.0
//}
//}
