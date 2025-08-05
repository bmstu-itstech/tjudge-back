package ports

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/program"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type ProgramRepository interface {
	Save(ctx context.Context, p *program.Program) error
	Program(ctx context.Context, id shared.ID) (*program.Program, error)
	LastTeamProgram(ctx context.Context, teamID shared.ID, game shared.ID) (*program.Program, bool, error)
}
