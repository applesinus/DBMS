package main

import (
	"Redis_Labs/task0"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func executeCommand(command string) string {
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
				ret := executeFile(words[0])
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
				return task0.CreatePool(name)
			} else if strings.ToLower(words[0]) == "deletepool" {
				return task0.DeletePool(name)
			}

			fmt.Println("Wrong command")
			return "error"

		// Create and Delete Scheme
		case "createscheme", "deletescheme":
			if len(words) > 4 {
				fmt.Println("Too many arguments")
				return "error"
			}

			name := words[1]
			if words[2] != "in" {
				fmt.Println("Wrong command on scheme name declaration")
				return "error"
			}
			poolName := words[3]

			if strings.ToLower(words[0]) == "deletescheme" {
				return task0.DeleteScheme(name, poolName)
			} else if strings.ToLower(words[0]) == "createscheme" {
				return task0.CreateScheme(name, poolName)
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
				fmt.Println("Wrong command on scheme name declaration")
				return "error"
			}
			schemeAndPool := strings.Split(words[3], ".")
			if len(schemeAndPool) != 2 {
				fmt.Println("Wrong command on scheme name declaration")
				return "error"
			}
			poolName := schemeAndPool[0]
			schemeName := schemeAndPool[1]

			if strings.ToLower(words[0]) == "createcollection" {
				return task0.CreateCollection(name, poolName, schemeName)
			} else if strings.ToLower(words[0]) == "deletecollection" {
				return task0.DeleteCollection(name, poolName, schemeName)
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
			poolAndSchemeAndCollection := strings.Split(words[len(words)-1], ".")
			if len(poolAndSchemeAndCollection) != 3 {
				fmt.Println("Wrong command on collection name declaration")
				return "error"
			}
			poolName := poolAndSchemeAndCollection[0]
			schemeName := poolAndSchemeAndCollection[1]
			collectionName := poolAndSchemeAndCollection[2]

			if strings.ToLower(words[0]) == "set" {
				return task0.Set(key, value, poolName, schemeName, collectionName)
			} else if strings.ToLower(words[0]) == "update" {
				return task0.Update(key, value, poolName, schemeName, collectionName)
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
			poolAndSchemeAndCollection := strings.Split(words[len(words)-1], ".")
			if len(poolAndSchemeAndCollection) != 3 {
				fmt.Println("Wrong command on collection name declaration")
				return "error"
			}
			poolName := poolAndSchemeAndCollection[0]
			schemeName := poolAndSchemeAndCollection[1]
			collectionName := poolAndSchemeAndCollection[2]

			if strings.ToLower(words[0]) == "get" {
				return task0.Get(key, poolName, schemeName, collectionName)
			} else if strings.ToLower(words[0]) == "delete" {
				return task0.Delete(key, poolName, schemeName, collectionName)
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
			poolAndSchemeAndCollection := strings.Split(words[len(words)-1], ".")
			if len(poolAndSchemeAndCollection) != 3 {
				fmt.Println("Wrong command on collection name declaration")
				return "error"
			}
			poolName := poolAndSchemeAndCollection[0]
			schemeName := poolAndSchemeAndCollection[1]
			collectionName := poolAndSchemeAndCollection[2]

			return task0.GetRange(leftBound, rightBound, poolName, schemeName, collectionName)
		}
	}

	fmt.Println("Something went wrong")
	return "error"
}

func executeFile(filePath string) string {
	file, err := os.Open(filePath)
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			command := scanner.Text()
			ret := executeCommand(command)
			if ret == "exit" {
				return "ok"
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

	// if there's a filename
	if len(args) == 1 {
		ret := executeFile(args[0])
		if ret == "error" {
			fmt.Println("Could not execute file")
		} else if ret == "exit" {
			return
		}
	} else if len(args) > 1 {
		fmt.Println("Too many arguments")
	}

	help()
	line := ""

	for line != "exit" {
		fmt.Scanf("%s", &line)
		ret := executeCommand(line)

		if ret == "exit" {
			line = "exit"
		}
	}
}
