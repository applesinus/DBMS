package task6

import (
	"DBMS/task4"
	"fmt"
)

type Btree struct {
	root          *nodeBtree
	secondaryRoot *nodeBtree
	t             int
}

// len(children) == len(keys) + 1
type nodeBtree struct {
	isLeaf   bool
	keys     []string
	children []*nodeBtree
	values   []*task4.TrieWord
}

func (tree Btree) Copy() tree {
	newTree := &Btree{root: nil, secondaryRoot: nil, t: tree.t}
	newTree.root = tree.root.copy()
	newTree.secondaryRoot = tree.secondaryRoot.copy()
	return newTree
}

func (node *nodeBtree) copy() *nodeBtree {
	newNode := &nodeBtree{
		isLeaf:   node.isLeaf,
		keys:     make([]string, len(node.keys)),
		children: make([]*nodeBtree, len(node.children)),
		values:   make([]*task4.TrieWord, len(node.values)),
	}

	copy(newNode.keys, node.keys)
	copy(newNode.values, node.values)

	if node.isLeaf {
		return newNode
	}

	for i := 0; i < len(node.children); i++ {
		newNode.children[i] = node.children[i].copy()
	}

	return newNode
}

func (tree *Btree) search(key string) (*nodeBtree, int) {
	if tree.root == nil {
		return nil, -1
	}
	return tree.root.search(key)
}

func (tree *Btree) searchBySecondaryKey(secondaryKey string) (*nodeBtree, int) {
	if tree.secondaryRoot == nil {
		return nil, -1
	}
	return tree.secondaryRoot.search(secondaryKey)
}

func (node *nodeBtree) search(key string) (*nodeBtree, int) {
	if node == nil || len(node.keys) == 0 {
		return nil, -1
	}

	for i := 0; i < len(node.keys); i++ {
		if key < node.keys[i] {
			if node.isLeaf {
				return node, -1
			} else {
				return node.children[i].search(key)
			}
		} else if key == node.keys[i] {
			return node, i
		}
	}

	if !node.isLeaf {
		return node.children[len(node.keys)].search(key)
	}

	return nil, -1
}

func (tree *Btree) searchPLR(key string) (*nodeBtree, int, *nodeBtree, *nodeBtree, int, *nodeBtree, int) {
	if tree.root == nil {
		return nil, -1, nil, nil, -1, nil, -1
	}
	return tree.root.searchPLR(key, nil, nil, -1, nil, -1)
}

func (tree *Btree) searchPLRbySecondary(secondaryKey string) (*nodeBtree, int, *nodeBtree, *nodeBtree, int, *nodeBtree, int) {
	if tree.secondaryRoot == nil {
		return nil, -1, nil, nil, -1, nil, -1
	}
	return tree.secondaryRoot.searchPLR(secondaryKey, nil, nil, -1, nil, -1)
}

func (node *nodeBtree) searchPLR(key string, parent *nodeBtree, left *nodeBtree, leftIndex int, right *nodeBtree, rightIndex int) (*nodeBtree, int, *nodeBtree, *nodeBtree, int, *nodeBtree, int) {
	if len(node.keys) == 0 {
		return nil, -1, nil, nil, -1, nil, -1
	}

	for i := 0; i < len(node.keys); i++ {
		if key < node.keys[i] {
			if node.isLeaf {
				return nil, -1, nil, nil, -1, nil, -1
			} else {
				left = node.children[i]
				right = node.children[i+1]
				return node.children[i].searchPLR(key, node, left, i-1, right, i)
			}
		} else if key == node.keys[i] {
			return node, i, parent, left, leftIndex, right, rightIndex
		}
	}

	if key > node.keys[len(node.keys)-1] {
		if node.isLeaf {
			return nil, -1, nil, nil, -1, nil, -1
		} else {
			left = node.children[len(node.keys)]
			right = nil
			return node.children[len(node.keys)].searchPLR(key, node, left, len(node.keys), right, -1)
		}
	}

	return nil, -1, nil, nil, -1, nil, -1
}

func (tree *Btree) update(key string, value string) string {

	if tree.root == nil {
		return "does not exist"
	}

	node, index := tree.search(key)
	if node == nil || index == -1 {
		return "does not exist"
	}

	newVal, ok := task4.Pool.Insert(value)
	if ok != "ok" {
		return "error"
	}
	node.values[index] = newVal
	return "ok"
}

func (tree *Btree) set(key string, secondaryKey string, value string) string {
	if tree.root == nil {
		tree.root = &nodeBtree{
			isLeaf:   true,
			keys:     make([]string, 0),
			children: make([]*nodeBtree, 0),
			values:   make([]*task4.TrieWord, 0),
		}
		tree.root.keys = append(tree.root.keys, key)

		tree.secondaryRoot = &nodeBtree{
			isLeaf:   true,
			keys:     make([]string, 0),
			children: make([]*nodeBtree, 0),
			values:   make([]*task4.TrieWord, 0),
		}
		tree.secondaryRoot.keys = append(tree.secondaryRoot.keys, secondaryKey)

		newVal, ok := task4.Pool.Insert(value)
		if ok != "ok" {
			return "error"
		}
		tree.root.values = append(tree.root.values, newVal)
		tree.secondaryRoot.values = append(tree.secondaryRoot.values, newVal)
		return "ok"
	}

	node, index := tree.search(key)
	if node != nil && index != -1 {
		return "exist"
	}

	root := tree.root
	if len(root.keys) == 2*tree.t-1 {
		newRoot := &nodeBtree{
			isLeaf:   false,
			keys:     make([]string, 0),
			values:   make([]*task4.TrieWord, 0),
			children: make([]*nodeBtree, 0),
		}
		newRoot.children = append(newRoot.children, root)
		tree.splitChild(newRoot, 0, root)
		tree.insertNonFull(newRoot, key, value)
		tree.root = newRoot
	} else {
		tree.insertNonFull(root, key, value)
	}

	node, index = tree.searchBySecondaryKey(secondaryKey)
	if node != nil && index != -1 {
		return "exist"
	}

	root = tree.secondaryRoot
	if len(root.keys) == 2*tree.t-1 {
		newRoot := &nodeBtree{
			isLeaf:   false,
			keys:     make([]string, 0),
			values:   make([]*task4.TrieWord, 0),
			children: make([]*nodeBtree, 0),
		}
		newRoot.children = append(newRoot.children, root)
		tree.splitChild(newRoot, 0, root)
		tree.insertNonFull(newRoot, secondaryKey, value)
		tree.secondaryRoot = newRoot
	} else {
		tree.insertNonFull(root, secondaryKey, value)
	}

	return "ok"
}

func (tree *Btree) insertNonFull(node *nodeBtree, key string, value string) {
	i := len(node.keys) - 1
	if node.isLeaf {
		node.keys = append(node.keys, "")
		node.values = append(node.values, nil)
		for i >= 0 && key < node.keys[i] {
			node.keys[i+1] = node.keys[i]
			node.values[i+1] = node.values[i]
			i--
		}
		node.keys[i+1] = key
		newVal, ok := task4.Pool.Insert(value)
		if ok != "ok" {
			return
		}
		node.values[i+1] = newVal
	} else {
		for i >= 0 && key < node.keys[i] {
			i--
		}
		i++
		if len(node.children[i].keys) == 2*tree.t-1 {
			tree.splitChild(node, i, node.children[i])
			if key > node.keys[i] {
				i++
			}
		}
		tree.insertNonFull(node.children[i], key, value)
	}
}

func (tree *Btree) splitChild(parent *nodeBtree, i int, fullNode *nodeBtree) {
	t := tree.t
	newNode := &nodeBtree{
		isLeaf:   fullNode.isLeaf,
		keys:     make([]string, 0),
		children: make([]*nodeBtree, 0),
		values:   make([]*task4.TrieWord, 0),
	}
	parent.children = append(parent.children[:i+1], append([]*nodeBtree{newNode}, parent.children[i+1:]...)...)
	parent.keys = append(parent.keys[:i], append([]string{fullNode.keys[t-1]}, parent.keys[i:]...)...)
	parent.values = append(parent.values[:i], append([]*task4.TrieWord{fullNode.values[t-1]}, parent.values[i:]...)...)
	newNode.keys = append(newNode.keys, fullNode.keys[t:]...)
	newNode.values = append(newNode.values, fullNode.values[t:]...)
	fullNode.keys = fullNode.keys[:t-1]
	fullNode.values = fullNode.values[:t-1]
	if !fullNode.isLeaf {
		newNode.children = append(newNode.children, fullNode.children[t:]...)
		fullNode.children = fullNode.children[:t]
	}
}

func (tree *Btree) remove(key string) string {
	if tree.root == nil {
		return "does not exist"
	}
	node, index, parent, left, leftIndex, right, rightIndex := tree.searchPLR(key)
	if node == nil || index == -1 {
		return "does not exist"
	}

	response := tree.removeHelper(key, node, index, parent, tree.root, left, leftIndex, right, rightIndex)
	if response != "ok" {
		return response
	}

	return tree.removeHelper(key, node, index, parent, tree.secondaryRoot, left, leftIndex, right, rightIndex)
}

func (tree *Btree) removeHelper(key string, node *nodeBtree, index int, root *nodeBtree, parent *nodeBtree, left *nodeBtree, leftIndex int, right *nodeBtree, rightIndex int) string {
	if node.isLeaf {
		if len(node.keys) > tree.t {
			// case 1
			node.keys = append(node.keys[:index], node.keys[index+1:]...)
			node.values = append(node.values[:index], node.values[index+1:]...)
			return "ok"
		} else {
			if right != nil && len(right.keys) >= tree.t {
				// delete key
				node.keys = append(node.keys[:index], node.keys[index+1:]...)
				node.values = append(node.values[:index], node.values[index+1:]...)
				// move separator to the node
				node.keys = append(node.keys, parent.keys[leftIndex+1])
				node.values = append(node.values, parent.values[leftIndex+1])
				// move 1st right key to the separator position
				parent.keys[rightIndex] = right.keys[0]
				parent.values[rightIndex] = right.values[0]
				// remove 1st right key
				right.keys = right.keys[1:]
				right.values = right.values[1:]

				return "ok"
			} else if left != nil && len(left.keys) >= tree.t {
				// delete key
				node.keys = append(node.keys[:index], node.keys[index+1:]...)
				node.values = append(node.values[:index], node.values[index+1:]...)
				// move separator to the node
				node.keys = append(node.keys, parent.keys[leftIndex])
				node.values = append(node.values, parent.values[leftIndex])
				// move last left key to the separator position
				parent.keys[leftIndex] = left.keys[len(left.keys)-1]
				parent.values[leftIndex] = left.values[len(left.keys)-1]
				// remove last left key
				left.keys = left.keys[:len(left.keys)-1]
				left.values = left.values[:len(left.keys)-1]

				return "ok"
			} else if right != nil {
				// delete key
				node.keys = append(node.keys[:index], node.keys[index+1:]...)
				node.values = append(node.values[:index], node.values[index+1:]...)
				// move separator to the node
				fmt.Printf("right separator for %v: %v\n", key, parent.values[rightIndex])
				node.keys = append(node.keys, parent.keys[rightIndex])
				node.values = append(node.values, parent.values[rightIndex])
				// merge with right
				node.keys = append(node.keys, right.keys...)
				node.values = append(node.values, right.values...)
				// check if parent has no more keys
				if len(parent.keys) == 1 {
					if parent == root {
						root = node
					}
					*parent = *node
				} else {
					parent.keys = append(parent.keys[:leftIndex+1], parent.keys[leftIndex+2:]...)
					parent.values = append(parent.values[:leftIndex+1], parent.values[leftIndex+2:]...)
					parent.children = append(parent.children[:leftIndex+1], parent.children[leftIndex+2:]...)
				}
				return "ok"
			} else if left != nil {
				// delete key
				node.keys = append(node.keys[:index], node.keys[index+1:]...)
				node.values = append(node.values[:index], node.values[index+1:]...)
				// move separator to the node
				node.keys = append(node.keys, parent.keys[leftIndex])
				node.values = append(node.values, parent.values[leftIndex])
				// merge with left
				node.keys = append(node.keys, left.keys...)
				node.values = append(node.values, left.values...)
				// check if parent has no more keys
				if len(parent.keys) == 1 {
					if parent == root {
						root = node
					}
					*parent = *node
				} else {
					parent.keys = append(parent.keys[:leftIndex], parent.keys[leftIndex+1:]...)
					parent.values = append(parent.values[:leftIndex], parent.values[leftIndex+1:]...)
					parent.children = append(parent.children[:leftIndex], parent.children[leftIndex+1:]...)
				}
				return "ok"
			} else if node == root {
				node.keys = append(node.keys[:index], node.keys[index+1:]...)
				node.values = append(node.values[:index], node.values[index+1:]...)
				return "ok"
			}
		}
	} else {
		if left != nil && len(left.keys) >= tree.t {
			// case 2a
			node.keys[index] = left.keys[len(left.keys)-1]
			node.values[index] = left.values[len(left.keys)-1]
			left.keys = left.keys[:len(left.keys)-1]
			return "ok"
		} else if right != nil && len(right.keys) >= tree.t {
			// case 2b
			node.keys[index] = right.keys[0]
			node.values[index] = right.values[0]
			right.keys = right.keys[1:]
			return "ok"
		} else if left != nil && right != nil {
			// case 2c
			left.keys = append(left.keys, right.keys...)
			left.values = append(left.values, right.values...)
			left.children = append(left.children, right.children...)

			parent.keys = append(parent.keys[:rightIndex], parent.keys[rightIndex+1:]...)
			parent.values = append(parent.values[:rightIndex], parent.values[rightIndex+1:]...)
			parent.children = append(parent.children[:rightIndex], parent.children[rightIndex+1:]...)
			return "ok"
		} else {
			if node == root {
				newKey := node.children[index+1].keys[0]
				newValue := node.children[index+1].values[0]
				res := tree.remove(newKey)
				node.keys[index] = newKey
				node.values[index] = newValue
				return res
			}
		}
	}
	return "error"
}

func (tree *Btree) print() {
	fmt.Println("\nB Tree:\n")
	fmt.Println("Primary keys sorted:")
	tree.root.printHelper()
	fmt.Println("Secondary keys sorted:")
	tree.secondaryRoot.printHelper()
}

func (node *nodeBtree) printHelper() {
	if node == nil {
		return
	}

	fmt.Printf("Node (%v), keys: %v, values: %v, children: %v\n", node.isLeaf, node.keys, node.values, node.children)
	for i := 0; i < len(node.children); i++ {
		node.children[i].printHelper()
	}
}

func (tree *Btree) getter(key string, root *nodeBtree) (string, bool) {
	node, index := root.search(key)
	if node == nil || index == -1 {
		return "", false
	}

	val, ok := node.values[index].String()

	if !ok {
		return "", false
	}
	return val, true
}

func (tree *Btree) get(key string) (string, bool) {
	return tree.getter(key, tree.root)
}

func (tree *Btree) getBySecondaryKey(secondaryKey string) (string, bool) {
	return tree.getter(secondaryKey, tree.secondaryRoot)
}

func (tree *Btree) getRange(leftBound string, rightBound string) (*map[string]string, string) {
	result := make(map[string]string)
	return &result, tree.root.getRange(leftBound, rightBound, &result)
}

func (tree *Btree) getRangeBySecondaryKey(leftBound string, rightBound string) (*map[string]string, string) {
	result := make(map[string]string)
	return &result, tree.secondaryRoot.getRange(leftBound, rightBound, &result)
}

func (node *nodeBtree) getRange(leftBound string, rightBound string, result *map[string]string) (ret string) {
	defer func() {
		if ret != "ok" {
			ret = "error"
		}
	}()

	ret = "start"

	if !node.isLeaf {
		for i := 0; i < len(node.keys); i++ {
			if node.keys[i] >= leftBound && node.keys[i] <= rightBound {
				ret = node.children[i].getRange(leftBound, rightBound, result)
				if ret != "ok" {
					return
				}
			}
		}
		if node.keys[len(node.keys)-1] >= leftBound && node.keys[len(node.keys)-1] <= rightBound {
			ret = node.children[len(node.children)-1].getRange(leftBound, rightBound, result)
			if ret != "ok" {
				return
			}
		}
	}

	for i := 0; i < len(node.keys); i++ {
		if node.keys[i] >= leftBound && node.keys[i] <= rightBound {
			val, ok := node.values[i].String()

			if !ok {
				return "error"
			}
			(*result)[node.keys[i]] = val
		}
	}

	return "ok"
}
