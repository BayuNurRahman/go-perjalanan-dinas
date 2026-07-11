# 📋 Rencana Setup Frontend - Sistem Perjalanan Dinas

## 🧱 Tech Stack
| Teknologi | Versi | Kegunaan |
|-----------|-------|----------|
| **Vue.js** | 3.x (Composition API) | Framework utama frontend |
| **Vite** | 5.x | Build tool & dev server |
| **Vue Router** | 4.x | Client-side routing & navigation guard |
| **Pinia** | 2.x | State management (token, user data) |
| **Axios** | 1.x | HTTP client ke backend API |
| **Tailwind CSS** | 3.x | Styling utility-first |

---

## 🗂️ Struktur Folder Target

```
frontend/
├── public/
│   └── favicon.ico
├── src/
│   ├── api/                    # Semua pemanggilan API axios per domain
│   │   ├── auth.js             # login, logout, register
│   │   ├── trips.js            # CRUD perjalanan dinas
│   │   ├── claims.js           # CRUD reimbursement
│   │   ├── departments.js      # CRUD departemen
│   │   ├── users.js            # CRUD user
│   │   └── roles.js            # CRUD role
│   ├── components/             # Komponen reusable
│   │   ├── layout/
│   │   │   ├── Navbar.vue
│   │   │   ├── Sidebar.vue
│   │   │   └── AppLayout.vue
│   │   ├── ui/
│   │   │   ├── BaseButton.vue
│   │   │   ├── BaseModal.vue
│   │   │   ├── BaseTable.vue
│   │   │   ├── BaseInput.vue
│   │   │   ├── BaseBadge.vue   # Badge status (PENDING, APPROVED, dll)
│   │   │   └── LoadingSpinner.vue
│   │   └── trip/
│   │       ├── TripCard.vue
│   │       ├── TripStatusBadge.vue
│   │       └── TripForm.vue
│   ├── views/                  # Halaman per role
│   │   ├── auth/
│   │   │   └── LoginView.vue
│   │   ├── employee/
│   │   │   ├── DashboardView.vue       # GET /trips/dashboard
│   │   │   ├── MyTripsView.vue         # GET /trips/me
│   │   │   ├── CreateTripView.vue      # POST /trips/
│   │   │   ├── TripDetailView.vue      # GET trip detail + klaim
│   │   │   └── SubmitClaimView.vue     # POST /claims/
│   │   ├── manager/
│   │   │   ├── DashboardView.vue       # GET /trips/manager/dashboard
│   │   │   ├── ApplicationsView.vue    # GET /trips/manager/applications
│   │   │   ├── TeamDistributionView.vue # GET /trips/manager/team-distribution
│   │   │   └── TripDetailView.vue      # PATCH /trips/:id/status
│   │   ├── finance/
│   │   │   ├── AllTripsView.vue        # GET /trips/
│   │   │   ├── TripFinancialView.vue   # PATCH /trips/:id/review-financial
│   │   │   ├── ClaimsView.vue          # GET /claims/trip/:trip_id
│   │   │   └── DisburseView.vue        # PATCH /trips/:id/disburse
│   │   └── superadmin/
│   │       ├── UsersView.vue           # GET/PUT/DELETE /users/
│   │       ├── DepartmentsView.vue     # CRUD /departments/
│   │       ├── RolesView.vue           # CRUD /roles/
│   │       └── RegisterUserView.vue    # POST /auth/register
│   ├── router/
│   │   └── index.js            # Route definitions + navigation guard
│   ├── stores/
│   │   └── auth.js             # Pinia store: token, user info, role
│   ├── utils/
│   │   └── axios.js            # Axios instance + interceptor JWT
│   ├── App.vue
│   └── main.js
├── index.html
├── vite.config.js
├── tailwind.config.js
├── postcss.config.js
└── package.json
```

---

## 🔐 Daftar Halaman & Role Akses

| Route | Komponen | Role yang Diizinkan |
|-------|----------|---------------------|
| `/login` | `LoginView.vue` | Public |
| `/dashboard` | `employee/DashboardView.vue` | EMPLOYEE, MANAGER, SUPER_ADMIN |
| `/trips` | `employee/MyTripsView.vue` | EMPLOYEE |
| `/trips/create` | `employee/CreateTripView.vue` | EMPLOYEE |
| `/trips/:id` | `employee/TripDetailView.vue` | EMPLOYEE |
| `/trips/:id/claim` | `employee/SubmitClaimView.vue` | EMPLOYEE |
| `/manager/dashboard` | `manager/DashboardView.vue` | MANAGER, SUPER_ADMIN |
| `/manager/applications` | `manager/ApplicationsView.vue` | MANAGER, SUPER_ADMIN |
| `/manager/team` | `manager/TeamDistributionView.vue` | MANAGER, SUPER_ADMIN |
| `/finance/trips` | `finance/AllTripsView.vue` | ADMIN_FIN |
| `/finance/claims` | `finance/ClaimsView.vue` | ADMIN_FIN |
| `/admin/users` | `superadmin/UsersView.vue` | SUPER_ADMIN |
| `/admin/departments` | `superadmin/DepartmentsView.vue` | SUPER_ADMIN |
| `/admin/roles` | `superadmin/RolesView.vue` | SUPER_ADMIN |
| `/admin/register` | `superadmin/RegisterUserView.vue` | SUPER_ADMIN |

---

## 🔄 Alur Autentikasi

```
User → Login (email + password)
     ↓
POST /api/v1/auth/login
     ↓
Response: { token, role, name }
     ↓
Simpan ke Pinia store + localStorage
     ↓
Axios interceptor otomatis sisipkan:
  Authorization: Bearer <token>
     ↓
Navigation Guard router cek role
     ↓
Redirect ke dashboard sesuai role
```

---

## 🗺️ Langkah Implementasi (Urutan)

### Phase 1 - Setup Awal
- [ ] Scaffold project dengan `npm create vite@latest . -- --template vue`
- [ ] Install dependencies: `vue-router`, `pinia`, `axios`, `tailwindcss`
- [ ] Konfigurasi Tailwind CSS
- [ ] Buat `src/utils/axios.js` (instance + interceptor JWT)
- [ ] Buat Pinia store `src/stores/auth.js`

### Phase 2 - Routing & Layout
- [ ] Buat file `src/router/index.js` dengan semua route + navigation guard
- [ ] Buat komponen layout: `AppLayout.vue`, `Navbar.vue`, `Sidebar.vue`
- [ ] Sidebar menampilkan menu yang berbeda sesuai role

### Phase 3 - Auth
- [ ] Buat `LoginView.vue` dengan form email + password
- [ ] Hubungkan ke `POST /api/v1/auth/login`
- [ ] Simpan token & redirect sesuai role setelah login
- [ ] Implementasi logout (call `POST /auth/logout`, hapus token)

### Phase 4 - Employee Pages
- [ ] `DashboardView.vue` → statistik perjalanan dinas pribadi
- [ ] `MyTripsView.vue` → daftar trip + pagination
- [ ] `CreateTripView.vue` → form multipart + upload file
- [ ] `TripDetailView.vue` → detail trip + list klaim
- [ ] `SubmitClaimView.vue` → form klaim + upload bukti

### Phase 5 - Manager Pages
- [ ] `manager/DashboardView.vue` → statistik tim
- [ ] `manager/ApplicationsView.vue` → daftar pengajuan + tombol approve/reject
- [ ] `manager/TeamDistributionView.vue` → distribusi perjalanan per anggota tim
- [ ] `manager/TripDetailView.vue` → form update status trip

### Phase 6 - Finance Pages
- [ ] `finance/AllTripsView.vue` → semua trip + filter
- [ ] `finance/TripFinancialView.vue` → review finansial trip
- [ ] `finance/ClaimsView.vue` → daftar klaim + approve/reject
- [ ] `finance/DisburseView.vue` → form pencairan dana

### Phase 7 - Super Admin Pages
- [ ] `superadmin/UsersView.vue` → tabel user + aksi edit/hapus
- [ ] `superadmin/DepartmentsView.vue` → CRUD departemen
- [ ] `superadmin/RolesView.vue` → CRUD role
- [ ] `superadmin/RegisterUserView.vue` → form daftarkan user baru

### Phase 8 - Polish & UX
- [ ] Loading state di setiap request API
- [ ] Error handling global (token expired → auto redirect login)
- [ ] Notifikasi toast (sukses/gagal)
- [ ] Responsive layout (mobile friendly)

---

## ⚙️ Konfigurasi API

```javascript
// src/utils/axios.js
Base URL : http://localhost:8080/api/v1
Header   : Authorization: Bearer <token>
```

---

## 🌐 Konfigurasi CORS Backend

Backend sudah mendukung CORS. Saat `APP_ENV=development`, semua origin diizinkan (`*`).
Frontend dev server Vite berjalan di `http://localhost:5173`.

---

## 📦 Perintah yang Akan Dijalankan

```bash
# 1. Scaffold project di folder frontend/
cd frontend
npm create vite@latest . -- --template vue

# 2. Install semua dependencies
npm install
npm install vue-router pinia axios
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p

# 3. Jalankan dev server
npm run dev
# Tersedia di: http://localhost:5173
```
