package dto

import (
	"io"
	"mime/multipart"
)

type ShortcutRequestDTO struct {
	Name         string                `form:"name"`
	ShortcutType string                `form:"shortcut_type"`
	Image        *multipart.FileHeader `form:"image"`
	Link         string                `form:"link"`
}

type ShortcutRequestUpdateDTO struct {
	Name  *string               `form:"name"`
	Image *multipart.FileHeader `form:"image"`
	Link  *string               `form:"link"`
}

type ShortcutCreateDTO struct {
	Name         string
	ShortcutType string
	Image        io.ReadSeeker
	Link         string
}

type ShortcutUpdateDTO struct {
	Name  *string
	Image io.ReadSeeker
	Link  *string
}

type ShortcutResponseDTO struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	Link     string `json:"link"`
}
