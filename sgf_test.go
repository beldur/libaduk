package libaduk

import (
    "testing"
    "io/ioutil"
    "fmt"
)

const (
    Testgame9x9 = "testing/Batora-okao.sgf"
    TestgameMultipleNodes = "testing/Multiple-Nodes.sgf"
    TestgameSmall = "testing/Small.sgf"
)

func TestSgfRead(t *testing.T) {
    sgfData, _ := ioutil.ReadFile(TestgameSmall)
    fmt.Println(string(sgfData))

    cursor, _ := NewCursor(sgfData)

    tree := cursor.tree

    for tree != nil {
        fmt.Printf("%p, %+v\n", tree, tree)

        down := tree.Down
        for down != nil {
            fmt.Printf("\t%p, %+v\n", down, down)
            down = down.Down
        }

        tree = tree.Next
    }
}
