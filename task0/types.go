package task0

import (
	"DBMS/interfaces"
	"DBMS/task4"
)

// Database
type Database struct {
	pools map[string]Pool
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
	value map[string]task4.TrieWord
}
