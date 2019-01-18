package server

import (
	"fmt"
	"log"
	"net/http"
)

func Start() {

	host := "127.0.0.1:8000"
	fmt.Printf("HTTP server started at %s\n", host)

	router := createRouter()
	err := http.ListenAndServe(host, router)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
