package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Yamil2717/todo"
)

type testCase struct {
	name            string
	path            string
	expectedCode    int
	expectedContent string
	expItems        int
}

func TestAdd(t *testing.T) {
	url := setupAPI(t)
	taskName := "Nueva tarea"

	if !t.Run("Add", func(t *testing.T) {
		var body bytes.Buffer
		item := struct {
			Task string `json:"task"`
		}{Task: taskName}
		if err := json.NewEncoder(&body).Encode(item); err != nil {
			t.Fatal(err)
		}

		res, err := http.Post(url+"/todo", "application/json", &body)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("Se esperaba status 201, se obtuvo %d", res.StatusCode)
		}
	}) {
		return
	}

	t.Run("CheckAdd", func(t *testing.T) {
		res, err := http.Get(url + "/todo/3")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("Se esperaba status 200, se obtuvo %d", res.StatusCode)
		}

		var response struct {
			Results todo.List `json:"results"`
		}
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}
		if len(response.Results) != 1 || response.Results[0].Task != taskName {
			t.Errorf("Se esperaba la tarea %q, se obtuvo %+v", taskName, response.Results)
		}
	})
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

func setupAPI(t *testing.T) string {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), "todos-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString("[]"); err != nil {
		file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(newMux(file.Name()))
	cleaner := func() {
		server.Close()
		if err := os.Remove(file.Name()); err != nil {
			t.Errorf("No se pudo eliminar el archivo temporal: %v", err)
		}
	}
	t.Cleanup(cleaner)

	for i := 0; i < 3; i++ {
		var body bytes.Buffer
		item := struct {
			Task string `json:"task"`
		}{Task: fmt.Sprintf("Task %d", i)}
		if err := json.NewEncoder(&body).Encode(item); err != nil {
			t.Fatal(err)
		}
		res, err := http.Post(server.URL+"/todo", "application/json", &body)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("Se esperaba status 201, se obtuvo %d", res.StatusCode)
		}
	}

	return server.URL
}

func TestGet(t *testing.T) {
	url := setupAPI(t)

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
		{
			name:            "GetAll",
			path:            "/todo",
			expectedCode:    http.StatusOK,
			expectedContent: "Task 0",
			expItems:        3,
		},
		{
			name:            "GetOne",
			path:            "/todo/1",
			expectedCode:    http.StatusOK,
			expectedContent: "Task 1",
			expItems:        1,
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

			switch res.Header.Get("Content-Type") {
			case "text/plain", "text/plain; charset=utf-8":
				body, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("No se pudo leer el cuerpo de la respuesta: %v", err)
				}
				if !strings.Contains(string(body), tc.expectedContent) {
					t.Errorf("Se esperaba %q, se obtuvo %q", tc.expectedContent, body)
				}
			case "application/json":
				var response struct {
					Results      todo.List `json:"results"`
					Date         time.Time `json:"date"`
					TotalResults int       `json:"total_results"`
				}
				if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
					t.Fatalf("No se pudo decodificar la respuesta: %v", err)
				}
				if response.TotalResults != tc.expItems || len(response.Results) != tc.expItems {
					t.Errorf("Se esperaban %d tareas, se obtuvo total %d y %d resultados", tc.expItems, response.TotalResults, len(response.Results))
				}
				if len(response.Results) == 0 || response.Results[0].Task != tc.expectedContent {
					t.Errorf("Se esperaba la primera tarea %q, se obtuvo %+v", tc.expectedContent, response.Results)
				}
			default:
				t.Fatalf("Unsupported Content-Type: %q", res.Header.Get("Content-Type"))
			}
		})
	}
}
