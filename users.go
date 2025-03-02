package main

import "context"
import "encoding/json"
import "net/http"
import "time"
import "github.com/google/uuid"
import "fmt"
import "internal/auth"
import "internal/database"

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type userParams struct {
	Email string `json:"email"`
	Password string `json:"password"`
	//Expires *int `json:"expires_in_seconds,omitempty"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding user")
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, "Error hashing password")
		return
	}

	user, err := cfg.Db.CreateUser(context.Background(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashed,
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, 500, "Error creating user")
		return
	}
	respondWithJson(w, http.StatusCreated, response{
		User: User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
	})
}

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Token missing or broken")
		return
	}

	user_id, err := auth.ValidateJWT(tokenStr, cfg.Secret)
	if err != nil {
		respondWithError(w, 401, "Token not valid")
		return
	}

	// user, err := cfg.Db.GetUserFromToken(r.Context(), tokenStr)
	// if err != nil {
	// 	respondWithError(w, 401, "User not found from token")
	// 	return
	// }


	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding user")
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, "Error hashing password")
		return
	}

	updated, err := cfg.Db.UpdateUser(context.Background(), database.UpdateUserParams{
		ID: user_id,
		Email: params.Email,
		HashedPassword: hashed,
	})
	if err != nil {
		respondWithError(w, 500, "Unable to update user")
		return
	}

	respondWithJson(w, http.StatusOK, User{
		ID: updated.ID,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
		Email: updated.Email,
	})

}
