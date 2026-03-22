package shortcut

import (
	"context"

	"github.com/google/uuid"
)

type ShortcutRepository interface {
	Store(ctx context.Context, shortcut *Shortcut) (*Shortcut, error)
	FindById(ctx context.Context, id uuid.UUID) (*Shortcut, error)
	FindAll(ctx context.Context) ([]*Shortcut, error)
	Update(ctx context.Context, shortcut *Shortcut) error
	Delete(ctx context.Context, id uuid.UUID) error
}
