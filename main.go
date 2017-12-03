package main

import (
	"flag"
	"mazey/maze"
	"os"
)

var width = flag.Int("w", 30, "Width of the maze")
var height = flag.Int("h", 20, "Height of the maze")

func main() {
	flag.Parse()
	m := maze.MakeMaze(*height, *width)
	maze.Solve(m)
	maze.PrintMaze(m, os.Stdout)
}
