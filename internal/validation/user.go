package validation

func ValidatePassword(password string) bool {
	return len(password) >= 6
}

func ValidateUsername(username string) bool {
	return len(username) > 0
}
