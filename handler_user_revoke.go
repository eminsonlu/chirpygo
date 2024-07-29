package main

import (
	"github.com/eminsonlu/chirpygo/internal/auth"
	"net/http"
	"strconv"
	"time"
)

func (cfg *apiConfig) handlerUsersRevoke(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
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

	err = cfg.DB.RevokeUserToken(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token")
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
