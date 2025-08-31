package tjudge

import "fmt"

type TeamFactory interface {
	Create(name string, code string, contest ContestID) (*Team, error)
}

type FixedSizeFactory struct {
	maxSize int
}

func NewFixedSizeFactory(maxSize int) (*FixedSizeFactory, error) {
    if maxSize <= 0 {
        return nil, fmt.Errorf("%w: maxSize must be positive", ErrInvalidInput)
    }
    return &FixedSizeFactory{maxSize: maxSize}, nil
}

func (f *FixedSizeFactory) Create(name string, code string, contest ContestID) (*Team, error) {
    return newTeam(name, TeamCode(code), nil, contest, f.maxSize)
}