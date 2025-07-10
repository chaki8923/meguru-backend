package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"meguru-backend/internal/domain/entity"
	"meguru-backend/internal/domain/repository"
)

type PushSubscriptionUsecase struct {
	subscriptionRepo repository.PushSubscriptionRepository
	webPushService   WebPushService
}

type WebPushService interface {
	SendNotification(ctx context.Context, subscription *entity.PushSubscription, payload []byte) error
}

type CreateSubscriptionRequest struct {
	Endpoint string `json:"endpoint" binding:"required"`
	P256dh   string `json:"p256dh" binding:"required"`
	Auth     string `json:"auth" binding:"required"`
}

type CreateSubscriptionResponse struct {
	ID        int       `json:"id"`
	Endpoint  string    `json:"endpoint"`
	CreatedAt time.Time `json:"created_at"`
}

type SendNotificationRequest struct {
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func NewPushSubscriptionUsecase(subscriptionRepo repository.PushSubscriptionRepository, webPushService WebPushService) *PushSubscriptionUsecase {
	return &PushSubscriptionUsecase{
		subscriptionRepo: subscriptionRepo,
		webPushService:   webPushService,
	}
}

func (u *PushSubscriptionUsecase) CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*CreateSubscriptionResponse, error) {
	// Check if subscription already exists
	existing, err := u.subscriptionRepo.FindByEndpoint(ctx, req.Endpoint)
	if err != nil && err.Error() != "record not found" {
		return nil, err
	}
	if existing != nil {
		// Update existing subscription if it's different
		if existing.P256dh != req.P256dh || existing.Auth != req.Auth {
			existing.P256dh = req.P256dh
			existing.Auth = req.Auth
			if err := u.subscriptionRepo.Update(ctx, existing); err != nil {
				return nil, err
			}
		}
		return &CreateSubscriptionResponse{
			ID:        existing.ID,
			Endpoint:  existing.Endpoint,
			CreatedAt: existing.CreatedAt,
		}, nil
	}

	subscription := &entity.PushSubscription{
		Endpoint:  req.Endpoint,
		P256dh:    req.P256dh,
		Auth:      req.Auth,
		CreatedAt: time.Now(),
	}

	if err := u.subscriptionRepo.Create(ctx, subscription); err != nil {
		return nil, err
	}

	return &CreateSubscriptionResponse{
		ID:        subscription.ID,
		Endpoint:  subscription.Endpoint,
		CreatedAt: subscription.CreatedAt,
	}, nil
}

func (u *PushSubscriptionUsecase) SendNotificationToAll(ctx context.Context, req *SendNotificationRequest) error {
	subscriptions, err := u.subscriptionRepo.FindAll(ctx)
	if err != nil {
		return err
	}

	if len(subscriptions) == 0 {
		return errors.New("no subscriptions found")
	}

		payload, err := json.Marshal(map[string]string{
		"title": req.Title,
		"body":  req.Message,
	})
	if err != nil {
		return err
	}

	for _, sub := range subscriptions {
		err := u.webPushService.SendNotification(ctx, sub, payload)
		if err != nil {
			// Log error but continue to send to other subscriptions
			// If the error is due to a gone subscription, delete it from DB
			// For simplicity, we'll just log for now.
			// In a real application, you'd want to handle specific webpush errors
			// e.g., webpush.ErrSubscriptionExpired, webpush.ErrSubscriptionGone
			// and delete the subscription from the database.
			// For now, we'll just log the error.
			// fmt.Printf("Failed to send notification to %s: %v\n", sub.Endpoint, err)
			// If the subscription is no longer valid, delete it
			if err.Error() == "webpush: subscription is no longer valid" || err.Error() == "webpush: subscription is gone" {
				u.subscriptionRepo.Delete(ctx, sub.Endpoint)
			}
		}
	}

	return nil
}
