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
	id         TeamId
	name       string
	contest_id ContestId
	created_at time.Time
	join_code  string
}

func (t Team) Id() TeamId {
	return t.id
}

func (t Team) Name() string {
	return t.name
}

func (t Team) ContestId() ContestId {
	return t.contest_id
}

func (t Team) CreatedAt() time.Time {
	return t.created_at
}

func (t Team) JoinCode() string {
	return t.join_code
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
