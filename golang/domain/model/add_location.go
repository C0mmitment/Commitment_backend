package model

import (
	"github.com/google/uuid"
)

type AddLocation struct {
	LocationId uuid.UUID
	UserId     uuid.UUID
	Lat        float64
	Lng        float64
	Geo        string
}
