package shortcut

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
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
	ErrorInvalidName   = errors.New("O campo invalido")
	ErrorEmptyImageUrl = errors.New("URL da imagem está vazio")
	ErrorInvalidLink   = errors.New("Link invalido")
)

type ShortcutType string

const (
	ClientLink   ShortcutType = "CLIENT_LINK"
	InternalLink ShortcutType = "INTERNAL_LINK"
)

func NewShortcut(name, imageUrl, link string, shortcutType ShortcutType) (*Shortcut, error) {
	if err := validateName(name); err != nil {
		return nil, ErrorInvalidName
	}

	if err := validateLink(link); err != nil {
		return nil, ErrorInvalidLink
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
		if err := validateName(*name); err != nil {
			return ErrorInvalidName
		}
		s.Name = *name
	}

	if link != nil {
		if err := validateLink(*link); err != nil {
			return ErrorInvalidLink
		}
		s.Link = *link
	}

	if imageUrl != nil {
		s.ImageUrl = *imageUrl
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

func validateLink(rawURL string) error {
	if len(rawURL) >= 2000 {
		return fmt.Errorf("url muito longa")
	}

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("formato de URL invalido")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("apenas URLs http:// ou https:// são permitidas")
	}

	if parsed.Host == "" {
		return errors.New("URL sem host válido")
	}

	if strings.ContainsAny(rawURL, " \t\n") {
		return errors.New("URL contém espaços inválidos")
	}

	return nil
}

func validateName(name string) error {
	var validNameRegex = regexp.MustCompile(`^[\p{L}\p{N}\s\-_.,!?()']+$`)

	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New("nome não pode ser vazio")
	}

	if len([]rune(name)) > 30 {
		return errors.New("nome muito longo, máximo 30 caracteres")
	}

	if !validNameRegex.MatchString(name) {
		return errors.New("nome contém caracteres inválidos")
	}

	return nil
}
