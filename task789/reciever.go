package task789

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
	ok, username, _ := isLoggedIn(w, r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fmt.Printf("Received data: %v\n", r.Body)
	var recieved recieved

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&recieved)
	if err != nil {
		fmt.Printf("Error: %v\nData: %v\n", err, recieved)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, ok := users.List[username]
	if !ok {
		http.Error(w, "not a user", http.StatusBadRequest)
		return
	}

	//fmt.Printf("Received data: %v\n", recieved)

	response := ""

	switch recieved.Operation {
	case "createPool", "deletePool":
		if u.Role != "admin" && u.Role != "superuser" && u.Role != "editor" {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		} else if u.Role == "editor" {
			ok, res := u.checkCollectionAccess(recieved.Collection, "w")
			if res != "ok" {
				http.Error(w, res, http.StatusBadRequest)
				return
			}
			if !ok {
				http.Error(w, "Access denied", http.StatusBadRequest)
				return
			}
		}

		if recieved.Pool == "" {
			http.Error(w, "Pool name is required", http.StatusBadRequest)
			return
		}
		response = database.ExecuteCommand(recieved.Operation + " " + recieved.Pool)

	case "createSchema":
		if u.Role != "admin" && u.Role != "superuser" && u.Role != "editor" {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		} else if u.Role == "editor" {
			ok, res := u.checkCollectionAccess(recieved.Collection, "w")
			if res != "ok" {
				http.Error(w, res, http.StatusBadRequest)
				return
			}
			if !ok {
				http.Error(w, "Access denied", http.StatusBadRequest)
				return
			}
		}

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
		if u.Role != "admin" && u.Role != "superuser" && u.Role != "editor" {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		} else if u.Role == "editor" {
			ok, res := u.checkCollectionAccess(recieved.Collection, "w")
			if res != "ok" {
				http.Error(w, res, http.StatusBadRequest)
				return
			}
			if !ok {
				http.Error(w, "Access denied", http.StatusBadRequest)
				return
			}
		}

		if recieved.Schema == "" {
			http.Error(w, "Schema name is required", http.StatusBadRequest)
			return
		}
		sc := strings.Split(recieved.Schema, ".")
		response = database.ExecuteCommand(recieved.Operation + " " + sc[1] + " in " + sc[0])

	case "createCollection":
		if u.Role != "admin" && u.Role != "superuser" && u.Role != "editor" {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		} else if u.Role == "editor" {
			ok, res := u.checkCollectionAccess(recieved.Collection, "w")
			if res != "ok" {
				http.Error(w, res, http.StatusBadRequest)
				return
			}
			if !ok {
				http.Error(w, "Access denied", http.StatusBadRequest)
				return
			}
		}

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
		if u.Role != "admin" && u.Role != "superuser" && u.Role != "editor" {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		} else if u.Role == "editor" {
			ok, res := u.checkCollectionAccess(recieved.Collection, "w")
			if res != "ok" {
				http.Error(w, res, http.StatusBadRequest)
				return
			}
			if !ok {
				http.Error(w, "Access denied", http.StatusBadRequest)
				return
			}
		}

		if recieved.Collection == "" {
			http.Error(w, "Collection name is required", http.StatusBadRequest)
			return
		}
		col := strings.Split(recieved.Collection, ".")
		response = database.ExecuteCommand(recieved.Operation + " " + col[2] + " in " + col[0] + "." + col[1])

	case "set":
		if u.Role != "admin" && u.Role != "superuser" && u.Role != "editor" && u.Role != "user" {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		} else if u.Role == "editor" || u.Role == "user" {
			ok, res := u.checkCollectionAccess(recieved.Collection, "w")
			if res != "ok" {
				http.Error(w, res, http.StatusBadRequest)
				return
			}
			if !ok {
				http.Error(w, "Access denied", http.StatusBadRequest)
				return
			}
		}

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
		if u.Role != "admin" && u.Role != "superuser" && u.Role != "editor" && u.Role != "user" {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		} else if u.Role == "editor" || u.Role == "user" {
			ok, res := u.checkCollectionAccess(recieved.Collection, "w")
			if res != "ok" {
				http.Error(w, res, http.StatusBadRequest)
				return
			}
			if !ok {
				http.Error(w, "Access denied", http.StatusBadRequest)
				return
			}
		}

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
		if u.Role != "admin" && u.Role != "superuser" && u.Role != "editor" && u.Role != "user" {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		} else if u.Role == "editor" || u.Role == "user" {
			ok, res := u.checkCollectionAccess(recieved.Collection, "w")
			if res != "ok" {
				http.Error(w, res, http.StatusBadRequest)
				return
			}
			if !ok {
				http.Error(w, "Access denied", http.StatusBadRequest)
				return
			}
		}

		if recieved.Key == "" {
			http.Error(w, "Key is required", http.StatusBadRequest)
			return
		}
		key := strings.Split(recieved.Key, ".")
		response = database.ExecuteCommand(recieved.Operation + " " + key[3] + " in " + key[0] + "." + key[1] + "." + key[2])

	case "get":
		if u.Role != "admin" && u.Role != "superuser" && u.Role != "editor" && u.Role != "user" {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		} else if u.Role == "editor" || u.Role == "user" {
			ok, res := u.checkCollectionAccess(recieved.Collection, "r")
			if res != "ok" {
				http.Error(w, res, http.StatusBadRequest)
				return
			}
			if !ok {
				http.Error(w, "Access denied", http.StatusBadRequest)
				return
			}
		}

		if recieved.Key == "" {
			http.Error(w, "Key is required", http.StatusBadRequest)
			return
		}
		key := strings.Split(recieved.Key, ".")
		response = database.ExecuteCommand(recieved.Operation + " " + key[3] + " in " + key[0] + "." + key[1] + "." + key[2])

	case "getSecondary":
		if u.Role != "admin" && u.Role != "superuser" && u.Role != "editor" && u.Role != "user" {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		} else if u.Role == "editor" || u.Role == "user" {
			ok, res := u.checkCollectionAccess(recieved.Collection, "r")
			if res != "ok" {
				http.Error(w, res, http.StatusBadRequest)
				return
			}
			if !ok {
				http.Error(w, "Access denied", http.StatusBadRequest)
				return
			}
		}

		if recieved.Secondary == "" {
			http.Error(w, "Secondary key is required", http.StatusBadRequest)
			return
		}
		key := strings.Split(recieved.Secondary, ".")
		response = database.ExecuteCommand(recieved.Operation + " " + key[3] + " in " + key[0] + "." + key[1] + "." + key[2])

	case "getRange", "getRangeSecondary":
		if u.Role != "admin" && u.Role != "superuser" && u.Role != "editor" && u.Role != "user" {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		} else if u.Role == "editor" || u.Role == "user" {
			ok, res := u.checkCollectionAccess(recieved.Collection, "r")
			if res != "ok" {
				http.Error(w, res, http.StatusBadRequest)
				return
			}
			if !ok {
				http.Error(w, "Access denied", http.StatusBadRequest)
				return
			}
		}

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
		if u.Role != "admin" && u.Role != "superuser" && u.Role != "editor" && u.Role != "user" {
			http.Error(w, "Permission denied", http.StatusForbidden)
			return
		} else if u.Role == "editor" || u.Role == "user" {
			ok, res := u.checkCollectionAccess(recieved.Collection, "r")
			if res != "ok" {
				http.Error(w, res, http.StatusBadRequest)
				return
			}
			if !ok {
				http.Error(w, "Access denied", http.StatusBadRequest)
				return
			}
		}

		if recieved.Key == "" {
			http.Error(w, "Key is required", http.StatusBadRequest)
			return
		}
		if recieved.Time == "" {
			http.Error(w, "Time is required", http.StatusBadRequest)
			return
		}
		key := strings.Split(recieved.Key, ".")
		fTime := recieved.Time[:10] + " " + recieved.Time[11:] + " MSK"

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
