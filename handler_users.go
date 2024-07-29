package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	users, err := cfg.DB.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve users")
		return
	}

	for _, user := range users {
		if user.Email == params.Email {
			respondWithError(w, http.StatusBadRequest, "Email already exists")
			return
		}
	}

	hashedPass, err := cfg.DB.BcyrptPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, hashedPass)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (cfg *apiConfig) handlerUsersRetrieve(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := cfg.DB.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve users")
		return
	}

	users := []User{}
	for _, dbUser := range dbUsers {
		users = append(users, User{
			ID:    dbUser.ID,
			Email: dbUser.Email,
		})
	}

	respondWithJSON(w, http.StatusOK, users)
}
