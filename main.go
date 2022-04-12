package main

import (
	"fmt"
	"math"
)

type Point struct {
	x float64 // x-coordinate in units
	y float64 // y-coordinate in units
}

type Vector struct {
	// TODO add periodicity to dir so accepted values are 0 <= dir < 2*pi
	len float64 // Length (magnitude) in units
	dir float64 // Angle in radians, defined clockwise from the +x axis
}

func (p *Point) PointToVector() Vector {
	var dir, len float64
	len = math.Sqrt(math.Pow(p.x, 2) + math.Pow(p.y, 2))

	// Calculate dir in the principal value ranges of arccosine or arcsine
	if len == 0 {
		dir = 0
	} else {
		dir = math.Acos(p.x / len) // Principal value range:     0 <= dir <= pi
	}

	// If p is below the x-axis, arccosine returns the counterclockwise angle,
	// not the clockwise angle, so we have to flip it.
	if p.y < 0 {
		dir = 2*math.Pi - dir
	}

	return Vector{
		len: len,
		dir: dir,
	}
}

func main() {
	newPoint := Point{x: 1, y: -1}
	newVector := newPoint.PointToVector()
	fmt.Println("len:", newVector.len)
	fmt.Println("dir:", newVector.dir)
}
