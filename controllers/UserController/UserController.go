package UserController

import (
	"io"
	"library/database"
	"library/models"
	"library/models/UserModels"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetAll(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	query := `SELECT 
				users.id,
				users.name,
				users.email,
				roles.name as role,
				users.avatar,
				profiles.identity_number,
				profiles.gender,
				profiles.birth_date,
				profiles.address,
				profiles.phone_number,
				users.created_at,
				users.updated_at
			FROM users
				JOIN roles ON roles.id = users.role_id
				JOIN profiles ON profiles.user_id = users.id
			ORDER BY users.id ASC`

	rows, err := db.Query(query)

	if err != nil {
		return err
	}

	var data_users []UserModels.GetUser
	for rows.Next() {
		var user UserModels.GetUser

		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.RoleId, &user.Avatar, &user.IdentityNumber, &user.Gender, &user.BirthDate, &user.Address, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			return err
		}

		data_users = append(data_users, user)
	}

	var total_data int
	query_paginate := "SELECT COUNT(id) FROM users"

	err = db.QueryRow(query_paginate).Scan(&total_data)

	if err != nil {
		return err
	}

	response := models.Response{
		Data:     data_users,
		Message:  "Get Users list successfully",
		Paginate: map[string]int{"total_data": total_data},
	}

	return ctx.JSON(http.StatusOK, response)
}

func Create(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	var user UserModels.CreateUser
	if err := ctx.Bind(&user); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&user); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)"
	err := db.QueryRow(query, &user.Email).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		response := models.ResponseDetail{
			Message: "This Email has been used",
		}
		return ctx.JSON(http.StatusBadRequest, response)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query = "INSERT INTO users (name,email,password) VALUES ($1,$2,$3) RETURNING id,name,email"

	err = db.QueryRow(query, &user.Name, &user.Email, string(hashedPassword)).Scan(&user.Id, &user.Name, &user.Email)

	if err != nil {
		return err
	}

	query = "INSERT INTO profiles (user_id) VALUES ($1)"

	_, err = db.Exec(query, user.Id)

	if err != nil {
		return err
	}

	response := models.ResponseDetail{
		Data:    user,
		Message: "Create Users successfully",
	}

	return ctx.JSON(http.StatusCreated, response)
}

func Show(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	var user UserModels.GetUser

	query := `SELECT 
				users.id,
				users.name,
				users.email,
				roles.name as role,
				users.avatar,
				profiles.identity_number,
				profiles.gender,
				profiles.birth_date,
				profiles.address,
				profiles.phone_number,
				users.created_at,
				users.updated_at
			FROM users
				JOIN roles ON roles.id = users.role_id
				JOIN profiles ON profiles.user_id = users.id
			WHERE users.id = $1`

	err := db.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Email, &user.RoleId, &user.Avatar, &user.IdentityNumber, &user.Gender, &user.BirthDate, &user.Address, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		response := models.ResponseDetail{
			Message: "User not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	response := models.ResponseDetail{
		Data:    user,
		Message: "Get user detail successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}

func Update(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	id := ctx.Param("id")

	avatar, err := ctx.FormFile("avatar")
	if err != nil {
		return err
	}

	// Buka file gambar
	src, err := avatar.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	// Simpan gambar ke folder
	folder := "./assets/images/avatar"
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.MkdirAll(folder, os.ModePerm)
	}

	currentTime := time.Now()
	fileName := currentTime.Format("20060102_150405") + "_" + avatar.Filename

	// Buat path lengkap ke file
	filePath := filepath.Join(folder, fileName)

	fileName = "http://localhost:3333/avatar/" + fileName

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

	var user UserModels.UpdateUser

	if err := ctx.Bind(&user); err != nil {
		return err
	}

	query := `UPDATE users SET avatar = $1,updated_at = CURRENT_TIMESTAMP WHERE id = $2 RETURNING id,name,email,role_id,avatar,created_at,updated_at`

	err = db.QueryRow(query, fileName, id).Scan(&user.Id, &user.Name, &user.Email, &user.RoleId, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		response := models.ResponseDetail{
			Message: "User not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	query = `UPDATE profiles SET identity_number = $1, gender = $2,birth_date = $3,address = $4,phone_number = $5 WHERE user_id = $6 RETURNING identity_number, gender, birth_date, address, phone_number`

	err = db.QueryRow(query, &user.IdentityNumber, &user.Gender, &user.BirthDate, &user.Address, &user.PhoneNumber, id).Scan(&user.IdentityNumber, &user.Gender, &user.BirthDate, &user.Address, &user.PhoneNumber)

	if err != nil {
		return err
	}

	response := models.ResponseDetail{
		Data:    user,
		Message: "Update user successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}
