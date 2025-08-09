package commands

import (
	"context"
	"log/slog"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/contest"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
	"github.com/google/uuid"
)

type UpsertContest struct {
	Id      string
	Name    string
	Starts  time.Time
	Ends    time.Time
	GameIds []string
	TeamIds []string
}

type UpsertContestHandler decorator.CommandHandler[UpsertContest]

type upsertContestHandler struct {
	cr ports.ContestRepository
	gr ports.GameRepository
	tr ports.TeamRepository
}

func (h upsertContestHandler) Handle(ctx context.Context, cmd UpsertContest) error {
	ids, err := idsFromDto(cmd.GameIds)
	if err != nil {
		return err
	}
	games := make(map[shared.ID]*contest.Game)
	for _, id := range ids {
		game, err := h.gr.Game(ctx, id)
		if err != nil {
			return err
		}
		games[id] = game
	}
	ids, err = idsFromDto(cmd.TeamIds)
	if err != nil {
		return err
	}
	teams := make(map[shared.ID]*contest.Team)
	for _, id := range ids {
		team, err := h.tr.Team(ctx, id)
		if err != nil {
			return err
		}
		teams[id] = team
	}
	id, err := uuid.Parse(cmd.Id)
	if err != nil {
		return err
	}
	contest, err := contest.ParseContest(
		shared.ID(id),
		cmd.Name,
		cmd.Starts,
		cmd.Ends,
		games,
		teams,
	)
	if err != nil {
		return err
	}
	return h.cr.Upsert(ctx, &contest)
}

func NewUpsertContestHandler(
	cr ports.ContestRepository,
	gr ports.GameRepository,
	tr ports.TeamRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) UpsertContestHandler {
	return decorator.ApplyCommandDecorators(upsertContestHandler{cr, gr, tr}, l, mc)
}
