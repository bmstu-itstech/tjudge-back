package app

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type Register struct {
	Username	 string
	Fullname	 string
	Password 	 string
}

type RegisterHandler struct {
	userRep tjudge.UserRepository
}

func NewRegisterHandler(userRep tjudge.UserRepository) RegisterHandler {
	return RegisterHandler{userRep}
}

func (h *RegisterHandler) Execute(ctx context.Context, cmd Register) error {
	newuser, err := tjudge.NewUserPrototype(cmd.Username, cmd.Fullname, cmd.Password, false)
	if err == nil {
		_, err = h.userRep.Create(ctx, newuser)
	}
	return err;
}
