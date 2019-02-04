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

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"result"`
	Error   string      `json:"error"`
}

func createRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/initiate", InitiateSiteHandler).Methods("GET")
	r.HandleFunc("/audit", AuditSiteHandler).Methods("GET")
	r.HandleFunc("/participate", participateSiteHandler).Methods("GET")
	r.HandleFunc("/redeem", RedemptionSiteHandler).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/audit", AuditHandler).Methods("POST")
	api.HandleFunc("/initiate", InitiateHandler).Methods("POST")
	api.HandleFunc("/participate", ParticipateHandler).Methods("POST")
	api.HandleFunc("/redeem", RedemptionHandler).Methods("POST")
	api.HandleFunc("/extractsecret", SecretHandler).Methods("POST")
	http.Handle("/", r)

	return r
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		fmt.Println("error template")
	}
}

func InitiateSiteHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "initiate.gohtml", nil)
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
	respond(w, contract, err)
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
		respond(w, nil, err)
	}

	contract, err := atomic.Participate(req.FormValue("coin"), req.FormValue("counterPartyAddr"), req.FormValue("wif"), amount, req.FormValue("secret"))
	respond(w, contract, err)
}

func RedemptionSiteHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "redeem.gohtml", nil)
	if err != nil {
		fmt.Println("error template")
	}
}

func RedemptionHandler(w http.ResponseWriter, req *http.Request) {
	redemption, err := atomic.Redeem(req.FormValue("coin"), req.FormValue("contractHex"), req.FormValue("contractTransaction"), req.FormValue("secret"), req.FormValue("wif"))
	respond(w, redemption, err)
}

func SecretHandler(w http.ResponseWriter, req *http.Request) {
	secret, err := atomic.ExtractSecret(req.FormValue("redemptionTransaction"), req.FormValue("secretHash"))
	respond(w, secret, err)
}

func AuditSiteHandler(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "audit.gohtml", nil)
	if err != nil {
		fmt.Println("error template")
	}
}

//type UserInput struct {
//	Coin                string `json:"coin"`
//	ContractHex         string `json:"contractHex"`
//	ContractTransaction string `json:"contractTransaction"`
//}

// audit a contract by giving the coin symbol, contract hex and contract transaction
// from the contract which needs to be audited
func AuditHandler(w http.ResponseWriter, req *http.Request) {
	//form := req.Form
	contract, err := atomic.AuditContract(req.FormValue("coin"), req.FormValue("contractHex"), req.FormValue("contractTransaction"))
	respond(w, contract, err)
}

func respond(w http.ResponseWriter, data interface{}, err error) {
	response := Response{Data: data, Success: true}
	if err != nil {
		response.Data = nil
		response.Success = false
		response.Error = err.Error()
		///w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return
}
