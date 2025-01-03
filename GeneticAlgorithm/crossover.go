package main

import (
	"math/rand"
)

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
			child1 = append(child1, (*gene2).path[i])
		}
		if _, ok := ids2[(*gene1).path[i].id]; !ok {
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

type adjMatrix struct {
	adj map[string]bool
}

func union(s1, s2 map[string]bool) map[string]bool {
	s_union := map[string]bool{}
	for k, _ := range s1 {
		s_union[k] = true
	}
	for k, _ := range s2 {
		s_union[k] = true
	}
	//k is the key
	return s_union
}

func EdgeRecombination(gene1 *Trip, gene2 *Trip, children *[]Trip, cityMap *map[string]City, numCities int) {
	path1 := (*gene1).path
	path2 := (*gene2).path
	adjGene1 := make(map[string]adjMatrix)
	adjGene2 := make(map[string]adjMatrix)

	//add cities to adjGene1 and Gene2
	temp1 := map[string]bool{path1[0].id: true, path1[numCities-1].id: true}
	temp2 := map[string]bool{path2[0].id: true, path2[numCities-1].id: true}
	adjGene1[path1[0].id] = adjMatrix{temp1}
	adjGene2[path2[0].id] = adjMatrix{temp2}
	for i := 1; i < numCities-1; i++ {
		temp1 := map[string]bool{path1[i-1].id: true, path1[i+1].id: true}
		temp2 := map[string]bool{path2[i-1].id: true, path2[i+1].id: true}
		adjGene1[path1[i].id] = adjMatrix{temp1}
		adjGene2[path2[i].id] = adjMatrix{temp2}
	}
	temp1 = map[string]bool{path1[0].id: true, path1[numCities-2].id: true}
	temp2 = map[string]bool{path2[0].id: true, path2[numCities-2].id: true}
	adjGene1[path1[numCities-1].id] = adjMatrix{temp1}
	adjGene2[path2[numCities-1].id] = adjMatrix{temp2}

	//make union of sets
	masterMap := make(map[string]adjMatrix)
	for i := 0; i < numCities; i++ {
		key := path1[i].id
		set1 := adjGene1[key].adj
		set2 := adjGene2[key].adj
		newSet := union(set1, set2)
		masterMap[key] = adjMatrix{newSet}
	}

	var K []City
	var N City
	selectParent := rand.Float32()
	if selectParent >= 0.5 {
		N = path1[0]
	} else {
		N = path2[0]
	}

	for {
		if len(K) == numCities {
			break
		}
		K = append(K, N)
		remId := N.id
		//remove N from all neighbor lists
		for _, set := range masterMap {
			delete(set.adj, remId)
		}
		//if N's neighbor list is non-empty then do stuff
		var NStar City
		var nextCity string
		var equalSmallest []string
		Nlist := masterMap[remId].adj
		if len(Nlist) > 0 {
			bestLen := 10
			for key := range Nlist {
				curLength := len(masterMap[key].adj)
				if curLength < bestLen {
					equalSmallest = append(equalSmallest, key)
					bestLen = curLength
				} else if curLength == bestLen {
					equalSmallest = append(equalSmallest, key)
				}
			}
			nextNRand := rand.Intn(len(equalSmallest))
			nextCity = equalSmallest[nextNRand]
			NStar = (*cityMap)[nextCity]
		} else { //if no neighbors in N, choose a random one, need to see what's not currently in K
			for cityName := range *cityMap {
				found := 0
				for c := range K {
					if K[c].id == cityName {
						found = 1
						break
					}
				}
				if found == 0 { //if not found then add
					NStar = (*cityMap)[cityName]
					break
				}
			}
		}
		N = NStar
	}
	newTrip := Trip{
		path:       K,
		pathLength: 0,
		prob:       0,
	}
	newTrip.calcPathLength()
	(*children) = append((*children), newTrip)
}

//adjGene1 needs to have a string key that correlates to a map[string]bool of the cities the key is adjacent to
//then we can more easily union the sets
//for ERX we need to make a mapping of all the cities for this to work
