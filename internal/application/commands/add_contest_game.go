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

type AddContestGame struct {
	GameId    string
	ContestId string
}

type AddContestGameHandler decorator.CommandHandler[AddContestGame]

type addContestGameHandler struct {
	cr ports.ContestRepository
	gr ports.GameRepository
}

func (h addContestGameHandler) Handle(ctx context.Context, cmd AddContestGame) error {
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
	_, already := cont.Games[shared.ID(gameId)]
	if already {
		return contest.ErrContestAlreadyContainGame
	}
	game, err := h.gr.Game(ctx, shared.ID(gameId))
	if err != nil {
		return err
	}
	cont.Games[shared.ID(gameId)] = game
	return h.cr.Upsert(ctx, cont)
}

func NewAddContestGameHandler(
	cr ports.ContestRepository,
	gr ports.GameRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) AddContestGameHandler {
	return decorator.ApplyCommandDecorators(addContestGameHandler{cr, gr}, l, mc)
}
