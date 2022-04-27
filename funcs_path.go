package main

import (
	"math/rand"
)

type Node struct {
	position  Coordinate
	parent    Coordinate
	neighbors []Coordinate
	distance  int
	visited   bool
}

const (
	Infinity = 9999999
)

// Returns true if coordinate contains a collidable object, or a tree.
func (game *Game) IsBlocked(coordinate Coordinate) bool {
	if content, exists := game.world.content[coordinate]; exists {
		switch content := content.(type) {
		case Object:
			if content.collidable {
				return true
			}
		case *Tree:
			return true
		}
	}

	return false
}

// Returns true if coordinate contains a collidable object, or a tree, or fire.
func (game *Game) IsPathBlocked(coordinate Coordinate) bool {
	if content, exists := game.world.content[coordinate]; exists {
		switch content := content.(type) {
		case Object:
			if content.collidable {
				return true
			}
		case *Fire:
			return true
		case *Tree:
			return true
		}
	}

	return false
}

// Returns true if coordinate contains a collidable or non-plantable object, or a tree.
func (game *Game) IsUnplantable(coordinate Coordinate) bool {
	if content, exists := game.world.content[coordinate]; exists {
		switch content := content.(type) {
		case Object:
			if content.collidable || !content.plantable {
				return true
			}
		case *Tree:
			return true
		}
	}

	return false
}

// Returns true if coordinate contains a collidable object, or fire.
func (game *Game) IsUnflammable(coordinate Coordinate) bool {
	if content, exists := game.world.content[coordinate]; exists {
		switch content := content.(type) {
		case Object:
			if !content.flammable {
				return true
			}
		case *Fire:
			return true
		}
	}

	return false
}

func (game *Game) GetRandomAvailableCoordinate() Coordinate {
	coordinate := Coordinate{rand.Intn(game.world.width), rand.Intn(game.world.height)}
	iterations := 0
	for {
		if iterations >= MaxIterations {
			panic("Reached max iterations in GetRandomAvailableCoordinate()")
		}
		iterations++

		if game.IsPathBlocked(coordinate) {
			coordinate = Coordinate{rand.Intn(game.world.width), rand.Intn(game.world.height)}
		} else {
			break
		}
	}

	return coordinate
}

func (game *Game) GetRandomFlammableCoordinate() Coordinate {
	coordinate := Coordinate{rand.Intn(game.world.width), rand.Intn(game.world.height)}
	iterations := 0
	for {
		if iterations >= MaxIterations {
			panic("Reached max iterations in GetRandomFlammableCoordinate()")
		}
		iterations++

		if game.IsUnflammable(coordinate) {
			coordinate = Coordinate{rand.Intn(game.world.width), rand.Intn(game.world.height)}
		} else {
			break
		}
	}

	return coordinate
}

func (game *Game) GetRandomPlantableCoordinate() Coordinate {
	coordinate := Coordinate{rand.Intn(game.world.width), rand.Intn(game.world.height)}
	iterations := 0
	for {
		if iterations >= MaxIterations {
			panic("Reached max iterations in GetRandomPlantableCoordinate()")
		}
		iterations++

		if game.IsUnplantable(coordinate) {
			coordinate = Coordinate{rand.Intn(game.world.width), rand.Intn(game.world.height)}
		} else {
			break
		}
	}

	return coordinate
}

func GetRandomDirection() int {
	randInt := rand.Intn(4)
	switch randInt {
	case 0:
		return DirUp
	case 1:
		return DirRight
	case 2:
		return DirDown
	case 3:
		return DirLeft
	default:
		return DirNone
	}
}

func (node *Node) GetNeighbors(graph map[Coordinate]*Node) []Coordinate {
	// Get slice of neighbors based, but only ones that are within the world dimensions.
	var neighbors []Coordinate
	if _, exists := graph[Coordinate{node.position.x, node.position.y - 1}]; exists {
		neighbors = append(neighbors, Coordinate{node.position.x, node.position.y - 1})
	}
	if _, exists := graph[Coordinate{node.position.x + 1, node.position.y}]; exists {
		neighbors = append(neighbors, Coordinate{node.position.x + 1, node.position.y})
	}
	if _, exists := graph[Coordinate{node.position.x, node.position.y + 1}]; exists {
		neighbors = append(neighbors, Coordinate{node.position.x, node.position.y + 1})
	}
	if _, exists := graph[Coordinate{node.position.x - 1, node.position.y}]; exists {
		neighbors = append(neighbors, Coordinate{node.position.x - 1, node.position.y})
	}

	return neighbors
}

func SelectNextNode(graph map[Coordinate]*Node) *Node {
	var nextNode *Node
	lowestDistance := Infinity

	// Determine the lowest distance among unvisited nodes.
	for _, node := range graph {
		if !node.visited && node.distance <= lowestDistance {
			lowestDistance = node.distance
		}
	}

	// Select any unvisited node with the lowest distance.
	for _, node := range graph {
		if !node.visited && node.distance == lowestDistance {
			nextNode = node
			break
		}
	}

	return nextNode
}

func (game *Game) FindPath(start Coordinate, end Coordinate) map[int]Coordinate {
	// Create a map of all coordinates.
	// We set the distance to all other nodes to "inf", except the distance
	// to the start node which we set at 0.
	graph := make(map[Coordinate]*Node, (2*game.world.height)*(2*game.world.width))
	for c := 0; c < game.world.width; c++ {
		for r := 0; r < game.world.height; r++ {
			// Only add non-collidable points, since actors cannot move throgh blocked points.
			if !game.IsPathBlocked(Coordinate{c, r}) {
				graph[Coordinate{c, r}] = &Node{position: Coordinate{c, r}, distance: Infinity, visited: false}
			}
		}
	}

	// Add neighbors to each node of the graph.
	for _, node := range graph {
		node.neighbors = node.GetNeighbors(graph)
	}

	// Initialize the starting coordinate as the currentNode.
	graph[start].parent = start
	graph[start].distance = 0
	var currentNode *Node
	iterations := 0
	// Now begin the algorithm:
	// While there are unvisited nodes left in the graph,
	// select an unknown node currentNode with the loewst distance.

	// and update neighboringNode's parent to currentNode.
	// Continue until there are no unvisited nodes left in the graph, or until max iterations is reached.

	for {
		if iterations >= MaxIterations {
			panic("Reached max iterations in FindPath()")
		}
		iterations++

		// Select an unvisited node with the lowest distance.
		// If there are no more unvisited nodes, or we reach the maximum number of iterations, we're done.
		currentNode = SelectNextNode(graph)
		if currentNode == nil {
			break
		}
		currentNode.visited = true

		// For each node neighboring currentNode,
		// if currentNode distance + distance from currentNode to neighbor is less than neighbor's distance,
		// update neighbor's distance to currentNode distance + distance from currentNode to neighbor.
		// Also mark currentNode as neighbor's parent, so it knows how to follow the shortest path.
		for _, neighbor := range currentNode.neighbors {
			if currentNode.distance+1 < graph[neighbor].distance {
				graph[neighbor].distance = currentNode.distance + 1
				graph[neighbor].parent = currentNode.position
			}
		}
	}

	// Create map where key is an int representing the path order (e.g. key = 1 is first step, key = 2 is second step)
	// and value is the target coordinate (e.g. value for key = 1 is the coordinate for first step in path).
	pathNode, exists := graph[end]
	path := make(map[int]Coordinate)
	if !exists || pathNode.distance == Infinity {
		// Was unable to find a path.
		path[1] = start
		return path
	}

	stepCount := 1
	for {
		if pathNode.parent == start {
			break
		} else {
			pathNode = graph[pathNode.parent]
		}
		stepCount++
	}

	pathNode = graph[end]
	for i := stepCount; i > 0; i-- {
		path[i] = pathNode.position
		pathNode = graph[pathNode.parent]
	}

	return path
}

func (game *Game) FindNextDirection(squirrelKey int) int {
	// Loop through current path and check if any are blocked. If not, move to next coord in path.
	// Otherwise, run pathfinding again.
	squirrel := game.squirrels[squirrelKey]
	findNewPath := false
	for _, coord := range squirrel.path {
		if game.IsPathBlocked(coord) {
			findNewPath = true
			break
		}
	}

	if findNewPath {
		squirrel.path = game.FindPath(squirrel.position, squirrel.destination)
	}

	start := squirrel.position
	next := squirrel.path[1]
	var dir int
	if start.y > next.y {
		dir = DirUp
	} else if start.x < next.x {
		dir = DirRight
	} else if start.y < next.y {
		dir = DirDown
	} else if start.x > next.x {
		dir = DirLeft
	} else if start.x == next.x && start.y == next.y {
		// This is true when FindNextCoordinate was unable to find a path,
		// or start and end are equal.
		// Actors that use this pathfinding function should update their destinations
		// if this value is received.
		dir = DirNone
	} else {
		panic("FindNextDirection: Unable to determine proper direction!")
	}

	return dir
}

func (game *Game) UpdatePath(squirrelKey int) {
	squirrel := game.squirrels[squirrelKey]
	newPath := make(map[int]Coordinate)
	for key, coord := range squirrel.path {
		if key != 1 {
			newPath[key-1] = coord
		}
	}

	squirrel.path = newPath
}
