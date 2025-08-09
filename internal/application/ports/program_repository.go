package ports

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/program"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type ProgramRepository interface {
	Upsert(ctx context.Context, p *program.Program) error
	Program(ctx context.Context, id shared.ID) (*program.Program, error)
	Active(ctx context.Context, teamId shared.ID, gameId shared.ID) (*program.Program, bool, error)
}
