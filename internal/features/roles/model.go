package roles

type Role struct {
	ID          int    `json:"ID"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}
