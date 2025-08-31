package http_api

import (
	"errors"
	"net/http"

	"github.com/bmstu-itstech/tjudge-back/internal/app"
	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/jwtauth"
	"github.com/go-chi/render"
)

type HTTPServer struct {
	app *app.App
}

func NewHTTPServer(app *app.App) HTTPServer {
	return HTTPServer{app: app}
}

func (s HTTPServer) PostLogin(w http.ResponseWriter, r *http.Request) {
	userLogin := UserLogin{}
	if err := render.Decode(r, &userLogin); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	token, err := s.app.UserLogin.Execute(r.Context(), app.Login{Username: userLogin.Username, Password: userLogin.Password})
	if errors.Is(err, tjudge.ErrUserNotFound) {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	if errors.Is(err, tjudge.ErrPasswordMismatch) {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, AuthToken{AccessToken: token})
}

func (s HTTPServer) PostRegister(w http.ResponseWriter, r *http.Request) {
	userRegister := UserRegister{}
	if err := render.Decode(r, &userRegister); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	err := s.app.UserRegister.Execute(r.Context(), app.Register{Username: userRegister.Username, Fullname: userRegister.Fullname, Password: userRegister.Password})
	if errors.Is(err, tjudge.ErrInvalidInput) {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	if errors.Is(err, tjudge.ErrUserAlreadyExists) {
		httpError(w, r, err, http.StatusConflict)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s HTTPServer) DeleteUser(w http.ResponseWriter, r *http.Request, id int64) {
	userID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	err = s.app.UserDelete.Execute(r.Context(), app.Delete{AuthorID: int(userID), Id: id})
	if errors.Is(err, tjudge.ErrUserNotFound) {
		httpError(w, r, err, http.StatusNotFound)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (s HTTPServer) GetUser(w http.ResponseWriter, r *http.Request, id int64) {
	_, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	user, err := s.app.UserGet.Execute(r.Context(), app.GetUser{Id: id})
	if errors.Is(err, tjudge.ErrUserNotFound) {
		httpError(w, r, err, http.StatusNotFound)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, convertUserToApi(user))
}

func (s HTTPServer) GetUsers(w http.ResponseWriter, r *http.Request) {
	_, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	users, err := s.app.UsersGet.Execute(r.Context())
	if errors.Is(err, tjudge.ErrUserNotFound) {
		httpError(w, r, err, http.StatusNotFound)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	apiUsers := make([]UserResponse, len(users))
	for i, user := range users {
		apiUsers[i] = convertUserToApi(*user)
	}
	render.JSON(w, r, apiUsers)
}

func (s HTTPServer) RegisterTeam(w http.ResponseWriter, r *http.Request) {
	teamRegister := TeamRegister{}
	if err := render.Decode(r, &teamRegister); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	err := s.app.TeamRegister.Execute(r.Context(), app.RegisterTeam{Name: teamRegister.Name, Code: teamRegister.Code, Contest: int(teamRegister.Contest), MaxSize: *teamRegister.MaxSize})
	if errors.Is(err, tjudge.ErrInvalidInput) {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	if errors.Is(err, tjudge.ErrTeamAlreadyExist) {
		httpError(w, r, err, http.StatusConflict)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s HTTPServer) DeleteTeam(w http.ResponseWriter, r *http.Request, code string) {
	userID, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	err = s.app.TeamDelete.Execute(r.Context(), app.DeleteTeam{AuthorID: int(userID), Code: code})
	if errors.Is(err, tjudge.ErrTeamNotFound) {
		httpError(w, r, err, http.StatusNotFound)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func (s HTTPServer) GetTeam(w http.ResponseWriter, r *http.Request, code string) {
	_, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	team, err := s.app.TeamGet.Execute(r.Context(), app.GetTeam{Code: code})
	if errors.Is(err, tjudge.ErrTeamNotFound) {
		httpError(w, r, err, http.StatusNotFound)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, convertTeamToApi(*team))
}

func (s HTTPServer) GetTeams(w http.ResponseWriter, r *http.Request) {
	_, err := jwtauth.UserIDFromContext(r.Context())
	if err != nil {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	}
	teams, err := s.app.TeamsGet.Execute(r.Context())
	if errors.Is(err, tjudge.ErrTeamNotFound) {
		httpError(w, r, err, http.StatusNotFound)
		return
	}
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	apiTeams := make([]Team, len(teams))
	for i, team := range teams {
		apiTeams[i] = convertTeamToApi(*team)
	}
	render.JSON(w, r, apiTeams)
}

func convertUserToApi(user app.User) UserResponse {
	return UserResponse{
		Id:       int(user.Id),
		Username: user.Username,
		Fullname: user.Fullname,
	}
}

func convertTeamToApi(team app.Team) Team {
	amembers := make([]UserResponse, len(team.Members))
	for i, member := range team.Members {
		amembers[i] = convertUserToApi(*member)
	}
	var leader UserResponse
	if team.Leader != nil {
		leader = convertUserToApi(*team.Leader)
	} else {
		leader = UserResponse{}
	}
	return Team{
		Code:    team.Code,
		Name:    team.Name,
		Leader:  &leader,
		Contest: int64(team.Contest),
		MaxSize: team.MaxSize,
		Members: amembers,
	}
}

func httpError(w http.ResponseWriter, r *http.Request, err error, code int) {
	render.Status(r, code)
	render.JSON(w, r, Error{Message: err.Error()})
}
