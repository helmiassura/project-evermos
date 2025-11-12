package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"evermos-project/controllers"
	"evermos-project/middleware"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {

	auth := controllers.NewAuthController(db)
	user := controllers.NewUserController(db)
	alamat := controllers.NewAlamatController(db)
	category := controllers.NewCategoryController(db)
	toko := controllers.NewTokoController(db)        
	product := controllers.NewProductController(db)  
	trx := controllers.NewTrxController(db)          
	provcity := controllers.NewProvCityController()

	// Base group
	api := app.Group("/api/v1")

	// ========================
	// AUTH (PUBLIC)
	// ========================
	api.Post("/auth/register", auth.Register)
	api.Post("/auth/login", auth.Login)

	// ========================
	// USER (PROTECTED)
	// ========================
	userGroup := api.Group("/user", middleware.AuthMiddleware(db))
	userGroup.Get("/", user.GetMyProfile)
	userGroup.Put("/", user.UpdateProfile)

	// ========================
	// ALAMAT (PROTECTED)
	// ========================
	userGroup.Get("/alamat", alamat.GetMyAlamat)
	userGroup.Get("/alamat/:id", alamat.GetAlamatByID)
	userGroup.Post("/alamat", alamat.CreateAlamat)
	userGroup.Put("/alamat/:id", alamat.UpdateAlamat)
	userGroup.Delete("/alamat/:id", alamat.DeleteAlamat)

	// ========================
	// CATEGORY
	// ========================

	// Public
	api.Get("/category", category.GetAllCategories)
	api.Get("/category/:id", category.GetCategoryByID)

	// Admin-only
	categoryAdmin := api.Group("/category", middleware.AuthMiddleware(db), middleware.AdminMiddleware(db))
	categoryAdmin.Post("/", category.CreateCategory)
	categoryAdmin.Put("/:id", category.UpdateCategory)
	categoryAdmin.Delete("/:id", category.DeleteCategory)

	// ========================
	// TOKO (PUBLIC + PROTECTED)
	// ========================

	// Gunakan single group, terapkan middleware per-route
	tokoGroup := api.Group("/toko")
	
	// Public routes
	tokoGroup.Get("/", toko.GetAllToko)                                              // GET /api/v1/toko
	
	// Protected routes - route spesifik HARUS didaftarkan SEBELUM route dinamis
	tokoGroup.Get("/my", middleware.AuthMiddleware(db), toko.GetMyToko)             // GET /api/v1/toko/my
	
	// Public route dengan parameter dinamis
	tokoGroup.Get("/:id_toko", toko.GetTokoByID)                                    // GET /api/v1/toko/:id_toko
	
	// Protected route dengan parameter dinamis
	tokoGroup.Put("/:id_toko", middleware.AuthMiddleware(db), toko.UpdateToko)      // PUT /api/v1/toko/:id_toko

	// ========================
	// PRODUCT (PUBLIC + PROTECTED)
	// ========================

	// Gunakan single group dengan middleware per-route
	productGroup := api.Group("/product")
	
	// Public routes
	productGroup.Get("/", product.GetAllProducts)                                        // GET /api/v1/product
	productGroup.Get("/:id", product.GetProductByID)                                     // GET /api/v1/product/:id
	
	// Protected routes
	productGroup.Post("/", middleware.AuthMiddleware(db), product.CreateProduct)         // POST /api/v1/product
	productGroup.Put("/:id", middleware.AuthMiddleware(db), product.UpdateProduct)       // PUT /api/v1/product/:id
	productGroup.Delete("/:id", middleware.AuthMiddleware(db), product.DeleteProduct)    // DELETE /api/v1/product/:id

	// ========================
	// TRANSAKSI (PROTECTED)
	// ========================

	trxGroup := api.Group("/trx", middleware.AuthMiddleware(db))
	trxGroup.Get("/", trx.GetAllTrx)
	trxGroup.Get("/:id", trx.GetTrxByID)
	trxGroup.Post("/", trx.CreateTrx)

	// ========================
	// PROV & CITY (PUBLIC)
	// ========================

	api.Get("/provcity/listprovinces", provcity.GetListProvinces)
	api.Get("/provcity/listcities/:prov_id", provcity.GetListCities)
	api.Get("/provcity/detailprovince/:prov_id", provcity.GetDetailProvince)
	api.Get("/provcity/detailcity/:city_id", provcity.GetDetailCity)
	api.Get("/provcity/listdistricts/:city_id", provcity.GetListDistricts)
	api.Get("/provcity/detaildistrict/:district_id", provcity.GetDetailDistrict)
	api.Get("/provcity/listvillages/:district_id", provcity.GetListVillages)
	api.Get("/provcity/detailvillage/:village_id", provcity.GetDetailVillage)
}