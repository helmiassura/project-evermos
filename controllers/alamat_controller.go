package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"evermos-project/models"
	"evermos-project/utils"
)

type AlamatController struct {
	DB *gorm.DB
}

func NewAlamatController(db *gorm.DB) *AlamatController {
	return &AlamatController{DB: db}
}

func (ac *AlamatController) GetMyAlamat(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var alamatList []models.Alamat
	query := ac.DB.Where("id_user = ?", userID)

	if search := c.Query("judul_alamat"); search != "" {
		query = query.Where("judul_alamat LIKE ?", "%"+search+"%")
	}

	if err := query.Find(&alamatList).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusInternalServerError, false, "Failed to GET data", []string{err.Error()}, nil)
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to GET data", nil, alamatList)
}

func (ac *AlamatController) GetAlamatByID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")

	var alamat models.Alamat
	if err := ac.DB.Where("id = ? AND id_user = ?", id, userID).First(&alamat).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusNotFound, false, "Failed to GET data", []string{"Alamat not found"}, nil)
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to GET data", nil, alamat)
}

func (ac *AlamatController) CreateAlamat(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req struct {
		JudulAlamat  string `json:"judul_alamat"`
		NamaPenerima string `json:"nama_penerima"`
		NoTelp       string `json:"no_telp"`
		DetailAlamat string `json:"detail_alamat"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.RespondJSON(c, fiber.StatusBadRequest, false, "Failed to POST data", []string{err.Error()}, nil)
	}

	alamat := models.Alamat{
		IDUser:       userID,
		JudulAlamat:  req.JudulAlamat,
		NamaPenerima: req.NamaPenerima,
		NoTelp:       req.NoTelp,
		DetailAlamat: req.DetailAlamat,
	}

	if err := ac.DB.Create(&alamat).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusInternalServerError, false, "Failed to POST data", []string{err.Error()}, nil)
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to POST data", nil, alamat.ID)
}

func (ac *AlamatController) UpdateAlamat(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")

	var alamat models.Alamat
	if err := ac.DB.Where("id = ? AND id_user = ?", id, userID).First(&alamat).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusBadRequest, false, "Failed to GET data", []string{"record not found"}, nil)
	}

	var req struct {
		NamaPenerima string `json:"nama_penerima"`
		NoTelp       string `json:"no_telp"`
		DetailAlamat string `json:"detail_alamat"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.RespondJSON(c, fiber.StatusBadRequest, false, "Failed to UPDATE data", []string{err.Error()}, nil)
	}

	alamat.NamaPenerima = req.NamaPenerima
	alamat.NoTelp = req.NoTelp
	alamat.DetailAlamat = req.DetailAlamat

	if err := ac.DB.Save(&alamat).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusInternalServerError, false, "Failed to UPDATE data", []string{err.Error()}, nil)
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to GET data", nil, "")
}

func (ac *AlamatController) DeleteAlamat(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")

	var alamat models.Alamat
	if err := ac.DB.Where("id = ? AND id_user = ?", id, userID).First(&alamat).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusBadRequest, false, "Failed to GET data", []string{"record not found"}, nil)
	}

	if err := ac.DB.Delete(&alamat).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusInternalServerError, false, "Failed to DELETE data", []string{err.Error()}, nil)
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to GET data", nil, "")
}
