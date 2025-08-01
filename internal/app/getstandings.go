package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type GetStandings struct {
	GameId string
}

type GetStandingsHandler decorator.QueryHandler[GetStandings, map[string]int]

type getStandingsHandler struct {
	r tjudge.RoundEventStore
	p tjudge.ProgramRepository
}

func (h getStandingsHandler) Handle(ctx context.Context, a GetStandings) (map[string]int, error) {
	standings := make(map[string]int)
	active, err := h.p.Actives(ctx, tjudge.GameId(a.GameId))
	if err != nil {
		return nil, err
	}
	events, err := h.r.Load(ctx, tjudge.GameId(a.GameId))
	if err != nil {
		return nil, err
	}
	for _, v := range events {
		if v.Type != "finished" {
			continue
		}
		stale := false
		for team_id, prog_id := range v.Programs {
			if active[team_id] != prog_id {
				stale = true
				break
			}
		}
		if !stale {
			for team_id, result := range v.Results {
				standings[string(team_id)] += result
			}
		}
	}
	return standings, nil
}

func NewGetStandingsHandler(
	r tjudge.RoundEventStore,
	p tjudge.ProgramRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetStandingsHandler {
	return decorator.ApplyQueryDecorators(getStandingsHandler{r, p}, l, mc)
}
