// Copyright (c) 2019 Romano (Viacoin developer)
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"github.com/gorilla/handlers"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("website/*"))
}

func Start() {

	host := "127.0.0.1:8000"
	fmt.Printf("HTTP server started at %s\n", host)

	router := createRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	err := http.ListenAndServe(host, handlers.CORS(headersOk)(router))

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
