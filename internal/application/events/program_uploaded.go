package events

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/program"
)

type ProgramUploadedHandler struct {
	repos     ports.ContestRepository
	publisher ports.EventPublisher
}

func (h *ProgramUploadedHandler) Handle(ctx context.Context, event program.UploadedEvent) error {
	c, err := h.repos.Contest(ctx, event.ContestID)
	if err != nil {
		return err
	}

	evs, err := c.ScheduleMatchesFor(event.TeamID, event.GameID)
	if err != nil {
		return err
	}

	for _, ev := range evs {
		err = h.publisher.Publish(ctx, ev)
		if err != nil {
			return err
		}
	}

	return nil
}
