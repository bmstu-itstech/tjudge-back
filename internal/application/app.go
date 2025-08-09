package application

import (
	"github.com/bmstu-itstech/tjudge-back/internal/application/commands"
	"github.com/bmstu-itstech/tjudge-back/internal/application/queries"
)

type Commands struct {
	CreateGame    commands.CreateGameHandler
	CreateContest commands.CreateContestHandler
	DeleteContest commands.DeleteContestHandler
	DeleteGame    commands.DeleteGameHandler
	UploadProgram commands.UploadProgramHandler
	UpsertContest commands.UpsertContestHandler
	UpsertGame    commands.UpsertGameHandler
}

type Queries struct {
	ActiveContests queries.ActiveContestsHandler
	ActiveProgram  queries.ActiveProgramHandler
	AllContests    queries.AllContestsHandler
	AllGames       queries.AllGamesHandler
	GetContest     queries.GetContestHandler
	GetGame        queries.GetGameHandler
	GetProgram     queries.GetProgramHandler
	GetStandings   queries.GetStandingsHandler
}

type App struct {
	Commands Commands
	Queries  Queries
}
