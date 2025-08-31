package tjudge_test

import (
	"errors"
	"testing"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

func createTestUser(id tjudge.UserID, username string, fullname string) *tjudge.User {
	proto := tjudge.MustNewUserPrototype(username, fullname, "password123", true)
	user, _ := proto.Build(id)
	return user
}

func TestTeamCreation(t *testing.T) {
	leader := createTestUser(1, "leader1", "Team Leader")
	contestID := tjudge.ContestID(1)

	tests := []struct {
		name      string
		code      tjudge.TeamCode
		teamName  string
		leader    *tjudge.User
		contest   tjudge.ContestID
		maxSize   int
		wantError error
	}{
		{
			name:      "valid team",
			code:      "TEAM1",
			teamName:  "Team One",
			leader:    leader,
			contest:   contestID,
			maxSize:   5,
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team := tjudge.RestoreTeam(tt.code, tt.teamName, tt.leader, tt.contest, tt.maxSize)

			if tt.wantError == nil {
				if team.Code() != tt.code {
					t.Errorf("Code() = %v, want %v", team.Code(), tt.code)
				}
				if team.Name() != tt.teamName {
					t.Errorf("Name() = %v, want %v", team.Name(), tt.teamName)
				}
				if team.Leader() != tt.leader {
					t.Errorf("Leader() = %v, want %v", team.Leader(), tt.leader)
				}
				if team.Contest() != tt.contest {
					t.Errorf("Contest() = %v, want %v", team.Contest(), tt.contest)
				}
				if team.MaxSize() != tt.maxSize {
					t.Errorf("MaxSize() = %v, want %v", team.MaxSize(), tt.maxSize)
				}
				if len(team.Members()) != 1 || team.Members()[0] != tt.leader {
					t.Errorf("Members not initialized correctly")
				}
			}
		})
	}
}

func TestTeamAddRemoveMember(t *testing.T) {
	leader := createTestUser(1, "leader1", "Team Leader")
	member := createTestUser(2, "member1", "Team Member")
	team := tjudge.RestoreTeam("TEAM1", "Test Team", leader, 1, 5)

	t.Run("add member", func(t *testing.T) {
		err := team.AddMember(member)
		if err != nil {
			t.Errorf("AddMember() unexpected error: %v", err)
		}
		if len(team.Members()) != 2 || team.Members()[1] != member {
			t.Errorf("Member not added correctly")
		}
	})

	t.Run("remove member", func(t *testing.T) {
		err := team.RemoveMember(member)
		if err != nil {
			t.Errorf("RemoveMember() unexpected error: %v", err)
		}
		if len(team.Members()) != 1 {
			t.Errorf("Member not removed correctly")
		}
	})
}

func TestTeamMaxSize(t *testing.T) {
	leader := createTestUser(1, "leader1", "Team Leader")
	team := tjudge.RestoreTeam("TEAM1", "Test Team", leader, 1, 2)

	member1 := createTestUser(2, "member1", "Member One")
	member2 := createTestUser(3, "member2", "Member Two")

	t.Run("add up to max size", func(t *testing.T) {
		err := team.AddMember(member1)
		if err != nil {
			t.Errorf("AddMember() unexpected error: %v", err)
		}
	})

	t.Run("exceed max size", func(t *testing.T) {
		err := team.AddMember(member2)
		if !errors.Is(err, tjudge.ErrTeamFull) {
			t.Errorf("AddMember() error = %v, want %v", err, tjudge.ErrTeamFull)
		}
	})
}

func TestTeamGetters(t *testing.T) {
	leader := createTestUser(1, "leader1", "Team Leader")
	team := tjudge.RestoreTeam("TEAM1", "Test Team", leader, 1, 3)

	t.Run("get code", func(t *testing.T) {
		if team.Code() != "TEAM1" {
			t.Errorf("Code() = %v, want %v", team.Code(), "TEAM1")
		}
	})

	t.Run("get name", func(t *testing.T) {
		if team.Name() != "Test Team" {
			t.Errorf("Name() = %v, want %v", team.Name(), "Test Team")
		}
	})

	t.Run("get leader", func(t *testing.T) {
		if team.Leader() != leader {
			t.Errorf("Leader() = %v, want %v", team.Leader(), leader)
		}
	})

	t.Run("get contest", func(t *testing.T) {
		if team.Contest() != 1 {
			t.Errorf("Contest() = %v, want %v", team.Contest(), 1)
		}
	})

	t.Run("get max size", func(t *testing.T) {
		if team.MaxSize() != 3 {
			t.Errorf("MaxSize() = %v, want %v", team.MaxSize(), 3)
		}
	})

	t.Run("get members", func(t *testing.T) {
		members := team.Members()
		if len(members) != 1 || members[0] != leader {
			t.Errorf("Members() = %v, want %v", members, []*tjudge.User{leader})
		}
	})
}
