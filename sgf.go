package libaduk

import (
    "fmt"
)

const (
    SEQUENCE_START = '('
    SEQUENCE_END = ')'
    NODE_START = ';'
    PROPERTY_START = '['
    PROPERTY_END = ']'
)

type Node struct {
    Previous *Node // Parent Node
    Next *Node // Child Node
    Up *Node // Upper sibling Node
    Down *Node // Lower sibling Node
    sgfData string
    numChildren int
    level int
}

func NewNode(prev *Node) *Node {
    return &Node { prev, nil, nil, nil, "", 0, 0 }
}

// The cursor parses Sgf data and provides methods to traverse the game
type Cursor struct {
    tree *Node
}

func NewCursor(sgf []byte) (*Cursor, error) {
    tree, err := parse(string(sgf))

    if err != nil {
        return nil, err
    }

    return &Cursor { tree }, nil
}

// Begin parse an sgf string
func parse(sgf string) (*Node, error) {
    tree := NewNode(nil)
    currentNode := tree
    sequenceNodes := make([]*Node, 0)
    nodeStartIndex := -1
    lastParsedType := SEQUENCE_END

    // range on string handles unicode automatically
    for i, value := range sgf {

        // If value is not a control character, ignore it
        if !(value == SEQUENCE_START || value == SEQUENCE_END || value == PROPERTY_START ||
                value == PROPERTY_END || value == NODE_START) {
            continue
        }

        // Sequence starts
        if value == SEQUENCE_START {
            // Safe sgf string to current node before creating a new one
            if lastParsedType != SEQUENCE_END && nodeStartIndex != -1 {
                currentNode.sgfData = sgf[nodeStartIndex:i]
            }

            // Create new Node for Sequence
            node := NewNode(currentNode)

            // Has current node already a child, than node is a sibling of currentNode
            if currentNode.Next != nil {
                last := currentNode.Next

                for last.Down != nil {
                    last = last.Down
                }

                node.Up = last
                last.Down = node
                node.level = last.level + 1
            } else {
                currentNode.Next = node
            }

            // Update current to new sequence
            currentNode.numChildren++

            // Add sequence to stack
            sequenceNodes = append(sequenceNodes, currentNode)

            currentNode = node
            nodeStartIndex = -1
            lastParsedType = SEQUENCE_START
        }

        // Sequence ends
        if value == SEQUENCE_END {
            // Safe sgf string to current node before creating a new one
            if lastParsedType != SEQUENCE_END && nodeStartIndex != -1 {
                currentNode.sgfData = sgf[nodeStartIndex:i]
            }

            // If we had sequences in the stack, set current node to last in stack
            if len(sequenceNodes) > 0 {
                currentNode = sequenceNodes[len(sequenceNodes) - 1]
                sequenceNodes = sequenceNodes[:len(sequenceNodes) - 1]
            }

            lastParsedType = SEQUENCE_END
        }

        // Node starts
        if value == NODE_START {
            if nodeStartIndex != -1 {
                // Safe sgf string to current node before creating a new one
                currentNode.sgfData = sgf[nodeStartIndex:i]

                // Create new node and update current
                node := NewNode(currentNode)
                currentNode.numChildren = 1
                currentNode.Next = node
                currentNode = node

            }

            nodeStartIndex = i
        }

        fmt.Printf("") //"%d:%+v ", i, string(value))
    }

    return tree, nil
}
