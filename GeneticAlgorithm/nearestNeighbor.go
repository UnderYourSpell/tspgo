package main

import (
	"math"
	"math/rand"
)

//1. new random path
// V is the new path
// add to U

func findDist(city1 City, city2 City) float64 {
	x1 := city1.x
	y1 := city1.y
	x2 := city2.x
	y2 := city2.y
	distance := math.Sqrt(math.Pow((x2-x1), 2) + math.Pow((y2-y1), 2))
	return distance
}

// greedy nearest neighbor
func nearestNeighbor(initCities []City) Trip {
	var V, U []City

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
		V = append(V, initCities[numToAdd])
	}

	U = append(U, V[0])
	V = remove(V, 0)
	var dist float64
	var closestCity City
	var deleteIndex int
	for {
		dist = 999999999
		for i := 0; i < len(V); i++ {
			curDist := findDist(U[len(U)-1], V[i])
			if curDist < dist {
				closestCity = V[i]
				dist = curDist
				deleteIndex = i
			}
		}
		U = append(U, closestCity)
		V = remove(V, deleteIndex)
		if len(U) == len(initCities) {
			break
		}
	}

	newTrip := Trip{
		path:       U,
		pathLength: 0,
		prob:       0,
	}
	newTrip.calcPathLength()
	return newTrip
}
