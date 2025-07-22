package slogpretty

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"
)

type PlainHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

// PlainHandler is the PrettyHandler we all know with colours removed
type PlainHandler struct {
	opts PlainHandlerOptions
	slog.Handler
	l     *stdLog.Logger
	attrs []slog.Attr
}

func (opts PlainHandlerOptions) NewPlainHandler(
	out io.Writer,
) *PlainHandler {
	h := &PlainHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		l:       stdLog.New(out, "", 0),
	}

	return h
}

func (h *PlainHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[02/01 15:04:05.000]")

	h.l.Println(
		timeStr,
		level,
		r.Message,
		string(b),
	)

	return nil
}

func (h *PlainHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PlainHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *PlainHandler) WithGroup(name string) slog.Handler {
	return &PlainHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}
