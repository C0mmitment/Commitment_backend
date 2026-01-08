package dto

import "github.com/google/uuid"

type DeleteHeatmapRequest struct {
	UserId uuid.UUID `param:"user_uuid"`
}
