package main

func (coordinate *Coordinate) Translate(deltaX int, deltaY int) {
	coordinate.x = coordinate.x + deltaX
	coordinate.y = coordinate.y + deltaY
}

func Translate(coordinate Coordinate, deltaX int, deltaY int) Coordinate {
	return Coordinate{coordinate.x + deltaX, coordinate.y + deltaY}
}
