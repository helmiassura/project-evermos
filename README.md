# 🛍️ Evermos Mini Project - Backend Golang
<div align="center">

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)]()
[![Fiber](https://img.shields.io/badge/Fiber-Framework-2C8EBB?logo=fiber&logoColor=white)]()
[![MySQL](https://img.shields.io/badge/MySQL-Database-00758F?logo=mysql&logoColor=white)]()
[![JWT](https://img.shields.io/badge/Auth-JWT-000000?logo=jsonwebtokens&logoColor=white)]()

### RESTful API E-Commerce Backend  
Dibangun dengan **Golang**, **Fiber Framework**, **GORM ORM**, dan **MySQL**

---

🎓 **Proyek Akhir Virtual Internship**  
**Rakamin Academy x Evermos**  
_Backend Developer Program 2025_

👨‍💻 **Dibuat oleh:** M Helmi Assura

📘 [Tentang Proyek](#-tentang-proyek) 🧩 [API Documentation](#-dokumentasi-api) 🚀 [Quick Start](#-instalasi--setup)

</div>

---

## 📌 Tentang Proyek

**Evermos Mini Project** adalah sistem backend REST API untuk platform e-commerce sederhana  
yang mencakup fitur manajemen **user, toko, produk, kategori, alamat, dan transaksi**  
dengan sistem autentikasi **JWT** serta role-based access control.

Proyek ini dikembangkan mengikuti spesifikasi **Rakamin Evermos Postman Collection**,  
serta menerapkan **Clean Architecture** dan **best practices backend modern**  
menggunakan *modular design* berbasis *controller-service-model*.

---

## 🎯 Tujuan Pembelajaran

- Implementasi REST API dengan Golang  
- Penerapan Clean Architecture & Modular Design  
- Manajemen relasi database (One-to-Many, Many-to-One)  
- Implementasi JWT Authentication & Authorization  
- Penerapan file upload dan pagination  
- Penggunaan GORM sebagai ORM modern di Golang

---

## ⚙️ Teknologi yang Digunakan

| Komponen | Teknologi |
|-----------|------------|
| Bahasa | Go 1.21+ |
| Framework | Fiber v2 |
| ORM | GORM |
| Database | MySQL |
| Auth | JWT |
| Helper | godotenv, bcrypt |
| Dokumentasi | Postman |

---

## 🚀 Instalasi & Setup

### 1️⃣ Clone Repository
```bash
git clone https://github.com/helmiassura/evermos-project.git
cd evermos-project
```

### 2️⃣ Buat file `.env`
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=evermos_db
JWT_SECRET=your-secret-key-change-this-in-production
PORT=8000
```

### 3️⃣ Jalankan Server
```bash
go run main.go
```

Server akan berjalan di:
```
http://127.0.0.1:8000
```

Cek status API:
```
GET /
```
Response:
```json
{
  "status": true,
  "message": "API is running",
  "data": {
    "service": "Evermos API",
    "status": "ok"
  }
}
```

---

## 🧭 Dokumentasi API

### 🔐 Auth
| Method | Endpoint | Deskripsi |
|--------|-----------|-----------|
| POST | `/api/v1/auth/register` | Register user baru (otomatis buat toko) |
| POST | `/api/v1/auth/login` | Login dan mendapatkan JWT token |

### 👤 User
| Method | Endpoint | Deskripsi |
|--------|-----------|-----------|
| GET | `/api/v1/user/` | Get profil user login |
| PUT | `/api/v1/user/` | Update profil user login |

### 🏠 Alamat
| Method | Endpoint | Deskripsi |
|--------|-----------|-----------|
| GET | `/api/v1/user/alamat` | Get semua alamat user |
| GET | `/api/v1/user/alamat/:id` | Get alamat by ID |
| POST | `/api/v1/user/alamat` | Tambah alamat baru |
| PUT | `/api/v1/user/alamat/:id` | Update alamat user |
| DELETE | `/api/v1/user/alamat/:id` | Hapus alamat user |

### 🏪 Toko
| Method | Endpoint | Deskripsi |
|--------|-----------|-----------|
| GET | `/api/v1/toko` | Get semua toko (pagination + filter nama) |
| GET | `/api/v1/toko/:id_toko` | Get toko by ID |
| GET | `/api/v1/toko/my` | Get toko milik user login |
| PUT | `/api/v1/toko/:id_toko` | Update toko (upload foto) |

### 🛒 Produk
| Method | Endpoint | Deskripsi |
|--------|-----------|-----------|
| GET | `/api/v1/product` | Get semua produk (pagination + filter) |
| GET | `/api/v1/product/:id` | Get produk by ID |
| POST | `/api/v1/product` | Tambah produk (upload foto dengan key `photos`) |
| PUT | `/api/v1/product/:id` | Update produk |
| DELETE | `/api/v1/product/:id` | Hapus produk |

### 📂 Kategori (Admin only)
| Method | Endpoint | Deskripsi |
|--------|-----------|-----------|
| GET | `/api/v1/category` | List kategori |
| GET | `/api/v1/category/:id` | Detail kategori |
| POST | `/api/v1/category` | Tambah kategori (Admin only) |
| PUT | `/api/v1/category/:id` | Update kategori |
| DELETE | `/api/v1/category/:id` | Hapus kategori |

### 💳 Transaksi
| Method | Endpoint | Deskripsi |
|--------|-----------|-----------|
| GET | `/api/v1/trx` | Get semua transaksi milik user |
| GET | `/api/v1/trx/:id` | Get detail transaksi |
| POST | `/api/v1/trx` | Buat transaksi baru (isi log produk otomatis) |

### 🌍 Provinsi & Kota (Public API)
| Method | Endpoint | Deskripsi |
|--------|-----------|-----------|
| GET | `/api/v1/provcity/listprovinces` | List semua provinsi |
| GET | `/api/v1/provcity/detailprovince/:prov_id` | Detail provinsi berdasarkan ID |
| GET | `/api/v1/provcity/listcities/:prov_id` | List semua kota berdasarkan ID provinsi |
| GET | `/api/v1/provcity/detailcity/:city_id` | Detail kota berdasarkan ID |
| GET | `/api/v1/provcity/listdistricts/:city_id` | List semua kecamatan berdasarkan ID kota |
| GET | `/api/v1/provcity/detaildistrict/:district_id` | Detail kecamatan berdasarkan ID |
| GET | `/api/v1/provcity/listvillages/:district_id` | List semua desa/kelurahan berdasarkan ID kecamatan |
| GET | `/api/v1/provcity/detailvillage/:village_id` | Detail desa/kelurahan berdasarkan ID |

---

## 📂 Struktur Folder

```
evermos-project/
├── config/                             # Database configuration
│   └── database.go                     # GORM initialization & migration
├── controllers/                        # Business logic & request handlers
│   ├── auth_controller.go
│   ├── user_controller.go
│   ├── alamat_controller.go
│   ├── category_controller.go
│   ├── product_controller.go
│   ├── toko_controller.go
│   ├── trx_controller.go
│   └── provcity_controller.go
├── middleware/                         # Authentication & authorization
│   └── auth.go
├── models/                             # Data structures & ORM models
│   └── models.go
├── routes/                             # API route definitions
│   └── routes.go
├── utils/                              # Helper functions
│   ├── jwt.go
│   ├── password.go
│   ├── response.go
│   ├── slug.go
│   └── file.go
├── uploads/                            # Directory untuk uploaded files
├── main.go                             # Application entry point
├── go.mod                              # Go module definition
├── go.sum                              # Dependency checksum file
├── Rakamin Evermos.....collection.json # Postman API Collection
├── .env                                # Environment variables
└── README.md                           # Project documentation
```

---

## 🧱 Clean Architecture Principles

- **Controller** → menerima HTTP request & validasi input  
- **Service / Usecase** → menampung logika bisnis  
- **Model** → representasi data & ORM (GORM)  
- **Utils** → helper reusable seperti response JSON, file upload, JWT  
- **Middleware** → validasi token & role user  
- **Routes** → peta endpoint API

---

## 🧪 Testing

Gunakan Postman Collection:  
**Rakamin Evermos Virtual Internship.postman_collection.json**

Langkah:
1. Import collection ke Postman  
2. Jalankan urutan:
   - Register → Login → Copy token
   - Test endpoint `/user`, `/toko/my`, `/trx`
3. Cek respons dan status code

---

## 🧑‍💻 Author

**M Helmi Assura**  
Rakamin x Evermos Virtual Internship 2025  
📧 Email: lerkud600@gmail.com  
🔗 GitHub: [helmiassura](https://github.com/helmiassura)

---

> _“Clean code always looks like it was written by someone who cares.”_ 💡
