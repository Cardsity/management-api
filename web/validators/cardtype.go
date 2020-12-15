package validators

import "github.com/go-playground/validator/v10"

// Represents the type of a card. Can be either a black or a white card.
type CardType uint

const (
	BlackCard CardType = iota + 1
	WhiteCard
)

// Returns if the supplied CardType is valid. This means it should be a 1 (for a black card) or a 2 (for a white card).
func IsValidCardType(c CardType) bool {
	return c == BlackCard || c == WhiteCard
}

// A validator that checks if the field is a valid CardType. This is done by using the IsValidCardType method.
var cardTypeValidator validator.Func = func(fl validator.FieldLevel) bool {
	v, ok := fl.Field().Interface().(CardType)
	if ok {
		return IsValidCardType(v)
	}
	return false
}
