package contest

import (
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
)

type Score int

type Team struct {
	ID    shared.ID
	Score Score
}

func (t *Team) ResetScore() {
	t.Score = 0
}

func (t *Team) AddScore(score Score) {
	t.Score += score
}
