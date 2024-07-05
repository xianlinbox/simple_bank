package util

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) (string, error) {
	Password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil	{
		return "", err
	}
	return string(Password), nil
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
