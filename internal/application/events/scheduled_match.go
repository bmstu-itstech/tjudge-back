package events

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/contest"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type MatchScheduledEventConsumer decorator.EventConsumer

type matchScheduledEventHandler struct {
	contestRepos ports.ContestRepository
	programRepos ports.ProgramRepository
	judge        ports.Judge
}

func (h matchScheduledEventHandler) Handle(ctx context.Context, event contest.MatchScheduledEvent) error {
	c, err := h.contestRepos.Contest(ctx, event.ContestID)
	if err != nil {
		return err
	}

	p1, found, err := h.programRepos.Active(ctx, event.Team1ID, event.GameID)
	if err != nil {
		return err
	}

	// Skip if not found
	if !found {
		return nil
	}

	p2, found, err := h.programRepos.Active(ctx, event.Team2ID, event.GameID)
	if err != nil {
		return err
	}

	// Skip if not found
	if !found {
		return nil
	}

	r1, r2, err := h.judge.Run(ctx, p1.Path, p2.Path)
	if err != nil {
		return err
	}

	err = c.FinishMatch(event.GameID, event.MatchID, r1, r2)
	if err != nil {
		return err
	}

	return h.contestRepos.Upsert(ctx, c)
}

func NewMatchScheduledConsumer(
	contestRepos ports.ContestRepository,
	programRepos ports.ProgramRepository,
	judge ports.Judge,
	l *slog.Logger,
	mc decorator.MetricsClient,
) decorator.EventConsumer {
	return decorator.ApplyConsumerDecorators(matchScheduledEventHandler{contestRepos, programRepos, judge}, l, mc)
}
