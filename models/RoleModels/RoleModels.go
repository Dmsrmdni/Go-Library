package RoleModels

type GetRole struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CreateRole struct {
	Id   int    `json:"id"`
	Name string `json:"name" validate:"required"`
}

type UpdateRole struct {
	Id   int    `json:"id"`
	Name string `json:"name" validate:"required"`
}
