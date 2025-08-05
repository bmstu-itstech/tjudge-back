package contest

import (
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type Team struct {
	ID        shared.ID
	ContestID shared.ID
	Name      shared.ID
	CreatedAt time.Time
	JoinCode  string
}
