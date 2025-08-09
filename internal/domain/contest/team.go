package contest

import (
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type Team struct {
	Id        shared.ID
	ContestId shared.ID
	Name      shared.ID
	CreatedAt time.Time
	JoinCode  string
}
