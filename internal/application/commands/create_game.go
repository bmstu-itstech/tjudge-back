package commands

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/contest"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type CreateGame struct {
	Name     string
	RulesUrl string
}

type CreateGameHandler decorator.CommandHandler[CreateGame]

type createGameHandler struct {
	repo ports.GameRepository
}

func (h createGameHandler) Handle(ctx context.Context, cmd CreateGame) error {
	game, err := contest.NewGame(
		cmd.Name,
		cmd.RulesUrl,
	)
	if err != nil {
		return err
	}
	return h.repo.Upsert(ctx, &game)
}

func NewCreateGameHandler(
	repo ports.GameRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) CreateGameHandler {
	return decorator.ApplyCommandDecorators(createGameHandler{repo}, l, mc)
}
