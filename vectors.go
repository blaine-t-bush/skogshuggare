package main

import (
	"math"
)

type Point struct {
	x float64 // x-coordinate in arbitrary units
	y float64 // y-coordinate in arbitrary units
}

type Vector struct {
	x float64 // x component in arbitrary units
	y float64 // y component in arbitrary units
}

// Create a vector to a target point assuming origin at 0, 0.
func (p *Point) ToVector() Vector {
	return Vector{
		x: p.x,
		y: p.y,
	}
}

// Create a vector connecting an origin point to a target point.
func (p1 *Point) ToConnectingVector(p2 *Point) Vector {
	return Vector{
		x: p2.x - p1.x,
		y: p2.y - p1.y,
	}
}

// Return a new point equal to the original point translated by the given vector.
func (p *Point) Translate(v *Vector) Point {
	return Point{
		x: p.x + v.x,
		y: p.y + v.y,
	}
}

// Return the magnitude (length) of v.
func (v *Vector) Len() float64 {
	return math.Sqrt(math.Pow(v.x, 2) + math.Pow(v.y, 2))
}

// Return the direction of v in radians, defined clockwise from the +x-axis.
func (v *Vector) Dir() float64 {
	var dir float64
	len := v.Len()

	// Calculate dir in the principal value ranges of arccosine or arcsine
	if len == 0 {
		dir = 0
	} else {
		dir = math.Acos(v.x / len) // Principal value range: 0 <= dir <= pi
	}

	// If v points below the x-axis (assuming origin at 0, 0) then arccosine
	// returns the counterclockwise angle, not the clockwise angle, so we have
	// to flip it.
	if v.y < 0 {
		dir = 2*math.Pi - dir
	}

	return dir
}

// Return a new vector that is equal to the sum of two input vectors.
func (v1 *Vector) Add(v2 *Vector) Vector {
	return Vector{
		x: v1.x + v2.x,
		y: v1.y + v2.y,
	}
}
