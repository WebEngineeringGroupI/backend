package pipeline

import (
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Validator struct {
	validators []url.Validator
}

func (v *Validator) ValidateURLs(url []string) (bool, error) {
	for _, validator := range v.validators {
		areURLsValid, err := validator.ValidateURLs(url)
		if err != nil {
			return false, err
		}
		if !areURLsValid {
			return areURLsValid, err
		}
	}
	return true, nil
}

func NewValidator(validators ...url.Validator) *Validator {
	return &Validator{
		validators: validators,
	}
}
