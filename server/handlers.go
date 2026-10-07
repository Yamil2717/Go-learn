package main

import "net/http"

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorReply(w, r, http.StatusNotFound, "404 page not found")
		return
	}

	textReply(w, r, http.StatusOK, "Hello World")
}
