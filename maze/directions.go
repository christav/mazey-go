package maze

// Direction to move in the maze
type Direction int

const (
	// None - Direction that is no move
	None Direction = iota
	// Up - Move up
	Up
	// Down - Move down
	Down
	// Left - Move left
	Left
	// Right - Move right
	Right
)

var allDirs = []Direction{Up, Down, Left, Right}

// AllDirections - Return slice containing all directions to move in the maze except None
func AllDirections() []Direction {
	return allDirs[:]
}

// NextDirection - Given a direction, return the next direction to process
func NextDirection(d Direction) Direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Left
	case Left:
		return Right
	default:
		return None
	}
}

// ToDyDx - Convert a direction to a row, col offset
func ToDyDx(d Direction) (int, int) {
	switch d {
	case None:
		return 0, 0
	case Up:
		return -1, 0
	case Down:
		return +1, 0
	case Left:
		return 0, -1
	case Right:
		return 0, +1
	default:
		return 0, 0
	}
}

func toDoorMask(d Direction) int {
	switch d {
	case Up:
		return 1
	case Down:
		return 2
	case Left:
		return 4
	case Right:
		return 8
	default:
		return 255
	}
}

func (d Direction) toOppositeDirection() Direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	default:
		return None
	}
}
