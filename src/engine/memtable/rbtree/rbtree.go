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

func (this *Node) Get(key string) (string, error) {
	if this == nil {
		return "", vars.GET_FAIL_ERROR
	}
	if this.key == key {
		return this.value, nil
	}

	if this.key < key {
		return this.right.Get(key)
	}
	return this.left.Get(key)
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

func (this *RedBlackTree) Get(key string) (string, error) {
	node := this.root
	return node.Get(key)
}

func (this *RedBlackTree) Put(key string, value string) error {
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

	x := this.root
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
		this.root = node
	} else if key < y.key {
		y.left = node
	} else if key > y.key {
		y.right = node
	} else { // if key == y.key, overwrite the value. + Tombstone will be a special string
		y.value = value
		// return nil
	}
	if node.parent == nil {
		node.colour = black
		return nil
	}
	if node.parent.parent == nil {
		return nil
	}
	this.insertFix1(node)
	return nil
}

func (this *RedBlackTree) insertFix1(node *Node) {
	if node.parent == nil {
		node.colour = black
	} else {
		this.insertFix2(node)
	}
}

func (this *RedBlackTree) insertFix2(node *Node) {
	if nodeColour(node.parent) == black {
		return
	}
	this.insertFix3(node)
}

func (this *RedBlackTree) insertFix3(node *Node) {
	uncle := nodeUncle(node)
	if nodeColour(uncle) == red {
		node.parent.colour = black
		uncle.colour = black
		nodeGrandparent(node).colour = red
		this.insertFix1(nodeGrandparent(node))
	} else {
		this.insertFix4(node)
	}
}

func (this *RedBlackTree) insertFix4(node *Node) {
	grandparent := nodeGrandparent(node)
	if node == node.parent.right && node.parent == grandparent.left {
		this.leftRotate(node.parent)
		node = node.left
	} else if node == node.parent.left && node.parent == grandparent.right {
		this.rightRotate(node.parent)
		node = node.right
	}
	this.insertFix5(node)
}

func (this *RedBlackTree) insertFix5(node *Node) {
	node.parent.colour = black
	grandparent := nodeGrandparent(node)
	grandparent.colour = red
	if node == node.parent.left && node.parent == grandparent.left {
		this.rightRotate(grandparent)
	} else if node == node.parent.right && node.parent == grandparent.right {
		this.leftRotate(grandparent)
	}
}

func (this *RedBlackTree) leftRotate(node *Node) {
	y := node.right
	node.right = y.left
	if y.left != nil {
		y.left.parent = node
	}
	y.parent = node.parent
	if node.parent == nil {
		this.root = y
	} else if node == node.parent.left {
		node.parent.left = y
	} else {
		node.parent.right = y
	}
	y.left = node
	node.parent = y
}

func (this *RedBlackTree) rightRotate(node *Node) {
	y := node.left
	node.left = y.right
	if y.right != nil {
		y.right.parent = node
	}
	y.parent = node.parent
	if node.parent == nil {
		this.root = y
	} else if node == node.parent.right {
		node.parent.right = y
	} else {
		node.parent.left = y
	}
	y.right = node
	node.parent = y
}

func (this *RedBlackTree) Flush() []vars.KeyValue {
	ret := []vars.KeyValue{}
	flushHelper(this.root, &ret)
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
