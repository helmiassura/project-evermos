package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"evermos-project/models"
	"evermos-project/utils"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

type RegisterRequest struct {
    Nama         string `json:"nama"`
    KataSandi    string `json:"kata_sandi"`
    NoTelp       string `json:"no_telp"`
    TanggalLahir string `json:"tanggal_Lahir"`
    JenisKelamin string `json:"jenis_kelamin"`
    Tentang      string `json:"tentang"`
    Pekerjaan    string `json:"pekerjaan"`
    Email        string `json:"email"`
    IdProvinsi   string `json:"id_provinsi"`
    IdKota       string `json:"id_kota"`
}

type LoginRequest struct {
	NoTelp    string `json:"no_telp"`
	KataSandi string `json:"kata_sandi"`
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.RespondJSON(c, fiber.StatusBadRequest, false, "Failed to POST data", []string{err.Error()}, nil)
	}

	fmt.Printf("DEBUG Register payload: %+v\n", req)

	hashedPassword, err := utils.HashPassword(req.KataSandi)
	if err != nil {
		return utils.RespondJSON(c, fiber.StatusInternalServerError, false, "Failed to POST data", []string{err.Error()}, nil)
	}

	var tanggalLahirPtr *time.Time
	if req.TanggalLahir != "" {
		t, err := time.Parse("2006-01-02", req.TanggalLahir)
		if err != nil {
			return utils.RespondJSON(c, fiber.StatusBadRequest, false, "Failed to POST data", []string{"Invalid tanggal_Lahir format, use YYYY-MM-DD"}, nil)
		}
		tanggalLahirPtr = &t
	}

	tx := ac.DB.Begin()

	user := models.User{
		Nama:         req.Nama,
		KataSandi:    hashedPassword,
		NoTelp:       req.NoTelp,
		TanggalLahir: tanggalLahirPtr,
		JenisKelamin: req.JenisKelamin,
		Tentang:      req.Tentang,
		Pekerjaan:    req.Pekerjaan,
		Email:        req.Email,
		IdProvinsi:   req.IdProvinsi,
		IdKota:       req.IdKota,
		IsAdmin:      false,
	}

	fmt.Printf("DEBUG will create user: %+v\n", user)

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return utils.RespondJSON(c, fiber.StatusBadRequest, false, "Failed to POST data", []string{fmt.Sprintf("Error 1062: %v", err)}, nil)
	}

	// Auto create toko
	namaToko := utils.GenerateSlug(req.Nama)
	toko := models.Toko{
		IDUser:   user.ID,
		NamaToko: namaToko,
		UrlFoto:  "",
	}

	if err := tx.Create(&toko).Error; err != nil {
		tx.Rollback()
		return utils.RespondJSON(c, fiber.StatusInternalServerError, false, "Failed to POST data", []string{err.Error()}, nil)
	}

	tx.Commit()

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to POST data", nil, "Register Succeed")
}


func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.RespondJSON(c, fiber.StatusBadRequest, false, "Failed to POST data", []string{err.Error()}, nil)
	}

	var user models.User
	if err := ac.DB.Where("no_telp = ?", req.NoTelp).First(&user).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusUnauthorized, false, "Failed to POST data", []string{"No Telp atau kata sandi salah"}, nil)
	}

	if !utils.CheckPasswordHash(req.KataSandi, user.KataSandi) {
		return utils.RespondJSON(c, fiber.StatusUnauthorized, false, "Failed to POST data", []string{"No Telp atau kata sandi salah"}, nil)
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return utils.RespondJSON(c, fiber.StatusInternalServerError, false, "Failed to POST data", []string{err.Error()}, nil)
	}

	response := fiber.Map{
		"nama":          user.Nama,
		"no_telp":       user.NoTelp,
		"tanggal_Lahir": user.TanggalLahir,
		"tentang":       user.Tentang,
		"pekerjaan":     user.Pekerjaan,
		"email":         user.Email,
		"id_provinsi": fiber.Map{
			"id":   user.IdProvinsi,
			"name": "",
		},
		"id_kota": fiber.Map{
			"id":          user.IdKota,
			"province_id": user.IdProvinsi,
			"name":        "",
		},
		"token": token,
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to POST data", nil, response)
}
