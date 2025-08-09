package commands

import (
	"context"
	"log/slog"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/contest"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type CreateContest struct {
	Name    string
	Starts  time.Time
	Ends    time.Time
	GameIds []string
}

type CreateContestHandler decorator.CommandHandler[CreateContest]

type createContestHandler struct {
	cr ports.ContestRepository
	gr ports.GameRepository
}

func (h createContestHandler) Handle(ctx context.Context, cmd CreateContest) error {
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
	contest, err := contest.NewContest(
		cmd.Name,
		cmd.Starts,
		cmd.Ends,
		games,
	)
	if err != nil {
		return err
	}
	return h.cr.Upsert(ctx, &contest)
}

func NewCreateContestHandler(
	cr ports.ContestRepository,
	gr ports.GameRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) CreateContestHandler {
	return decorator.ApplyCommandDecorators(createContestHandler{cr, gr}, l, mc)
}
