package commands

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
	"github.com/google/uuid"
)

type DeleteContest struct {
	Id string
}

type DeleteContestHandler decorator.CommandHandler[DeleteContest]

type deleteContestHandler struct {
	cr ports.ContestRepository
	tr ports.TeamRepository
}

func (h deleteContestHandler) Handle(ctx context.Context, cmd DeleteContest) error {
	uid, err := uuid.Parse(cmd.Id)
	if err != nil {
		return err
	}
	id := shared.ID(uid)
	// assuming we also want to delete all related teams
	teams, err := h.tr.ByContest(ctx, id)
	if err != nil {
		return err
	}
	for _, t := range teams {
		if err := h.tr.Delete(ctx, t.ID); err != nil {
			return err
		}
	}
	return h.cr.Delete(ctx, id)
}

func NewDeleteContestHandler(
	cr ports.ContestRepository,
	tr ports.TeamRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) DeleteContestHandler {
	return decorator.ApplyCommandDecorators(deleteContestHandler{cr, tr}, l, mc)
}
