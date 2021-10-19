package model

import (
	"github.com/WebEngineeringGroupI/backend/pkg/domain/click"
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

func ClickDetailsFromDomain(click *click.ClickDetails) Clickdetails {
	return Clickdetails{
		Hash: click.Hash,
		Ip: click.Ip,
	}
}

func ClickDetailsToDomain(clickModel *Clickdetails) *click.ClickDetails {
	return &click.ClickDetails {
		Hash: clickModel.Hash,
		Ip: clickModel.Ip,
	}
}
