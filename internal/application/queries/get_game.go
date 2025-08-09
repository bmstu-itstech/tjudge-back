package queries

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
	"github.com/google/uuid"
)

type GetGame struct {
	GameId string
}

type GetGameHandler decorator.QueryHandler[GetGame, Game]

type getGameHandler struct {
	r ports.GameRepository
}

func (h getGameHandler) Handle(ctx context.Context, q GetGame) (Game, error) {
	id, err := uuid.Parse(q.GameId)
	if err != nil {
		return Game{}, err
	}
	game, err := h.r.Game(ctx, shared.ID(id))
	if err != nil {
		return Game{}, err
	}
	return gameToDto(game), nil
}

func NewGetGameHandler(
	r ports.GameRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetGameHandler {
	return decorator.ApplyQueryDecorators(getGameHandler{r}, l, mc)
}
