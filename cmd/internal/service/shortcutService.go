package service

import (
	"database/sql"
	"fmt"
	"intellisys-backend/internal/types"
	"os"
	"strings"
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
	var oldImageUrl string
	_ = s.db.QueryRow("SELECT COALESCE(image_url, '') FROM shortcuts WHERE id=$1", id).Scan(&oldImageUrl)
	sqlStmt := `UPDATE shortcuts SET name=$1, url=$2, image_url=$3 WHERE id=$4 RETURNING id, name, url, image_url`

	fmt.Println(oldImageUrl)

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

	if oldImageUrl != "" && shortcut.ImageURL != oldImageUrl && !strings.HasPrefix(oldImageUrl, "http") {
		filePath := "." + oldImageUrl
		_ = os.Remove(filePath) // Deleta o arquivo antigo do disco
	}

	return *sResponse, nil
} // Fim UpdateShortcut

// Método para deletar um atalho pelo ID e remover o arquivo
func (s *ShortcutService) DeleteShortcut(id int) error {
	var imageUrl string
	query := "SELECT COALESCE(image_url, '') FROM shortcuts WHERE id = $1"
	err := s.db.QueryRow(query, id).Scan(&imageUrl)

	if err != nil {
		return err
	}

	sqlStmt := "DELETE FROM shortcuts WHERE id=$1"
	_, err = s.db.Exec(sqlStmt, id)

	if err != nil {
		return err
	}

	if imageUrl != "" && !strings.HasPrefix(imageUrl, "http") {
		filePath := "." + imageUrl
		_ = os.Remove(filePath)
	}

	return nil
} // Fim DeleteShortcut
