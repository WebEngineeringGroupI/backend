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

func LoadBalancedURLFromDomain(aURL *url.LoadBalancedURL) LoadBalancedUrlList {
	result := make(LoadBalancedUrlList, 0, len(aURL.LongURLs))
	for _, originalURL := range aURL.LongURLs {
		result = append(result, LoadBalancedUrl{
			Hash:        aURL.Hash,
			OriginalURL: originalURL.URL,
			IsValid:     originalURL.IsValid,
		})
	}
	return result
}

func LoadBalancedURLToDomain(aURL LoadBalancedUrlList) *url.LoadBalancedURL {
	result := &url.LoadBalancedURL{
		Hash:     aURL[0].Hash,
		LongURLs: []url.OriginalURL{},
	}
	for _, urlEntry := range aURL {
		result.LongURLs = append(result.LongURLs, url.OriginalURL{
			URL:     urlEntry.OriginalURL,
			IsValid: urlEntry.IsValid,
		})
	}
	return result
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
