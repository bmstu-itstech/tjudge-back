package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/api/http_api"
	"github.com/bmstu-itstech/tjudge-back/internal/app"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/internal/services"
	"github.com/bmstu-itstech/tjudge-back/internal/utils/postgres"
	"github.com/bmstu-itstech/tjudge-back/pkg/server"
	"github.com/go-chi/chi/v5"
)

func main() {
	db, closeFn := postgres.ConnectToDatabase()
	defer closeFn()
	userRep := services.NewPostgresUserRepository(db)
	teamRep := services.NewPostgresTeamRepository(db, userRep)
	teamFactory, _ := tjudge.NewFixedSizeFactory(3)

	secretKey := os.Getenv("SECRET_KEY")
	durationS, err := strconv.Atoi(os.Getenv("TOKEN_DURATION_S"))
	if err != nil {
		log.Fatalf("expected TOKEN_DURATION_S is valid integer, got %s", os.Getenv("TOKEN_DURATION_S"))
	}
	duration := time.Duration(durationS) * time.Second
	generator := services.NewJWTTokenGenerator(secretKey, duration)

	app := app.App{
		UserLogin:    app.NewLoginHandler(userRep, generator),
		UserRegister: app.NewRegisterHandler(userRep),
		UserDelete:   app.NewDeleteHandler(userRep),
		UserGet:      app.NewGetUserHandler(userRep),
		UsersGet:     app.NewGetUsersHandler(userRep),
		
		TeamRegister: app.NewRegisterTeamHandler(teamRep, teamFactory),
		TeamsGet:     app.NewGetTeamsHandler(teamRep),
		TeamGet:      app.NewGetTeamHandler(teamRep),
		TeamDelete:   app.NewDeleteTeamHandler(teamRep),
	}

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return http_api.HandlerFromMux(http_api.NewHTTPServer(&app), router)
	})
}