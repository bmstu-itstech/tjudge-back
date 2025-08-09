package program

import (
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type UploadedEvent struct {
	ContestId shared.ID
	GameId    shared.ID
	TeamId    shared.ID
	ProgramId shared.ID
	Timestamp time.Time
}

func (e UploadedEvent) Name() string {
	return "program_uploaded_event"
}

func (e UploadedEvent) IsEvent() {}
