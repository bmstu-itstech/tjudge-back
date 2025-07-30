package app

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
	"github.com/bmstu-itstech/tjudge-back/pkg/util"
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
	tor tjudge.TourRepository
	rr  tjudge.RoundRepository
	la  tjudge.Launcher
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

	tour, err := h.tor.Active(ctx, game.Id())
	if err != nil {
		return err
	}
	new_rounds := make([]tjudge.RoundId, 0)
	for _, round_id := range tour.RoundIds() {
		round, err := h.rr.Round(ctx, round_id)
		if err != nil {
			return err
		}
		rerun := false
		for _, result := range round.Results() {
			r_program, err := h.pr.Program(ctx, result.ProgramId())
			if err != nil {
				return err
			}
			if r_program.TeamId() == tjudge.TeamId(cmd.TeamId) {
				rerun = true
				break
			}
		}
		if !rerun {
			new_rounds = append(new_rounds, round.Id())
		}
	}
	contest_teams, err := h.ter.ByContest(ctx, contest.Id())
	if err != nil {
		return err
	}
	programs := make([]tjudge.Program, 0)
	for _, t := range contest_teams {
		pr, err := h.pr.Active(ctx, game.Id(), t.Id())
		if err != tjudge.ErrProgramNotExist {
			return err
		}
		if err == nil {
			programs = append(programs, pr)
		}
	}
	iterator := util.Iter{}
	iterator.N = len(programs)
	iterator.K = int(game.Players() - 1)
	for iterator.Next() {
		to_run := make([]tjudge.Program, 0, game.Players())
		to_run = append(to_run, program)
		for _, v := range iterator.Combination {
			to_run = append(to_run, programs[v])
		}
		results, err := h.la.Run(game, to_run)
		if err != nil {
			return err
		}
		round, err := tjudge.NewRound(results)
		if err != nil {
			return err
		}
		if err = h.rr.Upsert(ctx, round); err != nil {
			return err
		}
		new_rounds = append(new_rounds, round.Id())
	}
	new_tour, err := tjudge.NewTour(game.Id(), new_rounds)
	if err != nil {
		return err
	}
	return h.tor.Upsert(ctx, new_tour)
}

func NewUploadProgramHandler(
	cfg tjudge.ConfigRepository,
	pr tjudge.ProgramRepository,
	prs tjudge.ProgramSourceRepository,
	ter tjudge.TeamRepository,
	gr tjudge.GameRepository,
	cr tjudge.ContestRepository,
	tor tjudge.TourRepository,
	rr tjudge.RoundRepository,
	la tjudge.Launcher,
	l *slog.Logger,
	mc decorator.MetricsClient,
) UploadProgramHandler {
	return decorator.ApplyCommandDecorators(uploadProgramHandler{cfg, pr, prs, ter, gr, cr, tor, rr, la}, l, mc)
}
