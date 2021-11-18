package websocket

import (
	"encoding/json"

	"github.com/WebEngineeringGroupI/backend/pkg/domain"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

// The message types are defined in RFC 6455, section 11.8.
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

const (
	ShortURLsRequest = "short_urls"
)

type MessageHandler interface {
	HandleMessage(messageType int, message []byte) (responseType int, response []byte)
}

type messageHandler struct {
	urlShortener *url.SingleURLShortener
	wholeURL     *domain.WholeURL
}

func (m *messageHandler) HandleMessage(messageType int, message []byte) (responseType int, response []byte) {
	switch messageType {
	case PingMessage:
		return PongMessage, []byte("pong")
	case TextMessage:
		return m.handleTextMessage(message)
	}
	return CloseMessage, nil
}

func (m *messageHandler) handleTextMessage(message []byte) (responseType int, response []byte) {
	requestType, err := m.extractRequestType(message)
	if err != nil {
		return CloseMessage, nil
	}

	switch requestType {
	case ShortURLsRequest:
		return m.handleShortURLsMessage(message)
	}
	return CloseMessage, nil
}

func (m *messageHandler) extractRequestType(message []byte) (requestType string, err error) {
	var bareRequest = struct {
		RequestType string `json:"request_type"`
	}{}

	err = json.Unmarshal(message, &bareRequest)
	if err != nil {
		return "", err
	}

	return bareRequest.RequestType, nil
}

func NewMessageHandler(wholeURL *domain.WholeURL, shortener *url.SingleURLShortener) MessageHandler {
	return &messageHandler{
		wholeURL:     wholeURL,
		urlShortener: shortener,
	}
}
