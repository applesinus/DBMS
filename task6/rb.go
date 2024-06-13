package task6

import "fmt"

type Color bool

const (
	Red   Color = true
	Black Color = false
)

type RB struct {
	root *nodeRB
}

type nodeRB struct {
	color  Color
	key    string
	left   *nodeRB
	right  *nodeRB
	parent *nodeRB
	value  interface{}
}

// returns:
//
// {node, nil} if found
//
// {nil, parent node} if not found
//
// {nil, nil} if tree is empty (or other error)
func (tree *RB) search(key string) (*nodeRB, *nodeRB) {
	if tree.root == nil {
		return nil, nil
	}

	node := tree.root
	for node != nil {
		if key < node.key {
			if node.left == nil {
				return nil, node
			}
			node = node.left
		} else if key > node.key {
			if node.right == nil {
				return nil, node
			}
			node = node.right
		} else {
			return node, nil
		}
	}

	return nil, nil
}

func (tree *RB) update(key string, value interface{}) string {
	node, _ := tree.search(key)

	if node == nil {
		return "does not exist"
	}

	node.value = value
	return "ok"
}

func (tree *RB) insert(key string, value interface{}) string {
	node, parent := tree.search(key)

	if node != nil {
		return "exist"
	}

	node = &nodeRB{
		key:    key,
		value:  value,
		parent: parent,
		left:   nil,
		right:  nil,
		color:  Red,
	}

	if parent == nil {
		tree.root = node
		node.color = Black
		return "ok"
	}

	if key < parent.key {
		parent.left = node
	} else {
		parent.right = node
	}

	return tree.fixInsert(node)
}

func (tree *RB) remove(key string) string {
	node, _ := tree.search(key)

	if node == nil {
		return "does not exist"
	}
	parent := node.parent

	l, r := node.left, node.right
	if l == nil && r == nil {
		if parent == nil {
			tree.root = nil
		} else if parent.left == node {
			parent.left = nil
		} else {
			parent.right = nil
		}
	} else if l == nil {
		if parent == nil {
			tree.root = r
		} else if parent.left == node {
			parent.left = r
		} else {
			parent.right = r
		}
		r.parent = parent
		r.color = node.color
	} else if r == nil {
		if parent == nil {
			tree.root = l
		} else if parent.left == node {
			parent.left = l
		} else {
			parent.right = l
		}
		l.parent = parent
		l.color = node.color
	} else {
		min, minParent := tree.min(r)

		if min.right != nil {
			if minParent == node {
				minParent.right = min.right
			} else {
				minParent.left = min.right
			}
			min.right.parent = minParent
			min.right.color = min.color
		} else {
			minParent.left = nil
		}

		node.key = min.key
		node.value = min.value

		return tree.fixRemove(node)
	}

	return "ok"
}

func (tree *RB) fixRemove(node *nodeRB) string {
	if node.color == Black || node == tree.root {
		return "ok"
	}

	for node != nil && node.color == Black {
		if node == node.parent.left {
			sibling := node.parent.right
			if sibling.color == Red {
				sibling.color = Black
				node.parent.color = Red
				tree.leftRotate(node.parent)
				sibling = node.parent.right
			}
			if sibling.left.color == Black && sibling.right.color == Black {
				sibling.color = Red
				node = node.parent
			} else {
				if sibling.right.color == Black {
					sibling.left.color = Black
					sibling.color = Red
					tree.rightRotate(sibling)
					sibling = node.parent.right
				}
				sibling.color = node.parent.color
				node.parent.color = Black
				sibling.right.color = Black
				tree.leftRotate(node.parent)
				node = tree.root
			}
		} else {
			sibling := node.parent.left
			if sibling.color == Red {
				sibling.color = Black
				node.parent.color = Red
				tree.rightRotate(node.parent)
				sibling = node.parent.left
			}
			if sibling.right.color == Black && sibling.left.color == Black {
				sibling.color = Red
				node = node.parent
			} else {
				if sibling.left.color == Black {
					sibling.right.color = Black
					sibling.color = Red
					tree.leftRotate(sibling)
					sibling = node.parent.left
				}
				sibling.color = node.parent.color
				node.parent.color = Black
				sibling.left.color = Black
				tree.rightRotate(node.parent)
				node = tree.root
			}
		}
	}

	node.color = Black
	return "ok"
}

func (tree *RB) min(node *nodeRB) (*nodeRB, *nodeRB) {
	if node == nil {
		return nil, nil
	}
	current := node
	for current.left != nil {
		current = current.left
	}

	return current, current.parent
}

func (tree *RB) fixInsert(node *nodeRB) string {
	if node.parent.color == Black {
		return "ok"
	}

	for node.parent.color == Red {
		un, gp, pa := tree.uncle(node), tree.grandparent(node), node.parent
		if gp.left == pa {
			if un != nil && un.color == Red {
				pa.color = Black
				un.color = Black
				gp.color = Red
				node = gp
			} else {
				if node == pa.right {
					node = pa
					gp, pa = tree.grandparent(node), node.parent
					tree.leftRotate(node)
				}

				pa.color = Black
				gp.color = Red
				tree.rightRotate(gp)
			}
		} else {
			if un != nil && un.color == Red {
				pa.color = Black
				un.color = Black
				gp.color = Red
				node = gp
			} else {
				if node == pa.left {
					node = pa
					gp, pa = tree.grandparent(node), node.parent
					tree.rightRotate(node)
				}

				pa.color = Black
				gp.color = Red
				tree.leftRotate(gp)
			}
		}
	}

	tree.root.color = Black

	return "ok"
}

func (tree *RB) leftRotate(node *nodeRB) {
	parent := node.parent
	a := node
	b := a.right
	c := b.left

	if parent != nil {
		if parent.left == a {
			parent.left = b
		} else {
			parent.right = b
		}
	} else {
		tree.root = b
	}
	b.parent = parent

	b.left = a
	a.parent = b

	a.right = c
	if c != nil {
		c.parent = a
	}
}

func (tree *RB) rightRotate(node *nodeRB) {
	parent := node.parent
	a := node
	b := a.left
	c := b.right

	if parent != nil {
		if parent.left == a {
			parent.left = b
		} else {
			parent.right = b
		}
	} else {
		tree.root = b
	}
	b.parent = parent

	b.right = a
	a.parent = b

	a.left = c
	if c != nil {
		c.parent = a
	}
}

func (tree *RB) sibling(node *nodeRB) *nodeRB {
	if node == nil || node.parent == nil {
		return nil
	}

	if node.parent.left == node {
		return node.parent.right
	}

	return node.parent.left
}

func (tree *RB) grandparent(node *nodeRB) *nodeRB {
	if node == nil || node.parent == nil || node.parent.parent == nil {
		return nil
	}

	return node.parent.parent
}

func (tree *RB) uncle(node *nodeRB) *nodeRB {
	gp := tree.grandparent(node)
	if gp == nil {
		return nil
	}

	if node.parent == gp.left {
		return gp.right
	}

	return gp.left
}

func (tree *RB) print() {
	fmt.Println("\nRB Tree:\n")
	tree.printHelper(tree.root)
}

func (tree *RB) printHelper(node *nodeRB) {
	if node == nil {
		return
	}

	if node.parent != nil {
		fmt.Printf("Node: %v:%v, c=%v. (parent: %v, ", node.key, node.value, node.color, node.parent.key)
	} else {
		fmt.Printf("Node: %v:%v c=%v. (parent: nil, ", node.key, node.value, node.color)
	}
	if node.left != nil {
		fmt.Printf("Left: %v, ", node.left.key)
	}
	if node.right != nil {
		fmt.Printf("Right: %v, ", node.right.key)
	}
	fmt.Printf(")\n")

	tree.printHelper(node.left)
	tree.printHelper(node.right)
}
