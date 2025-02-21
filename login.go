package main

import "encoding/json"
import "github.com/google/uuid"
import "internal/auth"
import "net/http"
import "context"

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding user")
		return
	}

	user, err := cfg.Db.GetUser(context.Background(), params.Email)
	if user.ID == uuid.Nil {
		respondWithError(w, 404, "User not found")
		return
	}

	if err != nil {
		respondWithError(w, 500, "Error getting user")
		return
	}

	ok := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if ok != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	respondWithJson(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})





}
