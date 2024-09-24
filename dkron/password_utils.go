package dkron

import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) (string, error) {
	// 生成盐值，bcrypt 会自动生成
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
