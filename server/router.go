package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/romanornr/AtomicOTCswap/atomic"
	btcutil "github.com/viacoin/viautil"
	"log"
	"net/http"
	"strconv"
)


func createRouter() *mux.Router {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/audit/{coin}/{contractHex}/{contractTransaction}", AuditHandler).Methods("GET")
	api.HandleFunc("/initiate", InitiateHandler).Methods("POST")
	api.HandleFunc("/participate", ParticipateHandler).Methods("POST")
	api.HandleFunc("/redeem", RedemptionHandler).Methods("POST")
	api.HandleFunc("/extractsecret", SecretHandler).Methods("POST")
	http.Handle("/", r)
	return r
}

// initiate a contract by parsing the post request
// it parses the coin symbol, counter party address, amount and the wif
func InitiateHandler(w http.ResponseWriter, req *http.Request) {

	amount, err := strconv.ParseFloat(req.FormValue("amount"), 64)
	if err != nil {
		log.Printf("amount should be a float. example: 0.02")
	}

	wif, err := btcutil.DecodeWIF(req.FormValue("wif"))
	if err != nil {
		log.Printf("error decoding private key in wif format: %s\n", err)
	}

	contract, err := atomic.Initiate(req.FormValue("coin"), req.FormValue("counter_party_2_addr"), amount, wif)
	if err != nil {
		log.Printf("erorr initiating contract: %s\n", err)
	}

	json.NewEncoder(w).Encode(contract)
}

func ParticipateHandler(w http.ResponseWriter, req *http.Request) {

	amount, err := strconv.ParseFloat(req.FormValue("amount"), 64)
	if err != nil {
		log.Printf("amount should be a float. example: 0.02")
	}

	wif, err := btcutil.DecodeWIF(req.FormValue("wif"))
	if err != nil {
		log.Printf("error decoding private key in wif format: %s\n", err)
	}

	contract, err := atomic.Participate(req.FormValue("coin"), req.FormValue("initiatorAddress"), wif, amount, req.FormValue("secret"))
	if err != nil {
		log.Printf("error participating contract: %s\n", err)
	}

	json.NewEncoder(w).Encode(contract)
}

func RedemptionHandler(w http.ResponseWriter, req *http.Request) {

	wif, err := btcutil.DecodeWIF(req.FormValue("wif"))
	if err != nil {
		log.Printf("error decoding private key in wif format: %s\n", err)
	}

	redemption, err := atomic.Redeem(req.FormValue("coin"), req.FormValue("contractHex"), req.FormValue("contractTransaction"), req.FormValue("secretHex"), wif)
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
	params := mux.Vars(req)
	coin, contractHex, contractTransaction := params["coin"], params["contractHex"], params["contractTransaction"]
	contract, err := atomic.AuditContract(coin, contractHex, contractTransaction)
	if err != nil {
		fmt.Sprintf("%s\n", err)
	}
	json.NewEncoder(w).Encode(&contract)
}