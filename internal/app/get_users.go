package app

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)


type GetUsersHandler struct {
	userRep tjudge.UserRepository
}

func NewGetUsersHandler(userRep tjudge.UserRepository) GetUsersHandler {
	return GetUsersHandler{userRep}
}

func (h *GetUsersHandler) Execute(ctx context.Context) ([]*User, error) {
	users, err := h.userRep.Users(ctx)
	if err != nil {
		return nil, err
	}
	var dusers = []*User{}
	for _, user := range users {
		duser := User{int(user.Id()), user.Username(), user.Fullname()}
		dusers = append(dusers, &duser)
	}
	return dusers, nil
}
