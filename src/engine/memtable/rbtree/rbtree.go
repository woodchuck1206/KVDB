package rbtree

import (
	"errors"

	myVar "github.com/woodchuckchoi/KVDB/src/engine/vars"
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
		return "", errors.New("NOT EXIST")
	}

	if node.key < key {
		return node.left.Get(key)
	}
	return node.right.Get(key)
}

// NODE ---

// RBTREE ---

type RedBlackTree struct {
	root *Node
}

func (rbtree *RedBlackTree) Get(key string) (string, error) {
	node := rbtree.root
	return node.Get(key)
}

func (rbtree *RedBlackTree) Insert(key string, value string) {
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
		return
	}
	if node.parent == nil {
		node.colour = black
		return
	}
	if node.parent.parent == nil {
		return
	}
	rbtree.insertFix(node)
}

func (rbtree *RedBlackTree) insertFix(node *Node) {
	for node.parent != nil && node.parent.colour == red {
		if node.parent.parent != nil && node.parent.parent.left != nil && node.parent == node.parent.parent.left {
			uncle := node.parent.parent.right
			if uncle != nil && uncle.colour == red {
				node.parent.colour = black
				uncle.colour = black
				node.parent.parent.colour = red
				node = node.parent.parent
			} else if node.parent.right != nil && node == node.parent.right {
				node = node.parent
				rbtree.leftRotate(node)
			}
			node.parent.colour = black
			node.parent.parent.colour = red
			rbtree.rightRotate(node.parent.parent)
		} else if node.parent.parent != nil && node.parent.parent.right != nil && node.parent == node.parent.parent.right {
			uncle := node.parent.parent.left
			if uncle != nil && uncle.colour == red {
				node.parent.colour = black
				uncle.colour = black
				node.parent.parent.colour = red
				node = node.parent.parent
			} else if node.parent.left != nil && node == node.parent.left {
				node = node.parent
				rbtree.rightRotate(node)
			}
			node.parent.colour = black
			node.parent.parent.colour = red
			rbtree.leftRotate(node.parent.parent)
		}
		rbtree.root.colour = black
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

func (rbtree *RedBlackTree) Flush() []myVar.KeyValue {
	ret := []myVar.KeyValue{}
	flushHelper(rbtree.root, &ret)
	return ret
}

func flushHelper(node *Node, out *[]myVar.KeyValue) {
	if node == nil {
		return
	}
	flushHelper(node.left, out)
	*out = append(*out, myVar.KeyValue{Key: node.key, Value: node.value})
	flushHelper(node.right, out)
}

// RBTREE ---
