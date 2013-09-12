package libaduk

import (
    "fmt"
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

    for x := uint8(0); x < board.BoardSize; x++ {
        for y := uint8(0); y < board.BoardSize; y++ {
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

func (board *AbstractBoard) getStatus(x uint8, y uint8) BoardStatus {
    return board.data[board.BoardSize * x + y]
}

func (board *AbstractBoard) SetStatus(x uint8, y uint8, status BoardStatus) {
    board.data[board.BoardSize * x + y] = status
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
