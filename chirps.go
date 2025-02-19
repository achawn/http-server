package main

import "encoding/json"
import "github.com/google/uuid"
import "time"
import "net/http"
import "context"
import "internal/database"

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	type chirpParams struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := chirpParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding chirp")
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "{\"error\": \"Chirp is too long\"}")
		return
	}

	cleaned := removeProfanity(params.Body)
	params.Body = cleaned

	// create chirp
	chirp, err := cfg.Db.CreateChirp(context.Background(), database.CreateChirpParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Body: params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, 500, "Error creating chirp")
		return
	}

	respondWithJson(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.Db.GetChirps(context.Background())
	if err != nil {
		respondWithError(w, 500, "Error getting chirps")
		return
	}

	response := make([]Chirp, len(chirps))
	for i, c := range chirps {
		response[i] = Chirp{
			ID: c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body: c.Body,
			UserID: c.UserID,
		}
	}


	respondWithJson(w, http.StatusOK, response)
}
