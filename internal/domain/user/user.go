package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin      = "ADMIN"
	RoleNormalUser = "NORMAL_USER"
)

type User struct {
	Id        uuid.UUID
	Name      string
	Password  string
	Role      Role
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

var (
	ErrorUserNotFound      = errors.New("Usuário não encontrado!")
	ErrorPasswordShort     = errors.New("Senha muito curta, deve conter no minimo 6 caracteres!")
	ErrorInvalidRole       = errors.New("Role Inválida!")
	ErrorInvalidNome       = errors.New("Nome Inválido!")
	ErrorNameAlreadyExists = errors.New("Nome já existe!")
)

func NewUsuario(name, password string, role Role) (*User, error) {
	if err := isValidName(name); err != nil {
		return nil, err
	}

	if len(password) < 6 {
		return nil, ErrorPasswordShort
	}

	return &User{
		Name:     name,
		Password: password,
		Role:     role,
	}, nil
}

func (u *User) UpdateUsuario(name *string, password *string) error {
	if name != nil {
		if err := isValidName(*name); err != nil {
			return err
		}
		u.Name = *name
	}

	if password != nil && *password != "" {
		if len(*password) < 6 {
			return ErrorPasswordShort
		}
	}

	return nil
}

func isValidName(name string) error {
	name = strings.TrimSpace(name)

	if len(name) < 4 {
		return errors.New("nome deve ter pelo menos 4 caracteres")
	}

	if strings.Contains(name, " ") {
		return errors.New("nome não pode conter espaços")
	}

	lower := strings.ToLower(name)

	if lower == "admin" || lower == "administrador" {
		return errors.New("nome não permitido")
	}

	return nil
}

func ParseRole(role string) (Role, error) {
	switch role {
	case string(RoleAdmin):
		return RoleAdmin, nil
	case string(RoleNormalUser):
		return RoleNormalUser, nil
	default:
		return "", fmt.Errorf("tipo de role não encontrado")
	}
}
