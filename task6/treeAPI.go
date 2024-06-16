package task6

import (
	"fmt"
	"strconv"
	"strings"
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

func CreateDB(dbType string) *Tree {
	words := strings.Fields(dbType)
	if words[0] == "btree" {
		switch len(words) {
		case 1:
			return &Tree{self: &Btree{root: nil, t: 2}, variant: "btree"}
		case 2:
			t, err := strconv.Atoi(words[1])
			if err != nil {
				return nil
			}
			return &Tree{self: &Btree{root: nil, t: t}, variant: "btree"}
		default:
			return nil
		}

	} else if dbType == "avl" {
		return &Tree{self: &AVL{root: nil}, variant: "avl"}
	} else if dbType == "rb" {
		return &Tree{self: &RB{root: nil}, variant: "rb"}
	}
	return nil
}

func (t Tree) CreatePool(settings map[string]string, name string) string {
	_, ok := t.self.find(name)
	if ok {
		fmt.Printf("Pool %s already exists\n", name)
		return "error"
	}

	var schema tree
	switch t.variant {
	case "btree":
		schema = &Btree{root: nil, t: t.self.(*Btree).t}
	case "avl":
		schema = &AVL{root: nil}
	case "rb":
		schema = &RB{root: nil}
	default:
		fmt.Printf("Unknown database type %s\n", t.variant)
		return "error"
	}

	return t.self.insert(name, schema)
}

func (t Tree) DeletePool(settings map[string]string, name string) string {
	_, ok := t.self.find(name)
	if !ok {
		fmt.Printf("Pool %s does not exist\n", name)
		return "error"
	}

	return t.self.remove(name)
}

func (t Tree) CreateSchema(settings map[string]string, name string, poolName string) string {
	pool, ok := t.self.find(poolName)
	if !ok {
		fmt.Printf("Pool %s does not exist\n", poolName)
		return "error"
	}
	_, ok = pool.(tree).find(name)
	if ok {
		fmt.Printf("Schema %s already exists in pool %s\n", name, poolName)
		return "error"
	}

	var schema tree
	switch t.variant {
	case "btree":
		schema = &Btree{root: nil, t: pool.(*Btree).t}
	case "avl":
		schema = &AVL{root: nil}
	case "rb":
		schema = &RB{root: nil}
	}

	return pool.(tree).insert(name, schema)
}

func (t Tree) DeleteSchema(settings map[string]string, name string, poolName string) string {
	pool, ok := t.self.find(poolName)
	if !ok {
		fmt.Printf("Pool %s does not exist\n", poolName)
		return "error"
	}

	return pool.(tree).remove(name)
}

func (t Tree) CreateCollection(settings map[string]string, name string, poolName string, schemaName string) string {
	pool, ok := t.self.find(poolName)
	if !ok {
		fmt.Printf("Pool %s does not exist\n", poolName)
		return "error"
	}

	schema, ok := pool.(tree).find(schemaName)
	if !ok {
		fmt.Printf("Schema %s does not exist in pool %s\n", schemaName, poolName)
		return "error"
	}

	_, ok = schema.(tree).find(name)
	if ok {
		fmt.Printf("Collection %s already exists in schema %s in pool %s\n", name, schemaName, poolName)
		return "error"
	}

	var collection tree
	switch t.variant {
	case "btree":
		collection = &Btree{root: nil, t: schema.(*Btree).t}
	case "avl":
		collection = &AVL{root: nil}
	case "rb":
		collection = &RB{root: nil}
	}

	return schema.(tree).insert(name, collection)
}

func (t Tree) DeleteCollection(settings map[string]string, name string, poolName string, schemaName string) string {
	pool, ok := t.self.find(poolName)
	if !ok {
		fmt.Printf("Pool %s does not exist\n", poolName)
		return "error"
	}

	schema, ok := pool.(tree).find(schemaName)
	if !ok {
		fmt.Printf("Schema %s does not exist in pool %s\n", schemaName, poolName)
		return "error"
	}

	return schema.(tree).remove(name)
}

func (t Tree) Set(settings map[string]string, key string, value []string, poolName string, schemaName string, collectionName string) string {
	pool, ok := t.self.find(poolName)
	if !ok {
		fmt.Printf("Pool %s does not exist\n", poolName)
		return "error"
	}

	schema, ok := pool.(tree).find(schemaName)
	if !ok {
		fmt.Printf("Schema %s does not exist in pool %s\n", schemaName, poolName)
		return "error"
	}

	collection, ok := schema.(tree).find(collectionName)
	if !ok {
		fmt.Printf("Collection %s does not exist in schema %s in pool %s\n", collectionName, schemaName, poolName)
		return "error"
	}

	_, ok = collection.(tree).find(key)
	if ok {
		fmt.Printf("Key %s already exists in collection %s in schema %s in pool %s\n", key, collectionName, schemaName, poolName)
		return "error"
	}

	return collection.(tree).insert(key, value)
}

func (t Tree) Update(settings map[string]string, key string, value []string, poolName string, schemaName string, collectionName string) string {
	pool, ok := t.self.find(poolName)
	if !ok {
		fmt.Printf("Pool %s does not exist\n", poolName)
		return "error"
	}

	schema, ok := pool.(tree).find(schemaName)
	if !ok {
		fmt.Printf("Schema %s does not exist in pool %s\n", schemaName, poolName)
		return "error"
	}

	collection, ok := schema.(tree).find(collectionName)
	if !ok {
		fmt.Printf("Collection %s does not exist in schema %s in pool %s\n", collectionName, schemaName, poolName)
		return "error"
	}

	_, ok = collection.(tree).find(key)
	if !ok {
		fmt.Printf("Key %s does not exist in collection %s in schema %s in pool %s\n", key, collectionName, schemaName, poolName)
		return "error"
	}

	return collection.(tree).update(key, value)
}

func (t Tree) Get(settings map[string]string, key string, poolName string, schemaName string, collectionName string) (string, string) {
	pool, ok := t.self.find(poolName)
	if !ok {
		fmt.Printf("Pool %s does not exist\n", poolName)
		return "", "error"
	}

	schema, ok := pool.(tree).find(schemaName)
	if !ok {
		fmt.Printf("Schema %s does not exist in pool %s\n", schemaName, poolName)
		return "", "error"
	}

	collection, ok := schema.(tree).find(collectionName)
	if !ok {
		fmt.Printf("Collection %s does not exist in schema %s in pool %s\n", collectionName, schemaName, poolName)
		return "", "error"
	}

	value, ok := collection.(tree).find(key)
	if !ok {
		fmt.Printf("Key %s does not exist in collection %s in schema %s in pool %s\n", key, collectionName, schemaName, poolName)
		return "", "error"
	}

	return fmt.Sprintf("%s", value), "ok"
}

func (t Tree) GetRange(settings map[string]string, leftBound string, rightBound string, poolName string, schemaName string, collectionName string) (map[string]string, string) {
	// TODO

	return nil, "ok"
}

func (t Tree) Delete(settings map[string]string, key string, poolName string, schemaName string, collectionName string) string {
	pool, ok := t.self.find(poolName)
	if !ok {
		fmt.Printf("Pool %s does not exist\n", poolName)
		return "error"
	}

	schema, ok := pool.(tree).find(schemaName)
	if !ok {
		fmt.Printf("Schema %s does not exist in pool %s\n", schemaName, poolName)
		return "error"
	}

	collection, ok := schema.(tree).find(collectionName)
	if !ok {
		fmt.Printf("Collection %s does not exist in schema %s in pool %s\n", collectionName, schemaName, poolName)
		return "error"
	}

	return collection.(tree).remove(key)
}
