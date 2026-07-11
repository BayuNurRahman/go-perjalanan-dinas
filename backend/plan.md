# Rencana Struktur Folder Proyek TA (Clean Architecture)

**Judul TA:** Perancangan Enterprise RESTful API Sistem Monitoring Perjalanan Dinas Berbasis Clean Architecture Golang dengan Validasi Dokumen dan Role-Based Access Control (RBAC)

---

## 1. Pohon Struktur Direktori

```text
go-perjalanan-dinas/
├── config/             # Folder Utama Konfigurasi Sistem
│   ├── database.go     # Inisialisasi & Koneksi PostgreSQL via GORM
│   └── jwt.go          # Pengaturan Secret Key & Durasi Token JWT
├── docs/               # Dokumentasi API Swagger (Hasil otomatis dari 'swag init')
├── middleware/         # Berkas Keamanan & Otorisasi (Interceptor)
│   ├── auth_jwt.go     # Memvalidasi token masuk
│   └── rbac.go         # Membatasi hak akses Peran (Role)
├── routes/             # Manajemen Rute Terpusat
│   └── routes.go       # Berkas tunggal wajib untuk registrasi semua endpoint API
├── uploads/            # Direktori penyimpanan berkas fisik PDF/Gambar di PC lokal
├── src/
│   ├── handler/        # HTTP Layer (Mengurusi Kontroler Gin Gonic & Request Binding)
│   │   ├── auth_handler.go
│   │   └── trip_handler.go
│   ├── service/        # Business Logic Layer (Aturan bisnis & validasi berkas)
│   │   ├── auth_service.go
│   │   └── trip_service.go
│   └── repository/     # Data Layer (Kueri database PostgreSQL menggunakan GORM)
│       ├── user_repository.go
│       └── trip_repository.go
├── models/             # GORM Struct Models / Cetak Biru Tabel Database
│   ├── user.go         # Skema tabel Users (ID, Nama, Email, Password, Role)
│   └── trip.go         # Skema tabel Perjalanan Dinas & Dokumen
├── dto/                # Data Transfer Object (Wrapper Request & Response JSON)
│   ├── auth_dto.go
│   └── trip_dto.go
├── main.go             # Entry Point Utama Aplikasi (Bootstrapping & Dependency Injection)
├── .env          # Spesifikasi Kontainerisasi Aplikasi Go
└── docker-compose.yml  # Orkestrasi Lokal Server Go + Database PostgreSQL