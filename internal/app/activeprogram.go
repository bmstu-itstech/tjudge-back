package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type ActiveProgram struct {
	TeamId string
	GameId string
}

type ActiveProgramHandler decorator.QueryHandler[ActiveProgram, Program]

type activeProgramHandler struct {
	r tjudge.ProgramRepository
}

func (h activeProgramHandler) Handle(ctx context.Context, q ActiveProgram) (Program, error) {
	program, err := h.r.Active(ctx, tjudge.GameId(q.GameId), tjudge.TeamId(q.TeamId))
	if err != nil {
		return Program{}, err
	}
	return programToDto(program), nil
}

func NewActiveProgramHandler(
	r tjudge.ProgramRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) ActiveProgramHandler {
	return decorator.ApplyQueryDecorators(activeProgramHandler{r}, l, mc)
}
