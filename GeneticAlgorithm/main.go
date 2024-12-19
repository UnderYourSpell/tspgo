package main

/*
NOTES
Due to the speed at which this program runs, I believe it to be a good idea to try and recreate the Edge Recombination Operator
Also, due to how well go does with mapping, trying to implement Partially Mapped Crossover would be cool
But, I think trying to write a fast Roulette Wheel would be the first thing I should do

TODO
1. Nearest Neighbor - DONE
2. Roulette Wheel - DONE
3. Edge Recombination
4. Partially Mapped Crossover
*/

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main() {
	const maxGens = 3000
	const popSize = 128
	const mutateRate = 0.8
	const elitism = 2
	const nn = 1
	selection := "LRS"

	//read in file
	fp := "./tsp/eil51.tsp"
	data := readEucTSPFile(fp)
	initCities := data.initCities
	numCities := data.dimension

	var genePool []Trip

	if nn == 1 { //nearest Neighbor
		nnTrip := nearestNeighbor(initCities)
		for i := 0; i < popSize; i++ {
			genePool = append(genePool, nnTrip)
		}
	} else {
		for i := 0; i < popSize; i++ {
			var newTrip Trip
			n := len(initCities)

			var picks []int
			for l := 0; l < n; l++ {
				picks = append(picks, l)
			}

			for j := 0; j < len(initCities); j++ {
				randIndex := rand.Intn(n)
				numToAdd := picks[randIndex]
				picks[randIndex] = picks[n-1]
				n--
				newTrip.path = append(newTrip.path, initCities[numToAdd])
			}
			newTrip.calcPathLength()
			genePool = append(genePool, newTrip)
		}
	}

	start := time.Now() //start the clock
	//main loop
	for p := 0; p < maxGens; p++ {
		var parents []Trip

		//LRS for the start
		if selection == "LRS" {
			LRS(&genePool, &parents, 2)
		} else if selection == "RWS" {
			RWS(&genePool, &parents)
		}
		//SPX
		var children []Trip
		for i := 0; i < popSize/2; i += 2 {
			SPX(&parents[i], &parents[i+1], &children, numCities)
		}

		//for children mutate given a threshold
		for i := 0; i < popSize/2; i++ {
			mutateThreshold := rand.Float32()
			if mutateThreshold >= mutateRate {
				swapMutate(&children[i], numCities)
			}
		}

		//Sort original gene pool
		sort.Slice(genePool, func(i, j int) bool {
			return genePool[i].pathLength < genePool[j].pathLength
		})

		var newGen []Trip
		newGen = append(newGen, genePool[:elitism]...)
		parents = append(parents, children...)

		sort.Slice(parents, func(i, j int) bool {
			return parents[i].pathLength < parents[j].pathLength
		})
		newGen = append(newGen, parents[:popSize-elitism]...)

		genePool = nil
		genePool = newGen

	}
	elapsed := time.Since(start)
	fmt.Printf("Runtime: %s\n", elapsed)
	fmt.Print("Best Path Found: ")
	fmt.Print("Length: ", genePool[0].pathLength, " ")

	//Generic Print Statment
	/*
		for i := 0; i < len(parents); i++ {
			fmt.Print("Path ", i)
			fmt.Print(" ", parents[i].pathLength)
			parents[i].printPath()
		}
	*/
}
