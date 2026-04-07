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
	s.log.Info("Convertendo tipo do shortcut")
	shortcutType, err := ParseShortcutType(request.ShortcutType)

	if err != nil {
		return nil, err
	}

	fileName, err := s.storageImage.Save(request.Image)

	if err != nil {
		return nil, err
	}

	newShortcut, err := NewShortcut(request.Name, fileName, shortcutType)

	if err != nil {
		_ = s.storageImage.Delete(fileName)
		return nil, err
	}

	shortcutSaved, err := s.repository.Store(ctx, newShortcut)

	if err != nil {
		_ = s.storageImage.Delete(fileName)
		return nil, err
	}

	return &dto.ShortcutResponseDTO{
		Id:       shortcutSaved.Id.String(),
		Name:     shortcutSaved.Name,
		ImageURL: "/uploads/" + shortcutSaved.ImageUrl,
	}, nil
}

func (s *ShortcutService) UpdateShortcut(ctx context.Context, request dto.ShortcutUpdateDTO, id string) error {
	shortcutId, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("UUID inválido: %w", err)
	}

	shortcutFound, err := s.repository.FindById(ctx, shortcutId)
	if err != nil {
		return err
	}

	var newImageName *string
	oldImage := shortcutFound.ImageUrl

	if request.Image != nil {
		fileName, err := s.storageImage.Save(request.Image)
		if err != nil {
			return err
		}

		newImageName = &fileName
	}

	err = shortcutFound.UpdateShortcut(request.Name, newImageName)
	if err != nil {
		if newImageName != nil {
			_ = s.storageImage.Delete(*newImageName)
		}
		return err
	}

	err = s.repository.Update(ctx, shortcutFound)
	if err != nil {
		if newImageName != nil {
			_ = s.storageImage.Delete(*newImageName)
			shortcutFound.ImageUrl = oldImage
		}
		return err
	}

	if newImageName != nil && oldImage != "" {
		_ = s.storageImage.Delete(oldImage)
	}

	return nil
}

func (s *ShortcutService) GetAllShortcuts(ctx context.Context) ([]*dto.ShortcutResponseDTO, error) {
	s.log.Info("Obtendo Shortcuts", "TESTE")

	shortcuts, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.ShortcutResponseDTO, 0, len(shortcuts))

	for _, sc := range shortcuts {
		response = append(response, &dto.ShortcutResponseDTO{
			Id:       sc.Id.String(),
			Name:     sc.Name,
			ImageURL: s.urlServer + "/uploads/" + sc.ImageUrl,
		})
	}

	return response, nil
}

func (s *ShortcutService) DeleteShortcut(ctx context.Context, id string) error {
	shortcutId, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("UUID inválido: %w", err)
	}

	shortcutFound, err := s.repository.FindById(ctx, shortcutId)
	if err != nil {
		return err
	}

	err = s.repository.Delete(ctx, shortcutId)
	if err != nil {
		return err
	}

	if shortcutFound.ImageUrl != "" {
		_ = s.storageImage.Delete(shortcutFound.ImageUrl)
	}

	return nil
}
