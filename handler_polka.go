package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerPolka(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		}
	}

	authHeader := r.Header.Get("Authorization")
	key := strings.Replace(authHeader, "ApiKey ", "", 1)

	if key != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	users, err := cfg.DB.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get users")
		return
	}

	for _, user := range users {
		if user.ID == params.Data.UserId {
			err := cfg.DB.UpdateUserChirpyRed(user.ID)
			if err != nil {
				respondWithError(w, http.StatusNotFound, "Couldn't update user")
				return
			}
			respondWithJSON(w, http.StatusNoContent, nil)
		}
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
