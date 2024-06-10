package AuthorModels

import "time"

type GetAuthor struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type CreateAuthor struct {
	Id   int    `json:"id"`
	Name string `json:"name" validate:"required"`
}

type UpdateAuthor struct {
	Id        int        `json:"id"`
	Name      string     `json:"name" validate:"required"`
	UpdatedAt *time.Time `json:"updated_at"`
}
