package query_model

import (
	"time"

	"github.com/google/uuid"
)

type Stores struct {
	ID          int64
	StoreID     uuid.UUID
	Name        string
	Email       string
	PhoneNumber string
	Zipcode     string
	Prefecture  string
	City        string
	Street      string
	CreatedAt   time.Time
}
