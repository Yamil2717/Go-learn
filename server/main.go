package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	host := flag.String("h", "localhost", "Host donde escucha el servidor")
	port := flag.Int("p", 8080, "Puerto donde escucha el servidor")
	datafile := flag.String("f", "datafile.json", "Archivo con la lista de todos")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)

	server := &http.Server{
		Addr:    addr,
		Handler: newMux(*datafile),
	}

	log.Fatal(server.ListenAndServe())
}
