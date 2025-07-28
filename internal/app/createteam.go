package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

type CreateTeam struct {
	Name      string
	ContestId string
}

type CreateTeamHandler decorator.CommandHandler[CreateTeam]

type createTeamHandler struct {
	cfg tjudge.ConfigRepository
	tr  tjudge.TeamRepository
	cr  tjudge.ContestRepository
}

func (h createTeamHandler) Handle(ctx context.Context, cmd CreateTeam) error {
	// As a side effect, we confirm that the contest exists
	con, err := h.cr.Contest(ctx, tjudge.ContestId(cmd.ContestId))
	if err != nil {
		return err
	}
	cfg, err := h.cfg.Get(ctx)
	if err != nil {
		return err
	}
	if !cfg.AllowLateReg() && time.Now().After(con.Starts()) {
		return tjudge.ErrNoLateJoin
	}
	team, err := tjudge.NewTeam(cmd.Name, tjudge.ContestId(cmd.ContestId))
	if err != nil {
		return err
	}
	return h.tr.Upsert(ctx, team)
}

func NewCreateTeamHandler(
	cfg tjudge.ConfigRepository,
	tr tjudge.TeamRepository,
	cr tjudge.ContestRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) CreateTeamHandler {
	return decorator.ApplyCommandDecorators(createTeamHandler{cfg, tr, cr}, l, mc)
}
