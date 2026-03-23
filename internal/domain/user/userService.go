package user

import (
	"context"
	"fmt"
	"intellisys-backend/internal/api/dto"

	"github.com/google/uuid"
)

type UserService struct {
	repository     UserRepository
	passwordHasher PasswordHasher
}

func NewUserService(repository UserRepository, colaborador, passwordHasher PasswordHasher) *UserService {
	return &UserService{
		repository:     repository,
		passwordHasher: passwordHasher,
	}
}

func (s *UserService) UpdateUser(ctx context.Context, request dto.UserRequestDTO, id string) error {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("UUID do usuário inválido: %w", err)
	}

	existingUser, err := s.repository.FindById(ctx, userUUID)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return ErrorUserNotFound
	}

	if request.Name != nil {
		exists, err := s.repository.ExistsNameExcludingId(ctx, *request.Name, existingUser.Id)
		if err != nil {
			return err
		}

		if exists {
			return ErrorNameAlreadyExists
		}
	}

	if err := existingUser.UpdateUsuario(request.Name, request.Password); err != nil {
		return err
	}

	if request.Password != nil && *request.Password != "" {
		hashedPassword, err := s.passwordHasher.HashPassword(*request.Password)
		if err != nil {
			return fmt.Errorf("erro ao hashear senha: %w", err)
		}
		existingUser.Password = hashedPassword
	}

	return s.repository.Update(ctx, existingUser)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	userUUID, err := uuid.Parse(id)

	if err != nil {
		return fmt.Errorf("UUID do usuário inválido: %v", err)
	}

	existingUser, err := s.repository.FindById(ctx, userUUID)

	if existingUser == nil {
		return ErrorUserNotFound
	}

	if err != nil {
		return err
	}

	return s.repository.Delete(ctx, existingUser.Id)
}

func (s *UserService) GetUserById(ctx context.Context, id string) (*dto.UserResponseDTO, error) {
	userUUID, err := uuid.Parse(id)

	if err != nil {
		return nil, fmt.Errorf("UUID do usuário inválido: %v", err)
	}

	existingUser, err := s.repository.FindById(ctx, userUUID)

	if existingUser == nil {
		return nil, ErrorUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &dto.UserResponseDTO{
		Id:   existingUser.Id.String(),
		Name: existingUser.Name,
	}, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]*dto.UserResponseDTO, error) {
	u, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.UserResponseDTO, 0, len(u))

	for _, usr := range u {
		response = append(response, &dto.UserResponseDTO{
			Id:   usr.Id.String(),
			Name: usr.Name,
		})
	}

	return response, nil
}

func (s *UserService) ExistsUserById(ctx context.Context, id string) error {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("UUID do usuário inválido: %v", err)
	}

	exists, err := s.repository.ExistsId(ctx, userUUID)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("usuário não existe")
	}

	return nil
}
