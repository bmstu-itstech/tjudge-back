package decorator

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

func ApplyConsumerDecorators[E shared.Event](
	handler EventConsumerHandler[E],
	logger *slog.Logger,
	metricsClient MetricsClient,
) EventConsumer {
	return consumerFilter[E]{
		base: consumerLoggingDecorator[E]{
			base: consumerMetricsDecorator[E]{
				base:   handler,
				client: metricsClient,
			},
			logger: logger,
		}}
}

type consumerFilter[E shared.Event] struct {
	base EventConsumerHandler[E]
}

func (c consumerFilter[E]) Handle(ctx context.Context, e shared.Event) error {
	e_cast, ok := e.(E)
	if !ok {
		return nil
	}
	return c.base.Handle(ctx, e_cast)
}

type EventConsumer interface {
	Handle(ctx context.Context, e shared.Event) error
}

type EventConsumerHandler[E shared.Event] interface {
	Handle(ctx context.Context, e E) error
}
