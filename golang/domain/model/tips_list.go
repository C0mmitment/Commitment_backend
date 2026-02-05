package model

import (
	"github.com/google/uuid"
)

type TipsList struct {
	TipsId   uuid.UUID
	Title    string
	Category string
	Content  string
}
