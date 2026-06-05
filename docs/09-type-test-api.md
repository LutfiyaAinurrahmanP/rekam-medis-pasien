# Test Types API Documentation

## Overview

API untuk manajemen data test types (jenis tes laboratorium) dalam sistem rekam medis. Test Types adalah master data untuk berbagai jenis pemeriksaan laboratorium yang tersedia.

**Base URL:** `/api/v1/test-types`

---

## Table of Contents

- [Authentication](#authentication)
- [Authorization](#authorization)
- [Endpoints Summary](#endpoints-summary)
- [Public Endpoints](#public-endpoints)
- [Admin Endpoints](#admin-endpoints)
- [Super Admin Endpoints](#super-admin-endpoints)
- [Database Model](#database-model)
- [Request & Response Examples](#request--response-examples)
- [Error Responses](#error-responses)

---

## Authentication

Semua endpoints memerlukan JWT token di header:

```
Authorization: Bearer <your-jwt-token>
```

---

## Authorization

| Endpoint                           | Patient | Doctor | Receptionist | Admin | Super Admin |
| ---------------------------------- | ------- | ------ | ------------ | ----- | ----------- |
| GET /test-types                    | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /test-types/active             | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /test-types/inactive           | ❌      | ✅     | ✅           | ✅    | ✅          |
| GET /test-types/deleted            | ❌      | ❌     | ❌           | ✅    | ✅          |
| GET /test-types/:id                | ✅      | ✅     | ✅           | ✅    | ✅          |
| POST /test-types                   | ❌      | ❌     | ❌           | ✅    | ✅          |
| PUT /test-types/:id                | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /test-types/:id/activate     | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /test-types/:id/deactivate   | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /test-types/:id             | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /test-types/:id/restore      | ❌      | ❌     | ❌           | ✅    | ���         |
| DELETE /test-types/:id/hard-delete | ❌      | ❌     | ❌           | ❌    | ✅          |

---

## Endpoints Summary

### Public Endpoints (All Authenticated)

| Method | Endpoint                         | Description              | Role Required                            |
| ------ | -------------------------------- | ------------------------ | ---------------------------------------- |
| GET    | `/test-types`                    | List all test types      | All Authenticated                        |
| GET    | `/test-types/active`             | List active test types   | All Authenticated                        |
| GET    | `/test-types/inactive`           | List inactive test types | Doctor, Receptionist, Admin, Super Admin |
| GET    | `/test-types/deleted`            | List deleted test types  | Admin, Super Admin                       |
| GET    | `/test-types/:id`                | Get test type by ID      | All Authenticated                        |
| GET    | `/test-types/category/:category` | Get by category          | All Authenticated                        |

### Admin Endpoints

| Method | Endpoint                     | Description               | Role Required      |
| ------ | ---------------------------- | ------------------------- | ------------------ |
| POST   | `/test-types`                | Create test type          | Admin, Super Admin |
| PUT    | `/test-types/:id`            | Update test type          | Admin, Super Admin |
| PATCH  | `/test-types/:id/activate`   | Activate test type        | Admin, Super Admin |
| PATCH  | `/test-types/:id/deactivate` | Deactivate test type      | Admin, Super Admin |
| DELETE | `/test-types/:id`            | Soft delete test type     | Admin, Super Admin |
| PATCH  | `/test-types/:id/restore`    | Restore deleted test type | Admin, Super Admin |

### Super Admin Endpoints

| Method | Endpoint                      | Description                  | Role Required |
| ------ | ----------------------------- | ---------------------------- | ------------- |
| DELETE | `/test-types/:id/hard-delete` | Permanently delete test type | Super Admin   |

---

## Public Endpoints

### 1. List Test Types

**Endpoint:** `GET /api/v1/test-types`

**Description:** Mendapatkan daftar semua test types dengan pagination, search, dan filter.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter   | Type    | Default | Description                                          |
| ----------- | ------- | ------- | ---------------------------------------------------- |
| `page`      | integer | 1       | Halaman                                              |
| `page_size` | integer | 10      | Jumlah data per halaman (max: 100)                   |
| `search`    | string  | -       | Cari berdasarkan name, code, description             |
| `category`  | string  | -       | Filter by category                                   |
| `is_active` | boolean | -       | Filter by active status                              |
| `min_price` | decimal | -       | Filter minimum price                                 |
| `max_price` | decimal | -       | Filter maximum price                                 |
| `sort_by`   | string  | name    | Sort field (name, code, category, price, created_at) |
| `sort_dir`  | string  | asc     | Sort direction (asc, desc)                           |

**Example Request:**

```
GET /api/v1/test-types?page=1&page_size=10&category=Hematologi&is_active=true&sort_by=name&sort_dir=asc
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Test types retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "Complete Blood Count (CBC)",
        "code": "LAB-HEM-001",
        "category": "Hematologi",
        "description": "Pemeriksaan darah lengkap untuk menghitung jumlah sel darah",
        "price": 150000.0,
        "is_active": true,
        "created_at": "2024-01-19T10:00:00Z"
      },
      {
        "id": 2,
        "name": "Hemoglobin (Hb)",
        "code": "LAB-HEM-002",
        "category": "Hematologi",
        "description": "Pemeriksaan kadar hemoglobin dalam darah",
        "price": 50000.0,
        "is_active": true,
        "created_at": "2024-01-19T10:05:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total_items": 2,
      "total_pages": 1
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/test-types?page=1&page_size=10&category=Hematologi" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 2. List Active Test Types

**Endpoint:** `GET /api/v1/test-types/active`

**Description:** Mendapatkan daftar test types yang aktif saja (untuk order lab test).

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**
Same as List Test Types, but automatically filters `is_active=true`.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Active test types retrieved successfully",
  "data": {
    "total_active_tests": 50,
    "categories": [
      {
        "category": "Hematologi",
        "count": 15
      },
      {
        "category": "Kimia Darah",
        "count": 20
      },
      {
        "category": "Mikrobiologi",
        "count": 10
      },
      {
        "category": "Urinalisis",
        "count": 5
      }
    ],
    "data": [
      {
        "id": 1,
        "name": "Complete Blood Count (CBC)",
        "code": "LAB-HEM-001",
        "category": "Hematologi",
        "price": 150000.0,
        "is_active": true
      }
    ],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total_items": 50,
      "total_pages": 5
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/test-types/active?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Use Case:**

- Doctor ordering lab tests
- Patient melihat available tests
- Price list untuk billing

---

### 3. List Inactive Test Types

**Endpoint:** `GET /api/v1/test-types/inactive`

**Description:** Mendapatkan daftar test types yang tidak aktif (discontinued, obsolete).

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**
Same as List Test Types.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Inactive test types retrieved successfully",
  "data": {
    "data": [
      {
        "id": 100,
        "name": "Old Test Method",
        "code": "LAB-OLD-001",
        "category": "Obsolete",
        "description": "Test method yang sudah tidak digunakan",
        "is_active": false,
        "created_at": "2020-01-01T10:00:00Z",
        "updated_at": "2024-01-19T10:00:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total_items": 1,
      "total_pages": 1
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/test-types/inactive" \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN"
```

---

### 4. List Deleted Test Types

**Endpoint:** `GET /api/v1/test-types/deleted`

**Description:** Mendapatkan daftar test types yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**
Same as List Test Types.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Deleted test types retrieved successfully",
  "data": {
    "data": [
      {
        "id": 150,
        "name": "Deleted Test",
        "code": "LAB-DEL-001",
        "category": "Hematologi",
        "is_active": false,
        "created_at": "2023-01-01T10:00:00Z",
        "updated_at": "2024-01-19T10:00:00Z",
        "deleted_at": "2024-01-19T10:00:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total_items": 1,
      "total_pages": 1
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/test-types/deleted" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 5. Get Test Type by ID

**Endpoint:** `GET /api/v1/test-types/:id`

**Description:** Mendapatkan detail test type berdasarkan ID, termasuk statistik usage.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Test Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Test type retrieved successfully",
  "data": {
    "id": 1,
    "name": "Complete Blood Count (CBC)",
    "code": "LAB-HEM-001",
    "category": "Hematologi",
    "description": "Pemeriksaan darah lengkap untuk menghitung jumlah sel darah merah, sel darah putih, dan trombosit",
    "price": 150000.0,
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z",
    "statistics": {
      "total_orders": 1500,
      "orders_this_month": 150,
      "average_orders_per_day": 5,
      "revenue_this_month": 22500000.0
    },
    "components": [
      "Hemoglobin (Hb)",
      "Hematokrit (Ht)",
      "Eritrosit",
      "Leukosit",
      "Trombosit",
      "MCV, MCH, MCHC",
      "Hitung Jenis Leukosit"
    ],
    "sample_type": "Darah EDTA",
    "sample_volume": "2-3 ml",
    "preparation": "Tidak perlu puasa",
    "turnaround_time": "2-4 jam"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Test type not found",
  "error": "test type not found"
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/test-types/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Admin Endpoints

### 6. Create Test Type

**Endpoint:** `POST /api/v1/test-types`

**Description:** Admin membuat test type baru.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "name": "Complete Blood Count (CBC)",
  "code": "LAB-HEM-001",
  "category": "Hematologi",
  "description": "Pemeriksaan darah lengkap untuk menghitung jumlah sel darah merah, sel darah putih, dan trombosit",
  "price": 150000.0,
  "is_active": true
}
```

**Field Rules:**

- `name`: required, max 200 characters, indexed
- `code`: required, unique, max 50 characters, indexed
- `category`: optional, max 100 characters (Hematologi, Kimia Darah, Mikrobiologi, Urinalisis, dll)
- `description`: optional, text
- `price`: optional, decimal(10,2)
- `is_active`: optional, boolean (default: true)

**Response Success (201 Created):**

```json
{
  "success": true,
  "message": "Test type created successfully",
  "data": {
    "id": 1,
    "name": "Complete Blood Count (CBC)",
    "code": "LAB-HEM-001",
    "category": "Hematologi",
    "description": "Pemeriksaan darah lengkap untuk menghitung jumlah sel darah merah, sel darah putih, dan trombosit",
    "price": 150000.0,
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z"
  }
}
```

**Response Error (400 Bad Request):**

```json
{
  "success": false,
  "message": "Failed to create test type",
  "error": "code already exists"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/test-types \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Complete Blood Count (CBC)",
    "code": "LAB-HEM-001",
    "category": "Hematologi",
    "description": "Pemeriksaan darah lengkap",
    "price": 150000
  }'
```

**Notes:**

- Code harus unique dan mengikuti naming convention
- Category untuk grouping tests
- Price dalam IDR (Indonesian Rupiah)
- Format code: LAB-[CATEGORY_CODE]-[NUMBER]

---

### 7. Update Test Type

**Endpoint:** `PUT /api/v1/test-types/:id`

**Description:** Admin mengupdate data test type berdasarkan ID.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Test Type ID (integer)

**Request Body:**

```json
{
  "name": "Complete Blood Count (CBC) - Updated",
  "category": "Hematologi",
  "description": "Pemeriksaan darah lengkap dengan diferensial",
  "price": 175000.0,
  "is_active": true
}
```

**Field Rules:**

- All fields optional
- Validation same as create
- Code cannot be changed if already used in lab tests

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Test type updated successfully",
  "data": {
    "id": 1,
    "name": "Complete Blood Count (CBC) - Updated",
    "code": "LAB-HEM-001",
    "category": "Hematologi",
    "description": "Pemeriksaan darah lengkap dengan diferensial",
    "price": 175000.0,
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T15:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/test-types/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "price": 175000,
    "description": "Pemeriksaan darah lengkap dengan diferensial"
  }'
```

---

### 8. Activate Test Type

**Endpoint:** `PATCH /api/v1/test-types/:id/activate`

**Description:** Admin mengaktifkan test type yang sedang inactive.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Test Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Test type activated successfully",
  "data": null
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/test-types/1/activate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

**Use Case:**

- Test equipment kembali tersedia
- Test method updated dan ready to use
- New reagent arrived

---

### 9. Deactivate Test Type

**Endpoint:** `PATCH /api/v1/test-types/:id/deactivate`

**Description:** Admin menonaktifkan test type (equipment rusak, reagent habis, dll).

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Test Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Test type deactivated successfully",
  "data": null
}
```

**Notes:**

- Test type yang deactivated tidak bisa di-order
- Lab tests yang sudah di-order tetap bisa diproses
- Tidak muncul di available test list

**⚠️ Business Rules:**

- Pending lab tests dengan test type ini tetap valid
- Historical data tetap accessible
- Akan muncul notifikasi ke doctors

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/test-types/1/deactivate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 10. Soft Delete Test Type

**Endpoint:** `DELETE /api/v1/test-types/:id`

**Description:** Admin menghapus test type (soft delete). Data tetap ada dan bisa di-restore.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Test Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Test type deleted successfully",
  "data": null
}
```

**Notes:**

- Test type yang dihapus tidak muncul di list normal
- Lab test records tetap tersimpan
- Bisa di-restore dengan endpoint restore
- Automatically set is_active = false

**⚠️ Business Rules:**

- Historical lab tests tetap valid
- Tidak bisa di-order lagi
- Billing records tetap intact

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/test-types/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 11. Restore Test Type

**Endpoint:** `PATCH /api/v1/test-types/:id/restore`

**Description:** Admin me-restore test type yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Test Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Test type restored successfully",
  "data": null
}
```

**Notes:**

- Test type di-restore dengan status inactive
- Perlu activate manual jika ingin langsung available

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/test-types/1/restore \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

## Super Admin Endpoints

### 12. Hard Delete Test Type

**Endpoint:** `DELETE /api/v1/test-types/:id/hard-delete`

**Description:** Super Admin menghapus test type secara permanen dari database.

**Authentication:** Required (Super Admin Only)

**Request Headers:**

```
Authorization: Bearer <super-admin-token>
```

**URL Parameters:**

- `id`: Test Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Test type permanently deleted",
  "data": null
}
```

**⚠️ WARNING:**

- Data dihapus permanen dari database
- Tidak bisa di-restore
- Gunakan dengan sangat hati-hati

**⚠️ Business Rules:**

- Tidak bisa hard delete jika masih ada lab test records
- Must archive all historical data first
- Requires special approval

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/test-types/1/hard-delete \
  -H "Authorization: Bearer SUPER_ADMIN_JWT_TOKEN"
```

---

## Database Model

### Table: type_tests

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| name | VARCHAR(200) | NOT NULL, INDEX | Nama tes (e.g., Complete Blood Count, Urinalisis) |
| code | VARCHAR(50) | UNIQUE, NOT NULL, INDEX | Kode tes unik (e.g., CBC, UA, LFT) |
| category | VARCHAR(100) | INDEX | Kategori tes / reference ke master category |
| description | TEXT | NULLABLE | Deskripsi detail dan prosedur tes |
| price | DECIMAL(10,2) | DEFAULT 0 | Harga tes dalam IDR |
| is_active | BOOLEAN | NOT NULL, DEFAULT true, INDEX | Status aktif tes |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

**Indexes:**
- Primary Key: id
- Unique Index: code
- Regular Index: category, is_active, deleted_at

**Relationships:**
- Has Many Lab Tests (one-to-many)
- References Test Category (many-to-one, via category field)

**Notes:**
- Code sebaiknya uppercase dan singkat (CBC, UA, LFT, BIO, CRP, dll)
- Category bisa direct string atau foreign key ke test_categories table
- Price dapat disesuaikan per test dan dapat berubah
- is_active untuk hide discontinued tests
- Description harus include indikasi dan normal range values

---

## Error Responses

### 400 Bad Request

```json
{
  "success": false,
  "message": "Validation error",
  "error": "code already exists"
}
```

### 401 Unauthorized

```json
{
  "success": false,
  "message": "Authorization header is required",
  "error": null
}
```

### 403 Forbidden

```json
{
  "success": false,
  "message": "Access denied: insufficient permissions",
  "error": null
}
```

### 404 Not Found

```json
{
  "success": false,
  "message": "Test type not found",
  "error": "test type not found"
}
```

### 409 Conflict

```json
{
  "success": false,
  "message": "Cannot delete test type",
  "error": "test type has associated lab test records"
}
```

---

## Data Models

### Test Type Object (Full)

```json
{
  "id": 1,
  "name": "Complete Blood Count (CBC)",
  "code": "LAB-HEM-001",
  "category": "Hematologi",
  "description": "Pemeriksaan darah lengkap untuk menghitung jumlah sel darah merah, sel darah putih, dan trombosit",
  "price": 150000.0,
  "is_active": true,
  "created_at": "2024-01-19T10:00:00Z",
  "updated_at": "2024-01-19T10:00:00Z",
  "deleted_at": null
}
```

### Test Categories

```json
{
  "categories": [
    {
      "name": "Hematologi",
      "description": "Pemeriksaan darah dan komponen darah",
      "examples": ["CBC", "Hemoglobin", "Leukosit", "Trombosit"]
    },
    {
      "name": "Kimia Darah",
      "description": "Pemeriksaan kimia dalam darah",
      "examples": ["Glukosa", "Kolesterol", "Trigliserida", "SGOT/SGPT"]
    },
    {
      "name": "Mikrobiologi",
      "description": "Pemeriksaan mikroorganisme",
      "examples": ["Kultur Darah", "Gram Stain", "Antibiogram"]
    },
    {
      "name": "Urinalisis",
      "description": "Pemeriksaan urine",
      "examples": ["Urine Lengkap", "Protein Urine", "Glukosa Urine"]
    },
    {
      "name": "Serologi",
      "description": "Pemeriksaan antibodi dan antigen",
      "examples": ["HBsAg", "Anti-HIV", "VDRL", "Widal"]
    },
    {
      "name": "Imunoserologi",
      "description": "Pemeriksaan sistem imun",
      "examples": ["IgG", "IgM", "Complement"]
    },
    {
      "name": "Hormon",
      "description": "Pemeriksaan kadar hormon",
      "examples": ["TSH", "T3", "T4", "Insulin"]
    },
    {
      "name": "Tumor Marker",
      "description": "Penanda tumor",
      "examples": ["CEA", "AFP", "CA 125", "PSA"]
    }
  ]
}
```

---

## Business Rules

1. **Code Uniqueness**: Test code harus unik di seluruh sistem
2. **Naming Convention**: Format LAB-[CAT]-[NUM] (e.g., LAB-HEM-001)
3. **Active Status**: Hanya test active yang bisa di-order
4. **Price Management**: Price dalam IDR, optional (untuk internal tests)
5. **Category Grouping**: Category untuk organize tests
6. **Soft Delete Protection**: Historical data preserved
7. **Equipment Dependency**: Test availability tergantung equipment
8. **Reagent Management**: Deactivate jika reagent habis
9. **Quality Control**: Test harus QC approved sebelum activate
10. **Accreditation**: Follow lab accreditation standards

---

## Common Test Types (Indonesia)

### Hematologi

```json
[
  { "code": "LAB-HEM-001", "name": "Darah Lengkap (CBC)", "price": 150000 },
  { "code": "LAB-HEM-002", "name": "Hemoglobin (Hb)", "price": 50000 },
  { "code": "LAB-HEM-003", "name": "Hematokrit (Ht)", "price": 50000 },
  {
    "code": "LAB-HEM-004",
    "name": "Golongan Darah ABO/Rhesus",
    "price": 75000
  },
  { "code": "LAB-HEM-005", "name": "LED (Laju Endap Darah)", "price": 60000 },
  {
    "code": "LAB-HEM-006",
    "name": "Waktu Perdarahan/Pembekuan",
    "price": 80000
  }
]
```

### Kimia Darah

```json
[
  { "code": "LAB-KIM-001", "name": "Glukosa Darah Puasa", "price": 50000 },
  { "code": "LAB-KIM-002", "name": "Glukosa 2 Jam PP", "price": 50000 },
  { "code": "LAB-KIM-003", "name": "HbA1c", "price": 200000 },
  { "code": "LAB-KIM-004", "name": "Kolesterol Total", "price": 75000 },
  { "code": "LAB-KIM-005", "name": "Trigliserida", "price": 75000 },
  { "code": "LAB-KIM-006", "name": "HDL Kolesterol", "price": 80000 },
  { "code": "LAB-KIM-007", "name": "LDL Kolesterol", "price": 80000 },
  { "code": "LAB-KIM-008", "name": "SGOT (AST)", "price": 75000 },
  { "code": "LAB-KIM-009", "name": "SGPT (ALT)", "price": 75000 },
  { "code": "LAB-KIM-010", "name": "Ureum", "price": 70000 },
  { "code": "LAB-KIM-011", "name": "Kreatinin", "price": 70000 },
  { "code": "LAB-KIM-012", "name": "Asam Urat", "price": 65000 }
]
```

### Serologi

```json
[
  { "code": "LAB-SER-001", "name": "HBsAg", "price": 120000 },
  { "code": "LAB-SER-002", "name": "Anti-HCV", "price": 150000 },
  { "code": "LAB-SER-003", "name": "Anti-HIV", "price": 200000 },
  { "code": "LAB-SER-004", "name": "VDRL/RPR", "price": 100000 },
  { "code": "LAB-SER-005", "name": "Widal", "price": 120000 }
]
```

### Urinalisis

```json
[
  { "code": "LAB-URI-001", "name": "Urine Lengkap", "price": 80000 },
  { "code": "LAB-URI-002", "name": "Protein Urine Kualitatif", "price": 50000 },
  { "code": "LAB-URI-003", "name": "Glukosa Urine", "price": 45000 }
]
```

---

## Common Use Cases

### Use Case 1: Doctor Orders Lab Test

```bash
# 1. Browse available tests by category
GET /api/v1/test-types/category/Hematologi?is_active=true

# 2. View test details
GET /api/v1/test-types/1

# 3. Order lab test (different endpoint)
POST /api/v1/lab-tests
{
  "test_type_id": 1,
  "patient_id": 5,
  "medical_record_id": 10
}
```

### Use Case 2: Patient Views Available Tests

```bash
# 1. Search tests
GET /api/v1/test-types/search?keyword=darah&max_price=200000

# 2. View test details and price
GET /api/v1/test-types/1
```

### Use Case 3: Lab Admin Manages Tests

```bash
# 1. Create new test type
POST /api/v1/test-types

# 2. Update price
PUT /api/v1/test-types/1

# 3. Deactivate (reagent habis)
PATCH /api/v1/test-types/1/deactivate

# 4. Activate (reagent available)
PATCH /api/v1/test-types/1/activate
```

---

## Testing Examples

### Test 1: Complete Test Type Management Flow

```bash
# 1. Create Test Type
curl -X POST http://localhost:8080/api/v1/test-types \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hemoglobin (Hb)",
    "code": "LAB-HEM-002",
    "category": "Hematologi",
    "price": 50000
  }'

# 2. List Active Tests
curl -X GET "http://localhost:8080/api/v1/test-types/active" \
  -H "Authorization: Bearer PATIENT_TOKEN"

# 3. Search by Category
curl -X GET "http://localhost:8080/api/v1/test-types/category/Hematologi" \
  -H "Authorization: Bearer DOCTOR_TOKEN"

# 4. Update Price
curl -X PUT http://localhost:8080/api/v1/test-types/1 \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"price": 55000}'

# 5. Deactivate
curl -X PATCH http://localhost:8080/api/v1/test-types/1/deactivate \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

---

## Notes

- Code format: LAB-[CATEGORY]-[NUMBER]
- Price dalam IDR (Indonesian Rupiah)
- Category untuk grouping dan reporting
- Active status controls ordering availability
- Soft delete preserves historical data
- Test type linked to equipment and reagents
- Support for test packages/panels
- Integration with LIS (Lab Information System)
- Quality control requirements
- Turnaround time tracking

---

**Last Updated:** 2024-01-19  
**API Version:** 1.0.0
