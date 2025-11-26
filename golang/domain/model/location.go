package model

import (
	"github.com/google/uuid"
)

type Location struct {
	LocationId uuid.UUID
	UserId     uuid.UUID
	Lat        float64
	Lng        float64
	Geo        string
}
