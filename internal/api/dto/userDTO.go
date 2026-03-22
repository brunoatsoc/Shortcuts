package dto

type UserRequestDTO struct {
	Name     *string `form:"name"`
	Password *string `form:"password"`
}

type UserResponseDTO struct {
	Id   string `form:"id"`
	Name string `form:"name"`
}

type LoginRequestDTO struct {
	Name     *string `form:"name"`
	Password *string `form:"password"`
}
