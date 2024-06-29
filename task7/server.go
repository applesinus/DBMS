package task7

import (
	"DBMS/database"
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
	Datas []valuesData
}

type valuesData struct {
	Key          string
	SecondaryKey string
	Value        string
}

func initHandlers() {
	http.HandleFunc("/", mainPage)
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

	return server
}

func StopHTTPServer(server *http.Server) {
	fmt.Println("Stopping HTTP Server")
	err := server.Shutdown(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
}

func testDB() data {
	return data{
		Pools: []pool{
			{
				Name: "pool1",
				Schemas: []schema{
					{
						Name: "schema1",
						Collections: []collection{
							{
								Name: "collection1",
								Datas: []valuesData{
									{
										Key:   "key1",
										Value: "value1",
									},
									{
										Key:   "key2",
										Value: "value2",
									},
								},
							},
							{
								Name: "collection2",
								Datas: []valuesData{
									{
										Key:   "key1",
										Value: "value1",
									},
									{
										Key:   "key2",
										Value: "value2",
									},
								},
							},
						},
					},
					{
						Name: "schema2",
						Collections: []collection{
							{
								Name: "collection1",
								Datas: []valuesData{
									{
										Key:   "key1",
										Value: "value1",
									},
									{
										Key:   "key2",
										Value: "value2",
									},
								},
							},
							{
								Name: "collection2",
								Datas: []valuesData{
									{
										Key:   "key1",
										Value: "value1",
									},
									{
										Key:   "key2",
										Value: "value2",
									},
								},
							},
						},
					},
				},
			},
			{
				Name: "pool2",
				Schemas: []schema{
					{
						Name: "schema1",
						Collections: []collection{
							{
								Name: "collection1",
								Datas: []valuesData{
									{
										Key:   "key1",
										Value: "value1",
									},
									{
										Key:   "key2",
										Value: "value2",
									},
								},
							},
							{
								Name: "collection2",
								Datas: []valuesData{
									{
										Key:   "key1",
										Value: "value1",
									},
									{
										Key:   "key2",
										Value: "value2",
									},
								},
							},
						},
					},
					{
						Name: "schema2",
						Collections: []collection{
							{
								Name: "collection1",
								Datas: []valuesData{
									{
										Key:   "key1",
										Value: "value1",
									},
									{
										Key:   "key2",
										Value: "value2",
									},
								},
							},
							{
								Name: "collection2",
								Datas: []valuesData{
									{
										Key:   "key1",
										Value: "value1",
									},
									{
										Key:   "key2",
										Value: "value2",
									},
								},
							},
						},
					},
				},
			},
		},
		RecieverURL: "http://localhost:8080/receiver",
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	//db := testDB()
	db := data{
		RecieverURL:     "http://localhost:8080/receiver",
		CollectionTypes: database.AvaliableCollectionTypes(),
		Pools:           []pool{},
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
					newData := valuesData{Key: d.Key, SecondaryKey: d.SecondaryKey, Value: d.Value}
					newCollection.Datas = append(newCollection.Datas, newData)
				}
				newSchema.Collections = append(newSchema.Collections, newCollection)
			}
			newPool.Schemas = append(newPool.Schemas, newSchema)
		}
		db.Pools = append(db.Pools, newPool)
	}

	t, _ := template.ParseFiles("web/template.html", "web/mainSU.html")
	err := t.Execute(w, db)
	if err != nil {
		println(err.Error())
	}
}
