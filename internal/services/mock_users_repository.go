package services

import (
	"context"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type MockUserRepository struct {
	m     map[tjudge.UserID]*tjudge.User
	index tjudge.UserID
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		m:     make(map[tjudge.UserID]*tjudge.User),
		index: 1,
	}
}

func (r *MockUserRepository) Create(ctx context.Context, proto *tjudge.UserPrototype) (*tjudge.User, error) {
	for _, user := range r.m {
		if user.Username() == proto.Username() {
			return nil, tjudge.ErrUserAlreadyExists
		}
	}
	
	user, err := proto.Build(r.index)
	r.index++
	if err != nil {
		return user, err
	}
	r.m[user.Id()] = user
	return user, nil
}

func (r *MockUserRepository) User(ctx context.Context, id tjudge.UserID) (*tjudge.User, error) {
	got, ok := r.m[id]
	if !ok {
		return got, tjudge.ErrUserNotFound
	}
	return got, nil
}

func (r *MockUserRepository) Update(ctx context.Context, user *tjudge.User) error {
	_, ok := r.m[user.Id()]
	if !ok {
		return tjudge.ErrUserNotFound
	}
	r.m[user.Id()] = user
	return nil
}

func (r *MockUserRepository) Delete(ctx context.Context, id tjudge.UserID) error {
	_, ok := r.m[id]
	if !ok {
		return tjudge.ErrUserNotFound
	}
	delete(r.m, id)
	return nil
}

func (r *MockUserRepository) UserByUsername(ctx context.Context, username string) (*tjudge.User, error) {
	for _, user := range r.m {
		if user.Username() == username {
			return user, nil
		}
	}
	return nil, tjudge.ErrUserNotFound
}

func (r *MockUserRepository) Users(ctx context.Context) ([]*tjudge.User, error) {
	if len(r.m) == 0 {
		return nil, tjudge.ErrUserNotFound
	}
	//users := make([]*tjudge.User, 0, len(r.m))
	var users []*tjudge.User
    for _, user := range r.m {
        users = append(users, user)
    }
	return users, nil
}