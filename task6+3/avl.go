package task6

import (
	"DBMS/task4"
	"fmt"
)

type AVL struct {
	root          *nodeAVL
	secondaryRoot *nodeAVL
}

type nodeAVL struct {
	key    string
	left   *nodeAVL
	right  *nodeAVL
	parent *nodeAVL
	height int

	secondaryKey    string
	secondaryLeft   *nodeAVL
	secondaryRight  *nodeAVL
	secondaryParent *nodeAVL
	secondaryHeight int

	value task4.TrieWord
}

func (tree AVL) Copy() tree {
	newTree := &AVL{root: nil, secondaryRoot: nil}
	newTree.root = tree.root.copy()
	return newTree
}

func (node *nodeAVL) copy() *nodeAVL {
	if node == nil {
		return nil
	}

	return &nodeAVL{
		key:    node.key,
		left:   node.left.copy(),
		right:  node.right.copy(),
		parent: node.parent.copy(),
		height: node.height,

		secondaryKey:    node.secondaryKey,
		secondaryLeft:   node.secondaryLeft.copy(),
		secondaryRight:  node.secondaryRight.copy(),
		secondaryParent: node.secondaryParent.copy(),
		secondaryHeight: node.secondaryHeight,

		value: node.value,
	}
}

func (tree *AVL) needToRebalance(node *nodeAVL) bool {
	return (tree.height(node.left)-tree.height(node.right) > 1) || (tree.height(node.right)-tree.height(node.left) > 1)
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
	for node != nil {
		node.height = max(tree.height(node.left), tree.height(node.right)) + 1

		if tree.needToRebalance(node) {
			ok := tree.rebalance(node)
			if ok != "ok" {
				return "balance error: " + ok
			}
			break
		}

		node = node.parent
	}

	return "ok"
}

func (tree *AVL) set(key string, secondaryKey string, value string) string {
	node, parent := tree.search(key)
	secondaryNode, secondaryParent := tree.searchBySecondaryKey(secondaryKey)

	if node != nil || secondaryNode != nil {
		return "exist"
	}

	newVal, ok := task4.Pool.Insert(value)

	if ok != "ok" {
		return "error"
	}

	node = &nodeAVL{
		key:    key,
		parent: parent,
		left:   nil,
		right:  nil,
		height: 1,

		secondaryKey:    secondaryKey,
		secondaryLeft:   nil,
		secondaryRight:  nil,
		secondaryParent: secondaryParent,
		secondaryHeight: 1,

		value: *newVal,
	}

	if parent == nil {
		tree.root = node
		return "ok"
	}

	if key < parent.key {
		parent.left = node
	} else {
		parent.right = node
	}

	return tree.changeHeights(node.parent)
}

func (tree *AVL) update(key string, value string) string {
	node, _ := tree.search(key)

	if node == nil {
		return "does not exist"
	}

	newVal, ok := task4.Pool.Insert(value)

	if ok != "ok" {
		return "error"
	}

	node.value = *newVal
	return "ok"
}

func (tree *AVL) remove(key string) string {
	node, _ := tree.search(key)
	if node == nil {
		return "does not exist"
	}

	response := tree.secondaryRemove(key, node)

	if response != "ok" {
		return response
	}

	parent := node.parent
	if node.left == nil && node.right == nil {
		// Node with no children
		if parent == nil {
			tree.root = nil
			return "ok"
		} else {
			if parent.left == node {
				parent.left = nil
			} else {
				parent.right = nil
			}

			return tree.changeHeights(parent)
		}
	} else if node.left == nil {
		// Node with only right child
		if parent == nil {
			tree.root = node.right
			node.right.parent = nil
			return "ok"
		} else {
			if parent.left == node {
				parent.left = node.right
			} else {
				parent.right = node.right
			}
			node.right.parent = parent

			parent.height = max(tree.height(node.parent.left), tree.height(node.parent.right)) + 1
			return tree.changeHeights(parent)
		}
	} else if node.right == nil {
		// Node with only left child
		if parent == nil {
			tree.root = node.left
			node.left.parent = nil
			return "ok"
		} else {
			if parent.left == node {
				parent.left = node.left
			} else {
				parent.right = node.left
			}
			node.left.parent = parent

			parent.height = max(tree.height(parent.left), tree.height(parent.right)) + 1
			return tree.changeHeights(parent)
		}
	} else {
		// Node with two children
		successor, successorParent := tree.min(node.right)
		if successor == nil && successorParent == nil {
			return "not found children... somehow"
		}
		if successorParent == nil {
			successorParent = successor.parent
		}

		node.key = successor.key
		node.value = successor.value

		if successorParent.left == successor {
			if successorParent.left.right != nil {
				successorParent.left = successorParent.left.right
			} else {
				successorParent.left = nil
			}
		} else {
			successorParent.right = nil
		}

		successorParent.height = max(tree.height(successorParent.left), tree.height(successorParent.right)) + 1
		return tree.changeHeights(successorParent)
	}
}

func (tree *AVL) min(node *nodeAVL) (*nodeAVL, *nodeAVL) {
	if node == nil {
		return nil, nil
	}
	current := node
	for current.left != nil {
		current = current.left
	}

	return current, current.parent
}

func (tree *AVL) print() {
	fmt.Println("\nAVL Tree:\n")
	tree.root.printHelper()
}

func (node *nodeAVL) printHelper() {
	if node == nil {
		return
	}

	if node.parent != nil {
		fmt.Printf("Node: %v/%v:%v, h=%v. (parent: %v, ", node.key, node.secondaryKey, node.value, node.height, node.parent.key)
	} else {
		fmt.Printf("Node: %v/%v:%v h=%v. (parent: nil, ", node.key, node.secondaryKey, node.value, node.height)
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

func (tree *AVL) get(key string) (string, bool) {
	node, _ := tree.search(key)
	if node == nil {
		return "", false
	}
	val, ok := node.value.String()
	if !ok {
		return "", false
	}
	return val, true
}

func (tree *AVL) getRange(leftBound string, rightBound string) (*map[string]string, string) {
	result := make(map[string]string)
	return &result, tree.root.getRange(leftBound, rightBound, &result)
}

func (node *nodeAVL) getRange(leftBound string, rightBound string, result *map[string]string) (ret string) {
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