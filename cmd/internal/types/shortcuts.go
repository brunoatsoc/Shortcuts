package types

// Struct para pegar requisição de atalho
type ShortcutsRequest struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	ImageURL string `json:"image_url" form:"image_url"`
}

// Struct para resposta de atalho
type ShortcutsResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	ImageURL string `json:"image_url"`
}

// Interface para manipulação dos atalhos
type ShortcutsInterface interface {
	CreateShortcut(shortcut ShortcutsRequest) (ShortcutsResponse, error)
	GetAllShortcuts() ([]ShortcutsResponse, error)
	UpdateShortcut(id int, shortcut ShortcutsRequest) (ShortcutsResponse, error)
	DeleteShortcut(id int) error
}
