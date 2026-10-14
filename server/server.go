package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func newMux(datafile string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/todo", getAllHandler(datafile))
	mux.HandleFunc("/todo/", getAllHandler(datafile))

	return mux
}

func textReply(w http.ResponseWriter, r *http.Request, status int, payload string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(payload))
}

func errorReply(w http.ResponseWriter, r *http.Request, status int, payload string) {
	log.Printf("Error %d: %v", status, payload)
	http.Error(w, payload, status)
}

func jsonReply(w http.ResponseWriter, r *http.Request, status int, payload *todoResponse) {
	body, err := json.Marshal(payload)
	if err != nil {
		errorReply(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}
