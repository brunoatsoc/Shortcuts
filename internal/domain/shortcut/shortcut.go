package shortcut

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Shortcut struct {
	Id           uuid.UUID
	Name         string
	Link         string
	ImageUrl     string
	ShortcutType ShortcutType
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var (
	ErrorEmptyName     = errors.New("O campo nome está vazio")
	ErrorEmptyImageUrl = errors.New("URL da imagem está vazio")
	ErrorEmptyLink     = errors.New("Link está vazio")
)

type ShortcutType string

const (
	ClientLink   ShortcutType = "CLIENT_LINK"
	InternalLink ShortcutType = "INTERNAL_LINK"
)

func NewShortcut(name, imageUrl, link string, shortcutType ShortcutType) (*Shortcut, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrorEmptyName
	}

	if strings.TrimSpace(imageUrl) == "" {
		return nil, ErrorEmptyImageUrl
	}

	if strings.TrimSpace(link) == "" {
		return nil, ErrorEmptyLink
	}

	return &Shortcut{
		Name:         name,
		ImageUrl:     imageUrl,
		Link:         link,
		ShortcutType: shortcutType,
	}, nil
}

func (s *Shortcut) UpdateShortcut(name, imageUrl, link *string) error {
	if name != nil {
		if strings.TrimSpace(*name) == "" {
			return ErrorEmptyName
		}

		s.Name = *name
	}

	if imageUrl != nil {
		s.ImageUrl = *imageUrl
	}

	if link != nil {
		s.Link = *link
	}

	return nil
}

func ParseShortcutType(shortcutType string) (ShortcutType, error) {
	switch shortcutType {
	case string(ClientLink):
		return ClientLink, nil
	case string(InternalLink):
		return InternalLink, nil
	default:
		return "", fmt.Errorf("tipo de shortcut não encontrado")
	}
}
