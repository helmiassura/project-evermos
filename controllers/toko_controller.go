package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"evermos-project/middleware"
	"evermos-project/models"
	"evermos-project/utils"
)

type TokoController struct {
	DB *gorm.DB
}

func NewTokoController(db *gorm.DB) *TokoController {
	return &TokoController{DB: db}
}

// ===========================
// GET MY TOKO (debug + robust)
// ===========================
func (tc *TokoController) GetMyToko(c *fiber.Ctx) error {
	auth := string(c.Request().Header.Peek("Authorization"))
	fmt.Printf("[DEBUG] Authorization header: %s\n", auth)
	fmt.Printf("[DEBUG] c.Locals(\"userID\"): %#v\n", c.Locals("userID"))
	fmt.Printf("[DEBUG] c.Locals(\"user_id\"): %#v\n", c.Locals("user_id"))
	fmt.Printf("[DEBUG] c.Locals(\"id\"): %#v\n", c.Locals("id"))
	fmt.Printf("[DEBUG] c.Locals(middleware.UserIDKey): %#v\n", c.Locals(middleware.UserIDKey))

	var userIDUint uint

	if v := c.Locals(middleware.UserIDKey); v != nil {
		switch val := v.(type) {
		case uint:
			userIDUint = val
		case int:
			if val > 0 {
				userIDUint = uint(val)
			}
		case float64:
			if val > 0 {
				userIDUint = uint(val)
			}
		case string:
			if id, err := strconv.Atoi(val); err == nil && id > 0 {
				userIDUint = uint(id)
			}
		}
	}

	if userIDUint == 0 {
		if v := c.Locals("userID"); v != nil {
			switch val := v.(type) {
			case uint:
				userIDUint = val
			case int:
				if val > 0 {
					userIDUint = uint(val)
				}
			case float64:
				if val > 0 {
					userIDUint = uint(val)
				}
			case string:
				if id, err := strconv.Atoi(val); err == nil && id > 0 {
					userIDUint = uint(id)
				}
			}
		}
	}
	if userIDUint == 0 {
		if v := c.Locals("user_id"); v != nil {
			switch val := v.(type) {
			case uint:
				userIDUint = val
			case int:
				if val > 0 {
					userIDUint = uint(val)
				}
			case float64:
				if val > 0 {
					userIDUint = uint(val)
				}
			case string:
				if id, err := strconv.Atoi(val); err == nil && id > 0 {
					userIDUint = uint(id)
				}
			}
		}
	}
	if userIDUint == 0 {
		if v := c.Locals("id"); v != nil {
			switch val := v.(type) {
			case uint:
				userIDUint = val
			case int:
				if val > 0 {
					userIDUint = uint(val)
				}
			case float64:
				if val > 0 {
					userIDUint = uint(val)
				}
			case string:
				if id, err := strconv.Atoi(val); err == nil && id > 0 {
					userIDUint = uint(id)
				}
			}
		}
	}

	fmt.Printf("[DEBUG] resolved userIDUint = %d\n", userIDUint)

	if userIDUint == 0 {
		return utils.RespondJSON(c, fiber.StatusUnauthorized, false, "Failed to GET data", []string{"Unauthorized or userID not found in context"}, nil)
	}

	var toko models.Toko
	if err := tc.DB.Where("id_user = ?", userIDUint).First(&toko).Error; err != nil {
		fmt.Printf("[DEBUG] db query error: %v\n", err)
		return utils.RespondJSON(c, fiber.StatusNotFound, false, "Failed to GET data", []string{"Toko tidak ditemukan"}, nil)
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to GET data", nil, toko)
}

// ===========================
// UPDATE TOKO (type-safe userID)
// ===========================
func (tc *TokoController) UpdateToko(c *fiber.Ctx) error {
	rawID := c.Locals(middleware.UserIDKey)
	var userID uint
	switch v := rawID.(type) {
	case uint:
		userID = v
	case int:
		if v > 0 {
			userID = uint(v)
		}
	case float64:
		if v > 0 {
			userID = uint(v)
		}
	case string:
		if id, err := strconv.Atoi(v); err == nil && id > 0 {
			userID = uint(id)
		}
	default:
		userID = 0
	}
	if userID == 0 {
		if v := c.Locals("userID"); v != nil {
			switch val := v.(type) {
			case uint:
				userID = val
			case int:
				if val > 0 {
					userID = uint(val)
				}
			case float64:
				if val > 0 {
					userID = uint(val)
				}
			case string:
				if id, err := strconv.Atoi(val); err == nil && id > 0 {
					userID = uint(id)
				}
			}
		}
	}
	if userID == 0 {
		return utils.RespondJSON(c, fiber.StatusUnauthorized, false, "Unauthorized", []string{"userID not found"}, nil)
	}

	idToko, _ := strconv.Atoi(c.Params("id_toko"))

	var toko models.Toko
	if err := tc.DB.Where("id = ? AND id_user = ?", idToko, userID).First(&toko).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusNotFound, false, "Failed to UPDATE data", []string{"Toko not found or unauthorized"}, nil)
	}

	form, err := c.MultipartForm()
	if err != nil && err != fiber.ErrUnprocessableEntity {
		return utils.RespondJSON(c, fiber.StatusBadRequest, false, "Invalid form data", []string{err.Error()}, nil)
	}

	if nama := c.FormValue("nama_toko"); nama != "" {
		toko.NamaToko = nama
	}

	if form != nil {
		files := form.File["photo"]
		if len(files) > 0 {
			file := files[0]

			utils.DeleteFile(toko.UrlFoto)

			filename, err := utils.SaveFiberFile(file)
			if err != nil {
				return utils.RespondJSON(c, fiber.StatusInternalServerError, false, "Failed to upload file", []string{err.Error()}, nil)
			}
			toko.UrlFoto = filename
		}
	}

	if err := tc.DB.Save(&toko).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusInternalServerError, false, "Failed to UPDATE data", []string{err.Error()}, nil)
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to UPDATE data", nil, toko)
}

// ===========================
// GET TOKO BY ID
// ===========================
func (tc *TokoController) GetTokoByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id_toko"))

	var toko models.Toko
	if err := tc.DB.First(&toko, id).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusNotFound, false, "Failed to GET data", []string{"Toko tidak ditemukan"}, nil)
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to GET data", nil, toko)
}

// ===========================
// GET ALL TOKO
// ===========================
func (tc *TokoController) GetAllToko(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	nama := c.Query("nama")

	offset := (page - 1) * limit

	query := tc.DB.Model(&models.Toko{})
	if nama != "" {
		query = query.Where("nama_toko LIKE ?", "%"+nama+"%")
	}

	var tokoList []models.Toko
	query.Limit(limit).Offset(offset).Find(&tokoList)

	result := fiber.Map{
		"page":  page,
		"limit": limit,
		"data":  tokoList,
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to GET data", nil, result)
}
