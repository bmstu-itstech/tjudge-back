package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type CreateGame struct {
	Name        string
	Players     uint
	RulesUrl    string
	AllowedExts []string
}

type CreateGameHandler decorator.CommandHandler[CreateGame]

type createGameHandler struct {
	r tjudge.GameRepository
}

func (h createGameHandler) Handle(ctx context.Context, cmd CreateGame) error {
	Game, err := tjudge.NewGame(
		cmd.Name,
		cmd.Players,
		cmd.RulesUrl,
		cmd.AllowedExts)
	if err != nil {
		return err
	}
	return h.r.Upsert(ctx, Game)
}

func NewCreateGameHandler(
	r tjudge.GameRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) CreateGameHandler {
	return decorator.ApplyCommandDecorators(createGameHandler{r}, l, mc)
}
