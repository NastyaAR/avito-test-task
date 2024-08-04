package pkg

import (
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string, lg *zap.Logger) (string, error) {
	bytePassword := []byte(password)
	encryptedBytePassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		lg.Warn("user usecase: register error", zap.Error(err))
		return "", fmt.Errorf("user usecase: reqister error: %v", err.Error())
	}
	encryptedPassword := string(encryptedBytePassword)
	return encryptedPassword, nil
}

func IsEqualPasswords(encryptedPassword string, expectedPassword string) error {
	bytesEncryptedPassword := []byte(encryptedPassword)
	bytesExpectedPassword := []byte(expectedPassword)
	err := bcrypt.CompareHashAndPassword(bytesEncryptedPassword, bytesExpectedPassword)
	return err
}
