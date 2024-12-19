package main

import (
	"fmt"
	"math"
)

type City struct {
	x  float64
	y  float64
	id string
}

type Triper interface {
	getDistance()
	calcPathLength()
}

type Trip struct {
	path       []City
	pathLength float32
	prob       float64
}

func (x *Trip) getDistance(city1 City, city2 City) float64 {
	x1 := city1.x
	y1 := city1.y
	x2 := city2.x
	y2 := city2.y
	distance := math.Sqrt(math.Pow((x2-x1), 2) + math.Pow((y2-y1), 2))
	return distance
}

func (x *Trip) calcPathLength() {
	pathL := 0.0
	N := len(x.path)
	for i := 0; i < N; i++ {
		if i+1 != N {
			pathL += x.getDistance(x.path[i], x.path[i+1])
		} else {
			pathL += x.getDistance(x.path[i], x.path[0])
		}
	}
	x.pathLength = float32(pathL)
}

func (x *Trip) printPath() {
	for i := 0; i < len(x.path); i++ {
		fmt.Print(x.path[i].id)
		fmt.Print(" ")
	}
	fmt.Println()
}
