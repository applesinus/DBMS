package task6

import (
	"fmt"
	"strconv"
)

type tree interface {
	print()
	insert(key string, value interface{}) string
	update(key string, value interface{}) string
	find(key string) (interface{}, bool)
	remove(key string) string
}

type Tree struct {
	self    tree
	variant string
}

func NewTree(variant string) *Tree {
	if variant == "AVL" {
		return &Tree{
			self:    &AVL{root: nil},
			variant: "AVL",
		}
	}

	if variant == "RB" {
		return &Tree{
			self:    &RB{root: nil},
			variant: "RB",
		}
	}

	if len(variant) > 5 && variant[:5] == "Btree" {
		t, err := strconv.Atoi(variant[5:])
		if err != nil {
			t = 2
		}
		return &Tree{
			self:    &Btree{root: nil, t: t},
			variant: variant,
		}
	}

	return nil
}

func (t Tree) Set(key string, value string) string {
	_, ok := t.self.find(key)
	if ok {
		fmt.Printf("Key %s already exists\n", key)
		return "error"
	}

	return t.self.insert(key, value)
}

func (t Tree) Update(key string, value string) string {
	_, ok := t.self.find(key)
	if !ok {
		fmt.Printf("Key %s does not exist\n", key)
		return "error"
	}

	return t.self.update(key, value)
}

func (t Tree) Get(key string) (string, string) {
	value, ok := t.self.find(key)
	if !ok {
		fmt.Printf("Key %s does not exist\n", key)
		return "", "error"
	}

	return fmt.Sprintf("%s", value), "ok"
}

func (t Tree) GetRange(leftBound string, rightBound string) (map[string]string, string) {
	// TODO

	return nil, "ok"
}

func (t Tree) Delete(key string) string {
	return t.self.remove(key)
}
