package main

/*
	Christofides Algorithm
	By Moro Bamber
*/

import (
	"fmt"
	"time"
)

func main() {
	//step 1. create a graph with all edges in the tree
	//read in file
	fp := "./tsp/eil101.tsp"
	data := readEucTSPFile(fp)

	//create list of edges and weights
	numCities := data.dimension
	nodes := data.initCities
	graph := make(map[string]Node) //master map
	for i := range numCities {
		graph[nodes[i].id] = nodes[i]
	}
	edges := createEdges(nodes)

	start := time.Now() //start the clock for the tree - this is the meat of thw algorithm, everything else was set up

	//Create Minimum Spanning Tree
	var tree []Edge
	tree, graph = createMST(numCities, nodes, edges, graph)

	//Create Min Cost Perfect Matching Eulerian Circuit
	nodeOrder := minCostPerfMatch(numCities, nodes, tree, graph, edges)

	//Create Hamiltonian Path
	finalGraph := hamiltonianPath(nodeOrder)

	pathLength := pathLengthEdges(finalGraph)
	fmt.Println("Path Length w/o 2opt swap:", pathLength)

	//Path is the sequence of nodes created from the edge sequence of the final graph
	path := createPath(finalGraph)

	var bestPath []Node
	parallel := false
	tempCreate := true //are we using temperature controlled 2opt (simulated annealing)?
	if parallel {
		//work in progress
		temp := make(chan float64, 1)
		temp <- 1
		globalBest := make(chan float64, 1)
		globalBest <- calcPathLength(path)
		iteration := make(chan int, 1)
		iteration <- 0
		go twoOptPathCreateParallel(path, temp, globalBest, iteration)
		fmt.Println("Global Best", <-globalBest)
		<-temp
		<-iteration
	} else if tempCreate { //twoOpt Sequential with temp control, works the best
		bestPath = twoOptPathCreateTemp(path)
		elapsed := time.Since(start)
		fmt.Println("Final Length with 2opt swap:", calcPathLength(bestPath))
		fmt.Printf("Runtime: %s\n", elapsed)
	} else { //base 2opt algorithm
		bestPath = twoOptPathCreateSequential(path)
		elapsed := time.Since(start)
		fmt.Println("Final Length with 2opt swap:", calcPathLength(bestPath))
		fmt.Printf("Runtime: %s\n", elapsed)
	}

	//create output graph -- Optional
	var finalEdges []Edge
	for i := 0; i < len(bestPath)-1; i++ {
		newEdge := createEdge(bestPath[i], bestPath[i+1])
		finalEdges = append(finalEdges, newEdge)
	}
	finalEdges = append(finalEdges, createEdge(bestPath[len(bestPath)-1], bestPath[0]))
	//treeOutput(finalEdges, "graph.txt")
}
