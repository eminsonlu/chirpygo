package main

import (
	"github.com/eminsonlu/chirpygo/internal/auth"
	"net/http"
	"strconv"
	"time"
)

func (cfg *apiConfig) handlerUsersRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token")
		return
	}

	user, err := cfg.DB.GetUserByRefreshToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user")
		return
	}

	i, err := strconv.ParseInt(strconv.FormatInt(user.ExpiresAt, 10), 10, 64)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse expires_at")
		return
	}
	tm := time.Unix(i, 0)

	if tm.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired")
		return
	}

	token, err = auth.MakeJWT(user.ID, cfg.JWTSecret, time.Duration(1)*time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}
