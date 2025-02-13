package main

import "net/http"
import "sync/atomic"
import "fmt"
import "encoding/json"

type params struct {
	Body string `json:"body"`
	CleanedBody string `json:"cleaned_body"`
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMertrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
	<html>
		<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
		</body>
	</html>
	`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handlerValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}


	decoder := json.NewDecoder(r.Body)
	chirp := params{}
	err := decoder.Decode(&chirp)
	if err != nil {
		http.Error(w, "Error decoding json", http.StatusInternalServerError)
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, 400, "{\"error\": \"Chirp is too long\"}")
		return
	}

	f := removeProfanity(chirp.Body)
	//fmt.Println(f)
	chirp.CleanedBody = f

	jsonData, err := json.Marshal(chirp)
	if err != nil {
		respondWithError(w, 500, "{\"error\": \"Error marshalling json\"}")
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(jsonData))
}


func main() {
	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/api/healthz", handlerReadiness)
	mux.HandleFunc("/admin/metrics", apiCfg.handlerMertrics)
	mux.HandleFunc("/api/validate_chirp", apiCfg.handlerValidate)
	mux.HandleFunc("/admin/reset", apiCfg.handlerReset)
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
