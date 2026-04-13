package postgres

import (
	"context"
	"errors"
	"fmt"
	"intellisys-backend/internal/domain/shortcut"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShortcutRepository struct {
	pool *pgxpool.Pool
}

func NewShortcutRepository(pool *pgxpool.Pool) *ShortcutRepository {
	return &ShortcutRepository{pool: pool}
}

func (r *ShortcutRepository) Store(ctx context.Context, shortcut *shortcut.Shortcut) (*shortcut.Shortcut, error) {
	query := `
		INSERT INTO shortcuts (name, link, image_url, shortcut_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx,
		query,
		shortcut.Name,
		shortcut.Link,
		shortcut.ImageUrl,
		string(shortcut.ShortcutType),
	).Scan(
		&shortcut.Id,
		&shortcut.CreatedAt,
		&shortcut.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("erro ao inserir shortcut: %w", err)
	}

	return shortcut, nil
}

func (r *ShortcutRepository) FindById(ctx context.Context, id uuid.UUID) (*shortcut.Shortcut, error) {
	query := `
		SELECT id, name, link, image_url, shortcut_type, created_at, updated_at
		FROM shortcuts
		WHERE id = $1
	`

	var shortcut shortcut.Shortcut

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&shortcut.Id,
		&shortcut.Name,
		&shortcut.Link,
		&shortcut.ImageUrl,
		&shortcut.ShortcutType,
		&shortcut.CreatedAt,
		&shortcut.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("shortcut não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar shortcut: %w", err)
	}

	return &shortcut, nil
}

func (r *ShortcutRepository) FindAll(ctx context.Context) ([]*shortcut.Shortcut, error) {
	query := `
		SELECT id, name, link, image_url, shortcut_type, created_at, updated_at
		FROM shortcuts
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar shortcuts: %w", err)
	}
	defer rows.Close()

	shortcuts := make([]*shortcut.Shortcut, 0)

	for rows.Next() {
		var s shortcut.Shortcut
		var shortcutType string

		err := rows.Scan(
			&s.Id,
			&s.Name,
			&s.Link,
			&s.ImageUrl,
			&shortcutType,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler shortcut: %w", err)
		}

		s.ShortcutType = shortcut.ShortcutType(shortcutType)

		shortcuts = append(shortcuts, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração: %w", err)
	}

	return shortcuts, nil
}

func (r *ShortcutRepository) Update(ctx context.Context, shortcut *shortcut.Shortcut) error {
	query := `
		UPDATE shortcuts
		SET name = $1,
		    image_url = $2,
			link = $3,
		    shortcut_type = $4,
		    updated_at = NOW()
		WHERE id = $5
	`

	cmdTag, err := r.pool.Exec(
		ctx,
		query,
		shortcut.Name,
		shortcut.ImageUrl,
		shortcut.Link,
		string(shortcut.ShortcutType),
		shortcut.Id,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar shortcut: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("shortcut não encontrado")
	}

	return nil
}

func (r *ShortcutRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM shortcuts
		WHERE id = $1
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar shortcut: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("shortcut não encontrado")
	}

	return nil
}
