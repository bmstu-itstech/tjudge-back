package contest

import (
	"errors"
	"fmt"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

var ErrInvalidContest = errors.New("invalid contest passed")
var ErrContestNotExist = errors.New("contest doesn't exist")

var ErrContestAlreadyContainGame = errors.New("game is already part of contest")
var ErrContestNotContainGame = errors.New("game isn't part of contest")
var ErrContestNotContainTeam = errors.New("team isn't part of contest")
var ErrContestNotContainMatch = errors.New("match isn't part of contest")

type Contest struct {
	Id     shared.ID
	Name   string
	Starts time.Time
	Ends   time.Time
	Games  map[shared.ID]*Game
	Teams  map[shared.ID]*Team
}

func (c *Contest) ScheduleMatchesFor(teamId shared.ID, gameId shared.ID) ([]MatchScheduledEvent, error) {
	game, ok := c.Games[gameId]
	if !ok {
		return nil, fmt.Errorf("%w (id: %s)", ErrContestNotContainGame, gameId)
	}
	_, ok = c.Teams[gameId] // just to check that the team exists?
	if !ok {
		return nil, fmt.Errorf("%w (id: %s)", ErrContestNotContainTeam, teamId)
	}

	events := make([]MatchScheduledEvent, 0)
	for id, match := range game.Matches {
		if match.Team1ID != teamId && match.Team2ID != teamId {
			continue
		}
		var adversaryId shared.ID
		if match.Team1ID == teamId {
			adversaryId = match.Team1ID
		} else {
			adversaryId = match.Team2ID
		}
		m, ev, err := ScheduleMatch(game.Id, teamId, adversaryId)
		if err != nil {
			return nil, err
		}
		events = append(events, ev)
		delete(game.Matches, id)
		game.Matches[m.ID] = m
	}
	return events, nil
}

func (c *Contest) ScheduleAllMatchesFor(gameId shared.ID) ([]MatchScheduledEvent, error) {
	game, ok := c.Games[gameId]
	if !ok {
		return nil, fmt.Errorf("%w (id: %s)", ErrContestNotContainGame, gameId)
	}
	game.Matches = make(map[shared.ID]*Match)
	events := make([]MatchScheduledEvent, 0)
	for t1 := range c.Teams {
		for t2 := range c.Teams {
			if t1 != t2 {
				m, ev, err := ScheduleMatch(game.Id, t1, t2)
				if err != nil {
					return nil, err
				}
				events = append(events, ev)
				game.Matches[m.ID] = m
			}
		}
	}
	return events, nil
}

func (c *Contest) ScheduleAllMatches() ([]MatchScheduledEvent, error) {
	events := make([]MatchScheduledEvent, 0)
	for _, game := range c.Games {
		new_events, err := c.ScheduleAllMatchesFor(game.Id)
		if err != nil {
			return nil, err
		}
		events = append(events, new_events...)
	}
	return events, nil
}

func (c *Contest) FinishMatch(gameId shared.ID, matchId shared.ID, r1 Result, r2 Result) error {
	game, ok := c.Games[gameId]
	if !ok {
		return fmt.Errorf("%w (id: %s)", ErrContestNotContainGame, gameId)
	}

	match, ok := game.Matches[matchId]
	if !ok {
		return fmt.Errorf("%w (id: %s)", ErrContestNotContainMatch, matchId)
	}
	return match.Finish(r1, r2)
}

func (c *Contest) Standings(gameId shared.ID) (map[shared.ID]Score, error) {
	game, ok := c.Games[gameId]
	if !ok {
		return nil, fmt.Errorf("%w (id: %s)", ErrContestNotContainGame, gameId)
	}

	standings := make(map[shared.ID]Score)
	for _, match := range game.Matches {
		if match.Result.R1.ErrMsg == nil {
			standings[match.Team1ID] += match.Result.R1.Score
		}
		if match.Result.R2.ErrMsg == nil {
			standings[match.Team2ID] += match.Result.R2.Score
		}
	}
	return standings, nil
}

func ParseContest(id shared.ID, name string, starts time.Time, ends time.Time,
	games map[shared.ID]*Game, teams map[shared.ID]*Team,
) (Contest, error) {
	if len(id) == 0 || name == "" || starts.After(ends) || starts.IsZero() ||
		ends.IsZero() || games == nil || teams == nil {
		return Contest{}, ErrInvalidContest
	}
	return Contest{id, name, starts, ends, games, teams}, nil
}

func MustParseContest(id shared.ID, name string, starts time.Time, ends time.Time,
	games map[shared.ID]*Game, teams map[shared.ID]*Team,
) Contest {
	c, err := ParseContest(id, name, starts, ends, games, teams)
	if err != nil {
		panic(err)
	}
	return c
}

func NewContest(name string, starts time.Time, ends time.Time, games map[shared.ID]*Game) (Contest, error) {
	return ParseContest(shared.NewID(), name, starts, ends, games, make(map[shared.ID]*Team))
}

func MustNewContest(name string, starts time.Time, ends time.Time, games map[shared.ID]*Game) Contest {
	c, err := NewContest(name, starts, ends, games)
	if err != nil {
		panic(err)
	}
	return c
}
