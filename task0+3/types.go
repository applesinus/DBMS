package task0

import (
	"DBMS/interfaces"
	"DBMS/task4"
	"sync"
)

// Database
type Database struct {
	pools map[string]Pool
	Mutex *sync.Mutex
}

// Pool
type Pool struct {
	name   string
	schema map[string]Schema
}

// Schema
type Schema struct {
	name       string
	collection map[string]interfaces.CollectionInterface
}

// CollectionBI - Built-in representation of collection
type CollectionBI struct {
	name  string
	value []Value
}

type Value struct {
	key          string
	secondaryKey string
	value        task4.TrieWord
}

type Datas struct {
	Key          string
	SecondaryKey string
	Value        string
}
