package discord

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type DiscordClient struct {
	config *DiscordClientConfig
}

func NewDiscordClient(config *DiscordClientConfig) *DiscordClient {
	return &DiscordClient{config: config}
}

func (c *DiscordClient) SendMessage(_message string) (err error) {
	message := map[string]string{
		"content": _message,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = http.Post(c.config.WebhookUrl, "application/json", bytes.NewBuffer(jsonMessage))
	return err
}
