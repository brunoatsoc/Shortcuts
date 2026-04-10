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

func validateURL(rawURL string) error {
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
		return errors.New("nome muito longo, máximo 100 caracteres")
	}

	if !validNameRegex.MatchString(name) {
		return errors.New("nome contém caracteres inválidos")
	}

	return nil
}
