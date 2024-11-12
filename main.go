package main

import "fmt"

func main() {
	//h := NewHopfield(7)

	n := 20
	target := make([]int, n)
	for i := 0; i < n; i++ {
		target[i] = 1
	}
	generations := 100
	pop := NewPopulation(1000, n, 1, 1000, NewProbabilitySelector([]float64{0.25, 0.25, 0.5}))
	correct, plastic := pop.Evolve(generations, target, 0)
	c := float64(correct[generations-1]) / float64(n)
	p := float64(plastic[generations-1]) / float64(n)
	i := 1 - c - p
	result := fmt.Sprintf("correct %f\nplastic %f\nincorrect %f", c, p, i)
	fmt.Println(result)
}
