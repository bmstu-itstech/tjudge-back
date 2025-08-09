package commands

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/contest"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
	"github.com/google/uuid"
)

type RemoveContestGame struct {
	GameId    string
	ContestId string
}

type RemoveContestGameHandler decorator.CommandHandler[RemoveContestGame]

type removeContestGameHandler struct {
	cr ports.ContestRepository
}

func (h removeContestGameHandler) Handle(ctx context.Context, cmd RemoveContestGame) error {
	contestId, err := uuid.Parse(cmd.ContestId)
	if err != nil {
		return err
	}
	cont, err := h.cr.Contest(ctx, shared.ID(contestId))
	if err != nil {
		return err
	}

	gameId, err := uuid.Parse(cmd.GameId)
	if err != nil {
		return err
	}
	_, present := cont.Games[shared.ID(gameId)]
	if !present {
		return contest.ErrContestNotContainGame
	}
	delete(cont.Games, shared.ID(gameId))
	return h.cr.Upsert(ctx, cont)
}

func NewRemoveContestGameHandler(
	cr ports.ContestRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) RemoveContestGameHandler {
	return decorator.ApplyCommandDecorators(removeContestGameHandler{cr}, l, mc)
}
