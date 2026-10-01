package main

import (
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the home page."))
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the about page."))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/about", aboutHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
