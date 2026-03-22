package user

import (
	"context"
	"errors"
)

var ErrInvalidCredentials = errors.New("Senha ou email inválidos")

type AuthService struct {
	repository     UserRepository
	password       PasswordHasher
	tokenGenerator TokenGenerator
}

func NewAuthService(repository UserRepository, password PasswordHasher, tokenGenerator TokenGenerator) *AuthService {
	return &AuthService{
		repository:     repository,
		password:       password,
		tokenGenerator: tokenGenerator,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, name, password string) (*string, error) {
	user, err := s.repository.FindByName(ctx, name)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !s.password.ComparePassword(user.Password, password) {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokenGenerator.GenerateToken(user.Id.String(), string(user.Role))

	if err != nil {
		return nil, err
	}

	return &token, nil
}
