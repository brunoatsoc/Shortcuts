package shortcut

import (
	"context"
	"fmt"
	"intellisys-backend/internal/api/dto"
	"intellisys-backend/internal/domain/log"
	"intellisys-backend/internal/infra/config"

	"github.com/google/uuid"
)

type ShortcutService struct {
	repository   ShortcutRepository
	storageImage FileStorage
	urlServer    string
	log          log.Logger
}

func NewshortcutService(repository ShortcutRepository, storeImage FileStorage, cfg *config.Config, log log.Logger) *ShortcutService {
	return &ShortcutService{
		repository:   repository,
		storageImage: storeImage,
		urlServer:    "http://" + cfg.Server.Host + ":" + cfg.Server.Port,
		log:          log,
	}
}

func (s *ShortcutService) CreateShortcut(ctx context.Context, request dto.ShortcutCreateDTO) (*dto.ShortcutResponseDTO, error) {
	s.log.Info("Criando novo Shortcut")
	s.log.Info("Convertendo tipo do shortcut")

	shortcutType, err := ParseShortcutType(request.ShortcutType)
	if err != nil {
		s.log.Error("Erro ao converter tipo de Shortcut", "ID", request.ShortcutType, "error", err)
		return nil, err
	}

	var fileName string

	s.log.Info("Verificando imagem")

	if request.Image != nil {
		s.log.Info("Salvando imagem")
		fileName, err = s.storageImage.Save(request.Image)

		if err != nil {
			s.log.Error("Erro ao salvar imagem", "error", err)
			return nil, err
		}
	}

	newShortcut, err := NewShortcut(request.Name, fileName, request.Link, shortcutType)

	if err != nil {
		s.log.Error("Erro ao criar nova instância Shortcut", "error", err)
		s.log.Info("Deletando imagem salva no servidor")

		_ = s.storageImage.Delete(fileName)
		return nil, err
	}

	shortcutSaved, err := s.repository.Store(ctx, newShortcut)

	s.log.Info("Salvando Shortcut no banco de dados")

	if err != nil {
		s.log.Error("Erro ao salvar Shortcut no banco", "error", err)
		s.log.Info("Deletando imagem salva no servidor")

		_ = s.storageImage.Delete(fileName)
		return nil, err
	}

	s.log.Info("Criação de Shortcut concluido com sucesso")

	return &dto.ShortcutResponseDTO{
		Id:       shortcutSaved.Id.String(),
		Name:     shortcutSaved.Name,
		Link:     shortcutSaved.Link,
		ImageURL: "/uploads/" + shortcutSaved.ImageUrl,
	}, nil
}

func (s *ShortcutService) UpdateShortcut(ctx context.Context, request dto.ShortcutUpdateDTO, id string) error {
	s.log.Info("Atualizando Shortcut")
	s.log.Info("Convertendo ID do Shortcut")

	shortcutId, err := uuid.Parse(id)
	if err != nil {
		s.log.Error("Erro ao converter string em UUID", "ID", id, "error", err)

		return fmt.Errorf("UUID inválido: %w", err)
	}

	s.log.Info("Procurando Shortcut para ser atualizado")

	shortcutFound, err := s.repository.FindById(ctx, shortcutId)
	if err != nil {
		s.log.Error("Erro ao procurar Shortcut no banco de dados", "error", err)

		return err
	}

	var newImageName *string
	oldImage := shortcutFound.ImageUrl

	s.log.Info("Verificando existência da imagem")

	if request.Image != nil {
		s.log.Info("Atualizando imagem no servidor")

		fileName, err := s.storageImage.Save(request.Image)
		if err != nil {
			s.log.Error("Erro ao salvar imagem no servidor", "error", err)

			return err
		}

		s.log.Info("Imagem atualizada com sucesso")

		newImageName = &fileName
	}

	s.log.Info("Atualizando informações do Shortcut")

	err = shortcutFound.UpdateShortcut(request.Name, newImageName, request.Link)
	if err != nil {
		s.log.Error("Erro ao atualizar informações do Shortuct", "error", err)

		if newImageName != nil {
			s.log.Info("Deletando imagem nova do servidor")

			_ = s.storageImage.Delete(*newImageName)
		}
		return err
	}

	s.log.Info("Atualizando Shortcut no banco de dados")

	err = s.repository.Update(ctx, shortcutFound)
	if err != nil {
		s.log.Error("Erro ao atualizar Shortcut no banco de dados", "error", err)

		if newImageName != nil {
			s.log.Info("Deletando imagem nova do servidor")

			_ = s.storageImage.Delete(*newImageName)
			shortcutFound.ImageUrl = oldImage
		}
		return err
	}

	if newImageName != nil && oldImage != "" {
		s.log.Info("Deletando imagem antiga do servidor")

		_ = s.storageImage.Delete(oldImage)
	}

	s.log.Info("Atualização de Shortcut concluida com sucesso")

	return nil
}

func (s *ShortcutService) GetAllShortcuts(ctx context.Context) ([]*dto.ShortcutResponseDTO, error) {
	s.log.Info("Obtendo todos os Shrotcuts")
	s.log.Info("Obtendo Shortcuts do banco de dados")

	shortcuts, err := s.repository.FindAll(ctx)
	if err != nil {
		s.log.Error("Erro ao obter Shortcuts do banco de dados", "error", err)

		return nil, err
	}

	response := make([]*dto.ShortcutResponseDTO, 0, len(shortcuts))

	s.log.Info("Iterando pelos shortcuts obtidos do banco de dados")

	for _, sc := range shortcuts {
		var imageURL string

		if sc.ImageUrl != "" {
			imageURL = s.urlServer + "/uploads/" + sc.ImageUrl
		}

		response = append(response, &dto.ShortcutResponseDTO{
			Id:           sc.Id.String(),
			Name:         sc.Name,
			ShortcutType: string(sc.ShortcutType),
			Link:         sc.Link,
			ImageURL:     imageURL,
		})
	}

	s.log.Info("Shortcuts obtidos com sucesso")

	return response, nil
}

func (s *ShortcutService) DeleteShortcut(ctx context.Context, id string) error {
	s.log.Info("Deletando Shortcut")
	s.log.Info("Convertendo ID do Shortcut")

	shortcutId, err := uuid.Parse(id)
	if err != nil {
		s.log.Error("Erro ao converter string em UUID", "error", err)

		return fmt.Errorf("UUID inválido: %w", err)
	}

	s.log.Info("Procurando Shortcut para deleção")

	shortcutFound, err := s.repository.FindById(ctx, shortcutId)
	if err != nil {
		s.log.Error("Erro ao procurar Shortcut no banco de dados", "error", err)

		return err
	}

	s.log.Info("Deletando Shortcut do banco de dados")

	err = s.repository.Delete(ctx, shortcutId)
	if err != nil {
		s.log.Error("Erro ao deletar Shortcut do banco de dados")

		return err
	}

	if shortcutFound.ImageUrl != "" {
		s.log.Info("Deletando imagem do servidor")

		_ = s.storageImage.Delete(shortcutFound.ImageUrl)
	}

	s.log.Info("Shortcut deletado com sucesso")

	return nil
}
