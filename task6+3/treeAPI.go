package task6

import (
	"DBMS/interfaces"
	"fmt"
	"strconv"
)

type tree interface {
	print()
	set(key string, secondaryKey string, value string) string
	update(key string, value string) string
	get(key string) (string, bool)
	getBySecondaryKey(secondaryKey string) (string, bool)
	getRange(leftBound string, rightBound string) (*map[string]string, string)
	getRangeBySecondaryKey(leftBound string, rightBound string) (*map[string]string, string)
	getAll() (*[]string, *[]string, *[]string, string)
	remove(key string) string
	Copy() tree
}

type Tree struct {
	self    tree
	variant string
}

func (t Tree) Copy() interfaces.CollectionInterface {
	return Tree{self: t.self.Copy(), variant: t.variant}
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
			self:    &Btree{root: nil, secondaryRoot: nil, t: t},
			variant: variant,
		}
	}

	return nil
}

func (t Tree) Set(key string, secondaryKey string, value string) string {
	_, ok := t.self.get(key)
	if ok {
		fmt.Printf("Key %s already exists\n", key)
		return "error"
	}

	return t.self.set(key, secondaryKey, value)
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

func (t Tree) GetAll() (*[]string, *[]string, *[]string, string) {
	return t.self.getAll()
}

func (t Tree) GetBySecondaryKey(secondaryKey string) (string, string) {
	value, ok := t.self.getBySecondaryKey(secondaryKey)
	if !ok {
		fmt.Printf("Secondary key %s does not exist\n", secondaryKey)
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

func (t Tree) GetRangeBySecondaryKey(leftBound string, rightBound string) (*map[string]string, string) {
	result, ok := t.self.getRangeBySecondaryKey(leftBound, rightBound)

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
