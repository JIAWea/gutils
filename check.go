package gutils

import "regexp"

func IsPhone(phone string) bool {
	var pattern = "^1[345789]{1}\\d{9}$"
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(phone)
}

func IsPassword(value string) bool {
	pattern := `^[a-zA-Z0-9_-]{6,32}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(value)
}

func IsEmail(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
