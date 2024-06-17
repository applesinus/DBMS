package task0

import (
	"DBMS/interfaces"
	"DBMS/task4"
	"fmt"
)

func (collection *CollectionBI) Set(key string, value string) string {
	if _, ok := collection.value[key]; ok {
		fmt.Printf("Key %s already exists\n", key)
		return "error"
	}

	word, ok := task4.Pool.Insert(value)
	if ok != "ok" {
		return "error"
	}
	collection.value[key] = *word

	return "ok"
}

func (collection *CollectionBI) Update(key string, value string) string {
	if _, ok := collection.value[key]; !ok {
		fmt.Printf("Key %s does not exist\n", key)
		return "error"
	}

	word, ok := task4.Pool.Insert(value)
	if ok != "ok" {
		return "error"
	}
	collection.value[key] = *word

	return "ok"
}

func (collection *CollectionBI) Get(key string) (string, string) {
	if _, ok := collection.value[key]; !ok {
		fmt.Printf("Key %s does not exist\n", key)
		return "", "error"
	}

	str, ok := collection.value[key].String()

	if !ok {
		return "", "error"
	}

	return str, "ok"
}

func (collection *CollectionBI) GetRange(leftBound string, rightBound string) (*map[string]string, string) {
	result := make(map[string]string, 0)
	for k, v := range collection.value {
		if k >= leftBound && k <= rightBound {
			res, ok := v.String()
			if !ok {
				return nil, "error"
			}
			result[k] = res
		}
	}

	return &result, "ok"
}

func (collection *CollectionBI) Delete(key string) string {
	if _, ok := collection.value[key]; !ok {
		fmt.Printf("Key %s does not exist\n", key)
		return "error"
	}

	delete(collection.value, key)

	return "ok"
}

func (collection *CollectionBI) Copy() interfaces.CollectionInterface {
	return &CollectionBI{name: collection.name, value: collection.value}
}
