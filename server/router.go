package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/romanornr/AtomicOTCswap/atomic"
	"net/http"
)


func createRouter() *mux.Router {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/audit/{coin}/{contractHex}/{contractTransaction}", AuditHandler).Methods("GET")
	http.Handle("/", r)
	return r
}

func AuditHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	coin, contractHex, contractTransaction := params["coin"], params["contractHex"], params["contractTransaction"]
	contract, _ := atomic.AuditContract(coin, contractHex, contractTransaction)
	fmt.Println(contract)
		json.NewEncoder(w).Encode(&contract)
}