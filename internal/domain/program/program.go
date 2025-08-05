package program

import (
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type Path string

type Program struct {
	ID        shared.ID
	ContestID shared.ID
	GameID    shared.ID
	TeamID    shared.ID
	Path      Path
}

func New(teamID shared.ID, contestID shared.ID, gameID shared.ID, path Path) (*Program, UploadedEvent, error) {
	p := &Program{
		ID:     shared.NewID(),
		TeamID: teamID,
		Path:   path,
	}

	e := UploadedEvent{
		ContestID: contestID,
		GameID:    gameID,
		TeamID:    teamID,
		ProgramID: p.ID,
		Timestamp: time.Now(),
	}

	return p, e, nil
}
