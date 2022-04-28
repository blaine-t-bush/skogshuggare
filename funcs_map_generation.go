package main

import (
	"math/rand"
	"os"
	"time"
)

// Generate river to place in map
func GenerateRiver() {

}

// Generate lake to place in map
func GenerateLake() {

}

// Generate island map
func GenerateIsland() {

}

func IsEdge(coord Coordinate, width int, height int) bool {
	return (coord.x == 0 || coord.x == width) || (coord.y == 0 || coord.y == height)
}

// Insert value at index in slice
func insert(a []byte, c byte, i int) []byte {
	return append(a[:i], append([]byte{c}, a[i:]...)...)
}

func GenerateBorders(width int, height int, mapData *[]byte) {

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if IsEdge(Coordinate{x, y}, width-1, height-1) {
				(*mapData)[To1DIndex(Coordinate{x, y}, width)] = '#'
			} else {
				(*mapData)[To1DIndex(Coordinate{x, y}, width)] = ' '
			}
		}
	}

	// Inser the newlines at the correct positions based on the width
	// width * (y + 1) + y
	for y := 0; y < height; y++ {
		*mapData = insert(*mapData, '\n', (width*(y+1) + y))
	}
}

func To2DIndex(point, width, height int) Coordinate {
	return Coordinate{point / width, point % height}

}

/*
Draws from (0,0) -> (0,1) -> (0,2)... etc
	|....
	v....
	 ....

*/
func To1DIndex(coord Coordinate, width int) int {
	return (coord.y * width) + coord.x
}

func WriteMapToFile(data []byte) {
	err := os.WriteFile("kartor/artificiell_skog.karta", data, 0644)
	if err != nil {
		panic(err)
	}
}

func PlacePlayer(width, height int, mapData *[]byte) {
	rand.Seed(time.Now().UnixNano())
	for {
		x := rand.Intn(width-1) + 1 // max - min + min
		y := rand.Intn(height-1) + 1
		coord := Coordinate{x, y}

		if (*mapData)[To1DIndex(coord, width)] == ' ' {
			(*mapData)[To1DIndex(coord, width)] = 'p'
			break
		}
	}
}

func PlaceSquirrel(width, height int, mapData *[]byte) {
	rand.Seed(time.Now().UnixNano())

	numberOfSquirrels := rand.Intn(10-1) + 1
	for i := 0; i < numberOfSquirrels; i++ {
		for {
			x := rand.Intn(width-1) + 1 // max - min + min
			y := rand.Intn(height-1) + 1
			coord := Coordinate{x, y}

			if (*mapData)[To1DIndex(coord, width)] == ' ' {
				(*mapData)[To1DIndex(coord, width)] = 's'
				break
			}
		}
	}
}

func GenerateMap(width int, height int, maptype string) {
	if maptype == "Island" {
		GenerateIsland()
	}

	mapData := make([]byte, width*height)

	GenerateBorders(width, height, &mapData)

	// Place the player and squirrels on the map after generating the terrain
	PlacePlayer(width, height, &mapData)
	//PlaceSquirrel(width, height, &mapData)

	/*for _, b := range mapData {
		fmt.Print(string(b))
	}*/

	WriteMapToFile(mapData)
}
