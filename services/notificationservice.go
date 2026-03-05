package services

import (
	"context"

	"firebase.google.com/go/v4/messaging"
	"github.com/cuappdev/hustle-backend/auth"
	"github.com/cuappdev/hustle-backend/models"
)

type NotificationPayload struct {
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Data  map[string]string `json:"data,omitempty"`
}

// sends notification to all user's devices
func SendToUser(userID uint, payload NotificationPayload) error {
	tokens, err := models.GetUserTokens(userID)
	if err != nil || len(tokens) == 0 {
		return err
	}

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: payload.Title,
			Body:  payload.Body,
		},
		Data:   payload.Data,
		Tokens: tokens,
	}

	client := auth.GetMessagingClient()
	response, err := client.SendEachForMulticast(context.Background(), message)
	if err != nil {
		return err
	}

	// Remove invalid tokens
	if response.FailureCount > 0 {
		for idx, resp := range response.Responses {
			if !resp.Success {
				models.DeleteToken(tokens[idx])
			}
		}
	}

	return nil
}

// sends to a specific token
func SendToToken(token string, payload NotificationPayload) error {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: payload.Title,
			Body:  payload.Body,
		},
		Data:  payload.Data,
		Token: token,
	}

	client := auth.GetMessagingClient()
	_, err := client.Send(context.Background(), message)

	if err != nil {
		// Remove invalid token
		models.DeleteToken(token)
	}

	return err
}
