package queries

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type AllGames struct{}

type AllGamesHandler decorator.QueryHandler[AllGames, []Game]

type allGamesHandler struct {
	repo ports.GameRepository
}

func (h allGamesHandler) Handle(ctx context.Context, a AllGames) ([]Game, error) {
	games, err := h.repo.All(ctx)
	if err != nil {
		return nil, err
	}
	return batchGamesToDto(games), nil
}

func NewAllGamesHandler(
	r ports.GameRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) AllGamesHandler {
	return decorator.ApplyQueryDecorators(allGamesHandler{r}, l, mc)
}
