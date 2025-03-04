package main

import "context"
import "net/http"
import "github.com/google/uuid"
import "encoding/json"
import "internal/auth"

type webhookParams struct {
	Event string `json:"event"`
	Data struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Error getting api key")
		return
	}

	if key != cfg.Polka {
		respondWithError(w, 401, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := webhookParams{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding webhook")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	user, err := cfg.Db.GetUserFromID(context.Background(), params.Data.UserID)
	if user.ID == uuid.Nil {
		respondWithError(w, 404, "User not found")
		return
	}

	_, err = cfg.Db.UpgradeUser(context.Background(), user.ID)
	if err != nil {
		respondWithError(w, 500, "Error upgrading user")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
