package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
	Handles .tsp file input
*/

type TSPProbelmEUC struct {
	name           string
	comment        string
	probType       string
	dimension      int
	edgeWeightType string
	initCities     []Node
}

func convertFloatXY(x string, y string) (float64, float64) {
	xf, err1 := strconv.ParseFloat(x, 64)
	if err1 != nil {
		panic(err1)
	}
	yf, err2 := strconv.ParseFloat(y, 64)
	if err2 != nil {
		panic(err2)
	}
	return xf, yf
}

func convInt(x string) int {
	xi, err1 := strconv.Atoi(x)
	if err1 != nil {
		panic(err1)
	}
	return xi
}

func readEucTSPFile(fp string) TSPProbelmEUC {
	f, err := os.Open(fp)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var pname, pcomment, ptype, pedgewtype string
	var pdimension int

	scanner := bufio.NewScanner(f)
	var path []Node
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "NAME") {
			parts := strings.Split(line, ":")
			pname = parts[1]
		} else if strings.Contains(line, "COMMENT") {
			parts := strings.Split(line, ":")
			pcomment = parts[1]
		} else if strings.Contains(line, "DIMENSION") {
			parts := strings.Split(line, ":")
			pdimension = convInt(strings.ReplaceAll(parts[1], " ", "")) //need to convert
		} else if strings.Contains(line, "TYPE") {
			parts := strings.Split(line, ":")
			ptype = parts[1] //need to convert
		} else if strings.Contains(line, "EDGE_WEIGHT_TYPE") {
			parts := strings.Split(line, ":")
			pedgewtype = parts[1]
		} else if strings.Contains(line, "NODE_COORD_SECTION") {
			fmt.Print("")
		} else if strings.Contains(line, "EOF") {
			break
		} else {
			parts := strings.Split(line, " ")
			x, y := convertFloatXY(parts[1], parts[2])
			newCity := Node{
				id: parts[0],
				x:  x,
				y:  y,
			}
			path = append(path, newCity)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	data := TSPProbelmEUC{
		name:           pname,
		comment:        pcomment,
		probType:       ptype,
		dimension:      pdimension,
		edgeWeightType: pedgewtype,
		initCities:     path,
	}
	return data
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
