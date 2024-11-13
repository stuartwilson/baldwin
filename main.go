package main

import (
	"fmt"
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

func main() {

	//IndividualType := "Basic"
	IndividualType := "Hopfield"

	generations := 500
	ps := NewProbabilitySelector([]float64{0.25, 0.25, 0.5})

	P := make(Population, 0)
	populationSize := 1000

	var n int
	switch IndividualType {
	case "Basic":
		n = 20
		for i := 0; i < populationSize; i++ {
			P = append(P, NewBasic(n, ps, 1, 1000))
		}
	case "Hopfield":
		n = 7
		for i := 0; i < populationSize; i++ {
			P = append(P, NewHoppy(n, ps, 50))
		}
	}

	p, f := Evolve(P, generations, sameInts(n, 1), 0)

	fmt.Println("plastic = np.array([" + formatArray(p) + "])")
	fmt.Println("fitness = np.array([" + formatArray(f) + "])")
}

//result := fmt.Sprintf()
/*
	c := float64(correct[generations-1]) / float64(n)
	p := float64(plastic[generations-1]) / float64(n)
	i := 1 - c - p
	result := fmt.Sprintf("correct %f\nplastic %f\nincorrect %f", c, p, i)

*/
//fmt.Println(result)
