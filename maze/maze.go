package maze

//
// Cell is a flyweight that represents a single cell in the maze
//
type Cell struct {
	row int
	col int
	m   *Maze
}

// CellIterator - interface allowing iteration over sets
// of cells
type CellIterator interface {
	// Next - Get the next cell from the iterator, ok returns false when at the end of collection
	Next() (c Cell, ok bool)
	// Reset - Go back to the start of iteration
	Reset()
}

// Cell2DIterator is used to iterate over all the rows
// or columns one at a time.
type Cell2DIterator interface {
	Next() (iter *CellIterator, ok bool)
}

// Row - What row is this cell in?
func (c Cell) Row() int { return c.row }

// Col - what column is this cell in?
func (c Cell) Col() int { return c.col }

// Mark - Return the current mark in this cell
func (c Cell) Mark() int {
	return c.m.MarkAt(c.row, c.col)
}

// SetMark - set the mark for this cell
func (c Cell) SetMark(newMark int) {
	c.m.SetMarkAt(c.row, c.col, newMark)
}

//
// Go - Return a new Cell that represents a move
// from this cell one step in the given direction.
// DOES NOT PAY ATTENTION TO OPEN DOORS.
//
// If this cell is on the edge of the maze,
// an attempt to move out of the maze will
// return the current cell instead.
//
func (c Cell) Go(d Direction) Cell {
	if c.IsInMaze() {
		dRow, dCol := ToDyDx(d)
		return c.m.CellAt(c.row+dRow, c.col+dCol)
	}
	return c
}

//
// CanGo - Checks if there is an open door in the direction given.
//
func (c Cell) CanGo(d Direction) bool {
	return c.m.CanGo(c.row, c.col, d)
}

//
// OpenDoor - Open the door in the given direction.
//
func (c Cell) OpenDoor(d Direction) {
	c.m.OpenDoor(c.row, c.col, d)
}

//
// Iterator to retrieve the accessible neighbors
// from this cell - neighboring cells that have
// open doors to them.
//
type neighborIterator struct {
	c             Cell
	nextDirection Direction
}

func (iter *neighborIterator) Next() (c Cell, ok bool) {
	if iter.nextDirection == None {
		return iter.c, false
	}

	for d := iter.nextDirection; d != None; d = NextDirection(d) {
		if iter.c.CanGo(d) {
			n := iter.c.Go(d)
			iter.nextDirection = NextDirection(d)
			return n, true
		}
	}
	iter.nextDirection = None
	return iter.c, false
}

func (iter *neighborIterator) Reset() {
	iter.nextDirection = Up
}

// Neighbors - returns an iterator over the set of cells you
// can get to from this one through open doors
func (c Cell) Neighbors() CellIterator {
	return &neighborIterator{c, Up}
}

//
// IsInMaze - Check that this cell is within the bounds of the maze.
//
func (c Cell) IsInMaze() bool {
	return isInMaze(c.m, c.row, c.col)
}

//
// IsEntrance - Is this cell the maze entrance?
//
func (c Cell) IsEntrance() bool {
	return c.CanGo(Left) && c.col == 0
}

//
// IsExit - Is this cell the maze exit?
//
func (c Cell) IsExit() bool {
	return c.CanGo(Right) && c.col == c.m.cols-1
}

//
// Maze - this is basically an array of MazeCell structs, which
// track which doors are open from that cell, and what the current
// mark in the cell is.
//
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

// Rows - number of rows in the maze
func (m *Maze) Rows() int { return m.rows }

// Cols - number of columsn in the maze
func (m *Maze) Cols() int { return m.cols }

func isInMaze(m *Maze, row, col int) bool {
	return row >= 0 && row < m.rows && col >= 0 && col < m.cols
}

//
// CellAt - Construct a cell object for the given coordinates
//
func (m *Maze) CellAt(row, col int) Cell {
	return Cell{row, col, m}
}

//
// Various iterators to get combinations of cells.
// One for retrieving all cells in a row, one for
// retrieving all cells in a column, and one for
// getting all the cells in the maze (col major order)
//

type mazeRowIterator struct {
	m     *Maze
	row   int
	index int
}

func (iter *mazeRowIterator) Next() (c Cell, ok bool) {
	if iter.index >= iter.m.Cols() {
		return iter.m.CellAt(iter.row, iter.m.Cols()-1), false
	}

	c = iter.m.CellAt(iter.row, iter.index)
	iter.index++
	return c, true
}

func (iter *mazeRowIterator) Reset() {
	iter.index = 0
}

// Row - get the cells in the given row
func (m *Maze) Row(row int) CellIterator {
	return &mazeRowIterator{m, row, 0}
}

type mazeColIterator struct {
	m     *Maze
	col   int
	index int
}

func (iter *mazeColIterator) Next() (c Cell, ok bool) {
	if iter.index >= iter.m.Rows() {
		return iter.m.CellAt(iter.m.Rows()-1, iter.col), false
	}

	c = iter.m.CellAt(iter.index, iter.col)
	iter.index++
	return c, true
}

func (iter *mazeColIterator) Reset() {
	iter.index = 0
}

// Col - get the cells in the given column
func (m *Maze) Col(col int) CellIterator {
	return &mazeColIterator{m, col, 0}
}

type allCellsIterator struct {
	m   *Maze
	row int
	col int
}

func (iter *allCellsIterator) Next() (c Cell, ok bool) {
	if iter.row >= iter.m.Rows() {
		return Cell{-1, -1, iter.m}, false
	}

	c = iter.m.CellAt(iter.row, iter.col)

	iter.col++
	if iter.col >= iter.m.Cols() {
		iter.row++
		iter.col = 0
	}

	return c, true
}

func (iter *allCellsIterator) Reset() {
	iter.row = 0
	iter.col = 0
}

// AllCells - return all cells in the maze from top-left to bottom-right.
func (m *Maze) AllCells() CellIterator {
	return &allCellsIterator{m, 0, 0}
}

// Iterators that get either all rows or all columns

type allRowsIterator struct {
	m   *Maze
	row int
}

func (iter *allRowsIterator) Next() (ci *CellIterator, ok bool) {
	if iter.row >= iter.m.Rows() {
		return nil, false
	}
	row := iter.m.Row(iter.row)
	iter.row++
	return &row, true
}

// AllRows - Get all the rows in the maze.
func (m *Maze) AllRows() Cell2DIterator {
	return &allRowsIterator{m, 0}
}

type allColsIterator struct {
	m   *Maze
	col int
}

func (iter *allColsIterator) Next() (ci *CellIterator, ok bool) {
	if iter.col >= iter.m.Cols() {
		return nil, false
	}

	col := iter.m.Col(iter.col)
	iter.col++
	return &col, true
}

// AllCols - Get all the columns in the maze
func (m *Maze) AllCols() Cell2DIterator {
	return &allColsIterator{m, 0}
}

// Entrance - Get the entrance cell
func (m *Maze) Entrance() Cell {
	for r := 0; r < m.rows; r++ {
		if m.CanGo(r, 0, Left) {
			return m.CellAt(r, 0)
		}
	}
	return Cell{-1, -1, m}
}

// Exit - Get the exit cell
func (m *Maze) Exit() Cell {
	for r := 0; r < m.rows; r++ {
		if m.CanGo(r, m.cols-1, Right) {
			return m.CellAt(r, m.cols-1)
		}
	}
	return Cell{-1, -1, m}
}

// MarkAt - get the mark in the cell at given row, col
func (m *Maze) MarkAt(row, col int) int {
	if !isInMaze(m, row, col) {
		return 0
	}

	return m.cells[row][col].mark
}

// SetMarkAt - set the mark in the cell at given row, col
func (m *Maze) SetMarkAt(row, col int, mark int) {
	if isInMaze(m, row, col) {
		m.cells[row][col].mark = mark
	}
}

// SetAllMarks - mark every cell in the maze with the given value.
// Useful to set to common value for initialization.
func (m *Maze) SetAllMarks(mark int) {
	iter := m.AllCells()
	for c, ok := iter.Next(); ok; c, ok = iter.Next() {
		c.SetMark(mark)
	}
}

// CanGo - Can you move from row, col in the given direction
// Checks if destination is in the maze and there is an open door.
func (m *Maze) CanGo(row, col int, d Direction) bool {
	if !isInMaze(m, row, col) {
		return false
	}

	doorOpen := m.cells[row][col].doors & toDoorMask(d)
	return doorOpen != 0
}

// OpenDoor - open a door in the cell row, col in the given direction.
func (m *Maze) OpenDoor(row, col int, d Direction) {
	if isInMaze(m, row, col) {
		m.cells[row][col].doors |= toDoorMask(d)

		dRow, dCol := ToDyDx(d)
		if isInMaze(m, row+dRow, col+dCol) {
			m.cells[row+dRow][col+dCol].doors |= toDoorMask(d.toOppositeDirection())
		}
	}
}
