package dto

import "github.com/google/uuid"

type DeleteHeatmapRequest struct {
	UserId uuid.UUID `json:"user_uuid"`
}
