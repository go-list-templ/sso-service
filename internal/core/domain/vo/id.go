package vo

import (
	"errors"
	"github.com/google/uuid"
)

var ErrInvalidID = errors.New("invalid id")

type ID struct {
	value uuid.UUID
}

func NewID() ID {
	return ID{
		value: uuid.New(),
	}
}

func FromStr(id string) (ID, error) {
	parse, err := uuid.Parse(id)
	if err != nil {
		return ID{}, ErrInvalidID
	}

	return ID{parse}, err
}

func UnsafeID(id uuid.UUID) ID {
	return ID{value: id}
}

func (i *ID) Value() uuid.UUID {
	return i.value
}
