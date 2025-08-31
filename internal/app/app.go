package app

type App struct {
	UserLogin LoginHandler
	UserRegister RegisterHandler 
	UserDelete DeleteHandler
	UserGet GetUserHandler
	UsersGet GetUsersHandler

	TeamRegister RegisterTeamHandler
	TeamGet *GetTeamHandler
	TeamsGet *GetTeamsHandler
	TeamDelete DeleteTeamHandler
}