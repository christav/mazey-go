package main

import (
	"flag"
	"mazey/maze"
	"os"
)

var width = flag.Int("w", 30, "Width of the maze")
var height = flag.Int("h", 20, "Height of the maze")
var ascii = flag.Bool("a", false, "Print maze as ASCII")
var noSolve = flag.Bool("n", false, "Don't solve the maze")

func main() {
	flag.Parse()
	m := maze.MakeMaze(*height, *width)
	if !*noSolve {
		maze.Solve(m)
	}
	maze.PrintMaze(m, *ascii, os.Stdout)
}
