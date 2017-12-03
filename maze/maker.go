package maze

import (
	"math/rand"
	"time"
)

func MakeMaze(rows, cols int) *Maze {
	rand.Seed(time.Now().Unix())
	m := NewMaze(rows, cols)
	startRow := rand.Intn(rows)
	startCol := rand.Intn(cols)
	openCells(m.CellAt(startRow, startCol))

	m.CellAt(rand.Intn(rows), 0).OpenDoor(Left)
	m.CellAt(rand.Intn(rows), cols-1).OpenDoor(Right)
	return m
}

type neighborInDirection struct {
	neighbor  Cell
	direction Direction
}

func availableNeighbors(c Cell) []neighborInDirection {
	neighbors := make([]neighborInDirection, 0, 4)
	for d := range AllDirections() {
		neighbor := c.Go(d)
		if neighbor.IsInMaze() {
			neighbors = append(neighbors, neighborInDirection{neighbor, d})
		}
	}
	return neighbors
}

func openCells(c Cell) {
	c.SetMark(1)

	n := availableNeighbors(c)
	for _, i := range rand.Perm(len(n)) {
		if n[i].neighbor.Mark() == 0 {
			c.OpenDoor(n[i].direction)
			openCells(n[i].neighbor)
		}
	}
}
