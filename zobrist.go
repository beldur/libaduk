package libaduk

import (
	"fmt"
	"math/rand"
	"time"
)

type ZobristHash struct {
	table     [][]int64
	hash      int64
	boardsize uint8
}

// Create a new Zobrist struct for given boardsize
func NewZobristHash(boardSize uint8) *ZobristHash {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	table := make([][]int64, boardSize*boardSize)

	for i, _ := range table {
		table[i] = []int64{rnd.Int63(), rnd.Int63()}
	}

	return &ZobristHash{
		table,
		0,
		boardSize,
	}
}

// Update the hash for the played move
func (zob *ZobristHash) Hash(x uint8, y uint8, status BoardStatus) (int64, error) {
	var index int

	if status == WHITE {
		index = 1
	} else if status == BLACK {
		index = 0
	} else {
		return -1, fmt.Errorf("The provided status (%d) is not valid!", status)
	}

	zob.hash ^= zob.table[zob.boardsize*x+y][index]

	return zob.hash, nil
}

// Returns the current hash value
func (zob *ZobristHash) GetHash() int64 {
	return zob.hash
}
