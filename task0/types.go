package task0

// collectionInterface - basic collection interface
type CollectionInterface interface {
	Set(key string, value string) string
	Update(key string, value string) string
	Get(key string) (string, string)
	GetRange(leftBound string, rightBound string) (*map[string]string, string)
	Delete(key string) string
}

// CollectionBI - Built-in representation of collection
type CollectionBI struct {
	name  string
	value map[string]string
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
