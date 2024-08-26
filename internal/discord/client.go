package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type JSONMarshaller interface {
	Marshal(v any) ([]byte, error)
}

type DefaultJSONMarshaller struct{}

func (m DefaultJSONMarshaller) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

type HTTPClient interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

type DefaultHTTPClient struct{}

func (c DefaultHTTPClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return http.Post(url, contentType, body)
}

type DiscordClient struct {
	config         *DiscordClientConfig
	jsonMarshaller JSONMarshaller
	httpClient     HTTPClient
}

func NewDiscordClient(config *DiscordClientConfig, jsonMarshaller JSONMarshaller, httpClient HTTPClient) *DiscordClient {
	return &DiscordClient{config: config, jsonMarshaller: jsonMarshaller, httpClient: httpClient}
}

func (c *DiscordClient) SendMessage(_message string) (err error) {
	message := map[string]string{
		"content": _message,
	}

	jsonMessage, err := c.jsonMarshaller.Marshal(message)
	if err != nil {
		return err
	}

	response, err := c.httpClient.Post(c.config.WebhookUrl, "application/json", bytes.NewBuffer(jsonMessage))
	if err != nil {
		log.Println(err)
		return err
	}
	if response.StatusCode >= 300 {
		return fmt.Errorf("error response from discord server: %v", response.StatusCode)
	}
	return err
}
