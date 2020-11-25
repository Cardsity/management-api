package response

const (
	ErrorBadRequest                 = "ERR_BAD_REQUEST"
	ErrorDuplicateUsername          = "ERR_DUPLICATE_USERNAME"
	ErrorInternal                   = "ERR_INTERNAL"
	ErrorNotFound                   = "ERR_NOT_FOUND"
	ErrorForbidden                  = "ERR_FORBIDDEN"
	ErrorPasswordRequirementsNotMet = "ERR_PASSWORD_REQUIREMENTS_NOT_MET"
	ErrorNoValidBearerSupplied      = "ERR_NO_VALID_BEARER_SUPPLIED"
	ErrorCardTextInvalid            = "ERR_CARD_TEXT_INVALID"
	ErrorDeckOfficialButNotAdmin    = "ERR_DECK_OFFICIAL_BUT_NOT_ADMIN"
	ErrorCardAmountInvalid          = "ERR_CARD_AMOUNT_INVALID"
)
