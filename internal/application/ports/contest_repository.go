package ports

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/contest"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type ContestRepository interface {
	Save(ctx context.Context, c *contest.Contest) error
	Contest(ctx context.Context, id shared.ID) (*contest.Contest, error)
}
