package domain

import (
	"fmt"
	"strings"
)

type WholeURL struct {
	baseDomain string
}

func (w *WholeURL) FromHash(hash string) string {
	return fmt.Sprintf("%s/r/%s", strings.TrimSuffix(w.baseDomain, "/"), hash)
}

func NewWholeURL(baseDomain string) *WholeURL {
	return &WholeURL{baseDomain: baseDomain}
}
