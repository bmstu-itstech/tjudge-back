package app

type Commands struct {
	CreateContestHandler
	UpsertContestHandler
}

type Queries struct {
	GetContestHandler
	ActiveContestsHandler
}

type Application struct {
	Commands Commands
	Queries  Queries
}
