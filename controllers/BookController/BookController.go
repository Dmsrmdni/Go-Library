package BookController

import (
	"io"
	"library/database"
	"library/models/BookModels"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"library/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func GetAll(ctx echo.Context) error {

	db := database.Init()

	defer db.Close()

	search := "%" + ctx.QueryParam("search") + "%"
	category := "%" + ctx.QueryParam("category") + "%"

	limit := ctx.QueryParam("limit")

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
				books.id, 
				books.title,
				categories.name as category,
				books.publication_year,
				books.description,
				books.code,
				books.thumbnail,
				books.created_at,
				books.updated_at 
			FROM 
				books
			JOIN 
				categories ON categories.id = books.category_id
			WHERE
				books.title ILIKE $1 
			AND
				categories.name ILIKE $2
			AND 
				books.deleted_at is NULL
			ORDER BY id
			LIMIT $3 OFFSET $4
			`

	rows, err := db.Query(query, search, category, limit, page)

	if err != nil {
		return err
	}

	var data_books []BookModels.GetBook
	for rows.Next() {
		var book BookModels.GetBook

		err := rows.Scan(&book.Id, &book.Title, &book.Category, &book.PublicationYear, &book.Description, &book.Code, &book.Thumbnail, &book.CreatedAt, &book.UpdatedAt)

		if err != nil {
			return err
		}

		data_books = append(data_books, book)
	}

	var total_data int

	query_paginate := `
					SELECT 
						COUNT(books.id)
					FROM 
						books
					JOIN 
						categories ON categories.id = books.category_id
					WHERE
						books.title ILIKE $1 
					AND
						categories.name ILIKE $2
					AND 
						books.deleted_at is NULL
					`

	err = db.QueryRow(query_paginate, search, category).Scan(&total_data)

	if err != nil {
		return err
	}

	response := models.Response{
		Data:     data_books,
		Message:  "Get books list successfully",
		Paginate: map[string]int{"total_data": total_data},
	}

	return ctx.JSON(http.StatusOK, response)
}

func Create(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	var book BookModels.CreateBook

	if err := ctx.Bind(&book); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&book); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM books WHERE code = $1 AND deleted_at IS NULL)"
	err := db.QueryRow(query, book.Code).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		response := models.ResponseDetail{
			Message: "This code has been used",
		}
		return ctx.JSON(http.StatusBadRequest, response)
	}

	thumbnail, err := ctx.FormFile("thumbnail")
	if err != nil {
		response := models.ResponseDetail{
			Message: "Thumbnail required",
		}
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Buka file gambar
	src, err := thumbnail.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	// Simpan gambar ke folder
	folder := "./assets/images/thumbnail"
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, os.ModePerm)
	}

	currentTime := time.Now()
	fileName := currentTime.Format("20060102_150405") + "_" + thumbnail.Filename

	// Buat path lengkap ke file
	filePath := filepath.Join(folder, fileName)

	fileName = "http://localhost:3333/thumbnail/" + fileName

	// Buat file di path yang telah ditentukan
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Salin data gambar dari src ke dst
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	query = "INSERT INTO books(title, category_id, publication_year, description, code, thumbnail) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, title, category_id, publication_year, description, code, thumbnail, created_at, updated_at"

	err = db.QueryRow(query, &book.Title, &book.CategoryId, &book.PublicationYear, &book.Description, &book.Code, fileName).Scan(&book.Id, &book.Title, &book.CategoryId, &book.PublicationYear, &book.Description, &book.Code, &book.Thumbnail, &book.CreatedAt, &book.UpdatedAt)

	if err != nil {
		return err
	}

	for _, authorId := range book.AuthorId {
		query := "INSERT INTO author_book(author_id, book_id) VALUES ($1, $2)"
		_, err := db.Exec(query, authorId, &book.Id)
		if err != nil {
			return err
		}
	}

	response := models.ResponseDetail{
		Data:    book,
		Message: "Create Book successfully",
	}

	return ctx.JSON(http.StatusCreated, response)
}

func Show(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	var book BookModels.GetBookDetail

	query := `
			SELECT 
				books.id, 
				books.title,
				categories.name as category,
				books.publication_year,
				books.description,
				books.code,
				books.thumbnail,
				books.created_at,
				books.updated_at 
			FROM 
				books
			JOIN 
				categories ON categories.id = books.category_id
			WHERE
				books.id = $1
			`

	err := db.QueryRow(query, id).Scan(&book.Id, &book.Title, &book.Category, &book.PublicationYear, &book.Description, &book.Code, &book.Thumbnail, &book.CreatedAt, &book.UpdatedAt)

	if err != nil {
		response := models.ResponseDetail{
			Message: "book not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	query_total_inventory := `
					SELECT 
						COUNT(id)
					FROM 
						inventories 
					WHERE
						status != 'scrap'
					AND 
						book_id = $1
					`

	err = db.QueryRow(query_total_inventory, id).Scan(&book.TotalInventory)

	if err != nil {
		return err
	}

	query_total_available := `
						SELECT 
							COUNT(id)
						FROM 
							inventories 
						WHERE
							status = 'available'
						AND 
							book_id = $1
						`

	err = db.QueryRow(query_total_available, id).Scan(&book.TotalAvailable)

	if err != nil {
		return err
	}

	query = `
		SELECT 
			author.id, 
			author.name
		FROM 
			author_book
		JOIN 
			author ON author.id = author_book.author_id
		WHERE
			author_book.book_id = $1
	`

	rows, err := db.Query(query, id)

	if err != nil {
		return err
	}

	var data_author []BookModels.Author
	for rows.Next() {
		var author BookModels.Author

		err := rows.Scan(&author.Id, &author.Name)

		if err != nil {
			return err
		}

		data_author = append(data_author, author)
	}

	book.Author = data_author

	response := models.ResponseDetail{
		Data:    book,
		Message: "Get book detail successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}

func Update(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	var book BookModels.UpdateBook

	if err := ctx.Bind(&book); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&book); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM books WHERE code = $1 AND id != $2 AND deleted_at IS NULL)"
	err := db.QueryRow(query, book.Code, id).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		response := models.ResponseDetail{
			Message: "This code has been used",
		}
		return ctx.JSON(http.StatusBadRequest, response)
	}

	thumbnail, err := ctx.FormFile("thumbnail")

	var Filename string
	if err != nil {
		query := "SELECT thumbnail from books WHERE id = $1"

		err := db.QueryRow(query, id).Scan(&Filename)

		if err != nil {
			return err
		}
	} else {

		// Ambil nama file gambar lama dari database
		var oldFilename string
		query := "SELECT thumbnail FROM books WHERE id = $1"
		err := db.QueryRow(query, id).Scan(&oldFilename)
		if err != nil {
			return err
		}

		// Hapus gambar lama jika ada
		if oldFilename != "" {
			oldFilePath := filepath.Join("./assets/images/thumbnail", filepath.Base(oldFilename))
			if err := os.Remove(oldFilePath); err != nil && !os.IsNotExist(err) {
				return err
			}
		}

		src, err := thumbnail.Open()
		if err != nil {
			return err
		}

		defer src.Close()

		// Simpan gambar ke folder
		folder := "./assets/images/thumbnail"
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			os.MkdirAll(folder, os.ModePerm)
		}

		currentTime := time.Now()
		file := currentTime.Format("20060102_150405") + "_" + thumbnail.Filename

		// Buat path lengkap ke file
		filePath := filepath.Join(folder, file)

		file = "http://localhost:3333/thumbnail/" + file

		Filename = file

		// Buat file di path yang telah ditentukan
		dst, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Salin data gambar dari src ke dst
		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}
	}

	query = "UPDATE books SET title = $1, category_id = $2 ,publication_year = $3, description = $4, code = $5, thumbnail = $6, updated_at = CURRENT_TIMESTAMP WHERE id = $7 RETURNING id, title, category_id, publication_year, description, code, thumbnail, created_at, updated_at"

	err = db.QueryRow(query, &book.Title, &book.CategoryId, &book.PublicationYear, &book.Description, &book.Code, Filename, id).Scan(&book.Id, &book.Title, &book.CategoryId, &book.PublicationYear, &book.Description, &book.Code, &book.Thumbnail, &book.CreatedAt, &book.UpdatedAt)

	if err != nil {
		return err
	}

	if len(book.AuthorId) > 0 {

		query := "DELETE FROM author_book WHERE book_id = $1"
		_, err := db.Exec(query, &book.Id)
		if err != nil {
			return err
		}

		for _, authorId := range book.AuthorId {
			query := "INSERT INTO author_book(author_id, book_id) VALUES ($1, $2)"
			_, err := db.Exec(query, authorId, &book.Id)
			if err != nil {
				return err
			}
		}
	}

	response := models.ResponseDetail{
		Data:    book,
		Message: "Update Book successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}

func Delete(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	var oldFilename string
	query := "SELECT thumbnail FROM books WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&oldFilename)
	if err != nil {
		return err
	}

	// Hapus gambar lama jika ada
	if oldFilename != "" {
		oldFilePath := filepath.Join("./assets/images/thumbnail", filepath.Base(oldFilename))
		if err := os.Remove(oldFilePath); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	query = "UPDATE books SET deleted_at = CURRENT_TIMESTAMP  WHERE id = $1"

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
			Message: "books not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	response := models.ResponseDetail{
		Message: "books deleted Successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}

func ShowInvetory(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

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
			WHERE
				books.id = $1
			`

	rows, err := db.Query(query, id)

	if err != nil {
		return err
	}

	var data_inventory []BookModels.GetInventory
	for rows.Next() {
		var inventory BookModels.GetInventory

		err := rows.Scan(&inventory.Id, &inventory.Book, &inventory.EntryTime, &inventory.ScrapTime, &inventory.Status, &inventory.CreatedAt, &inventory.UpdatedAt)

		if err != nil {
			return err
		}

		data_inventory = append(data_inventory, inventory)
	}

	var total_data int
	query_paginate := "SELECT COUNT(id) FROM inventories WHERE book_id = $1"

	err = db.QueryRow(query_paginate, id).Scan(&total_data)

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
