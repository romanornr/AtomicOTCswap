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

// audit a contract by giving the coin symbol, contract hex and contract transaction
// from the contract which needs to be audited
func AuditHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	coin, contractHex, contractTransaction := params["coin"], params["contractHex"], params["contractTransaction"]
	contract, err := atomic.AuditContract(coin, contractHex, contractTransaction)
	if err != nil {
		fmt.Sprintf("%s\n", err)
	}
	json.NewEncoder(w).Encode(&contract)
}