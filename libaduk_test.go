package libaduk

import (
    "testing"
)

// Tests invertColor functionality
func TestInvertColor(t *testing.T) {
    if BLACK.invert() != WHITE {
        t.Errorf("Inverted Color of %d should be %d but was %d!", BLACK, WHITE, BLACK.invert())
    }

    if WHITE.invert() != BLACK {
        t.Errorf("Inverted Color of %d should be %d but was %d!", WHITE, BLACK, WHITE.invert())
    }
}
