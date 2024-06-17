package task6

import (
	"fmt"
	"strconv"
)

type tree interface {
	print()
	set(key string, value string) string
	update(key string, value string) string
	get(key string) (string, bool)
	getRange(leftBound string, rightBound string) (*map[string]string, string)
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

	if len(variant) >= 5 && variant[:5] == "Btree" {
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
	_, ok := t.self.get(key)
	if ok {
		fmt.Printf("Key %s already exists\n", key)
		return "error"
	}

	return t.self.set(key, value)
}

func (t Tree) Update(key string, value string) string {
	_, ok := t.self.get(key)
	if !ok {
		fmt.Printf("Key %s does not exist\n", key)
		return "error"
	}

	return t.self.update(key, value)
}

func (t Tree) Get(key string) (string, string) {
	value, ok := t.self.get(key)
	if !ok {
		fmt.Printf("Key %s does not exist\n", key)
		return "", "error"
	}

	return fmt.Sprintf("%s", value), "ok"
}

func (t Tree) GetRange(leftBound string, rightBound string) (*map[string]string, string) {
	result, ok := t.self.getRange(leftBound, rightBound)

	if ok != "ok" {
		return nil, ok
	}

	ret := make(map[string]string)

	for key, value := range *result {
		ret[key] = fmt.Sprintf("%s", value)
	}

	return &ret, "ok"
}

func (t Tree) Delete(key string) string {
	return t.self.remove(key)
}
