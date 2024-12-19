package main

import (
	"math"
	"math/rand"
)

// Linear Rank Selection
func LRS(genePool *[]Trip, parents *[]Trip, selectionPressure int) {
	//gene pool is sorted before it comes to this function
	sp := float64(selectionPressure)
	n := float64(len(*genePool))
	intN := int(n)
	firstFactor := math.Pow(n, -1)
	secondFactor := 2*sp - 2

	p := firstFactor * sp
	(*genePool)[0].prob = p
	for i := 1; i < intN; i++ {
		p = firstFactor*(sp-secondFactor*(float64(i-1)/(n-1))) + (*genePool)[i-1].prob
		(*genePool)[i].prob = p
	}

	for i := 0; i < intN/2; i++ {
		a := rand.Float64()
		for j := 0; j < intN; j++ {
			if (*genePool)[j].prob >= a {
				if j != 0 {
					(*parents) = append((*parents), (*genePool)[j-1])
					break
				}
			}
		}
	}
}

/*
	Roulette Wheel Selection

What I did differently in the go version compared to the C++ version was
calculating the probability - or range - of each gene before we enter the selection loop.
This makes it so we are doing far less division. If we had a large population,
we would be dividing each path length by total inverted path length everytime we were calculating cumulative probabilty
in the selection loop. Float64 division is slow, so we're not doing as much of it anymore.
*/
func RWS(genePool *[]Trip, parents *[]Trip) {
	n := len(*genePool)
	var S float32
	S = 0
	for i := 0; i < n; i++ {
		S += (*genePool)[i].pathLength
	}

	for i := 0; i < n; i++ {
		(*genePool)[i].prob = float64((*genePool)[i].pathLength) / float64(S)
	}

	for i := 0; i < (n / 2); i++ {
		a := rand.Float64()
		var iSum float64
		j := 0
		for {
			iSum = iSum + (*genePool)[j].prob
			if iSum > a {
				break
			}
			j++
		}
		(*parents) = append((*parents), (*genePool)[j])
	}
}
