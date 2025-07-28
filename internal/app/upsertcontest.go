package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type UpsertContest struct {
	ContestId string
	Name      string
	TeamLimit uint
	Starts    time.Time
	Ends      time.Time
	GameIds   []string
}

type UpsertContestHandler decorator.CommandHandler[UpsertContest]

type upsertContestHandler struct {
	r tjudge.ContestRepository
}

func (h upsertContestHandler) Handle(ctx context.Context, cmd UpsertContest) error {
	contest, err := tjudge.ParseContest(
		tjudge.ContestId(cmd.ContestId),
		cmd.Name,
		cmd.TeamLimit,
		cmd.Starts, cmd.Ends,
		gameIdsFromDto(cmd.GameIds)) // TODO: verify games exist?
	if err != nil {
		return err
	}
	return h.r.Upsert(ctx, contest)
}

func NewUpsertContestHandler(
	r tjudge.ContestRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) UpsertContestHandler {
	return decorator.ApplyCommandDecorators(upsertContestHandler{r}, l, mc)
}
