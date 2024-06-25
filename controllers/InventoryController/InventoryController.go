package InventoryController

import (
	"library/database"
	"library/models/InventoryModels"
	"net/http"

	"library/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func GetAll(ctx echo.Context) error {

	db := database.Init()

	defer db.Close()

	// search := "%" + ctx.QueryParam("search") + "%"

	// limit := ctx.QueryParam("limit")

	// if limit == "" {
	// 	limit = "10"
	// }

	// page := ctx.QueryParam("page")

	// if page == "" {
	// 	page = "0"
	// }

	query := `
			SELECT 
				inventories.id, 
				books.title,
				inventories.entry_time,
				inventories.scrap_time,
				inventories.status, 
				inventories.created_at,
				inventories.updated_at 
			FROM 
				inventories
			JOIN 
				books ON books.id = inventories.book_id
			`

	rows, err := db.Query(query)

	if err != nil {
		return err
	}

	var data_inventory []InventoryModels.GetInventory
	for rows.Next() {
		var inventory InventoryModels.GetInventory

		err := rows.Scan(&inventory.Id, &inventory.Book, &inventory.EntryTime, &inventory.ScrapTime, &inventory.Status, &inventory.CreatedAt, &inventory.UpdatedAt)

		if err != nil {
			return err
		}

		data_inventory = append(data_inventory, inventory)
	}

	var total_data int
	query_paginate := "SELECT COUNT(id) FROM inventories"

	err = db.QueryRow(query_paginate).Scan(&total_data)

	if err != nil {
		return err
	}

	response := models.Response{
		Data:     data_inventory,
		Message:  "Get inventory list successfully",
		Paginate: map[string]int{"total_data": total_data},
	}

	return ctx.JSON(http.StatusOK, response)
}

func Create(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	var inventory InventoryModels.CreateInventory

	if err := ctx.Bind(&inventory); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&inventory); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	query := "INSERT INTO inventories (id, book_id, entry_time, status) VALUES ($1, $2, $3, $4) RETURNING id, book_id, entry_time, scrap_time, status, created_at, updated_at"

	err := db.QueryRow(query, &inventory.Id, &inventory.BookId, &inventory.EntryTime, &inventory.Status).Scan(&inventory.Id, &inventory.BookId, &inventory.EntryTime, &inventory.ScrapTime, &inventory.Status, &inventory.CreatedAt, &inventory.UpdatedAt)

	if err != nil {
		return err
	}

	response := models.ResponseDetail{
		Data:    inventory,
		Message: "Create inventory successfully",
	}

	return ctx.JSON(http.StatusCreated, response)
}

func Update(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	ids := ctx.QueryParams()["id"]

	var data_inventory []InventoryModels.UpdateInventory

	var inventory InventoryModels.UpdateInventory

	// Bind data untuk setiap iterasi
	if err := ctx.Bind(&inventory); err != nil {
		return err
	}

	for _, id := range ids {

		// Jalankan query untuk setiap ID
		query := "UPDATE inventories SET scrap_time = $1, status = $2 WHERE id = $3 RETURNING id, book_id, entry_time, scrap_time, status, created_at, updated_at"
		err := db.QueryRow(query, &inventory.ScrapTime, "scrap", id).Scan(&inventory.Id, &inventory.BookId, &inventory.EntryTime, &inventory.ScrapTime, &inventory.Status, &inventory.CreatedAt, &inventory.UpdatedAt)

		if err != nil {
			return err
		}

		data_inventory = append(data_inventory, inventory)
	}

	response := models.ResponseDetail{
		Data:    data_inventory,
		Message: "Inventory updated successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}
