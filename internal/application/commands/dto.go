package commands

import (
	"github.com/bmstu-itstech/tjudge-back/internal/domain/shared"
	"github.com/google/uuid"
)

func idsFromDto(ids []string) ([]shared.ID, error) {
	out := make([]shared.ID, 0)
	for _, id_str := range ids {
		id, err := uuid.Parse(id_str)
		if err != nil {
			return nil, err
		}
		out = append(out, shared.ID(id))
	}
	return out, nil
}
