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
		x      float64
		y      float64
		vector Vector
	}

	tests := []test{
		{x: 0, y: 0, vector: Vector{x: 0, y: 0}},
		{x: 1, y: 0, vector: Vector{x: 1, y: 0}},
		{x: 0, y: 1, vector: Vector{x: 0, y: 1}},
		{x: 2, y: 2, vector: Vector{x: 2, y: 2}},
		{x: 2, y: -5.7, vector: Vector{x: 2, y: -5.7}},
		{x: 2.31, y: -2, vector: Vector{x: 2.31, y: -2}},
	} // FIXME dynamically generate Points instead of writing values twice.

	for _, testCase := range tests {
		if reflect.TypeOf(testCase.vector).Name() != "Vector" {
			t.Errorf("v is not of type Vector, got type %s", reflect.TypeOf(testCase.vector).Name())
		}

		if !EqualFloats(testCase.vector.x, testCase.x) {
			t.Errorf("v does not have x value %f, got value %f", testCase.x, testCase.vector.x)
		}

		if !EqualFloats(testCase.vector.y, testCase.y) {
			t.Errorf("v does not have y value %f, got value %f", testCase.y, testCase.vector.y)
		}
	}
}
func TestToVector(t *testing.T) {
	tests := []Point{
		{x: 0, y: 0},
		{x: 1, y: 0},
		{x: 0, y: 1},
		{x: 2, y: 2},
		{x: 2, y: -5.7},
		{x: 2.31, y: -2},
	} // FIXME dynamically generate Points instead of writing values twice.

	for _, testCase := range tests {
		got := testCase.ToVector()

		if reflect.TypeOf(got).Name() != "Vector" {
			t.Errorf("v1 is not of type Vector, got type %s", reflect.TypeOf(got).Name())
		}

		if !EqualFloats(got.x, testCase.x) {
			t.Errorf("v1 does not have x value %f, got value %f", testCase.x, got.x)
		}

		if !EqualFloats(got.y, testCase.y) {
			t.Errorf("v1 does not have y value %f, got value %f", testCase.y, got.y)
		}

	}
}

func TestTranslate(t *testing.T) {
	type test struct {
		point  Point
		vector Vector
		want   Point
	}

	tests := []test{
		{point: Point{x: 0, y: 0}, vector: Vector{x: 0, y: 0}, want: Point{x: 0, y: 0}},
		{point: Point{x: 0, y: 0}, vector: Vector{x: 1, y: 0}, want: Point{x: 1, y: 0}},
		{point: Point{x: 0, y: 0}, vector: Vector{x: 1, y: 7.2}, want: Point{x: 1, y: 7.2}},
	}

	for _, testCase := range tests {
		got := testCase.point.Translate(&testCase.vector)

		if !EqualFloats(testCase.want.x, got.x) {
			t.Errorf("point does not have x value %v, got value %v instead", testCase.want.x, got.x)
		}

		if !EqualFloats(testCase.want.y, got.y) {
			t.Errorf("point does not have y value %v, got value %v instead", testCase.want.y, got.y)
		}
	}
}

func TestLen(t *testing.T) {
	type test struct {
		vector Vector
		len    float64
	}

	tests := []test{
		{vector: Vector{x: 0, y: 0}, len: 0},
		{vector: Vector{x: 1, y: 0}, len: 1},
		{vector: Vector{x: 0, y: 2.5}, len: 2.5},
		{vector: Vector{x: -7.9, y: -7.9}, len: math.Sqrt(124.82)},
	}

	for _, testCase := range tests {
		got := testCase.vector.Len()

		if !EqualFloats(got, testCase.len) {
			t.Errorf("vector does not have length %f, got %f instead", testCase.len, got)
		}
	}
}

func TestDir(t *testing.T) {
	type test struct {
		vector Vector
		dir    float64
	}

	tests := []test{
		{vector: Vector{x: 0, y: 0}, dir: 0},
		{vector: Vector{x: 1, y: 0}, dir: 0},
		{vector: Vector{x: 0, y: 1}, dir: math.Pi / 2},
		{vector: Vector{x: -7.9, y: -7.9}, dir: 5 * math.Pi / 4},
	}

	for _, testCase := range tests {
		got := testCase.vector.Dir()

		if !EqualFloats(got, testCase.dir) {
			t.Errorf("vector does not have direction %f, got %f instead", testCase.dir, got)
		}
	}
}

func TestAdd(t *testing.T) {
	type test struct {
		vector1 Vector
		vector2 Vector
		want    Vector
	}

	tests := []test{
		{vector1: Vector{x: 0, y: 0}, vector2: Vector{x: 0, y: 0}, want: Vector{x: 0, y: 0}},
		{vector1: Vector{x: 0, y: 0}, vector2: Vector{x: 1, y: 0}, want: Vector{x: 1, y: 0}},
		{vector1: Vector{x: 0, y: 0}, vector2: Vector{x: 1, y: 7.2}, want: Vector{x: 1, y: 7.2}},
	}

	for _, testCase := range tests {
		got := testCase.vector1.Add(&testCase.vector2)

		if !EqualFloats(testCase.want.x, got.x) {
			t.Errorf("vector does not have x value %f, got value %f instead", testCase.want.x, got.x)
		}

		if !EqualFloats(testCase.want.y, got.y) {
			t.Errorf("vector does not have y value %f, got value %f instead", testCase.want.y, got.y)
		}
	}
}
