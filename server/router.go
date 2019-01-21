package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/romanornr/AtomicOTCswap/atomic"
	"log"
	"net/http"
	"strconv"
)


func createRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/initiate", InitiateSiteHandler).Methods("GET")
	r.HandleFunc("/audit", AuditSiteHandler).Methods("GET")
	r.HandleFunc("/participate", participateSiteHandler).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	//api.HandleFunc("/audit/{asset}/{contractHex}/{contractTransaction}", AuditHandler).Methods("GET")
	api.HandleFunc("/audit", AuditHandler).Methods("POST")
	api.HandleFunc("/initiate", InitiateHandler).Methods("POST")
	api.HandleFunc("/participate", ParticipateHandler).Methods("POST")
	api.HandleFunc("/redeem", RedemptionHandler).Methods("POST")
	api.HandleFunc("/extractsecret", SecretHandler).Methods("POST")
	http.Handle("/", r)
	return r
}

func InitiateSiteHandler(w http.ResponseWriter, r *http.Request) {
	err :=  tpl.ExecuteTemplate(w, "initiate.gohtml", nil)
	if err != nil {
		fmt.Println("error template")
	}
}

// initiate a contract by parsing the post request
// it parses the coin symbol, counter party address, amount and the wif
func InitiateHandler(w http.ResponseWriter, req *http.Request) {

	amount, err := strconv.ParseFloat(req.FormValue("amount"), 64)
	if err != nil {
		log.Printf("amount should be a float. example: 0.02")
	}

	contract, err := atomic.Initiate(req.FormValue("coin"), req.FormValue("counterPartyAddr"), amount, req.FormValue("wif"))
	if err != nil {
		log.Printf("erorr initiating contract: %s\n", err)
	}

	json.NewEncoder(w).Encode(contract)
}

func participateSiteHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "participate.gohtml", nil)
	if err != nil {
		fmt.Println("error template")
	}
}

func ParticipateHandler(w http.ResponseWriter, req *http.Request) {

	amount, err := strconv.ParseFloat(req.FormValue("amount"), 64)
	if err != nil {
		log.Printf("amount should be a float. example: 0.02")
	}

	contract, err := atomic.Participate(req.FormValue("coin"), req.FormValue("counterPartyAddr"), req.FormValue("wif"), amount, req.FormValue("secret"))
	if err != nil {
		log.Printf("error participating contract: %s\n", err)
	}

	json.NewEncoder(w).Encode(contract)
}

func RedemptionHandler(w http.ResponseWriter, req *http.Request) {

	redemption, err := atomic.Redeem(req.FormValue("coin"), req.FormValue("contractHex"), req.FormValue("contractTransaction"), req.FormValue("secretHex"), req.FormValue("wif"))
	if err != nil {
		log.Printf("error redemption: %s\n", err)
	}

	json.NewEncoder(w).Encode(redemption)
}

func SecretHandler(w http.ResponseWriter, req *http.Request) {
	secret, err := atomic.ExtractSecret(req.FormValue("redemptionTransaction"), req.FormValue("secretHash"))
	if err != nil {
		log.Printf("error extracting secret: %s\n", err)
	}
	json.NewEncoder(w).Encode(secret)
}

// audit a contract by giving the coin symbol, contract hex and contract transaction
// from the contract which needs to be audited
func AuditHandler(w http.ResponseWriter, req *http.Request) {
	//params := mux.Vars(req)
	//coin, contractHex, contractTransaction := params["coin"], params["contractHex"], params["contractTransaction"]
	contract, err := atomic.AuditContract(req.FormValue("coin"), req.FormValue("contractHex"), req.FormValue("contractTransaction"))
	if err != nil {
		fmt.Sprintf("%s\n", err)
	}
	json.NewEncoder(w).Encode(&contract)
}

func AuditSiteHandler(w http.ResponseWriter, req *http.Request) {
	err :=  tpl.ExecuteTemplate(w, "audit.gohtml", nil)
	if err != nil {
		fmt.Println("error template")
	}
}