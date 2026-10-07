package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testCase struct {
	name            string
	path            string
	expectedCode    int
	expectedContent string
}

func setupAPI(t *testing.T) (url string, cleaner func()) {
	t.Helper()

	server := httptest.NewServer(newMux())

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
