package ports

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/contest"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

// This wouldn't be required if games were *directly*
// tied to contests 1:1? But I don't think we want that.
type GameRepository interface {
	Upsert(ctx context.Context, g *contest.Game) error
	Game(ctx context.Context, id shared.ID) (*contest.Game, error)
	All(ctx context.Context) ([]*contest.Game, error)
	Delete(ctx context.Context, id shared.ID) error
}
