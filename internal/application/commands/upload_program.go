package commands

import (
	"context"
	"io"

	"github.com/bmstu-itstech/tjudge-back/internal/application/ports"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/program"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type UploadProgram struct {
	ContestID shared.ID
	TeamID    shared.ID
	GameID    shared.ID
	Reader    io.Reader
}

type UploadProgramHandler struct {
	storage   ports.FileStorage
	repos     ports.ProgramRepository
	publisher ports.EventPublisher
}

func (h *UploadProgramHandler) Handle(ctx context.Context, cmd UploadProgram) error {
	path, err := h.storage.Upload(ctx, cmd.Reader)
	if err != nil {
		return err
	}

	p, ev, err := program.New(cmd.ContestID, cmd.GameID, cmd.TeamID, path)
	if err != nil {
		_ = h.storage.Delete(ctx, path)
		return err
	}

	if err = h.repos.Upsert(ctx, p); err != nil {
		_ = h.storage.Delete(ctx, path)
		return err
	}

	return h.publisher.Publish(ctx, ev)
}
