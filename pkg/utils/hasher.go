package utils

import "golang.org/x/crypto/bcrypt"

// NormalizePassword converts the password to a byte slice.
func NormalizePassword(password string) []byte {
	return []byte(password)
}

// GeneratePassword func for hashing user password.
func GeneratePassword(password string) (string, error) {
	bytePassword := NormalizePassword(password)

	hash, err := bcrypt.GenerateFromPassword(bytePassword, 12)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// ComparePasswords func for comparing passwords.
func ComparePasswords(hashedPassword, unhashedPassword string) bool {
	hashedBytes := NormalizePassword(hashedPassword)
	unhashedBytes := NormalizePassword(unhashedPassword)

	return bcrypt.CompareHashAndPassword(hashedBytes, unhashedBytes) == nil
}
