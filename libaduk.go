package libaduk

type BoardStatus uint8

// Inverts Black to White or White to Black
func (bs BoardStatus) invert() BoardStatus {
	if bs == WHITE {
		return BLACK
	}

	if bs == BLACK {
		return WHITE
	}

	return EMPTY
}

const (
	EMPTY BoardStatus = iota
	BLACK
	WHITE
	PASS
)

type Position struct {
	X uint8
	Y uint8
}

// Checks if the position has the same coordinates as b
func (a *Position) isSamePosition(b Position) bool {
	return a.X == b.X && a.Y == b.Y
}

// Represents a Move on the board
type Move struct {
	X        uint8
	Y        uint8
	Color    BoardStatus
	Captures []Position
}
