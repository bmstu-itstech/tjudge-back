package services_test

import (
	"context"
	"testing"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/internal/services"
	"github.com/bmstu-itstech/tjudge-back/internal/utils/postgres"
	"github.com/stretchr/testify/require"
)

func setupPostgresRepository(t *testing.T) *services.PostgresUserRepository {
	db, closeFn := postgres.ConnectToDatabase()
	t.Cleanup(closeFn)
	return services.NewPostgresUserRepository(db)
}

func TestPostgresRep(t *testing.T) {
	repo := setupPostgresRepository(t)
	testRepository(t, repo)
	testRepository_User(t, repo)
	testRepository_Update(t, repo)
	testRepository_Delete(t, repo)
}

func TestMockRep(t *testing.T) {
	repo := services.NewMockUserRepository()
	testRepository(t, repo)
	testRepository_User(t, repo)
	testRepository_Update(t, repo)
	testRepository_Delete(t, repo)
}

func testRepository(t *testing.T, repo tjudge.UserRepository) {
	ctx := context.Background()

	proto := tjudge.MustNewUserPrototype("Ivan Petrov", "Ivan", "12345678", true)
	user, err := repo.Create(ctx, proto)
	require.NoError(t, err)
	require.Equal(t, proto.Username(), user.Username())
	require.Equal(t, proto.Fullname(), user.Fullname())

	_, err = repo.Create(ctx, proto)
	require.ErrorIs(t, err, tjudge.ErrUserAlreadyExists)

	_ = user.SetFullname("Vasya")

	repo.Update(ctx, user)

	user, err = repo.User(ctx, user.Id())
	require.NoError(t, err)
	require.Equal(t, "Vasya", user.Fullname())

	err = repo.Delete(ctx, user.Id())
	require.NoError(t, err)

	_, err = repo.User(ctx, user.Id())
	require.ErrorIs(t, err, tjudge.ErrUserNotFound)
}

func testRepository_User(t *testing.T, repo tjudge.UserRepository) {
	ctx := context.Background()
	_, err := repo.User(ctx, 100)
	require.ErrorIs(t, err, tjudge.ErrUserNotFound)

}

func testRepository_Update(t *testing.T, repo tjudge.UserRepository) {
	ctx := context.Background()
	proto := tjudge.MustNewUserPrototype("Andrey22", "Andrey Vas", "12345678", true)
	user, err := repo.Create(ctx, proto)
	require.NoError(t, err)
	err = repo.Delete(ctx, user.Id())
	require.NoError(t, err)

	err = repo.Update(ctx, user)
	require.ErrorIs(t, err, tjudge.ErrUserNotFound)
}

func testRepository_Delete(t *testing.T, repo tjudge.UserRepository) {
	ctx := context.Background()
	err := repo.Delete(ctx, 100)
	require.ErrorIs(t, err, tjudge.ErrUserNotFound)
}
