package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type GetTour struct {
	TourId string
}

type GetTourHandler decorator.QueryHandler[GetTour, Tour]

type getTourHandler struct {
	r tjudge.TourRepository
}

func (h getTourHandler) Handle(ctx context.Context, q GetTour) (Tour, error) {
	tour, err := h.r.Tour(ctx, tjudge.TourId(q.TourId))
	if err != nil {
		return Tour{}, err
	}
	return tourToDto(tour), nil
}

func NewGetTourHandler(
	r tjudge.TourRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetTourHandler {
	return decorator.ApplyQueryDecorators(getTourHandler{r}, l, mc)
}
