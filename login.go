package main

import "encoding/json"
import "github.com/google/uuid"
import "internal/auth"
import "net/http"
import "context"
import "time"


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

	expiration := 3600
	if params.Expires != nil {
		if *params.Expires > 3600 {
			expiration = 3600
		} else if *params.Expires > 0 {
			expiration = *params.Expires
		}
	}
	tk, err := auth.MakeJWT(user.ID, cfg.Secret, time.Duration(expiration) * time.Second)
	if err != nil {
		respondWithError(w, 500, "Error making jwt")
		return
	}

	respondWithJson(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: tk,
	})





}
