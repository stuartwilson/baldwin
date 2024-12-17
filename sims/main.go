package main

import (
	"fmt"
	"github.com/stuartwilson/baldwin"
	"os"
	"strconv"
)

// main ...
// filename 			xxxx.json
// individualType 		Basic, Hopfield, or GRN
// network nodes
// population size
// generations
// adaptation trials

func main() {

	args := os.Args[1:]

	filename := args[0]
	IndividualType := args[1]
	n, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	populationSize, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println(err)
		return
	}

	generations, err := strconv.Atoi(args[4])
	if err != nil {
		fmt.Println(err)
		return
	}

	trials, err := strconv.Atoi(args[5])
	if err != nil {
		fmt.Println(err)
		return
	}

	unstable, err := strconv.Atoi(args[6])
	if err != nil {
		fmt.Println(err)
		return
	}

	initialPlastic, err := strconv.ParseFloat(args[7], 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	seed, err := strconv.Atoi(args[8])
	if err != nil {
		fmt.Println(err)
		return
	}

	extra := make([]string, 0)
	if len(args) > 9 {
		for i := 9; i < len(args); i++ {
			extra = append(extra, args[i])
		}
	}

	baldwin.Run(filename, IndividualType, n, populationSize, generations, trials, unstable, initialPlastic, int64(seed), extra)
}
