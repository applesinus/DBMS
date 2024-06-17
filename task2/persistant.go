package task2

import (
	"DBMS/interfaces"
	"DBMS/task0"
	"fmt"
	"strings"
	"time"
)

type collection struct {
	firstState task0.CollectionInterface
	commands   []string
	timings    []time.Time
}

var PersistantCollections = make(map[string]collection)

func SaveCollection(collName string, collElements task0.CollectionInterface) {
	PersistantCollections[collName] = collection{
		firstState: collElements,
		commands:   make([]string, 0),
		timings:    make([]time.Time, 0),
	}
}

func DeleteCollection(collName string) {
	delete(PersistantCollections, collName)
}

func AddCommand(collName string, command string) {
	coll := PersistantCollections[collName]

	coll.commands = append(coll.commands, command)
	coll.timings = append(coll.timings, time.Now())

	PersistantCollections[collName] = coll
}

func jumpToState(coll collection, index int) interfaces.CollectionInterface {
	temp := coll.firstState.

	for i := 0; i < index; i++ {
		words := strings.Fields(coll.commands[i])
		switch strings.ToLower(words[0]) {
		// Set and Update value
		case "set", "update":
			key := words[1]
			value := strings.Join(words[2:len(words)-2], " ")

			if strings.ToLower(words[0]) == "set" {
				temp.Set(key, value)
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
				value := db.Get(settings, key, poolName, schemaName, collectionName)
				if value != "" {
					fmt.Printf("%v = %v\n", key, value)
					return "ok"
				}
				return "error"
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

			result := db.GetRange(settings, leftBound, rightBound, poolName, schemaName, collectionName)
			if result != nil && len(*result) != 0 {
				for k, v := range *result {
					fmt.Printf("%v = %v\n", k, v)
				}
				return "ok"
			}
			return "error"
		}
	}
}

func GetValueByTime(collName string, key string, timeStamp time.Time) string {
	coll := PersistantCollections[collName]
	temp := coll.firstState

	index := 0
	for ; index < len(coll.timings); index++ {
		if coll.timings[index].After(timeStamp) {
			break
		}
	}
	if index != 0 {
		index--
	}

}
