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
	fmt.Printf("removing %v\n", key)
	node, _ := tree.search(key)

	if node == nil {
		return "does not exist"
	}
	parent := node.parent

	l, r := node.left, node.right
	if l == nil && r == nil {
		if parent == nil {
			tree.root = nil
			return "ok"
		} else if parent.left == node {
			parent.left = nil
			return tree.fixRemove(node, true)
		} else {
			parent.right = nil
			return tree.fixRemove(node, false)
		}
	} else if l == nil {
		wasLeft := parent.left == node
		if parent == nil {
			tree.root = r
		} else if parent.left == node {
			parent.left = r
		} else {
			parent.right = r
		}
		r.parent = parent
		r.color = node.color

		return tree.fixRemove(node, wasLeft)
	} else if r == nil {
		wasLeft := parent.left == node
		if parent == nil {
			tree.root = l
		} else if parent.left == node {
			parent.left = l
		} else {
			parent.right = l
		}
		l.parent = parent
		l.color = node.color

		return tree.fixRemove(node, wasLeft)
	} else {
		min := tree.min(r)
		if min == nil {
			return "error"
		}
		newKey := min.key
		newValue := min.value
		res := tree.remove(min.key)
		if res != "ok" {
			return res
		}

		node.key = newKey
		node.value = newValue

		return res
	}
}

func (tree *RB) fixRemove(node *nodeRB, wasLeft bool) string {
	if getColor(node) == Red || node == tree.root {
		return "ok"
	}
	fmt.Printf("fixing node: %v\n", node.key)

	for node != tree.root && getColor(node) == Black {
		if wasLeft {
			var sibling *nodeRB
			if node == tree.root {
				sibling = nil
			} else {
				sibling = node.parent.right
			}
			if getColor(sibling) == Red {
				sibling.color = Black
				node.parent.color = Red
				tree.leftRotate(node.parent)
				sibling = node.parent.right
			}
			if getColor(sibling.left) == Black && getColor(sibling.right) == Black {
				sibling.color = Red
				if node == node.parent.left {
					wasLeft = true
				} else {
					wasLeft = false
				}
				node = node.parent
			} else {
				if getColor(sibling.right) == Black {
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
			var sibling *nodeRB
			if node == tree.root {
				sibling = nil
			} else {
				sibling = node.parent.left
			}
			if getColor(sibling) == Red {
				sibling.color = Black
				node.parent.color = Red
				tree.rightRotate(node.parent)
				sibling = node.parent.left
			}
			if getColor(sibling.right) == Black && getColor(sibling.left) == Black {
				sibling.color = Red
				if node == node.parent.right {
					wasLeft = false
				} else {
					wasLeft = true
				}
				node = node.parent
			} else {
				if getColor(sibling.left) == Black {
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
	tree.root.color = Black
	return "ok"
}

func getColor(node *nodeRB) Color {
	if node == nil {
		return Black
	}
	return node.color
}

func (tree *RB) min(node *nodeRB) *nodeRB {
	if node == nil {
		return nil
	}
	current := node
	for current.left != nil {
		current = current.left
	}

	return current
}

func (tree *RB) fixInsert(node *nodeRB) string {
	if getColor(node.parent) == Black || node.parent == nil {
		return "ok"
	}

	for node.parent != nil && getColor(node.parent) == Red {
		un, gp, pa := tree.uncle(node), tree.grandparent(node), node.parent
		if gp.left == pa {
			if un != nil && getColor(un) == Red {
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
			if un != nil && getColor(un) == Red {
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
	tree.root.printHelper()
}

func (node *nodeRB) printHelper() {
	if node == nil {
		return
	}

	var c string
	if getColor(node) == Red {
		c = "Red"
	} else {
		c = "Black"
	}
	if node.parent != nil {
		fmt.Printf("Node: %v:%v, c=%v. (parent: %v, ", node.key, node.value, c, node.parent.key)
	} else {
		fmt.Printf("Node: %v:%v c=%v. (parent: nil, ", node.key, node.value, c)
	}
	if node.left != nil {
		fmt.Printf("Left: %v, ", node.left.key)
	}
	if node.right != nil {
		fmt.Printf("Right: %v, ", node.right.key)
	}
	fmt.Printf(")\n")

	node.left.printHelper()
	node.right.printHelper()
}

func (tree *RB) find(key string) (interface{}, bool) {
	node, _ := tree.search(key)
	if node == nil {
		return nil, false
	}
	return node.value, true
}