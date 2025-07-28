package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type ActiveTour struct {
	GameId string
}

type ActiveTourHandler decorator.QueryHandler[ActiveTour, Tour]

type activeTourHandler struct {
	r tjudge.TourRepository
}

func (h activeTourHandler) Handle(ctx context.Context, q ActiveTour) (Tour, error) {
	tour, err := h.r.Active(ctx, tjudge.GameId(q.GameId))
	if err != nil {
		return Tour{}, err
	}
	return tourToDto(tour), nil
}

func NewActiveTourHandler(
	r tjudge.TourRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) ActiveTourHandler {
	return decorator.ApplyQueryDecorators(activeTourHandler{r}, l, mc)
}
