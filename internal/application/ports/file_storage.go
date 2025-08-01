package ports

import (
	"context"
	"io"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/program"
)

type FileStorage interface {
	Upload(ctx context.Context, r io.Reader) (program.Path, error)
	Delete(ctx context.Context, path program.Path) error
}
