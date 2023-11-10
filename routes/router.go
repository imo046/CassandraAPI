package routes

import (
	"CassandraAPI/api_handlers"
	"github.com/gorilla/mux"
)

type ApiRouter struct {
	Router  *mux.Router
	Handler *api_handlers.Handler
}

func (r *ApiRouter) Home() {
	r.Router.HandleFunc("/", r.Handler.HomeLink)
}

func (r *ApiRouter) Create() {
	r.Router.HandleFunc("/create", r.Handler.CreateEntry).Methods("POST")
}

func (r *ApiRouter) Count() {
	r.Router.HandleFunc("/count", r.Handler.CountEntries).Methods("GET")
}

func (r *ApiRouter) GetOne() {
	r.Router.HandleFunc("/get/{entry_id}", r.Handler.GetEntry).Methods("GET")
}

func (r *ApiRouter) GetAll() {
	r.Router.HandleFunc("/getAll", r.Handler.GetEntries).Methods("GET")
}

func (r *ApiRouter) Delete() {
	r.Router.HandleFunc("/delete/{entry_id}", r.Handler.DeleteOne).Methods("DELETE")
}

func (r *ApiRouter) DeleteAll() {
	r.Router.HandleFunc("/deleteAll", r.Handler.DeleteAll).Methods("DELETE")
}

func (r *ApiRouter) Update() {
	r.Router.HandleFunc("/update/{entry_id}", r.Handler.Update).Methods("PATCH")
}
