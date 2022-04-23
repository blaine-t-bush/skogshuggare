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

// Returns true if coordinate contains a collidable non-actor.
func (game *Game) IsBlocked(coordinate Coordinate) bool {
	if content, exists := game.world.content[coordinate]; exists {
		switch object := content.(type) {
		case Object:
			if object.collidable {
				return true
			}
		case *Tree:
			return true
		}
		return true
	}

	return false
}

// Returns true if coordinate contains an actor.
func (game *Game) IsOccupied(coordinate Coordinate) bool {
	if game.player.position == coordinate || game.squirrel.position == coordinate {
		return true
	}

	return false
}

// Returns true if coordinate contains a collidable non-actor, or an actor.
func (game *Game) IsBlockedOrOccupied(coordinate Coordinate) bool {
	return game.IsBlocked(coordinate) || game.IsOccupied(coordinate)
}

func (game *Game) GetRandomAvailableCoordinate() Coordinate {
	coordinate := Coordinate{rand.Intn(game.world.width), rand.Intn(game.world.height)}
	for {
		if game.IsBlocked(coordinate) {
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

func (game *Game) FindNextCoordinate(start Coordinate, end Coordinate) Coordinate {
	// Create a map of all coordinates.
	// We set the distance to all other nodes to "inf", except the distance
	// to the start node which we set at 0.
	graph := make(map[Coordinate]*Node, (2*game.world.height)*(2*game.world.width))
	for c := 0; c < game.world.width; c++ {
		for r := 0; r < game.world.height; r++ {
			// Only add non-collidable points, since actors cannot move throgh blocked points.
			if !game.IsBlocked(Coordinate{c, r}) {
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
	maxIterations := 1000
	// Now begin the algorithm:
	// While there are unvisited nodes left in the graph,
	// select an unknown node currentNode with the loewst distance.

	// and update neighboringNode's parent to currentNode.
	// Continue until there are no unvisited nodes left in the graph, or until max iterations is reached.

	for {
		// Select an unvisited node with the lowest distance.
		// If there are no more unvisited nodes, or we reach the maximum number of iterations, we're done.
		currentNode = SelectNextNode(graph)
		if currentNode == nil || iterations >= maxIterations {
			break
		}
		iterations++
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

	pathNode := graph[end]
	if pathNode.distance == Infinity {
		// Was unable to find a path.
		return start
	}
	for {
		if pathNode.parent == start {
			return pathNode.position // If success, return the coordinate of first step in optimal path.
		}
		pathNode = graph[pathNode.parent]
	}
}

func (game *Game) FindNextDirection(start Coordinate, end Coordinate) int {
	next := game.FindNextCoordinate(start, end)

	if start.x != next.x && start.y != next.y {
		panic("FindNextDirection: Next target coordinate is not adjacent to current coordinate!")
	}

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
