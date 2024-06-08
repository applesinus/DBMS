package task0

import (
	"fmt"
	"strings"
)

func CreateDB() DatabaseBI {
	return DatabaseBI{pools: make(map[string]PoolBI)}
}

func (db DatabaseBI) CreatePool(settings map[string]string, name string) string {
	if _, ok := db.pools[name]; ok {
		fmt.Printf("Pool %s already exists\n", name)
		return "error"
	}

	db.pools[name] = PoolBI{
		name:   name,
		schema: make(map[string]SchemaBI),
	}
	return "ok"
}

func (db DatabaseBI) DeletePool(settings map[string]string, name string) string {
	if _, ok := db.pools[name]; !ok {
		fmt.Printf("Pool %s does not exist\n", name)
		return "error"
	}

	delete(db.pools, name)
	return "ok"
}

func (db DatabaseBI) CreateSchema(settings map[string]string, name string, pool string) string {
	if _, ok := db.pools[pool]; !ok {
		fmt.Printf("Pool %s does not exist\n", pool)
		return "error"
	}

	if _, ok := db.pools[pool].schema[name]; ok {
		fmt.Printf("Schema %s already exists in pool %s\n", name, pool)
		return "error"
	}

	db.pools[pool].schema[name] = SchemaBI{
		name:       name,
		collection: make(map[string]CollectionBI),
	}
	return "ok"
}

func (db DatabaseBI) DeleteSchema(settings map[string]string, name string, pool string) string {
	if _, ok := db.pools[pool]; !ok {
		fmt.Printf("Pool %s does not exist\n", pool)
		return "error"
	}

	if _, ok := db.pools[pool].schema[name]; !ok {
		fmt.Printf("Schema %s does not exist in pool %s\n", name, pool)
		return "error"
	}

	delete(db.pools[pool].schema, name)
	return "ok"
}

func (db DatabaseBI) CreateCollection(settings map[string]string, name string, pool string, schema string) string {
	if _, ok := db.pools[pool]; !ok {
		fmt.Printf("Pool %s does not exist\n", pool)
		return "error"
	}

	if _, ok := db.pools[pool].schema[schema]; !ok {
		fmt.Printf("Schema %s does not exist in pool %s\n", schema, pool)
		return "error"
	}

	if _, ok := db.pools[pool].schema[schema].collection[name]; ok {
		fmt.Printf("Collection %s already exists in schema %s in pool %s\n", name, schema, pool)
		return "error"
	}

	db.pools[pool].schema[schema].collection[name] = CollectionBI{
		name:  name,
		value: make(map[string]string),
	}
	return "ok"
}

func (db DatabaseBI) DeleteCollection(settings map[string]string, name string, pool string, schema string) string {
	if _, ok := db.pools[pool]; !ok {
		fmt.Printf("Pool %s does not exist\n", pool)
		return "error"
	}

	if _, ok := db.pools[pool].schema[schema]; !ok {
		fmt.Printf("Schema %s does not exist in pool %s\n", schema, pool)
		return "error"
	}

	if _, ok := db.pools[pool].schema[schema].collection[name]; !ok {
		fmt.Printf("Collection %s does not exist in schema %s in pool %s\n", name, schema, pool)
		return "error"
	}

	delete(db.pools[pool].schema[schema].collection, name)
	return "ok"
}

func (db DatabaseBI) Set(settings map[string]string, key string, value []string, pool string, schema string, collection string) string {
	if _, ok := db.pools[pool]; !ok {
		fmt.Printf("Pool %s does not exist\n", pool)
		return "error"
	}
	if _, ok := db.pools[pool].schema[schema]; !ok {
		fmt.Printf("Schema %s does not exist in pool %s\n", schema, pool)
		return "error"
	}
	if _, ok := db.pools[pool].schema[schema].collection[collection]; !ok {
		fmt.Printf("Collection %s does not exist in schema %s in pool %s\n", collection, schema, pool)
		return "error"
	}

	if _, ok := db.pools[pool].schema[schema].collection[collection].value[key]; ok {
		fmt.Printf("Key %s already exists in collection %s in schema %s in pool %s\n", key, collection, schema, pool)
		return "error"
	}

	db.pools[pool].schema[schema].collection[collection].value[key] = strings.Join(value, ";")

	return "ok"
}

func (db DatabaseBI) Update(settings map[string]string, key string, value []string, pool string, schema string, collection string) string {
	if _, ok := db.pools[pool]; !ok {
		fmt.Printf("Pool %s does not exist\n", pool)
		return "error"
	}
	if _, ok := db.pools[pool].schema[schema]; !ok {
		fmt.Printf("Schema %s does not exist in pool %s\n", schema, pool)
		return "error"
	}
	if _, ok := db.pools[pool].schema[schema].collection[collection]; !ok {
		fmt.Printf("Collection %s does not exist in schema %s in pool %s\n", collection, schema, pool)
		return "error"
	}

	if _, ok := db.pools[pool].schema[schema].collection[collection].value[key]; !ok {
		fmt.Printf("Key %s does not exist in collection %s in schema %s in pool %s\n", key, collection, schema, pool)
		return "error"
	}

	db.pools[pool].schema[schema].collection[collection].value[key] = strings.Join(value, ";")

	return "ok"
}

func (db DatabaseBI) Get(settings map[string]string, key string, pool string, schema string, collection string) (string, string) {
	if _, ok := db.pools[pool]; !ok {
		fmt.Printf("Pool %s does not exist\n", pool)
		return "", "error"
	}
	if _, ok := db.pools[pool].schema[schema]; !ok {
		fmt.Printf("Schema %s does not exist in pool %s\n", schema, pool)
		return "", "error"
	}
	if _, ok := db.pools[pool].schema[schema].collection[collection]; !ok {
		fmt.Printf("Collection %s does not exist in schema %s in pool %s\n", collection, schema, pool)
		return "", "error"
	}

	if _, ok := db.pools[pool].schema[schema].collection[collection].value[key]; !ok {
		fmt.Printf("Key %s does not exist in collection %s in schema %s in pool %s\n", key, collection, schema, pool)
		return "", "error"
	}

	return db.pools[pool].schema[schema].collection[collection].value[key], "ok"
}

func (db DatabaseBI) GetRange(settings map[string]string, leftBound string, rightBound string, pool string, schema string, collection string) (map[string]string, string) {
	if _, ok := db.pools[pool]; !ok {
		fmt.Printf("Pool %s does not exist\n", pool)
		return nil, "error"
	}
	if _, ok := db.pools[pool].schema[schema]; !ok {
		fmt.Printf("Schema %s does not exist in pool %s\n", schema, pool)
		return nil, "error"
	}
	if _, ok := db.pools[pool].schema[schema].collection[collection]; !ok {
		fmt.Printf("Collection %s does not exist in schema %s in pool %s\n", collection, schema, pool)
		return nil, "error"
	}

	result := make(map[string]string, 0)
	for k, v := range db.pools[pool].schema[schema].collection[collection].value {
		if k >= leftBound && k <= rightBound {
			result[k] = v
		}
	}

	return result, "ok"
}

func (db DatabaseBI) Delete(settings map[string]string, key string, pool string, schema string, collection string) string {
	if _, ok := db.pools[pool]; !ok {
		fmt.Printf("Pool %s does not exist\n", pool)
		return "error"
	}
	if _, ok := db.pools[pool].schema[schema]; !ok {
		fmt.Printf("Schema %s does not exist in pool %s\n", schema, pool)
		return "error"
	}
	if _, ok := db.pools[pool].schema[schema].collection[collection]; !ok {
		fmt.Printf("Collection %s does not exist in schema %s in pool %s\n", collection, schema, pool)
		return "error"
	}
	if _, ok := db.pools[pool].schema[schema].collection[collection].value[key]; !ok {
		fmt.Printf("Key %s does not exist in collection %s in schema %s in pool %s\n", key, collection, schema, pool)
		return "error"
	}

	delete(db.pools[pool].schema[schema].collection[collection].value, key)

	return "ok"
}
