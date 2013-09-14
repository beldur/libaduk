package libaduk

import (
    "testing"
)

// Tests clearing the board
func TestClearBoard(t *testing.T) {
    board, _ := NewBoard(9)

    board.setStatus(4, 6, BLACK)
    board.Clear()

    if board.getStatus(4, 6) != EMPTY {
        t.Errorf("Position 4,6 should be %d but was %d!", EMPTY, board.getStatus(4, 6))
    }
}

// Tests internal set/get status
func TestGetStatusAndSetStatus(t *testing.T) {
    board, _ := NewBoard(9)

    board.setStatus(4, 6, BLACK)
    board.setStatus(3, 5, WHITE)

    if board.getStatus(4, 6) != BLACK {
        t.Errorf("Position 4,6 should be %d but was %d!", BLACK, board.getStatus(4, 6))
    }

    if board.getStatus(3, 5) != WHITE {
        t.Errorf("Position 3,5 should be %d but was %d!", WHITE, board.getStatus(3, 5))
    }
}

// Tests invertColor functionality
func TestInvertColor(t *testing.T) {
    board, _ := NewBoard(9)

    if board.invertColor(BLACK) != WHITE {
        t.Errorf("Inverted Color of %d should be %d but was %d!", BLACK, WHITE, board.invertColor(BLACK))
    }

    if board.invertColor(WHITE) != BLACK {
        t.Errorf("Inverted Color of %d should be %d but was %d!", WHITE, BLACK, board.invertColor(WHITE))
    }
}

// Tests getNeighbours functionality
func TestNeighbours(t *testing.T) {
    board, _ := NewBoard(9)

    board.setStatus(0, 0, BLACK)
    board.setStatus(1, 0, WHITE)
    neighbours := board.getNeighbours(0, 0)

    if len(neighbours) != 2 {
        t.Errorf("Position 0,0 should have 2 neighbours but had %d!", len(neighbours))
    }

    for i := 0; i < len(neighbours); i++ {
        neighbourStatus := board.getStatus(neighbours[i].X, neighbours[i].Y)

        if neighbours[i].X == 1 && neighbourStatus != WHITE {
            t.Errorf("Right neighbour should be %d but was %d!", WHITE, neighbourStatus)
        }

        if neighbours[i].X == 0 && neighbourStatus != EMPTY {
            t.Errorf("Bottom neighbour should be %d but was %d!", EMPTY, neighbourStatus)
        }
    }
}

// Tests to play on an occupied position or outside of the board
func TestPlayOccupiedPostionAndOffBoardsize(t *testing.T) {
    board, _ := NewBoard(9)

    err := board.Play(9, 9, BLACK)

    if err == nil {
        t.Errorf("A play on 9, 9 should be illegal!")
    }

    board.setStatus(3, 4, BLACK)
    err = board.Play(3, 4, BLACK)

    if err == nil {
        t.Errorf("A play on an occupied position should be illegal!")
    }
}

// Tests if the played move is recognized as a suicide
func TestPlayIsSuicide(t *testing.T) {
    board, _ := NewBoard(9)

    board.Play(0, 1, BLACK)
    board.Play(0, 3, BLACK)
    board.Play(1, 2, BLACK)
    err := board.Play(0, 2, WHITE)

    if err == nil {
        t.Errorf("A suicide move should be illegal!")
    }
}

// Tests if the played move captures the stones with no liberties left
func TestPlayAndCaptureGroup(t *testing.T) {
    board, _ := NewBoard(9)

    board.Play(0, 2, BLACK)
    board.Play(0, 1, BLACK)
    board.Play(0, 0, WHITE)
    board.Play(1, 1, WHITE)
    board.Play(1, 2, WHITE)

    // This move should capture the two black stones
    board.Play(0, 3, WHITE)

    if board.getStatus(0, 2) != EMPTY || board.getStatus(0, 1) != EMPTY {
        t.Errorf("There should be no stones on the captures position!")
    }
}


// Tests if Undo turns board to the previous position
func TestPlayAndCaptureAndUndu(t *testing.T) {
    board, _ := NewBoard(9)

    board.Play(0, 2, BLACK)
    board.Play(0, 1, BLACK)
    board.Play(0, 0, WHITE)
    board.Play(1, 1, WHITE)
    board.Play(1, 2, WHITE)

    boardPositionBeforeUndo := board.ToString()

    // Play and undo move immediatly 
    board.Play(0, 3, WHITE)
    board.Undo(1)

    boardPositionAfterUndo := board.ToString()

    if boardPositionBeforeUndo != boardPositionAfterUndo {
        t.Errorf("Undo should have recovered old board position!")
    }
}
