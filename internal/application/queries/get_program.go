package queries

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
	"github.com/google/uuid"
)

type GetProgram struct {
	Id string
}

type GetProgramHandler decorator.QueryHandler[GetProgram, Program]

type getProgramHandler struct {
	r  ports.ProgramRepository
	fs ports.FileStorage
}

func (h getProgramHandler) Handle(ctx context.Context, q GetProgram) (Program, error) {
	id, err := uuid.Parse(q.Id)
	if err != nil {
		return Program{}, err
	}

	prog, err := h.r.Program(ctx, shared.ID(id))
	if err != nil {
		return Program{}, err
	}

	file, err := h.fs.Read(ctx, prog.Path)
	if err != nil {
		return Program{}, err
	}

	return programToDto(*prog, file), nil
}

func NewGetProgramHandler(
	r ports.ProgramRepository,
	fs ports.FileStorage,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetProgramHandler {
	return decorator.ApplyQueryDecorators(getProgramHandler{r, fs}, l, mc)
}
