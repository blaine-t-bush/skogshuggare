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

func (p *Point) ConvertToVector() Vector {
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

// Return the magnitude of the projection of v along the x-axis (i.e. its horizontal component)
func (v *Vector) X() float64 {
	x := v.len * math.Cos(v.dir)
	return x
}

// Return the magnitude of the projection of v along the y-axis (i.e. its vertical component)
func (v *Vector) Y() float64 {
	y := v.len * math.Sin(v.dir)
	return y
}

// Return a new point equal to the original point translated by the given vector.
func (p *Point) AddVector(v *Vector) Point {
	return Point{x: p.x + v.X(), y: p.y + v.Y()}
}

func main() {
	newPoint := Point{x: 1, y: -1}
	newVector := newPoint.ConvertToVector()
	fmt.Println("len:", newVector.len)
	fmt.Println("dir:", newVector.dir)
}
