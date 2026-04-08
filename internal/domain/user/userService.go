package user

import (
	"context"
	"fmt"
	"intellisys-backend/internal/api/dto"
	"intellisys-backend/internal/domain/log"

	"github.com/google/uuid"
)

type UserService struct {
	repository     UserRepository
	passwordHasher PasswordHasher
	log            log.Logger
}

func NewUserService(repository UserRepository, colaborador, passwordHasher PasswordHasher, log log.Logger) *UserService {
	return &UserService{
		repository:     repository,
		passwordHasher: passwordHasher,
		log:            log,
	}
}

func (s *UserService) UpdateUser(ctx context.Context, request dto.UserRequestDTO, id string) error {
	s.log.Info("Atualizando Usuário")
	s.log.Info("Convertendo UUID do Usuário")

	userUUID, err := uuid.Parse(id)
	if err != nil {
		s.log.Error("Erro ao converter string em UUId", "ID", id, "error", err)
		return fmt.Errorf("UUID do usuário inválido: %w", err)
	}

	s.log.Info("Procurando Usuário para atualização")

	existingUser, err := s.repository.FindById(ctx, userUUID)
	if err != nil {
		s.log.Error("Erro ao procurar Usuário", "error", err)
		return err
	}

	if existingUser == nil {
		s.log.Warn("Usuário não encontrado")

		return ErrorUserNotFound
	}

	if request.Name != nil {
		s.log.Info("Verificando se nome já existe no banco de dados")

		exists, err := s.repository.ExistsNameExcludingId(ctx, *request.Name, existingUser.Id)
		if err != nil {
			s.log.Error("Erro ao procurar usuário")
			return err
		}

		if exists {
			s.log.Warn("Usuário já existe")
			return ErrorNameAlreadyExists
		}
	}

	s.log.Info("Atualizando informações do Usuário")

	if err := existingUser.UpdateUsuario(request.Name, request.Password); err != nil {
		s.log.Error("Erro ao atualizar informações do Usuário")
		return err
	}

	if request.Password != nil && *request.Password != "" {
		hashedPassword, err := s.passwordHasher.HashPassword(*request.Password)
		if err != nil {
			s.log.Error("Erro ao criptografar senha do usuário", "error", err)

			return fmt.Errorf("erro ao hashear senha: %w", err)
		}
		existingUser.Password = hashedPassword
	}

	s.log.Info("Usuário atualizado com sucesso")

	return s.repository.Update(ctx, existingUser)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	s.log.Info("Deletando Usuário")
	s.log.Info("Convertendo ID do Usuário")

	userUUID, err := uuid.Parse(id)
	if err != nil {
		s.log.Error("Erro ao converter string em UUID", "ID", id, "error", err)

		return fmt.Errorf("UUID do usuário inválido: %v", err)
	}

	s.log.Info("Procurando Usuário no banco de dados")

	existingUser, err := s.repository.FindById(ctx, userUUID)
	if existingUser == nil {
		s.log.Warn("Usuário não encontrado")

		return ErrorUserNotFound
	}

	if err != nil {
		s.log.Error("Erro ao procurar Usuário no banco de dados")

		return err
	}

	s.log.Info("Usuário deletado com sucesso")

	return s.repository.Delete(ctx, existingUser.Id)
}

func (s *UserService) GetUserById(ctx context.Context, id string) (*dto.UserResponseDTO, error) {
	s.log.Info("Obtendo Usuário pelo ID")
	s.log.Info("Convertendo ID do Usuário")

	userUUID, err := uuid.Parse(id)
	if err != nil {
		s.log.Error("Erro ao converter string em UUID", "ID", id, "error", err)

		return nil, fmt.Errorf("UUID do usuário inválido: %v", err)
	}

	s.log.Info("Procurando Usuário")

	existingUser, err := s.repository.FindById(ctx, userUUID)
	if existingUser == nil {
		s.log.Warn("Usuário não encontrado")

		return nil, ErrorUserNotFound
	}

	if err != nil {
		s.log.Error("Erro ao procurar Usuário")

		return nil, err
	}

	s.log.Info("Usuário encontrado com sucesso")

	return &dto.UserResponseDTO{
		Id:   existingUser.Id.String(),
		Name: existingUser.Name,
	}, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]*dto.UserResponseDTO, error) {
	s.log.Info("Obtendo todos os Usuários")

	u, err := s.repository.FindAll(ctx)
	if err != nil {
		s.log.Error("Erro ao obter todos os usuários do banco")

		return nil, err
	}

	response := make([]*dto.UserResponseDTO, 0, len(u))

	s.log.Info("Iterando pelos Usuários do banco de dados")

	for _, usr := range u {
		response = append(response, &dto.UserResponseDTO{
			Id:   usr.Id.String(),
			Name: usr.Name,
		})
	}

	s.log.Info("Usuários obtidos com sucesso")

	return response, nil
}

func (s *UserService) ExistsUserById(ctx context.Context, id string) error {
	s.log.Info("Verificando se o Usuário existe pelo ID")
	s.log.Info("Convertendo ID do Usuário")

	userUUID, err := uuid.Parse(id)
	if err != nil {
		s.log.Error("Erro ao converter string em UUID", "ID", id, "error", err)

		return fmt.Errorf("UUID do usuário inválido: %v", err)
	}

	s.log.Info("Verificando existência do Usuário pelo ID", "ID", id)

	exists, err := s.repository.ExistsId(ctx, userUUID)
	if err != nil {
		s.log.Error("Erro ao verificar existência do Usuário pelo ID", "ID", id, "error", err)

		return err
	}

	if !exists {
		s.log.Warn("Usuário não existe")

		return fmt.Errorf("usuário não existe")
	}

	s.log.Info("Usuário verificado com sucesso")

	return nil
}
