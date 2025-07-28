package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type GetRound struct {
	RoundId string
}

type GetRoundHandler decorator.QueryHandler[GetRound, Round]

type getRoundHandler struct {
	r tjudge.RoundRepository
}

func (h getRoundHandler) Handle(ctx context.Context, q GetRound) (Round, error) {
	round, err := h.r.Round(ctx, tjudge.RoundId(q.RoundId))
	if err != nil {
		return Round{}, err
	}
	return roundToDto(round), nil
}

func NewGetRoundHandler(
	r tjudge.RoundRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetRoundHandler {
	return decorator.ApplyQueryDecorators(getRoundHandler{r}, l, mc)
}
