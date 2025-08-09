package queries

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type ActiveContests struct{}

type ActiveContestsHandler decorator.QueryHandler[ActiveContests, []Contest]

type activeContestsHandler struct {
	repo ports.ContestRepository
}

func (h activeContestsHandler) Handle(ctx context.Context, a ActiveContests) ([]Contest, error) {
	contests, err := h.repo.Active(ctx)
	if err != nil {
		return nil, err
	}
	return batchContestsToDto(contests), nil
}

func NewActiveContestsHandler(
	r ports.ContestRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) ActiveContestsHandler {
	return decorator.ApplyQueryDecorators(activeContestsHandler{r}, l, mc)
}
