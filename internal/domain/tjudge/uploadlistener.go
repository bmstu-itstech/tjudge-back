package tjudge

import "context"

type UploadListener interface {
	OnUpload(context.Context, GameId) error
}
