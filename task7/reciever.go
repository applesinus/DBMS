package task7

import (
	"DBMS/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type recieved struct {
	Operation      string `json:"operation"`
	Key            string `json:"key"`
	Secondary      string `json:"secondaryKey"`
	Value          string `json:"value"`
	Pool           string `json:"pool"`
	Schema         string `json:"schema"`
	Collection     string `json:"collection"`
	Time           string `json:"time"`
	CollectionType string `json:"collectionType"`
}

func receiver(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received data: %v\n", r.Body)
	var recieved recieved

	// Parse the JSON request body
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&recieved)
	if err != nil {
		fmt.Printf("Error: %v\nData: %v\n", err, recieved)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Received data: %v\n", recieved)

	response := ""

	switch recieved.Operation {
	case "createPool", "deletePool":
		if recieved.Pool == "" {
			http.Error(w, "Pool name is required", http.StatusBadRequest)
			return
		}
		response = database.ExecuteCommand(recieved.Operation + " " + recieved.Pool)

	case "createSchema":
		if recieved.Schema == "" {
			http.Error(w, "Schema name is required", http.StatusBadRequest)
			return
		}
		if recieved.Pool == "" {
			http.Error(w, "Pool name is required", http.StatusBadRequest)
			return
		}
		response = database.ExecuteCommand(recieved.Operation + " " + recieved.Schema + " in " + recieved.Pool)

	case "deleteSchema":
		if recieved.Schema == "" {
			http.Error(w, "Schema name is required", http.StatusBadRequest)
			return
		}
		sc := strings.Split(recieved.Schema, ".")
		response = database.ExecuteCommand(recieved.Operation + " " + sc[1] + " in " + sc[0])

	case "createCollection":
		if recieved.Collection == "" {
			http.Error(w, "Collection name is required", http.StatusBadRequest)
			return
		}
		if recieved.Schema == "" {
			http.Error(w, "Pool.Schema name is required", http.StatusBadRequest)
			return
		}
		if recieved.CollectionType == "" {
			http.Error(w, "Collection type is required", http.StatusBadRequest)
			return
		}
		response = database.ExecuteCommand(recieved.Operation + " " + recieved.CollectionType + " " + recieved.Collection + " in " + recieved.Schema)

	case "deleteCollection":
		if recieved.Collection == "" {
			http.Error(w, "Collection name is required", http.StatusBadRequest)
			return
		}
		col := strings.Split(recieved.Collection, ".")
		response = database.ExecuteCommand(recieved.Operation + " " + col[2] + " in " + col[1] + "." + col[0])

	case "set":
		if recieved.Key == "" {
			http.Error(w, "Key is required", http.StatusBadRequest)
			return
		}
		if recieved.Secondary == "" {
			http.Error(w, "Secondary key is required", http.StatusBadRequest)
			return
		}
		if recieved.Value == "" {
			http.Error(w, "Value is required", http.StatusBadRequest)
			return
		}
		if recieved.Collection == "" {
			http.Error(w, "Pool.Schema.Collection name is required", http.StatusBadRequest)
			return
		}
		response = database.ExecuteCommand(recieved.Operation + " " + recieved.Key + " " + recieved.Secondary + " " + recieved.Value + " in " + recieved.Collection)

	case "update":
		if recieved.Key == "" {
			http.Error(w, "Key is required", http.StatusBadRequest)
			return
		}
		if recieved.Value == "" {
			http.Error(w, "Value is required", http.StatusBadRequest)
			return
		}
		key := strings.Split(recieved.Key, ".")
		response = database.ExecuteCommand(recieved.Operation + " " + key[3] + " " + recieved.Value + " in " + key[0] + "." + key[1] + "." + key[2])

	case "delete":
		if recieved.Key == "" {
			http.Error(w, "Key is required", http.StatusBadRequest)
			return
		}
		key := strings.Split(recieved.Key, ".")
		response = database.ExecuteCommand(recieved.Operation + " " + key[3] + " in " + key[0] + "." + key[1] + "." + key[2])

	case "get":
		if recieved.Key == "" {
			http.Error(w, "Key is required", http.StatusBadRequest)
			return
		}
		key := strings.Split(recieved.Key, ".")
		response = database.ExecuteCommand(recieved.Operation + " " + key[3] + " in " + key[0] + "." + key[1] + "." + key[2])

	case "getSecondary":
		if recieved.Secondary == "" {
			http.Error(w, "Secondary key is required", http.StatusBadRequest)
			return
		}
		key := strings.Split(recieved.Secondary, ".")
		response = database.ExecuteCommand(recieved.Operation + " " + key[3] + " in " + key[0] + "." + key[1] + "." + key[2])

	case "getRange", "getRangeSecondary":
		if recieved.Key == "" || recieved.Secondary == "" {
			http.Error(w, "Left and right bounds are required", http.StatusBadRequest)
			return
		}
		if recieved.Collection == "" {
			http.Error(w, "Pool.Schema.Collection name is required", http.StatusBadRequest)
			return
		}
		response = database.ExecuteCommand(recieved.Operation + " " + recieved.Key + " " + recieved.Secondary + " in " + recieved.Collection)

	case "getAt":
		if recieved.Key == "" {
			http.Error(w, "Key is required", http.StatusBadRequest)
			return
		}
		if recieved.Time == "" {
			http.Error(w, "Time is required", http.StatusBadRequest)
			return
		}
		key := strings.Split(recieved.Key, ".")
		fTime := ""
		// TODO convert time format to time.Time
		response = database.ExecuteCommand(recieved.Operation + " " + fTime + " " + key[3] + " in " + key[0] + "." + key[1] + "." + key[2])

	default:
		http.Error(w, "Invalid operation", http.StatusBadRequest)
		return
	}

	fmt.Printf("Response: %v\n", response)

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": response})
}
