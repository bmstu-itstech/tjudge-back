package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type GetProgramSource struct {
	ProgramId string
}

type GetProgramSourceHandler decorator.QueryHandler[GetProgramSource, ProgramSource]

type getProgramSourceHandler struct {
	r tjudge.ProgramSourceRepository
}

func (h getProgramSourceHandler) Handle(ctx context.Context, q GetProgramSource) (ProgramSource, error) {
	source, err := h.r.ProgramSource(ctx, tjudge.ProgramId(q.ProgramId))
	if err != nil {
		return ProgramSource{}, err
	}
	return sourceToDto(source), nil
}

func NewGetProgramSourceHandler(
	r tjudge.ProgramSourceRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetProgramSourceHandler {
	return decorator.ApplyQueryDecorators(getProgramSourceHandler{r}, l, mc)
}
