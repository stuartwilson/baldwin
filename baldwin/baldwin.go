package baldwin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
)

func formatArray(x []float64) string {
	s := ""
	for i := 0; i < len(x); i++ {
		s += fmt.Sprint(x[i], ",")
	}
	return s[:len(s)-1]
}

func randInts(n int) []int {
	x := make([]int, n)
	for i := 0; i < n; i++ {
		if rand.Float64() < 0.5 {
			x[i] = 1
		}
	}
	return x
}

func sameInts(n, v int) []int {
	x := make([]int, n)
	for i := 0; i < n; i++ {
		x[i] = v
	}
	return x
}

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
	fmt.Println(1.0 - pow((1.0-pow(0.5, a)), trials))
	return rand.Float64() < 1.0-pow((1.0-pow(0.5, a)), trials)
}

type Result struct {
	N              int       `json:"n"`
	Ngenes         int       `json:"genes"`
	Trials         int       `json:"trials"`
	Plastic        []float64 `json:"plastic"`
	Fitness        []float64 `json:"fitness"`
	Genomes        []int     `json:"genomes"`
	NumUnique      int       `json:"numUnique"`
	Unique         [][]int   `json:"unique"`
	PopulationSize int       `json:"populationSize"`
	Generations    int       `json:"generations"`
	GensEvaluated  int       `json:"gensEvaluated"`
	IndividualType string    `json:"individualType"`
	Filename       string    `json:"filename"`
	Probs          []float64 `json:"probs"`
	Unstable       int       `json:"unstable"`
	Seed           int64     `json:"seed"`
	Extra          []string  `json:"extra"`
	PerUnique      []int     `json:"perUnique"`
	InitialPlastic float64   `json:"initialPlastic"`
}

func Run(filename, IndividualType string, n, populationSize, generations, trials, unstable int, initialPlastic float64, seed int64, extra []string) {

	rand.Seed(seed)

	//IndividualType := "Basic"
	//IndividualType := "Hopfield"
	//IndividualType := "GRN"

	//generations := 500
	ps := NewProbabilitySelector([]float64{(1 - initialPlastic) / 2, (1 - initialPlastic) / 2, initialPlastic})

	P := make(Population, 0)
	//populationSize := 1000

	//var n int
	//var trials int
	switch IndividualType {
	case "Basic":
		//n = 20
		//trials = 1000
		for i := 0; i < populationSize; i++ {
			P = append(P, NewBasic(n, ps, trials, extra))
		}
	case "Hopfield":
		//n = 7
		//trials = 50
		for i := 0; i < populationSize; i++ {
			P = append(P, NewHoppy(n, ps, trials, extra))
		}
	case "GRN":
		//n = 6
		//trials = 50
		for i := 0; i < populationSize; i++ {
			P = append(P, NewGRN(n, ps, trials, extra))
		}
	default:
		fmt.Println("Invalid individual type: ", IndividualType)
		return
	}

	gens, p, f, genomes, unique, perUnique := Evolve(P, generations, sameInts(n, 1), unstable)

	r := Result{
		N:              n,
		Trials:         trials,
		Plastic:        p,
		Fitness:        f,
		Genomes:        genomes,
		Unique:         unique,
		PerUnique:      perUnique,
		NumUnique:      len(unique),
		Ngenes:         len(P[0].GetGenome()),
		PopulationSize: populationSize,
		Generations:    generations,
		GensEvaluated:  gens,
		IndividualType: IndividualType,
		Filename:       filename,
		Probs:          ps.Probs,
		Seed:           seed,
		Unstable:       unstable,
		Extra:          extra,
		InitialPlastic: initialPlastic,
	}

	jsonData, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("Data successfully written to json output file")

	//fmt.Println("plastic = np.array([" + formatArray(p) + "])")
	//fmt.Println("fitness = np.array([" + formatArray(f) + "])")
}

//result := fmt.Sprintf()
/*
	c := float64(correct[generations-1]) / float64(n)
	p := float64(plastic[generations-1]) / float64(n)
	i := 1 - c - p
	result := fmt.Sprintf("correct %f\nplastic %f\nincorrect %f", c, p, i)

*/
//fmt.Println(result)
