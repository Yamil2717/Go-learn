package main

import "net/http"

func newMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	return mux
}
