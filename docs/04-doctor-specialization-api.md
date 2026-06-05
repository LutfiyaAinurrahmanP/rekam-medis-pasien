\*\*\*\*# Doctor Specializations API Documentation

## Overview

API untuk manajemen master data spesialisasi dokter dalam sistem rekam medis. Specializations adalah master data yang dapat direferensikan oleh Doctors untuk menstandarisasi spesialisasi dokter di seluruh sistem (menghindari duplikasi data seperti "Jantung" vs "jantung" atau "Kardiologi").

**Base URL:** `/api/v1/doctor-specializations`

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

| Endpoint                                       | Patient | Doctor | Receptionist | Admin | Super Admin |
| ---------------------------------------------- | ------- | ------ | ------------ | ----- | ----------- |
| GET /doctor-specializations                    | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /doctor-specializations/active             | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /doctor-specializations/inactive      | ❌      | ❌     | ❌           | ✅    | ✅          |
| GET /doctor-specializations/:id                | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /doctor-specializations/deleted            | ❌      | ❌     | ❌           | ✅    | ✅          |
| POST /doctor-specializations                   | ❌      | ❌     | ❌           | ✅    | ✅          |
| PUT /doctor-specializations/:id                | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /doctor-specializations/:id/activate     | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /doctor-specializations/:id/deactivate   | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /doctor-specializations/:id             | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /doctor-specializations/:id/restore      | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /doctor-specializations/:id/hard-delete | ❌      | ❌     | ❌           | ❌    | ✅          |

---

## Endpoints Summary

### Public Endpoints (All Authenticated)

| Method | Endpoint                         | Description                        | Role Required     |
| ------ | -------------------------------- | ---------------------------------- | ----------------- |
| GET    | `/doctor-specializations`        | List all doctor specializations    | All Authenticated |
| GET    | `/doctor-specializations/active` | List active doctor specializations | All Authenticated |
| GET    | `/doctor-specializations/:id`    | Get doctor specialization by ID    | All Authenticated |

### Admin Endpoints

### X. List Inactive Doctor Specializations

**Endpoint:** `GET /api/v1/doctor-specializations/inactive`

**Description:** Mendapatkan daftar doctor specializations yang tidak aktif.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**Query Parameters:**

| Parameter   | Type    | Default | Description                          |
| ----------- | ------- | ------- | ------------------------------------ |
| `page`      | integer | 1       | Halaman                              |
| `page_size` | integer | 10      | Jumlah data per halaman              |

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Inactive doctor specializations retrieved successfully",
  "data": {
    "data": [],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total_items": 0,
      "total_pages": 0
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/doctor-specializations/inactive" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

| Method | Endpoint                                 | Description                       | Role Required      |
| ------ | ---------------------------------------- | --------------------------------- | ------------------ |
| GET    | `/doctor-specializations/inactive` | List inactive doctor specializations | Admin, Super Admin |
| GET    | `/doctor-specializations/deleted`        | List deleted specializations      | Admin, Super Admin |
| POST   | `/doctor-specializations`                | Create doctor specialization      | Admin, Super Admin |
| PUT    | `/doctor-specializations/:id`            | Update doctor specialization      | Admin, Super Admin |
| PATCH  | `/doctor-specializations/:id/activate`   | Activate specialization           | Admin, Super Admin |
| PATCH  | `/doctor-specializations/:id/deactivate` | Deactivate specialization         | Admin, Super Admin |
| DELETE | `/doctor-specializations/:id`            | Soft delete doctor specialization | Admin, Super Admin |
| PATCH  | `/doctor-specializations/:id/restore`    | Restore deleted specialization    | Admin, Super Admin |

### Super Admin Endpoints

| Method | Endpoint                                  | Description                              | Role Required |
| ------ | ----------------------------------------- | ---------------------------------------- | ------------- |
| DELETE | `/doctor-specializations/:id/hard-delete` | Permanently delete doctor specialization | Super Admin   |

---

## Public Endpoints

### 1. List Doctor Specializations

**Endpoint:** `GET /api/v1/doctor-specializations`

**Description:** Mendapatkan daftar semua doctor specializations dengan pagination dan filter.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter   | Type    | Default | Description                        |
| ----------- | ------- | ------- | ---------------------------------- |
| `page`      | integer | 1       | Halaman                            |
| `page_size` | integer | 10      | Jumlah data per halaman (max: 100) |
| `search`    | string  | -       | Cari berdasarkan nama spesialisasi |
| `is_active` | boolean | -       | Filter by active status            |
| `sort_by`   | string  | name    | Sort field (name, created_at)      |
| `sort_dir`  | string  | asc     | Sort direction (asc, desc)         |

**Example Request:**

```
GET /api/v1/doctor-specializations?page=1&page_size=10&is_active=true&sort_by=name&sort_dir=asc
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor specializations retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "Kardiologi",
        "description": "Spesialisasi penyakit jantung dan pembuluh darah",
        "code": "CARDIO",
        "is_active": true,
        "created_at": "2024-01-19T10:00:00Z"
      },
      {
        "id": 2,
        "name": "Ortopedi",
        "description": "Spesialisasi tulang dan sendi",
        "code": "ORTHO",
        "is_active": true,
        "created_at": "2024-01-19T10:05:00Z"
      },
      {
        "id": 3,
        "name": "Neurologi",
        "description": "Spesialisasi saraf dan otak",
        "code": "NEURO",
        "is_active": true,
        "created_at": "2024-01-19T10:10:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total_items": 3,
      "total_pages": 1
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/doctor-specializations?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 2. List Active Doctor Specializations

**Endpoint:** `GET /api/v1/doctor-specializations/active`

**Description:** Mendapatkan daftar doctor specializations yang aktif saja.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter   | Type    | Default | Description                        |
| ----------- | ------- | ------- | ---------------------------------- |
| `page`      | integer | 1       | Halaman                            |
| `page_size` | integer | 10      | Jumlah data per halaman (max: 100) |
| `search`    | string  | -       | Cari berdasarkan nama spesialisasi |
| `sort_by`   | string  | name    | Sort field (name, created_at)      |
| `sort_dir`  | string  | asc     | Sort direction (asc, desc)         |

**Example Request:**

```
GET /api/v1/doctor-specializations/active?page=1&page_size=10&sort_by=name&sort_dir=asc
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Active doctor specializations retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "Kardiologi",
        "description": "Spesialisasi penyakit jantung dan pembuluh darah",
        "code": "CARDIO",
        "is_active": true,
        "created_at": "2024-01-19T10:00:00Z"
      },
      {
        "id": 2,
        "name": "Ortopedi",
        "description": "Spesialisasi tulang dan sendi",
        "code": "ORTHO",
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
curl -X GET "http://localhost:8080/api/v1/doctor-specializations/active?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 3. Get Doctor Specialization by ID

**Endpoint:** `GET /api/v1/doctor-specializations/:id`

**Description:** Mendapatkan detail doctor specialization berdasarkan ID.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Doctor Specialization ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor specialization retrieved successfully",
  "data": {
    "id": 1,
    "name": "Kardiologi",
    "description": "Spesialisasi penyakit jantung dan pembuluh darah",
    "code": "CARDIO",
    "is_active": true,
    "total_doctors": 5,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Doctor specialization not found",
  "error": "specialization not found"
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/doctor-specializations/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Admin Endpoints

### X. List Inactive Doctor Specializations

**Endpoint:** `GET /api/v1/doctor-specializations/inactive`

**Description:** Mendapatkan daftar doctor specializations yang tidak aktif.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**Query Parameters:**

| Parameter   | Type    | Default | Description                          |
| ----------- | ------- | ------- | ------------------------------------ |
| `page`      | integer | 1       | Halaman                              |
| `page_size` | integer | 10      | Jumlah data per halaman              |

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Inactive doctor specializations retrieved successfully",
  "data": {
    "data": [],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total_items": 0,
      "total_pages": 0
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/doctor-specializations/inactive" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 4. List Deleted Doctor Specializations

**Endpoint:** `GET /api/v1/doctor-specializations/deleted`

**Description:** Mendapatkan daftar doctor specializations yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter   | Type    | Default | Description                        |
| ----------- | ------- | ------- | ---------------------------------- |
| `page`      | integer | 1       | Halaman                            |
| `page_size` | integer | 10      | Jumlah data per halaman (max: 100) |
| `search`    | string  | -       | Cari berdasarkan nama spesialisasi |
| `sort_by`   | string  | name    | Sort field (name, created_at)      |
| `sort_dir`  | string  | asc     | Sort direction (asc, desc)         |

**Example Request:**

```
GET /api/v1/doctor-specializations/deleted?page=1&page_size=10
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Deleted doctor specializations retrieved successfully",
  "data": {
    "data": [
      {
        "id": 5,
        "name": "Pediatri",
        "description": "Spesialisasi anak-anak",
        "code": "PEDIA",
        "is_active": false,
        "created_at": "2024-01-19T10:00:00Z",
        "updated_at": "2024-01-20T10:00:00Z",
        "deleted_at": "2024-01-20T10:00:00Z"
      },
      {
        "id": 6,
        "name": "Dermatologi",
        "description": "Spesialisasi kulit",
        "code": "DERM",
        "is_active": false,
        "created_at": "2024-01-19T10:00:00Z",
        "updated_at": "2024-01-20T10:00:00Z",
        "deleted_at": "2024-01-20T10:00:00Z"
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
curl -X GET "http://localhost:8080/api/v1/doctor-specializations/deleted?page=1&page_size=10" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 5. Create Doctor Specialization

**Endpoint:** `POST /api/v1/doctor-specializations`

**Description:** Admin membuat spesialisasi dokter baru.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "name": "Kardiologi",
  "description": "Spesialisasi penyakit jantung dan pembuluh darah",
  "code": "CARDIO",
  "is_active": true
}
```

**Field Rules:**

- `name`: required, unique, max 100 characters
- `description`: optional, text
- `code`: optional, unique, max 20 characters
- `is_active`: optional, boolean (default: true)

**Response Success (201 Created):**

```json
{
  "success": true,
  "message": "Doctor specialization created successfully",
  "data": {
    "id": 1,
    "name": "Kardiologi",
    "description": "Spesialisasi penyakit jantung dan pembuluh darah",
    "code": "CARDIO",
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
  "message": "Failed to create doctor specialization",
  "error": "name already exists"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/doctor-specializations \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Kardiologi",
    "description": "Spesialisasi penyakit jantung dan pembuluh darah",
    "code": "CARDIO",
    "is_active": true
  }'
```

---

### 6. Update Doctor Specialization

**Endpoint:** `PUT /api/v1/doctor-specializations/:id`

**Description:** Admin mengupdate data doctor specialization berdasarkan ID.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Doctor Specialization ID (integer)

**Request Body:**

```json
{
  "name": "Kardiologi - Updated",
  "description": "Spesialisasi penyakit jantung, pembuluh darah, dan sistem sirkulasi",
  "code": "CARDIO",
  "is_active": true
}
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor specialization updated successfully",
  "data": {
    "id": 1,
    "name": "Kardiologi - Updated",
    "description": "Spesialisasi penyakit jantung, pembuluh darah, dan sistem sirkulasi",
    "code": "CARDIO",
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T15:30:00Z"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Doctor specialization not found",
  "error": "specialization not found"
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/doctor-specializations/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Kardiologi - Updated",
    "description": "Updated description",
    "code": "CARDIO",
    "is_active": true
  }'
```

---

### 7. Activate Doctor Specialization

**Endpoint:** `PATCH /api/v1/doctor-specializations/:id/activate`

**Description:** Admin mengaktifkan doctor specialization yang tidak aktif.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Doctor Specialization ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor specialization activated successfully",
  "data": {
    "id": 1,
    "name": "Kardiologi",
    "code": "CARDIO",
    "is_active": true,
    "updated_at": "2024-01-19T16:00:00Z"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Doctor specialization not found",
  "error": "specialization not found"
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/doctor-specializations/1/activate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 8. Deactivate Doctor Specialization

**Endpoint:** `PATCH /api/v1/doctor-specializations/:id/deactivate`

**Description:** Admin menonaktifkan doctor specialization.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Doctor Specialization ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor specialization deactivated successfully",
  "data": {
    "id": 1,
    "name": "Kardiologi",
    "code": "CARDIO",
    "is_active": false,
    "updated_at": "2024-01-19T16:00:00Z"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Doctor specialization not found",
  "error": "specialization not found"
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/doctor-specializations/1/deactivate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 9. Delete Doctor Specialization

**Endpoint:** `DELETE /api/v1/doctor-specializations/:id`

**Description:** Admin melakukan soft delete doctor specialization.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Doctor Specialization ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor specialization deleted successfully",
  "data": {
    "id": 1,
    "name": "Kardiologi",
    "is_active": false,
    "deleted_at": "2024-01-19T15:30:00Z"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Doctor specialization not found",
  "error": "specialization not found"
}
```

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/doctor-specializations/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

## Super Admin Endpoints

### 10. Restore Doctor Specialization

**Endpoint:** `PATCH /api/v1/doctor-specializations/:id/restore`

**Description:** Admin me-restore doctor specialization yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Doctor Specialization ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor specialization restored successfully",
  "data": {
    "id": 1,
    "name": "Kardiologi",
    "code": "CARDIO",
    "is_active": false,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-20T15:30:00Z",
    "deleted_at": null
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Doctor specialization not found",
  "error": "specialization not found"
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/doctor-specializations/1/restore \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

**Notes:**

- Setelah di-restore, specialization masih dalam status `is_active: false`
- Perlu activate manual jika ingin mengaktifkan kembali

---

## Super Admin Endpoints

### 11. Hard Delete Doctor Specialization

### 12. Hard Delete Doctor Specialization

**Endpoint:** `DELETE /api/v1/doctor-specializations/:id/hard-delete`

**Description:** Super Admin melakukan permanent delete doctor specialization.

**Authentication:** Required (Super Admin)

**Request Headers:**

```
Authorization: Bearer <super-admin-token>
```

**URL Parameters:**

- `id`: Doctor Specialization ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor specialization permanently deleted",
  "data": null
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Doctor specialization not found",
  "error": "specialization not found"
}
```

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/doctor-specializations/1/hard-delete \
  -H "Authorization: Bearer SUPER_ADMIN_JWT_TOKEN"
```

---

## Database Model

### Table: doctor_specializations

| Field       | Type         | Constraints                                           | Description                                     |
| ----------- | ------------ | ----------------------------------------------------- | ----------------------------------------------- |
| id          | BIGINT       | PRIMARY KEY, AUTO_INCREMENT                           | Unique identifier                               |
| name        | VARCHAR(100) | NOT NULL, UNIQUE, INDEX                               | Nama spesialisasi (e.g., Kardiologi, Neurologi) |
| code        | VARCHAR(20)  | UNIQUE, INDEX                                         | Kode spesialisasi (e.g., KARDIO, NEURO)         |
| description | TEXT         | NULLABLE                                              | Deskripsi detail spesialisasi                   |
| is_active   | BOOLEAN      | NOT NULL, DEFAULT true, INDEX                         | Status aktif                                    |
| created_at  | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP                             | Waktu pembuatan record                          |
| updated_at  | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir                           |
| deleted_at  | TIMESTAMP    | INDEX, NULLABLE                                       | Soft delete timestamp                           |

**Indexes:**

- Primary Key: id
- Unique Index: name, code
- Regular Index: is_active, deleted_at

**Relationships:**

- Has Many Doctors (one-to-many) - doctors reference specialization

**Notes:**

- Master data untuk standardisasi spesialisasi dokter
- Code sebaiknya uppercase (KARDIO, NEURO, ORTHO, dll)
- is_active untuk hide inactive specializations dari list

---

## Error Responses

### Common Error Codes

| Status Code | Message               | Description                                  |
| ----------- | --------------------- | -------------------------------------------- |
| 400         | Bad Request           | Invalid request body or query parameters     |
| 401         | Unauthorized          | Missing or invalid JWT token                 |
| 403         | Forbidden             | User does not have permission                |
| 404         | Not Found             | Resource not found                           |
| 409         | Conflict              | Duplicate data (e.g., specialization exists) |
| 500         | Internal Server Error | Unexpected server error                      |

### Example Error Response

```json
{
  "success": false,
  "message": "Invalid request",
  "error": "name is required"
}
```

---

## Notes

- Spesialisasi dokter digunakan untuk menstandarisasi spesialisasi pada data dokter
- Nama spesialisasi bersifat unique dan case-sensitive
- Soft delete hanya mengubah status, data tetap tersimpan
- Hard delete hanya bisa dilakukan oleh Super Admin dan bersifat permanent
