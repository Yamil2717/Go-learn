package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupAPI(t *testing.T) (url string, cleaner func()) {
	t.Helper()

	server := httptest.NewServer(newMux())

	url = server.URL
	cleaner = func() {
		server.Close()
	}

	return url, cleaner
}

func TestGetRoot(t *testing.T) {
	path := "/"
	expectedCode := http.StatusOK
	expectedContent := "Hello World"

	url, cleaner := setupAPI(t)
	defer cleaner()

	res, err := http.Get(url + path)
	if err != nil {
		t.Fatalf("No se pudo hacer la petición GET a %s: %v", url+path, err)
	}
	defer res.Body.Close()

	if res.StatusCode != expectedCode {
		t.Errorf("Se esperaba el status %s, se obtuvo %s", http.StatusText(expectedCode), http.StatusText(res.StatusCode))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("No se pudo leer el cuerpo de la respuesta: %v", err)
	}

	if !strings.Contains(string(body), expectedContent) {
		t.Errorf("Se esperaba %q, se obtuvo %q", expectedContent, body)
	}
}

func TestGetNotFound(t *testing.T) {
	path := "/about"
	expectedCode := http.StatusNotFound
	expectedContent := "404 page not found"

	url, cleaner := setupAPI(t)
	defer cleaner()

	res, err := http.Get(url + path)
	if err != nil {
		t.Fatalf("No se pudo hacer la petición GET a %s: %v", url+path, err)
	}
	defer res.Body.Close()

	if res.StatusCode != expectedCode {
		t.Errorf("Se esperaba el status %s, se obtuvo %s", http.StatusText(expectedCode), http.StatusText(res.StatusCode))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("No se pudo leer el cuerpo de la respuesta: %v", err)
	}

	if !strings.Contains(string(body), expectedContent) {
		t.Errorf("Se esperaba %q, se obtuvo %q", expectedContent, body)
	}
}
