package service

import (
	"database/sql"
	"intellisys-backend/internal/types"
)

type UsersService struct {
	db *sql.DB
}

func NewUsersService(db *sql.DB) *UsersService {
	return &UsersService{db: db}
}

// Método para obter usuário pelo nome de usuário
func (s *UsersService) GetUserByUsername(user types.UsersRequest) (types.UsersResponse, error) {
	sqlStmt := `SELECT id, username, password FROM users WHERE username=$1`

	uResponse := new(types.UsersResponse)

	err := s.db.QueryRow(
		sqlStmt,
		user.Username,
	).Scan(
		&uResponse.ID,
		&uResponse.Username,
		&uResponse.Password,
	)

	if err != nil {
		return types.UsersResponse{}, err
	}

	return *uResponse, nil
} // Fim GetUserByUsername
