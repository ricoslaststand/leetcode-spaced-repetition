package services

import (
	"context"

	"leetcode-spaced-repetition/repositories"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository repositories.UserRepository
	logger         *zap.Logger
}

func NewAuthService(logger *zap.Logger) *AuthService {
	return &AuthService{logger: logger}
}

func (a AuthService) Login(ctx context.Context, email, password string) (bool, error) {
	passwordHash, err := a.userRepository.GetPasswordHashByEmail(ctx, email)
	if err != nil {
		a.logger.Error("failed to get password hash during login", zap.String("email", email), zap.Error(err))
		return false, err
	}

	if passwordHash == nil {
		return false, nil
	}

	result := a.comparePasswords(*passwordHash, []byte(password))
	return result, nil
}

func (a AuthService) Logout() {

}

func (a AuthService) RegisterUser(ctx context.Context, email, password string) error {
	hash, err := a.hashAndSaltPassword(password)
	if err != nil {
		a.logger.Error("failed to hash password during registration", zap.String("email", email), zap.Error(err))
		return err
	}

	if err := a.userRepository.CreateUser(ctx, email, hash); err != nil {
		a.logger.Error("failed to create user", zap.String("email", email), zap.Error(err))
		return err
	}

	return nil
}

func (a AuthService) hashAndSaltPassword(plainPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.MinCost)
	if err != nil {
		return "", nil
	}

	return string(hash), nil
}

func (a AuthService) comparePasswords(hashedPassword string, plainPassword []byte) bool {
	byteHash := []byte(hashedPassword)

	if err := bcrypt.CompareHashAndPassword(byteHash, plainPassword); err != nil {
		return false
	}

	return true
}
