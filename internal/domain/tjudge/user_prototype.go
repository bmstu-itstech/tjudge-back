package tjudge

import (
	"crypto/sha256"
	"fmt"
)

type UserPrototype struct {
	username		string
	fullname		string	
	passwordHash	[]byte
	isAdmin			bool
}

func (u *UserPrototype) Username() string {
	return u.username
}

func (u *UserPrototype) Fullname() string {
	return u.fullname
}

func (u *UserPrototype) Passhash() []byte {
	return u.passwordHash
}

func (u *UserPrototype) IsAdmin() bool {
	return u.isAdmin
}

func NewUserPrototype(username string, fullname	string, password string, isAdmin bool) (*UserPrototype, error) {
	if len(password) < 8 {
		return nil, fmt.Errorf("%w: the password contains less than 8 characters", ErrInvalidInput)
	}
	if username == "" {
		return nil, fmt.Errorf("%w: expected not empty name", ErrInvalidInput)
	}
	passwordHash := hashSHA256(password)
	return &UserPrototype {username, fullname, passwordHash, isAdmin}, nil
}

func MustNewUserPrototype(username string, fullname	string, password string, isAdmin bool) *UserPrototype {
	proto, err := NewUserPrototype(username, fullname, password, isAdmin)
	if err != nil {
		panic(err)
	}
	return proto
}

func (p *UserPrototype) Build(id UserID) (*User, error) {
	return newUser(id, *p)
}

func hashSHA256(input string) []byte {
	hash := sha256.Sum256([]byte(input))
	return hash[:]
}
