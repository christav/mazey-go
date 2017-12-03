package main

import (
	"mazey/maze"
	"os"
)

func main() {
	m := maze.MakeMaze(20, 30)
	maze.Solve(m)
	maze.PrintMaze(m, os.Stdout)
}
