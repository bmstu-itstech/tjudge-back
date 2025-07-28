package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type GetPrograms struct {
	TeamId string
	GameId string
}

type GetProgramsHandler decorator.QueryHandler[GetPrograms, []Program]

type getProgramsHandler struct {
	r tjudge.ProgramRepository
}

func (h getProgramsHandler) Handle(ctx context.Context, q GetPrograms) ([]Program, error) {
	programs, err := h.r.Programs(ctx, tjudge.GameId(q.GameId), tjudge.TeamId(q.TeamId))
	if err != nil {
		return nil, err
	}
	return batchProgramsToDto(programs), nil
}

func NewGetProgramsHandler(
	r tjudge.ProgramRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetProgramsHandler {
	return decorator.ApplyQueryDecorators(getProgramsHandler{r}, l, mc)
}
