package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

// I hope you brought enough RAM.

type MockProgramFileRepository struct {
	sync.RWMutex
	m map[tjudge.ProgramId][]byte
}

func (r *MockProgramFileRepository) Result(ctx context.Context, id tjudge.ProgramId) ([]byte, error) {
	r.RLock()
	defer r.RUnlock()
	source, ok := r.m[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", tjudge.ErrProgramNotExist, id)
	}
	return source, nil
}

func (r *MockProgramFileRepository) Upload(ctx context.Context, id tjudge.ProgramId, source []byte) error {
	r.Lock()
	defer r.Unlock()
	r.m[id] = source
	return nil
}
