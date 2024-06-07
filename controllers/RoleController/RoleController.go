package RoleController

import (
	"library/database"
	"library/models/RoleModels"
	"net/http"

	"library/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func GetAll(ctx echo.Context) error {

	db := database.Init()

	defer db.Close()

	query := "SELECT id, name FROM roles ORDER BY id"

	rows, err := db.Query(query)

	if err != nil {
		return err
	}

	var data_roles []RoleModels.GetRole
	for rows.Next() {
		var role RoleModels.GetRole

		err := rows.Scan(&role.Id, &role.Name)

		if err != nil {
			return err
		}

		data_roles = append(data_roles, role)
	}

	var total_data int
	query_paginate := "SELECT COUNT(id) FROM roles"

	err = db.QueryRow(query_paginate).Scan(&total_data)

	if err != nil {
		return err
	}

	response := models.Response{
		Data:     data_roles,
		Message:  "Get role list successfully",
		Paginate: map[string]int{"total_data": total_data},
	}

	return ctx.JSON(http.StatusOK, response)
}

func Create(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	var role RoleModels.CreateRole

	if err := ctx.Bind(&role); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&role); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	query := "INSERT INTO roles (name) VALUES ($1) returning id, name"

	err := db.QueryRow(query, &role.Name).Scan(&role.Id, &role.Name)

	if err != nil {
		return err
	}

	response := models.ResponseDetail{
		Data:    role,
		Message: "Create roles successfully",
	}

	return ctx.JSON(http.StatusCreated, response)
}

func Show(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	var role RoleModels.GetRole

	if err := ctx.Bind(&role); err != nil {
		return err
	}

	query := "SELECT id,name FROM roles WHERE id = $1"

	err := db.QueryRow(query, id).Scan(&role.Id, &role.Name)

	if err != nil {
		response := models.ResponseDetail{
			Message: "Role not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	response := models.ResponseDetail{
		Data:    role,
		Message: "Get role detail successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}

func Update(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	var role RoleModels.UpdateRole

	if err := ctx.Bind(&role); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&role); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	query := "UPDATE roles SET name = $1 WHERE id = $2 returning id,name"

	err := db.QueryRow(query, &role.Name, id).Scan(&role.Id, &role.Name)

	if err != nil {
		response := models.ResponseDetail{
			Message: "Role not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	response := models.ResponseDetail{
		Data:    role,
		Message: "Role updated successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}

func Delete(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	query := "DELETE FROM roles WHERE id = $1"

	result, err := db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		response := models.ResponseDetail{
			Message: "Role not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	response := models.ResponseDetail{
		Message: "Roles deleted Successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}
