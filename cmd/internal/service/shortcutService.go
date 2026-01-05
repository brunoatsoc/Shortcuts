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

// CREATE
func (s *ShortcutService) CreateShortcut(shortcut types.ShortcutsRequest) (types.ShortcutsResponse, error) {
	// Adicionei image_url na query
	sqlStmt := `INSERT INTO shortcuts (name, url, description, image_url) VALUES ($1, $2, $3, $4) RETURNING id, name, url, description, image_url`

	sResponse := new(types.ShortcutsResponse)

	err := s.db.QueryRow(
		sqlStmt,
		shortcut.Name,
		shortcut.URL,
		shortcut.Description,
		shortcut.ImageURL, // Passando o novo campo
	).Scan(
		&sResponse.ID,
		&sResponse.Name,
		&sResponse.URL,
		&sResponse.Description,
		&sResponse.ImageURL, // Lendo o novo campo de volta
	)

	if err != nil {
		return types.ShortcutsResponse{}, err
	}

	return *sResponse, nil
}

// GET ALL
func (s *ShortcutService) GetAllShortcuts() ([]types.ShortcutsResponse, error) {
	// Adicionei COALESCE para evitar erro se o campo for NULL no banco antigo
	sqlStmt := `SELECT id, name, url, description, COALESCE(image_url, '') FROM shortcuts order by created_at desc`

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
			&shortcut.Description,
			&shortcut.ImageURL, // Lendo o novo campo
		)

		if err != nil {
			return nil, err
		}

		shortcuts = append(shortcuts, shortcut)
	}

	return shortcuts, nil
}

// UPDATE
func (s *ShortcutService) UpdateShortcut(id int, shortcut types.ShortcutsRequest) (types.ShortcutsResponse, error) {
	// Atualizei a query de Update
	sqlStmt := `UPDATE shortcuts SET name=$1, url=$2, description=$3, image_url=$4 WHERE id=$5 RETURNING id, name, url, description, image_url`

	sResponse := new(types.ShortcutsResponse)
	err := s.db.QueryRow(
		sqlStmt,
		shortcut.Name,
		shortcut.URL,
		shortcut.Description,
		shortcut.ImageURL, // Novo campo
		id,                // O ID agora é o 5º parametro ($5)
	).Scan(
		&sResponse.ID,
		&sResponse.Name,
		&sResponse.URL,
		&sResponse.Description,
		&sResponse.ImageURL,
	)

	if err != nil {
		return types.ShortcutsResponse{}, err
	}

	return *sResponse, nil
}

// DELETE (Não muda nada)
func (s *ShortcutService) DeleteShortcut(id int) error {
	sqlStmt := `DELETE FROM shortcuts WHERE id=$1`
	_, err := s.db.Exec(sqlStmt, id)
	return err
}
