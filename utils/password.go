package utils

// Checks if a password meets the password requirements. These are:
// - At least 8 characters long
func MeetsPasswordRequirements(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}
