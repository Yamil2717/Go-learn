package main

import (
	"log"
	"net/http"
)

type homeHandler struct {
}

func (hh homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello About"))
}

func helpHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Help"))
}

func main() {
	mux := http.NewServeMux()

	home := homeHandler{}
	mux.Handle("/", home)

	about := http.HandlerFunc(aboutHandler)
	mux.Handle("/about", about)

	mux.HandleFunc("/help", helpHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
