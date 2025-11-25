package req

import (
	"github.com/go-playground/validator/v10"
)

func IsValide[T any](payload T) error {
	valideta := validator.New()
	err := valideta.Struct(payload)
	return err
}
