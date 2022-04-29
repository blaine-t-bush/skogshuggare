package main

func (coordinate *Coordinate) Translate(deltaX int, deltaY int) {
	coordinate.x = coordinate.x + deltaX
	coordinate.y = coordinate.y + deltaY
}

func Translate(coordinate Coordinate, deltaX int, deltaY int) Coordinate {
	return Coordinate{coordinate.x + deltaX, coordinate.y + deltaY}
}

func TranslateByDir(coordinate Coordinate, dir int, len int) Coordinate {
	switch dir {
	case DirUp:
		return Translate(coordinate, 0, -len)
	case DirRight:
		return Translate(coordinate, len, 0)
	case DirDown:
		return Translate(coordinate, 0, len)
	case DirLeft:
		return Translate(coordinate, -len, 0)
	default:
		return coordinate
	}
}
