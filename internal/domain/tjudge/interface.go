package tjudge

import (
	"context"
	"errors"
)

var ErrUserAlreadyExists = errors.New("user already exist")
var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, proto *UserPrototype) (*User, error)
	User(ctx context.Context, id UserID) (*User, error)
	UserByUsername(ctx context.Context, username string) (*User, error)
	Users(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id UserID) error
}

type Token = string

type TokenGenerator interface {
	Generate(user_id UserID) (Token, error)
}

var ErrTeamAlreadyExist = errors.New("command already exist")
var ErrTeamNotFound = errors.New("command not found")
var ErrTeamFull = errors.New("command full")
var ErrUserAlreadyInTeam = errors.New("user already in team")

type TeamRepository interface {
	Save(ctx context.Context, team *Team) (*Team, error)
	Team(ctx context.Context, code string) (*Team, error)
	Teams(ctx context.Context) ([]*Team, error)
	TeamsByContest(ctx context.Context, contest ContestID) ([]*Team, error)
	Update(ctx context.Context, team *Team) error
	Delete(ctx context.Context, code string) error
}
