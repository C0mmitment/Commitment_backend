package dto

import (
	"github.com/google/uuid"
)

type TipsList struct {
	TipsId   uuid.UUID `json:"tips_id"`
	Title    string    `json:"title"`
	Category string    `json:"category"`
	Content  string    `json:"content"`
}

type TipesListResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Tips    []TipsList `json:"tips_list"`
	Error   string     `json:"error,omitempty"`
}
