package utils

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

func NewGegerateUuid() (uuid.UUID, error) {
	u7, err := uuid.NewV7()
	if err != nil {
		log.Printf("UUIDの生成に失敗しました: %v", err)
		return uuid.Nil, fmt.Errorf("UUIDの生成に失敗しました: %w", err)
	}

	return u7, nil
}
