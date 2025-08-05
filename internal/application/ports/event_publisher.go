package ports

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type EventPublisher interface {
	Publish(ctx context.Context, event shared.Event) error
	Subscribe(ctx context.Context, name string, consumer decorator.EventConsumer) error
}
