package ports

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/contest"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type TeamRepository interface {
	Upsert(context.Context, *contest.Team) error
	Team(ctx context.Context, id shared.ID) (*contest.Team, error)
	ByContest(ctx context.Context, id shared.ID) ([]*contest.Team, error)
	ByJoinCode(ctx context.Context, code string) (*contest.Team, error)
	Delete(ctx context.Context, id shared.ID) error
}
