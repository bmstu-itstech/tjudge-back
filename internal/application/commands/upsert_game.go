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

type UpsertGame struct {
	Id       string
	Name     string
	RulesUrl string
}

type UpsertGameHandler decorator.CommandHandler[UpsertGame]

type upsertGameHandler struct {
	repo ports.GameRepository
}

func (h upsertGameHandler) Handle(ctx context.Context, cmd UpsertGame) error {
	id, err := uuid.Parse(cmd.Id)
	if err != nil {
		return err
	}
	old_game, err := h.repo.Game(ctx, shared.ID(id))
	var matches map[shared.ID]*contest.Match
	if err != nil && err != contest.ErrGameNotExist {
		return err
	} else if err == contest.ErrGameNotExist {
		matches = make(map[shared.ID]*contest.Match)
	} else {
		matches = old_game.Matches
	}

	game, err := contest.ParseGame(
		shared.ID(id),
		cmd.Name,
		cmd.RulesUrl,
		matches,
	)
	if err != nil {
		return err
	}
	return h.repo.Upsert(ctx, &game)
}

func NewUpsertGameHandler(
	repo ports.GameRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) UpsertGameHandler {
	return decorator.ApplyCommandDecorators(upsertGameHandler{repo}, l, mc)
}
