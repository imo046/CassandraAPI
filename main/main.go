package main

import (
	"CassandraAPI/api_handlers"
	"CassandraAPI/routes"
	"CassandraAPI/utils"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

var Session *gocql.Session

func main() {

	var err error
	cluster := gocql.NewCluster("127.0.0.1:9042")
	cluster.Keyspace = "myapi"
	Session, err = cluster.CreateSession()
	utils.Panic(err, "Failed to create session")
	fmt.Println("Cassandra successfully initialized")

	//Check(Session)
	defer Session.Close()

	myHandler := api_handlers.Handler{S: Session}

	r := routes.ApiRouter{
		Router:  mux.NewRouter().StrictSlash(true),
		Handler: &myHandler,
	}
	//routes
	r.Home("/")
	r.Create("/create")
	r.GetOne("/get/{entry_id")
	r.GetAll("/getAll")
	r.Count("/count")
	r.Delete("/delete/{entry_id}")
	r.DeleteAll("/deleteAll")
	r.Update("/update/{entry_id}")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	r.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	if serveErr := http.ListenAndServe(":4040", handlers.CORS(headers, methods, origins)(r.Router)); serveErr != nil {
		utils.Panic(serveErr, "Failed to run a backend")
	}

}
