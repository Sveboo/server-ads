package users

import (
	"regexp"
)

type User struct {
	ID    int64
	Name  string
	Email string
}

var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

func New(name string, email string) *User {
	if !IsLetter(name) {
		name = "auto-generated-string"
	}
	if !IsLetter(email) {
		email = "auto-generated-email"
	}
	return &User{
		ID:    0,
		Name:  name,
		Email: email,
	}
}
