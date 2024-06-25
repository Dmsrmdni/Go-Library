package InventoryModels

import "time"

type GetInventory struct {
	Id        string     `json:"id"`
	Book      string     `json:"book"`
	EntryTime *time.Time `json:"entry_time"`
	ScrapTime *time.Time `json:"scrap_time"`
	Status    string     `json:"status"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type CreateInventory struct {
	Id        string     `json:"id" validate:"required"`
	BookId    int        `json:"book_id" validate:"required"`
	EntryTime *string    `json:"entry_time" validate:"required"`
	ScrapTime *string    `json:"scrap_time"`
	Status    string     `json:"status" validate:"required"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type UpdateInventory struct {
	Id        string     `json:"id" validate:"required"`
	BookId    int        `json:"book_id" validate:"required"`
	EntryTime *string    `json:"entry_time" validate:"required"`
	ScrapTime *string    `json:"scrap_time"`
	Status    string     `json:"status" validate:"required"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
