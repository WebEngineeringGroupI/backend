package click

import (
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Logger struct {
	repository ClickLoggerRepository
}

func NewLogger(repository ClickLoggerRepository) *Logger {
	return &Logger{repository: repository}
}

func (l *Logger) LogIP(shortURL *url.ShortURL, ip string) {
	l.repository.Save(shortURL, ip)
}

type ClickLoggerRepository interface {
	FindByHash(hash string) []string
	Save(url *url.ShortURL, ip string)
}
