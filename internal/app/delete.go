package app

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type Delete struct {
	AuthorID int
	Id int64
}

type DeleteHandler struct {
	userRep tjudge.UserRepository
}

func NewDeleteHandler(userRep tjudge.UserRepository) DeleteHandler {
	return DeleteHandler{userRep}
}


func (h *DeleteHandler) Execute(ctx context.Context, cmd Delete) error {
	author, _ := h.userRep.User(ctx, tjudge.UserID(cmd.AuthorID))
	err := author.CanDelete()
	if err != nil {
		return err
	}
	err = h.userRep.Delete(ctx, tjudge.UserID(cmd.Id))
	return err
}
