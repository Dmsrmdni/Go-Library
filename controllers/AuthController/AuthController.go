package AuthController

import (
	"io"
	"library/database"
	"library/models"
	"library/models/AuthModels"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var SigningKey = []byte("secret")

func Login(ctx echo.Context) error {

	db := database.Init()

	defer db.Close()

	var login AuthModels.Login

	if err := ctx.Bind(&login); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&login); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	query := "SELECT id,password FROM users WHERE email = $1"

	var password string
	err := db.QueryRow(query, &login.Email).Scan(&login.Id, &password)
	if err != nil {
		response := models.ResponseDetail{
			Message: "Account not found",
		}

		return ctx.JSON(http.StatusNotFound, response)

	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(login.Password))

	if err != nil {
		return err
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = login.Id
	claims["email"] = login.Email
	exp := time.Now().Add(time.Hour * 24).Unix()
	claims["exp"] = exp

	t, err := token.SignedString(SigningKey)
	if err != nil {
		return err
	}

	response := models.ResponseDetail{
		Data: map[string]any{
			"access_token":         t,
			"expired_access_token": exp,
		},
		Message: "Login successfully",
	}
	return ctx.JSON(http.StatusOK, response)
}

func Register(ctx echo.Context) error {

	db := database.Init()

	defer db.Close()

	var register AuthModels.Register

	if err := ctx.Bind(&register); err != nil {
		return err
	}

	validate := validator.New()

	if err := validate.Struct(&register); err != nil {
		response := models.ResponseDetail{
			Data:    err.Error(),
			Message: "Validation Error",
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)"
	err := db.QueryRow(query, &register.Email).Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		response := models.ResponseDetail{
			Message: "This Email has been used",
		}
		return ctx.JSON(http.StatusBadRequest, response)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query = "INSERT INTO users (name,email,password) VALUES ($1,$2,$3) RETURNING id,name,email"

	err = db.QueryRow(query, &register.Name, &register.Email, string(hashedPassword)).Scan(&register.Id, &register.Name, &register.Email)

	if err != nil {
		return err
	}

	query = "INSERT INTO profiles (user_id) VALUES ($1)"

	_, err = db.Exec(query, register.Id)

	if err != nil {
		return err
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = register.Id
	claims["email"] = register.Email
	exp := time.Now().Add(time.Hour * 24).Unix()
	claims["exp"] = exp

	t, err := token.SignedString(SigningKey)
	if err != nil {
		return err
	}

	response := models.ResponseDetail{
		Data: map[string]any{
			"access_token":         t,
			"expired_access_token": exp,
		},
		Message: "Register successfully",
	}
	return ctx.JSON(http.StatusOK, response)
}

func Profile(ctx echo.Context) error {
	db := database.Init()
	defer db.Close()

	// Extract the token from the context
	user := ctx.Get("user")
	token := user.(*jwt.Token)

	// Get the claims from the token
	claims := token.Claims.(jwt.MapClaims)

	id := claims["id"]

	var profile AuthModels.GetProfile

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

	err := db.QueryRow(query, id).Scan(&profile.Id, &profile.Name, &profile.Email, &profile.RoleId, &profile.Avatar, &profile.IdentityNumber, &profile.Gender, &profile.BirthDate, &profile.Address, &profile.PhoneNumber, &profile.CreatedAt, &profile.UpdatedAt)

	if err != nil {
		response := models.ResponseDetail{
			Message: "User not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	response := models.ResponseDetail{
		Data:    profile,
		Message: "Get profil data successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}

func UpdateProfile(ctx echo.Context) error {
	db := database.Init()

	defer db.Close()

	user := ctx.Get("user")
	token := user.(*jwt.Token)

	// Get the claims from the token
	claims := token.Claims.(jwt.MapClaims)

	id := claims["id"]

	avatar, err := ctx.FormFile("avatar")
	var Filename *string

	if err != nil {
		query := "SELECT avatar from users WHERE id = $1"

		_ = db.QueryRow(query, id).Scan(&Filename)

	} else {
		var oldFilename string
		query := "SELECT avatar FROM users WHERE id = $1"
		err := db.QueryRow(query, id).Scan(&oldFilename)
		if err == nil {
			oldFilePath := filepath.Join("./assets/images/avatar", filepath.Base(oldFilename))
			if err := os.Remove(oldFilePath); err != nil && !os.IsNotExist(err) {
				return err
			}
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
		file := currentTime.Format("20060102_150405") + "_" + avatar.Filename

		// Buat path lengkap ke file
		filePath := filepath.Join(folder, file)

		file = "http://localhost:3333/avatar/" + file

		Filename = &file

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

	var profile AuthModels.UpdateProfile

	if err := ctx.Bind(&profile); err != nil {
		return err
	}

	query := `UPDATE users SET avatar = $1,updated_at = CURRENT_TIMESTAMP WHERE id = $2 RETURNING id,name,email,role_id,avatar,created_at,updated_at`

	err = db.QueryRow(query, Filename, id).Scan(&profile.Id, &profile.Name, &profile.Email, &profile.RoleId, &profile.Avatar, &profile.CreatedAt, &profile.UpdatedAt)

	if err != nil {
		response := models.ResponseDetail{
			Message: "User not found",
		}

		return ctx.JSON(http.StatusNotFound, response)
	}

	query = `UPDATE profiles SET identity_number = $1, gender = $2,birth_date = $3,address = $4,phone_number = $5, updated_at = CURRENT_TIMESTAMP WHERE user_id = $6 RETURNING identity_number, gender, birth_date, address, phone_number`

	err = db.QueryRow(query, &profile.IdentityNumber, &profile.Gender, &profile.BirthDate, &profile.Address, &profile.PhoneNumber, id).Scan(&profile.IdentityNumber, &profile.Gender, &profile.BirthDate, &profile.Address, &profile.PhoneNumber)

	if err != nil {
		return err
	}

	response := models.ResponseDetail{
		Data:    profile,
		Message: "Update profile successfully",
	}

	return ctx.JSON(http.StatusOK, response)
}
