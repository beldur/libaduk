package libaduk

type BoardStatus uint8

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

// Represents a Move on the board
type Move struct {
    X uint8
    Y uint8
    Color BoardStatus
    Captures []Position
}
