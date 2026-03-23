package postgres

import (
	"context"
	"fmt"
	"intellisys-backend/internal/domain/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users
		SET name = $1,
			password = $2
		WHERE id = $3
	`

	cmdTag, err := r.pool.Exec(ctx, query,
		u.Name,
		u.Password,
		u.Id,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("usuário não encontrado")
	}

	return nil
}

func (r *UserRepository) FindByName(ctx context.Context, name string) (*user.User, error) {
	query := `
		SELECT id, name, password, role, createdat, updatedat
		FROM users
		WHERE name = $1
	`

	var u user.User

	err := r.pool.QueryRow(ctx, query, name).Scan(
		&u.Id,
		&u.Name,
		&u.Password,
		&u.Role,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("usuário não encontrado: %w", err)
	}

	return &u, nil
}

func (r *UserRepository) FindById(ctx context.Context, userId uuid.UUID) (*user.User, error) {
	query := `
		SELECT id, name, password, createdat, updatedat
		FROM users
		WHERE id = $1
	`

	var u user.User

	err := r.pool.QueryRow(ctx, query, userId).Scan(
		&u.Id,
		&u.Name,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("usuário não encontrado: %w", err)
	}

	return &u, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]*user.User, error) {
	query := `
		SELECT id, name, createdat, updatedat
		FROM users
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao reculperar usuários: %w", err)
	}
	defer rows.Close()

	var users []*user.User

	for rows.Next() {
		var u user.User

		err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear usuários: %w", err)
		}

		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar usuários: %w", err)
	}

	return users, nil
}

func (r *UserRepository) ExistsNameExcludingId(ctx context.Context, name string, userId uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM users
			WHERE name = $1 AND id <> $2
		)
	`

	var exists bool

	err := r.pool.QueryRow(ctx, query, name, userId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência do nome do usuário: %w", err)
	}

	return exists, nil
}

func (r *UserRepository) ExistsId(ctx context.Context, userId uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM users
			WHERE id = $1
		)
	`

	var exists bool

	err := r.pool.QueryRow(ctx, query, userId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência id: %w", err)
	}

	return exists, nil
}

func (r *UserRepository) Delete(ctx context.Context, userId uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	cmdTag, err := r.pool.Exec(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("usuário não encontrado")
	}

	return nil
}
