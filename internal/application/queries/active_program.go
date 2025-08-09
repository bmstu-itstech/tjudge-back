package queries

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/program"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
	"github.com/google/uuid"
)

type ActiveProgram struct {
	TeamId string
	GameId string
}

type ActiveProgramHandler decorator.QueryHandler[ActiveProgram, Program]

type activeProgramHandler struct {
	r  ports.ProgramRepository
	fs ports.FileStorage
}

func (h activeProgramHandler) Handle(ctx context.Context, q ActiveProgram) (Program, error) {
	teamId, err := uuid.Parse(q.TeamId)
	if err != nil {
		return Program{}, err
	}
	gameId, err := uuid.Parse(q.GameId)
	if err != nil {
		return Program{}, err
	}

	prog, ok, err := h.r.Active(ctx, shared.ID(teamId), shared.ID(gameId))
	if err != nil {
		return Program{}, err
	} else if !ok {
		return Program{}, program.ErrNoActiveProgram
	}

	file, err := h.fs.Read(ctx, prog.Path)
	if err != nil {
		return Program{}, err
	}

	return programToDto(*prog, file), nil
}

func NewActiveProgramHandler(
	r ports.ProgramRepository,
	fs ports.FileStorage,
	l *slog.Logger,
	mc decorator.MetricsClient,
) ActiveProgramHandler {
	return decorator.ApplyQueryDecorators(activeProgramHandler{r, fs}, l, mc)
}
