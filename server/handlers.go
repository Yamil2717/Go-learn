package main

import (
	"encoding/json"
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

func router(datafile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list := todo.List{}
		if err := list.Get(datafile); err != nil {
			errorReply(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		switch r.Method {
		case http.MethodGet:
			getAllHandler(w, r, &list)
		case http.MethodPost:
			addHandler(w, r, &list, datafile)
		default:
			errorReply(w, r, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}
}

func getAllHandler(w http.ResponseWriter, r *http.Request, list *todo.List) {
	response := todoResponse{Results: *list}
	jsonReply(w, r, http.StatusOK, &response)
}

func addHandler(w http.ResponseWriter, r *http.Request, list *todo.List, datafile string) {
	type NewTask struct {
		Task string `json:"task"`
	}

	item := NewTask{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&item); err != nil {
		errorReply(w, r, http.StatusBadRequest, err.Error())
		return
	}

	list.Add(item.Task)
	if err := list.Save(datafile); err != nil {
		errorReply(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	textReply(w, r, http.StatusCreated, "")
}
