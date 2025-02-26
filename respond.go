package main

import "net/http"
import "encoding/json"

func respondWithError(w http.ResponseWriter, code int, msg string) {
	http.Error(w, msg, code)
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonData, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, "{\"error\": \"Error marshalling json\"}")
		return
	}
	w.Write([]byte(jsonData))
}
