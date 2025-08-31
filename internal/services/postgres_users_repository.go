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

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, proto *tjudge.UserPrototype) (*tjudge.User, error) {
	if proto == nil {
		return nil, fmt.Errorf("user prototype is nil")
	}

	var id int64
	err := pgutils.Get(ctx, r.db, &id, `
		INSERT INTO users (username, fullname, password, isadmin) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`,
		proto.Username(), proto.Fullname(), proto.Passhash(), proto.IsAdmin())

	if pgutils.IsUniqueViolationError(err) {
		return nil, tjudge.ErrUserAlreadyExists
	} else if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return tjudge.RestoreUser(id, proto.Username(), proto.Fullname(), proto.Passhash()), nil
}

func (r *PostgresUserRepository) User(ctx context.Context, id tjudge.UserID) (*tjudge.User, error) {
	var row userRow
	err := pgutils.Get(ctx, r.db, &row, `
		SELECT username, fullname, password, isadmin 
		FROM users 
		WHERE id = $1`, int(id))

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %d", tjudge.ErrUserNotFound, id)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return tjudge.RestoreUser(int64(id), row.Username, row.Fullname, row.Password), nil
}

func (r *PostgresUserRepository) UserByUsername(ctx context.Context, username string) (*tjudge.User, error) {
	var row userRow
	err := pgutils.Get(ctx, r.db, &row, `
		SELECT id, username, fullname, password, isadmin 
		FROM users 
		WHERE username = $1`, username)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %s", tjudge.ErrUserNotFound, username)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return tjudge.RestoreUser(int64(row.Id), row.Username, row.Fullname, row.Password), nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *tjudge.User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	err := pgutils.RequireAffected(pgutils.Exec(ctx, r.db, `
		UPDATE users 
		SET username = $1, fullname = $2, password = $3, isadmin = $4 
		WHERE id = $5`,
		user.Username(), user.Fullname(), user.Passhash(), user.IsAdmin(), int(user.Id())))

	if errors.Is(err, pgutils.ErrNoAffectedRows) {
		return fmt.Errorf("%w: %d", tjudge.ErrUserNotFound, user.Id())
	} else if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id tjudge.UserID) error {
	err := pgutils.RequireAffected(pgutils.Exec(ctx, r.db, `
		DELETE FROM users 
		WHERE id = $1`, int(id)))

	if errors.Is(err, pgutils.ErrNoAffectedRows) {
		return fmt.Errorf("%w: %d", tjudge.ErrUserNotFound, id)
	} else if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) Users(ctx context.Context) ([]*tjudge.User, error) {
	var rows []userRow
	err := pgutils.Select(ctx, r.db, &rows, `
		SELECT id, username, fullname, password, isadmin 
		FROM users`)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*tjudge.User{}, nil
		}
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	users := make([]*tjudge.User, len(rows))
	for i, row := range rows {
		users[i] = tjudge.RestoreUser(int64(row.Id), row.Username, row.Fullname, row.Password)
	}

	return users, nil
}

type userRow struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Fullname string `db:"fullname"`
	Password []byte `db:"password"`
	IsAdmin  bool   `db:"isadmin"`
}