package task0

import (
	"DBMS/interfaces"
	"DBMS/task4"
	"fmt"
)

func (collection *CollectionBI) Set(key string, secondaryKey string, value string) string {
	for _, k := range collection.value {
		if k.key == key {
			fmt.Printf("Key %s already exists\n", key)
			return "error"
		}
	}

	word, ok := task4.Pool.Insert(value)
	if ok != "ok" {
		return "error"
	}
	collection.value = append(collection.value, Value{key, secondaryKey, *word})

	return "ok"
}

func (collection *CollectionBI) Update(key string, value string) string {
	index := -1
	for i, k := range collection.value {
		if k.key == key {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Printf("Key %s does not exist\n", key)
		return "error"
	}

	word, ok := task4.Pool.Insert(value)
	if ok != "ok" {
		return "error"
	}
	collection.value[index].value = *word

	return "ok"
}

func (collection *CollectionBI) Get(key string) (string, string) {
	index := -1
	for i, k := range collection.value {
		if k.key == key {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Printf("Key %s does not exist\n", key)
		return "", "error"
	}

	str, ok := collection.value[index].value.String()

	if !ok {
		return "", "error"
	}

	return str, "ok"
}

func (collection *CollectionBI) GetBySecondaryKey(secondaryKey string) (string, string) {
	index := -1
	for i, k := range collection.value {
		if k.secondaryKey == secondaryKey {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Printf("Secondary key %s does not exist\n", secondaryKey)
		return "", "error"
	}

	str, ok := collection.value[index].value.String()

	if !ok {
		return "", "error"
	}

	return str, "ok"
}

func (collection *CollectionBI) GetRange(leftBound string, rightBound string) (*map[string]string, string) {
	result := make(map[string]string, 0)
	for _, value := range collection.value {
		if value.key >= leftBound && value.key <= rightBound {
			res, ok := value.value.String()
			if !ok {
				return nil, "error"
			}
			result[value.key] = res
		}
	}

	return &result, "ok"
}

func (collection *CollectionBI) GetRangeBySecondaryKey(leftBound string, rightBound string) (*map[string]string, string) {
	result := make(map[string]string, 0)
	for _, value := range collection.value {
		if value.secondaryKey >= leftBound && value.secondaryKey <= rightBound {
			res, ok := value.value.String()
			if !ok {
				return nil, "error"
			}
			result[value.secondaryKey] = res
		}
	}

	return &result, "ok"
}

func (collection *CollectionBI) GetAll() (*[]string, *[]string, *[]string, string) {
	keys := make([]string, 0)
	secondaryKeys := make([]string, 0)
	values := make([]string, 0)
	for _, value := range collection.value {
		keys = append(keys, value.key)
		secondaryKeys = append(secondaryKeys, value.secondaryKey)
		res, ok := value.value.String()
		if !ok {
			return nil, nil, nil, "error"
		}
		values = append(values, res)
	}

	return &keys, &secondaryKeys, &values, "ok"
}

func (collection *CollectionBI) Delete(key string) string {
	index := -1
	for i, k := range collection.value {
		if k.key == key {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Printf("Key %s does not exist\n", key)
		return "error"
	}

	if index == len(collection.value)-1 {
		collection.value = collection.value[:index]
	} else {
		collection.value = append(collection.value[:index], collection.value[index+1:]...)
	}

	return "ok"
}

func (collection *CollectionBI) Copy() interfaces.CollectionInterface {
	return &CollectionBI{name: collection.name, value: collection.value}
}
