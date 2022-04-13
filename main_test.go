package main

import (
	"math"
	"reflect"
	"testing"
)

// Helper function for comparing closeness of float64 values.
// Returns true if the difference between two numbers is within a small threshold value,
// and false otherwise.
func EqualFloats(a, b float64) bool {
	return math.Abs(a-b) <= 1e-10
}

func TestPoint(t *testing.T) {
	type test struct {
		x     float64
		y     float64
		point Point
	}

	tests := []test{
		{x: 0, y: 0, point: Point{x: 0, y: 0}},
		{x: 1, y: 0, point: Point{x: 1, y: 0}},
		{x: 0, y: 1, point: Point{x: 0, y: 1}},
		{x: 2, y: 2, point: Point{x: 2, y: 2}},
		{x: 2, y: -5.7, point: Point{x: 2, y: -5.7}},
		{x: 2.31, y: -2, point: Point{x: 2.31, y: -2}},
	} // FIXME dynamically generate Points instead of writing values twice.

	for _, testCase := range tests {
		if reflect.TypeOf(testCase.point).Name() != "Point" {
			t.Errorf("p is not of type Point, got type %s", reflect.TypeOf(testCase.point).Name())
		}

		if !EqualFloats(testCase.point.x, testCase.x) {
			t.Errorf("p does not have x value %f, got value %f", testCase.x, testCase.point.x)
		}

		if !EqualFloats(testCase.point.y, testCase.y) {
			t.Errorf("p does not have y value %f, got value %f", testCase.y, testCase.point.y)
		}
	}
}

func TestVector(t *testing.T) {
	type test struct {
		len    float64
		dir    float64
		vector Vector
	}

	tests := []test{
		{len: 0, dir: 0, vector: Vector{len: 0, dir: 0}},
		{len: 1, dir: 7, vector: Vector{len: 1, dir: 7}},
		{len: 99999.9999, dir: math.Pi, vector: Vector{len: 99999.9999, dir: math.Pi}},
		{len: -5.2, dir: -9e-5, vector: Vector{len: -5.2, dir: -9e-5}},
		{len: -9e-5, dir: 2, vector: Vector{len: -9e-5, dir: 2}},
		{len: 7.5e7, dir: 3, vector: Vector{len: 7.5e7, dir: 3}},
		{len: 2, dir: 5, vector: Vector{len: 2, dir: 5}},
	}

	for _, testCase := range tests {
		if reflect.TypeOf(testCase.vector).Name() != "Vector" {
			t.Errorf("v is not of type Vector, got type %s", reflect.TypeOf(testCase.vector).Name())
		}

		if !EqualFloats(testCase.vector.len, testCase.len) {
			t.Errorf("v does not have len value %f, got value %f", testCase.len, testCase.vector.len)
		}

		if !EqualFloats(testCase.vector.dir, testCase.dir) {
			t.Errorf("v does not have dir value %f, got value %f", testCase.dir, testCase.vector.dir)
		}
	}
}
func TestConvertToVector(t *testing.T) {
	type test struct {
		point Point
		len   float64
		dir   float64
	}

	tests := []test{
		{point: Point{x: 0, y: 1}, len: 1, dir: math.Pi / 2},
		{point: Point{x: -2, y: -2}, len: math.Sqrt(8), dir: 5 * math.Pi / 4},
	}

	for _, testCase := range tests {
		got := testCase.point.ConvertToVector()

		if reflect.TypeOf(got).Name() != "Vector" {
			t.Errorf("v1 is not of type Vector, got type %s", reflect.TypeOf(got).Name())
		}

		if !EqualFloats(got.len, testCase.len) {
			t.Errorf("v1 does not have len value %f, got value %f", testCase.len, got.len)
		}

		if !EqualFloats(got.dir, testCase.dir) {
			t.Errorf("v1 does not have dir value %f, got value %f", testCase.dir, got.dir)
		}

	}
}

func TestAddVector(t *testing.T) {
	type test struct {
		point  Point
		vector Vector
		want   Point
	}

	tests := []test{
		{point: Point{x: 0, y: 0}, vector: Vector{len: 0, dir: 0}, want: Point{x: 0, y: 0}},
		{point: Point{x: 0, y: 0}, vector: Vector{len: 1, dir: 0}, want: Point{x: 1, y: 0}},
		{point: Point{x: 0, y: 0}, vector: Vector{len: 1, dir: math.Pi / 2}, want: Point{x: 0, y: 1}},
	}

	for _, testCase := range tests {
		got := testCase.point.AddVector(&testCase.vector)
		if !EqualFloats(testCase.want.x, got.x) || !EqualFloats(testCase.want.y, got.y) {
			t.Errorf("point does not have value %v, got value %v instead", testCase.want, got)
		}
	}
}
