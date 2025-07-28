package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type CreateContest struct {
	Name      string
	TeamLimit uint
	Starts    time.Time
	Ends      time.Time
	GameIds   []string
}

type CreateContestHandler decorator.CommandHandler[CreateContest]

type createContestHandler struct {
	r tjudge.ContestRepository
}

func (h createContestHandler) Handle(ctx context.Context, cmd CreateContest) error {
	contest, err := tjudge.NewContest(
		cmd.Name,
		cmd.TeamLimit,
		cmd.Starts, cmd.Ends,
		gameIdsFromDto(cmd.GameIds))
	if err != nil {
		return err
	}
	return h.r.Upsert(ctx, contest)
}

func NewCreateContestHandler(
	r tjudge.ContestRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) CreateContestHandler {
	return decorator.ApplyCommandDecorators(createContestHandler{r}, l, mc)
}
