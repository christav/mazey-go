package maze

type Cell struct {
	row int
	col int
	m   *Maze
}

func (c Cell) Row() int { return c.row }
func (c Cell) Col() int { return c.col }
func (c Cell) Mark() int {
	return c.m.MarkAt(c.row, c.col)
}
func (c Cell) SetMark(newMark int) {
	c.m.SetMarkAt(c.row, c.col, newMark)
}

func (c Cell) Go(d Direction) Cell {
	if c.IsInMaze() {
		dRow, dCol := ToDyDx(d)
		return c.m.CellAt(c.row+dRow, c.col+dCol)
	}
	return c
}

func (c Cell) CanGo(d Direction) bool {
	return c.m.CanGo(c.row, c.col, d)
}

func (c Cell) OpenDoor(d Direction) {
	c.m.OpenDoor(c.row, c.col, d)
}

func (c Cell) Neighbors() <-chan Cell {
	ch := make(chan Cell)
	go func() {
		for d := range AllDirections() {
			if c.CanGo(d) {
				ch <- c.Go(d)
			}
		}
		close(ch)
	}()
	return ch
}

func (c Cell) IsInMaze() bool {
	return isInMaze(c.m, c.row, c.col)
}

func (c Cell) IsEntrance() bool {
	return c.CanGo(Left) && c.col == 0
}

func (c Cell) IsExit() bool {
	return c.CanGo(Right) && c.col == c.m.cols-1
}

type Maze struct {
	rows  int
	cols  int
	cells [][]mazeCell
}

type mazeCell struct {
	mark  int
	doors int
}

// NewMaze create a new maze of given height and width
func NewMaze(rows, cols int) *Maze {
	cellRows := make([][]mazeCell, rows)
	for row := 0; row < rows; row++ {
		cellRows[row] = make([]mazeCell, cols)
	}
	return &Maze{rows, cols, cellRows}
}

func (m *Maze) Rows() int { return m.rows }
func (m *Maze) Cols() int { return m.cols }

func isInMaze(m *Maze, row, col int) bool {
	return row >= 0 && row < m.rows && col >= 0 && col < m.cols
}

func (m *Maze) CellAt(row, col int) Cell {
	return Cell{row, col, m}
}

func (m *Maze) Row(row int) <-chan Cell {
	ch := make(chan Cell)
	go func() {
		if row >= 0 && row < m.rows {
			for c := 0; c < m.cols; c++ {
				ch <- Cell{row, c, m}
			}
		}
		close(ch)
	}()
	return ch
}

func (m *Maze) AllRows() <-chan []Cell {
	ch := make(chan []Cell)
	go func() {
		defer close(ch)
		for r := 0; r < m.rows; r++ {
			ch <- toSlice(m.Row(r))
		}
	}()
	return ch
}

func (m *Maze) Col(col int) <-chan Cell {
	ch := make(chan Cell)
	go func() {
		defer close(ch)
		if col >= 0 && col < m.cols {
			for r := 0; r < m.rows; r++ {
				ch <- Cell{r, col, m}
			}
		}
	}()
	return ch
}

func (m *Maze) AllCols() <-chan []Cell {
	ch := make(chan []Cell)
	go func() {
		defer close(ch)
		for c := 0; c < m.cols; c++ {
			ch <- toSlice(m.Col(c))
		}
	}()
	return ch
}

func (m *Maze) AllCells() <-chan Cell {
	ch := make(chan Cell)
	go func() {
		for r := 0; r < m.rows; r++ {
			for c := 0; c < m.cols; c++ {
				ch <- Cell{r, c, m}
			}
		}
		close(ch)
	}()

	return ch
}

func toSlice(cells <-chan Cell) []Cell {
	s := make([]Cell, 0)
	for c := range cells {
		s = append(s, c)
	}
	return s
}

func (m *Maze) Entrance() Cell {
	for r := 0; r < m.rows; r++ {
		if m.CanGo(r, 0, Left) {
			return m.CellAt(r, 0)
		}
	}
	return Cell{-1, -1, m}
}

func (m *Maze) MarkAt(row, col int) int {
	if !isInMaze(m, row, col) {
		return 0
	}

	return m.cells[row][col].mark
}

func (m *Maze) SetMarkAt(row, col int, mark int) {
	if isInMaze(m, row, col) {
		m.cells[row][col].mark = mark
	}
}

func (m *Maze) CanGo(row, col int, d Direction) bool {
	if !isInMaze(m, row, col) {
		return false
	}

	doorOpen := m.cells[row][col].doors & toDoorMask(d)
	return doorOpen != 0
}

func (m *Maze) OpenDoor(row, col int, d Direction) {
	if isInMaze(m, row, col) {
		m.cells[row][col].doors |= toDoorMask(d)

		dRow, dCol := ToDyDx(d)
		if isInMaze(m, row+dRow, col+dCol) {
			m.cells[row+dRow][col+dCol].doors |= toDoorMask(d.toOppositeDirection())
		}
	}
}
