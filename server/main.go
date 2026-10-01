package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello World"))
}

func main() {
	port := flag.Int("p", 8080, "Puerto donde escucha el servidor")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	log.Fatal(http.ListenAndServe(addr, mux))
}
