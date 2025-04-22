package main

import "net/http"

func main() {
	httpMux := http.NewServeMux()
	server := http.Server{
		Handler: httpMux,
		Addr: ":8080",
	}
	http.ListenAndServe(server.Addr, server.Handler)

}