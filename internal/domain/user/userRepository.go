package user

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, userId uuid.UUID) error
	FindById(ctx context.Context, userId uuid.UUID) (*User, error)
	FindAll(ctx context.Context) ([]*User, error)
	FindByName(ctx context.Context, name string) (*User, error)
	ExistsId(ctx context.Context, userId uuid.UUID) (bool, error)
	ExistsNameExcludingId(ctx context.Context, name string, userId uuid.UUID) (bool, error)
}
