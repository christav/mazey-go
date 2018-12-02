package maze

import (
	"fmt"
	"io"
)

// CharSet represents the characters used to print the maze
type CharSet struct {
	corners       [16]rune
	solutionChars [16]string
}

// UnicodeCharSet is a CharSet using unicode line drawing characters
var UnicodeCharSet = CharSet{
	corners: [...]rune{
		' ', '╹', '╺', '┗', '╻', '┃', '┏', '┣',
		'╸', '┛', '━', '┻', '┓', '┫', '┳', '╋'},

	solutionChars: [...]string{
		"   ", "   ", "   ", " ╰┄",
		"   ", " ┆ ", " ╭┄", "   ",
		"   ", "┄╯ ", "┄┄┄", "   ",
		"┄╮ ", "   ", "   ", "   "},
}

// ASCIICharSet is a Charset using only basic ASCII characters
var ASCIICharSet = CharSet{
	corners: [...]rune{
		' ', '+', '+', '+',
		'+', '|', '+', '+',
		'+', '+', '-', '+',
		'+', '+', '+', '+'},

	solutionChars: [...]string{
		"   ", "   ", "   ", " XX",
		"   ", " X ", " XX", "   ",
		"   ", "XX ", "XXX", "   ",
		"XX ", "   ", "   ", "   "},
}

type printer struct {
	out           io.Writer
	charSet       *CharSet
	horizontalBar string
}

// PrintMaze prints the maze to the given output using either unicode or ascii charsets
func PrintMaze(m *Maze, printASCII bool, out io.Writer) {
	if printASCII {
		printASCIIMaze(m, out)
	} else {
		printUnicodeMaze(m, out)
	}
}

func printUnicodeMaze(m *Maze, out io.Writer) {
	PrintMazeWithCharSet(m, out, &UnicodeCharSet)
}

func printASCIIMaze(m *Maze, out io.Writer) {
	PrintMazeWithCharSet(m, out, &ASCIICharSet)
}

// PrintMazeWithCharSet prints the maze to the given output using the given CharSet
func PrintMazeWithCharSet(m *Maze, out io.Writer, charSet *CharSet) {
	bar := string(charSet.corners[10])
	horizontalBar := bar + bar + bar
	p := printer{out, charSet, horizontalBar}
	p.print(m)
}

func (p *printer) print(m *Maze) {
	rowIter := m.AllRows()
	for r, ok := rowIter.Next(); ok; r, ok = rowIter.Next() {
		p.printRowSeparator(*r)
		(*r).Reset()
		p.printRow(*r)
	}
	p.printMazeBottom(m)
}

func (p *printer) printRowSeparator(row CellIterator) {
	var c Cell
	var ok bool
	for c, ok = row.Next(); ok; c, ok = row.Next() {
		fmt.Fprintf(p.out, "%c", p.cornerChar(c))
		if c.CanGo(Up) {
			if IsSolutionCell(c) && IsSolutionCell(c.Go(Up)) {
				fmt.Fprintf(p.out, "%s", p.charSet.solutionChars[5])
			} else {
				fmt.Fprintf(p.out, "   ")
			}
		} else {
			fmt.Fprintf(p.out, p.horizontalBar)
		}
	}
	fmt.Fprintf(p.out, "%c\n", p.rowSeparatorEnd(c))
}

func (p *printer) cornerChar(c Cell) rune {
	neighbors := [...]Cell{c.Go(Up), c.Go(Left)}
	index := 0
	if !(neighbors[0].IsInMaze() && neighbors[0].CanGo(Left)) {
		index |= 1
	}
	if !c.CanGo(Up) {
		index |= 2
	}
	if !(c.IsEntrance() || c.CanGo(Left)) {
		index |= 4
	}

	if !(neighbors[1].IsInMaze() && neighbors[1].CanGo(Up)) {
		index |= 8
	}

	if c.Row() == 0 {
		index &= 0xe
	}

	if c.Col() == 0 {
		index &= 0x7
	}
	return p.charSet.corners[index]
}

func (p *printer) rowSeparatorEnd(cell Cell) rune {
	upCell := cell.Go(Up)
	index := 0
	if !(upCell.IsInMaze() && upCell.IsExit()) {
		index |= 1
	}

	if !(cell.IsExit()) {
		index |= 4
	}
	if !(cell.CanGo(Up)) {
		index |= 8
	}

	if cell.Row() == 0 {
		index &= 0xe
	}
	return p.charSet.corners[index]
}

func (p *printer) printRow(row CellIterator) {
	var c Cell
	var ok bool
	for c, ok = row.Next(); ok; c, ok = row.Next() {
		if c.IsEntrance() {
			if IsSolutionCell(c) {
				fmt.Fprintf(p.out, "%c", []rune(p.charSet.solutionChars[10])[1])
			} else {
				fmt.Fprintf(p.out, " ")
			}
		} else if c.CanGo(Left) {
			if IsSolutionCell(c) && IsSolutionCell(c.Go(Left)) {
				fmt.Fprintf(p.out, "%c", []rune(p.charSet.solutionChars[10])[1])
			} else {
				fmt.Fprintf(p.out, " ")
			}
		} else {
			fmt.Fprintf(p.out, "%c", p.charSet.corners[5])
		}
		fmt.Fprintf(p.out, p.cellContents(c))
	}

	lastCell := c
	if lastCell.IsExit() {
		if IsSolutionCell(lastCell) {
			fmt.Fprintf(p.out, "%c", []rune(p.charSet.solutionChars[10])[1])
		} else {
			fmt.Fprintf(p.out, " ")
		}
	} else {
		fmt.Fprintf(p.out, "%c", p.charSet.corners[5])
	}
	fmt.Fprintln(p.out)
}

func (p *printer) cellContents(c Cell) string {
	if !IsSolutionCell(c) {
		return "   "
	}

	index := 0
	if c.CanGo(Up) && IsSolutionCell(c.Go(Up)) {
		index |= 1
	}

	if c.IsExit() || c.CanGo(Right) && IsSolutionCell(c.Go(Right)) {
		index |= 2
	}

	if c.CanGo(Down) && IsSolutionCell(c.Go(Down)) {
		index |= 4
	}

	if c.IsEntrance() || c.CanGo(Left) && IsSolutionCell(c.Go(Left)) {
		index |= 8
	}
	return p.charSet.solutionChars[index]
}

func (p *printer) printMazeBottom(m *Maze) {
	lastRow := m.Row(m.Rows() - 1)
	var c Cell
	var ok bool
	for c, ok = lastRow.Next(); ok; c, ok = lastRow.Next() {
		index := 0xa
		if !c.CanGo(Left) {
			index |= 1
		}
		if c.Col() == 0 {
			index &= 0x7
		}
		fmt.Fprintf(p.out, "%c", p.charSet.corners[index])
		fmt.Fprintf(p.out, p.horizontalBar)
	}
	lastCell := c // lastRow[len(lastRow)-1]
	index := 0x8
	if !lastCell.CanGo(Right) {
		index |= 1
	}
	fmt.Fprintf(p.out, "%c\n", p.charSet.corners[index])
}
