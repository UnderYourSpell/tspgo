package main

import "math"

/*
	File Containing necessary struct definitions and helper functions pertaining to the structs
*/

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

func calcPathLength(nodes []Node) float64 {
	var length float64
	for i := 0; i < len(nodes)-1; i++ {
		length += calcDistance(nodes[i], nodes[i+1])
	}
	length += calcDistance(nodes[len(nodes)-1], nodes[0])
	return length
}

func createEdge(origin Node, dest Node) Edge {
	newEdge := Edge{
		origin: origin,
		dest:   dest,
		wt:     calcDistance(origin, dest),
		index:  0,
	}
	return newEdge
}

func pathLengthEdges(edges []Edge) float64 {
	var pathLength float64
	for i := range edges {
		pathLength += edges[i].wt
	}
	return pathLength
}

func createPath(edges []Edge) []Node {
	var path []Node
	path = append(path, edges[0].origin, edges[0].dest)
	for i := 1; i < len(edges)-1; i++ {
		path = append(path, edges[i].dest)
	}
	return path
}
