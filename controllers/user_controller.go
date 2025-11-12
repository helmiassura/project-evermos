package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"evermos-project/models"
	"evermos-project/utils"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) GetMyProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var user models.User
	if err := uc.DB.First(&user, userID).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusNotFound, false, "Failed to GET data", []string{"User not found"}, nil)
	}

	response := fiber.Map{
		"nama":          user.Nama,
		"no_telp":       user.NoTelp,
		"tanggal_Lahir": user.TanggalLahir,
		"tentang":       user.Tentang,
		"pekerjaan":     user.Pekerjaan,
		"email":         user.Email,
		"id_provinsi":   fiber.Map{"id": user.IdProvinsi, "name": ""},
		"id_kota":       fiber.Map{"id": user.IdKota, "province_id": user.IdProvinsi, "name": ""},
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to GET data", nil, response)
}

func (uc *UserController) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req struct {
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

	if err := c.BodyParser(&req); err != nil {
		return utils.RespondJSON(c, fiber.StatusBadRequest, false, "Failed to UPDATE data", []string{err.Error()}, nil)
	}

	var user models.User
	if err := uc.DB.First(&user, userID).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusNotFound, false, "Failed to UPDATE data", []string{"User not found"}, nil)
	}

	var tPtr *time.Time
	if req.TanggalLahir != "" {
		t, err := time.Parse("2006-01-02", req.TanggalLahir) // ubah layout jika format berbeda
		if err != nil {
			return utils.RespondJSON(c, fiber.StatusBadRequest, false, "Failed to UPDATE data", []string{err.Error()}, nil)
		}
		tPtr = &t
	}

	// Update fields
	if req.Nama != "" {
		user.Nama = req.Nama
	}
	if req.NoTelp != "" {
		user.NoTelp = req.NoTelp
	}
	user.TanggalLahir = tPtr
	if req.Pekerjaan != "" {
		user.Pekerjaan = req.Pekerjaan
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.IdProvinsi != "" {
		user.IdProvinsi = req.IdProvinsi
	}
	if req.IdKota != "" {
		user.IdKota = req.IdKota
	}
	if req.Tentang != "" {
		user.Tentang = req.Tentang
	}

	if req.KataSandi != "" {
		hashedPassword, err := utils.HashPassword(req.KataSandi)
		if err != nil {
			return utils.RespondJSON(c, fiber.StatusInternalServerError, false, "Failed to UPDATE data", []string{err.Error()}, nil)
		}
		user.KataSandi = hashedPassword
	}

	if err := uc.DB.Save(&user).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusInternalServerError, false, "Failed to UPDATE data", []string{err.Error()}, nil)
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to UPDATE data", nil, "")
}
