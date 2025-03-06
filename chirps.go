package main

import "encoding/json"
import "github.com/google/uuid"
import "time"
import "net/http"
import "context"
import "internal/database"
import "internal/auth"
import "slices"

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

	tk, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, "Error getting token")
		return
	}

	user_id, err := auth.ValidateJWT(tk, cfg.Secret)
	if err != nil {
		respondWithError(w, 401, "Token not valid")
		return
	}

	type chirpParams struct {
		Body string `json:"body"`
		//UserID uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := chirpParams{}
	err = decoder.Decode(&params)
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
		UserID: user_id,
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

type Sort string
const (
	SortAsc Sort = "asc"
	SortDesc Sort = "desc"
)

func getSort(s string) Sort {
	switch s {
	case string(SortAsc):
		return SortAsc
	case string(SortDesc):
		return SortDesc
	default:
		return SortAsc
	}
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	a_id := r.URL.Query().Get("author_id")
	s := r.URL.Query().Get("sort")
	sort := getSort(s)
	var chirps []database.Chirp
	var err error
	if a_id != "" {
		id, err := uuid.Parse(a_id)
		if err != nil {
			respondWithError(w, 500, "Error parsing id")
			return
		}
		chirps, err = cfg.Db.GetChirpsByUser(context.Background(), id)
	} else {
		chirps, err = cfg.Db.GetChirps(context.Background())
	}

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

	if sort == "desc" {
		slices.Reverse(response)
	}

	respondWithJson(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, 500, "Error parsing id")
		return
	}

	chirp, err := cfg.Db.GetChirp(context.Background(), id)

	if chirp.ID == uuid.Nil {
		respondWithError(w, 404, "Chirp not found")
		return
	}

	if err != nil {
		respondWithError(w, 500, "Error getting chirp")
		return
	}

	respondWithJson(w, http.StatusOK, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})

}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	tk, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Error getting token")
		return
	}

	user_id, err := auth.ValidateJWT(tk, cfg.Secret)
	if err != nil {
		respondWithError(w, 401, "Token not valid")
		return
	}

	idStr := r.PathValue("chirpID")
	id, err := uuid.Parse(idStr)

	chirp, err := cfg.Db.GetChirp(context.Background(), id)

	if chirp.ID == uuid.Nil {
		respondWithError(w, 404, "Chirp not found")
		return
	}

	if err != nil {
		respondWithError(w, 500, "Error getting chirp")
		return
	}

	if user_id != chirp.UserID {
		respondWithError(w, 403, "Forbidden")
		return
	}

	err = cfg.Db.DeleteChirp(context.Background(), id)
	if err != nil {
		respondWithError(w, 500, "Error deleting chirp")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
