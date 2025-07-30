package app

type Commands struct {
	UploadProgramHandler
	CreateContestHandler
	UpsertContestHandler
	CreateGameHandler
	UpsertGameHandler
}

type Queries struct {
	GetProgramHandler
	GetProgramsHandler
	ActiveProgramHandler
	GetProgramSourceHandler
	GetContestHandler
	ActiveContestsHandler
	GetGameHandler
	GetTourHandler
	GetToursHandler
	ActiveTourHandler
	GetRoundHandler
	GetTeamHandler
	CreateTeamHandler
	GetPlayerHandler
}

type Application struct {
	Commands Commands
	Queries  Queries
}
