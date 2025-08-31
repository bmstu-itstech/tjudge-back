package services_test

import (
	"context"
	"testing"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
	"github.com/bmstu-itstech/tjudge-back/internal/services"
	"github.com/bmstu-itstech/tjudge-back/internal/utils/postgres"
	"github.com/stretchr/testify/require"
)

func setupPostgresTeamRepository(t *testing.T) (*services.PostgresTeamRepository, *services.PostgresUserRepository, *tjudge.User, *tjudge.FixedSizeFactory) {
	db, closeFn := postgres.ConnectToDatabase()
	t.Cleanup(closeFn)
	
	userRepo := services.NewPostgresUserRepository(db)
	teamRepo := services.NewPostgresTeamRepository(db, userRepo)
	
	// Создаем тестового пользователя для лидера команды
	ctx := context.Background()
	userProto := tjudge.MustNewUserPrototype("Team Leader5", "leader", "password555", true)
	leader, err := userRepo.Create(ctx, userProto)
	require.NoError(t, err)

	// Создаем фабрику команд
	factory, err := tjudge.NewFixedSizeFactory(5)
	require.NoError(t, err)
	
	return teamRepo, userRepo, leader, factory
}

func TestPostgresTeamRep(t *testing.T) {
	teamRepo, _, leader, factory := setupPostgresTeamRepository(t)
	
	testTeamRepository(t, teamRepo, leader, factory)
	testTeamRepository_Team(t, teamRepo, leader, factory)
	testTeamRepository_Update(t, teamRepo, leader, factory)
	testTeamRepository_Delete(t, teamRepo, leader, factory)
}

func testTeamRepository(t *testing.T, repo tjudge.TeamRepository, leader *tjudge.User, factory *tjudge.FixedSizeFactory) {
	ctx := context.Background()

	// Создаем команду через фабрику
	team, err := factory.Create("Dream Team11", "DT665", 5)
	require.NoError(t, err)

	// Сохраняем команду в репозитории
	createdTeam, err := repo.Save(ctx, team)
	require.NoError(t, err)
	require.Equal(t, "Dream Team11", createdTeam.Name())
	require.Equal(t, tjudge.TeamCode("DT665"), createdTeam.Code())
	require.Equal(t, leader.Id(), createdTeam.Leader().Id())
	require.Equal(t, tjudge.ContestID(5), createdTeam.Contest())
	require.Equal(t, 5, createdTeam.MaxSize())

	// Попытка создать команду с тем же кодом
	_, err = repo.Save(ctx, team)
	require.ErrorIs(t, err, tjudge.ErrTeamAlreadyExist)

	// Обновляем данные команды
	createdTeam.SetName("Super Team5")
	err = repo.Update(ctx, createdTeam)
	require.NoError(t, err)

	// Проверяем обновленные данные
	updatedTeam, err := repo.Team(ctx, string(createdTeam.Code()))
	require.NoError(t, err)
	require.Equal(t, "Super Team5", updatedTeam.Name())

	// Удаляем команду
	err = repo.Delete(ctx, string(createdTeam.Code()))
	require.NoError(t, err)

	// Проверяем что команда удалена
	_, err = repo.Team(ctx, string(createdTeam.Code()))
	require.ErrorIs(t, err, tjudge.ErrTeamNotFound)
}

func testTeamRepository_Team(t *testing.T, repo tjudge.TeamRepository, leader *tjudge.User, factory *tjudge.FixedSizeFactory) {
	ctx := context.Background()
	
	// Проверяем несуществующую команду
	_, err := repo.Team(ctx, "nonexistent-code")
	require.ErrorIs(t, err, tjudge.ErrTeamNotFound)

	// Создаем тестовую команду
	team, err := factory.Create("Test Team5", "TT555", 5)
	require.NoError(t, err)
	
	createdTeam, err := repo.Save(ctx, team)
	require.NoError(t, err)

	// Проверяем что можем получить команду по коду
	foundTeam, err := repo.Team(ctx,string(createdTeam.Code()))
	require.NoError(t, err)
	require.Equal(t, createdTeam.Code(), foundTeam.Code())

	// Удаляем тестовую команду
	err = repo.Delete(ctx, string(createdTeam.Code()))
	require.NoError(t, err)
}

func testTeamRepository_Update(t *testing.T, repo tjudge.TeamRepository, leader *tjudge.User, factory *tjudge.FixedSizeFactory) {
	ctx := context.Background()
	
	// Создаем тестовую команду
	team, err := factory.Create("Update Test5", "UT789", 5)
	require.NoError(t, err)
	
	createdTeam, err := repo.Save(ctx, team)
	require.NoError(t, err)

	// Удаляем команду
	err = repo.Delete(ctx, string(createdTeam.Code()))
	require.NoError(t, err)

	// Пытаемся обновить удаленную команду
	err = repo.Update(ctx, createdTeam)
	require.ErrorIs(t, err, tjudge.ErrTeamNotFound)
}

func testTeamRepository_Delete(t *testing.T, repo tjudge.TeamRepository, leader *tjudge.User, factory *tjudge.FixedSizeFactory) {
	ctx := context.Background()
	
	// Пытаемся удалить несуществующую команду
	err := repo.Delete(ctx, "nonexistent-code")
	require.ErrorIs(t, err, tjudge.ErrTeamNotFound)

	// Создаем и сразу удаляем тестовую команду
	team, err := factory.Create("Temp Team5", "TEMP555", 1)
	require.NoError(t, err)
	
	createdTeam, err := repo.Save(ctx, team)
	require.NoError(t, err)

	err = repo.Delete(ctx, string(createdTeam.Code()))
	require.NoError(t, err)
}
