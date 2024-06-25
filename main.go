package main

import (
	"DBMS/task0"
	"DBMS/task2"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func executeCommand(db *task0.Database, settings map[string]string, command string) string {
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
		case "set", "update":
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

			if strings.ToLower(words[0]) == "set" {
				response := db.Set(settings, key, value, poolName, schemaName, collectionName)
				if response != "ok" {
					return response
				}
				if settings["persistant"] == "on" {
					task2.AddCommand(collectionName, command)
				}
				return response
			} else if strings.ToLower(words[0]) == "update" {
				response := db.Update(settings, key, value, poolName, schemaName, collectionName)
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
				value := db.Get(settings, key, poolName, schemaName, collectionName)
				if value != "" {
					fmt.Printf("%v = %v\n", key, value)
					return "ok"
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

			result := db.GetRange(settings, leftBound, rightBound, poolName, schemaName, collectionName)
			if result != nil && len(*result) != 0 {
				for k, v := range *result {
					fmt.Printf("%v = %v\n", k, v)
				}
				return "ok"
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
				fmt.Printf("%v = %v\n", key, value)
				return "ok"
			}
			return "error"
		}
	}

	fmt.Println("Something went wrong")
	return "error"
}

func executeFile(db *task0.Database, settings map[string]string, filePath string) string {
	file, err := os.Open(filePath)
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			command := scanner.Text()
			ret := executeCommand(db, settings, command)
			fmt.Printf("command: %v, ret: %v\n", command, ret)
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
	settings["persistant"] = "on"

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

	in := bufio.NewReader(os.Stdin)
	for line != "exit" {
		fmt.Printf("Enter a command: ")
		line, err := in.ReadString('\n')

		if err != nil {
			fmt.Println(err)
			return
		}

		ret := executeCommand(db, settings, line)

		if ret == "exit" {
			line = "exit"
		}
	}
}
