package task6

import "fmt"

// returns:
//
// {node, nil} if found
//
// {nil, parent node} if not found
//
// {nil, nil} if tree is empty (or other error)
func (tree *RB) searchBySecondaryKey(secondaryKey string) (*nodeRB, *nodeRB) {
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

func (tree *RB) fixSecondaryInsert(node *nodeRB) string {
	if getSecondaryColor(node.secondaryParent) == Black || node.secondaryParent == nil {
		return "ok"
	}

	for node.secondaryParent != nil && getSecondaryColor(node.secondaryParent) == Red {
		un, gp, pa := tree.secondaryUncle(node), tree.secondaryGrandparent(node), node.secondaryParent
		if gp.secondaryLeft == pa {
			if un != nil && getSecondaryColor(un) == Red {
				pa.secondaryColor = Black
				un.secondaryColor = Black
				gp.secondaryColor = Red
				node = gp
			} else {
				if node == pa.secondaryRight {
					node = pa
					gp, pa = tree.secondaryGrandparent(node), node.secondaryParent
					tree.secondaryLeftRotate(node)
				}

				pa.secondaryColor = Black
				gp.secondaryColor = Red
				tree.secondaryRightRotate(gp)
			}
		} else {
			if un != nil && getSecondaryColor(un) == Red {
				pa.secondaryColor = Black
				un.secondaryColor = Black
				gp.secondaryColor = Red
				node = gp
			} else {
				if node == pa.secondaryLeft {
					node = pa
					gp, pa = tree.secondaryGrandparent(node), node.secondaryParent
					tree.secondaryRightRotate(node)
				}

				pa.secondaryColor = Black
				gp.secondaryColor = Red
				tree.secondaryLeftRotate(gp)
			}
		}
	}

	tree.secondaryRoot.secondaryColor = Black

	return "ok"
}

func (tree *RB) secondaryLeftRotate(node *nodeRB) {
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
}

func (tree *RB) secondaryRightRotate(node *nodeRB) {
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
}

func getSecondaryColor(node *nodeRB) Color {
	if node == nil {
		return Black
	}
	return node.secondaryColor
}

func (tree *RB) secondaryMin(node *nodeRB) *nodeRB {
	if node == nil {
		return nil
	}
	current := node
	for current.secondaryLeft != nil {
		current = current.secondaryLeft
	}

	return current
}

func (tree *RB) secondarySibling(node *nodeRB) *nodeRB {
	if node == nil || node.secondaryParent == nil {
		return nil
	}

	if node.secondaryParent.secondaryLeft == node {
		return node.secondaryParent.secondaryRight
	}

	return node.secondaryParent.secondaryLeft
}

func (tree *RB) secondaryGrandparent(node *nodeRB) *nodeRB {
	if node == nil || node.secondaryParent == nil || node.secondaryParent.secondaryParent == nil {
		return nil
	}

	return node.secondaryParent.secondaryParent
}

func (tree *RB) secondaryUncle(node *nodeRB) *nodeRB {
	gp := tree.secondaryGrandparent(node)
	if gp == nil {
		return nil
	}

	if node.secondaryParent == gp.secondaryLeft {
		return gp.secondaryRight
	}

	return gp.secondaryLeft
}

func (tree *RB) getBySecondaryKey(secondaryKey string) (string, bool) {
	node, _ := tree.searchBySecondaryKey(secondaryKey)
	if node == nil {
		return "", false
	}

	val, ok := node.value.String()
	if !ok {
		return "", false
	}
	return val, true
}

func (tree *RB) getRangeBySecondaryKey(leftBound string, rightBound string) (*map[string]string, string) {
	result := make(map[string]string)
	return &result, tree.secondaryRoot.getRangeBySecondaryKey(leftBound, rightBound, &result)
}

func (node *nodeRB) getRangeBySecondaryKey(leftBound string, rightBound string, result *map[string]string) (ret string) {
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
		node.secondaryLeft.getRangeBySecondaryKey(leftBound, rightBound, result)
	}

	if node.secondaryRight != nil && node.secondaryKey <= rightBound {
		node.secondaryRight.getRangeBySecondaryKey(leftBound, rightBound, result)
	}

	ret = "ok"
	return
}

func (tree *RB) fixSecondaryRemove(node *nodeRB, wasLeft bool) string {
	if getSecondaryColor(node) == Red || node == tree.secondaryRoot {
		return "ok"
	}
	fmt.Printf("fixing node: %v\n", node.secondaryKey)

	for node != tree.secondaryRoot && getSecondaryColor(node) == Black {
		if wasLeft {
			var sibling *nodeRB
			if node == tree.secondaryRoot {
				sibling = nil
			} else {
				sibling = node.secondaryParent.secondaryRight
			}
			if getSecondaryColor(sibling) == Red {
				sibling.secondaryColor = Black
				node.secondaryParent.secondaryColor = Red
				tree.secondaryLeftRotate(node.secondaryParent)
				sibling = node.secondaryParent.secondaryRight
			}
			if getSecondaryColor(sibling.secondaryLeft) == Black && getSecondaryColor(sibling.secondaryRight) == Black {
				sibling.secondaryColor = Red
				if node == node.secondaryParent.secondaryLeft {
					wasLeft = true
				} else {
					wasLeft = false
				}
				node = node.secondaryParent
			} else {
				if getSecondaryColor(sibling.secondaryRight) == Black {
					sibling.secondaryLeft.secondaryColor = Black
					sibling.secondaryColor = Red
					tree.secondaryRightRotate(sibling)
					sibling = node.secondaryParent.secondaryRight
				}
				sibling.secondaryColor = node.secondaryParent.secondaryColor
				node.secondaryParent.secondaryColor = Black
				sibling.secondaryRight.secondaryColor = Black
				tree.secondaryLeftRotate(node.secondaryParent)
				node = tree.secondaryRoot
			}
		} else {
			var sibling *nodeRB
			if node == tree.secondaryRoot {
				sibling = nil
			} else {
				sibling = node.secondaryParent.secondaryLeft
			}
			if getSecondaryColor(sibling) == Red {
				sibling.secondaryColor = Black
				node.secondaryParent.secondaryColor = Red
				tree.secondaryRightRotate(node.secondaryParent)
				sibling = node.secondaryParent.secondaryLeft
			}
			if getSecondaryColor(sibling.secondaryRight) == Black && getSecondaryColor(sibling.secondaryLeft) == Black {
				sibling.secondaryColor = Red
				if node == node.secondaryParent.secondaryRight {
					wasLeft = false
				} else {
					wasLeft = true
				}
				node = node.secondaryParent
			} else {
				if getSecondaryColor(sibling.secondaryLeft) == Black {
					sibling.secondaryRight.secondaryColor = Black
					sibling.secondaryColor = Red
					tree.secondaryLeftRotate(sibling)
					sibling = node.secondaryParent.secondaryLeft
				}
				sibling.secondaryColor = node.secondaryParent.secondaryColor
				node.secondaryParent.secondaryColor = Black
				sibling.secondaryLeft.secondaryColor = Black
				tree.secondaryRightRotate(node.secondaryParent)
				node = tree.secondaryRoot
			}
		}
	}

	node.secondaryColor = Black
	tree.secondaryRoot.secondaryColor = Black
	return "ok"
}

func (tree *RB) secondaryRemove(secondaryKey string) string {
	node, _ := tree.searchBySecondaryKey(secondaryKey)

	if node == nil {
		return "does not exist"
	}
	parent := node.secondaryParent

	l, r := node.secondaryLeft, node.secondaryRight
	if l == nil && r == nil {
		if parent == nil {
			tree.secondaryRoot = nil
			return "ok"
		} else if parent.secondaryLeft == node {
			parent.secondaryLeft = nil
			return tree.fixSecondaryRemove(node, true)
		} else {
			parent.secondaryRight = nil
			return tree.fixSecondaryRemove(node, false)
		}
	} else if l == nil && node != tree.secondaryRoot && r != nil {
		wasLeft := parent.secondaryLeft == node
		if parent == nil {
			tree.secondaryRoot = r
		} else if parent.secondaryLeft == node {
			parent.secondaryLeft = r
		} else {
			parent.secondaryRight = r
		}
		r.secondaryParent = parent
		r.secondaryColor = node.secondaryColor

		return tree.fixSecondaryRemove(node, wasLeft)
	} else if r == nil && node != tree.secondaryRoot && l != nil {
		wasLeft := parent.secondaryLeft == node
		if parent == nil {
			tree.secondaryRoot = l
		} else if parent.secondaryLeft == node {
			parent.secondaryLeft = l
		} else {
			parent.secondaryRight = l
		}
		l.secondaryParent = parent
		l.secondaryColor = node.secondaryColor

		return tree.fixSecondaryRemove(node, wasLeft)
	} else {
		toDelete := tree.secondaryMin(r)
		if toDelete == nil {
			toDelete = tree.secondaryMax(l)
		}
		newKey := toDelete.secondaryKey
		newValue := toDelete.value
		res := tree.secondaryRemove(toDelete.secondaryKey)
		if res != "ok" {
			return res
		}

		node.secondaryKey = newKey
		node.value = newValue

		return res
	}
}

func (tree *RB) secondaryMax(node *nodeRB) *nodeRB {
	if node == nil {
		return nil
	}
	current := node
	for current.secondaryRight != nil {
		current = current.secondaryRight
	}

	return current
}
