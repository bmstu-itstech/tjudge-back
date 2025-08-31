package app

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type GetUser struct {
	Id int64
}

type GetUserHandler struct {
	userRep tjudge.UserRepository
}

func NewGetUserHandler(userRep tjudge.UserRepository) GetUserHandler {
	return GetUserHandler{userRep}
}

func (h *GetUserHandler) Execute(ctx context.Context, cmd GetUser) (User, error) {
	user, err := h.userRep.User(ctx, tjudge.UserID(cmd.Id))
	if err != nil {
		return User{}, err
	}
	duser := User{int(user.Id()), user.Username(), user.Fullname()}
	return duser, nil
}
