package task0

// DatabaseBI - Built-in representations of database
type DatabaseBI struct {
	pools map[string]PoolBI
}

// PoolBI - Built-in representation of pool
type PoolBI struct {
	name   string
	schema map[string]SchemaBI
}

// SchemaBI - Built-in representation of schema
type SchemaBI struct {
	name       string
	collection map[string]CollectionBI
}

// CollectionBI - Built-in representation of collection
type CollectionBI struct {
	name  string
	value map[string]string
}
