package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Yamil2717/todo"
)

type testCase struct {
	name            string
	path            string
	expectedCode    int
	expectedContent string
}

func TestComplete(t *testing.T) {
	list := todo.List{}
	list.Add("Primera tarea")
	list.Add("Segunda tarea")
	datafile := filepath.Join(t.TempDir(), "todos.json")
	if err := list.Save(datafile); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(newMux(datafile))
	defer server.Close()

	for _, tc := range []struct {
		path   string
		status int
	}{
		{"/todo/0", http.StatusBadRequest},
		{"/todo/2?complete", http.StatusBadRequest},
		{"/todo/1?complete", http.StatusNoContent},
	} {
		req, err := http.NewRequest(http.MethodPatch, server.URL+tc.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != tc.status {
			t.Errorf("Para %s se esperaba status %d, se obtuvo %d", tc.path, tc.status, res.StatusCode)
		}
	}

	var saved todo.List
	if err := saved.Get(datafile); err != nil {
		t.Fatal(err)
	}
	if len(saved) != 2 || saved[0].Done || !saved[1].Done || saved[1].CompletedAt.IsZero() {
		t.Errorf("Lista guardada inesperada: %+v", saved)
	}
}

func TestDelete(t *testing.T) {
	list := todo.List{}
	list.Add("Primera tarea")
	list.Add("Segunda tarea")
	datafile := filepath.Join(t.TempDir(), "todos.json")
	if err := list.Save(datafile); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(newMux(datafile))
	defer server.Close()

	for _, tc := range []struct {
		path   string
		status int
	}{
		{"/todo/2", http.StatusBadRequest},
		{"/todo/0", http.StatusNoContent},
	} {
		req, err := http.NewRequest(http.MethodDelete, server.URL+tc.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != tc.status {
			t.Errorf("Para %s se esperaba status %d, se obtuvo %d", tc.path, tc.status, res.StatusCode)
		}
	}

	var saved todo.List
	if err := saved.Get(datafile); err != nil {
		t.Fatal(err)
	}
	if len(saved) != 1 || saved[0].Task != "Segunda tarea" {
		t.Errorf("Lista guardada inesperada: %+v", saved)
	}
}

func TestGetOne(t *testing.T) {
	list := todo.List{}
	list.Add("Primera tarea")
	list.Add("Segunda tarea")
	datafile := filepath.Join(t.TempDir(), "todos.json")
	if err := list.Save(datafile); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(newMux(datafile))
	defer server.Close()

	res, err := http.Get(server.URL + "/todo/1")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	var response struct {
		Results      []struct{ Task string }
		TotalResults int `json:"total_results"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK || response.TotalResults != 1 || len(response.Results) != 1 || response.Results[0].Task != "Segunda tarea" {
		t.Errorf("Respuesta inesperada: status %d, resultado %+v", res.StatusCode, response)
	}

	for _, path := range []string{"/todo/-1", "/todo/2", "/todo/abc", "/todo/"} {
		t.Run(path, func(t *testing.T) {
			res, err := http.Get(server.URL + path)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != http.StatusBadRequest {
				t.Errorf("Se esperaba status 400, se obtuvo %d", res.StatusCode)
			}
		})
	}

	res, err = http.Post(server.URL+"/todo/0", "application/json", strings.NewReader(`{"task":"Otra tarea"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Se esperaba status 405, se obtuvo %d", res.StatusCode)
	}
}

func setupAPI(t *testing.T) (url string, cleaner func()) {
	t.Helper()

	server := httptest.NewServer(newMux("datafile.json"))

	url = server.URL
	cleaner = func() {
		server.Close()
	}

	return url, cleaner
}

func TestGet(t *testing.T) {
	url, cleaner := setupAPI(t)
	defer cleaner()

	testCases := []testCase{
		{
			name:            "TestGetRoot",
			path:            "/",
			expectedCode:    http.StatusOK,
			expectedContent: "Hello World",
		},
		{
			name:            "TestGetNotFound",
			path:            "/about",
			expectedCode:    http.StatusNotFound,
			expectedContent: "404 page not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := http.Get(url + tc.path)
			if err != nil {
				t.Fatalf("No se pudo hacer la petición GET a %s: %v", url+tc.path, err)
			}
			defer res.Body.Close()

			if res.StatusCode != tc.expectedCode {
				t.Errorf("Se esperaba el status %s, se obtuvo %s", http.StatusText(tc.expectedCode), http.StatusText(res.StatusCode))
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("No se pudo leer el cuerpo de la respuesta: %v", err)
			}

			switch res.Header.Get("Content-Type") {
			case "text/plain", "text/plain; charset=utf-8":
				if !strings.Contains(string(body), tc.expectedContent) {
					t.Errorf("Se esperaba %q, se obtuvo %q", tc.expectedContent, body)
				}
			default:
				t.Fatalf("Unsupported Content-Type: %q", res.Header.Get("Content-Type"))
			}
		})
	}
}
