package redirect

import (
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"strings"
)

type Redirector struct {
	repository url.ShortURLRepository
}

func NewRedirector(repository url.ShortURLRepository) *Redirector {
	return &Redirector{repository: repository}
}

func (r *Redirector) ReturnOriginalURL (hash string) string {
	//FindByHash
	aLongURL := r.repository.FindByHash(hash).LongURL
	if !strings.HasPrefix(aLongURL, "http://") && !strings.HasPrefix(aLongURL, "https://") {
		return ""
	}
	return aLongURL
}
