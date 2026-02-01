# Departments API Documentation

## Overview

API untuk manajemen data departments (departemen/bagian) dalam sistem rekam medis. Department merupakan divisi atau bagian dalam rumah sakit seperti Kardiologi, Neurologi, dll.

**Base URL:** `/api/v1/departments`

---

## Table of Contents

- [Authentication](#authentication)
- [Authorization](#authorization)
- [Endpoints Summary](#endpoints-summary)
- [Department Endpoints](#department-endpoints)
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
| GET /departments                    | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /departments/:id                | ✅      | ✅     | ✅           | ✅    | ✅          |
| POST /departments                   | ❌      | ❌     | ❌           | ✅    | ✅          |
| PUT /departments/:id                | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /departments/:id             | ❌      | ❌     | ❌           | ✅    | ✅          |
| GET /departments/deleted            | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /departments/:id/restore      | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /departments/:id/hard-delete | ❌      | ❌     | ❌           | ❌    | ✅          |

---

## Endpoints Summary

| Method | Endpoint                       | Description                   | Auth Required      |
| ------ | ------------------------------ | ----------------------------- | ------------------ |
| GET    | `/departments`                 | List active departments       | All Authenticated  |
| GET    | `/departments/deleted`         | List deleted departments      | Admin, Super Admin |
| GET    | `/departments/:id`             | Get department by ID          | All Authenticated  |
| POST   | `/departments`                 | Create department             | Admin, Super Admin |
| PUT    | `/departments/:id`             | Update department             | Admin, Super Admin |
| DELETE | `/departments/:id`             | Soft delete department        | Admin, Super Admin |
| PATCH  | `/departments/:id/restore`     | Restore deleted department    | Admin, Super Admin |
| DELETE | `/departments/:id/hard-delete` | Permanently delete department | Super Admin        |

---

## Department Endpoints

### 1. List Departments

**Endpoint:** `GET /api/v1/departments`

**Description:** Mendapatkan daftar department aktif dengan pagination, search, dan filter.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter   | Type    | Default    | Description                                   |
| ----------- | ------- | ---------- | --------------------------------------------- |
| `page`      | integer | 1          | Halaman                                       |
| `page_size` | integer | 10         | Jumlah data per halaman (max: 100)            |
| `search`    | string  | -          | Cari berdasarkan name, code, atau description |
| `sort_by`   | string  | created_at | Field untuk sorting (created_at, name, code)  |
| `sort_dir`  | string  | desc       | Arah sorting (asc, desc)                      |

**Example Request:**

```
GET /api/v1/departments?page=1&page_size=10&search=kardio&sort_by=name&sort_dir=asc
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Departments retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "Kardiologi",
        "code": "KARDIO",
        "description": "Departemen jantung dan pembuluh darah",
        "floor_location": "Lantai 3",
        "created_at": "2024-01-19T10:00:00Z",
        "updated_at": "2024-01-19T10:00:00Z"
      },
      {
        "id": 2,
        "name": "Neurologi",
        "code": "NEURO",
        "description": "Departemen saraf",
        "floor_location": "Lantai 2",
        "created_at": "2024-01-19T10:05:00Z",
        "updated_at": "2024-01-19T10:05:00Z"
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
curl -X GET "http://localhost:8080/api/v1/departments?page=1&page_size=10&search=kardio" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 2. List Deleted Departments

**Endpoint:** `GET /api/v1/departments/deleted`

**Description:** Mendapatkan daftar department yang sudah di-soft delete.

**Authentication:** Required (Admin/Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**Query Parameters:**
Same as List Departments endpoint.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Deleted departments retrieved successfully",
  "data": {
    "data": [
      {
        "id": 5,
        "name": "Deleted Department",
        "code": "DEL01",
        "description": "Department yang dihapus",
        "floor_location": "Lantai 1",
        "created_at": "2024-01-18T10:00:00Z",
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
curl -X GET "http://localhost:8080/api/v1/departments/deleted?page=1&page_size=10" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 3. Get Department by ID

**Endpoint:** `GET /api/v1/departments/:id`

**Description:** Mendapatkan detail department berdasarkan ID, termasuk daftar doctors dan rooms.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Department ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Department retrieved successfully",
  "data": {
    "id": 1,
    "name": "Kardiologi",
    "code": "KARDIO",
    "description": "Departemen jantung dan pembuluh darah",
    "floor_location": "Lantai 3",
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z",
    "doctors": [
      {
        "id": 1,
        "full_name": "Dr. John Smith",
        "specialization": "Cardiologist",
        "license_number": "DOC-001"
      }
    ],
    "rooms": [
      {
        "id": 1,
        "room_number": "301-A",
        "room_type": "class_1",
        "bed_capacity": 2,
        "available_beds": 1
      }
    ]
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Department not found",
  "error": "department not found"
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/departments/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 4. Create Department

**Endpoint:** `POST /api/v1/departments`

**Description:** Admin membuat department baru.

**Authentication:** Required (Admin/Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "name": "Kardiologi",
  "code": "KARDIO",
  "description": "Departemen jantung dan pembuluh darah",
  "floor_location": "Lantai 3"
}
```

**Field Rules:**

- `name`: required, max 100 characters, indexed
- `code`: required, unique, max 20 characters, indexed
- `description`: optional, text
- `floor_location`: optional, max 50 characters

**Response Success (201 Created):**

```json
{
  "success": true,
  "message": "Department created successfully",
  "data": {
    "id": 1,
    "name": "Kardiologi",
    "code": "KARDIO",
    "description": "Departemen jantung dan pembuluh darah",
    "floor_location": "Lantai 3",
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z"
  }
}
```

**Response Error (400 Bad Request):**

```json
{
  "success": false,
  "message": "Failed to create department",
  "error": "code already exists"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/departments \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Kardiologi",
    "code": "KARDIO",
    "description": "Departemen jantung dan pembuluh darah",
    "floor_location": "Lantai 3"
  }'
```

---

### 5. Update Department

**Endpoint:** `PUT /api/v1/departments/:id`

**Description:** Admin mengupdate data department berdasarkan ID.

**Authentication:** Required (Admin/Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Department ID (integer)

**Request Body:**

```json
{
  "name": "Kardiologi Updated",
  "code": "KARDIO-NEW",
  "description": "Updated description",
  "floor_location": "Lantai 4"
}
```

**Field Rules:**

- All fields optional
- Validation same as create
- Code must be unique if changed

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Department updated successfully",
  "data": {
    "id": 1,
    "name": "Kardiologi Updated",
    "code": "KARDIO-NEW",
    "description": "Updated description",
    "floor_location": "Lantai 4",
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T12:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/departments/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Kardiologi Updated",
    "floor_location": "Lantai 4"
  }'
```

---

### 6. Soft Delete Department

**Endpoint:** `DELETE /api/v1/departments/:id`

**Description:** Admin menghapus department (soft delete). Data tetap ada dan bisa di-restore.

**Authentication:** Required (Admin/Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Department ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Department deleted successfully",
  "data": null
}
```

**Notes:**

- Data tidak dihapus dari database
- Bisa di-restore dengan endpoint restore
- Doctors dan Rooms tetap ada (tidak terpengaruh)

**⚠️ Business Rules:**

- Tidak bisa delete department yang masih memiliki active doctors
- Tidak bisa delete department yang masih memiliki active rooms

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/departments/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 7. Restore Department

**Endpoint:** `PATCH /api/v1/departments/:id/restore`

**Description:** Admin me-restore department yang sudah di-soft delete.

**Authentication:** Required (Admin/Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Department ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Department restored successfully",
  "data": null
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/departments/1/restore \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 8. Hard Delete Department

**Endpoint:** `DELETE /api/v1/departments/:id/hard-delete`

**Description:** Super Admin menghapus department secara permanen dari database.

**Authentication:** Required (Super Admin Only)

**Request Headers:**

```
Authorization: Bearer <super-admin-token>
```

**URL Parameters:**

- `id`: Department ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Department permanently deleted",
  "data": null
}
```

**⚠️ WARNING:**

- Data dihapus permanen dari database
- Tidak bisa di-restore
- Gunakan dengan hati-hati
- Pastikan tidak ada relasi dengan tabel lain

**⚠️ Business Rules:**

- Tidak bisa hard delete jika masih ada doctors terkait
- Tidak bisa hard delete jika masih ada rooms terkait

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/departments/1/hard-delete \
  -H "Authorization: Bearer SUPER_ADMIN_JWT_TOKEN"
```

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
  "message": "Department not found",
  "error": "department not found"
}
```

### 409 Conflict

```json
{
  "success": false,
  "message": "Cannot delete department",
  "error": "department has active doctors or rooms"
}
```

### 500 Internal Server Error

```json
{
  "success": false,
  "message": "Internal server error",
  "error": "database connection failed"
}
```

---

## Data Models

### Department Object

```json
{
  "id": 1,
  "name": "Kardiologi",
  "code": "KARDIO",
  "description": "Departemen jantung dan pembuluh darah",
  "floor_location": "Lantai 3",
  "created_at": "2024-01-19T10:00:00Z",
  "updated_at": "2024-01-19T10:00:00Z",
  "deleted_at": null,
  "doctors": [],
  "rooms": []
}
```

### Pagination Meta

```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 20,
  "total_pages": 2
}
```

---

## Business Rules

1. **Code Uniqueness**: Department code harus unik
2. **Name Required**: Department name wajib diisi
3. **Soft Delete Protection**: Tidak bisa delete department dengan active doctors/rooms
4. **Hard Delete Protection**: Super Admin only, pastikan tidak ada relasi
5. **Restore**: Department yang di-soft delete bisa di-restore
6. **Floor Location**: Opsional, untuk memudahkan navigasi
7. **Description**: Opsional, untuk informasi detail department

---

## Common Department Codes (Indonesia)

| Code     | Name                    | Description                |
| -------- | ----------------------- | -------------------------- |
| `KARDIO` | Kardiologi              | Jantung dan pembuluh darah |
| `NEURO`  | Neurologi               | Saraf                      |
| `PEDIA`  | Pediatri                | Anak                       |
| `OBGYN`  | Obstetri & Ginekologi   | Kandungan                  |
| `ORTHO`  | Orthopedi               | Tulang dan sendi           |
| `RADIO`  | Radiologi               | Pencitraan medis           |
| `LABKLI` | Laboratorium Klinik     | Analisis laboratorium      |
| `IGD`    | Instalasi Gawat Darurat | Emergency                  |
| `BEDAH`  | Bedah Umum              | General surgery            |
| `INTER`  | Penyakit Dalam          | Internal medicine          |

---

## Testing Examples

### Test 1: Create and Manage Department

```bash
# 1. Create Department
curl -X POST http://localhost:8080/api/v1/departments \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Kardiologi","code":"KARDIO","description":"Departemen jantung","floor_location":"Lantai 3"}'

# 2. List Departments
curl -X GET "http://localhost:8080/api/v1/departments?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 3. Get Department by ID
curl -X GET http://localhost:8080/api/v1/departments/1 \
  -H "Authorization: Bearer YOUR_TOKEN"

# 4. Update Department
curl -X PUT http://localhost:8080/api/v1/departments/1 \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"floor_location":"Lantai 4"}'

# 5. Soft Delete
curl -X DELETE http://localhost:8080/api/v1/departments/1 \
  -H "Authorization: Bearer ADMIN_TOKEN"

# 6. List Deleted
curl -X GET "http://localhost:8080/api/v1/departments/deleted" \
  -H "Authorization: Bearer ADMIN_TOKEN"

# 7. Restore
curl -X PATCH http://localhost:8080/api/v1/departments/1/restore \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

---

## Notes

- Semua timestamps menggunakan format ISO 8601 (UTC)
- Pagination maksimal 100 items per page
- Soft deleted departments muncul di endpoint `/departments/deleted`
- Hard delete hanya bisa dilakukan oleh Super Admin
- Department code case-sensitive (recommended: UPPERCASE)

---

**Last Updated:** 2024-01-19  
**API Version:** 1.0.0
