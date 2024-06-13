package main

import (
	"DBMS/task0"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// database - basic database interface
type database interface {
	CreatePool(settings map[string]string, name string) string
	DeletePool(settings map[string]string, name string) string
	CreateSchema(settings map[string]string, name string, pool string) string
	DeleteSchema(settings map[string]string, name string, pool string) string
	CreateCollection(settings map[string]string, name string, pool string, schema string) string
	DeleteCollection(settings map[string]string, name string, pool string, schema string) string
	Set(settings map[string]string, key string, value []string, pool string, schema string, collection string) string
	Update(settings map[string]string, key string, value []string, pool string, schema string, collection string) string
	Get(settings map[string]string, key string, pool string, schema string, collection string) (string, string)
	GetRange(settings map[string]string, leftBound string, rightBound string, pool string, schema string, collection string) (map[string]string, string)
	Delete(settings map[string]string, key string, pool string, schema string, collection string) string
}

func executeCommand(db database, settings map[string]string, command string) string {
	words := strings.Fields(command)
	switch len(words) {
	// Empty line
	case 0:
		return "empty"

	// Single command (application side)
	case 1:
		switch strings.ToLower(command) {
		case "stop", "exit", "-1":
			return "exit"
		case "help":
			help()
			return "ok"
		default:
			_, err := os.Open(words[0])
			if err == nil {
				ret := executeFile(db, settings, words[0])
				if ret == "error" {
					fmt.Println("Could not execute file")
				} else {
					return ret
				}
			}
			fmt.Printf("%s\n", command)
			return "error"
		}

	// Multiple commands (database side)
	default:
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

		// Create and Delete Collection
		case "createcollection", "deletecollection":
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

			if strings.ToLower(words[0]) == "createcollection" {
				return db.CreateCollection(settings, name, poolName, schemaName)
			} else if strings.ToLower(words[0]) == "deletecollection" {
				return db.DeleteCollection(settings, name, poolName, schemaName)
			}

			fmt.Println("Wrong command")
			return "error"

		// Set and Update value
		case "set", "update":
			if len(words) < 5 {
				fmt.Println("Too few arguments")
				return "error"
			}

			key := words[1]
			value := words[2 : len(words)-2]
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

			if strings.ToLower(words[0]) == "set" {
				return db.Set(settings, key, value, poolName, schemaName, collectionName)
			} else if strings.ToLower(words[0]) == "update" {
				return db.Update(settings, key, value, poolName, schemaName, collectionName)
			}

			fmt.Println("Wrong command")
			return "error"

		// Get and Delete value
		case "get", "delete":
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
				value, ok := db.Get(settings, key, poolName, schemaName, collectionName)
				if ok == "ok" {
					fmt.Printf("%v = %v\n", key, value)
				}
				return ok
			} else if strings.ToLower(words[0]) == "delete" {
				return db.Delete(settings, key, poolName, schemaName, collectionName)
			}

			fmt.Println("Wrong command")
			return "error"

		// Get range
		case "getrange":
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

			result, ok := db.GetRange(settings, leftBound, rightBound, poolName, schemaName, collectionName)
			if ok == "ok" {
				for k, v := range result {
					fmt.Printf("%v = %v\n", k, v)
				}
			}
			return ok
		}
	}

	fmt.Println("Something went wrong")
	return "error"
}

func executeFile(db database, settings map[string]string, filePath string) string {
	file, err := os.Open(filePath)
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			command := scanner.Text()
			ret := executeCommand(db, settings, command)
			if ret == "exit" {
				return "exit"
			} else if ret == "error" {
				return "error"
			}
		}
	} else {
		fmt.Printf("Could not open file. %v\n", err)
		return "error"
	}

	return "ok"
}

func help() {
	fmt.Println("Commands:")
	fmt.Println("\t'stop', 'exit' or '-1' to stop the program")
	fmt.Println("\t'help' to show this help")
}

func main() {
	args := os.Args[1:]

	settings := make(map[string]string)
	settings["DBtype"] = "BI"

	db := task0.CreateDB()

	// if there's a filename
	if len(args) == 1 {
		ret := executeFile(db, settings, args[0])
		if ret == "error" {
			fmt.Println("Could not execute file")
		} else if ret == "exit" {
			return
		}
	} else if len(args) > 1 {
		fmt.Println("Too many arguments, only one is allowed")
	}

	help()
	line := ""

	for line != "exit" {
		fmt.Println("Enter a command: ")
		fmt.Scanf("%s", &line)
		ret := executeCommand(db, settings, line)

		if ret == "exit" {
			line = "exit"
		}
	}
}