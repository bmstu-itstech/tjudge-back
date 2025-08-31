package tjudge

import (
	"bytes"
	"errors"
	"fmt"
)

type UserID int

type User struct {
	UserPrototype
	id UserID
}

func (u *User) Id() UserID {
	return u.id
}

func newUser(id UserID, prototype UserPrototype) (*User, error) {
	if id == 0 {
		return nil, fmt.Errorf("%w: invalid ID", ErrInvalidInput)
	}
	return &User{
		prototype,
		id,
	}, nil
}

func (u *User) SetFullname(name string) error {
	if name == "" {
		return fmt.Errorf("%w: expected not empty name", ErrInvalidInput)
	}
	u.UserPrototype.fullname = name
	return nil
}

func (u *User) CanDelete() error {
	if !u.isAdmin {
		return fmt.Errorf("%w: only admins can delete", ErrNoAccess)
	}
	return nil
}

func RestoreUser(id int64, username string, fullname string, passhash []byte) *User {
	return &User{
		UserPrototype: UserPrototype{username, fullname, passhash, false},
		id:            UserID(id),
	}
}

var ErrPasswordMismatch = errors.New("passwords mismatch")

func (u *User) PasswordCompare(password string) error {
	hash := hashSHA256(password)
	if !bytes.Equal(hash, u.passwordHash) {
		return ErrPasswordMismatch
	}
	return nil
}
