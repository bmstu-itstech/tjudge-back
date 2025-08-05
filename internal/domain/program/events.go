package program

import (
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type UploadedEvent struct {
	ContestID shared.ID
	GameID    shared.ID
	TeamID    shared.ID
	ProgramID shared.ID
	Timestamp time.Time
}

func (e UploadedEvent) Name() string {
	return "program_uploaded_event"
}

func (e UploadedEvent) IsEvent() {}
