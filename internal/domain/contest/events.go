package contest

import (
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type ScheduledMatchEvent struct {
	MatchID   shared.ID
	ContestID shared.ID
	GameID    shared.ID
	Team1ID   shared.ID
	Team2ID   shared.ID
	Timestamp time.Time
}

func (e ScheduledMatchEvent) Name() string {
	return "contest_scheduled_match_event"
}

func (e ScheduledMatchEvent) IsEvent() {}
