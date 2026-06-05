# Medicine Types API Documentation

## Overview

API untuk manajemen master data tipe obat dalam sistem rekam medis. Medicine Types adalah master data yang dapat direferensikan oleh Medicines untuk menstandarisasi tipe obat (menghindari duplikasi data seperti "tablet" vs "Tablet").

**Base URL:** `/api/v1/medicine-types`

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

| Endpoint                            | Patient | Doctor | Receptionist | Admin | Super Admin |
| ----------------------------------- | ------- | ------ | ------------ | ----- | ----------- |
| GET /medicine-types                 | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /medicine-types/active          | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /medicine-types/inactive              | ❌      | ❌     | ❌           | ✅    | ✅          |
| GET /medicine-types/:id             | ✅      | ✅     | ✅           | ✅    | ✅          |
| POST /medicine-types                | ❌      | ❌     | ❌           | ✅    | ✅          |
| PUT /medicine-types/:id             | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /medicine-types/:id          | ❌      | ❌     | ❌           | ✅    | ✅          |
| GET /medicine-types/deleted             | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /medicine-types/:id/activate      | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /medicine-types/:id/deactivate    | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /medicine-types/:id/restore       | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /medicine-types/:id/hard-delete | ❌      | ❌     | ❌           | ❌    | ✅          |

---

## Endpoints Summary

### Public Endpoints (All Authenticated)

| Method | Endpoint                  | Description                    | Role Required                            |
| ------ | ------------------------- | ------------------------------ | ---------------------------------------- |
| GET    | `/medicine-types`         | List all medicine types        | All Authenticated                        |
| GET    | `/medicine-types/active`  | List active medicine types     | All Authenticated                        |
| GET    | `/medicine-types/:id`     | Get medicine type by ID        | All Authenticated                        |

### Admin Endpoints

### X. List Inactive Medicine Types

**Endpoint:** `GET /api/v1/medicine-types/inactive`

**Description:** Mendapatkan daftar medicine types yang tidak aktif.

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
  "message": "Inactive medicine types retrieved successfully",
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
curl -X GET "http://localhost:8080/api/v1/medicine-types/inactive" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

| Method | Endpoint                | Description               | Role Required      |
| ------ | ----------------------- | ------------------------- | ------------------ |
| GET    | `/medicine-types/inactive` | List inactive medicine types | Admin, Super Admin |
| POST   | `/medicine-types`       | Create medicine type      | Admin, Super Admin |
| PUT    | `/medicine-types/:id`   | Update medicine type      | Admin, Super Admin |
| DELETE | `/medicine-types/:id`   | Soft delete medicine type | Admin, Super Admin |

### Super Admin Endpoints

| Method | Endpoint                         | Description                      | Role Required |
| ------ | -------------------------------- | -------------------------------- | ------------- |
| DELETE | `/medicine-types/:id/hard-delete` | Permanently delete medicine type | Super Admin   |

---

## Public Endpoints

### 1. List Medicine Types

**Endpoint:** `GET /api/v1/medicine-types`

**Description:** Mendapatkan daftar semua medicine types dengan pagination dan filter.

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
| `search`    | string  | -       | Cari berdasarkan nama tipe obat      |
| `is_active` | boolean | -       | Filter by active status              |
| `sort_by`   | string  | name    | Sort field (name, created_at)        |
| `sort_dir`  | string  | asc     | Sort direction (asc, desc)           |

**Example Request:**

```
GET /api/v1/medicine-types?page=1&page_size=10&is_active=true&sort_by=name&sort_dir=asc
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine types retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "Tablet",
        "description": "Bentuk obat padat dalam bentuk tablet",
        "code": "TBL",
        "is_active": true,
        "created_at": "2024-01-19T10:00:00Z"
      },
      {
        "id": 2,
        "name": "Capsule",
        "description": "Bentuk obat yang dibungkus dengan gelatin",
        "code": "CAP",
        "is_active": true,
        "created_at": "2024-01-19T10:05:00Z"
      },
      {
        "id": 3,
        "name": "Syrup",
        "description": "Bentuk obat cair dengan rasa manis",
        "code": "SYR",
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
curl -X GET "http://localhost:8080/api/v1/medicine-types?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 2. List Active Medicine Types

**Endpoint:** `GET /api/v1/medicine-types/active`

**Description:** Mendapatkan daftar medicine types yang aktif saja.

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
GET /api/v1/medicine-types/active?page=1&page_size=10
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Active medicine types retrieved successfully",
  "data": {
    "total_active": 6,
    "data": [
      {
        "id": 1,
        "name": "Tablet",
        "description": "Bentuk obat padat dalam bentuk tablet",
        "code": "TBL",
        "is_active": true,
        "created_at": "2024-01-19T10:00:00Z"
      },
      {
        "id": 2,
        "name": "Capsule",
        "description": "Bentuk obat yang dibungkus dengan gelatin",
        "code": "CAP",
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
curl -X GET "http://localhost:8080/api/v1/medicine-types/active" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 3. Get Medicine Type by ID

**Endpoint:** `GET /api/v1/medicine-types/:id`

**Description:** Mendapatkan detail medicine type berdasarkan ID.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Medicine Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine type retrieved successfully",
  "data": {
    "id": 1,
    "name": "Tablet",
    "description": "Bentuk obat padat dalam bentuk tablet",
    "code": "TBL",
    "is_active": true,
    "total_medicines": 25,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Medicine type not found",
  "error": "medicine type not found"
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/medicine-types/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Admin Endpoints

### X. List Inactive Medicine Types

**Endpoint:** `GET /api/v1/medicine-types/inactive`

**Description:** Mendapatkan daftar medicine types yang tidak aktif.

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
  "message": "Inactive medicine types retrieved successfully",
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
curl -X GET "http://localhost:8080/api/v1/medicine-types/inactive" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 4. Create Medicine Type

**Endpoint:** `POST /api/v1/medicine-types`

**Description:** Admin membuat tipe obat baru.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "name": "Tablet",
  "description": "Bentuk obat padat dalam bentuk tablet",
  "code": "TBL",
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
  "message": "Medicine type created successfully",
  "data": {
    "id": 1,
    "name": "Tablet",
    "description": "Bentuk obat padat dalam bentuk tablet",
    "code": "TBL",
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
  "message": "Failed to create medicine type",
  "error": "name already exists"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/medicine-types \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tablet",
    "description": "Bentuk obat padat dalam bentuk tablet",
    "code": "TBL",
    "is_active": true
  }'
```

---

### 5. Update Medicine Type

**Endpoint:** `PUT /api/v1/medicine-types/:id`

**Description:** Admin mengupdate data medicine type berdasarkan ID.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Medicine Type ID (integer)

**Request Body:**

```json
{
  "name": "Tablet - Updated",
  "description": "Bentuk obat padat dalam bentuk tablet dengan coating",
  "code": "TBL",
  "is_active": true
}
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine type updated successfully",
  "data": {
    "id": 1,
    "name": "Tablet - Updated",
    "description": "Bentuk obat padat dalam bentuk tablet dengan coating",
    "code": "TBL",
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
  "message": "Medicine type not found",
  "error": "medicine type not found"
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/medicine-types/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Tablet - Updated",
    "description": "Updated description",
    "code": "TBL",
    "is_active": true
  }'
```


---

### 6. Activate Medicine Type

**Endpoint:** `PATCH /api/v1/medicine-types/:id/activate`

**Description:** Admin mengaktifkan kembali medicine type.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Medicine Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine type activated successfully",
  "data": {
    "id": 1,
    "name": "Tablet",
    "is_active": true,
    "updated_at": "2024-01-19T16:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/medicine-types/1/activate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 7. Deactivate Medicine Type

**Endpoint:** `PATCH /api/v1/medicine-types/:id/deactivate`

**Description:** Admin menonaktifkan medicine type.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Medicine Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine type deactivated successfully",
  "data": {
    "id": 1,
    "name": "Tablet",
    "is_active": false,
    "updated_at": "2024-01-19T16:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/medicine-types/1/deactivate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 8. Soft Delete Medicine Type

**Endpoint:** `DELETE /api/v1/medicine-types/:id`

**Description:** Admin melakukan soft delete medicine type.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Medicine Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine type deleted successfully",
  "data": {
    "id": 1,
    "name": "Tablet",
    "is_active": false,
    "deleted_at": "2024-01-19T15:30:00Z"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Medicine type not found",
  "error": "medicine type not found"
}
```

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/medicine-types/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```


---

### 9. List Deleted Medicine Types

**Endpoint:** `GET /api/v1/medicine-types/deleted`

**Description:** Mendapatkan daftar medicine types yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
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
  "message": "Deleted medicine types retrieved successfully",
  "data": {
    "data": [
      {
        "id": 5,
        "name": "Serbuk Lama",
        "description": "Tidak dipakai lagi",
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
curl -X GET "http://localhost:8080/api/v1/medicine-types/deleted" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 10. Restore Medicine Type

**Endpoint:** `PATCH /api/v1/medicine-types/:id/restore`

**Description:** Admin memulihkan (restore) medicine type yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Medicine Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine type restored successfully",
  "data": {
    "id": 5,
    "name": "Serbuk Lama",
    "is_active": false,
    "deleted_at": null,
    "updated_at": "2024-01-19T16:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/medicine-types/5/restore \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

## Super Admin Endpoints

### 11. Hard Delete Medicine Type

**Endpoint:** `DELETE /api/v1/medicine-types/:id/hard-delete`

**Description:** Super Admin melakukan permanent delete medicine type.

**Authentication:** Required (Super Admin)

**Request Headers:**

```
Authorization: Bearer <super-admin-token>
```

**URL Parameters:**

- `id`: Medicine Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine type permanently deleted",
  "data": null
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Medicine type not found",
  "error": "medicine type not found"
}
```

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/medicine-types/1/hard-delete \
  -H "Authorization: Bearer SUPER_ADMIN_JWT_TOKEN"
```

---

## Database Model

### Table: medicine_types

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| name | VARCHAR(100) | NOT NULL, UNIQUE, INDEX | Nama tipe obat (e.g., Tablet, Kapsul, Suntik) |
| code | VARCHAR(20) | UNIQUE, INDEX | Kode tipe obat (e.g., TAB, KAP, SUN, KRIM) |
| description | TEXT | NULLABLE | Deskripsi detail tipe obat |
| is_active | BOOLEAN | NOT NULL, DEFAULT true, INDEX | Status aktif |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

**Indexes:**
- Primary Key: id
- Unique Index: name, code
- Regular Index: is_active, deleted_at

**Relationships:**
- Has Many Medicines (one-to-many) - medicines reference type

**Notes:**
- Master data untuk standardisasi tipe obat/format
- Code sebaiknya uppercase (TAB, KAP, SUN, KRIM, SIRUP, dll)
- Contoh tipe: Tablet, Kapsul, Suntik, Krim, Sirup, Salep
- is_active untuk hide inactive types

---

## Error Responses

### Common Error Codes

| Status Code | Message                           | Description                               |
| ----------- | --------------------------------- | ----------------------------------------- |
| 400         | Bad Request                       | Invalid request body or query parameters  |
| 401         | Unauthorized                      | Missing or invalid JWT token              |
| 403         | Forbidden                         | User does not have permission             |
| 404         | Not Found                         | Resource not found                        |
| 409         | Conflict                          | Duplicate data (e.g., medicine type exists) |
| 500         | Internal Server Error             | Unexpected server error                   |

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

- Tipe obat digunakan untuk menstandarisasi tipe obat pada data medicines
- Nama tipe obat bersifat unique dan case-sensitive
- Soft delete hanya mengubah status, data tetap tersimpan
- Hard delete hanya bisa dilakukan oleh Super Admin dan bersifat permanent
