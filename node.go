package libaduk

import ()

type Node struct {
	Previous    *Node // Parent Node
	Next        *Node // Child Node
	Up          *Node // Upper sibling Node
	Down        *Node // Lower sibling Node
	sgfData     string
	numChildren int
	level       int
}

func NewNode(prev *Node) *Node {
	return &Node{prev, nil, nil, nil, "", 0, 0}
}

func (node *Node) ToString() string {
	return "TODO"
}
