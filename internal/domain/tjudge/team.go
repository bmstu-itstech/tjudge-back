package tjudge

import (
	"context"
	"errors"
	"time"
)

type TeamId uuid

var ErrTeamNotExist = errors.New("team doesn't exist")
var ErrInvalidTeam = errors.New("invalid team passed")

type Team struct {
	Id        TeamId
	Name      string
	CreatedAt time.Time
	Contest   ContestId
	JoinCode  string
}

type TeamRepository interface {
	Team(context.Context, TeamId) (Team, error)
	Add(context.Context, Team) error
}
