package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"evermos-project/models"
	"evermos-project/utils"
)

type CategoryController struct {
	DB *gorm.DB
}

func NewCategoryController(db *gorm.DB) *CategoryController {
	return &CategoryController{DB: db}
}

func (cc *CategoryController) GetAllCategories(c *fiber.Ctx) error {
	var categories []models.Category
	if err := cc.DB.Find(&categories).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusInternalServerError, false, "Failed to GET data", []string{err.Error()}, nil)
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to GET data", nil, categories)
}

func (cc *CategoryController) GetCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var category models.Category
	if err := cc.DB.First(&category, id).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusNotFound, false, "Failed to GET data", []string{"No Data Category"}, nil)
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to GET data", nil, category)
}

func (cc *CategoryController) CreateCategory(c *fiber.Ctx) error {
	var req struct {
		NamaCategory string `json:"nama_category"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.RespondJSON(c, fiber.StatusBadRequest, false, "Failed to POST data", []string{err.Error()}, nil)
	}

	category := models.Category{NamaCategory: req.NamaCategory}

	if err := cc.DB.Create(&category).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusInternalServerError, false, "Failed to POST data", []string{err.Error()}, nil)
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to POST data", nil, category.ID)
}

func (cc *CategoryController) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	var category models.Category
	if err := cc.DB.First(&category, id).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusNotFound, false, "Failed to UPDATE data", []string{"Category not found"}, nil)
	}

	var req struct {
		NamaCategory string `json:"nama_category"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.RespondJSON(c, fiber.StatusBadRequest, false, "Failed to UPDATE data", []string{err.Error()}, nil)
	}

	category.NamaCategory = req.NamaCategory

	if err := cc.DB.Save(&category).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusInternalServerError, false, "Failed to UPDATE data", []string{err.Error()}, nil)
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to UPDATE data", nil, "")
}

func (cc *CategoryController) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	var category models.Category
	if err := cc.DB.First(&category, id).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusBadRequest, false, "Failed to GET data", []string{"record not found"}, nil)
	}

	if err := cc.DB.Delete(&category).Error; err != nil {
		return utils.RespondJSON(c, fiber.StatusInternalServerError, false, "Failed to DELETE data", []string{err.Error()}, nil)
	}

	return utils.RespondJSON(c, fiber.StatusOK, true, "Succeed to DELETE data", nil, "")
}
