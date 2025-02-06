package validatorutil

import (
	"sync/atomic"

	"github.com/go-playground/validator/v10"
)

var defaultValidate atomic.Pointer[validator.Validate]

func Default() *validator.Validate {
	if result := defaultValidate.Load(); result != nil {
		return result
	}

	result := validator.New(validator.WithRequiredStructEnabled())
	defaultValidate.Store(result)

	return result
}

func FormatErrors(err error) []string {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	result := []string{}

	for _, err := range errs {
		result = append(result, err.Error())
	}

	return result
}
