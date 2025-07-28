package app

type Commands struct {
	CreateContestHandler
	UpsertContestHandler
	CreateGameHandler
	UpsertGameHandler
}

type Queries struct {
	GetContestHandler
	ActiveContestsHandler
	GetGameHandler
	GetTourHandler
	GetToursHandler
	ActiveTourHandler
	GetResultHandler
	GetRoundHandler
	GetTeamHandler
	CreateTeamHandler
	GetPlayerHandler
}

type Application struct {
	Commands Commands
	Queries  Queries
}
