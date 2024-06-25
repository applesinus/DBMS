package task2

import (
	"DBMS/interfaces"
	"strings"
	"time"
)

type collection struct {
	firstState interfaces.CollectionInterface
	commands   []string
	timings    []time.Time
}

var PersistantCollections = make(map[string]collection)

func SaveCollection(collName string, collElements interfaces.CollectionInterface) {
	PersistantCollections[collName] = collection{
		firstState: collElements,
		commands:   make([]string, 0),
		timings:    make([]time.Time, 0),
	}
}

func DeleteCollection(collName string) {
	AddCommand(collName, "deleteCollection "+collName)
}

func AddCommand(collName string, command string) {
	coll := PersistantCollections[collName]

	coll.commands = append(coll.commands, command)
	coll.timings = append(coll.timings, time.Now())

	PersistantCollections[collName] = coll
}

func jumpToState(coll collection, index int) interfaces.CollectionInterface {
	temp := coll.firstState.Copy()

	for i := 0; i <= index; i++ {
		words := strings.Fields(coll.commands[i])
		switch strings.ToLower(words[0]) {
		// Set and Update value
		case "set", "update":
			key := words[1]
			secondaryKey := words[2]
			value := strings.Join(words[2:len(words)-2], " ")

			if strings.ToLower(words[0]) == "set" {
				temp.Set(key, secondaryKey, value)
			} else if strings.ToLower(words[0]) == "update" {
				temp.Update(key, value)
			}

		// Get and Delete value
		case "delete":
			key := words[1]
			temp.Delete(key)
		}
	}

	return temp
}

func GetValueByTime(collName string, key string, timeStamp time.Time) string {
	coll := PersistantCollections[collName]

	index := 0
	for ; index < len(coll.timings); index++ {
		if coll.timings[index].After(timeStamp) {
			break
		}
	}
	if index == 0 {
		return "collection was created at " + coll.timings[index].String()
	}

	index--

	if coll.commands[index] == "deleteCollection "+collName {
		return "collection has been deleted at " + coll.timings[index].String()
	}

	temp := jumpToState(coll, index)
	result, response := temp.Get(key)

	if response != "ok" {
		return key + "did not exist at " + timeStamp.String()
	}

	return result
}
