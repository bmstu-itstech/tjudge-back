package tjudge

import (
	"context"
	"errors"
	"time"

	"github.com/bmstu-itstech/tjudge-back/pkg/uuid"
)

type TeamId shortUuid

var ErrTeamNotExist = errors.New("team doesn't exist")
var ErrInvalidTeam = errors.New("invalid team passed")

type Team struct {
	Id        TeamId
	Name      string
	Contest   ContestId
	CreatedAt time.Time
	JoinCode  string
}

type TeamRepository interface {
	Team(context.Context, TeamId) (Team, error)
	Upsert(context.Context, Team) error
	ByContest(context.Context, ContestId) ([]Team, error)
	ByJoinCode(context.Context, string) (Team, error)
}

func ParseTeam(id TeamId, name string, contest ContestId, created time.Time, code string) (Team, error) {
	if id == "" || name == "" || contest == "" || code == "" || created.IsZero() {
		return Team{}, ErrInvalidTeam
	}
	return Team{
		id,
		name,
		contest,
		created,
		code,
	}, nil
}

func MustParseTeam(id TeamId, name string, contest ContestId, created time.Time, code string) Team {
	t, err := ParseTeam(id, name, contest, created, code)
	if err != nil {
		panic(err)
	}
	return t
}

func NewTeam(name string, contest ContestId) (Team, error) {
	id := uuid.GenerateShort()
	created := time.Now()
	code := uuid.GenerateShort() // we don't just use the id because... reasons
	return ParseTeam(TeamId(id), name, contest, created, code)
}

func MustNewTeam(name string, contest ContestId) Team {
	t, err := NewTeam(name, contest)
	if err != nil {
		panic(err)
	}
	return t
}
