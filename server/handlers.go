package main

import (
	"net/http"

	"github.com/Yamil2717/todo"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorReply(w, r, http.StatusNotFound, "404 page not found")
		return
	}

	textReply(w, r, http.StatusOK, "Hello World")
}

func getAllHandler(datafile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list := todo.List{}
		if err := list.Get(datafile); err != nil {
			errorReply(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		response := todoResponse{Results: list}
		jsonReply(w, r, http.StatusOK, &response)
	}
}
