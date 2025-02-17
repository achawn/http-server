package main

import "context"
import "encoding/json"
import "net/http"
import "time"
import "github.com/google/uuid"
import "fmt"

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type userParams struct {
		Email string `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding user")
	}

	user, err := cfg.Db.CreateUser(context.Background(), params.Email)
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
