package maze

import (
	"fmt"
	"io"
)

type CharSet struct {
	corners       [16]rune
	solutionChars [16]string
}

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

var ASCIICharSet = CharSet{
	corners: [...]rune{
		' ', '+', '+', '+',
		'+', '|', '+', '+',
		'+', '+', '-', '+',
		'+', '+', '+', '+'},

	solutionChars: [...]string{
		"   ", "   ", "   ", "XXX",
		"   ", "XXX", "XXX", "   ",
		"   ", "XXX", "XXX", "   ",
		"XXX", "   ", "   ", "   "},
}

type printer struct {
	out           io.Writer
	charSet       *CharSet
	horizontalBar string
}

func PrintMaze(m *Maze, out io.Writer) {
	PrintUnicodeMaze(m, out)
}

func PrintUnicodeMaze(m *Maze, out io.Writer) {
	PrintMazeWithCharSet(m, out, &UnicodeCharSet)
}

func PrintASCIIMaze(m *Maze, out io.Writer) {
	PrintMazeWithCharSet(m, out, &ASCIICharSet)
}

func PrintMazeWithCharSet(m *Maze, out io.Writer, charSet *CharSet) {
	bar := string(charSet.corners[10])
	horizontalBar := bar + bar + bar
	p := printer{out, charSet, horizontalBar}
	p.Print(m)
}

func (p *printer) Print(m *Maze) {
	for _, row := range m.AllRows() {
		p.printRowSeparator(row)
		p.printRow(row)
	}
	p.printMazeBottom(m)
}

func (p *printer) printRowSeparator(row []Cell) {
	for _, c := range row {
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
	fmt.Fprintf(p.out, "%c\n", p.rowSeparatorEnd(row))
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

func (p *printer) rowSeparatorEnd(row []Cell) rune {
	cell := row[len(row)-1]
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

func (p *printer) printRow(row []Cell) {
	for _, c := range row {
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

	lastCell := row[len(row)-1]
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
	for _, c = range lastRow {
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
