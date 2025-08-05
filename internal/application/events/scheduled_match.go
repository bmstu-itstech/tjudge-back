package events

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/contest"
)

type ScheduledMatchEventHandler struct {
	contestRepos ports.ContestRepository
	programRepos ports.ProgramRepository
	judge        ports.Judge
}

func (h *ScheduledMatchEventHandler) Handle(ctx context.Context, event contest.ScheduledMatchEvent) error {
	c, err := h.contestRepos.Contest(ctx, event.ContestID)
	if err != nil {
		return err
	}

	p1, found, err := h.programRepos.LastTeamProgram(ctx, event.Team1ID)
	if err != nil {
		return err
	}

	// Skip if not found
	if !found {
		return nil
	}

	p2, found, err := h.programRepos.LastTeamProgram(ctx, event.Team2ID)
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

	return c.FinishMatch(event.MatchID, r1, r2)
}
