package validators

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Registers custom validators in gins binding.Validator.
func RegisterValidators() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("cardtype", cardTypeValidator)
		if err != nil {
			return err
		}
	}
	return nil
}
