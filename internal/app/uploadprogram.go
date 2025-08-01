package app

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
)

// Welcome to the most fun file in this whole project.
// Buckle up!

type UploadProgram struct {
	TeamId string
	GameId string
	Code   []byte
	Ext    string
}

type UploadProgramHandler decorator.CommandHandler[UploadProgram]

type uploadProgramHandler struct {
	cfg tjudge.ConfigRepository
	pr  tjudge.ProgramRepository
	prs tjudge.ProgramSourceRepository
	ter tjudge.TeamRepository
	gr  tjudge.GameRepository
	cr  tjudge.ContestRepository
	ul  tjudge.UploadListener
}

func (h uploadProgramHandler) Handle(ctx context.Context, cmd UploadProgram) error {
	config, err := h.cfg.Get(ctx)
	if err != nil {
		return err
	}
	if (len(cmd.Code) / 1024) > config.FileLimitKb() {
		return fmt.Errorf("%w (%dkb > %dkb)",
			tjudge.ErrProgramTooLarge, len(cmd.Code)/1024, config.FileLimitKb())
	}

	game, err := h.gr.Game(ctx, tjudge.GameId(cmd.GameId))
	if err != nil {
		return err
	}
	if !slices.Contains(game.AllowedExts(), cmd.Ext) {
		return fmt.Errorf("%w (%s)", tjudge.ErrWrongProgramExt, cmd.Ext)
	}

	team, err := h.ter.Team(ctx, tjudge.TeamId(cmd.TeamId))
	if err != nil {
		return err
	}
	contest, err := h.cr.Contest(ctx, team.ContestId())
	if err != nil {
		return err
	}
	if time.Now().Before(contest.Starts()) || time.Now().After(contest.Ends()) {
		return tjudge.ErrContestInactive
	}
	if !slices.Contains(contest.GameIds(), game.Id()) {
		return tjudge.ErrContestNotContainGame
	}
	if team.ContestId() != contest.Id() {
		return tjudge.ErrContestNotContainTeam
	}

	program, err := tjudge.NewProgram(team.Id(), game.Id())
	if err != nil {
		return err
	}
	src, err := tjudge.NewProgramSource(program.Id(), cmd.Code, cmd.Ext)
	if err != nil {
		return err
	}
	if err = h.prs.Upsert(ctx, src); err != nil {
		return err
	}
	if err = h.pr.Upsert(ctx, program); err != nil {
		err2 := h.prs.Delete(ctx, program.Id())
		if err2 != nil {
			// Whoopsie daisy : D
			return fmt.Errorf("loose program (couldn't save due to `%w`, followed by `%w` in src repo)", err, err2)
		}
		return err
	}

	return h.ul.OnUpload(ctx, game.Id())
}

func NewUploadProgramHandler(
	cfg tjudge.ConfigRepository,
	pr tjudge.ProgramRepository,
	prs tjudge.ProgramSourceRepository,
	ter tjudge.TeamRepository,
	gr tjudge.GameRepository,
	cr tjudge.ContestRepository,
	ul tjudge.UploadListener,
	l *slog.Logger,
	mc decorator.MetricsClient,
) UploadProgramHandler {
	return decorator.ApplyCommandDecorators(uploadProgramHandler{cfg, pr, prs, ter, gr, cr, ul}, l, mc)
}
