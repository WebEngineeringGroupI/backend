package url

import (
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
)

type LoadBalancedURLCreated struct {
	event.Base
	OriginalURLs []string
}

// TODO(fede): There is an event that comes from the message broker, from the network, that verifies a URL, implement it
type LoadBalancedURLVerified struct {
	event.Base
	VerifiedURL string
}

type ShortURLCreated struct {
	event.Base
	OriginalURL string
}

type ShortURLVerified struct {
	event.Base
}

type ShortURLClicked struct {
	event.Base
}
