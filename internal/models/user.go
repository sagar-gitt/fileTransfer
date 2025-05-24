package models

import (
	"encoding/json"
	"io"
)

type GoogleUser struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"verified_email"`
	Name            string `json:"name"`
	Avatar          string `json:"picture"`
}

func ParseGoogleUser(body io.Reader) (*GoogleUser, error) {
	var u GoogleUser
	err := json.NewDecoder(body).Decode(&u)
	return &u, err
}
