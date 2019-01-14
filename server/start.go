package server

import (
	"fmt"
	"log"
	"net/http"
)

func Start() {
	fmt.Println("HTTP server started...")

	router := createRouter()
	err := http.ListenAndServe("127.0.0.1:8000", router)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
