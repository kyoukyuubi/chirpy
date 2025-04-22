package main

import (
	"log"
	"net/http"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {
	const filepathRoot = "."
	const port = "8080"


	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	server := http.Server{
		Handler: mux,
		Addr: ":" + port,
	}
	mux.HandleFunc("/healthz", healthzHandler)
	
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}