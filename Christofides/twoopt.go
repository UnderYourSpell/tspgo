package main

import (
	"fmt"
	"math"
	"math/rand"
)

/*
	2Opt Swap code
*/

func twoOptSwap(nodes []Node, v1 int, v2 int) []Node {
	var newPath []Node
	newPath = append(newPath, nodes[:v1]...) //append in order from start to v1
	for i := v2; i >= v1; i-- {
		newPath = append(newPath, nodes[i]) //append in reverse order from v1 to v2
	}
	if v2+1 < len(nodes) {
		newPath = append(newPath, nodes[v2+1:]...) //append in order from v2 to end
	}
	return newPath
}

func exponentialCooling(TStart float64, alpha float64, iteration int) float64 {
	return TStart * (math.Pow(alpha, float64(iteration)))
}

// Work in progress. Attempting to create several go routines that recursivly try and do 2opt on suboptimal solutions
func twoOptPathCreateParallel(path []Node, temp chan float64, globalBest chan float64, iteration chan int) {
	localIter := <-iteration
	localBest := <-globalBest
	localTemp := <-temp
	bestDistance := calcPathLength(path)
	for {
		improved := false
		//start again
		for i := range path {
			for j := range path {
				newPath := twoOptSwap(path, i, j)
				newDistance := calcPathLength(newPath)
				if newDistance < bestDistance {
					path = newPath
					bestDistance = newDistance
					improved = true
					if bestDistance < localBest {
						globalBest <- bestDistance
					}
					break
				} else {
					rng := rand.Float64()
					pAccept := math.Pow(math.E, (-(newDistance - bestDistance) / localTemp)) //acceptance probability
					if pAccept > rng {
						path = newPath
						bestDistance = newDistance
						improved = true
						localIter++
						localTemp = exponentialCooling(localTemp, alpha, localIter)
						fmt.Println(localTemp)
						temp <- localTemp
						iteration <- localIter
						go twoOptPathCreateParallel(newPath, temp, globalBest, iteration)
						break
					}
				}
			}
			if improved {
				break
			}
		}
		if !improved {
			break
		}
	}
}

// Normal 2opt swap
func twoOptPathCreateSequential(path []Node) []Node {
	bestDistance := calcPathLength(path)
	for {
		improved := false
		//start again
		for i := range path {
			for j := range path {
				newPath := twoOptSwap(path, i, j)
				newDistance := calcPathLength(newPath)
				if newDistance < bestDistance {
					path = newPath
					bestDistance = newDistance
					improved = true
					break
				}
			}
			if improved {
				break
			}
		}
		if !improved {
			break
		}
	}
	return path
}

const alpha = 0.99 //temperature cooling alpha

// Temperature controlled 2opt, does simulated annealing in sequence
func twoOptPathCreateTemp(path []Node) []Node {
	bestDistance := calcPathLength(path)
	temp := 1.0
	iter := 1
	for {
		improved := false
		//start again
		for i := range path {
			for j := range path {
				newPath := twoOptSwap(path, i, j)
				newDistance := calcPathLength(newPath)
				if newDistance < bestDistance {
					path = newPath
					bestDistance = newDistance
					improved = true
					break
				} else {
					rng := rand.Float64()
					pAccept := math.Pow(math.E, (-(newDistance - bestDistance) / temp)) //acceptance probability
					if pAccept > rng {
						path = newPath
						bestDistance = newDistance
						improved = true
						iter++
						temp = exponentialCooling(temp, alpha, iter)
						break
					}
				}
			}
			if improved {
				break
			}
		}
		if !improved {
			break
		}
	}
	return path
}
