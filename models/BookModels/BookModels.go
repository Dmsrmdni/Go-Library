package BookModels

import "time"

type GetBook struct {
	Id              int        `json:"id"`
	Title           string     `json:"title"`
	Category        string     `json:"category"`
	PublicationYear int        `json:"publication_year"`
	Description     string     `json:"description"`
	Code            string     `json:"code"`
	Thumbnail       string     `json:"thumbnail"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type Author struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type GetBookDetail struct {
	Id              int        `json:"id"`
	Title           string     `json:"title"`
	Category        string     `json:"category"`
	PublicationYear int        `json:"publication_year"`
	Description     string     `json:"description"`
	Code            string     `json:"code"`
	Thumbnail       string     `json:"thumbnail"`
	Author          []Author   `json:"author"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type CreateBook struct {
	Id              int        `json:"id" form:"id"`
	Title           string     `json:"title" form:"title" validate:"required"`
	CategoryId      int        `json:"category_id" form:"category_id" validate:"required"`
	PublicationYear int        `json:"publication_year" form:"publication_year" validate:"required"`
	Description     string     `json:"description" form:"description"`
	Code            string     `json:"code" form:"code" validate:"required"`
	Thumbnail       string     `json:"thumbnail" form:"thumbnail"`
	AuthorId        []int      `json:"author_id" form:"author_id" validate:"required"`
	CreatedAt       *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at" form:"updated_at"`
}

type UpdateBook struct {
	Id              int        `json:"id" form:"id"`
	Title           string     `json:"title" form:"title" validate:"required"`
	CategoryId      int        `json:"category_id" form:"category_id" validate:"required"`
	PublicationYear int        `json:"publication_year" form:"publication_year" validate:"required"`
	Description     string     `json:"description" form:"description"`
	Code            string     `json:"code" form:"code" validate:"required"`
	Thumbnail       string     `json:"thumbnail" form:"thumbnail"`
	AuthorId        []int      `json:"author_id" form:"author_id"`
	CreatedAt       *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at" form:"updated_at"`
}
