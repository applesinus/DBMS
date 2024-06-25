package task6

func (tree *AVL) secondaryNeedToRebalance(node *nodeAVL) bool {
	return (tree.secondaryHeight(node.secondaryLeft)-tree.secondaryHeight(node.secondaryRight) > 1) || (tree.secondaryHeight(node.secondaryRight)-tree.secondaryHeight(node.secondaryLeft) > 1)
}

func (tree *AVL) secondaryHeight(node *nodeAVL) int {
	if node == nil {
		return 0
	}

	return node.secondaryHeight
}

func (tree *AVL) secondaryRebalance(node *nodeAVL) string {
	if node == nil {
		return "given nil node"
	}

	// Left rotations
	if tree.secondaryHeight(node.secondaryRight)-tree.secondaryHeight(node.secondaryLeft) > 1 {
		if tree.secondaryHeight(node.secondaryRight.secondaryRight) >= tree.secondaryHeight(node.secondaryRight.secondaryLeft) {
			// Small
			parent := node.secondaryParent
			a := node
			b := a.secondaryRight
			c := b.secondaryLeft

			if parent != nil {
				if parent.secondaryLeft == a {
					parent.secondaryLeft = b
				} else {
					parent.secondaryRight = b
				}
			} else {
				tree.secondaryRoot = b
			}
			b.secondaryParent = parent

			b.secondaryLeft = a
			a.secondaryParent = b

			a.secondaryRight = c
			if c != nil {
				c.secondaryParent = a
			}

			a.secondaryHeight -= 2

			return "ok"
		} else {
			// Big
			parent := node.secondaryParent
			a := node
			b := a.secondaryRight
			c := b.secondaryLeft
			m := c.secondaryLeft
			n := c.secondaryRight

			if parent != nil {
				if parent.secondaryLeft == a {
					parent.secondaryLeft = c
				} else {
					parent.secondaryRight = c
				}
			} else {
				tree.secondaryRoot = c
			}
			c.secondaryParent = parent

			c.secondaryLeft = a
			a.secondaryParent = c
			c.secondaryRight = b
			b.secondaryParent = c

			a.secondaryRight = m
			if m != nil {
				m.secondaryParent = a
			}
			b.secondaryLeft = n
			if n != nil {
				n.secondaryParent = b
			}

			a.secondaryHeight -= 2
			b.secondaryHeight--
			c.secondaryHeight++

			return "ok"
		}
	}

	// Right rotations
	if tree.secondaryHeight(node.secondaryLeft)-tree.secondaryHeight(node.secondaryRight) > 1 {
		if tree.secondaryHeight(node.secondaryLeft.secondaryLeft) >= tree.secondaryHeight(node.secondaryLeft.secondaryRight) {
			// Small
			parent := node.secondaryParent
			a := node
			b := a.secondaryLeft
			c := b.secondaryRight

			if parent != nil {
				if parent.secondaryLeft == a {
					parent.secondaryLeft = b
				} else {
					parent.secondaryRight = b
				}
			} else {
				tree.secondaryRoot = b
			}
			b.secondaryParent = parent

			b.secondaryRight = a
			a.secondaryParent = b

			a.secondaryLeft = c
			if c != nil {
				c.secondaryParent = a
			}

			a.secondaryHeight -= 2

			return "ok"
		} else {
			// Big
			parent := node.secondaryParent
			a := node
			b := a.secondaryLeft
			c := b.secondaryRight
			m := c.secondaryLeft
			n := c.secondaryRight

			if parent != nil {
				if parent.secondaryLeft == a {
					parent.secondaryLeft = c
				} else {
					parent.secondaryRight = c
				}
			} else {
				tree.secondaryRoot = c
			}
			c.secondaryParent = parent

			c.secondaryLeft = b
			b.secondaryParent = c
			c.secondaryRight = a
			a.secondaryParent = c

			b.secondaryRight = m
			if m != nil {
				m.secondaryParent = b
			}
			a.secondaryLeft = n
			if n != nil {
				n.secondaryParent = a
			}

			a.secondaryHeight -= 2
			b.secondaryHeight--
			c.secondaryHeight++

			return "ok"
		}
	}

	return "ok"
}

func (tree *AVL) searchBySecondaryKey(secondaryKey string) (*nodeAVL, *nodeAVL) {
	if tree.secondaryRoot == nil {
		return nil, nil
	}

	node := tree.secondaryRoot
	for node != nil {
		if secondaryKey < node.secondaryKey {
			if node.secondaryLeft == nil {
				return nil, node
			}
			node = node.secondaryLeft
		} else if secondaryKey > node.secondaryKey {
			if node.secondaryRight == nil {
				return nil, node
			}
			node = node.secondaryRight
		} else {
			return node, nil
		}
	}

	return nil, nil
}

func (tree *AVL) changeSecondaryHeights(node *nodeAVL) string {
	for node != nil {
		node.secondaryHeight = max(tree.secondaryHeight(node.secondaryLeft), tree.secondaryHeight(node.secondaryRight)) + 1

		if tree.needToRebalance(node) {
			ok := tree.rebalance(node)
			if ok != "ok" {
				return "balance error: " + ok
			}
			break
		}

		node = node.secondaryParent
	}

	return "ok"
}

func (tree *AVL) secondaryRemove(secondaryKey string, node *nodeAVL) string {
	parent := node.secondaryParent

	if node.secondaryLeft == nil && node.secondaryRight == nil {
		// Node with no children
		if parent == nil {
			tree.secondaryRoot = nil
			return "ok"
		} else {
			if parent.secondaryLeft == node {
				parent.secondaryLeft = nil
			} else {
				parent.secondaryRight = nil
			}

			return tree.changeHeights(parent)
		}
	} else if node.secondaryLeft == nil {
		// Node with only right child
		if parent == nil {
			tree.secondaryRoot = node.secondaryRight
			node.secondaryRight.secondaryParent = nil
			return "ok"
		} else {
			if parent.secondaryLeft == node {
				parent.secondaryLeft = node.secondaryRight
			} else {
				parent.secondaryRight = node.secondaryRight
			}
			node.secondaryRight.secondaryParent = parent

			parent.secondaryHeight = max(tree.secondaryHeight(node.secondaryParent.secondaryLeft), tree.secondaryHeight(node.secondaryParent.secondaryRight)) + 1
			return tree.changeHeights(parent)
		}
	} else if node.secondaryRight == nil {
		// Node with only left child
		if parent == nil {
			tree.secondaryRoot = node.secondaryLeft
			node.secondaryLeft.secondaryParent = nil
			return "ok"
		} else {
			if parent.secondaryLeft == node {
				parent.secondaryLeft = node.secondaryLeft
			} else {
				parent.secondaryRight = node.secondaryLeft
			}
			node.secondaryLeft.secondaryParent = parent

			parent.secondaryHeight = max(tree.secondaryHeight(parent.secondaryLeft), tree.secondaryHeight(parent.secondaryRight)) + 1
			return tree.changeHeights(parent)
		}
	} else {
		// Node with two children
		successor, successorParent := tree.secondaryMin(node.secondaryRight)
		if successor == nil && successorParent == nil {
			return "not found children... somehow"
		}
		if successorParent == nil {
			successorParent = successor.secondaryParent
		}

		node.secondaryKey = successor.secondaryKey
		node.value = successor.value

		if successorParent.secondaryLeft == successor {
			if successorParent.secondaryLeft.secondaryRight != nil {
				successorParent.secondaryLeft = successorParent.secondaryLeft.secondaryRight
			} else {
				successorParent.secondaryLeft = nil
			}
		} else {
			successorParent.secondaryRight = nil
		}

		successorParent.secondaryHeight = max(tree.secondaryHeight(successorParent.secondaryLeft), tree.secondaryHeight(successorParent.secondaryRight)) + 1
		return tree.changeHeights(successorParent)
	}
}

func (tree *AVL) secondaryMin(node *nodeAVL) (*nodeAVL, *nodeAVL) {
	if node == nil {
		return nil, nil
	}
	current := node
	for current.secondaryLeft != nil {
		current = current.secondaryLeft
	}

	return current, current.secondaryParent
}

// returns:
//
// {node, nil} if found
//
// {nil, parent node} if not found
//
// {nil, nil} if tree is empty (or other error)
func (tree *AVL) secondarySearch(secondaryKey string) (*nodeAVL, *nodeAVL) {
	if tree.secondaryRoot == nil {
		return nil, nil
	}

	node := tree.secondaryRoot
	for node != nil {
		if secondaryKey < node.secondaryKey {
			if node.secondaryLeft == nil {
				return nil, node
			}
			node = node.secondaryLeft
		} else if secondaryKey > node.secondaryKey {
			if node.secondaryRight == nil {
				return nil, node
			}
			node = node.secondaryRight
		} else {
			return node, nil
		}
	}

	return nil, nil
}

func (tree *AVL) getBySecondaryKey(secondaryKey string) (string, bool) {
	node, _ := tree.secondarySearch(secondaryKey)
	if node == nil {
		return "", false
	}
	val, ok := node.value.String()
	if !ok {
		return "", false
	}
	return val, true
}

func (tree *AVL) getRangeBySecondaryKey(leftBound string, rightBound string) (*map[string]string, string) {
	result := make(map[string]string)
	return &result, tree.secondaryRoot.getRangeBySecondaryKey(leftBound, rightBound, &result)
}

func (node *nodeAVL) getRangeBySecondaryKey(leftBound string, rightBound string, result *map[string]string) (ret string) {
	defer func() {
		if ret != "ok" {
			ret = "error"
		}
	}()

	ret = "start"

	if node.secondaryKey >= leftBound && node.secondaryKey <= rightBound {
		val, ok := node.value.String()
		if !ok {
			return
		}
		(*result)[node.secondaryKey] = val
	}

	if node.secondaryLeft != nil && node.secondaryKey >= leftBound {
		node.secondaryLeft.getRange(leftBound, rightBound, result)
	}

	if node.secondaryRight != nil && node.secondaryKey <= rightBound {
		node.secondaryRight.getRange(leftBound, rightBound, result)
	}

	ret = "ok"
	return
}
