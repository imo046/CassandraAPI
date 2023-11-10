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
	r.Home()
	r.Create()
	r.GetOne()
	r.GetAll()
	r.Count()
	r.Delete()
	r.DeleteAll()
	r.Update()

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	if serveErr := http.ListenAndServe(":4040", handlers.CORS(headers, methods, origins)(r.Router)); serveErr != nil {
		utils.Panic(serveErr, "Failed to run a backend")
	}

}