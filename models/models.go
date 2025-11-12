package models

import "time"

// ======================== USER ==========================
type User struct {
	ID           uint       `gorm:"column:id;primaryKey;autoIncrement;type:int unsigned" json:"id"`
	Nama         string     `gorm:"column:nama;type:varchar(255);not null" json:"nama"`
	KataSandi    string     `gorm:"column:kata_sandi;type:varchar(255);not null" json:"-"`
	NoTelp       string     `gorm:"column:no_telp;type:varchar(255);uniqueIndex;not null" json:"no_telp"`
	TanggalLahir *time.Time `gorm:"column:tanggal_lahir;type:date" json:"tanggal_lahir"`
	JenisKelamin string     `gorm:"column:jenis_kelamin;type:varchar(255)" json:"jenis_kelamin"`
	Tentang      string     `gorm:"column:tentang;type:text" json:"tentang"`
	Pekerjaan    string     `gorm:"column:pekerjaan;type:varchar(255)" json:"pekerjaan"`
	Email        string     `gorm:"column:email;type:varchar(255);uniqueIndex;not null" json:"email"`
	IdProvinsi   string     `gorm:"column:id_provinsi;type:varchar(255)" json:"id_provinsi"`
	IdKota       string     `gorm:"column:id_kota;type:varchar(255)" json:"id_kota"`
	IsAdmin      bool       `gorm:"column:is_admin;type:boolean;default:false" json:"is_admin"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// relations
	Tokos  []Toko   `gorm:"foreignKey:IDUser;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Alamats []Alamat `gorm:"foreignKey:IDUser;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (User) TableName() string {
	return "user"
}

// ======================== TOKO ==========================
type Toko struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement;type:int unsigned" json:"id"`
	IDUser    uint      `gorm:"column:id_user;type:int unsigned;not null;index" json:"id_user"`
	NamaToko  string    `gorm:"column:nama_toko;type:varchar(255);not null" json:"nama_toko"`
	UrlFoto   string    `gorm:"column:url_foto;type:varchar(255)" json:"url_foto"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	User   User      `gorm:"foreignKey:IDUser;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Produk []Product `gorm:"foreignKey:IDToko;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (Toko) TableName() string {
	return "toko"
}

// ======================== ALAMAT ==========================
type Alamat struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement;type:int unsigned" json:"id"`
	IDUser       uint      `gorm:"column:id_user;type:int unsigned;not null;index" json:"id_user"`
	JudulAlamat  string    `gorm:"column:judul_alamat;type:varchar(255);not null" json:"judul_alamat"`
	NamaPenerima string    `gorm:"column:nama_penerima;type:varchar(255);not null" json:"nama_penerima"`
	NoTelp       string    `gorm:"column:no_telp;type:varchar(255);not null" json:"no_telp"`
	DetailAlamat string    `gorm:"column:detail_alamat;type:varchar(255);not null" json:"detail_alamat"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	User User  `gorm:"foreignKey:IDUser;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Trxs []Trx `gorm:"foreignKey:AlamatKirim;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (Alamat) TableName() string {
	return "alamat"
}

// ======================== CATEGORY ==========================
type Category struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement;type:int unsigned" json:"id"`
	NamaCategory string    `gorm:"column:nama_category;type:varchar(255);not null" json:"nama_category"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Produk []Product `gorm:"foreignKey:IDCategory;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (Category) TableName() string {
	return "category"
}

// ======================== PRODUCT (produk) ==========================
type Product struct {
	ID            uint          `gorm:"column:id;primaryKey;autoIncrement;type:int unsigned" json:"id"`
	NamaProduk    string        `gorm:"column:nama_produk;type:varchar(255);not null" json:"nama_produk"`
	Slug          string        `gorm:"column:slug;type:varchar(255);not null;uniqueIndex" json:"slug"`
	HargaReseller string        `gorm:"column:harga_reseller;type:varchar(255);not null" json:"harga_reseller"`
	HargaKonsumen string        `gorm:"column:harga_konsumen;type:varchar(255);not null" json:"harga_konsumen"`
	Stok          int           `gorm:"column:stok;type:int;not null" json:"stok"`
	Deskripsi     string        `gorm:"column:deskripsi;type:text" json:"deskripsi"`
	IDToko        uint          `gorm:"column:id_toko;type:int unsigned;not null;index" json:"id_toko"`
	IDCategory    uint          `gorm:"column:id_category;type:int unsigned;not null;index" json:"id_category"`
	CreatedAt     time.Time     `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time     `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Toko     Toko          `gorm:"foreignKey:IDToko;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"toko"`
	Category Category      `gorm:"foreignKey:IDCategory;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"category"`
	Photos   []FotoProduct `gorm:"foreignKey:IDProduct;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"photos"`
	Logs     []LogProduct  `gorm:"foreignKey:IDProduct;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (Product) TableName() string {
	return "produk"
}

// ======================== FOTO PRODUCT (foto_produk) ==========================
type FotoProduct struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement;type:int unsigned" json:"id"`
	IDProduct uint      `gorm:"column:id_product;type:int unsigned;not null;index" json:"product_id"`
	Url       string    `gorm:"column:url;type:varchar(255);not null" json:"url"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Product Product `gorm:"foreignKey:IDProduct;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (FotoProduct) TableName() string {
	return "foto_produk"
}

// ======================== LOG PRODUCT (log_produk) ==========================
type LogProduct struct {
	ID            uint      `gorm:"column:id;primaryKey;autoIncrement;type:int unsigned" json:"id"`
	IDProduct     uint      `gorm:"column:id_product;type:int unsigned;not null;index" json:"id_product"`
	NamaProduk    string    `gorm:"column:nama_produk;type:varchar(255);not null" json:"nama_produk"`
	Slug          string    `gorm:"column:slug;type:varchar(255);not null" json:"slug"`
	HargaReseller string    `gorm:"column:harga_reseller;type:varchar(255);not null" json:"harga_reseller"`
	HargaKonsumen string    `gorm:"column:harga_konsumen;type:varchar(255);not null" json:"harga_konsumen"`
	Deskripsi     string    `gorm:"column:deskripsi;type:text" json:"deskripsi"`
	IDToko        uint      `gorm:"column:id_toko;type:int unsigned;not null" json:"id_toko"`
	IDCategory    uint      `gorm:"column:id_category;type:int unsigned;not null" json:"id_category"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (LogProduct) TableName() string {
	return "log_produk"
}

// ======================== TRANSAKSI (trx) ==========================
type Trx struct {
	ID          uint        `gorm:"column:id;primaryKey;autoIncrement;type:int unsigned" json:"id"`
	IDUser      uint        `gorm:"column:id_user;type:int unsigned;not null;index" json:"id_user"`
	AlamatKirim uint        `gorm:"column:alamat_kirim;type:int unsigned;not null;index" json:"alamat_kirim_id"`
	HargaTotal  int         `gorm:"column:harga_total;type:int;not null" json:"harga_total"`
	KodeInvoice string      `gorm:"column:kode_invoice;type:varchar(255);not null" json:"kode_invoice"`
	MethodBayar string      `gorm:"column:method_bayar;type:varchar(255);not null" json:"method_bayar"`
	CreatedAt   time.Time   `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	User      User        `gorm:"foreignKey:IDUser;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Alamat    Alamat      `gorm:"foreignKey:AlamatKirim;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"alamat"`
	DetailTrx []DetailTrx `gorm:"foreignKey:IDTrx;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"detail_trx"`
}

func (Trx) TableName() string {
	return "trx"
}

// ======================== DETAIL TRANSAKSI (detail_trx) ==========================
type DetailTrx struct {
	ID          uint       `gorm:"column:id;primaryKey;autoIncrement;type:int unsigned" json:"id"`
	IDTrx       uint       `gorm:"column:id_trx;type:int unsigned;not null;index" json:"id_trx"`
	IDLogProduk uint       `gorm:"column:id_log_produk;type:int unsigned;not null;index" json:"id_log_product"`
	IDToko      uint       `gorm:"column:id_toko;type:int unsigned;not null;index" json:"id_toko"`
	Kuantitas   int        `gorm:"column:kuantitas;type:int;not null" json:"kuantitas"`
	HargaTotal  int        `gorm:"column:harga_total;type:int;not null" json:"harga_total"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Trx       Trx        `gorm:"foreignKey:IDTrx;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	LogProduct LogProduct `gorm:"foreignKey:IDLogProduk;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"product"`
	Toko      Toko       `gorm:"foreignKey:IDToko;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"toko"`
}

func (DetailTrx) TableName() string {
	return "detail_trx"
}

// ======================== DTO Untuk API eksternal ==========================
type Province struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type City struct {
	ID         string `json:"id"`
	ProvinceID string `json:"province_id"`
	Name       string `json:"name"`
}

type District struct {
	ID       string `json:"id"`
	RegencyID string `json:"regency_id"`
	Name     string `json:"name"`
}

type Village struct {
	ID         string `json:"id"`
	DistrictID string `json:"district_id"`
	Name       string `json:"name"`
}