package main

/*
Prim's Algorithm to create a Minimum Spanning Tree on a TSP in Go
*/

import (
	"container/heap"
	"fmt"
	"math"
	"sort"
	"time"
)

type Graph struct {
	weight   float64
	vertices []Node
}

type Node struct {
	x     float64
	y     float64
	id    string
	cons  int
	edges []Edge
}

type Edge struct {
	origin Node
	dest   Node
	wt     float64
	index  int
}

func calcDistance(n1 Node, n2 Node) float64 {
	return math.Sqrt(math.Pow(n2.x-n1.x, 2) + math.Pow(n2.y-n1.y, 2))
}

func createEdges(nodes []Node) map[string][]Edge {
	//need to create a map with the city id
	edges := make(map[string][]Edge)
	for i := range nodes {
		var curEdges []Edge
		for j := range nodes {
			if nodes[i].id == nodes[j].id {
				continue
			}
			//the edge list for each node should be a priority queue this makes lookups faster
			distance := calcDistance(nodes[i], nodes[j])
			newEdge := Edge{
				origin: nodes[i],
				dest:   nodes[j],
				wt:     distance,
				index:  0,
			}
			curEdges = append(curEdges, newEdge)
		}
		sort.Slice(curEdges, func(i, j int) bool {
			return curEdges[i].wt < curEdges[j].wt
		})
		edges[nodes[i].id] = curEdges
	}
	return edges
}

func main() {
	//step 1. create a graph with all edges in the tree
	//read in file
	fp := "./tsp/fnl4461.tsp"
	data := readEucTSPFile(fp)

	//create list of edges and weights
	numCities := data.dimension
	nodes := data.initCities
	graph := make(map[string]Node) //master map
	for i := range numCities {
		graph[nodes[i].id] = nodes[i]
	}
	edges := createEdges(nodes)

	//Empty list of nodes and creating sets to track what in and whats not in the tree
	var V []Node
	V = append(V, nodes[0])
	Vset := make(map[string]bool)
	Vset[nodes[0].id] = true

	//Create MST
	startTreeCreate := time.Now() //start the clock - this is the meat of thw algorithm, everything else was set up
	var tree []Edge
	for j := 0; j < numCities; j++ {
		if len(V) == numCities {
			break
		}

		possibleEdges := make(PriorityQueue, len(V))
		for i := range V {
			curEdges := edges[V[i].id] //possible edges, this is a priority queue
			var edgeToAdd Edge
			for e := 0; e < numCities; e++ {
				edge := curEdges[e]
				_, ok := Vset[edge.dest.id]
				if !ok { //if other edge not found in V, then we have found the best edge
					edgeToAdd = edge
					break
				}
			}
			possibleEdges[i] = &edgeToAdd
		}
		heap.Init(&possibleEdges)
		bestEdge := heap.Pop(&possibleEdges).(*Edge)
		node := graph[bestEdge.origin.id]
		node.edges = append(node.edges, *bestEdge)
		graph[bestEdge.origin.id] = node //track what edges are connected to the vertex
		tree = append(tree, *bestEdge)
		V = append(V, bestEdge.dest) //can  maybe just use the set
		Vset[bestEdge.dest.id] = true
	}

	elapsedTree := time.Since(startTreeCreate)
	fmt.Printf("Runtime Create MST: %s\n", elapsedTree)

	//Now we do Matching
	//First we need to track connections
	//Should create a string map to a node to track its information
	//can have edges associated with the node
}
