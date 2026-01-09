package types

// Struct para pegar requisição de login do usuário
type UsersRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Struct para resposta de login do usuário
type UsersResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Interface para manipulação dos usuários
type UsersInterface interface {
	GetUserByUsername(user UsersRequest) (UsersResponse, error)
}
