package BorrowModels

import (
	"time"
)

type GetBorrow struct {
	Id             int        `json:"id"`
	UserId         int        `json:"user_id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	BorrowTime     *time.Time `json:"borrow_time"`
	DueDate        string     `json:"due_date"`
	ReturnTime     *time.Time `json:"return_time"`
	TotalBorrowing int        `json:"total_borrowing"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type Inventory struct {
	Id              string     `json:"id"`
	BookId          string     `json:"book_id"`
	Title           string     `json:"title"`
	PublicationYear int        `json:"publication_year"`
	Description     string     `json:"description"`
	Code            string     `json:"code"`
	Thumbnail       string     `json:"thumbnail"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type GetBorrowDetail struct {
	Id             int         `json:"id"`
	UserId         int         `json:"user_id"`
	Name           string      `json:"name"`
	Email          string      `json:"email"`
	BorrowTime     *time.Time  `json:"borrow_time"`
	DueDate        string      `json:"due_date"`
	ReturnTime     *time.Time  `json:"return_time"`
	TotalBorrowing int         `json:"total_borrowing"`
	Inventory      []Inventory `json:"inventories"`
	CreatedAt      *time.Time  `json:"created_at"`
	UpdatedAt      *time.Time  `json:"updated_at"`
}

type InventoryUpdate struct {
	BorrowId    int    `json:"borrow_id"`
	InventoryId string `json:"inventory_id"`
}

type CreateBorrow struct {
	Id         int        `json:"id"`
	UserId     int        `json:"user_id" validate:"required"`
	BorrowTime string     `json:"borrow_time" validate:"required"`
	DueDate    string     `json:"due_date" validate:"required"`
	ReturnTime *string    `json:"return_time"`
	Inventory  []string   `json:"inventories" validate:"required"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type Borrowing struct {
	Id         int        `json:"id"`
	UserId     int        `json:"user_id"`
	BorrowTime string     `json:"borrow_time" validate:"required"`
	DueDate    string     `json:"due_date"`
	ReturnTime *string    `json:"return_time"`
	Inventory  []string   `json:"inventories" validate:"required"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type BorrowingReturn struct {
	Id         int        `json:"id" validate:"required"`
	UserId     int        `json:"user_id"`
	BorrowTime string     `json:"borrow_time"`
	DueDate    string     `json:"due_date"`
	ReturnTime *string    `json:"return_time"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type UpdateBorrow struct {
	Id         int        `json:"id"`
	UserId     int        `json:"user_id" validate:"required"`
	BorrowTime string     `json:"borrow_time" validate:"required"`
	DueDate    string     `json:"due_date" validate:"required"`
	ReturnTime *string    `json:"return_time" validate:"required"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}
