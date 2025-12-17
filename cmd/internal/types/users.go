package types

type UsersRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UsersResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UsersInterface interface {
	GetUserByUsername(user UsersRequest) (UsersResponse, error)
}
