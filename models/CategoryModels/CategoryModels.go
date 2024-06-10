package CategoryModels

import "time"

type GetCategory struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type CreateCategory struct {
	Id   int    `json:"id"`
	Name string `json:"name" validate:"required"`
}

type UpdateCategory struct {
	Id        int        `json:"id"`
	Name      string     `json:"name" validate:"required"`
	UpdatedAt *time.Time `json:"updated_at"`
}
