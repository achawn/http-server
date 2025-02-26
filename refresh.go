package main

import "net/http"
import "internal/auth"
import "time"

type tokenParams struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, "Token missing or broken")
		return
	}

	user, err := cfg.Db.GetUserFromToken(r.Context(), tokenStr)
	if err != nil {
		respondWithError(w, 401, "User not found from token")
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.Secret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, 401, "Couldn't make token")
		return
	}

	respondWithJson(w, http.StatusOK, tokenParams{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, "Token missing or broken")
		return
	}

	_, err = cfg.Db.RevokeToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 400, "Couldn't revoke token")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
