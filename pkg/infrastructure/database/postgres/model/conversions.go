package model

import (
	"github.com/WebEngineeringGroupI/backend/pkg/domain/click"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

func ShortURLFromDomain(url *url.ShortURL) Shorturl {
	return Shorturl{
		Hash:    url.Hash,
		LongURL: url.OriginalURL.URL,
		IsValid: url.OriginalURL.IsValid,
	}
}

func ShortURLToDomain(shortURL Shorturl) *url.ShortURL {
	return &url.ShortURL{
		Hash: shortURL.Hash,
		OriginalURL: url.OriginalURL{
			URL:     shortURL.LongURL,
			IsValid: shortURL.IsValid,
		},
	}
}

func ClickDetailsFromDomain(click *click.Details) Clickdetails {
	return Clickdetails{
		Hash: click.Hash,
		IP:   click.IP,
	}
}

func ClickDetailsToDomain(clickModel *Clickdetails) *click.Details {
	return &click.Details{
		Hash: clickModel.Hash,
		IP:   clickModel.IP,
	}
}
