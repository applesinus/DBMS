package task0

import (
	"fmt"
)

func (collection *CollectionBI) Set(key string, value string) string {
	if _, ok := collection.value[key]; ok {
		fmt.Printf("Key %s already exists\n", key)
		return "error"
	}

	collection.value[key] = value

	return "ok"
}

func (collection *CollectionBI) Update(key string, value string) string {
	if _, ok := collection.value[key]; !ok {
		fmt.Printf("Key %s does not exist\n", key)
		return "error"
	}

	collection.value[key] = value

	return "ok"
}

func (collection *CollectionBI) Get(key string) (string, string) {
	if _, ok := collection.value[key]; !ok {
		fmt.Printf("Key %s does not exist\n", key)
		return "", "error"
	}

	return collection.value[key], "ok"
}

func (collection *CollectionBI) GetRange(leftBound string, rightBound string) (*map[string]string, string) {
	result := make(map[string]string, 0)
	for k, v := range collection.value {
		if k >= leftBound && k <= rightBound {
			result[k] = v
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
