package task789

import (
	"DBMS/database"
	"DBMS/task0+3"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

type data struct {
	Wpools          []pool
	Rpools          []pool
	RecieverURL     string
	CollectionTypes []string
	User            string
	Message         string
}

type pool struct {
	Name    string
	Schemas []schema
}

type schema struct {
	Name        string
	Collections []collection
}

type collection struct {
	Name  string
	Datas []ValuesData
}

type ValuesData task0.Datas

func initHandlers() {
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/admin", adminPage)

	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/register", registerPage)
	http.HandleFunc("/logout", logoutPage)

	http.HandleFunc("/receiver", receiver)

	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
}

func StartHTTPServer(port int) *http.Server {
	fmt.Println("Starting HTTP Server")
	fmt.Println("Listening on port", port)

	server := &http.Server{
		Addr: "localhost:" + strconv.Itoa(port),
	}
	serverIsRunning := make(chan bool)
	go func(serverIsRunning chan bool) {
		fmt.Println("\nSERVER IS RUNNING!")
		serverIsRunning <- true
		err := server.ListenAndServe()
		if err != nil {
			fmt.Println("\nTHE SERVER HAS BEEN STOPPED!\n" + err.Error())
		}
	}(serverIsRunning)
	<-serverIsRunning
	close(serverIsRunning)

	loadLogins()
	initHandlers()

	return server
}

func StopHTTPServer(server *http.Server) {
	fmt.Println("Stopping HTTP Server")
	err := server.Shutdown(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	ok, username, res := isLoggedIn(w, r)
	if res != "ok" || !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	u := users.List[username]

	db := data{
		RecieverURL:     "http://localhost:8080/receiver",
		CollectionTypes: database.AvaliableCollectionTypes(),
		Wpools:          []pool{},
		Rpools:          []pool{},
		User:            username,
		Message:         u.Role,
	}

	pools := database.ListPools()
	for _, p := range pools {
		newrPool := pool{Name: p}
		newwPool := pool{Name: p}

		if len(u.access) == 1 && u.access[0] == "none" {
			db.Wpools = append(db.Wpools, newwPool)
			db.Rpools = append(db.Rpools, newrPool)
			continue
		}

		schemas := database.ListSchemas(p)
		for _, s := range schemas {
			newrSchema := schema{Name: s}
			newwSchema := schema{Name: s}
			collections := database.ListCollections(p, s)
			for _, c := range collections {
				newrCollection := collection{Name: c}
				newwCollection := collection{Name: c}
				datas := database.ListDatas(p, s, c)
				for _, d := range datas {
					newData := ValuesData{Key: d.Key, SecondaryKey: d.SecondaryKey, Value: d.Value}
					newrCollection.Datas = append(newrCollection.Datas, newData)
				}

				if len(u.access) == 1 && u.access[0] == "all" {
					newwSchema.Collections = append(newwSchema.Collections, newwCollection)
					newrSchema.Collections = append(newrSchema.Collections, newrCollection)
					continue
				}

				poolAccess := ""
				schemaAccess := ""
				collectionAccess := ""
				for _, access := range u.access {
					words := strings.Split(access, "/")
					mod := words[1]
					words = strings.Split(words[0], ".")

					if len(words) == 1 && words[0] == p {
						poolAccess = mod
					} else if len(words) == 2 && words[0] == p && words[1] == s {
						schemaAccess = mod
					} else if len(words) == 3 && words[0] == p && words[1] == s && words[2] == c {
						collectionAccess = mod
					}
				}

				if collectionAccess == "w" || (collectionAccess == "" && schemaAccess == "w") ||
					(collectionAccess == "" && schemaAccess == "" && poolAccess == "w") {
					newwSchema.Collections = append(newwSchema.Collections, newwCollection)
					newrSchema.Collections = append(newrSchema.Collections, newrCollection)
				}

				if collectionAccess == "r" || (collectionAccess == "" && schemaAccess == "r") ||
					(collectionAccess == "" && schemaAccess == "" && poolAccess == "r") {
					newrSchema.Collections = append(newrSchema.Collections, newrCollection)
				}
			}
			newwPool.Schemas = append(newwPool.Schemas, newwSchema)
			newrPool.Schemas = append(newrPool.Schemas, newrSchema)
		}
		db.Rpools = append(db.Rpools, newrPool)
		db.Wpools = append(db.Wpools, newwPool)
	}
	t, _ := template.ParseFiles("web/template.html", "web/blocks_user.html", "web/main.html")
	err := t.Execute(w, db)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	ok, _, _ := isLoggedIn(w, r)
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		username := r.FormValue("login")
		password := r.FormValue("password")
		ok := login(w, username, password)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	t, _ := template.ParseFiles("web/template.html", "web/blocks_notuser.html", "web/login.html")
	err := t.Execute(w, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func registerPage(w http.ResponseWriter, r *http.Request) {
	ok, _, _ := isLoggedIn(w, r)
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		username := r.FormValue("login")
		if len(username) < 5 {
			http.Redirect(w, r, "/register?error=too short login", http.StatusSeeOther)
			return
		}
		if len(username) > 15 {
			http.Redirect(w, r, "/register?error=too long login", http.StatusSeeOther)
			return
		}
		for _, ch := range username {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')) {
				http.Redirect(w, r, "/register?error=invalid login", http.StatusSeeOther)
				return
			}
		}

		password := r.FormValue("password")
		if len(password) < 8 {
			http.Redirect(w, r, "/register?error=too short password", http.StatusSeeOther)
			return
		}

		res := register(w, username, password)
		if res != "ok" {
			http.Redirect(w, r, "/register?error="+res, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	message := data{
		Message: "",
	}
	args := r.URL.Query()
	if args.Get("error") != "" {
		message.Message = args.Get("error")
	}

	t, _ := template.ParseFiles("web/template.html", "web/blocks_notuser.html", "web/register.html")
	err := t.Execute(w, message)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func logoutPage(w http.ResponseWriter, r *http.Request) {
	ok, _, _ := isLoggedIn(w, r)
	if ok {
		logout(w)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func adminPage(w http.ResponseWriter, r *http.Request) {
	ok, username, res := isLoggedIn(w, r)
	if res != "ok" || !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	u := users.List[username]
	if u.Role != "admin" && u.Role != "superuser" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		targetUsername := r.FormValue("login")
		user, ok := users.List[targetUsername]
		if !ok {
			fmt.Println("Unknown user: " + targetUsername)
			return
		}
		if user.Role == "superuser" {
			fmt.Println("Can't change superuser")
			return
		}
		if u.Role == "admin" && user.Role == "admin" {
			fmt.Println("One admin can't change another admin")
			return
		}

		op := r.FormValue("op")

		switch op {
		case "delete":
			deleteUser(targetUsername)
		case "password":
			res := updatePassword(w, targetUsername, r.FormValue("password"))
			if res != "ok" {
				fmt.Println(res)
			}
		case "role":
			res := updateRole(w, targetUsername, r.FormValue("role"))
			if res != "ok" {
				fmt.Println(res)
			}
		case "access":
			res := updateAccess(w, targetUsername, r.FormValue("access"), r.FormValue("read"), r.FormValue("write"))
			if res != "ok" {
				fmt.Println(res)
			}
		default:
			fmt.Println("unknown admin operation: " + op)
		}
	}

	pools := database.ListPools()
	poolsList := []pool{}
	for _, p := range pools {
		newPool := pool{Name: p}
		schemas := database.ListSchemas(p)
		for _, s := range schemas {
			newSchema := schema{Name: s}
			collections := database.ListCollections(p, s)
			for _, c := range collections {
				newCollection := collection{Name: c}
				datas := database.ListDatas(p, s, c)
				for _, d := range datas {
					newData := ValuesData{Key: d.Key, SecondaryKey: d.SecondaryKey, Value: d.Value}
					newCollection.Datas = append(newCollection.Datas, newData)
				}
				newSchema.Collections = append(newSchema.Collections, newCollection)
			}
			newPool.Schemas = append(newPool.Schemas, newSchema)
		}
		poolsList = append(poolsList, newPool)
	}

	somaData := struct {
		Users logins
		User  string
		Pools []pool
	}{
		Users: users,
		User:  username,
		Pools: poolsList,
	}

	t, _ := template.ParseFiles("web/template.html", "web/blocks_user.html", "web/admin.html")

	err := t.Execute(w, somaData)
	if err != nil {
		fmt.Println(err.Error())
	}
}
