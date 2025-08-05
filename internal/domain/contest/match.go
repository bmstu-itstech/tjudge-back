package contest

import (
	"errors"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type MatchResult struct {
	R1 Result
	R2 Result
}

type Match struct {
	ID          shared.ID
	ContestID   shared.ID
	Team1ID     shared.ID
	Team2ID     shared.ID
	ScheduledAt time.Time
	Result      *MatchResult
	FinishedAt  *time.Time
}

func ScheduleMatch(contestID shared.ID, team1ID shared.ID, team2ID shared.ID) (*Match, ScheduledMatchEvent, error) {
	m := &Match{
		ID:          shared.NewID(),
		Team1ID:     team1ID,
		Team2ID:     team2ID,
		ScheduledAt: time.Now(),
		Result:      nil,
		FinishedAt:  nil,
	}

	e := ScheduledMatchEvent{
		MatchID:   m.ID,
		ContestID: contestID,
		Team1ID:   team1ID,
		Team2ID:   team2ID,
		Timestamp: time.Now(),
	}

	return m, e, nil
}

var ErrMatchAlreadyFinished = errors.New("match is already finished")

func (m *Match) Finish(r1 Result, r2 Result) error {
	if m.FinishedAt != nil {
		return ErrMatchAlreadyFinished
	}
	m.Result = &MatchResult{
		R1: r1,
		R2: r2,
	}
	now := time.Now()
	m.FinishedAt = &now
	return nil
}
