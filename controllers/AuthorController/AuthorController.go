package AuthorController

import (
	"library/database"
	"library/models/AuthorModels"
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
				author
			WHERE
				name ILIKE $1 AND deleted_at is NULL
			ORDER BY id
			LIMIT $2 OFFSET $3
			`

	rows, err := db.Query(query, search, limit, page)

	if err != nil {
		return err
	}

	var data_author []AuthorModels.GetAuthor
	for rows.Next() {
		var author AuthorModels.GetAuthor

		err := rows.Scan(&author.Id, &author.Name, &author.CreatedAt, &author.UpdatedAt)

		if err != nil {
			return err
		}

		data_author = append(data_author, author)
	}

	var total_data int
	query_paginate := "SELECT COUNT(id) FROM author WHERE deleted_at IS NULL"

	err = db.QueryRow(query_paginate).Scan(&total_data)

	if err != nil {
		return err
	}

	response := models.Response{
		Data:     data_author,
		Message:  "Get author list successfully",
		Paginate: map[string]int{"total_data": total_data},
	}

	return ctx.JSON(http.StatusOK, response)
}

func Create(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	var author AuthorModels.CreateAuthor

	if err := ctx.Bind(&author); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&author); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	query := "INSERT INTO author (name) VALUES ($1) returning id, name"

	err := db.QueryRow(query, &author.Name).Scan(&author.Id, &author.Name)

	if err != nil {
		return err
	}

	response := models.ResponseDetail{
		Data:    author,
		Message: "Create author successfully",
	}

	return ctx.JSON(http.StatusCreated, response)
}

func Show(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	var author AuthorModels.GetAuthor

	query := "SELECT id,name,created_at,updated_at FROM author WHERE id = $1"

	err := db.QueryRow(query, id).Scan(&author.Id, &author.Name, &author.CreatedAt, &author.UpdatedAt)

	if err != nil {
		response := models.ResponseDetail{
			Message: "Author not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	response := models.ResponseDetail{
		Data:    author,
		Message: "Get author detail successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}

func Update(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	var author AuthorModels.UpdateAuthor

	if err := ctx.Bind(&author); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&author); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	query := "UPDATE author SET name = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 returning id,name,updated_at"

	err := db.QueryRow(query, &author.Name, id).Scan(&author.Id, &author.Name, &author.UpdatedAt)

	if err != nil {
		response := models.ResponseDetail{
			Message: "author not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	response := models.ResponseDetail{
		Data:    author,
		Message: "Author updated successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}

func Delete(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	query := "UPDATE author SET deleted_at = CURRENT_TIMESTAMP  WHERE id = $1"

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
			Message: "author not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	response := models.ResponseDetail{
		Message: "author deleted Successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}
