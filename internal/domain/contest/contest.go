package contest

import (
	"errors"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type Contest struct {
	ID      shared.ID
	Teams   map[shared.ID]*Team
	Matches map[shared.ID]*Match
}

func (c *Contest) ScheduleAllMatches() ([]ScheduledMatchEvent, error) {
	c.Matches = make(map[shared.ID]*Match)

	for _, team := range c.Teams {
		team.ResetScore()
	}

	events := make([]ScheduledMatchEvent, 0)
	for t1 := range c.Teams {
		for t2 := range c.Teams {
			if t1 != t2 {
				m, ev, err := ScheduleMatch(c.ID, t1, t2)
				if err != nil {
					return nil, err
				}
				events = append(events, ev)
				c.Matches[m.ID] = m
			}
		}
	}

	return events, nil
}

var ErrMatchIsNotInContest = errors.New("match is not in contest")

func (c *Contest) FinishMatch(matchID shared.ID, r1 Result, r2 Result) error {
	match, ok := c.Matches[matchID]
	if !ok {
		return ErrMatchIsNotInContest
	}
	return match.Finish(r1, r2)
}
