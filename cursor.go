package libaduk

import (
	"fmt"
)

// The cursor provides methods to traverse the game
type Cursor struct {
	rootNode *Node
	// Pointer to current Node in tree
	currentNode *Node
}

// Create a new cursor struct for given sgf data
func NewCursor(sgf []byte) (*Cursor, error) {
	tree, err := parse(string(sgf))

	if err != nil {
		return nil, err
	}

	return &Cursor{tree, tree}, nil
}

// Returns the n'th root node. In a normal game there is only one root (0)
func (cursor *Cursor) getRootNode(n int) (*Node, error) {
	if n >= cursor.rootNode.numChildren {
		return nil, fmt.Errorf("Cant find %d'th Root Node!", n)
	}

	node := cursor.rootNode.Next

	for i := 0; i < n; i++ {
		node = node.Down
	}

	return node, nil
}

// Set the Cursor to the n'th game
func (cursor *Cursor) Game(n int) (*Node, error) {
	gameNode, err := cursor.getRootNode(n)

	if err != nil {
		return nil, err
	}

	cursor.currentNode = gameNode

	return cursor.currentNode, nil
}

// Returns the Cursors current node
func (cursor *Cursor) Current() *Node {
	return cursor.currentNode
}

// Set the cursor to the n'th next node
func (cursor *Cursor) Next(n int) (*Node, error) {
	if n >= cursor.currentNode.numChildren {
		return nil, fmt.Errorf("Can't find %d'th Next Node!", n)
	}

	cursor.currentNode = cursor.currentNode.Next

	for i := 0; i < n; i++ {
		cursor.currentNode = cursor.currentNode.Down
	}

	return cursor.currentNode, nil
}

// Set the cursor to the previous node
func (cursor *Cursor) Previous() (*Node, error) {
	if cursor.currentNode.Previous == nil {
		return nil, fmt.Errorf("Can't find Previous Node!")
	}

	cursor.currentNode = cursor.currentNode.Previous

	return cursor.currentNode, nil
}

// Deletes the given variation from the tree
func (cursor *Cursor) DeleteVariation(node *Node) {
	if node.Previous != nil {
		cursor.removeNode(node)
	} else {
		if node.Next != nil {
			n := node.Next
			for n.Down != nil {
				n = node.Down
				cursor.removeNode(n.Up)
			}
			cursor.removeNode(n)
		}
		node.Next = nil
	}
}

// Remove the given node from the tree
func (cursor *Cursor) removeNode(node *Node) {
	// Update Node Up/Previous to not include given node
	if node.Up != nil {
		node.Up.Down = node.Down
	} else {
		node.Previous.Next = node.Down
	}

	if node.Down != nil {
		// Update Node Down to not include given node
		node.Down.Up = node.Up

		// Update levels for all Down nodes bei -1
		n := node.Down
		for n != nil {
			n.level--
			n = node.Down
		}
	}

	// Update children count of parent and destroy node
	node.Previous.numChildren--
	node = nil
}
