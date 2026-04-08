package dto

type UserRequestDTO struct {
	Name     *string `json:"name"`
	Password *string `json:"password"`
}

type UserResponseDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type LoginRequestDTO struct {
	Name     *string `json:"name"`
	Password *string `json:"password"`
}
