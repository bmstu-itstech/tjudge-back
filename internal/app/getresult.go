package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type GetResult struct {
	ResultId string
}

type GetResultHandler decorator.QueryHandler[GetResult, Result]

type getResultHandler struct {
	r tjudge.ResultRepository
}

func (h getResultHandler) Handle(ctx context.Context, q GetResult) (Result, error) {
	result, err := h.r.Result(ctx, tjudge.ResultId(q.ResultId))
	if err != nil {
		return Result{}, err
	}
	return resultToDto(result), nil
}

func NewGetResultHandler(
	r tjudge.ResultRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) GetResultHandler {
	return decorator.ApplyQueryDecorators(getResultHandler{r}, l, mc)
}
