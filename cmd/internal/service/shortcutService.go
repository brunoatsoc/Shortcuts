package service

import (
	"database/sql"
	"intellisys-backend/internal/types"
)

type ShortcutService struct {
	db *sql.DB
}

func NewShortcutService(db *sql.DB) *ShortcutService {
	return &ShortcutService{db: db}
}

// Método para criar um novo atalho
func (s *ShortcutService) CreateShortcut(shortcut types.ShortcutsRequest) (types.ShortcutsResponse, error) {
	sqlStmt := `INSERT INTO shortcuts (name, url, image_url) VALUES ($1, $2, $3) RETURNING id, name, url, image_url`

	sResponse := new(types.ShortcutsResponse)

	err := s.db.QueryRow(
		sqlStmt,
		shortcut.Name,
		shortcut.URL,
		shortcut.ImageURL,
	).Scan(
		&sResponse.ID,
		&sResponse.Name,
		&sResponse.URL,
		&sResponse.ImageURL,
	)

	if err != nil {
		return types.ShortcutsResponse{}, err
	}

	return *sResponse, nil
} // Fim CreateShortcut

// Método para obter todos os atalhos
func (s *ShortcutService) GetAllShortcuts() ([]types.ShortcutsResponse, error) {
	sqlStmt := `SELECT id, name, url, COALESCE(image_url, '') FROM shortcuts order by created_at desc`

	rows, err := s.db.Query(sqlStmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var shortcuts []types.ShortcutsResponse

	for rows.Next() {
		var shortcut types.ShortcutsResponse
		err := rows.Scan(
			&shortcut.ID,
			&shortcut.Name,
			&shortcut.URL,
			&shortcut.ImageURL,
		)

		if err != nil {
			return nil, err
		}

		shortcuts = append(shortcuts, shortcut)
	}

	return shortcuts, nil
} // Fim GetAllShortcuts

// Método para atualizar um atalho existente
func (s *ShortcutService) UpdateShortcut(id int, shortcut types.ShortcutsRequest) (types.ShortcutsResponse, error) {
	sqlStmt := `UPDATE shortcuts SET name=$1, url=$2, image_url=$3 WHERE id=$4 RETURNING id, name, url, image_url`

	sResponse := new(types.ShortcutsResponse)
	err := s.db.QueryRow(
		sqlStmt,
		shortcut.Name,
		shortcut.URL,
		shortcut.ImageURL,
		id,
	).Scan(
		&sResponse.ID,
		&sResponse.Name,
		&sResponse.URL,
		&sResponse.ImageURL,
	)

	if err != nil {
		return types.ShortcutsResponse{}, err
	}

	return *sResponse, nil
} // Fim UpdateShortcut

// Método para deletar um atalho pelo ID
func (s *ShortcutService) DeleteShortcut(id int) error {
	sqlStmt := `DELETE FROM shortcuts WHERE id=$1`
	_, err := s.db.Exec(sqlStmt, id)
	return err
} // Fim DeleteShortcut
