package UserModels

import "time"

type GetUser struct {
	Id             int        `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	RoleId         string     `json:"role_id"`
	Avatar         *string    `json:"avatar"`
	IdentityNumber *string    `json:"identity_number"`
	Gender         *string    `json:"gender"`
	BirthDate      *string    `json:"birth_date"`
	Address        *string    `json:"address"`
	PhoneNumber    *string    `json:"phone_number"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type CreateUser struct {
	Id       int    `json:"id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUser struct {
	Id             int        `json:"id" form:"id"`
	Name           string     `json:"name" form:"name"`
	Email          string     `json:"email" form:"email"`
	RoleId         string     `json:"role_id" form:"role_id"`
	Avatar         *string    `json:"avatar" form:"avatar"`
	IdentityNumber *string    `json:"identity_number" form:"identity_number"`
	Gender         *string    `json:"gender" form:"gender"`
	BirthDate      *string    `json:"birth_date" form:"birth_date"`
	Address        *string    `json:"address" form:"address"`
	PhoneNumber    *string    `json:"phone_number" form:"phone_number"`
	CreatedAt      *time.Time `json:"created_at" form:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" form:"updated_at"`
}
