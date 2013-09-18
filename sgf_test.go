package libaduk

import (
    "testing"
    "io/ioutil"
)

const (
    Testgame9x9 = "testing/Batora-okao.sgf"
    TestgameEasy = "testing/Easy.sgf"
    TestgameSmall = "testing/Small.sgf"
    TestgameSmallMalformed = "testing/SmallMalformed.sgf"
)

/**
 * Small.sgf has this structure:
 *  
 *  X -> 1 -> 2 -> 3 -> 4 -> 5 -> 6
 *             \-> 3 -> 4
 *             |    \-> 4
 *             \-> 3 -> 4
 *
 * Checks if 2 has 3 children and 3.1 has 2 children
 */
func TestSgfReadAndNumChildren(t *testing.T) {
    sgfData, _ := ioutil.ReadFile(TestgameSmall)
    cursor, _ := NewCursor(sgfData)
    root, _ := cursor.GetRootNode(0)

    if root.numChildren != 1 {
        t.Errorf("Node 1 should have 1 children but was: %+v", root)
    }

    if root.Next.numChildren != 3 {
        t.Errorf("Node 2 should have 3 Children but was: %+v", root.Next)
    }

    if root.Next.Next.Down.numChildren != 2 {
        t.Errorf("Node 3.1 should have 2 children but was: %+v", root.Next.Next.Down)
    }
}

// Tests if sgf is correctly found malformed (it has a not closing SEQUENCE_START character)
func TestSgfReadMalformed(t *testing.T) {
    sgfData, _ := ioutil.ReadFile(TestgameSmallMalformed)
    _, err := NewCursor(sgfData)

    if err == nil {
        t.Errorf("Sgf should be malformed but was accepted!")
    }
}

// Tests reading easy sgf without errors
func TestEasySgfRead(t *testing.T) {
    sgfData, _ := ioutil.ReadFile(TestgameEasy)
    _, err := NewCursor(sgfData)

    if err != nil {
        t.Errorf("Reading an easy wellformed sgf should be successful but was %+v", err)
    }
}

// Tests if we can get the correct root games
func TestMultiRootGameAndGetRootNode(t *testing.T) {
    sgfData, _ := ioutil.ReadFile(Testgame9x9)
    cursor, _ := NewCursor(sgfData)

    _ = cursor.Game(0)
    if cursor.Current() == nil {
        t.Errorf("Cursor should have a root Node!")
    }

    _ = cursor.Game(1)
    if cursor.Current() == nil {
        t.Errorf("Cursor should have at least a second root Node!")
    }

    err := cursor.Game(2)
    if err == nil {
        t.Errorf("There should be no third root Node!")
    }
}
