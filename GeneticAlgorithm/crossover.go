package main

import "math/rand"

type void struct{}

var member void

// Single Point Crossover
func SPX(gene1 *Trip, gene2 *Trip, children *[]Trip, numCities int) {
	cut := rand.Intn(numCities)
	ids1 := make(map[string]void)
	ids2 := make(map[string]void)
	var child1 []City
	var child2 []City

	for i := 0; i < cut; i++ {
		child1 = append(child1, (*gene1).path[i])
		child2 = append(child2, (*gene2).path[i])
		ids1[(*gene1).path[i].id] = member
		ids2[(*gene2).path[i].id] = member
	}

	for i := 0; i < numCities; i++ {
		if _, ok := ids1[(*gene2).path[i].id]; !ok {
			//add to child 1
			child1 = append(child1, (*gene2).path[i])
		}
		if _, ok := ids2[(*gene1).path[i].id]; !ok {
			//add to child 2
			child2 = append(child2, (*gene1).path[i])
		}
	}

	//need to make some init function for this
	child1Trip := Trip{
		path:       child1,
		pathLength: 0,
		prob:       0,
	}
	child2Trip := Trip{
		path:       child2,
		pathLength: 0,
		prob:       0,
	}
	child1Trip.calcPathLength()
	child2Trip.calcPathLength()
	(*children) = append((*children), child1Trip, child2Trip)
}
