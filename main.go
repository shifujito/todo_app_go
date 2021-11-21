package main

import (
	"net/http"
	"text/template"
)

var tmpl = template.Must(template.ParseGlob("template/*"))

func main() {
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: nil,
	}
	http.HandleFunc("/", Index)
	server.ListenAndServe()
}

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "base", "")
}
