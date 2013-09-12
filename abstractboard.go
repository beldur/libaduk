package libaduk

import (
    "fmt"
    "log"
)

// Represents a Go board data structure
type AbstractBoard struct {
    BoardSize uint8
    data []BoardStatus
    undoStack []*Move
}

// Creates new Go Board
func NewBoard(boardSize uint8) (*AbstractBoard, error) {
    if boardSize < 1 {
        return nil, fmt.Errorf("Boardsize can not be less than 1!")
    }

    return &AbstractBoard {
        boardSize,
        make([]BoardStatus, boardSize * boardSize),
        make([]*Move, 0),
    }, nil
}

// Returns a string representation of the current board status
func (board *AbstractBoard) ToString() string {
    result := ""

    for y := uint8(0); y < board.BoardSize; y++ {
        for x := uint8(0); x < board.BoardSize; x++ {
            switch board.getStatus(x, y) {
            case EMPTY:
                result += ". "
            case BLACK:
                result += "X "
            case WHITE:
                result += "O "
            }
        }
        result += "\n"
    }

    return result
}

// Clears the board
func (board *AbstractBoard) Clear() {
    for i := 0; i < len(board.data); i++ {
        board.data[i] = EMPTY
    }
    board.undoStack = make([]*Move, 0)
}

// Returns the Top Position of the Undostack
func (board *AbstractBoard) UndostackTopPosition() Position {
    move := board.undoStack[len(board.undoStack) - 1]
    return Position{ move.X, move.Y }
}

// Returns the Top Color of the Undostack
func (board *AbstractBoard) UndostackTopColor() BoardStatus {
    move := board.undoStack[len(board.undoStack) - 1]
    return move.Color
}

// Returns the Top Captures of the Undostack
func (board *AbstractBoard) UndostackTopCaptures() []Position {
    move := board.undoStack[len(board.undoStack) - 1]
    return move.Captures
}

// Removes last Move from Undostack
func (board *AbstractBoard) UndostackPop() (move *Move) {
    if len(board.undoStack) > 0 {
        move = board.undoStack[len(board.undoStack) - 1]
        board.undoStack = board.undoStack[:len(board.undoStack) - 1]
    }

    return
}

// Adds the given Move to the Undostack
func (board *AbstractBoard) UndostackPush(move *Move) {
    board.undoStack = append(board.undoStack, move)
}

// Adds a Pass to the Undostack
func (board *AbstractBoard) UndostackAppendPass() {
    board.undoStack = append(board.undoStack, &Move { 255, 255, PASS, nil })
}

// Play move on board
func (board *AbstractBoard) PlayMove(move Move) (bool, error) {
    return board.Play(move.X, move.Y, move.Color)
}

// Play stone at given position
func (board *AbstractBoard) Play(x uint8, y uint8, color BoardStatus) (bool, error) {
    log.Printf("Play: X: %v, Y: %v, Color: %v", x, y, color)

    // Is move on the board?
    if x < 0 || x >= board.BoardSize || y < 0 || y >= board.BoardSize {
        return false, fmt.Errorf("Invalid move position!")
    }

    // Is already a stone on this position?
    if board.getStatus(x, y) != EMPTY {
        return false, fmt.Errorf("Position already occupied!")
    }

    captures, err := board.legal(x, y, color)
    if err != nil {
        return false, fmt.Errorf("Move is not legal!")
    }

    log.Printf("Captures: %+v", captures)

    // TODO: Remove captures and add them to undostack

    return true, nil
}

// Checks if move is legal and returns captured stones if necessary
func (board *AbstractBoard) legal(x uint8, y uint8, color BoardStatus) (captures []Position, err error) {
    captures = make([]Position, 0)
    neighbours := board.neighbours(x, y)

    // Check if we capture neighbouring stones
    for i := 0; i < len(neighbours); i++ {
        // Is neighbour from another color?
        if board.getStatus(neighbours[i].X, neighbours[i].Y) == board.invertColor(color) {
            log.Printf("Neighbour of (X: %d, Y: %d) at (X: %d, Y: %d) is %v",
                x, y, neighbours[i].X, neighbours[i].Y, board.invertColor(color))

            // Get enemy stones with no liberties left
            noLibertyStones := board.getNoLibertyStones(neighbours[i].X, neighbours[i].Y, int(board.BoardSize * x + y))
            for j := 0; j < len(noLibertyStones); j++ {
                captures = append(captures, noLibertyStones[j])
            }
        }
    }

    board.setStatus(x, y, color)

    // TODO: Delete Duplicates necessary????
    if len(captures) > 0 {
        return
    }

    // TODO: Check for suicide
    selfNoLiberties := board.getNoLibertyStones(x, y, 0)
    if len(selfNoLiberties) > 0 {
        // Take move back
        board.setStatus(x, y, EMPTY)
        err = fmt.Errorf("Invalid move (Suicide not allowed)!")
    }

    return
}

// Get all stones with no liberties left on given position
func (board *AbstractBoard) getNoLibertyStones(x uint8, y uint8, exc int) (noLibertyStones []Position) {
    noLibertyStones = make([]Position, 0)

    // TODO: Implement
    return
}

// Returns the neighbour array positions for a given point
func (board *AbstractBoard) neighbours(x uint8, y uint8) (neighbourIndexes []Position) {
    neighbourIndexes = make([]Position, 0)

    // Check for board borders
    if x >= 1 {
        neighbourIndexes = append(neighbourIndexes, Position { (x - 1), y })
    }
    if x < board.BoardSize - 1 {
        neighbourIndexes = append(neighbourIndexes, Position { (x + 1), y })
    }
    if y >= 1 {
        neighbourIndexes = append(neighbourIndexes, Position { x, y - 1 })
    }
    if y < board.BoardSize - 1 {
        neighbourIndexes = append(neighbourIndexes, Position { x, y + 1 })
    }

    log.Printf("Neighbours for (X: %d, Y: %d) are %+v", x, y, neighbourIndexes)

    return
}

func (board *AbstractBoard) getStatus(x uint8, y uint8) BoardStatus {
    return board.data[board.BoardSize * x + y]
}

func (board *AbstractBoard) setStatus(x uint8, y uint8, status BoardStatus) {
    board.data[board.BoardSize * x + y] = status
}

// Inverts Black to White or White to Black
func (board *AbstractBoard) invertColor(color BoardStatus) BoardStatus {
    if color == WHITE {
        return BLACK
    }

    if color == BLACK {
        return WHITE
    }

    return EMPTY
}
