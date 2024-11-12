package main

import "fmt"

func main() {

	n := 20
	target := make([]int, n)
	for i := 0; i < n; i++ {
		target[i] = 1
	}
	generations := 100
	ps := NewProbabilitySelector([]float64{0.25, 0.25, 0.5})
	P := make(Population, 0)
	for i := 0; i < 1000; i++ {
		//P = append(P, NewBasic(n, ps, 1, 1000))
		P = append(P, NewHoppy(n, ps, 1000))
	}

	correct, plastic := Evolve(P, generations, target, 0)
	c := float64(correct[generations-1]) / float64(n)
	p := float64(plastic[generations-1]) / float64(n)
	i := 1 - c - p
	result := fmt.Sprintf("correct %f\nplastic %f\nincorrect %f", c, p, i)
	fmt.Println(result)
}
