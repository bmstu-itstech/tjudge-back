package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type ActiveContests struct{}

type ActiveContestsHandler decorator.QueryHandler[ActiveContests, []Contest]

type activeContestsHandler struct {
	r tjudge.ContestRepository
}

func (h activeContestsHandler) Handle(ctx context.Context, a ActiveContests) ([]Contest, error) {
	contests, err := h.r.Active(ctx)
	if err != nil {
		return nil, err
	}
	return batchContestsToDto(contests), nil
}

func NewActiveContestsHandler(
	r tjudge.ContestRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) ActiveContestsHandler {
	return decorator.ApplyQueryDecorators(activeContestsHandler{r}, l, mc)
}
