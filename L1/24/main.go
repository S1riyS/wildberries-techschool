package main

import (
	"fmt"
	"math"
)

type Point struct {
	x, y float64
}

func (p *Point) Distance(other Point) float64 {
	xDiff := p.x - other.x
	yDiff := p.y - other.y

	return math.Sqrt(xDiff*xDiff + yDiff*yDiff)
}

func NewPoint(x, y float64) *Point {
	return &Point{x: x, y: y}
}

func main() {
	A := NewPoint(1, 1)
	B := NewPoint(3, 3)
	fmt.Println(A.Distance(*B))
}
