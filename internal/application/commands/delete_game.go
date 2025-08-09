package commands

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
	"github.com/google/uuid"
)

type DeleteGame struct {
	Id string
}

type DeleteGameHandler decorator.CommandHandler[DeleteGame]

type deleteGameHandler struct {
	gr ports.GameRepository
	cr ports.ContestRepository
}

func (h deleteGameHandler) Handle(ctx context.Context, cmd DeleteGame) error {
	uid, err := uuid.Parse(cmd.Id)
	if err != nil {
		return err
	}
	id := shared.ID(uid)
	// assuming we want to remove all mentions of the game
	contests, err := h.cr.With(ctx, id)
	if err != nil {
		return err
	}
	for _, c := range contests {
		delete(c.Games, id)
		if err := h.cr.Upsert(ctx, c); err != nil {
			return err
		}
	}
	return h.gr.Delete(ctx, id)
}

func NewDeleteGameHandler(
	gr ports.GameRepository,
	cr ports.ContestRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) DeleteGameHandler {
	return decorator.ApplyCommandDecorators(deleteGameHandler{gr, cr}, l, mc)
}
