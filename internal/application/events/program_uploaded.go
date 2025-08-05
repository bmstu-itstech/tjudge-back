package events

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/program"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type ProgramUploadedConsumer decorator.EventConsumer

type programUploadedHandler struct {
	repos     ports.ContestRepository
	publisher ports.EventPublisher
}

func (h programUploadedHandler) Handle(ctx context.Context, event program.UploadedEvent) error {
	c, err := h.repos.Contest(ctx, event.ContestID)
	if err != nil {
		return err
	}

	evs, err := c.ScheduleMatchesFor(event.TeamID, event.GameID)
	if err != nil {
		return err
	}

	err = h.repos.Upsert(ctx, c)
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

func NewProgramUploadedConsumer(
	repos ports.ContestRepository,
	publisher ports.EventPublisher,
	l *slog.Logger,
	mc decorator.MetricsClient,
) decorator.EventConsumer {
	return decorator.ApplyConsumerDecorators(programUploadedHandler{repos, publisher}, l, mc)
}
