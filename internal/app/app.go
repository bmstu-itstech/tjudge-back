package app

type Commands struct {
	UploadProgramHandler
	CreateContestHandler
	UpsertContestHandler
	CreateGameHandler
	UpsertGameHandler
}

type Queries struct {
	GetStandingsHandler
	GetProgramHandler
	GetProgramsHandler
	ActiveProgramHandler
	GetProgramSourceHandler
	GetContestHandler
	ActiveContestsHandler
	GetGameHandler
	GetTeamHandler
	CreateTeamHandler
	GetPlayerHandler
}

type Application struct {
	Commands Commands
	Queries  Queries
}
