package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type GetTours struct {
	GameId string
}

type GetToursHandler decorator.QueryHandler[GetTours, []Tour]

type getToursHandler struct {
	r tjudge.TourRepository
}

func (h getToursHandler) Handle(ctx context.Context, q GetTours) ([]Tour, error) {
	tours, err := h.r.Tours(ctx, tjudge.GameId(q.GameId))
	if err != nil {
		return nil, err
	}
	return batchToursToDto(tours), nil
}

func NewGetToursHandler(
	r tjudge.TourRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetToursHandler {
	return decorator.ApplyQueryDecorators(getToursHandler{r}, l, mc)
}
