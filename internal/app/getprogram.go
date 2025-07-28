package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type GetProgram struct {
	ProgramId string
}

type GetProgramHandler decorator.QueryHandler[GetProgram, Program]

type getProgramHandler struct {
	r tjudge.ProgramRepository
}

func (h getProgramHandler) Handle(ctx context.Context, q GetProgram) (Program, error) {
	program, err := h.r.Program(ctx, tjudge.ProgramId(q.ProgramId))
	if err != nil {
		return Program{}, err
	}
	return programToDto(program), nil
}

func NewGetProgramHandler(
	r tjudge.ProgramRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetProgramHandler {
	return decorator.ApplyQueryDecorators(getProgramHandler{r}, l, mc)
}
