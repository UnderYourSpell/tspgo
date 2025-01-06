package main

import (
	"fmt"
	"math"
	"math/rand"
)

func twoOptSwap(nodes []Node, v1 int, v2 int) []Node {
	var newPath []Node
	newPath = append(newPath, nodes[:v1]...)
	for i := v2; i >= v1; i-- {
		newPath = append(newPath, nodes[i])
	}
	if v2+1 < len(nodes) {
		newPath = append(newPath, nodes[v2+1:]...)
	}
	return newPath
}

func exponentialCooling(TStart float64, alpha float64, iteration int) float64 {
	return TStart * (math.Pow(alpha, float64(iteration)))
}

func twoOptPathCreateParallel(path []Node, temp chan float64) []Node {
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
					break
				} else {
					rng := rand.Float64()
					pAccept := math.Pow(math.E, (-(newDistance - bestDistance) / localTemp)) //acceptance probability
					if pAccept > rng {
						fmt.Println("Not finished :)")
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
