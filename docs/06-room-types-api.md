# Room Types API Documentation

## Overview

API untuk manajemen master data tipe ruangan dalam sistem rekam medis. Room Types adalah master data yang dapat direferensikan oleh Rooms untuk menstandarisasi tipe ruangan (menghindari duplikasi data seperti "VIP" vs "vip" atau "Class 1").

**Base URL:** `/api/v1/room-types`

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
| GET /room-types                    | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /room-types/:id                | ✅      | ✅     | ✅           | ✅    | ✅          |
| POST /room-types                   | ❌      | ❌     | ❌           | ✅    | ✅          |
| PUT /room-types/:id                | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /room-types/:id/activate     | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /room-types/:id/deactivate   | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /room-types/:id             | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /room-types/:id/hard-delete | ❌      | ❌     | ❌           | ❌    | ✅          |

---

## Endpoints Summary

### Public Endpoints (All Authenticated)

| Method | Endpoint          | Description         | Role Required     |
| ------ | ----------------- | ------------------- | ----------------- |
| GET    | `/room-types`     | List all room types | All Authenticated |
| GET    | `/room-types/:id` | Get room type by ID | All Authenticated |

### Admin Endpoints

| Method | Endpoint                     | Description           | Role Required      |
| ------ | ---------------------------- | --------------------- | ------------------ |
| POST   | `/room-types`                | Create room type      | Admin, Super Admin |
| PUT    | `/room-types/:id`            | Update room type      | Admin, Super Admin |
| PATCH  | `/room-types/:id/activate`   | Activate room type    | Admin, Super Admin |
| PATCH  | `/room-types/:id/deactivate` | Deactivate room type  | Admin, Super Admin |
| DELETE | `/room-types/:id`            | Soft delete room type | Admin, Super Admin |

### Super Admin Endpoints

| Method | Endpoint                      | Description                  | Role Required |
| ------ | ----------------------------- | ---------------------------- | ------------- |
| DELETE | `/room-types/:id/hard-delete` | Permanently delete room type | Super Admin   |

---

## Public Endpoints

### 1. List Room Types

**Endpoint:** `GET /api/v1/room-types`

**Description:** Mendapatkan daftar semua room types dengan pagination dan filter.

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
| `search`    | string  | -       | Cari berdasarkan nama tipe ruangan |
| `is_active` | boolean | -       | Filter by active status            |
| `sort_by`   | string  | name    | Sort field (name, created_at)      |
| `sort_dir`  | string  | asc     | Sort direction (asc, desc)         |

**Example Request:**

```
GET /api/v1/room-types?page=1&page_size=10&is_active=true&sort_by=name&sort_dir=asc
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Room types retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "VIP",
        "description": "Ruangan VIP dengan fasilitas premium",
        "code": "VIP",
        "is_active": true,
        "created_at": "2024-01-19T10:00:00Z"
      },
      {
        "id": 2,
        "name": "Class 1",
        "description": "Ruangan kelas 1 dengan fasilitas standar",
        "code": "CLS1",
        "is_active": true,
        "created_at": "2024-01-19T10:05:00Z"
      },
      {
        "id": 3,
        "name": "Class 2",
        "description": "Ruangan kelas 2 dengan fasilitas dasar",
        "code": "CLS2",
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
curl -X GET "http://localhost:8080/api/v1/room-types?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 2. Get Room Type by ID

**Endpoint:** `GET /api/v1/room-types/:id`

**Description:** Mendapatkan detail room type berdasarkan ID.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Room Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Room type retrieved successfully",
  "data": {
    "id": 1,
    "name": "VIP",
    "description": "Ruangan VIP dengan fasilitas premium",
    "code": "VIP",
    "is_active": true,
    "total_rooms": 10,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Room type not found",
  "error": "room type not found"
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/room-types/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Admin Endpoints

### 3. Create Room Type

**Endpoint:** `POST /api/v1/room-types`

**Description:** Admin membuat tipe ruangan baru.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "name": "VIP",
  "description": "Ruangan VIP dengan fasilitas premium",
  "code": "VIP",
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
  "message": "Room type created successfully",
  "data": {
    "id": 1,
    "name": "VIP",
    "description": "Ruangan VIP dengan fasilitas premium",
    "code": "VIP",
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
  "message": "Failed to create room type",
  "error": "name already exists"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/room-types \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "VIP",
    "description": "Ruangan VIP dengan fasilitas premium",
    "code": "VIP",
    "is_active": true
  }'
```

---

### 4. Update Room Type

**Endpoint:** `PUT /api/v1/room-types/:id`

**Description:** Admin mengupdate data room type berdasarkan ID.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Room Type ID (integer)

**Request Body:**

```json
{
  "name": "VIP - Updated",
  "description": "Ruangan VIP dengan fasilitas premium dan akses WiFi",
  "code": "VIP",
  "is_active": true
}
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Room type updated successfully",
  "data": {
    "id": 1,
    "name": "VIP - Updated",
    "description": "Ruangan VIP dengan fasilitas premium dan akses WiFi",
    "code": "VIP",
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
  "message": "Room type not found",
  "error": "room type not found"
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/room-types/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "VIP - Updated",
    "description": "Updated description",
    "code": "VIP",
    "is_active": true
  }'
```

---

### 5. Activate Room Type

**Endpoint:** `PATCH /api/v1/room-types/:id/activate`

**Description:** Admin mengaktifkan room type yang tidak aktif.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Room Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Room type activated successfully",
  "data": {
    "id": 1,
    "name": "VIP",
    "code": "VIP",
    "is_active": true,
    "updated_at": "2024-01-19T16:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/room-types/1/activate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 6. Deactivate Room Type

**Endpoint:** `PATCH /api/v1/room-types/:id/deactivate`

**Description:** Admin menonaktifkan room type.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Room Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Room type deactivated successfully",
  "data": {
    "id": 1,
    "name": "VIP",
    "code": "VIP",
    "is_active": false,
    "updated_at": "2024-01-19T16:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/room-types/1/deactivate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 7. Delete Room Type

**Endpoint:** `DELETE /api/v1/room-types/:id`

**Description:** Admin melakukan soft delete room type.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Room Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Room type deleted successfully",
  "data": {
    "id": 1,
    "name": "VIP",
    "is_active": false,
    "deleted_at": "2024-01-19T15:30:00Z"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Room type not found",
  "error": "room type not found"
}
```

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/room-types/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

## Super Admin Endpoints

### 8. Hard Delete Room Type

**Endpoint:** `DELETE /api/v1/room-types/:id/hard-delete`

**Description:** Super Admin melakukan permanent delete room type.

**Authentication:** Required (Super Admin)

**Request Headers:**

```
Authorization: Bearer <super-admin-token>
```

**URL Parameters:**

- `id`: Room Type ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Room type permanently deleted",
  "data": null
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Room type not found",
  "error": "room type not found"
}
```

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/room-types/1/hard-delete \
  -H "Authorization: Bearer SUPER_ADMIN_JWT_TOKEN"
```

---

## Database Model

### Table: room_types

| Field       | Type         | Constraints                                           | Description                                    |
| ----------- | ------------ | ----------------------------------------------------- | ---------------------------------------------- |
| id          | BIGINT       | PRIMARY KEY, AUTO_INCREMENT                           | Unique identifier                              |
| name        | VARCHAR(100) | NOT NULL, UNIQUE, INDEX                               | Nama tipe ruangan (e.g., Standard, VIP, ICU)   |
| code        | VARCHAR(20)  | UNIQUE, INDEX                                         | Kode tipe ruangan (e.g., STD, VIP, ICU)        |
| description | TEXT         | NULLABLE                                              | Deskripsi detail tipe ruangan                  |
| price_tier  | VARCHAR(20)  | NULLABLE                                              | Tier harga (Budget, Standard, Premium, Luxury) |
| is_active   | BOOLEAN      | NOT NULL, DEFAULT true, INDEX                         | Status aktif                                   |
| created_at  | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP                             | Waktu pembuatan record                         |
| updated_at  | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir                          |
| deleted_at  | TIMESTAMP    | INDEX, NULLABLE                                       | Soft delete timestamp                          |

**Indexes:**

- Primary Key: id
- Unique Index: name, code
- Regular Index: is_active, deleted_at

**Relationships:**

- Has Many Rooms (one-to-many) - rooms reference room type

**Notes:**

- Master data untuk standardisasi tipe ruangan
- Code sebaiknya uppercase (STD, VIP, ICU, dll)
- price_tier membantu mengelompokkan harga ruangan
- is_active untuk hide inactive types dari list

---

## Error Responses

### Common Error Codes

| Status Code | Message               | Description                              |
| ----------- | --------------------- | ---------------------------------------- |
| 400         | Bad Request           | Invalid request body or query parameters |
| 401         | Unauthorized          | Missing or invalid JWT token             |
| 403         | Forbidden             | User does not have permission            |
| 404         | Not Found             | Resource not found                       |
| 409         | Conflict              | Duplicate data (e.g., room type exists)  |
| 500         | Internal Server Error | Unexpected server error                  |

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

- Tipe ruangan digunakan untuk menstandarisasi tipe ruangan pada data rooms
- Nama tipe ruangan bersifat unique dan case-sensitive
- Soft delete hanya mengubah status, data tetap tersimpan
- Hard delete hanya bisa dilakukan oleh Super Admin dan bersifat permanent
