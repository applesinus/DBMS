package task789

import (
	"DBMS/database"
	"DBMS/task0+3"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

type data struct {
	Pools           []pool
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

	initHandlers()
	loadLogins()

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

	db := data{
		RecieverURL:     "http://localhost:8080/receiver",
		CollectionTypes: database.AvaliableCollectionTypes(),
		Pools:           []pool{},
		User:            username,
	}
	pools := database.ListPools()
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
		db.Pools = append(db.Pools, newPool)
	}

	t, _ := template.ParseFiles("web/template.html", "web/blocks_user.html", "web/mainSU.html")
	err := t.Execute(w, db)
	if err != nil {
		println(err.Error())
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
		println(err.Error())
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
		println(err.Error())
	}
}

func logoutPage(w http.ResponseWriter, r *http.Request) {
	logout(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
