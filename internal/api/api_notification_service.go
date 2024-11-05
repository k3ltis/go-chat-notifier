package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	openapi "go-teams-notifier/internal/generated/openapi/openapi"
)

// Implemented by the receiver of notifications forwarded by the notification server
type NotificationReceiver interface {
	SendMessage(msg string) (err error)
}

type NotificationAPIService struct {
	notificationReceiver NotificationReceiver
}

func NewNotificationAPIService(notificationReceiver NotificationReceiver) *NotificationAPIService {
	return &NotificationAPIService{notificationReceiver}
}

// Handler for the /notification POST endpoint
func (s *NotificationAPIService) PostNotification(ctx context.Context, notification openapi.Notification) (openapi.ImplResponse, error) {
	fmt.Printf("Received %s", notification)

	switch notification.Type {
	case notificationStateName[StateInfo]:
		return openapi.Response(http.StatusOK, nil), nil
	case notificationStateName[StateWarning]:
		message := fmt.Sprintf("-- %s: %s --\n\n%s\n----\n", notification.Type, notification.Name, notification.Description)
		err := s.notificationReceiver.SendMessage(message)

		if err != nil {
			log.Printf("Error sending message: %v", err)
			return openapi.Response(http.StatusInternalServerError, nil), errors.New("error while processsing the request")
		}

		return openapi.Response(http.StatusOK, nil), nil
	default:
		return openapi.Response(http.StatusBadRequest, nil), errors.New("unknown notification type")
	}
}
