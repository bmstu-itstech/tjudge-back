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
}

type Application struct {
	Commands Commands
	Queries  Queries
}
