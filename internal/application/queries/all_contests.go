package queries

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type AllContests struct{}

type AllContestsHandler decorator.QueryHandler[AllContests, []Contest]

type allContestsHandler struct {
	repo ports.ContestRepository
}

func (h allContestsHandler) Handle(ctx context.Context, a AllContests) ([]Contest, error) {
	contests, err := h.repo.All(ctx)
	if err != nil {
		return nil, err
	}
	return batchContestsToDto(contests), nil
}

func NewAllContestsHandler(
	r ports.ContestRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) AllContestsHandler {
	return decorator.ApplyQueryDecorators(allContestsHandler{r}, l, mc)
}
