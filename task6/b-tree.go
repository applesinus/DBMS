package task6

import "fmt"

type Btree struct {
	root *nodeBtree
	t    int
}

// len(children) == len(keys) + 1
type nodeBtree struct {
	isLeaf   bool
	keys     []string
	children []*nodeBtree
	values   []interface{}
}

func (tree *Btree) search(key string) (*nodeBtree, int) {
	if tree.root == nil {
		return nil, -1
	}
	return tree.root.search(key)
}

func (node *nodeBtree) search(key string) (*nodeBtree, int) {
	if len(node.keys) == 0 {
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

func (tree *Btree) searchPLR(key string) (*nodeBtree, int, *nodeBtree, *nodeBtree, *nodeBtree) {
	if tree.root == nil {
		return nil, -1, nil, nil, nil
	}
	return tree.root.searchPLR(key, nil, nil, nil)
}

func (node *nodeBtree) searchPLR(key string, parent *nodeBtree, left *nodeBtree, right *nodeBtree) (*nodeBtree, int, *nodeBtree, *nodeBtree, *nodeBtree) {
	if len(node.keys) == 0 {
		return nil, -1, nil, nil, nil
	}

	for i := 0; i < len(node.keys); i++ {
		if key < node.keys[i] {
			if node.isLeaf {
				return node, -1, nil, nil, nil
			} else {
				if i == 0 {
					left = nil
					right = node.children[i+1]
				} else {
					left = node.children[i-1]
					right = node.children[i+1]
				}
				return node.children[i].searchPLR(key, node, left, right)
			}
		} else if key == node.keys[i] {
			return node, i, parent, left, right
		}
	}

	if !node.isLeaf {
		left = node.children[len(node.children)-1]
		right = nil
		return node.children[len(node.keys)].searchPLR(key, node, left, right)
	}

	return nil, -1, nil, nil, nil
}

func (tree *Btree) update(key string, value interface{}) string {

	if tree.root == nil {
		return "does not exist"
	}

	node, index := tree.search(key)
	if node == nil || index == -1 {
		return "does not exist"
	}

	node.values[index] = value
	return "ok"
}

func (tree *Btree) insert(key string, value interface{}) string {
	if tree.root == nil {
		tree.root = &nodeBtree{
			isLeaf:   true,
			keys:     make([]string, 0),
			children: make([]*nodeBtree, 0),
			values:   make([]interface{}, 0),
		}
		tree.root.keys = append(tree.root.keys, key)
		tree.root.values = append(tree.root.values, value)
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
			values:   make([]interface{}, 0),
			children: make([]*nodeBtree, 0),
		}
		newRoot.children = append(newRoot.children, root)
		tree.splitChild(newRoot, 0, root)
		tree.insertNonFull(newRoot, key, value)
		tree.root = newRoot
	} else {
		tree.insertNonFull(root, key, value)
	}

	return "ok"
}

func (tree *Btree) insertNonFull(node *nodeBtree, key string, value interface{}) {
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
		node.values[i+1] = value
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
		values:   make([]interface{}, 0),
	}
	parent.children = append(parent.children[:i+1], append([]*nodeBtree{newNode}, parent.children[i+1:]...)...)
	parent.keys = append(parent.keys[:i], append([]string{fullNode.keys[t-1]}, parent.keys[i:]...)...)
	parent.values = append(parent.values[:i], append([]interface{}{fullNode.values[t-1]}, parent.values[i:]...)...)
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
	node, index, parent, left, right := tree.searchPLR(key)
	if node == nil || index == -1 {
		return "does not exist"
	}

	if node.isLeaf {
		if len(node.keys) > tree.t-1 {
			node.keys = append(node.keys[:index], node.keys[index+1:]...)
			node.values = append(node.values[:index], node.values[index+1:]...)
		} else {
			
		}
	}
}

func (tree *Btree) print() {
	fmt.Println("\nB Tree:\n")
	tree.root.printHelper()
}

func (node *nodeBtree) printHelper() {
	if node == nil {
		return
	}
	node.children[0].printHelper()

	for i := 0; i < len(node.keys); i++ {
		if node.isLeaf {
			fmt.Printf("%v: %v\n", node.keys[i], node.values[i])
		}

		node.children[i+1].printHelper()
	}
}
