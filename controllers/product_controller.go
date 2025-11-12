package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"evermos-project/middleware"
	"evermos-project/models"
	"evermos-project/utils"
)

type ProductController struct {
	DB *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{DB: db}
}

func (pc *ProductController) GetAllProducts(ctx *fiber.Ctx) error {
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))

	namaProduk := ctx.Query("nama_produk")
	categoryID, _ := strconv.Atoi(ctx.Query("category_id"))
	tokoID, _ := strconv.Atoi(ctx.Query("toko_id"))
	maxHarga, _ := strconv.Atoi(ctx.Query("max_harga"))
	minHarga, _ := strconv.Atoi(ctx.Query("min_harga"))

	offset := (page - 1) * limit

	query := pc.DB.Model(&models.Product{}).
		Preload("Toko").
		Preload("Category").
		Preload("Photos")

	if namaProduk != "" {
		query = query.Where("nama_produk LIKE ?", "%"+namaProduk+"%")
	}
	if categoryID > 0 {
		query = query.Where("id_category = ?", categoryID)
	}
	if tokoID > 0 {
		query = query.Where("id_toko = ?", tokoID)
	}
	if maxHarga > 0 {
		query = query.Where("CAST(harga_konsumen AS UNSIGNED) <= ?", maxHarga)
	}
	if minHarga > 0 {
		query = query.Where("CAST(harga_konsumen AS UNSIGNED) >= ?", minHarga)
	}

	var products []models.Product
	query.Limit(limit).Offset(offset).Find(&products)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Succeed to GET data",
		"data": fiber.Map{
			"data":  products,
			"page":  page,
			"limit": limit,
		},
	})
}

func (pc *ProductController) GetProductByID(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	var product models.Product
	if err := pc.DB.Preload("Toko").Preload("Category").Preload("Photos").First(&product, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Failed to GET data",
			"errors":  []string{"No Data Product"},
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Succeed to GET data",
		"data":    product,
	})
}

func (pc *ProductController) CreateProduct(ctx *fiber.Ctx) error {
	userID := ctx.Locals(middleware.UserIDKey).(uint)

	var toko models.Toko
	if err := pc.DB.Where("id_user = ?", userID).First(&toko).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Toko not found",
		})
	}

	namaProduk := ctx.FormValue("nama_produk")
	categoryID, _ := strconv.Atoi(ctx.FormValue("category_id"))
	hargaReseller := ctx.FormValue("harga_reseller")
	hargaKonsumen := ctx.FormValue("harga_konsumen")
	stok, _ := strconv.Atoi(ctx.FormValue("stok"))
	deskripsi := ctx.FormValue("deskripsi")

	hargaResellerInt := 0
	if hargaReseller != "" {
		v, err := strconv.Atoi(hargaReseller)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid harga_reseller",
				"errors":  []string{err.Error()},
			})
		}
		hargaResellerInt = v
	}

	hargaKonsumenInt := 0
	if hargaKonsumen != "" {
		v, err := strconv.Atoi(hargaKonsumen)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid harga_konsumen",
				"errors":  []string{err.Error()},
			})
		}
		hargaKonsumenInt = v
	}

	product := models.Product{
		NamaProduk:    namaProduk,
		Slug:          utils.GenerateSlug(namaProduk),
		HargaReseller: strconv.Itoa(hargaResellerInt), 
		HargaKonsumen: strconv.Itoa(hargaKonsumenInt), 
		Stok:          stok,
		Deskripsi:     deskripsi,
		IDToko:        toko.ID,
		IDCategory:    uint(categoryID),
	}

	tx := pc.DB.Begin()

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"errors":  []string{err.Error()},
		})
	}

	form, err := ctx.MultipartForm()
	if err == nil {
		files := form.File["photos"]

		for _, fileHeader := range files {
			filename, err := utils.SaveFiberFile(fileHeader)
			if err != nil {
				tx.Rollback()
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"errors":  []string{err.Error()},
				})
			}

			photo := models.FotoProduct{
				IDProduct: product.ID,
				Url:       filename,
			}
			tx.Create(&photo)
		}
	}

	tx.Commit()

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Succeed to POST data",
		"data":    product.ID,
	})
}

func (pc *ProductController) UpdateProduct(ctx *fiber.Ctx) error {
	userID := ctx.Locals(middleware.UserIDKey).(uint)
	id, _ := strconv.Atoi(ctx.Params("id"))

	var toko models.Toko
	if err := pc.DB.Where("id_user = ?", userID).First(&toko).Error; err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	var product models.Product
	if err := pc.DB.Where("id = ? AND id_toko = ?", id, toko.ID).First(&product).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Product not found or unauthorized",
		})
	}

	if namaProduk := ctx.FormValue("nama_produk"); namaProduk != "" {
		product.NamaProduk = namaProduk
		product.Slug = utils.GenerateSlug(namaProduk)
	}
	if categoryID := ctx.FormValue("category_id"); categoryID != "" {
		catID, _ := strconv.Atoi(categoryID)
		product.IDCategory = uint(catID)
	}
	if hargaReseller := ctx.FormValue("harga_reseller"); hargaReseller != "" {
		v, err := strconv.Atoi(hargaReseller)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid harga_reseller",
				"errors":  []string{err.Error()},
			})
		}
		product.HargaReseller = strconv.Itoa(v) 
	}
	if hargaKonsumen := ctx.FormValue("harga_konsumen"); hargaKonsumen != "" {
		v, err := strconv.Atoi(hargaKonsumen)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid harga_konsumen",
				"errors":  []string{err.Error()},
			})
		}
		product.HargaKonsumen = strconv.Itoa(v) 
	}
	if stok := ctx.FormValue("stok"); stok != "" {
		stokInt, _ := strconv.Atoi(stok)
		product.Stok = stokInt
	}
	if deskripsi := ctx.FormValue("deskripsi"); deskripsi != "" {
		product.Deskripsi = deskripsi
	}

	tx := pc.DB.Begin()

	form, err := ctx.MultipartForm()
	if err == nil {
		files := form.File["photos"]

		if len(files) > 0 {
			var oldPhotos []models.FotoProduct
			tx.Where("id_product = ?", product.ID).Find(&oldPhotos)
			for _, photo := range oldPhotos {
				utils.DeleteFile(photo.Url)
			}
			tx.Where("id_product = ?", product.ID).Delete(&models.FotoProduct{})

			for _, fileHeader := range files {
				filename, err := utils.SaveFiberFile(fileHeader)
				if err != nil {
					tx.Rollback()
					return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"success": false,
						"errors":  []string{err.Error()},
					})
				}

				photo := models.FotoProduct{
					IDProduct: product.ID,
					Url:       filename,
				}
				tx.Create(&photo)
			}
		}
	}

	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"errors":  []string{err.Error()},
		})
	}

	tx.Commit()

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Succeed to UPDATE data",
	})
}

func (pc *ProductController) DeleteProduct(ctx *fiber.Ctx) error {
	userID := ctx.Locals(middleware.UserIDKey).(uint)
	id, _ := strconv.Atoi(ctx.Params("id"))

	var toko models.Toko
	if err := pc.DB.Where("id_user = ?", userID).First(&toko).Error; err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
		})
	}

	var product models.Product
	if err := pc.DB.Preload("Photos").Where("id = ? AND id_toko = ?", id, toko.ID).First(&product).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "record not found",
		})
	}

	for _, photo := range product.Photos {
		utils.DeleteFile(photo.Url)
	}

	if err := pc.DB.Delete(&product).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"errors":  []string{err.Error()},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Succeed to DELETE data",
	})
}
