# Test Types Categories API Documentation

## Overview

API untuk manajemen master data kategori tes laboratorium dalam sistem rekam medis. Categories digunakan sebagai master data yang dapat direferensikan oleh Test Types untuk menjaga konsistensi data dan menghindari duplikasi data (lowercase/uppercase).

**Base URL:** `/api/v1/test-categories`

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

| Endpoint                                 | Patient | Doctor | Receptionist | Admin | Super Admin |
| ---------------------------------------- | ------- | ------ | ------------ | ----- | ----------- |
| GET /test-categories                     | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /test-categories/active              | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /test-categories/:id                 | ✅      | ✅     | ✅           | ✅    | ✅          |
| POST /test-categories                    | ❌      | ❌     | ❌           | ✅    | ✅          |
| PUT /test-categories/:id                 | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /test-categories/:id              | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /test-categories/:id/hard-delete  | ❌      | ❌     | ❌           | ❌    | ✅          |

---

## Endpoints Summary

### Public Endpoints (All Authenticated)

| Method | Endpoint                      | Description                        | Role Required                            |
| ------ | ----------------------------- | ---------------------------------- | ---------------------------------------- |
| GET    | `/test-categories`            | List all test categories           | All Authenticated                        |
| GET    | `/test-categories/active`     | List active test categories        | All Authenticated                        |
| GET    | `/test-categories/:id`        | Get test category by ID            | All Authenticated                        |

### Admin Endpoints

| Method | Endpoint               | Description               | Role Required      |
| ------ | ---------------------- | ------------------------- | ------------------ |
| POST   | `/test-categories`     | Create test category      | Admin, Super Admin |
| PUT    | `/test-categories/:id` | Update test category      | Admin, Super Admin |
| DELETE | `/test-categories/:id` | Soft delete test category | Admin, Super Admin |

### Super Admin Endpoints

| Method | Endpoint                            | Description                      | Role Required |
| ------ | ----------------------------------- | -------------------------------- | ------------- |
| DELETE | `/test-categories/:id/hard-delete` | Permanently delete test category | Super Admin   |

---

## Public Endpoints

### 1. List Test Categories

**Endpoint:** `GET /api/v1/test-categories`

**Description:** Mendapatkan daftar semua test categories dengan pagination dan filter.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter   | Type    | Default | Description                          |
| ----------- | ------- | ------- | ------------------------------------ |
| `page`      | integer | 1       | Halaman                              |
| `page_size` | integer | 10      | Jumlah data per halaman (max: 100)   |
| `search`    | string  | -       | Cari berdasarkan nama kategori       |
| `is_active` | boolean | -       | Filter by active status              |
| `sort_by`   | string  | name    | Sort field (name, created_at)        |
| `sort_dir`  | string  | asc     | Sort direction (asc, desc)           |

**Example Request:**

```
GET /api/v1/test-categories?page=1&page_size=10&is_active=true&sort_by=name&sort_dir=asc
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Test categories retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "Hematologi",
        "description": "Pemeriksaan yang berkaitan dengan darah dan komponen darah",
        "is_active": true,
        "created_at": "2024-01-19T10:00:00Z"
      },
      {
        "id": 2,
        "name": "Kimia Darah",
        "description": "Pemeriksaan kimia komponen darah",
        "is_active": true,
        "created_at": "2024-01-19T10:05:00Z"
      },
      {
        "id": 3,
        "name": "Mikrobiologi",
        "description": "Pemeriksaan mikroorganisme",
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
curl -X GET "http://localhost:8080/api/v1/test-categories?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 2. List Active Test Categories

**Endpoint:** `GET /api/v1/test-categories/active`

**Description:** Mendapatkan daftar test categories yang aktif saja.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter   | Type    | Description              |
| ----------- | ------- | ------------------------ |
| `page`      | integer | Page number              |
| `page_size` | integer | Items per page           |
| `sort_by`   | string  | Sort field (name, id)    |
| `sort_dir`  | string  | Sort direction (asc, desc) |

**Example Request:**

```
GET /api/v1/test-categories/active?page=1&page_size=10
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Active test categories retrieved successfully",
  "data": {
    "total_active": 3,
    "data": [
      {
        "id": 1,
        "name": "Hematologi",
        "description": "Pemeriksaan darah",
        "is_active": true,
        "created_at": "2024-01-19T10:00:00Z"
      },
      {
        "id": 2,
        "name": "Kimia Darah",
        "description": "Pemeriksaan kimia darah",
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
curl -X GET "http://localhost:8080/api/v1/test-categories/active" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 3. Get Test Category by ID

**Endpoint:** `GET /api/v1/test-categories/:id`

**Description:** Mendapatkan detail test category berdasarkan ID.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Test Category ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Test category retrieved successfully",
  "data": {
    "id": 1,
    "name": "Hematologi",
    "description": "Pemeriksaan yang berkaitan dengan darah dan komponen darah",
    "is_active": true,
    "total_tests": 5,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Test category not found",
  "error": "category not found"
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/test-categories/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Admin Endpoints

### 4. Create Test Category

**Endpoint:** `POST /api/v1/test-categories`

**Description:** Admin membuat kategori test type baru.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "name": "Hematologi",
  "description": "Pemeriksaan yang berkaitan dengan darah dan komponen darah",
  "is_active": true
}
```

**Field Rules:**

- `name`: required, unique, max 100 characters
- `description`: optional, text
- `is_active`: optional, boolean (default: true)

**Response Success (201 Created):**

```json
{
  "success": true,
  "message": "Test category created successfully",
  "data": {
    "id": 1,
    "name": "Hematologi",
    "description": "Pemeriksaan yang berkaitan dengan darah dan komponen darah",
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
  "message": "Failed to create test category",
  "error": "name already exists"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/test-categories \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hematologi",
    "description": "Pemeriksaan yang berkaitan dengan darah",
    "is_active": true
  }'
```

---

### 5. Update Test Category

**Endpoint:** `PUT /api/v1/test-categories/:id`

**Description:** Admin mengupdate data test category berdasarkan ID.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Test Category ID (integer)

**Request Body:**

```json
{
  "name": "Hematologi - Updated",
  "description": "Pemeriksaan darah lengkap dan detail",
  "is_active": true
}
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Test category updated successfully",
  "data": {
    "id": 1,
    "name": "Hematologi - Updated",
    "description": "Pemeriksaan darah lengkap dan detail",
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
  "message": "Test category not found",
  "error": "category not found"
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/test-categories/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hematologi - Updated",
    "description": "Updated description",
    "is_active": true
  }'
```

---

### 6. Delete Test Category

**Endpoint:** `DELETE /api/v1/test-categories/:id`

**Description:** Admin melakukan soft delete test category.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Test Category ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Test category deleted successfully",
  "data": {
    "id": 1,
    "name": "Hematologi",
    "is_active": false,
    "deleted_at": "2024-01-19T15:30:00Z"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Test category not found",
  "error": "category not found"
}
```

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/test-categories/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

## Super Admin Endpoints

### 7. Hard Delete Test Category

**Endpoint:** `DELETE /api/v1/test-categories/:id/hard-delete`

**Description:** Super Admin melakukan permanent delete test category.

**Authentication:** Required (Super Admin)

**Request Headers:**

```
Authorization: Bearer <super-admin-token>
```

**URL Parameters:**

- `id`: Test Category ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Test category permanently deleted",
  "data": null
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Test category not found",
  "error": "category not found"
}
```

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/test-categories/1/hard-delete \
  -H "Authorization: Bearer SUPER_ADMIN_JWT_TOKEN"
```

---

## Database Model

### Table: test_categories

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| name | VARCHAR(100) | NOT NULL, UNIQUE, INDEX | Nama kategori tes (e.g., Hematologi, Biokimia, Urinalisis) |
| code | VARCHAR(20) | UNIQUE, INDEX | Kode kategori (e.g., HEM, BIO, URI) |
| description | TEXT | NULLABLE | Deskripsi detail kategori |
| is_active | BOOLEAN | NOT NULL, DEFAULT true, INDEX | Status aktif |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

**Indexes:**
- Primary Key: id
- Unique Index: name, code
- Regular Index: is_active, deleted_at

**Relationships:**
- Has Many Test Types (one-to-many) - test types reference category

**Notes:**
- Master data untuk standardisasi kategori tes laboratorium
- Code sebaiknya uppercase (HEM, BIO, URI, MIK, dll)
- Digunakan untuk organisir berbagai jenis tes
- is_active untuk hide inactive categories

---

## Error Responses

### Common Error Codes

| Status Code | Message                           | Description                                |
| ----------- | --------------------------------- | ------------------------------------------ |
| 400         | Bad Request                       | Invalid request body or query parameters   |
| 401         | Unauthorized                      | Missing or invalid JWT token               |
| 403         | Forbidden                         | User does not have permission              |
| 404         | Not Found                         | Resource not found                         |
| 409         | Conflict                          | Duplicate data (e.g., category name exists)|
| 500         | Internal Server Error             | Unexpected server error                    |

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

- Kategori test digunakan untuk menstandarisasi kategori pada test types
- Nama kategori bersifat unique dan case-sensitive
- Soft delete hanya mengubah status, data tetap tersimpan
- Hard delete hanya bisa dilakukan oleh Super Admin dan bersifat permanent
