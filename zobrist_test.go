package libaduk

import (
    "testing"
)

// Test basic Zobrist hashing
func TestBasicZobristHash(t *testing.T) {
    zob := NewZobristHash(2)

    hashOne, _ := zob.Hash(0, 0, BLACK)
    _, _ = zob.Hash(1, 1, WHITE)
    _, _ = zob.Hash(0, 1, BLACK)
    _, _ = zob.Hash(1, 1, WHITE)
    hashTwo, _ := zob.Hash(0, 1, BLACK)

    if hashOne != hashTwo {
        t.Errorf("The second hash (%d) should be equal to the first (%d)", hashOne, hashTwo)
    }
}
