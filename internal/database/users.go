package database

import "errors"

func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:       id,
		Email:    email,
		Password: password,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("user not found")
}

func (db *DB) GetUserByRefreshToken(token string) (User, error) {
	dbStr, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStr.Users {
		if user.RefreshToken == token {
			return user, nil
		}
	}

	return User{}, errors.New("user not found")
}

func (db *DB) UpdateUser(id int, email, hashedPassword string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return User{}, errors.New("user not found")
	}

	user.Email = email
	user.Password = hashedPassword
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpdateUserToken(id int, token string, expiresAt int64) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return errors.New("user not found")
	}

	user.RefreshToken = token
	user.ExpiresAt = expiresAt
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RevokeUserToken(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return errors.New("user not found")
	}

	user.RefreshToken = ""
	user.ExpiresAt = 0
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
