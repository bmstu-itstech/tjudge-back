package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type UpsertGame struct {
	GameId      string
	Name        string
	Players     uint
	RulesUrl    string
	AllowedExts []string
}

type UpsertGameHandler decorator.CommandHandler[UpsertGame]

type upsertGameHandler struct {
	r tjudge.GameRepository
}

func (h upsertGameHandler) Handle(ctx context.Context, cmd UpsertGame) error {
	Game, err := tjudge.ParseGame(
		tjudge.GameId(cmd.GameId),
		cmd.Name,
		cmd.Players,
		cmd.RulesUrl,
		cmd.AllowedExts)
	if err != nil {
		return err
	}
	return h.r.Upsert(ctx, Game)
}

func NewUpsertGameHandler(
	r tjudge.GameRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) UpsertGameHandler {
	return decorator.ApplyCommandDecorators(upsertGameHandler{r}, l, mc)
}
