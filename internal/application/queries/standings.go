package queries

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
	"github.com/google/uuid"
)

type GetStandings struct {
	ContestId string
	GameId    string
}

type GetStandingsHandler decorator.QueryHandler[GetStandings, map[string]int]

type getStandingsHandler struct {
	r ports.ContestRepository
}

func (h getStandingsHandler) Handle(ctx context.Context, q GetStandings) (map[string]int, error) {
	id, err := uuid.Parse(q.ContestId)
	if err != nil {
		return nil, err
	}
	gameId, err := uuid.Parse(q.GameId)
	if err != nil {
		return nil, err
	}
	contest, err := h.r.Contest(ctx, shared.ID(id))
	if err != nil {
		return nil, err
	}
	standings, err := contest.Standings(shared.ID(gameId))
	if err != nil {
		return nil, err
	}
	return standingsToDto(standings), nil
}

func NewGetStandingsHandler(
	r ports.ContestRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetStandingsHandler {
	return decorator.ApplyQueryDecorators(getStandingsHandler{r}, l, mc)
}
