package types

type ShortcutsRequest struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

type ShortcutsResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

type ShortcutsInterface interface {
	CreateShortcut(shortcut ShortcutsRequest) (ShortcutsResponse, error)
	GetAllShortcuts() ([]ShortcutsResponse, error)
	UpdateShortcut(id int, shortcut ShortcutsRequest) (ShortcutsResponse, error)
	DeleteShortcut(id int) error
}
