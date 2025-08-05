package ports

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/contest"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/program"
)

type Judge interface {
	Run(ctx context.Context, p1 program.Path, p2 program.Path) (contest.Result, contest.Result, error)
}
