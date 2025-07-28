package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type GetPlayer struct {
	PlayerId string
}

type GetPlayerHandler decorator.QueryHandler[GetPlayer, Player]

type getPlayerHandler struct {
	r tjudge.PlayerRepository
}

func (h getPlayerHandler) Handle(ctx context.Context, q GetPlayer) (Player, error) {
	player, err := h.r.Player(ctx, tjudge.PlayerId(q.PlayerId))
	if err != nil {
		return Player{}, err
	}
	return playerToDto(player), nil
}

func NewGetPlayerHandler(
	r tjudge.PlayerRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetPlayerHandler {
	return decorator.ApplyQueryDecorators(getPlayerHandler{r}, l, mc)
}
