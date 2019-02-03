package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"github.com/gorilla/handlers"
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
