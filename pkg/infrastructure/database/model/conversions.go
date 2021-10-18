package model

import (
	`github.com/WebEngineeringGroupI/backend/pkg/domain/url`
)

func ShortURLFromDomain(url *url.ShortURL) Shorturl {
	return Shorturl{
		Hash:    url.Hash,
		LongURL: url.LongURL,
	}
}

func ShortURLToDomain(shortURL Shorturl) *url.ShortURL {
	return &url.ShortURL{
		Hash:    shortURL.Hash,
		LongURL: shortURL.LongURL,
	}
}
