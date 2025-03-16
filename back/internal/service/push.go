package service

import (
	"context"
	"encoding/json"
	"fmt"
	"loan-mgt/g-cram/internal/config"
	"loan-mgt/g-cram/internal/db"
	"net/http"

	webpush "github.com/SherClockHolmes/webpush-go"
)

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Icon  string `json:"icon"`
	Tag   string `json:"tag"`
}

type Subscription struct {
	Endpoint string   `json:"endpoint"`
	Keys     Keys     `json:"keys"`
	Encoding []string `json:"encoding,omitempty"`
	Detail   string   `json:"detail,omitempty"` // Added new field
}

type Keys struct {
	P256dh string `json:"p256dh"`
	Auth   string `json:"auth"`
}

type Options struct {
	TTL          int `json:"ttl"`
	VapidDetails struct {
		Subject    string `json:"subject"`
		PublicKey  string `json:"publicKey"`
		PrivateKey string `json:"privateKey"`
	} `json:"vapidDetails"`
}

// SendPush sends a web push notification
func SendPush(subscription Subscription, payload Notification, options Options) (*http.Response, error) {
	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	// Create webpush message
	s := &webpush.Subscription{
		Endpoint: subscription.Endpoint,
		Keys: webpush.Keys{
			P256dh: subscription.Keys.P256dh,
			Auth:   subscription.Keys.Auth,
		},
	}

	fmt.Println(s)

	// Send notification
	resp, err := webpush.SendNotification(
		payloadBytes,
		s,
		&webpush.Options{
			Subscriber:      options.VapidDetails.Subject,
			VAPIDPublicKey:  options.VapidDetails.PublicKey,
			VAPIDPrivateKey: options.VapidDetails.PrivateKey,
			TTL:             options.TTL,
		})

	if err != nil {
		return nil, fmt.Errorf("failed to send notification: %w", err)
	}

	return resp, nil
}

func PushNotification(db *db.Store, cfg *config.Config, id, body, tag string) error {
	user, err := db.GetUser(context.Background(), id)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	if user.Subscription.String == "" {
		fmt.Println("No subscription found")
		return nil
	}
	var sub Subscription

	if err := json.Unmarshal([]byte(user.Subscription.String), &sub); err != nil {
		return fmt.Errorf("error unmarshaling subscription: %w", err)
	}

	// Prepare notification payload
	payload := Notification{
		Title: "G-cram",
		Body:  body,
		Icon:  "https://ringwatchers.com/images/Main_Blank.png",
		Tag:   tag,
	}

	// Prepare options
	options := Options{
		TTL: 10000,
		VapidDetails: struct {
			Subject    string `json:"subject"`
			PublicKey  string `json:"publicKey"`
			PrivateKey string `json:"privateKey"`
		}{
			Subject:    cfg.VAPIDSubject,
			PublicKey:  cfg.VAPIDPublicKey,
			PrivateKey: cfg.VAPIDPrivateKey,
		},
	}

	// Extract endpoint ID for logging
	endpoint := sub.Endpoint
	endpointID := endpoint[len(endpoint)-8:]

	fmt.Printf("Endpoint: %s\n", endpoint)

	// Send push notification
	fmt.Println("Sending push notification")
	result, err := SendPush(sub, payload, options)
	if err != nil {
		fmt.Printf("Endpoint ID: %s\nError: %s\n", endpointID, err)
		return err
	}

	fmt.Printf("Endpoint ID: %s\nResult: %d\n", endpointID, result.StatusCode)
	return nil
}
