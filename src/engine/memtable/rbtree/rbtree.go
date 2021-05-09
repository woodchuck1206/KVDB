package rbtree

var (
	black colour = false
	red   colour = true
)

type colour bool

type RedBlackTree struct {
	root *Node
}

type Node struct {
	key    string
	value  string
	colour colour
	parent *Node
	left   *Node
	right  *Node
}

func (rbtree *RedBlackTree) Insert(key string, value string) {
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

	newNode := &Node{
		key:    key,
		value:  value,
		colour: red,
		parent: y,
		left:   nil,
		right:  nil,
	}

	if y == nil {
		rbtree.root = newNode
	} else if key < y.key {
		y.left = newNode
	} else {
		y.right = newNode
	} // if same -> should follow a different flow, maybe overlap

	rbtree.insertFix(newNode)
}

func (rbtree *RedBlackTree) insertFix(node *Node) {
	for node.parent != nil && node.parent.colour == red {
		if node.parent.parent != nil && node.parent == node.parent.parent.left {
			uncle := node.parent.parent.right
			if uncle.colour == red {
				node.parent.colour = black
				uncle.colour = black
				node.parent.parent.colour = red
				node = node.parent.parent
			} else if node == node.parent.right {
				node = node.parent
				leftRotate(rbtree, node)
			}
			node.parent.colour = black
			node.parent.parent.colour = red
			rightRotate(rbtree, node.parent.parent)
		} else if node.parent.parent != nil && node.parent == node.parent.parent.right {
			uncle := node.parent.parent.left
			if uncle.colour == red {
				node.parent.colour = black
				uncle.colour = black
				node.parent.parent.colour = red
				node = node.parent.parent
			} else if node == node.parent.left {
				node = node.parent
				leftRotate(rbtree, node)
			}
			node.parent.colour = black
			node.parent.parent.colour = red
			rightRotate(rbtree, node.parent.parent)
		}
		rbtree.root.colour = black
	}
}

func leftRotate(rbtree *RedBlackTree, node *Node) {

}

func rightRotate(rbtree *RedBlackTree, node *Node) {

}
