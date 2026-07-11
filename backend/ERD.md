# ERD Diagram - Travel Dinas System

```mermaid
erDiagram
    DEPARTMENT ||--o{ USER : "has"
    ROLE ||--o{ USER : "has"
    USER ||--o{ BUSINESS_TRIP : "submits"
    BUSINESS_TRIP ||--o{ REIMBURSEMENT : "has"

    DEPARTMENT {
        uint id PK
        string name
        string code
    }

    ROLE {
        uint id PK
        string name
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    USER {
        uint id PK
        string name
        string email
        string password
        string role
        uint role_id FK
        uint department_id FK
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    BUSINESS_TRIP {
        uint id PK
        uint user_id FK
        string description
        string destination
        datetime start_date
        datetime end_date
        string initiator
        string summary
        string nomor_surat
        string status
        string attachment_path
        string attachment_paths
        string notes
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    REIMBURSEMENT {
        uint id PK
        uint trip_id FK
        string title
        string description
        numeric amount
        string status
        string rejected_reason
        date transaction_date
        datetime reviewed_at
        datetime created_at
        datetime updated_at
        datetime deleted_at
    }

    BLACKLISTED_TOKEN {
        uint id PK
        string token
        datetime expires_at
    }
```

## Keterangan Hubungan Antar Tabel (Relasi)
1. **DEPARTMENT & USER**: Satu departemen dapat menaungi banyak user (`1 to many`). Kolom `department_id` di tabel `USER` bertindak sebagai Foreign Key ke tabel `DEPARTMENT`.
2. **ROLE & USER**: Satu role dapat dimiliki oleh banyak user (`1 to many`). Kolom `role_id` di tabel `USER` bertindak sebagai Foreign Key ke tabel `ROLE`.
3. **USER & BUSINESS_TRIP**: Satu user dapat mengajukan banyak perjalanan dinas (`1 to many`). Kolom `user_id` di tabel `BUSINESS_TRIP` bertindak sebagai Foreign Key ke tabel `USER`.
4. **BUSINESS_TRIP & REIMBURSEMENT**: Satu perjalanan dinas dapat memiliki banyak klaim reimbursement/pengeluaran (`1 to many`). Kolom `trip_id` di tabel `REIMBURSEMENT` bertindak sebagai Foreign Key ke tabel `BUSINESS_TRIP`.
5. **BLACKLISTED_TOKEN**: Tabel independen sistem yang digunakan untuk menyimpan daftar token JWT yang tidak valid setelah user melakukan logout (proses blacklist token) demi keamanan sesi.
