package webpush

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"meguru-backend/internal/domain/entity"
	"github.com/SherClockHolmes/webpush-go"
)

type webPushService struct {
	publicKey  string
	privateKey string
}

func NewWebPushService() *webPushService {
	return &webPushService{
		publicKey:  os.Getenv("VAPID_PUBLIC_KEY"),
		privateKey: os.Getenv("VAPID_PRIVATE_KEY"),
	}
}

func (s *webPushService) SendNotification(ctx context.Context, subscription *entity.PushSubscription, payload []byte) error {
	// Decode the subscription JSON
	var sub webpush.Subscription
	// The endpoint, p256dh, and auth are already base64 encoded from the client
	// and stored as such in the database. webpush-go expects them as is.
	sub.Endpoint = subscription.Endpoint
	sub.Keys.P256dh = subscription.P256dh
	sub.Keys.Auth = subscription.Auth

	// Send the notification
	resp, err := webpush.SendNotification(payload, &sub, &webpush.Options{
		VAPIDPublicKey:  s.publicKey,
		VAPIDPrivateKey: s.privateKey,
		TTL:             300,
		// You should replace this with your actual contact email
		Subscriber:      "mailto:your-email@example.com",
	})
	if err != nil {
		log.Printf("Failed to send notification: %v\n", err)
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to send notification: status code %d, body: %s", resp.StatusCode, string(body))
		return fmt.Errorf("failed to send notification: status code %d", resp.StatusCode)
	}

	return nil
}
