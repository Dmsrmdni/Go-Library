package CategoryController

import (
	"library/database"
	"library/models/CategoryModels"
	"net/http"

	"library/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func GetAll(ctx echo.Context) error {

	db := database.Init()

	defer db.Close()

	search := "%" + ctx.QueryParam("search") + "%"

	limit := ctx.QueryParam("limit")

	if limit == "" {
		limit = "10"
	}

	page := ctx.QueryParam("page")

	if page == "" {
		page = "0"
	}

	query := `
			SELECT 
				id, 
				name,
				created_at,
				updated_at 
			FROM 
				categories
			WHERE
				name ILIKE $1 AND deleted_at is NULL
			ORDER BY id
			LIMIT $2 OFFSET $3
			`

	rows, err := db.Query(query, search, limit, page)

	if err != nil {
		return err
	}

	var data_categories []CategoryModels.GetCategory
	for rows.Next() {
		var category CategoryModels.GetCategory

		err := rows.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)

		if err != nil {
			return err
		}

		data_categories = append(data_categories, category)
	}

	var total_data int
	query_paginate := "SELECT COUNT(id) FROM categories WHERE deleted_at IS NULL"

	err = db.QueryRow(query_paginate).Scan(&total_data)

	if err != nil {
		return err
	}

	response := models.Response{
		Data:     data_categories,
		Message:  "Get categories list successfully",
		Paginate: map[string]int{"total_data": total_data},
	}

	return ctx.JSON(http.StatusOK, response)
}

func Create(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	var category CategoryModels.CreateCategory

	if err := ctx.Bind(&category); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&category); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	query := "INSERT INTO categories (name) VALUES ($1) returning id, name"

	err := db.QueryRow(query, &category.Name).Scan(&category.Id, &category.Name)

	if err != nil {
		return err
	}

	response := models.ResponseDetail{
		Data:    category,
		Message: "Create category successfully",
	}

	return ctx.JSON(http.StatusCreated, response)
}

func Show(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	var category CategoryModels.GetCategory

	query := "SELECT id,name,created_at,updated_at FROM categories WHERE id = $1"

	err := db.QueryRow(query, id).Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		response := models.ResponseDetail{
			Message: "Categories not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	response := models.ResponseDetail{
		Data:    category,
		Message: "Get categories detail successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}

func Update(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	var category CategoryModels.UpdateCategory

	if err := ctx.Bind(&category); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&category); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	query := "UPDATE categories SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 returning id,name,updated_at"

	err := db.QueryRow(query, &category.Name, id).Scan(&category.Id, &category.Name, &category.UpdatedAt)

	if err != nil {
		response := models.ResponseDetail{
			Message: "Category not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	response := models.ResponseDetail{
		Data:    category,
		Message: "Category updated successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}

func Delete(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	query := "UPDATE categories SET deleted_at = CURRENT_TIMESTAMP  WHERE id = $1"

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
			Message: "Categories not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	response := models.ResponseDetail{
		Message: "Categories deleted Successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}
