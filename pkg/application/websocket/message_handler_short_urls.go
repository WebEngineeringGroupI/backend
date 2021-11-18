package websocket

import (
	"encoding/json"
	"fmt"
)

func (m *messageHandler) handleShortURLsMessage(message []byte) (responseType int, response []byte) {
	var data shortURLDataIn
	err := json.Unmarshal(message, &data)
	if err != nil {
		return CloseMessage, nil
	}

	dataResponse := shortURLDataOut{
		ResponseType: ShortURLsRequest,
	}
	for _, longURL := range data.Request.URLs {
		shortURL, err := m.urlShortener.HashFromURL(longURL)
		if err != nil {
			dataResponse.Response.URLs = append(dataResponse.Response.URLs, fmt.Sprintf("unable to short URL: %s", err))
			continue
		}
		dataResponse.Response.URLs = append(dataResponse.Response.URLs, m.wholeURL.FromHash(shortURL.Hash))
	}
	bytes, err := json.Marshal(dataResponse)
	if err != nil {
		return CloseMessage, nil
	}

	return TextMessage, bytes
}
