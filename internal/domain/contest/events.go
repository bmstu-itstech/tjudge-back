package contest

import (
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type MatchScheduledEvent struct {
	MatchID   shared.ID
	ContestID shared.ID
	GameID    shared.ID
	Team1ID   shared.ID
	Team2ID   shared.ID
	Timestamp time.Time
}

func (e MatchScheduledEvent) Name() string {
	return "contest_match_scheduled_event"
}

func (e MatchScheduledEvent) IsEvent() {}
