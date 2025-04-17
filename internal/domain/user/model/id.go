package model

import "github.com/google/uuid"

type ID struct {
	id uuid.UUID
}

func NewID(id uuid.UUID) ID {
	return ID{
		id: id,
	}
}

func (i ID) ID() uuid.UUID {
	return i.id
}
