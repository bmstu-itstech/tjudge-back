package queries

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
	"github.com/google/uuid"
)

type GetContest struct {
	ContestId string
}

type GetContestHandler decorator.QueryHandler[GetContest, Contest]

type getContestHandler struct {
	r ports.ContestRepository
}

func (h getContestHandler) Handle(ctx context.Context, q GetContest) (Contest, error) {
	id, err := uuid.Parse(q.ContestId)
	if err != nil {
		return Contest{}, err
	}
	contest, err := h.r.Contest(ctx, shared.ID(id))
	if err != nil {
		return Contest{}, err
	}
	return contestToDto(contest), nil
}

func NewGetContestHandler(
	r ports.ContestRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetContestHandler {
	return decorator.ApplyQueryDecorators(getContestHandler{r}, l, mc)
}
