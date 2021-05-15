package rbtree

import (
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

var (
	black colour = false
	red   colour = true
)

type colour bool

// NODE ---

// ver 1 only accepts string values
type Node struct {
	key    string
	value  string
	colour colour
	parent *Node
	left   *Node
	right  *Node
}

func (node *Node) Get(key string) (string, error) {
	if node.key == key {
		return node.value, nil
	}
	if node == nil {
		return "", vars.GET_FAIL_ERROR
	}

	if node.key < key {
		return node.left.Get(key)
	}
	return node.right.Get(key)
}

func nodeColour(node *Node) colour {
	if node == nil {
		return black
	}
	return node.colour
}

func nodeUncle(node *Node) *Node {
	if node == nil || node.parent == nil || node.parent.parent == nil {
		return nil
	}
	if node.parent == node.parent.parent.left {
		return node.parent.parent.right
	}
	return node.parent.parent.left
}

func nodeGrandparent(node *Node) *Node {
	if node != nil && node.parent != nil {
		return node.parent.parent
	}
	return nil
}

// NODE ---

// RBTREE ---

type RedBlackTree struct {
	root *Node
}

func NewTree() *RedBlackTree {
	return &RedBlackTree{
		root: nil,
	}
}

func (rbtree *RedBlackTree) Get(key string) (string, error) {
	node := rbtree.root
	return node.Get(key)
}

func (rbtree *RedBlackTree) Put(key string, value string) error {
	defer func() error {
		var ret error = nil
		if err := recover(); err != nil {
			ret = vars.PUT_FAIL_ERROR
		}
		return ret
	}()

	node := &Node{
		key:    key,
		value:  value,
		colour: red,
		parent: nil,
		left:   nil,
		right:  nil,
	}

	x := rbtree.root
	var y *Node = nil
	for x != nil {
		y = x
		if key < x.key {
			x = x.left
		} else {
			x = x.right
		}
	}

	node.parent = y
	if y == nil {
		rbtree.root = node
	} else if key < y.key {
		y.left = node
	} else if key > y.key {
		y.right = node
	} else { // if key == y.key, overwrite the value. + Tombstone will be a special string
		y.value = value
		return nil
	}
	if node.parent == nil {
		node.colour = black
		return nil
	}
	if node.parent.parent == nil {
		return nil
	}
	rbtree.insertFix1(node)
	return nil
}

func (rbtree *RedBlackTree) insertFix1(node *Node) {
	if node.parent == nil {
		node.colour = black
	} else {
		rbtree.insertFix2(node)
	}
}

func (rbtree *RedBlackTree) insertFix2(node *Node) {
	if nodeColour(node.parent) == black {
		return
	}
	rbtree.insertFix3(node)
}

func (rbtree *RedBlackTree) insertFix3(node *Node) {
	uncle := nodeUncle(node)
	if nodeColour(uncle) == red {
		node.parent.colour = black
		uncle.colour = black
		nodeGrandparent(node).colour = red
		rbtree.insertFix1(nodeGrandparent(node))
	} else {
		rbtree.insertFix4(node)
	}
}

func (rbtree *RedBlackTree) insertFix4(node *Node) {
	grandparent := nodeGrandparent(node)
	if node == node.parent.right && node.parent == grandparent.left {
		rbtree.leftRotate(node.parent)
		node = node.left
	} else if node == node.parent.left && node.parent == grandparent.right {
		rbtree.rightRotate(node.parent)
		node = node.right
	}
	rbtree.insertFix5(node)
}

func (rbtree *RedBlackTree) insertFix5(node *Node) {
	node.parent.colour = black
	grandparent := nodeGrandparent(node)
	grandparent.colour = red
	if node == node.parent.left && node.parent == grandparent.left {
		rbtree.rightRotate(grandparent)
	} else if node == node.parent.right && node.parent == grandparent.right {
		rbtree.leftRotate(grandparent)
	}
}

func (rbtree *RedBlackTree) leftRotate(node *Node) {
	y := node.right
	node.right = y.left
	if y.left != nil {
		y.left.parent = node
	}
	y.parent = node.parent
	if node.parent == nil {
		rbtree.root = y
	} else if node == node.parent.left {
		node.parent.left = y
	} else {
		node.parent.right = y
	}
	y.left = node
	node.parent = y
}

func (rbtree *RedBlackTree) rightRotate(node *Node) {
	y := node.left
	node.left = y.right
	if y.right != nil {
		y.right.parent = node
	}
	y.parent = node.parent
	if node.parent == nil {
		rbtree.root = y
	} else if node == node.parent.right {
		node.parent.right = y
	} else {
		node.parent.left = y
	}
	y.right = node
	node.parent = y
}

func (rbtree *RedBlackTree) Flush() []vars.KeyValue {
	ret := []vars.KeyValue{}
	flushHelper(rbtree.root, &ret)
	return ret
}

func flushHelper(node *Node, out *[]vars.KeyValue) {
	if node == nil {
		return
	}
	flushHelper(node.left, out)
	*out = append(*out, vars.KeyValue{Key: node.key, Value: node.value})
	flushHelper(node.right, out)
}

// RBTREE ---
