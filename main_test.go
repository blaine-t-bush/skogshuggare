package main

import (
	"math"
	"reflect"
	"testing"
)

func TestPoint(t *testing.T) {
	var x, y float64
	x = 0
	y = 1
	p := Point{x: x, y: y}

	if reflect.TypeOf(p).Name() != "Point" {
		t.Errorf("p is not of type Point, got type %s", reflect.TypeOf(p).Name())
	}

	if p.x != x {
		t.Errorf("p does not have x value %f, got value %f", x, p.x)
	}

	if p.y != y {
		t.Errorf("p does not have y value %f, got value %f", y, p.y)
	}
}

func TestVector(t *testing.T) {
	var len, dir float64
	len = 3
	dir = 6
	v := Vector{len: len, dir: dir}

	if reflect.TypeOf(v).Name() != "Vector" {
		t.Errorf("v is not of type Vector, got type %s", reflect.TypeOf(v).Name())
	}

	if v.len != len {
		t.Errorf("v does not have len value %f, got value %f", len, v.len)
	}

	if v.dir != dir {
		t.Errorf("v does not have dir value %f, got value %f", dir, v.dir)
	}
}
func TestPointToVector(t *testing.T) {
	// Test point (0, 1)
	var x1, y1, len1, dir1 float64
	x1 = 0
	y1 = 1
	len1 = 1
	dir1 = math.Pi / 2
	p1 := Point{x: x1, y: y1}
	v1 := p1.PointToVector()

	if reflect.TypeOf(v1).Name() != "Vector" {
		t.Errorf("v1 is not of type Vector, got type %s", reflect.TypeOf(v1).Name())
	}

	if v1.len != len1 {
		t.Errorf("v1 does not have len value %f, got value %f", len1, v1.len)
	}

	if v1.dir != dir1 {
		t.Errorf("v1 does not have dir value %f, got value %f", dir1, v1.dir)
	}

	// Test point (-2, -2)
	var x2, y2, len2, dir2 float64
	x2 = -2
	y2 = -2
	len2 = math.Sqrt(8)
	dir2 = 5 * math.Pi / 4
	p2 := Point{x: x2, y: y2}
	v2 := p2.PointToVector()

	if reflect.TypeOf(v2).Name() != "Vector" {
		t.Errorf("v2 is not of type Vector, got type %s", reflect.TypeOf(v2).Name())
	}

	if v2.len != len2 {
		t.Errorf("v2 does not have len value %f, got value %f", len2, v2.len)
	}

	if v2.dir != dir2 {
		t.Errorf("v2 does not have dir value %f, got value %f", dir2, v2.dir)
	}
}
