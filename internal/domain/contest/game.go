package contest

import "github.com/bmstu-itstech/tjudge-back/internal/domain/shared"

type Game struct {
	Id       shared.ID
	Name     string
	RulesUrl string

	Matches map[shared.ID]*Match
}
