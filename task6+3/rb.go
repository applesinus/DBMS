package task6

import (
	"DBMS/task4"
	"fmt"
)

type Color bool

const (
	Red   Color = true
	Black Color = false
)

type RB struct {
	root          *nodeRB
	secondaryRoot *nodeRB
}

type nodeRB struct {
	color  Color
	key    string
	altKey *nodeRB
	left   *nodeRB
	right  *nodeRB
	parent *nodeRB

	value *task4.TrieWord
}

func (tree RB) Copy() tree {
	newTree := &RB{root: nil, secondaryRoot: nil}
	newTree.root = tree.root.copy()
	newTree.secondaryRoot = tree.secondaryRoot.copy()
	return newTree
}

func (node *nodeRB) copy() *nodeRB {
	if node == nil {
		return nil
	}

	return &nodeRB{
		color:  node.color,
		key:    node.key,
		altKey: node.altKey.copy(),
		left:   node.left.copy(),
		right:  node.right.copy(),
		parent: node.parent.copy(),
		value:  node.value,
	}
}

// returns:
//
// {node, nil} if found
//
// {nil, parent node} if not found
//
// {nil, nil} if tree is empty (or other error)
func (tree *RB) search(node *nodeRB, key string) (*nodeRB, *nodeRB) {
	if node == nil {
		return nil, nil
	}

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

func (tree *RB) update(key string, value string) string {
	node, _ := tree.search(tree.root, key)
	secondaryNode := node.altKey

	if node == nil {
		return "does not exist"
	}

	newVal, ok := task4.Pool.Insert(value)
	if ok != "ok" {
		return "error"
	}

	secondaryNode.value = newVal
	node.value = newVal
	return "ok"
}

func (tree *RB) set(key string, secondaryKey string, value string) string {
	node, parent := tree.search(tree.root, key)
	secondaryNode, secondaryParent := tree.search(tree.secondaryRoot, secondaryKey)

	if node != nil || secondaryNode != nil {
		return "exist"
	}

	newVal, ok := task4.Pool.Insert(value)
	if ok != "ok" {
		return "error"
	}
	node = &nodeRB{
		key:    key,
		parent: parent,
		altKey: nil,
		left:   nil,
		right:  nil,
		color:  Red,
		value:  newVal,
	}
	secondaryNode = &nodeRB{
		key:    secondaryKey,
		parent: secondaryParent,
		altKey: node,
		left:   nil,
		right:  nil,
		color:  Red,
		value:  newVal,
	}
	node.altKey = secondaryNode

	if parent == nil {
		tree.root = node
		node.color = Black
	}
	if secondaryParent == nil {
		tree.secondaryRoot = secondaryNode
		node.color = Black
		return "ok"
	}

	if key < parent.key {
		parent.left = node
	} else {
		parent.right = node
	}
	if secondaryKey < secondaryParent.key {
		secondaryParent.left = secondaryNode
	} else {
		secondaryParent.right = secondaryNode
	}

	response := tree.fixInsert(tree.secondaryRoot, secondaryNode)
	if response != "ok" {
		return response
	}

	return tree.fixInsert(tree.root, node)
}

func (tree *RB) remove(key string) string {
	node, _ := tree.search(tree.root, key)
	secondaryNode := node.altKey
	secondaryParent := secondaryNode.parent

	if node == nil {
		return "does not exist"
	}
	parent := node.parent

	response := tree.remover(tree.root, node, parent)
	if response != "ok" {
		return response
	}

	return tree.remover(tree.secondaryRoot, secondaryNode, secondaryParent)
}

func (tree *RB) remover(root, node, parent *nodeRB) string {
	l, r := node.left, node.right
	if l == nil && r == nil {
		if parent == nil {
			root = nil
			return "ok"
		} else if parent.left == node {
			parent.left = nil
			return tree.fixRemove(root, node, true)
		} else {
			parent.right = nil
			return tree.fixRemove(root, node, false)
		}
	} else if l == nil && node != root && r != nil {
		wasLeft := parent.left == node
		if parent == nil {
			root = r
		} else if parent.left == node {
			parent.left = r
		} else {
			parent.right = r
		}
		r.parent = parent
		r.color = node.color

		return tree.fixRemove(root, node, wasLeft)
	} else if r == nil && node != root && l != nil {
		wasLeft := parent.left == node
		if parent == nil {
			root = l
		} else if parent.left == node {
			parent.left = l
		} else {
			parent.right = l
		}
		l.parent = parent
		l.color = node.color

		return tree.fixRemove(root, node, wasLeft)
	} else {
		toDelete := tree.min(r)
		if toDelete == nil {
			toDelete = tree.max(l)
		}
		newKey := toDelete.key
		newValue := toDelete.value
		res := tree.remove(toDelete.key)
		if res != "ok" {
			return res
		}

		node.key = newKey
		node.value = newValue

		return res
	}
}

func (tree *RB) fixRemove(root, node *nodeRB, wasLeft bool) string {
	if getColor(node) == Red || node == root {
		return "ok"
	}
	fmt.Printf("fixing node: %v\n", node.key)

	for node != root && getColor(node) == Black {
		if wasLeft {
			var sibling *nodeRB
			if node == root {
				sibling = nil
			} else {
				sibling = node.parent.right
			}
			if getColor(sibling) == Red {
				sibling.color = Black
				node.parent.color = Red
				tree.leftRotate(root, node.parent)
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
					tree.rightRotate(root, sibling)
					sibling = node.parent.right
				}
				sibling.color = node.parent.color
				node.parent.color = Black
				sibling.right.color = Black
				tree.leftRotate(root, node.parent)
				node = root
			}
		} else {
			var sibling *nodeRB
			if node == root {
				sibling = nil
			} else {
				sibling = node.parent.left
			}
			if getColor(sibling) == Red {
				sibling.color = Black
				node.parent.color = Red
				tree.rightRotate(root, node.parent)
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
					tree.leftRotate(root, sibling)
					sibling = node.parent.left
				}
				sibling.color = node.parent.color
				node.parent.color = Black
				sibling.left.color = Black
				tree.rightRotate(root, node.parent)
				node = root
			}
		}
	}

	node.color = Black
	root.color = Black
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

func (tree *RB) max(node *nodeRB) *nodeRB {
	if node == nil {
		return nil
	}
	current := node
	for current.right != nil {
		current = current.right
	}

	return current
}

func (tree *RB) fixInsert(root, node *nodeRB) string {
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
					tree.leftRotate(root, node)
				}

				pa.color = Black
				gp.color = Red
				tree.rightRotate(root, gp)
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
					tree.rightRotate(root, node)
				}

				pa.color = Black
				gp.color = Red
				tree.leftRotate(root, gp)
			}
		}
	}

	root.color = Black

	return "ok"
}

func (tree *RB) leftRotate(root, node *nodeRB) {
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
		root = b
	}
	b.parent = parent

	b.left = a
	a.parent = b

	a.right = c
	if c != nil {
		c.parent = a
	}
}

func (tree *RB) rightRotate(root, node *nodeRB) {
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
		root = b
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
	fmt.Println("\nRB Tree:")
	fmt.Println("\nOrdered by key:")
	tree.root.printHelper()
	fmt.Println("\nOrdered by secondary key:")
	tree.secondaryRoot.printHelper()
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

func (tree *RB) getAll() (*[]string, *[]string, *[]string, string) {
	keys := make([]string, 0)
	values := make([]string, 0)
	secondaryKeys := make([]string, 0)
	return &keys, &values, &secondaryKeys, tree.root.getAll(&keys, &values, &secondaryKeys)
}

func (node *nodeRB) getAll(keys, secondaryKeys, values *[]string) string {
	if node == nil {
		return "ok"
	}

	*keys = append(*keys, node.key)
	*secondaryKeys = append(*secondaryKeys, node.altKey.key)
	val, ok := node.value.String()
	if !ok {
		return "error"
	}
	*values = append(*values, val)

	if node.left.getAll(keys, secondaryKeys, values) != "ok" {
		return "error"
	}
	return node.right.getAll(keys, secondaryKeys, values)
}

func (tree *RB) get(key string) (string, bool) {
	node, _ := tree.search(tree.root, key)
	if node == nil {
		return "", false
	}

	val, ok := node.value.String()
	if !ok {
		return "", false
	}
	return val, true
}

func (tree *RB) getBySecondaryKey(secondaryKey string) (string, bool) {
	node, _ := tree.search(tree.secondaryRoot, secondaryKey)
	if node == nil {
		return "", false
	}

	val, ok := node.value.String()
	if !ok {
		return "", false
	}
	return val, true
}

func (tree *RB) getRange(leftBound string, rightBound string) (*map[string]string, string) {
	result := make(map[string]string)
	return &result, tree.root.getRange(leftBound, rightBound, &result)
}

func (tree *RB) getRangeBySecondaryKey(leftBound string, rightBound string) (*map[string]string, string) {
	result := make(map[string]string)
	return &result, tree.secondaryRoot.getRange(leftBound, rightBound, &result)
}

func (node *nodeRB) getRange(leftBound string, rightBound string, result *map[string]string) (ret string) {
	defer func() {
		if ret != "ok" {
			ret = "error"
		}
	}()

	ret = "start"

	if node.key >= leftBound && node.key <= rightBound {
		val, ok := node.value.String()
		if !ok {
			return
		}
		(*result)[node.key] = val
	}

	if node.left != nil && node.key >= leftBound {
		node.left.getRange(leftBound, rightBound, result)
	}

	if node.right != nil && node.key <= rightBound {
		node.right.getRange(leftBound, rightBound, result)
	}

	ret = "ok"
	return
}
