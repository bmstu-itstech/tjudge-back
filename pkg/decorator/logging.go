package decorator

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}

type commandLoggingDecorator[C any] struct {
	base   CommandHandler[C]
	logger *slog.Logger
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerType := generateActionName(cmd)

	logger := d.logger.With(
		slog.String("command", handlerType),
		slog.String("command_body", fmt.Sprintf("%v", cmd)),
	)

	logger.Debug("Executing command")
	defer func() {
		if err == nil {
			logger.Info("Command executed successfully")
		} else {
			logger.Error("Failed to execute command", "error", err.Error())
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type queryLoggingDecorator[C any, R any] struct {
	base   QueryHandler[C, R]
	logger *slog.Logger
}

func (d queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	logger := d.logger.With(
		slog.String("query", generateActionName(cmd)),
		slog.String("query_body", fmt.Sprintf("%v", cmd)),
	)

	logger.Debug("Executing query")
	defer func() {
		if err == nil {
			logger.Info("Query executed successfully")
		} else {
			logger.Error("Failed to execute query", "error", err.Error())
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type consumerLoggingDecorator[E shared.Event] struct {
	base   EventConsumerHandler[E]
	logger *slog.Logger
}

func (d consumerLoggingDecorator[E]) Handle(ctx context.Context, e E) (err error) {
	logger := d.logger.With(
		slog.String("consumer", generateActionName(e)),
		slog.String("consumer_body", fmt.Sprintf("%v", e)),
	)

	logger.Debug("Consuming event")
	defer func() {
		if err == nil {
			logger.Info("Event successfully consumed")
		} else {
			logger.Error("Failed to consume event", "error", err.Error())
		}
	}()

	return d.base.Handle(ctx, e)
}
