package discord

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type JSONMarshallerMock struct {
	mock.Mock
}

func (m *JSONMarshallerMock) Marshal(v any) ([]byte, error) {
	args := m.Called(v)
	return args.Get(0).([]byte), args.Error(1)
}

type HTTPClientMock struct {
	mock.Mock
}

func (m *HTTPClientMock) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	args := m.Called(url, contentType, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func Test_SendMessage_ReturnsErrorOnJsonError(t *testing.T) {
	config := NewDiscordClientConfig("dummyurl")
	jsonMarshallerMock := &JSONMarshallerMock{}
	jsonMarshallerMock.On("Marshal", map[string]string{
		"content": "some json string",
	}).Return(make([]byte, 0), errors.New("failed"))
	httpClientMock := &HTTPClientMock{}
	discordClient := NewDiscordClient(
		config,
		jsonMarshallerMock,
		httpClientMock,
	)

	err := discordClient.SendMessage("some json string")

	assert.Error(t, err)
	jsonMarshallerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func Test_SendMessage_ReturnsErrorOnNon200StatusCode(t *testing.T) {
	config := NewDiscordClientConfig("dummyurl")
	jsonMarshallerMock := &JSONMarshallerMock{}
	jsonMarshallerMock.On("Marshal", map[string]string{
		"content": "some json string",
	}).Return(make([]byte, 0), nil)
	httpClientMock := &HTTPClientMock{}
	httpResponse := http.Response{StatusCode: http.StatusBadRequest}
	httpClientMock.On("Post", "dummyurl", "application/json", bytes.NewBuffer([]byte(""))).Return(&httpResponse, nil)
	discordClient := NewDiscordClient(
		config,
		jsonMarshallerMock,
		httpClientMock,
	)

	err := discordClient.SendMessage("some json string")

	assert.ErrorContains(t, err, "error response from discord server: 400")
	jsonMarshallerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func Test_SendMessage_ReturnsErrorOnRequestError(t *testing.T) {
	config := NewDiscordClientConfig("dummyurl")
	jsonMarshallerMock := &JSONMarshallerMock{}
	jsonMarshallerMock.On("Marshal", map[string]string{
		"content": "some json string",
	}).Return(make([]byte, 0), nil)
	httpClientMock := &HTTPClientMock{}
	httpResponse := http.Response{StatusCode: http.StatusBadRequest}
	httpClientMock.On("Post", "dummyurl", "application/json", bytes.NewBuffer([]byte(""))).Return(&httpResponse, errors.New("failed post"))
	discordClient := NewDiscordClient(
		config,
		jsonMarshallerMock,
		httpClientMock,
	)

	err := discordClient.SendMessage("some json string")

	assert.ErrorContains(t, err, "failed post")
	jsonMarshallerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func Test_SendMessage_Passes(t *testing.T) {
	config := NewDiscordClientConfig("dummyurl")
	jsonMarshallerMock := &JSONMarshallerMock{}
	jsonMarshallerMock.On("Marshal", map[string]string{
		"content": "some json string",
	}).Return(make([]byte, 0), nil)
	httpClientMock := &HTTPClientMock{}
	httpResponse := http.Response{StatusCode: http.StatusOK}
	httpClientMock.On("Post", "dummyurl", "application/json", bytes.NewBuffer([]byte(""))).Return(&httpResponse, nil)
	discordClient := NewDiscordClient(
		config,
		jsonMarshallerMock,
		httpClientMock,
	)

	err := discordClient.SendMessage("some json string")

	assert.Nil(t, err)
	jsonMarshallerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}
