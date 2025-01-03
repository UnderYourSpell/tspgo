// This example demonstrates a priority queue built using the heap interface.
package main

import "container/heap"

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Edge

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	//want the lowest priority
	return pq[i].wt < pq[j].wt
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	edge := x.(*Edge)
	edge.index = n
	*pq = append(*pq, edge)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
/*
func (pq *PriorityQueue) update(item *Edge, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
*/

func createEdgesPQ(nodes []Node, numCities int) map[string]PriorityQueue {
	edges := make(map[string]PriorityQueue)
	for i := range nodes {
		pq := make(PriorityQueue, numCities-1)
		indexer := 0
		for j := range nodes {
			if nodes[i].id == nodes[j].id {
				continue
			}
			distance := calcDistance(nodes[i], nodes[j])
			pq[indexer] = &Edge{
				origin: nodes[i],
				dest:   nodes[j],
				wt:     distance,
				index:  indexer,
			}
			indexer++
		}
		heap.Init(&pq)
		edges[nodes[i].id] = pq
	}
	return edges
}
