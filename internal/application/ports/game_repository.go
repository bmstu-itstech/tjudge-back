package ports

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/contest"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

// MEMO: This wouldn't be required if games were *directly*
// tied to contests 1:1?
type GameRepository interface {
	Upsert(ctx context.Context, g *contest.Game) error
	Game(ctx context.Context, id shared.ID) (*contest.Game, error)
	All(ctx context.Context) ([]*contest.Game, error)
	Active(ctx context.Context) ([]*contest.Contest, error)
	Delete(ctx context.Context, id shared.ID) error
}
