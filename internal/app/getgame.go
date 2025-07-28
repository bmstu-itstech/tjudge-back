package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type GetGame struct {
	GameId string
}

type GetGameHandler decorator.QueryHandler[GetGame, Game]

type getGameHandler struct {
	r tjudge.GameRepository
}

func (h getGameHandler) Handle(ctx context.Context, q GetGame) (Game, error) {
	game, err := h.r.Game(ctx, tjudge.GameId(q.GameId))
	if err != nil {
		return Game{}, err
	}
	return gameToDto(game), nil
}

func NewGetGameHandler(
	r tjudge.GameRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetGameHandler {
	return decorator.ApplyQueryDecorators(getGameHandler{r}, l, mc)
}
