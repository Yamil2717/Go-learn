package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

		if r.URL.Path != "/todo" {
			id, err := validateID(strings.TrimPrefix(r.URL.Path, "/todo/"), &list)
			if err != nil {
				errorReply(w, r, http.StatusBadRequest, err.Error())
				return
			}

			switch r.Method {
			case http.MethodGet:
				getOneHandler(w, r, &list, id)
			case http.MethodDelete:
				deleteHandler(w, r, &list, id, datafile)
			default:
				errorReply(w, r, http.StatusMethodNotAllowed, "Method not allowed")
			}
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

func deleteHandler(w http.ResponseWriter, r *http.Request, list *todo.List, id int, datafile string) {
	if err := list.Delete(id); err != nil {
		errorReply(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := list.Save(datafile); err != nil {
		errorReply(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	textReply(w, r, http.StatusNoContent, "")
}

func validateID(value string, list *todo.List) (int, error) {
	id, err := strconv.Atoi(value)
	if err != nil || id < 0 || id >= len(*list) {
		return 0, fmt.Errorf("ID invalido: %s", value)
	}
	return id, nil
}

func getOneHandler(w http.ResponseWriter, r *http.Request, list *todo.List, id int) {
	response := todoResponse{Results: (*list)[id : id+1]}
	jsonReply(w, r, http.StatusOK, &response)
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
