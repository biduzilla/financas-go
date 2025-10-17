package model

import "time"

type Category struct {
	ID        int64
	CreatedAt time.Time
	Name      string
	Type      TypeCategoria
	Color     string
	User      *User
	Deleted   bool
	Version   int
}

type TypeCategoria int

const (
	RECEITA TypeCategoria = iota + 1
	DESPESA
)

func (t TypeCategoria) String() string {
	switch t {
	case RECEITA:
		return "RECEITA"
	case DESPESA:
		return "DESPESA"
	default:
		return "Unknown"
	}
}

func TypeCategoriaFromString(s string) TypeCategoria {
	switch s {
	case "RECEITA":
		return RECEITA
	case "DESPESA":
		return DESPESA
	default:
		return 0
	}
}

type CategoryDTO struct {
	ID        *int64     `json:"category_id"`
	CreatedAt *time.Time `json:"-"`
	Name      *string    `json:"name"`
	Type      *string    `json:"type"`
	Color     *string    `json:"color"`
	User      *UserDTO   `json:"user"`
	Version   *int       `json:"version"`
}

func (m *CategoryDTO) ToModel() *Category {
	category := &Category{}
	if m.ID != nil {
		category.ID = *m.ID
	}

	if m.CreatedAt != nil {
		category.CreatedAt = *m.CreatedAt
	}

	if m.Name != nil {
		category.Name = *m.Name
	}

	if m.Type != nil {
		category.Type = TypeCategoriaFromString(*m.Type)
	}

	if m.Color != nil {
		category.Color = *m.Color
	}

	if m.User != nil {
		category.User = m.User.ToModel()
	}

	if m.Version != nil {
		category.Version = *m.Version
	}

	return category
}

func (m *Category) ToDTO() *CategoryDTO {
	category := &CategoryDTO{}

	category.ID = &m.ID
	category.CreatedAt = &m.CreatedAt
	category.Name = &m.Name
	typeStr := m.Type.String()
	category.Type = &typeStr
	category.Color = &m.Color
	category.User = m.User.ToDTO()
	category.Version = &m.Version

	return category
}
