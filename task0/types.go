package task0

import (
	"DBMS/task4"
)

// collectionInterface - basic collection interface
type CollectionInterface interface {
	Set(key string, value string) string
	Update(key string, value string) string
	Get(key string) (string, string)
	GetRange(leftBound string, rightBound string) (*map[string]string, string)
	Delete(key string) string
}

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
	collection map[string]CollectionInterface
}

// CollectionBI - Built-in representation of collection
type CollectionBI struct {
	name  string
	value map[string]task4.TrieWord
}
