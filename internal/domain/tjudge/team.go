package tjudge

import (
	"fmt"
)

type ContestID int
type TeamCode string

type Team struct {
	code     	TeamCode
	name     	string
	leader 		*User
	contest 	ContestID
	maxSize 	int
	members 	[]*User
}

func (t *Team) Name() string {
	return t.name
}

func (t *Team) Code() TeamCode {
	return t.code
}

func (t *Team) Leader() *User {
	return t.leader
}

func (t *Team) Contest() ContestID {
	return t.contest
}

func (t *Team) MaxSize() int {
	return t.maxSize
}


func (t *Team) Members() []*User {
	return t.members
}


func newTeam(name string, code TeamCode, leader *User, contest ContestID, maxSize int) (*Team, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: expected not empty name", ErrInvalidInput)
	}
	members := []*User {leader}
	return &Team{
		code,
		name,
		leader,
		contest,
		maxSize,
		members,
	}, nil
}

func (t *Team) AddMember(user *User) error {
	if len(t.members) == t.MaxSize() {
		return fmt.Errorf("%w: team is already completed", ErrTeamFull)
	}
	for i := 0; i < len(t.members); i++ {  
		if t.members[i].Id() == user.Id() {
			return fmt.Errorf("%w: user already in team", ErrUserAlreadyInTeam)
		}
	}
	t.members = append(t.members, user)
	return nil
}

func (t *Team) RemoveMember(user *User) error {
    fl := 0
	for i := 0; i < len(t.members); i++ {  
        if t.members[i] != user {  
            t.members[i] = t.members[len(t.members)-1]
        } else {
			fl = 1
		}
    }
	if fl == 1 {
		t.members = t.members[:len(t.members)-1] 
		return nil
	}
	return ErrUserNotFound
}

func RestoreTeam(code TeamCode, name string, leader *User, contest ContestID, maxSize int) *Team {
	return &Team{
		code,
		name,
		leader,
		contest,
		maxSize,
		[]*User {leader},
	}
}

func (t *Team) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: expected not empty name", ErrInvalidInput)
	}
	t.name = name
	return nil
}

func (t *Team) SetLeader(leader *User) error {
	if leader == nil {
		return fmt.Errorf("%w: expected real user", ErrInvalidInput)
	}
	t.leader = leader
	return nil
}
