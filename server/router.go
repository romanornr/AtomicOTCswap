package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)


func createRouter() *mux.Router {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/audit/{coin}", AuditHandler).Methods("GET")
	http.Handle("/", r)
	return r
}

func AuditHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	coin := params["coin"]
	fmt.Println(coin)
}