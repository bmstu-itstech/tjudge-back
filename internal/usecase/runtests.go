package usecases

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/pkg/decorator"
	"github.com/bmstu-itstech/tjudge-back/pkg/util"
)

// Welcome to the *second* most fun file in this whole project.
// Buckle up again!

type RunRounds struct {
	GameId      string
	WithTimeout bool
}

type RunRoundsHandler decorator.CommandHandler[RunRounds]

type runRoundsHandler struct {
	cfg tjudge.ConfigRepository
	la  tjudge.Launcher
	es  tjudge.RoundEventStore
	pr  tjudge.ProgramRepository
	gr  tjudge.GameRepository
}

func (h runRoundsHandler) toRun(
	ctx context.Context,
	game tjudge.Game,
	withTimeout bool,
	timeOut time.Duration) ([]map[tjudge.TeamId]tjudge.ProgramId, error) {
	// we need to find rounds which should be run based on active programs,
	// but that haven't been scheduled or finished

	active, err := h.pr.Actives(ctx, game.Id())
	if err != nil {
		return nil, err
	}

	done := make(map[string]struct{})
	events, err := h.es.Load(ctx, game.Id())
	if err != nil {
		return nil, err
	}
	for _, v := range events {
		if withTimeout && v.Type == "scheduled" && time.Since(v.Time) > timeOut {
			continue
		}
		stale := false
		for team_id, prog_id := range v.Programs {
			if active[team_id] != prog_id {
				stale = true
				break
			}
		}
		if !stale {
			teamIDs := make([]tjudge.TeamId, 0, len(v.Programs))
			for tid := range v.Programs {
				teamIDs = append(teamIDs, tid)
			}
			slices.Sort(teamIDs)
			key := ""
			for _, tid := range teamIDs {
				key += string(v.Programs[tid]) + "|"
			}
			done[key] = struct{}{}
		}
	}

	teamIDs := make([]tjudge.TeamId, 0, len(active))
	for tid := range active {
		teamIDs = append(teamIDs, tid)
	}
	slices.Sort(teamIDs)

	toRun := make([]map[tjudge.TeamId]tjudge.ProgramId, 0)
	iterator := util.Iter{}
	iterator.N = len(active)
	iterator.K = int(game.Players())
	for iterator.Next() {
		comb := make(map[tjudge.TeamId]tjudge.ProgramId)
		key := ""
		for _, idx := range iterator.Combination {
			tid := teamIDs[idx]
			comb[tid] = active[tid]
			key += string(active[tid]) + "|"
		}
		if _, found := done[key]; !found {
			toRun = append(toRun, comb)
		}
	}
	return toRun, nil
}

func (h runRoundsHandler) Handle(ctx context.Context, cmd RunRounds) error {
	config, err := h.cfg.Get(ctx)
	if err != nil {
		return err
	}
	game, err := h.gr.Game(ctx, tjudge.GameId(cmd.GameId))
	if err != nil {
		return err
	}

	// we assume that the programs have already been uploaded and stored
	toRun, err := h.toRun(ctx, game, cmd.WithTimeout, config.JudgeTimeout())
	if err != nil {
		return err
	}
	for _, run := range toRun {
		event := tjudge.RoundEvent{
			Type:     "scheduled",
			GameId:   game.Id(),
			Programs: run,
			Results:  nil,
			Time:     time.Now(),
		}
		err := h.es.Append(ctx, event)
		if err != nil {
			return err
		}
		crew := make([]tjudge.ProgramId, 0, len(run))
		for _, p := range run {
			crew = append(crew, p)
		}
		go func(mission []tjudge.ProgramId) {
			// FIXME: And what if we get an error here?
			// I mean, it'll be caught by timeouts eventually, I guess.
			// But would still want to catch it.
			roundEvent := h.la.Run(tjudge.GameId(cmd.GameId), mission)
			_ = h.es.Append(ctx, roundEvent)
		}(crew)
	}
	return nil
}

func NewRunRoundsHandler(
	cfg tjudge.ConfigRepository,
	la tjudge.Launcher,
	es tjudge.RoundEventStore,
	pr tjudge.ProgramRepository,
	gr tjudge.GameRepository,
	l *slog.Logger,
	mc decorator.MetricsClient,
) RunRoundsHandler {
	return decorator.ApplyCommandDecorators(runRoundsHandler{cfg, la, es, pr, gr}, l, mc)
}
