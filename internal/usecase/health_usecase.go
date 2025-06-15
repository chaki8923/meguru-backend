package usecase

import (
	"time"
)

type HealthUsecase struct{}

func NewHealthUsecase() *HealthUsecase {
	return &HealthUsecase{}
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Timestamp int64     `json:"timestamp"`
	Time      time.Time `json:"time"`
	Uptime    string    `json:"uptime"`
}

func (h *HealthUsecase) GetHealthStatus() *HealthResponse {
	now := time.Now()
	return &HealthResponse{
		Status:    "OK!!!",
		Service:   "meguru-backend",
		Timestamp: now.Unix(),
		Time:      now,
		Uptime:    "running",
	}
} 