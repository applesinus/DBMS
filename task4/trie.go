package task4

import (
	"fmt"
	"strings"
)

type Trie struct {
	root *NodeTrie
}

type NodeTrie struct {
	isEnd bool
	value byte
	nexts []*NodeTrie
}

type TrieWord struct {
	chars []*NodeTrie
}

func (t *Trie) insert(word string) (*TrieWord, string) {
	if t.root == nil {
		t.root = &NodeTrie{}
	}
	return t.root.insert(word, &TrieWord{chars: make([]*NodeTrie, 0)})
}

func (node *NodeTrie) insert(word string, prev *TrieWord) (*TrieWord, string) {
	if len(word) == 0 {
		if node.isEnd {
			return nil, "already exist"
		} else {
			node.isEnd = true
			prev.chars = append(prev.chars, node)
			return prev, "ok"
		}
	}

	for _, next := range node.nexts {
		if next.value == word[0] {
			prev.chars = append(prev.chars, node)
			return next.insert(word[1:], prev)
		}
	}

	newNode := &NodeTrie{
		value: word[0],
		nexts: make([]*NodeTrie, 0),
		isEnd: false,
	}
	node.nexts = append(node.nexts, newNode)

	prev.chars = append(prev.chars, node)
	return newNode.insert(word[1:], prev)
}

func (trieWord *TrieWord) String() (res string, ok bool) {
	ok = true

	defer func(ok *bool) {
		*ok = !*ok
	}(&ok)

	var word strings.Builder
	for _, node := range trieWord.chars {
		word.Grow(1)
		word.WriteByte(node.value)
	}
	res = word.String()
	ok = false

	return
}

func (t *Trie) print() {
	if t.root == nil {
		return
	}
	t.root.print(t, make([]*NodeTrie, 0))
}

func (node *NodeTrie) print(tree *Trie, prev []*NodeTrie) {
	if node.isEnd {
		newprev := TrieWord{chars: append(prev, node)}
		word, ok := newprev.String()
		if !ok {
			fmt.Printf("error\n")
			return
		}
		fmt.Printf("%v\n", word)
	}
	for _, next := range node.nexts {
		newprev := append(prev, node)
		next.print(tree, newprev)
	}
}
