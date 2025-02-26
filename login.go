package main

import "encoding/json"
import "github.com/google/uuid"
import "internal/auth"
import "net/http"
import "context"
import "time"
import "internal/database"


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

	tk, err := auth.MakeJWT(user.ID, cfg.Secret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "Error making jwt")
		return
	}

	rtk, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "Error creating token")
		return

	}
	_, err = cfg.Db.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token: rtk,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, 500, "Error saving token")
		return
	}

	respondWithJson(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: tk,
		RefreshToken: rtk,
	})





}
