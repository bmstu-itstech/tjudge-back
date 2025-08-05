package ports

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/contest"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type ContestRepository interface {
	Upsert(ctx context.Context, c *contest.Contest) error
	Contest(ctx context.Context, id shared.ID) (*contest.Contest, error)
	All(ctx context.Context) ([]*contest.Contest, error)
	Active(ctx context.Context) ([]*contest.Contest, error)
	Delete(ctx context.Context, id shared.ID) error
}
