package tjudge

import (
	"context"
	"time"
)

// What happens when a program crashes? Who cares.
// 0 points, they'll figure it out?

type RoundEvent struct {
	Type     string // "scheduled"/"finished"/"aborted"/"timed_out"
	GameId   GameId
	Programs map[TeamId]ProgramId
	Results  map[TeamId]int
	Time     time.Time
}

type RoundEventConsumer interface {
	HandleEvent(context.Context, RoundEvent) error
}

type RoundEventStore interface {
	Append(context.Context, RoundEvent) error
	Load(ctx context.Context, game GameId) ([]RoundEvent, error)
	Subscribe(consumer RoundEventConsumer) error
}
