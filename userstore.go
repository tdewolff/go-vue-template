package main

import (
	"database/sql"
	"log"

	"github.com/tdewolff/auth"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db}
}

// Login logs in every user, creating a new account if it is a new user
func (s *UserStore) Get(emailAddress string) (int64, bool) {
	// Login or register
	var userID int64
	if err := s.db.QueryRow(`SELECT id FROM users WHERE email=?`, emailAddress).Scan(&userID); err != nil && err != sql.ErrNoRows {
		log.Println("userstore validation failed:", err)
		return 0, false
	} else if err == sql.ErrNoRows {
		return 0, false
	}
	return userID, true
}

func (s *UserStore) Set(user *auth.User) (int64, bool) {
	res, err := s.db.Exec(`INSERT INTO users (email) VALUES (?)`, user.Email)
	if err != nil {
		log.Println("userstore registration failed:", err)
		return 0, false
	}

	userID, err := res.LastInsertId()
	if err != nil {
		log.Println("userstore registration failed:", err)
		return 0, false
	}
	return userID, true
}

func (s *UserStore) SetToken(userID int64, provider string, token string) error {
	if _, err := s.db.Exec(`INSERT INTO social_tokens (user_id, provider, token) VALUES (?, ?, ?)`, userID, provider, token); err != nil {
		return err
	}
	return nil
}

func (s *UserStore) GetTokens(userID int64) (map[string]string, error) {
	rows, err := s.db.Query(`SELECT provider, token FROM social_tokens WHERE user_id=?`, userID)
	if err != nil {
		return nil, err
	}

	tokens := map[string]string{}
	for rows.Next() {
		var provider string
		var token string
		if err := rows.Scan(&provider, &token); err != nil {
			return nil, err
		}
		tokens[provider] = token
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}
