package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"
)

type PostgresTeamRepository struct {
	db       *sqlx.DB
	userRepo tjudge.UserRepository
}

func NewPostgresTeamRepository(db *sqlx.DB, userRepo tjudge.UserRepository) *PostgresTeamRepository {
	return &PostgresTeamRepository{db, userRepo}
}
func (r *PostgresTeamRepository) Save(ctx context.Context, team *tjudge.Team) (*tjudge.Team, error) {
	if team == nil {
		return nil, fmt.Errorf("team is nil")
	}

	var leaderID = 0
	var contestID = 0

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Используем leaderId только после проверки на nil
	if team.Leader() != nil {
		leaderID = int(team.Leader().Id())
	}

	// Более безопасное получение contest ID
	contestValue := team.Contest()
	if contestValue == 0 {
		return nil, fmt.Errorf("contest is required")
	}
	contestID = int(contestValue)

	_, err = pgutils.Exec(ctx, tx, `
		INSERT INTO teams (code, name, leader_id, contest_id, max_size) 
		VALUES ($1, $2, $3, $4, $5)`,
		string(team.Code()), team.Name(), leaderID, contestID, team.MaxSize())

	if pgutils.IsUniqueViolationError(err) {
		return nil, tjudge.ErrTeamAlreadyExist
	} else if err != nil {
		return nil, fmt.Errorf("failed to insert team: %w", err)
	}

	for _, member := range team.Members() {
		// Проверяем, что member не nil
		if member == nil {
			fmt.Println("Warning: nil member found in team members")
			continue
		}
		
		// Проверяем, что лидер не nil перед сравнением
		if team.Leader() != nil && member.Id() == team.Leader().Id() {
			continue
		}
		
		_, err := tx.ExecContext(ctx, `
			INSERT INTO team_members (user_id, team_code) 
			VALUES ($1, $2) 
			ON CONFLICT DO NOTHING`,
			int(member.Id()), string(team.Code()))
		if err != nil {
			return nil, fmt.Errorf("failed to insert team member: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return team, nil
}
func (r *PostgresTeamRepository) Team(ctx context.Context, code string) (*tjudge.Team, error) {
	var row teamRow

	err := pgutils.Get(ctx, r.db, &row, `
		SELECT code, name, leader_id, contest_id, max_size 
		FROM teams 
		WHERE code = $1`, code)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %s", tjudge.ErrTeamNotFound, code)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	// Обрабатываем случай, когда leader_id = 0
	var leader *tjudge.User
	if row.LeaderID > 0 {
		leader, err = r.userRepo.User(ctx, tjudge.UserID(row.LeaderID))
		if err != nil {
			return nil, fmt.Errorf("failed to get leader: %w", err)
		}
	}

	var memberIDs []int
	err = pgutils.Select(ctx, r.db, &memberIDs, `
		SELECT user_id 
		FROM team_members 
		WHERE team_code = $1`, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	members := make([]*tjudge.User, 0, len(memberIDs))
	for _, memberID := range memberIDs {
		member, err := r.userRepo.User(ctx, tjudge.UserID(memberID))
		if err != nil {
			return nil, fmt.Errorf("failed to get member %d: %w", memberID, err)
		}
		// Проверяем, что member не nil
		if member != nil {
			members = append(members, member)
		} else {
			fmt.Printf("Warning: user with ID %d not found\n", memberID)
		}
	}

	team := tjudge.RestoreTeam(
		tjudge.TeamCode(row.Code),
		row.Name,
		leader,
		tjudge.ContestID(row.ContestID),
		row.MaxSize,
	)

	for _, member := range members {
		if leader == nil || member.Id() != leader.Id() {
			if err := team.AddMember(member); err != nil {
				return nil, fmt.Errorf("failed to add member to team: %w", err)
			}
		}
	}

	return team, nil
}

func (r *PostgresTeamRepository) Update(ctx context.Context, team *tjudge.Team) error {
	if team == nil {
		return fmt.Errorf("team is nil")
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Используем 0 вместо NULL для пустого лидера
	leaderID := 0
	if team.Leader() != nil {
		leaderID = int(team.Leader().Id())
	}
	print(leaderID)

	err = pgutils.RequireAffected(pgutils.Exec(ctx, tx, `
		UPDATE teams 
		SET name = $1, leader_id = $2, contest_id = $3, max_size = $4 
		WHERE code = $5`,
		team.Name(),
		leaderID,
		int(team.Contest()),
		team.MaxSize(),
		string(team.Code())))

	if errors.Is(err, pgutils.ErrNoAffectedRows) {
		return fmt.Errorf("%w: %s", tjudge.ErrTeamNotFound, team.Code())
	} else if err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}

	// Удаляем всех членов команды, кроме лидера (если он есть)
	if leaderID > 0 {
		_, err = tx.ExecContext(ctx, `
			DELETE FROM team_members 
			WHERE team_code = $1 AND user_id != $2`,
			string(team.Code()), leaderID)
	} else {
		_, err = tx.ExecContext(ctx, `
			DELETE FROM team_members 
			WHERE team_code = $1`,
			string(team.Code()))
	}
	if err != nil {
		return fmt.Errorf("failed to clear team members: %w", err)
	}

	for _, member := range team.Members() {
		// Пропускаем лидера, если он есть
		if team.Leader() != nil && member.Id() == team.Leader().Id() {
			continue
		}

		_, err := tx.ExecContext(ctx, `
			INSERT INTO team_members (user_id, team_code) 
			VALUES ($1, $2) 
			ON CONFLICT DO NOTHING`,
			int(member.Id()), string(team.Code()))
		if err != nil {
			return fmt.Errorf("failed to add team member: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}


func (r *PostgresTeamRepository) Delete(ctx context.Context, code string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		DELETE FROM team_members 
		WHERE team_code = $1`, code)
	if err != nil {
		return fmt.Errorf("failed to delete team members: %w", err)
	}

	err = pgutils.RequireAffected(pgutils.Exec(ctx, tx, `
		DELETE FROM teams 
		WHERE code = $1`, code))

	if errors.Is(err, pgutils.ErrNoAffectedRows) {
		return fmt.Errorf("%w: %s", tjudge.ErrTeamNotFound, code)
	} else if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *PostgresTeamRepository) TeamsByContest(ctx context.Context, contest tjudge.ContestID) ([]*tjudge.Team, error) {
	var teamRows []teamRow

	err := pgutils.Select(ctx, r.db, &teamRows, `
		SELECT code, name, leader_id, contest_id, max_size 
		FROM teams 
		WHERE contest_id = $1`, int(contest))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*tjudge.Team{}, nil
		}
		return nil, fmt.Errorf("failed to get teams: %w", err)
	}

	teams := make([]*tjudge.Team, 0, len(teamRows))
	for _, row := range teamRows {
		// Обрабатываем случай, когда лидер может быть NULL
		var leader *tjudge.User
		if row.LeaderID > 0 {
			leader, err = r.userRepo.User(ctx, tjudge.UserID(row.LeaderID))
			if err != nil {
				return nil, fmt.Errorf("failed to get leader: %w", err)
			}
		}

		var memberIDs []int64
		err = pgutils.Select(ctx, r.db, &memberIDs, `
			SELECT user_id 
			FROM team_members 
			WHERE team_code = $1`, row.Code)
		if err != nil {
			return nil, fmt.Errorf("failed to get team members: %w", err)
		}

		members := make([]*tjudge.User, 0, len(memberIDs))
		for _, memberID := range memberIDs {
			member, err := r.userRepo.User(ctx, tjudge.UserID(memberID))
			if err != nil {
				return nil, fmt.Errorf("failed to get member %d: %w", memberID, err)
			}
			members = append(members, member)
		}

		team := tjudge.RestoreTeam(
			tjudge.TeamCode(row.Code),
			row.Name,
			leader,
			tjudge.ContestID(row.ContestID),
			row.MaxSize,
		)

		for _, member := range members {
			if leader == nil || member.Id() != leader.Id() {
				if err := team.AddMember(member); err != nil {
					return nil, fmt.Errorf("failed to add member to team: %w", err)
				}
			}
		}

		teams = append(teams, team)
	}

	return teams, nil
}

func (r *PostgresTeamRepository) Teams(ctx context.Context) ([]*tjudge.Team, error) {
	var teamRows []teamRow

	err := pgutils.Select(ctx, r.db, &teamRows, `
		SELECT code, name, leader_id, contest_id, max_size 
		FROM teams`)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*tjudge.Team{}, nil
		}
		return nil, fmt.Errorf("failed to get teams: %w", err)
	}

	teams := make([]*tjudge.Team, 0, len(teamRows))
	for _, row := range teamRows {
		// Обрабатываем случай, когда лидер может быть NULL
		var leader *tjudge.User
		if row.LeaderID > 0 {
			leader, err = r.userRepo.User(ctx, tjudge.UserID(row.LeaderID))
			if err != nil {
				return nil, fmt.Errorf("failed to get leader: %w", err)
			}
		}

		var memberIDs []int64
		err = pgutils.Select(ctx, r.db, &memberIDs, `
			SELECT user_id 
			FROM team_members 
			WHERE team_code = $1`, row.Code)
		if err != nil {
			return nil, fmt.Errorf("failed to get team members: %w", err)
		}

		members := make([]*tjudge.User, 0, len(memberIDs))
		for _, memberID := range memberIDs {
			member, err := r.userRepo.User(ctx, tjudge.UserID(memberID))
			if err != nil {
				return nil, fmt.Errorf("failed to get member %d: %w", memberID, err)
			}
			members = append(members, member)
		}

		team := tjudge.RestoreTeam(
			tjudge.TeamCode(row.Code),
			row.Name,
			leader,
			tjudge.ContestID(row.ContestID),
			row.MaxSize,
		)

		for _, member := range members {
			if leader == nil || member.Id() != leader.Id() {
				if err := team.AddMember(member); err != nil {
					return nil, fmt.Errorf("failed to add member to team: %w", err)
				}
			}
		}

		teams = append(teams, team)
	}

	return teams, nil
}

type teamRow struct {
	Code      string         `db:"code"`
	Name      string         `db:"name"`
	LeaderID  int  			`db:"leader_id"`
	ContestID int64          `db:"contest_id"`
	MaxSize   int            `db:"max_size"`
}