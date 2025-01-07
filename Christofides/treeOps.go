package main

import (
	"container/heap"
	"math/rand"
)

func createMST(numCities int, nodes []Node, edges map[string][]Edge, graph map[string]Node) ([]Edge, map[string]Node) {
	var V []Node
	Vset := make(map[string]bool)     //track what vertices have been added to the tree
	randStart := rand.Intn(numCities) //random start
	V = append(V, nodes[randStart])
	Vset[nodes[randStart].id] = true

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
				if !ok { //if destination vertex not found in V, then we have found the best edge
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
	return tree, graph
}

// Create Minimum Cost Perfect Matching on resulting MST. Returns the order of nodes
func minCostPerfMatch(numCities int, nodes []Node, tree []Edge, graph map[string]Node, edges map[string][]Edge) []Node {
	var oddVertices []Node
	oddSet := make(map[string]bool)
	for i := range numCities {
		id := nodes[i].id
		node := graph[id]
		if len(node.edges)%2 == 0 { //add any vertex with an odd amount of vertices,  % 2 == 0 for this is actually the odd ones
			oddVertices = append(oddVertices, nodes[i])
			oddSet[nodes[i].id] = true
		}
	}

	/*
		Search edges of odd vertices that connect to other odd vertices and add the closest one
		Part of the reason we don't do much searching here is because all the edges are already sorted in ascending order
		so the search is pretty fast, it will always find the closest odd Vertex
	*/
	for i := range oddVertices {
		originNode := oddVertices[i]
		searchEdges := edges[oddVertices[i].id]
		for j := range searchEdges {
			destID := searchEdges[j].dest.id //ID of destination edge in set of all possible edges from originNode
			_, ok := oddSet[destID]
			if ok {
				//if destination in odd degree veritce set (oddSet)
				newEdge := createEdge(originNode, searchEdges[j].dest)
				tree = append(tree, newEdge)
				originNode.edges = append(originNode.edges, newEdge) //add new edge to origin node
				graph[originNode.id] = originNode                    //update edges for source Node to the graph
				delete(oddSet, originNode.id)                        //delete both from Oddset - cannot be used for any more connection
				delete(oddSet, searchEdges[j].dest.id)
				break
			}
		}

	}

	//This doesn't quite give us minimum weight perfect matching but its good enough
	//have a multigraph with duplicate edges, all vertices are of even degree

	//create an order of nodes in the path
	var nodeOrder []Node
	nodeOrder = append(nodeOrder, tree[0].origin, tree[0].dest) //add the first two in the tree
	for i := 1; i < len(tree); i++ {
		nodeOrder = append(nodeOrder, tree[i].dest) //for all others add the destination
	}

	return nodeOrder
}

func hamiltonianPath(nodeOrder []Node) []Edge {
	visited := make(map[string]bool)
	var finalGraph []Edge
	//we always visit the first vertex
	firstEdge := createEdge(nodeOrder[0], nodeOrder[1])
	finalGraph = append(finalGraph, firstEdge)
	visited[nodeOrder[1].id] = true

	for i := 2; i < len(nodeOrder); i++ {
		_, ok := visited[nodeOrder[i].id]
		if !ok { //if not visited
			//add the edge and mark visited
			newEdge := createEdge(nodeOrder[i-1], nodeOrder[i])
			finalGraph = append(finalGraph, newEdge)
			visited[nodeOrder[i].id] = true
		}
	}
	lastEdge := createEdge(finalGraph[len(finalGraph)-1].dest, nodeOrder[0])
	finalGraph = append(finalGraph, lastEdge)
	return finalGraph
}
