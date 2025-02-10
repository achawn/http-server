package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		//w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
