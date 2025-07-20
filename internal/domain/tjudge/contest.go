package tjudge

import "time"

type ContestId string

// Объединение игр
type Contest struct {
	Id        ContestId
	Name      string
	TeamLimit uint
	Start     time.Time
	End       time.Time
	Games     []GameId
}
