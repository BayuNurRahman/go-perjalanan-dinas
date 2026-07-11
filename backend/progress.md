# Project Progress Timeline — up to 2026-07-11

## Overview
Semua fitur utama backend dan frontend telah diimplementasikan, diintegrasikan dengan database PostgreSQL dan objek storage MinIO, serta diverifikasi dengan pengujian unit/integrasi yang sukses.

---

## Progress Timeline

### 1) Project Initialization
- Menginisialisasi modul Go (`go.mod`) dan file utama `main.go`.
- Membuat struktur folder proyek (config, models, src, routes, migrations, docs).

### 2) Configuration and Environment
- Menambahkan konfigurasi database PostgreSQL, port aplikasi, JWT secret key, dan logger terpusat ke berkas `.env` dan `config/`.

### 3) Domain Modeling & Persistence
- Menyusun model data di folder `models/` (user, role, claim, trip, blacklisted_token, dll.).
- Mengimplementasikan layer repositori (`src/repository/`) untuk interaksi data via GORM.

### 4) Business Logic / Services
- Mengembangkan layer layanan (`src/service/`) yang mencakup manajemen user, autentikasi, perjalanan dinas, role, dan reimburse/klaim.

### 5) Middleware, Security & Access Control
- Membuat middleware autentikasi JWT (`AuthMiddleware`), CORS, dan RBAC (`RoleBlockMiddleware`).
- Menyelaraskan akses role `ADMIN_HR` (ke halaman Users/Register) dan `ADMIN_IT` (ke halaman Department/Roles).

### 6) Theme System (Light / Dark / System / Custom)
- Membangun `themeStore.js` (Pinia) di frontend untuk pelacakan tema secara persisten di `localStorage`.
- Mengimplementasikan modal pengaturan tema di `AppLayout.vue` dengan 4 opsi:
  - **Light Theme**: Mode terang dengan teks gelap.
  - **Dark Theme**: Mode gelap default.
  - **Sesuai Sistem**: Mengikuti preferensi mode sistem operasi.
  - **Custom Gradient**: Gradasi warna atas & bawah dengan *color picker*, input kode Hexa, dan 6 pilihan preset palet cepat dengan *live preview*.
- Menggunakan strategi CSS override terpusat di `style.css` berbasis variabel kustom (`[data-theme]`) untuk transisi warna yang halus tanpa merusak kode komponen yang sudah ada.

### 7) Claim Update & Delete Lifecycle
- Menyelesaikan fitur ubah (*update*) dan hapus (*delete*) klaim reimbursement.
- Mengintegrasikan pembersihan file lampiran fisik dari penyimpanan secara otomatis saat data klaim diperbarui atau dihapus.

### 8) MinIO S3 Object Storage Integration
- Mengintegrasikan penyimpanan berkas lampiran perjalanan dinas dan bukti klaim dengan server MinIO (port API `9000`).
- Mengotomatiskan pembuatan bucket `perjalanan-dinas` pada startup aplikasi.
- Merefaktor pengunggahan file di service `trip.go` dan `claim.go` untuk langsung mengalirkan data ke MinIO menggunakan method `PutObject` dari `minio-go/v7`.
- Menyediakan mekanisme **fallback otomatis ke disk lokal** agar semua unit/integration tests tetap dapat berjalan sukses tanpa ketergantungan wajib pada server MinIO eksternal yang aktif.
- Mengimplementasikan streaming pengunduhan file secara aman dari MinIO ke client melalui Gin `c.DataFromReader` pada handler `DownloadAttachment` dan `DownloadClaimAttachment`.

### 9) Trailing Slash Routing Refactor (Gin 307 CORS Fix)
- Menemukan dan mengatasi masalah browser memblokir upload file `multipart/form-data` karena redireksi Gin `307 Temporary Redirect` pada endpoint tanpa *trailing slash*.
- Merefaktor rute di `routes.go` untuk mendaftarkan endpoint utama dengan dan tanpa garis miring (seperti `/trips` dan `/trips/`), serta memperbarui helper API frontend agar memanggil endpoint berakhiran garis miring secara langsung.

### 10) Verification & Quality Assurance
- Menjalankan unit dan integration tests di backend secara berkala.
- ✅ Status pengujian: **LULUS (ok)** untuk semua package (`middleware`, `src/handler`, `src/service`).
