package app

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type Login struct {
	Username    string
	Password string
}

type LoginHandler struct {
	userRep tjudge.UserRepository
	generator tjudge.TokenGenerator
}

func NewLoginHandler(userRep tjudge.UserRepository, generator tjudge.TokenGenerator) LoginHandler {
	return LoginHandler{userRep, generator}
}

func (h *LoginHandler) Execute(ctx context.Context, cmd Login) (string, error) {
	user, err := h.userRep.UserByUsername(ctx, cmd.Username)
	if err != nil {
		return "", err
	}
	err = user.PasswordCompare(cmd.Password)
	if err != nil {
		return "", err
	}
	return h.generator.Generate(user.Id())
}
