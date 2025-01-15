package util

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func BcryptPassword(plainPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 11)
	return string(bytes), err
}

func BcryptCompare(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

// PasswordComplexityVerify 验证用户密码复杂度
func PasswordComplexityVerify(password string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(password) >= 8 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
