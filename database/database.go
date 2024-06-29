package database

import (
	"DBMS/task0+3"
	"DBMS/task2"
	"fmt"
	"strings"
	"time"
)

func ExecuteCommand(command string) string {
	words := strings.Fields(command)
	db := database

	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	switch strings.ToLower(words[0]) {

	// Create and Delete Pool
	case "createpool", "deletepool":
		if len(words) > 2 {
			fmt.Println("Too many arguments")
			return "error"
		}
		name := words[1]

		if strings.ToLower(words[0]) == "createpool" {
			return db.CreatePool(settings, name)
		} else if strings.ToLower(words[0]) == "deletepool" {
			return db.DeletePool(settings, name)
		}

		fmt.Println("Wrong command")
		return "error"

	// Create and Delete Schema
	case "createschema", "deleteschema":
		if len(words) > 4 {
			fmt.Println("Too many arguments")
			return "error"
		}

		name := words[1]
		if words[2] != "in" {
			fmt.Println("Wrong command on schema name declaration")
			return "error"
		}
		poolName := words[3]

		if strings.ToLower(words[0]) == "deleteschema" {
			return db.DeleteSchema(settings, name, poolName)
		} else if strings.ToLower(words[0]) == "createschema" {
			return db.CreateSchema(settings, name, poolName)
		}

		fmt.Println("Wrong command")
		return "error"

	// Create Collection
	case "createcollection":
		if len(words) > 5 {
			fmt.Println("Too many arguments")
			return "error"
		} else if len(words) < 5 {
			fmt.Println("Too few arguments")
			return "error"
		}

		collType := words[1]
		name := words[2]
		if words[3] != "in" {
			fmt.Println("Wrong command on schema name declaration")
			return "error"
		}
		schemaAndPool := strings.Split(words[4], ".")
		if len(schemaAndPool) != 2 {
			fmt.Println("Wrong command on schema name declaration")
			return "error"
		}
		poolName := schemaAndPool[0]
		schemaName := schemaAndPool[1]

		response := db.CreateCollection(settings, name, collType, poolName, schemaName)

		if response != "ok" {
			return response
		}

		if settings["persistant"] == "on" {
			task2.SaveCollection(name, db.GetCollection(settings, poolName, schemaName, name))
		}

		return response

	// Delete Collection
	case "deletecollection":
		if len(words) > 4 {
			fmt.Println("Too many arguments")
			return "error"
		} else if len(words) < 4 {
			fmt.Println("Too few arguments")
			return "error"
		}

		name := words[1]
		if words[2] != "in" {
			fmt.Println("Wrong command on schema name declaration")
			return "error"
		}
		schemaAndPool := strings.Split(words[3], ".")
		if len(schemaAndPool) != 2 {
			fmt.Println("Wrong command on schema name declaration")
			return "error"
		}
		poolName := schemaAndPool[0]
		schemaName := schemaAndPool[1]

		response := db.DeleteCollection(settings, name, poolName, schemaName)

		if response != "ok" {
			return response
		}

		if settings["persistant"] == "on" {
			task2.DeleteCollection(name)
		}

		return response

	// Set and Update value
	case "update":
		if len(words) < 5 {
			fmt.Println("Too few arguments")
			return "error"
		}

		key := words[1]
		value := strings.Join(words[2:len(words)-2], " ")
		if words[len(words)-2] != "in" {
			fmt.Println("Wrong command on collection name declaration")
			return "error"
		}
		poolAndSchemaAndCollection := strings.Split(words[len(words)-1], ".")
		if len(poolAndSchemaAndCollection) != 3 {
			fmt.Println("Wrong command on collection name declaration")
			return "error"
		}
		poolName := poolAndSchemaAndCollection[0]
		schemaName := poolAndSchemaAndCollection[1]
		collectionName := poolAndSchemaAndCollection[2]

		response := db.Update(settings, key, value, poolName, schemaName, collectionName)
		if response != "ok" {
			return response
		}
		if settings["persistant"] == "on" {
			task2.AddCommand(collectionName, command)
		}
		return response

	case "set":
		if len(words) < 6 {
			fmt.Println("Too few arguments")
			return "error"
		}

		key := words[1]
		secondaryKey := words[2]
		value := strings.Join(words[3:len(words)-2], " ")
		if words[len(words)-2] != "in" {
			fmt.Println("Wrong command on collection name declaration")
			return "error"
		}
		poolAndSchemaAndCollection := strings.Split(words[len(words)-1], ".")
		if len(poolAndSchemaAndCollection) != 3 {
			fmt.Println("Wrong command on collection name declaration")
			return "error"
		}
		poolName := poolAndSchemaAndCollection[0]
		schemaName := poolAndSchemaAndCollection[1]
		collectionName := poolAndSchemaAndCollection[2]

		response := db.Set(settings, key, secondaryKey, value, poolName, schemaName, collectionName)
		if response != "ok" {
			return response
		}
		if settings["persistant"] == "on" {
			task2.AddCommand(collectionName, command)
		}
		return response

	// Get and Delete value
	case "get", "getsecondary", "delete":
		if len(words) < 4 {
			fmt.Println("Too few arguments")
			return "error"
		}

		key := words[1]
		if words[2] != "in" {
			fmt.Println("Wrong command on collection name declaration")
			return "error"
		}
		poolAndSchemaAndCollection := strings.Split(words[len(words)-1], ".")
		if len(poolAndSchemaAndCollection) != 3 {
			fmt.Println("Wrong command on collection name declaration")
			return "error"
		}
		poolName := poolAndSchemaAndCollection[0]
		schemaName := poolAndSchemaAndCollection[1]
		collectionName := poolAndSchemaAndCollection[2]

		if strings.ToLower(words[0]) == "get" {
			value := db.Get(settings, key, poolName, schemaName, collectionName)
			if value != "" {
				return key + " = " + value
			}
			return "error"
		} else if strings.ToLower(words[0]) == "getsecondary" {
			value := db.GetBySecondaryKey(settings, key, poolName, schemaName, collectionName)
			if value != "" {
				return key + " = " + value
			}
			return "error"
		} else if strings.ToLower(words[0]) == "delete" {
			response := db.Delete(settings, key, poolName, schemaName, collectionName)
			if response != "ok" {
				return response
			}
			if settings["persistant"] == "on" {
				task2.AddCommand(collectionName, command)
			}
			return response
		}

		fmt.Println("Wrong command")
		return "error"

	// Get range
	case "getrange", "getrangesecondary":
		if len(words) < 5 {
			fmt.Println("Too few arguments")
			return "error"
		}

		leftBound := words[1]
		rightBound := words[2]
		if words[3] != "in" {
			fmt.Println("Wrong command on collection name declaration")
			return "error"
		}
		poolAndSchemaAndCollection := strings.Split(words[len(words)-1], ".")
		if len(poolAndSchemaAndCollection) != 3 {
			fmt.Println("Wrong command on collection name declaration")
			return "error"
		}
		poolName := poolAndSchemaAndCollection[0]
		schemaName := poolAndSchemaAndCollection[1]
		collectionName := poolAndSchemaAndCollection[2]

		if strings.ToLower(words[0]) == "getrange" {
			result := db.GetRange(settings, leftBound, rightBound, poolName, schemaName, collectionName)
			if result != nil {
				if len(*result) == 0 {
					return "no data in range " + leftBound + " - " + rightBound
				}
				str := new(strings.Builder)
				for k, v := range *result {
					str.WriteString(k + " = " + v + "\n")
				}
				return str.String()
			}
			return "error"
		} else if strings.ToLower(words[0]) == "getrangesecondary" {
			result := db.GetRangeBySecondaryKey(settings, leftBound, rightBound, poolName, schemaName, collectionName)
			if result != nil {
				if len(*result) == 0 {
					return "no data in range " + leftBound + " - " + rightBound
				}
				str := new(strings.Builder)
				for k, v := range *result {
					str.WriteString(k + " = " + v + "\n")
				}
				return str.String()
			}
			return "error"
		}
		return "error"

	// Get by time if persistant is on
	case "getat":
		if settings["persistant"] != "on" {
			fmt.Println("Persistant is not on")
			return "error"
		}

		if len(words) < 6 {
			fmt.Println("Too few arguments")
			return "error"
		}

		time, err := time.Parse("2006-01-02 15:04:05.000000 MST", words[1]+" "+words[2]+" "+words[3])
		if err != nil {
			fmt.Printf("Error on parsing time: %v\n", err.Error())
			return "error"
		}

		key := words[4]
		if words[5] != "in" {
			fmt.Println("Wrong command on collection name declaration")
			return "error"
		}
		poolAndSchemaAndCollection := strings.Split(words[len(words)-1], ".")
		if len(poolAndSchemaAndCollection) != 3 {
			fmt.Println("Wrong command on collection name declaration")
			return "error"
		}
		collectionName := poolAndSchemaAndCollection[2]

		value := task2.GetValueByTime(collectionName, key, time)
		if value != "" {
			return key + " = " + value + " (at " + time.Format("2006-01-02 15:04:05.000000 MST") + ")"
		}
		return "error"
	}

	fmt.Println("Something went wrong")
	return "error"
}

func ListPools() []string {
	return *Database().ListPools(settings)
}

func ListSchemas(pool string) []string {
	return *Database().ListSchemas(settings, pool)
}

func ListCollections(pool string, schema string) []string {
	return *Database().ListCollections(settings, pool, schema)
}

func ListDatas(pool string, schema string, collection string) [](task0.Datas) {
	return *Database().GetAll(settings, pool, schema, collection)
}
