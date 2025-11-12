package controllers

import (
	"fmt"
	"strconv"
	"time"

	"evermos-project/models"
	"evermos-project/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TrxController struct {
	DB *gorm.DB
}

func NewTrxController(db *gorm.DB) *TrxController {
	return &TrxController{DB: db}
}

type CreateTrxRequest struct {
	MethodBayar string         `json:"method_bayar"`
	AlamatKirim uint           `json:"alamat_kirim"`
	DetailTrx   []DetailTrxReq `json:"detail_trx"`
}

type DetailTrxReq struct {
	ProductID uint `json:"product_id"`
	Kuantitas int  `json:"kuantitas"`
}

func (tc *TrxController) GetAllTrx(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)
	search := c.Query("search", "")

	offset := (page - 1) * limit

	query := tc.DB.Model(&models.Trx{}).
		Where("id_user = ?", userID).
		Preload("Alamat").
		Preload("DetailTrx.LogProduct").
		Preload("DetailTrx.LogProduct.Category").
		Preload("DetailTrx.LogProduct.Photos").
		Preload("DetailTrx.Toko")

	if search != "" {
		query = query.Where("kode_invoice LIKE ?", "%"+search+"%")
	}

	var trxList []models.Trx
	query.Limit(limit).Offset(offset).Find(&trxList)

	result := fiber.Map{
		"data":  trxList,
		"page":  page,
		"limit": limit,
	}

	return utils.FiberSuccess(c, "Succeed to GET data", result)
}

// GET TRANSACTION BY ID
func (tc *TrxController) GetTrxByID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id := c.Params("id")

	var trx models.Trx
	if err := tc.DB.
		Where("id = ? AND id_user = ?", id, userID).
		Preload("Alamat").
		Preload("DetailTrx.LogProduct").
		Preload("DetailTrx.Toko").
		First(&trx).Error; err != nil {
		return utils.FiberErrorCustom(c, fiber.StatusNotFound, "Failed to GET data", "No Data Trx")
	}

	return utils.FiberSuccess(c, "Succeed to GET data", trx)
}

// CREATE TRANSACTION
func (tc *TrxController) CreateTrx(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	// Parse JSON body
	var req CreateTrxRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.FiberError(c, "Failed to parse JSON", err)
	}

	// Validate alamat belongs to user
	var alamat models.Alamat
	if err := tc.DB.Where("id = ? AND id_user = ?", req.AlamatKirim, userID).
		First(&alamat).Error; err != nil {
		return utils.FiberErrorCustom(c, fiber.StatusBadRequest, "Failed to POST data", "Alamat not found")
	}

	tx := tc.DB.Begin()

	// Generate invoice code
	kodeInvoice := fmt.Sprintf("INV-%d", time.Now().Unix())

	// Create transaction
	trx := models.Trx{
		IDUser:      userID,
		AlamatKirim: req.AlamatKirim,
		HargaTotal:  0,
		KodeInvoice: kodeInvoice,
		MethodBayar: req.MethodBayar,
	}

	if err := tx.Create(&trx).Error; err != nil {
		tx.Rollback()
		return utils.FiberError(c, "Failed to POST data", err)
	}

	totalHarga := 0

	// Process each detail
	for _, detail := range req.DetailTrx {

		var product models.Product
		if err := tx.Preload("Category").First(&product, detail.ProductID).Error; err != nil {
			tx.Rollback()
			return utils.FiberErrorCustom(c, fiber.StatusBadRequest, "Failed to POST data", "Product not found")
		}

		// Check stock
		if product.Stok < detail.Kuantitas {
			tx.Rollback()
			return utils.FiberErrorCustom(c, fiber.StatusBadRequest, "Failed to POST data", "Insufficient stock")
		}

		// Create log product
		logProduct := models.LogProduct{
			IDProduct:     product.ID,
			NamaProduk:    product.NamaProduk,
			Slug:          product.Slug,
			HargaReseller: product.HargaReseller,
			HargaKonsumen: product.HargaKonsumen,
			Deskripsi:     product.Deskripsi,
			IDToko:        product.IDToko,
			IDCategory:    product.IDCategory,
		}

		if err := tx.Create(&logProduct).Error; err != nil {
			tx.Rollback()
			return utils.FiberError(c, "Failed to POST data", err)
		}

		// Calculate price: convert HargaKonsumen (string) -> int
		hargaKonsumenInt, err := strconv.Atoi(product.HargaKonsumen)
		if err != nil {
			tx.Rollback()
			return utils.FiberErrorCustom(c, fiber.StatusInternalServerError, "Failed to POST data", "Invalid product price")
		}
		hargaTotalItem := hargaKonsumenInt * detail.Kuantitas
		totalHarga += hargaTotalItem

		// Create detail transaction
		detailTrx := models.DetailTrx{
			IDTrx:       trx.ID,
			IDLogProduk: logProduct.ID,
			IDToko:      product.IDToko,
			Kuantitas:   detail.Kuantitas,
			HargaTotal:  hargaTotalItem,
		}

		if err := tx.Create(&detailTrx).Error; err != nil {
			tx.Rollback()
			return utils.FiberError(c, "Failed to POST data", err)
		}

		// Update product stock
		product.Stok -= detail.Kuantitas
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return utils.FiberError(c, "Failed to update product stock", err)
		}
	}

	// Update total price
	trx.HargaTotal = totalHarga
	if err := tx.Save(&trx).Error; err != nil {
		tx.Rollback()
		return utils.FiberError(c, "Failed to update trx", err)
	}

	tx.Commit()

	return utils.FiberSuccess(c, "Succeed to POST data", fiber.Map{"trx_id": trx.ID})
}
