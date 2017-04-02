package libaduk

import (
	"fmt"
	"log"
)

// Represents a Go board data structure
type AbstractBoard struct {
	BoardSize uint8
	data      []BoardStatus
	undoStack []*Move
	zobrist   *ZobristHash
}

// Creates new Go Board
func NewBoard(boardSize uint8) (*AbstractBoard, error) {
	if boardSize < 1 {
		return nil, fmt.Errorf("Boardsize can not be less than 1!")
	}

	return &AbstractBoard{
		boardSize,
		make([]BoardStatus, boardSize*boardSize),
		make([]*Move, 0),
		NewZobristHash(boardSize),
	}, nil
}

func (board *AbstractBoard) Len() int {
	return len(board.data)
}

//Verify if a given position is within the boundaries of the board
func (board *AbstractBoard) Contains(position Position) bool {
	return position.X >= 0 &&
		position.X < board.BoardSize &&
		position.Y >= 0 &&
		position.Y < board.BoardSize
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
	board.zobrist.hash = 0
	board.undoStack = []*Move{}
}

// Returns the Top Move of the Undostack
func (board *AbstractBoard) UndostackTopMove() *Move {
	return board.undoStack[len(board.undoStack)-1]
}

// Removes last Move from Undostack
func (board *AbstractBoard) UndostackPop() (move *Move) {
	if len(board.undoStack) > 0 {
		move = board.undoStack[len(board.undoStack)-1]
		board.undoStack = board.undoStack[:len(board.undoStack)-1]
	}

	return
}

// Adds the given Move to the Undostack
func (board *AbstractBoard) UndostackPush(move *Move) {
	board.undoStack = append(board.undoStack, move)
}

// Adds a Pass to the Undostack
func (board *AbstractBoard) UndostackPushPass() {
	board.UndostackPush(&Move{255, 255, PASS, nil})
}

// Undo `count` moves on the board
func (board *AbstractBoard) Undo(count int) {
	for i := 0; i < count; i++ {
		if len(board.undoStack) > 0 {
			move := board.UndostackPop()

			// Remove stone from the board and update hash
			if move.Color == BLACK || move.Color == WHITE {
				board.zobrist.Hash(move.X, move.Y, move.Color)
				board.setStatus(move.X, move.Y, EMPTY)
			}

			// Add captures back to board if necessary and update hash
			for _, capture := range move.Captures {
				board.zobrist.Hash(capture.X, capture.Y, move.Color.invert())
				board.setStatus(capture.X, capture.Y, move.Color.invert())
			}
		}
	}
}

// Returns current board hash value
func (board *AbstractBoard) GetHash() int64 {
	return board.zobrist.GetHash()
}

// Play move on board
func (board *AbstractBoard) PlayMove(move Move) error {
	return board.Play(move.X, move.Y, move.Color)
}

// Play stone at given position
func (board *AbstractBoard) Play(x uint8, y uint8, color BoardStatus) error {
	log.Printf("Play: X: %v, Y: %v, Color: %v", x, y, color)

	// Is move on the board?
	if x < 0 || x >= board.BoardSize || y < 0 || y >= board.BoardSize {
		return fmt.Errorf("Invalid move position!")
	}

	// Is already a stone on this position?
	if board.getStatus(x, y) != EMPTY {
		return fmt.Errorf("Position already occupied!")
	}

	// Check if move is legal and get captures
	captures, err := board.legal(x, y, color)
	if err != nil {
		return err
	}

	// Remove captures
	for _, capture := range captures {
		board.zobrist.Hash(capture.X, capture.Y, color.invert())
		board.setStatus(capture.X, capture.Y, EMPTY)
	}

	// Add them to undostack
	board.UndostackPush(&Move{x, y, color, captures})

	return nil
}

// Checks if move is legal and returns captured stones if necessary
func (board *AbstractBoard) legal(x uint8, y uint8, color BoardStatus) (captures []Position, err error) {
	captures = []Position{}
	neighbours := board.getNeighbours(x, y)

	// Check if we capture neighbouring stones
	for _, neighbour := range neighbours {
		// Is neighbour from another color?
		if board.getStatus(neighbour.X, neighbour.Y) == color.invert() {
			// Get enemy stones with no liberties left
			noLibertyStones := board.getNoLibertyStones(neighbour.X, neighbour.Y, Position{x, y})
			for _, noLibertyStone := range noLibertyStones {
				captures = append(captures, noLibertyStone)
			}
		}
	}

	// Place stone on the board and update hash
	board.zobrist.Hash(x, y, color)
	board.setStatus(x, y, color)

	// TODO: Delete Duplicates necessary????
	if len(captures) > 0 {
		return
	}

	// Check if the played move has no liberties and therefore is a suicide
	selfNoLiberties := board.getNoLibertyStones(x, y, Position{})

	if len(selfNoLiberties) > 0 {
		// Take move back
		board.zobrist.Hash(x, y, color)
		board.setStatus(x, y, EMPTY)
		err = fmt.Errorf("Invalid move (Suicide not allowed)!")
	}

	log.SetPrefix("")
	return
}

// Get all stones with no liberties left on given position
func (board *AbstractBoard) getNoLibertyStones(x uint8, y uint8, orgPosition Position) (noLibertyStones []Position) {
	log.Printf("Get no liberty stones for (%d, %d)", x, y)

	noLibertyStones = []Position{}
	newlyFoundStones := []Position{Position{x, y}}
	foundNew := true
	var groupStones []Position = nil

	// Search until no new stones are found
	for foundNew == true {
		foundNew = false
		groupStones = []Position{}

		for _, newlyFoundStone := range newlyFoundStones {
			neighbours := board.getNeighbours(newlyFoundStone.X, newlyFoundStone.Y)

			// Check liberties of stone newlyFoundStone.X, newlyFoundStone.Y by checking the neighbours
			for _, neighbour := range neighbours {
				nbX := neighbour.X
				nbY := neighbour.Y

				// Has newlyFoundStone a free liberty?
				if board.getStatus(nbX, nbY) == EMPTY && !neighbour.isSamePosition(orgPosition) {
					// Neighbour is empty and not origPosition so newlyFoundStone has at least one liberty
					return noLibertyStones[:0]
				} else {
					// Is the neighbour of newlyFoundStone.X, newlyFoundStone.Y the same color? Then we have a group here
					if board.getStatus(newlyFoundStone.X, newlyFoundStone.Y) == board.getStatus(nbX, nbY) {
						foundNewHere := true
						nbGroupStone := Position{nbX, nbY}

						log.Printf("Found group stone for (%d, %d) at %+v", newlyFoundStone.X, newlyFoundStone.Y, nbGroupStone)

						// Check if found stone is already in our group list
						for _, groupStone := range groupStones {
							if groupStone.isSamePosition(nbGroupStone) {
								foundNewHere = false
								break
							}
						}

						// Check if found stone is already in result set list
						if foundNewHere {
							for _, noLibertyStone := range noLibertyStones {
								if noLibertyStone.isSamePosition(nbGroupStone) {
									foundNewHere = false
									break
								}
							}
						}

						// If groupStone is not known yet, add it
						if foundNewHere {
							groupStones = append(groupStones, nbGroupStone)
							foundNew = true
						}
					}
				}
			}
		}

		// Add newly found stones to the resultset
		noLibertyStones = append(noLibertyStones, newlyFoundStones...)

		// Now check the found group stones
		newlyFoundStones = groupStones
	}

	log.Printf("Found these stones with no liberties: %+v", noLibertyStones)

	return
}

// Returns the neighbour array positions for a given point
func (board *AbstractBoard) getNeighbours(x uint8, y uint8) (neighbourIndexes []Position) {
	neighbourIndexes = []Position{}
	possibleNeighbours := [...]Position{
		Position{(x - 1), y},
		Position{(x + 1), y},
		Position{x, (y - 1)},
		Position{x, (y + 1)},
	}

	for _, position := range possibleNeighbours {
		if board.Contains(position) {
			neighbourIndexes = append(neighbourIndexes, position)
		}
	}
	return
}

func (board *AbstractBoard) getStatus(x uint8, y uint8) BoardStatus {
	return board.data[board.BoardSize*x+y]
}

func (board *AbstractBoard) setStatus(x uint8, y uint8, status BoardStatus) {
	board.data[board.BoardSize*x+y] = status
}
