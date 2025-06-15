package query_model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int64
	UserID    uuid.UUID
	Email     string
	Name      string
	CreatedAt time.Time
}
