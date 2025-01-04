package main

/*
Christofides
*/

import (
	"container/heap"
	"fmt"
	"math"
	"sort"
	"time"
)

type Node struct {
	x     float64
	y     float64
	id    string
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
	fp := "./tsp/original10.tsp"
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
	start := time.Now() //start the clock for the tree - this is the meat of thw algorithm, everything else was set up
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

	/*
		Create Minimum Cost Perfect Matching
		Edges added to Tree and graph map
	*/

	//Add odd Vertices to an Array
	var oddVertices []Node
	oddSet := make(map[string]bool)
	for i := range numCities {
		id := nodes[i].id
		node := graph[id]
		if len(node.edges)%2 == 0 {
			oddVertices = append(oddVertices, nodes[i])
			oddSet[nodes[i].id] = true
		}
	}

	//Search edges of odd vertices that connect to other odd vertices and add the closest one
	//Part of the reason we don't do much searching here is because all the edges are already sorted in ascending order
	//so the search is pretty fast, it will always find the closest odd Vertex
	for i := range oddVertices {
		originNode := oddVertices[i]
		searchEdges := edges[oddVertices[i].id]
		for j := range searchEdges {
			destID := searchEdges[j].dest.id //ID of destination edge in set of all possible edges from originNode
			_, ok := oddSet[destID]
			if ok {
				//if destination in odd set
				newEdge := Edge{
					origin: originNode,
					dest:   searchEdges[j].dest,
					wt:     calcDistance(originNode, searchEdges[j].dest),
					index:  0,
				}
				//lets add the edge to the tree for visualization purposes
				tree = append(tree, newEdge)
				originNode.edges = append(originNode.edges, newEdge) //add edge to new Node
				graph[originNode.id] = originNode                    //update edges for source Node to the graph
				//delete both from Oddset - cannot be used
				delete(oddSet, originNode.id)
				delete(oddSet, searchEdges[j].dest.id)
				break
			}
		}

	}

	//its not minimum weight lol
	//have a multigraph with duplicate edges, all vertices are of even degree
	//take the Eulerian tour
	//traverse vertices visit only once

	//need an order of vertices from the tree edges
	//create an order of nodes
	var nodeOrder []Node
	nodeOrder = append(nodeOrder, tree[0].origin, tree[0].dest) //add the first two in the tree
	for i := 1; i < len(tree); i++ {
		nodeOrder = append(nodeOrder, tree[i].dest) //for all others add the destination
	}

	visited := make(map[string]bool)
	var finalGraph []Edge
	//we always visit the first guy
	firstEdge := Edge{
		origin: nodeOrder[0],
		dest:   nodeOrder[1],
		wt:     calcDistance(nodeOrder[0], nodeOrder[1]),
		index:  0,
	}
	finalGraph = append(finalGraph, firstEdge)
	visited[nodeOrder[1].id] = true

	for i := 2; i < len(nodeOrder); i++ {
		_, ok := visited[nodeOrder[i].id]
		if !ok { //if not visited
			//add the edge and mark visited
			newEdge := Edge{
				origin: nodeOrder[i-1],
				dest:   nodeOrder[i],
				wt:     calcDistance(nodeOrder[i-1], nodeOrder[i]),
				index:  0,
			}
			finalGraph = append(finalGraph, newEdge)
			visited[nodeOrder[i].id] = true
		}
	}

	lastEdge := Edge{
		origin: finalGraph[len(finalGraph)-1].dest,
		dest:   nodeOrder[0],
		wt:     calcDistance(nodeOrder[0], nodeOrder[1]),
		index:  0,
	}
	finalGraph = append(finalGraph, lastEdge)

	//now calculate final path length
	var pathLength float64
	for i := range finalGraph {
		pathLength += finalGraph[i].wt
	}
	fmt.Println("Path Length:", pathLength)
	elapsed := time.Since(start)
	fmt.Printf("Runtime Create MST: %s\n", elapsed)

	//next step is local search
	//subsequent graph is 2opt optimal
	//can use go routines for 2opt to do simulated annealing

	//treeOutput(finalGraph, "graph.txt")
}
