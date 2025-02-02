package main

import "net/http"
//import "fmt"


func main() {
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.Handle("/logo.png", http.FileServer(http.Dir("./assets")))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	})
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
