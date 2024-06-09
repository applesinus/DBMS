package task6

type AVL struct {
	root *nodeAVL
}

type nodeAVL struct {
	key    string
	value  interface{}
	left   *nodeAVL
	right  *nodeAVL
	parent *nodeAVL
	height int
}

func (tree *AVL) needToRebalance(node *nodeAVL) bool {
	if (tree.height(node.left)-tree.height(node.right) > 1) || (tree.height(node.right)-tree.height(node.left) > 1) {
		return true
	}

	return false
}

func (tree *AVL) height(node *nodeAVL) int {
	if node == nil {
		return 0
	}

	return node.height
}

func (tree *AVL) rebalance(node *nodeAVL) string {
	if node == nil {
		return "given nil node"
	}

	// Left rotations
	if tree.height(node.right)-tree.height(node.left) > 1 {
		if tree.height(node.right.right) >= tree.height(node.right.left) {
			// Small
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

			a.height -= 2

			return "ok"
		} else {
			// Big
			parent := node.parent
			a := node
			b := a.right
			c := b.left
			m := c.left
			n := c.right

			if parent != nil {
				if parent.left == a {
					parent.left = c
				} else {
					parent.right = c
				}
			} else {
				tree.root = c
			}
			c.parent = parent

			c.left = a
			a.parent = c
			c.right = b
			b.parent = c

			a.right = m
			if m != nil {
				m.parent = a
			}
			b.left = n
			if n != nil {
				n.parent = b
			}

			a.height -= 2
			b.height--
			c.height++

			return "ok"
		}
	}

	// Right rotations
	if tree.height(node.left)-tree.height(node.right) > 1 {
		if tree.height(node.left.left) >= tree.height(node.left.right) {
			// Small
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

			a.height -= 2

			return "ok"
		} else {
			// Big
			parent := node.parent
			a := node
			b := a.left
			c := b.right
			m := c.left
			n := c.right

			if parent != nil {
				if parent.left == a {
					parent.left = c
				} else {
					parent.right = c
				}
			} else {
				tree.root = c
			}
			c.parent = parent

			c.left = b
			b.parent = c
			c.right = a
			a.parent = c

			b.right = m
			if m != nil {
				m.parent = b
			}
			a.left = n
			if n != nil {
				n.parent = a
			}

			a.height -= 2
			b.height--
			c.height++

			return "ok"
		}
	}

	return "ok"
}

// returns:
//
// {node, nil} if found
//
// {nil, parent node} if not found
//
// {nil, nil} if tree is empty (or other error)
func (tree *AVL) search(key string) (*nodeAVL, *nodeAVL) {
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

func (tree *AVL) changeHeights(node *nodeAVL) string {
	for node.parent != nil {
		if tree.needToRebalance(node.parent) {
			ok := tree.rebalance(node.parent)
			if ok != "ok" {
				return "balance error: " + ok
			}
			break
		}

		node = node.parent
		node.height = max(tree.height(node.left), tree.height(node.right)) + 1
	}

	return "ok"
}

func (tree *AVL) insert(key string, value interface{}) string {
	node, parent := tree.search(key)

	if node != nil {
		return "exist"
	}

	newNode := &nodeAVL{
		key:    key,
		value:  value,
		parent: parent,
		left:   nil,
		right:  nil,
		height: 1,
	}

	if parent == nil {
		tree.root = newNode
		return "ok"
	}

	if key < parent.key {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	return tree.changeHeights(newNode)
}

func (tree *AVL) update(key string, value interface{}) string {
	node, _ := tree.search(key)

	if node == nil {
		return "does not exist"
	}

	node.value = value
	return "ok"
}

func (tree *AVL) delete(key string) string {
	node, parent := tree.search(key)

	if node == nil {
		return "does not exist"
	}

	// if node has no children
	if node.left == nil && node.right == nil {
		if parent == nil {
			tree.root = nil
		} else if parent.left == node {
			parent.left = nil
		} else {
			parent.right = nil
		}
		return "ok"
	}

	// if node has children
	left := node.left
	for left.right != nil || left != nil {
		left = left.right
	}
	right := node.right
	for right.left != nil || right != nil {
		right = right.left
	}

	target := node
	target = nil

	if node.left {
	}
	for i := 0; i < len(right.key) && i < len(left.key) && i < len(node.key); i++ {
		c := int(node.key[i])
		l := int(left.key[i])
		r := int(right.key[i])
		if l-c == r-l {
			continue
		} else if l-c > r-l {
			target = left
			break
		} else {
			target = right
			break
		}
	}
	if target == nil {
		if len(right.key) > len(left.key) {
			target = right
		} else {
			target = left
		}
	}

	return "ok"
}
