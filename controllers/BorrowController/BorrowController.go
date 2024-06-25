package BorrowController

import (
	"fmt"
	"library/database"
	"library/models/BorrowModels"
	"net/http"
	"strconv"
	"time"

	"library/models"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetAll(ctx echo.Context) error {

	db := database.Init()

	defer db.Close()

	user := "%" + ctx.QueryParam("user") + "%"

	limit := ctx.QueryParam("limit")

	status := ctx.QueryParam("status")

	if status == "true" {
		status = "IS NULL"
	} else if status == "false" {
		status = "IS NOT NULL"
	} else {
		status = "IS NOT NULL OR borrows.return_time IS NULL"
	}

	if limit == "" {
		limit = "10"
	}

	pageStr := ctx.QueryParam("page")
	page, _ := strconv.Atoi(pageStr)

	if page != 0 {
		page = page - 1
	}

	query := `
			SELECT 
				borrows.id, 
				borrows.user_id,
				users.name,
				users.email,
				borrows.borrow_time, 
				borrows.due_date, 
				borrows.return_time, 
				(SELECT COUNT(borrow_id) FROM borrow_inventory WHERE borrow_id = borrows.id) as total_borrowing,
				borrows.created_at,
				borrows.updated_at 
			FROM 
				borrows
			JOIN 
				users ON users.id = borrows.user_id
			WHERE
				users.name ILIKE $1 
			AND
				borrows.return_time ` + status + `
			AND
				borrows.deleted_at IS NULL
			ORDER BY borrows.id
			LIMIT $2 OFFSET $3
		`

	rows, err := db.Query(query, user, limit, page)

	if err != nil {
		return err
	}

	var data_borrow []BorrowModels.GetBorrow
	for rows.Next() {
		var borrow BorrowModels.GetBorrow

		err := rows.Scan(&borrow.Id, &borrow.UserId, &borrow.Name, &borrow.Email, &borrow.BorrowTime, &borrow.DueDate, &borrow.ReturnTime, &borrow.TotalBorrowing, &borrow.CreatedAt, &borrow.UpdatedAt)

		if err != nil {
			return err
		}

		data_borrow = append(data_borrow, borrow)
	}

	var total_data int

	query_paginate := `
					SELECT 
						COUNT(borrows.id)
					FROM 
						borrows
					JOIN 
						users ON users.id = borrows.user_id
					WHERE
						users.name ILIKE $1 
					AND
						borrows.deleted_at is NULL
					`

	err = db.QueryRow(query_paginate, user).Scan(&total_data)

	if err != nil {
		return err
	}

	response := models.Response{
		Data:     data_borrow,
		Message:  "Get borrowing list successfully",
		Paginate: map[string]int{"total_data": total_data},
	}

	return ctx.JSON(http.StatusOK, response)
}

func Create(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	var borrow BorrowModels.CreateBorrow

	if err := ctx.Bind(&borrow); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&borrow); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	var cek bool
	for _, inventory := range borrow.Inventory {
		query := "SELECT EXISTS (SELECT 1 FROM inventories WHERE status != 'available' AND id = $1)"
		err := db.QueryRow(query, inventory).Scan(&cek)

		if err != nil {
			return err
		}
	}

	if cek {
		response := models.ResponseDetail{
			Message: "Buku sedang di pinjam",
		}
		return ctx.JSON(http.StatusBadRequest, response)
	}

	query := "INSERT INTO borrows(user_id, borrow_time, due_date) VALUES ($1, $2 ,$3) RETURNING id, user_id, borrow_time, due_date, return_time, created_at, updated_at"

	err := db.QueryRow(query, &borrow.UserId, &borrow.BorrowTime, &borrow.DueDate).Scan(&borrow.Id, &borrow.UserId, &borrow.BorrowTime, &borrow.DueDate, &borrow.ReturnTime, &borrow.CreatedAt, &borrow.UpdatedAt)

	if err != nil {
		return err
	}

	for _, inventory := range borrow.Inventory {
		query := "INSERT INTO borrow_inventory(borrow_id, inventory_id) VALUES ($1, $2)"
		_, err := db.Exec(query, &borrow.Id, inventory)
		if err != nil {
			return err
		}

		query = "UPDATE inventories SET status = 'borrowed' WHERE id = $1"
		_, err = db.Exec(query, inventory)
		if err != nil {
			return err
		}
	}

	response := models.ResponseDetail{
		Data:    borrow,
		Message: "Create borrow successfully",
	}

	return ctx.JSON(http.StatusCreated, response)
}

func Show(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	var borrow BorrowModels.GetBorrowDetail

	query := `
			SELECT 
				borrows.id, 
				borrows.user_id,
				users.name,
				users.email,
				borrows.borrow_time, 
				borrows.due_date, 
				borrows.return_time, 
				(SELECT COUNT(borrow_id) FROM borrow_inventory) as total_borrowing,
				borrows.created_at,
				borrows.updated_at 
			FROM 
				borrows
			JOIN 
				users ON users.id = borrows.user_id
			WHERE
				borrows.id = $1
			`

	err := db.QueryRow(query, id).Scan(&borrow.Id, &borrow.UserId, &borrow.Name, &borrow.Email, &borrow.BorrowTime, &borrow.DueDate, &borrow.ReturnTime, &borrow.TotalBorrowing, &borrow.CreatedAt, &borrow.UpdatedAt)

	if err != nil {
		response := models.ResponseDetail{
			Message: "borrows not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	query = `
		SELECT 
			inventories.id,
			inventories.book_id,
			books.title,
			books.publication_year,
			books.description,
			books.code,
			books.thumbnail,
			inventories.created_at,
			inventories.updated_at
		FROM
			borrow_inventory
		JOIN
			inventories ON inventories.id = borrow_inventory.inventory_id
		JOIN
			books ON books.id = inventories.book_id
		WHERE
			borrow_inventory.borrow_id = $1
	`

	rows, err := db.Query(query, id)

	if err != nil {
		return err
	}

	var data_inventory []BorrowModels.Inventory
	for rows.Next() {
		var inventory BorrowModels.Inventory

		err := rows.Scan(&inventory.Id, &inventory.BookId, &inventory.Title, &inventory.PublicationYear, &inventory.Description, &inventory.Code, &inventory.Thumbnail, &inventory.CreatedAt, &inventory.UpdatedAt)

		if err != nil {
			return err
		}

		data_inventory = append(data_inventory, inventory)
	}

	borrow.Inventory = data_inventory

	response := models.ResponseDetail{
		Data:    borrow,
		Message: "Get borrow detail successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}

func Update(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	var borrow BorrowModels.UpdateBorrow

	if err := ctx.Bind(&borrow); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&borrow); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	query := "UPDATE borrows SET user_id = $1, borrow_time = $2, due_date = $3, return_time = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5 RETURNING id, user_id, borrow_time, due_date, return_time, created_at, updated_at"

	err := db.QueryRow(query, &borrow.UserId, &borrow.BorrowTime, &borrow.DueDate, &borrow.ReturnTime, id).Scan(&borrow.Id, &borrow.UserId, &borrow.BorrowTime, &borrow.DueDate, &borrow.ReturnTime, &borrow.CreatedAt, &borrow.UpdatedAt)

	if err != nil {
		response := models.ResponseDetail{
			Message: "borrows not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	query = `
			SELECT 
				borrow_id,
				inventory_id
			FROM
				borrow_inventory
			WHERE
				borrow_id = $1
		`
	rows, err := db.Query(query, id)

	if err != nil {
		return err
	}

	var data_inventory []BorrowModels.InventoryUpdate
	for rows.Next() {
		var inventory BorrowModels.InventoryUpdate

		err := rows.Scan(&inventory.BorrowId, &inventory.InventoryId)

		if err != nil {
			return err
		}

		data_inventory = append(data_inventory, inventory)
	}

	for _, inventory := range data_inventory {
		query := "UPDATE inventories SET status = 'available' WHERE id = $1"
		_, err := db.Exec(query, inventory.InventoryId)

		if err != nil {
			return err
		}
	}

	response := models.ResponseDetail{
		Data:    borrow,
		Message: "Update borrow successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}

func Delete(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	query := "UPDATE borrows SET deleted_at = CURRENT_TIMESTAMP  WHERE id = $1"

	result, err := db.Exec(query, id)

	if err != nil {
		return err
	}

	query = `
			SELECT 
				borrow_id,
				inventory_id
			FROM
				borrow_inventory
			WHERE
				borrow_id = $1
		`
	rows, err := db.Query(query, id)

	if err != nil {
		return err
	}

	var data_inventory []BorrowModels.InventoryUpdate
	for rows.Next() {
		var inventory BorrowModels.InventoryUpdate

		err := rows.Scan(&inventory.BorrowId, &inventory.InventoryId)

		if err != nil {
			return err
		}

		data_inventory = append(data_inventory, inventory)
	}

	for _, inventory := range data_inventory {
		query := "UPDATE inventories SET status = 'available' WHERE id = $1"
		_, err := db.Exec(query, inventory.InventoryId)

		if err != nil {
			return err
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		response := models.ResponseDetail{
			Message: "borrows not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	response := models.ResponseDetail{
		Message: "borrows deleted Successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}

func HistoryBorrowing(ctx echo.Context) error {
	db := database.Init()
	defer db.Close()

	// Extract the token from the context
	user := ctx.Get("user")
	token := user.(*jwt.Token)

	// Get the claims from the token
	claims := token.Claims.(jwt.MapClaims)

	id := claims["id"]

	limit := ctx.QueryParam("limit")

	status := ctx.QueryParam("status")

	if status == "true" {
		status = "IS NULL"
	} else if status == "false" {
		status = "IS NOT NULL"
	} else {
		status = "IS NOT NULL OR borrows.return_time IS NULL"
	}

	if limit == "" {
		limit = "10"
	}

	pageStr := ctx.QueryParam("page")
	page, _ := strconv.Atoi(pageStr)

	if page != 0 {
		page = page - 1
	}

	query := `
			SELECT 
				borrows.id, 
				borrows.user_id,
				users.name,
				users.email,
				borrows.borrow_time, 
				borrows.due_date, 
				borrows.return_time, 
				(SELECT COUNT(borrow_id) FROM borrow_inventory WHERE borrow_id = borrows.id) as total_borrowing,
				borrows.created_at,
				borrows.updated_at 
			FROM 
				borrows
			JOIN 
				users ON users.id = borrows.user_id
			WHERE
				borrows.user_id = $1 
			AND
				borrows.return_time ` + status + `
			AND
				borrows.deleted_at IS NULL
			ORDER BY borrows.id
			LIMIT $2 OFFSET $3
			`

	rows, err := db.Query(query, id, limit, page)

	if err != nil {
		return err
	}

	var data_borrow []BorrowModels.GetBorrowDetail
	for rows.Next() {
		var borrow BorrowModels.GetBorrowDetail

		err := rows.Scan(&borrow.Id, &borrow.UserId, &borrow.Name, &borrow.Email, &borrow.BorrowTime, &borrow.DueDate, &borrow.ReturnTime, &borrow.TotalBorrowing, &borrow.CreatedAt, &borrow.UpdatedAt)

		if err != nil {
			return err
		}

		query = `
			SELECT 
				inventories.id,
				inventories.book_id,
				books.title,
				books.publication_year,
				books.description,
				books.code,
				books.thumbnail,
				inventories.created_at,
				inventories.updated_at
			FROM
				borrow_inventory
			JOIN
				inventories ON inventories.id = borrow_inventory.inventory_id
			JOIN
				books ON books.id = inventories.book_id
			WHERE
				borrow_inventory.borrow_id = $1
		`

		rows, err := db.Query(query, &borrow.Id)

		if err != nil {
			return err
		}

		var data_inventory []BorrowModels.Inventory
		for rows.Next() {
			var inventory BorrowModels.Inventory

			err := rows.Scan(&inventory.Id, &inventory.BookId, &inventory.Title, &inventory.PublicationYear, &inventory.Description, &inventory.Code, &inventory.Thumbnail, &inventory.CreatedAt, &inventory.UpdatedAt)

			if err != nil {
				return err
			}

			data_inventory = append(data_inventory, inventory)
		}

		borrow.Inventory = data_inventory

		data_borrow = append(data_borrow, borrow)
	}

	var total_data int

	query_paginate := `
					SELECT 
						COUNT(borrows.id)
					FROM 
						borrows
					WHERE
						borrows.deleted_at is NULL
					AND
						borrows.user_id = $1
					AND
						borrows.return_time ` + status + `
					`

	err = db.QueryRow(query_paginate, id).Scan(&total_data)

	if err != nil {
		return err
	}

	response := models.Response{
		Data:     data_borrow,
		Message:  "Get history of borrowing books successfully",
		Paginate: map[string]int{"total_data": total_data},
	}

	return ctx.JSON(http.StatusOK, response)
}

func BorrowingBook(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	// Extract the token from the context
	user := ctx.Get("user")
	token := user.(*jwt.Token)

	// Get the claims from the token
	claims := token.Claims.(jwt.MapClaims)

	id := claims["id"]

	var borrow BorrowModels.Borrowing

	if err := ctx.Bind(&borrow); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&borrow); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	var cek bool
	for _, inventory := range borrow.Inventory {
		query := "SELECT EXISTS (SELECT 1 FROM inventories WHERE status != 'available' AND id = $1)"
		err := db.QueryRow(query, inventory).Scan(&cek)

		if err != nil {
			return err
		}
	}

	if cek {
		response := models.ResponseDetail{
			Message: "Buku sedang di pinjam",
		}
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Parse borrow time and calculate due date
	borrowTime, err := time.Parse("2006-01-02 15:04:05", borrow.BorrowTime)
	if err != nil {
		response := models.ResponseDetail{
			Message: "Invalid borrow_time format",
		}
		return ctx.JSON(http.StatusBadRequest, response)
	}
	dueDate := borrowTime.AddDate(0, 0, 7).Format("2006-01-02")

	// Insert into borrows table
	query := `
			INSERT INTO borrows(user_id, borrow_time, due_date)
			VALUES ($1, $2, $3)
			RETURNING id, user_id, borrow_time, due_date, return_time, created_at, updated_at`

	err = db.QueryRow(query, id, borrowTime, dueDate).Scan(&borrow.Id, &borrow.UserId, &borrow.BorrowTime, &borrow.DueDate, &borrow.ReturnTime, &borrow.CreatedAt, &borrow.UpdatedAt)

	if err != nil {
		return err
	}

	for _, inventory := range borrow.Inventory {
		query := "INSERT INTO borrow_inventory(borrow_id, inventory_id) VALUES ($1, $2)"
		_, err := db.Exec(query, &borrow.Id, inventory)
		if err != nil {
			return err
		}

		query = "UPDATE inventories SET status = 'borrowed' WHERE id = $1"
		_, err = db.Exec(query, inventory)
		if err != nil {
			return err
		}
	}

	response := models.ResponseDetail{
		Data:    borrow,
		Message: "Create borrow successfully",
	}

	return ctx.JSON(http.StatusCreated, response)
}

func ReturnBorrow(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	var borrow BorrowModels.BorrowingReturn

	if err := ctx.Bind(&borrow); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&borrow); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	query := `
		SELECT 
			borrows.id, 
			borrows.user_id,
			borrows.borrow_time, 
			borrows.due_date, 
			borrows.return_time, 
			borrows.created_at,
			borrows.updated_at 
		FROM 
			borrows
		WHERE
			borrows.id = $1
	`

	err := db.QueryRow(query, &borrow.Id).Scan(&borrow.Id, &borrow.UserId, &borrow.BorrowTime, &borrow.DueDate, &borrow.ReturnTime, &borrow.CreatedAt, &borrow.UpdatedAt)

	if err != nil {
		response := models.ResponseDetail{
			Message: "borrows not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	query = "UPDATE borrows SET return_time = CURRENT_TIMESTAMP WHERE id = $1"
	_, err = db.Exec(query, &borrow.Id)

	if err != nil {
		return err
	}

	layout := "2006-01-02T15:04:05Z"
	dueDate, err := time.Parse(layout, borrow.DueDate)
	if err != nil {
		fmt.Println("Error parsing due date:", err)
	}

	var message string

	daysLate := int(time.Since(dueDate).Hours() / 24)

	if daysLate > 0 {
		message = fmt.Sprintf("You're late returning the books for %d days, so you were fined Rp %d,-", daysLate, daysLate*1000)
	} else {
		message = "Thank you for returning the books on time"
	}

	query = `
			SELECT 
				borrow_id,
				inventory_id
			FROM
				borrow_inventory
			WHERE
				borrow_id = $1
		`
	rows, err := db.Query(query, borrow.Id)

	if err != nil {
		return err
	}

	var data_inventory []BorrowModels.InventoryUpdate
	for rows.Next() {
		var inventory BorrowModels.InventoryUpdate

		err := rows.Scan(&inventory.BorrowId, &inventory.InventoryId)

		if err != nil {
			return err
		}

		data_inventory = append(data_inventory, inventory)
	}

	for _, inventory := range data_inventory {
		query := "UPDATE inventories SET status = 'available' WHERE id = $1"
		_, err := db.Exec(query, inventory.InventoryId)

		if err != nil {
			return err
		}
	}

	response := models.ResponseDetail{
		Data:    borrow,
		Message: message,
	}

	return ctx.JSON(http.StatusCreated, response)
}
