package ports

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type EventPublisher interface {
	Publish(ctx context.Context, event shared.Event) error
}
