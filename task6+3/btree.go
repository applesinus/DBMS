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
	t        int
	n        int
	keys     []string
	altKeys  []string
	children []*nodeBtree
	values   []*task4.TrieWord
}

func (tree *Btree) set(key string, secondaryKey string, value string) string {
	word, res := task4.Pool.Insert(value)
	if res != "ok" {
		return res
	}

	_, idx := tree.root.search(key)
	if idx > -1 {
		return "key already exists"
	}
	_, secIdx := tree.secondaryRoot.search(secondaryKey)
	if secIdx > -1 {
		return "secondary key already exists"
	}

	tree.insert(key, secondaryKey, word)

	return "ok"
}

func (tree *Btree) update(key string, value string) string {
	word, res := task4.Pool.Insert(value)
	if res != "ok" {
		return res
	}

	node, idx := tree.root.search(key)
	if idx == -1 {
		return "key not found"
	}
	secondaryKey := node.altKeys[idx]
	secNode, secIdx := tree.secondaryRoot.search(secondaryKey)
	if secIdx == -1 {
		return "secondary key not found"
	}

	node.values[idx] = word
	secNode.values[secIdx] = word

	return "ok"
}

func (tree *Btree) get(key string) (string, bool) {
	if tree.root == nil || tree.secondaryRoot == nil {
		return "", false
	}

	node, i := tree.root.search(key)
	if i == -1 {
		return "", false
	}
	val, ok := node.values[i].String()
	if !ok {
		return "", false
	}

	return val, true
}

func (tree *Btree) getBySecondaryKey(secondaryKey string) (string, bool) {
	if tree.root == nil || tree.secondaryRoot == nil {
		return "", false
	}

	node, i := tree.secondaryRoot.search(secondaryKey)
	if i == -1 {
		return "", false
	}
	val, ok := node.values[i].String()
	if !ok {
		return "", false
	}

	return val, true
}

func (tree *Btree) print() {
	fmt.Println("\nB Tree:")
	fmt.Println("\nOrdered by keys:")
	tree.root.printHelper()
	fmt.Println("\nOrdered by secondary keys:")
	tree.secondaryRoot.printHelper()
}

func (node *nodeBtree) printHelper() {
	if node == nil {
		return
	}

	fmt.Printf("Node (%v), keys: %v, altKeys: %v, values: %v, children: %v\n", node.isLeaf, node.keys, node.altKeys, node.values, node.children)
	for i := 0; i < len(node.children); i++ {
		node.children[i].printHelper()
	}
}

func (tree *Btree) getAll() (*[]string, *[]string, *[]string, string) {
	keys := make([]string, 0)
	secondaryKeys := make([]string, 0)
	values := make([]string, 0)

	return &keys, &secondaryKeys, &values, tree.root.getAll(&keys, &secondaryKeys, &values)
}

func (node *nodeBtree) getAll(keys, secondaryKeys, values *[]string) string {
	if node == nil {
		return "ok"
	}

	for i := 0; i < node.n; i++ {
		*keys = append(*keys, node.keys[i])
		*secondaryKeys = append(*secondaryKeys, node.altKeys[i])
		val, ok := node.values[i].String()
		if !ok {
			return "error"
		}
		*values = append(*values, val)
	}
	for i := 0; i < node.n+1; i++ {
		res := node.children[i].getAll(keys, secondaryKeys, values)
		if res != "ok" {
			return res
		}
	}
	return "ok"
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
		if node.keys[0] > leftBound {
			ret = node.children[0].getRange(leftBound, rightBound, result)
			if ret != "ok" {
				return
			}
		}
		for i := 1; i < len(node.keys); i++ {
			if node.keys[i] >= leftBound && node.keys[i] <= rightBound {
				ret = node.children[i].getRange(leftBound, rightBound, result)
				if ret != "ok" {
					return
				}
			}
		}
		if node.keys[len(node.keys)-1] < rightBound {
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

func newNodeBtree(t int, isLeaf bool) *nodeBtree {
	return &nodeBtree{
		isLeaf:   isLeaf,
		keys:     make([]string, 2*t-1),
		altKeys:  make([]string, 2*t-1),
		children: make([]*nodeBtree, 2*t),
		values:   make([]*task4.TrieWord, 2*t-1),
	}
}

func (tree *Btree) insert(key, altKey string, value *task4.TrieWord) {
	if tree.root == nil {
		tree.root = newNodeBtree(tree.t, true)
		tree.secondaryRoot = newNodeBtree(tree.t, true)
		tree.root.keys[0] = key
		tree.root.altKeys[0] = altKey
		tree.root.values[0] = value
		tree.root.n = 1
		tree.secondaryRoot.keys[0] = altKey
		tree.secondaryRoot.altKeys[0] = key
		tree.secondaryRoot.values[0] = value
		tree.secondaryRoot.n = 1
	} else {
		if tree.root.n == 2*tree.t-1 {
			s := newNodeBtree(tree.t, false)
			s.children[0] = tree.root
			s.splitChild(0, tree.root, tree.t)
			i := 0
			if s.keys[0] < key {
				i++
			}
			s.children[i].insertNonFull(key, altKey, value, tree.t)
			tree.root = s
		} else {
			tree.root.insertNonFull(key, altKey, value, tree.t)
		}

		if tree.secondaryRoot.n == 2*tree.t-1 {
			s := newNodeBtree(tree.t, false)
			s.children[0] = tree.secondaryRoot
			s.splitChild(0, tree.secondaryRoot, tree.t)
			i := 0
			if s.keys[0] < altKey {
				i++
			}
			s.children[i].insertNonFull(altKey, key, value, tree.t)
			tree.secondaryRoot = s
		} else {
			tree.secondaryRoot.insertNonFull(altKey, key, value, tree.t)
		}
	}
}

func (node *nodeBtree) insertNonFull(key, altKey string, value *task4.TrieWord, t int) {
	i := node.n - 1

	if node.isLeaf {
		for i >= 0 && node.keys[i] > key {
			node.keys[i+1] = node.keys[i]
			node.altKeys[i+1] = node.altKeys[i]
			node.values[i+1] = node.values[i]
			i--
		}
		node.keys[i+1] = key
		node.altKeys[i+1] = altKey
		node.values[i+1] = value
		node.n++
	} else {
		for i >= 0 && node.keys[i] > key {
			i--
		}
		if node.children[i+1].n == 2*t-1 {
			node.splitChild(i+1, node.children[i+1], t)
			if node.keys[i+1] < key {
				i++
			}
		}
		node.children[i+1].insertNonFull(key, altKey, value, t)
	}
}

func (node *nodeBtree) splitChild(i int, y *nodeBtree, t int) {
	z := newNodeBtree(t, y.isLeaf)
	z.n = t - 1

	for j := 0; j < t-1; j++ {
		z.keys[j] = y.keys[j+t]
		z.altKeys[j] = y.altKeys[j+t]
		z.values[j] = y.values[j+t]
	}

	if !y.isLeaf {
		for j := 0; j < t; j++ {
			z.children[j] = y.children[j+t]
		}
	}

	y.n = t - 1

	for j := node.n; j >= i+1; j-- {
		node.children[j+1] = node.children[j]
	}
	node.children[i+1] = z

	for j := node.n - 1; j >= i; j-- {
		node.keys[j+1] = node.keys[j]
		node.altKeys[j+1] = node.altKeys[j]
		node.values[j+1] = node.values[j]
	}
	node.keys[i] = y.keys[t-1]
	node.altKeys[i] = y.altKeys[t-1]
	node.values[i] = y.values[t-1]

	node.n++
}

func (node *nodeBtree) search(key string) (*nodeBtree, int) {
	if node == nil {
		return nil, -1
	}
	i := 0
	for i < node.n && key > node.keys[i] {
		i++
	}

	if i < node.n && node.keys[i] == key {
		return node, i
	}

	if node.isLeaf {
		return node, -1
	}

	return node.children[i].search(key)
}

func (tree *Btree) remove(key string) string {
	if tree.root == nil {
		return "empty tree"
	}

	node, idx := tree.root.search(key)
	if idx == -1 {
		return "key not found"
	}
	altKey := node.altKeys[idx]

	tree.root.remove(key, tree.t)

	if tree.root.n == 0 {
		if tree.root.isLeaf {
			tree.root = nil
		} else {
			tree.root = tree.root.children[0]
		}
	}

	tree.secondaryRoot.remove(altKey, tree.t)

	if tree.secondaryRoot.n == 0 {
		if tree.secondaryRoot.isLeaf {
			tree.secondaryRoot = nil
		} else {
			tree.secondaryRoot = tree.secondaryRoot.children[0]
		}
	}

	return "ok"
}

func (node *nodeBtree) remove(key string, t int) {
	idx := node.findKey(key)

	if idx < node.n && node.keys[idx] == key {
		if node.isLeaf {
			node.removeFromLeaf(idx)
		} else {
			node.removeFromNonLeaf(idx, t)
		}
	} else {
		if node.isLeaf {
			fmt.Printf("The key %s does not exist in the tree\n", key)
			return
		}

		flag := (idx == node.n)

		if node.children[idx].n < t {
			node.fill(idx, t)
		}

		if flag && idx > node.n {
			node.children[idx-1].remove(key, t)
		} else {
			node.children[idx].remove(key, t)
		}
	}
}

func (node *nodeBtree) findKey(key string) int {
	idx := 0
	for idx < node.n && node.keys[idx] < key {
		idx++
	}
	return idx
}

func (node *nodeBtree) removeFromLeaf(idx int) {
	for i := idx + 1; i < node.n; i++ {
		node.keys[i-1] = node.keys[i]
		node.altKeys[i-1] = node.altKeys[i]
		node.values[i-1] = node.values[i]
	}
	node.n--
}

func (node *nodeBtree) removeFromNonLeaf(idx int, t int) {
	key := node.keys[idx]

	if node.children[idx].n >= t {
		pred := node.getPred(idx)
		node.keys[idx] = pred
		node.children[idx].remove(pred, t)
	} else if node.children[idx+1].n >= t {
		succ := node.getSucc(idx)
		node.keys[idx] = succ
		node.children[idx+1].remove(succ, t)
	} else {
		node.merge(idx, t)
		node.children[idx].remove(key, t)
	}
}

func (node *nodeBtree) getPred(idx int) string {
	cur := node.children[idx]
	for !cur.isLeaf {
		cur = cur.children[cur.n]
	}
	return cur.keys[cur.n-1]
}

func (node *nodeBtree) getSucc(idx int) string {
	cur := node.children[idx+1]
	for !cur.isLeaf {
		cur = cur.children[0]
	}
	return cur.keys[0]
}

func (node *nodeBtree) fill(idx, t int) {
	if idx != 0 && node.children[idx-1].n >= t {
		node.borrowFromPrev(idx)
	} else if idx != node.n && node.children[idx+1].n >= t {
		node.borrowFromNext(idx)
	} else {
		if idx != node.n {
			node.merge(idx, t)
		} else {
			node.merge(idx-1, t)
		}
	}
}

func (node *nodeBtree) borrowFromPrev(idx int) {
	child := node.children[idx]
	sibling := node.children[idx-1]

	for i := child.n - 1; i >= 0; i-- {
		child.keys[i+1] = child.keys[i]
		child.altKeys[i+1] = child.altKeys[i]
		child.values[i+1] = child.values[i]
	}

	if !child.isLeaf {
		for i := child.n; i >= 0; i-- {
			child.children[i+1] = child.children[i]
		}
	}

	child.keys[0] = node.keys[idx-1]
	child.altKeys[0] = node.altKeys[idx-1]
	child.values[0] = node.values[idx-1]

	if !child.isLeaf {
		child.children[0] = sibling.children[sibling.n]
	}

	node.keys[idx-1] = sibling.keys[sibling.n-1]
	node.altKeys[idx-1] = sibling.altKeys[sibling.n-1]
	node.values[idx-1] = sibling.values[sibling.n-1]

	child.n++
	sibling.n--
}

func (node *nodeBtree) borrowFromNext(idx int) {
	child := node.children[idx]
	sibling := node.children[idx+1]

	child.keys[child.n] = node.keys[idx]
	child.altKeys[child.n] = node.altKeys[idx]
	child.values[child.n] = node.values[idx]

	if !child.isLeaf {
		child.children[child.n+1] = sibling.children[0]
	}

	node.keys[idx] = sibling.keys[0]
	node.altKeys[idx] = sibling.altKeys[0]
	node.values[idx] = sibling.values[0]

	for i := 1; i < sibling.n; i++ {
		sibling.keys[i-1] = sibling.keys[i]
		sibling.altKeys[i-1] = sibling.altKeys[i]
		sibling.values[i-1] = sibling.values[i]
	}

	if !sibling.isLeaf {
		for i := 1; i <= sibling.n; i++ {
			sibling.children[i-1] = sibling.children[i]
		}
	}

	child.n++
	sibling.n--
}

func (node *nodeBtree) merge(idx, t int) {
	child := node.children[idx]
	sibling := node.children[idx+1]

	child.keys[t-1] = node.keys[idx]
	child.altKeys[t-1] = node.altKeys[idx]
	child.values[t-1] = node.values[idx]

	for i := 0; i < sibling.n; i++ {
		child.keys[i+t] = sibling.keys[i]
		child.altKeys[i+t] = sibling.altKeys[i]
		child.values[i+t] = sibling.values[i]
	}

	if !child.isLeaf {
		for i := 0; i <= sibling.n; i++ {
			child.children[i+t] = sibling.children[i]
		}
	}

	for i := idx + 1; i < node.n; i++ {
		node.keys[i-1] = node.keys[i]
		node.altKeys[i-1] = node.altKeys[i]
		node.values[i-1] = node.values[i]
	}

	for i := idx + 2; i <= node.n; i++ {
		node.children[i-1] = node.children[i]
	}

	child.n += sibling.n + 1
	node.n--

	sibling = nil
}

func (tree Btree) Copy() tree {
	newTree := &Btree{root: nil, secondaryRoot: nil, t: tree.t}
	newTree.root = tree.root.copy()
	newTree.secondaryRoot = tree.secondaryRoot.copy()
	return newTree
}

func (node *nodeBtree) copy() *nodeBtree {
	if node == nil {
		return nil
	}
	newNode := &nodeBtree{
		isLeaf:   node.isLeaf,
		keys:     make([]string, len(node.keys)),
		altKeys:  make([]string, len(node.altKeys)),
		children: make([]*nodeBtree, len(node.children)),
		values:   make([]*task4.TrieWord, len(node.values)),
		n:        node.n,
		t:        node.t,
	}

	copy(newNode.keys, node.keys)
	copy(newNode.altKeys, node.altKeys)
	copy(newNode.values, node.values)

	if node.isLeaf {
		return newNode
	}

	for i := 0; i < len(node.children); i++ {
		newNode.children[i] = node.children[i].copy()
	}

	return newNode
}
