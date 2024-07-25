package helper

func ValidateAuthType(authType string) bool {
	return authType == "basic" || authType == "digest"
}
