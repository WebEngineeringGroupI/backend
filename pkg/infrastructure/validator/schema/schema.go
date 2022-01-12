package schema

import (
	"context"
	"fmt"
	"strings"
)

type Validator struct {
	allowedPrefixes []string
}

func (v *Validator) ValidateURLs(ctx context.Context, urls []string) (bool, error) {
	for _, prefix := range v.allowedPrefixes {
		for _, url := range urls {
			if strings.HasPrefix(url, prefix) {
				return true, nil
			}
		}
	}
	return false, nil
}

func NewValidator(allowedSchemas ...string) *Validator {
	allowedPrefixes := make([]string, 0, len(allowedSchemas))
	for _, schema := range allowedSchemas {
		allowedPrefixes = append(allowedPrefixes, fmt.Sprintf("%s://", schema))
	}

	return &Validator{
		allowedPrefixes: allowedPrefixes,
	}
}
