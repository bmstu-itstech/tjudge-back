package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type GetContest struct {
	ContestId string
}

type GetContestHandler decorator.QueryHandler[GetContest, Contest]

type getContestHandler struct {
	r tjudge.ContestRepository
}

func (h getContestHandler) Handle(ctx context.Context, q GetContest) (Contest, error) {
	contest, err := h.r.Contest(ctx, tjudge.ContestId(q.ContestId))
	if err != nil {
		return Contest{}, err
	}
	return contestToDto(contest), nil
}

func NewGetContestHandler(
	r tjudge.ContestRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetContestHandler {
	return decorator.ApplyQueryDecorators(getContestHandler{r}, l, mc)
}
