package user

import (
	"context"
	"errors"
	"intellisys-backend/internal/api/dto"
	"intellisys-backend/internal/domain/log"
)

var ErrInvalidCredentials = errors.New("Senha ou email inválidos")

type AuthService struct {
	repository     UserRepository
	password       PasswordHasher
	tokenGenerator TokenGenerator
	log            log.Logger
}

func NewAuthService(repository UserRepository, password PasswordHasher, tokenGenerator TokenGenerator, log log.Logger) *AuthService {
	return &AuthService{
		repository:     repository,
		password:       password,
		tokenGenerator: tokenGenerator,
		log:            log,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, name, password string) (*dto.UserResponseDTO, *string, error) {
	s.log.Info("Autenticando Usuário")
	s.log.Info("Procurando usuário por Nome", "Nome", name)

	user, err := s.repository.FindByName(ctx, name)
	if err != nil {
		s.log.Error("Erro ao procurar usuário por Nome", "Nome", name, "error", err)

		return nil, nil, err
	}

	if user == nil {
		s.log.Warn("Credenciais fornecidar invalidas")

		return nil, nil, ErrInvalidCredentials
	}

	if !s.password.ComparePassword(user.Password, password) {
		s.log.Warn("Credenciais fornecidar invalidas")

		return nil, nil, ErrInvalidCredentials
	}

	s.log.Info("Gerando Token do Usuário")

	token, err := s.tokenGenerator.GenerateToken(user.Id.String(), string(user.Role))
	if err != nil {
		s.log.Error("Erro ao gerar Token do Usuário", "error", err)

		return nil, nil, err
	}

	s.log.Info("Usuário autenticado com sucesso")

	return &dto.UserResponseDTO{
		Id:   user.Id.String(),
		Name: user.Name,
	}, &token, nil
}
