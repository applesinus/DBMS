package task0

import (
	"DBMS/interfaces"
	"DBMS/task6+3"
	"fmt"
	"sync"
)

func CreateDB() *Database {
	return &Database{
		pools: make(map[string]Pool),
		Mutex: &sync.Mutex{},
	}
}

func (db *Database) CreatePool(settings map[string]string, name string) string {
	if _, ok := db.pools[name]; ok {
		fmt.Printf("Pool %s already exists\n", name)
		return "error"
	}

	db.pools[name] = Pool{
		name:   name,
		schema: make(map[string]Schema),
	}
	return "ok"
}

func (db *Database) checkPool(name string) bool {
	_, ok := db.pools[name]
	if !ok {
		fmt.Printf("Pool %s does not exist\n", name)
	}
	return ok
}

func (db *Database) DeletePool(settings map[string]string, name string) string {
	if !db.checkPool(name) {
		return "error"
	}

	delete(db.pools, name)
	return "ok"
}

func (db *Database) ListPools(settings map[string]string) *[]string {
	pools := make([]string, 0)
	for _, pool := range db.pools {
		pools = append(pools, pool.name)
	}
	return &pools
}

func (db *Database) CreateSchema(settings map[string]string, name string, pool string) string {
	if !db.checkPool(pool) {
		return "error"
	}

	if _, ok := db.pools[pool].schema[name]; ok {
		fmt.Printf("Schema %s already exists in pool %s\n", name, pool)
		return "error"
	}

	db.pools[pool].schema[name] = Schema{
		name:       name,
		collection: make(map[string]interfaces.CollectionInterface),
	}
	return "ok"
}

func (db *Database) checkSchema(pool string, name string) bool {
	if !db.checkPool(pool) {
		return false
	}
	_, ok := db.pools[pool].schema[name]
	if !ok {
		fmt.Printf("Schema %s does not exist in pool %s\n", name, pool)
	}
	return ok
}

func (db *Database) DeleteSchema(settings map[string]string, name string, pool string) string {
	if !db.checkSchema(pool, name) {
		return "error"
	}

	delete(db.pools[pool].schema, name)
	return "ok"
}

func (db *Database) ListSchemas(settings map[string]string, pool string) *[]string {
	schemas := make([]string, 0)
	for _, schema := range db.pools[pool].schema {
		schemas = append(schemas, schema.name)
	}
	return &schemas
}

func (db *Database) ListCollections(settings map[string]string, pool string, schema string) *[]string {
	collections := make([]string, 0)
	for collection := range db.pools[pool].schema[schema].collection {
		collections = append(collections, collection)
	}
	return &collections
}

func (db *Database) CreateCollection(settings map[string]string, name string, collType string, pool string, schema string) string {
	if !db.checkSchema(pool, schema) {
		return "error"
	}

	if _, ok := db.pools[pool].schema[schema].collection[name]; ok {
		fmt.Printf("Collection %s already exists in schema %s in pool %s\n", name, schema, pool)
		return "error"
	}

	if collType == "BI" {
		db.pools[pool].schema[schema].collection[name] = &CollectionBI{
			name:  name,
			value: make([]Value, 0),
		}
		return "ok"
	}

	if collType == "AVL" || collType == "RB" || (len(collType) >= 5 && collType[:5] == "Btree") {
		tree := task6.NewTree(collType)
		if tree == nil {
			return "error"
		}
		db.pools[pool].schema[schema].collection[name] = tree
		return "ok"
	}

	return "error"
}

func (db *Database) checkCollection(pool string, schema string, name string) bool {
	if !db.checkSchema(pool, schema) {
		return false
	}
	_, ok := db.pools[pool].schema[schema].collection[name]
	if !ok {
		fmt.Printf("Collection %s does not exist in schema %s in pool %s\n", name, schema, pool)
	}
	return ok
}

func (db *Database) DeleteCollection(settings map[string]string, name string, pool string, schema string) string {
	if !db.checkCollection(pool, schema, name) {
		return "error"
	}
	delete(db.pools[pool].schema[schema].collection, name)
	return "ok"
}

func (db *Database) GetCollection(settings map[string]string, pool string, schema string, coll string) interfaces.CollectionInterface {
	if !db.checkCollection(pool, schema, coll) {
		return nil
	}
	return db.pools[pool].schema[schema].collection[coll]
}

func (db *Database) Get(settings map[string]string, key string, pool string, schema string, coll string) string {
	if !db.checkCollection(pool, schema, coll) {
		return ""
	}
	res, ok := db.pools[pool].schema[schema].collection[coll].Get(key)
	if ok != "ok" {
		return ""
	}
	return res
}

func (db *Database) GetBySecondaryKey(settings map[string]string, secondaryKey string, pool string, schema string, coll string) string {
	if !db.checkCollection(pool, schema, coll) {
		return ""
	}
	res, ok := db.pools[pool].schema[schema].collection[coll].GetBySecondaryKey(secondaryKey)
	if ok != "ok" {
		return ""
	}
	return res
}

func (db *Database) GetRange(settings map[string]string, leftBound string, rightBound string, pool string, schema string, coll string) *map[string]string {
	if !db.checkCollection(pool, schema, coll) {
		ret := make(map[string]string)
		return &ret
	}
	res, ok := db.pools[pool].schema[schema].collection[coll].GetRange(leftBound, rightBound)
	if ok != "ok" {
		ret := make(map[string]string)
		return &ret
	}
	return res
}

func (db *Database) GetRangeBySecondaryKey(settings map[string]string, leftBound string, rightBound string, pool string, schema string, coll string) *map[string]string {
	if !db.checkCollection(pool, schema, coll) {
		ret := make(map[string]string)
		return &ret
	}
	res, ok := db.pools[pool].schema[schema].collection[coll].GetRangeBySecondaryKey(leftBound, rightBound)
	if ok != "ok" {
		ret := make(map[string]string)
		return &ret
	}
	return res
}

func (db *Database) GetAll(settings map[string]string, pool string, schema string, coll string) *[]Datas {
	if !db.checkCollection(pool, schema, coll) {
		return nil
	}
	keys, secs, vals, ok := db.pools[pool].schema[schema].collection[coll].GetAll()
	if ok != "ok" {
		return nil
	}
	res := make([]Datas, len(*keys))
	for i := 0; i < len(*keys); i++ {
		res[i] = Datas{Key: (*keys)[i], SecondaryKey: (*secs)[i], Value: (*vals)[i]}
	}
	return &res
}

func (db *Database) Set(settings map[string]string, key string, secondaryKey string, value string, pool string, schema string, coll string) string {
	if !db.checkCollection(pool, schema, coll) {
		return "error"
	}
	ok := db.pools[pool].schema[schema].collection[coll].Set(key, secondaryKey, value)
	if ok != "ok" {
		return "error"
	}
	return "ok"
}

func (db *Database) Update(settings map[string]string, key string, value string, pool string, schema string, coll string) string {
	if !db.checkCollection(pool, schema, coll) {
		return "error"
	}
	ok := db.pools[pool].schema[schema].collection[coll].Update(key, value)
	if ok != "ok" {
		return "error"
	}
	return "ok"
}

func (db *Database) Delete(settings map[string]string, key string, pool string, schema string, coll string) string {
	if !db.checkCollection(pool, schema, coll) {
		return "error"
	}
	ok := db.pools[pool].schema[schema].collection[coll].Delete(key)
	return ok
}
