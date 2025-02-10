package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		//w.WriteHeader(http.StatusMethodNotAllowed)
		//w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		http.Error(w, "Method Not ALlowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
