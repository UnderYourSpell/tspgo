package main

import "math/rand"

func swapMutate(gene *Trip, numCities int) {
	firstIndex := rand.Intn(numCities)
	secondIndex := rand.Intn(numCities)
	temp := (*gene).path[firstIndex]
	(*gene).path[firstIndex] = (*gene).path[secondIndex]
	(*gene).path[secondIndex] = temp
	gene.calcPathLength()
}
