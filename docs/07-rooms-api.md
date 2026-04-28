# Rooms API Documentation

## Overview

API untuk manajemen data rooms (ruangan/kamar) dalam sistem rekam medis. Rooms digunakan untuk hospitalisasi dan pelayanan pasien rawat inap.

**Base URL:** `/api/v1/rooms`

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

| Endpoint                       | Patient | Doctor | Receptionist | Admin | Super Admin |
| ------------------------------ | ------- | ------ | ------------ | ----- | ----------- |
| GET /rooms                     | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /rooms/available           | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /rooms/occupied            | ❌      | ✅     | ✅           | ✅    | ✅          |
| GET /rooms/inactive            | ❌      | ❌     | ✅           | ✅    | ✅          |
| GET /rooms/deleted             | ❌      | ❌     | ❌           | ✅    | ✅          |
| GET /rooms/:id                 | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /rooms/type/:room_type     | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /rooms/department/:dept_id | ✅      | ✅     | ✅           | ✅    | ✅          |
| POST /rooms                    | ❌      | ❌     | ❌           | ✅    | ✅          |
| PUT /rooms/:id                 | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /rooms/:id/activate      | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /rooms/:id/deactivate    | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /rooms/:id/occupy        | ❌      | ❌     | ✅           | ✅    | ✅          |
| PATCH /rooms/:id/release       | ❌      | ❌     | ✅           | ✅    | ✅          |
| DELETE /rooms/:id              | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /rooms/:id/restore       | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /rooms/:id/hard-delete  | ❌      | ❌     | ❌           | ❌    | ✅          |

---

## Endpoints Summary

### Public Endpoints (All Authenticated)

| Method | Endpoint                     | Description             | Role Required                            |
| ------ | ---------------------------- | ----------------------- | ---------------------------------------- |
| GET    | `/rooms`                     | List all rooms          | All Authenticated                        |
| GET    | `/rooms/available`           | List available rooms    | All Authenticated                        |
| GET    | `/rooms/occupied`            | List occupied rooms     | Doctor, Receptionist, Admin, Super Admin |
| GET    | `/rooms/inactive`            | List inactive rooms     | Receptionist, Admin, Super Admin         |
| GET    | `/rooms/deleted`             | List deleted rooms      | Admin, Super Admin                       |
| GET    | `/rooms/:id`                 | Get room by ID          | All Authenticated                        |
| GET    | `/rooms/type/:room_type`     | Get rooms by type       | All Authenticated                        |
| GET    | `/rooms/department/:dept_id` | Get rooms by department | All Authenticated                        |

### Admin Endpoints

| Method | Endpoint                | Description          | Role Required                    |
| ------ | ----------------------- | -------------------- | -------------------------------- |
| POST   | `/rooms`                | Create room          | Admin, Super Admin               |
| PUT    | `/rooms/:id`            | Update room          | Admin, Super Admin               |
| PATCH  | `/rooms/:id/activate`   | Activate room        | Admin, Super Admin               |
| PATCH  | `/rooms/:id/deactivate` | Deactivate room      | Admin, Super Admin               |
| PATCH  | `/rooms/:id/occupy`     | Occupy bed           | Receptionist, Admin, Super Admin |
| PATCH  | `/rooms/:id/release`    | Release bed          | Receptionist, Admin, Super Admin |
| DELETE | `/rooms/:id`            | Soft delete room     | Admin, Super Admin               |
| PATCH  | `/rooms/:id/restore`    | Restore deleted room | Admin, Super Admin               |

### Super Admin Endpoints

| Method | Endpoint                 | Description             | Role Required |
| ------ | ------------------------ | ----------------------- | ------------- |
| DELETE | `/rooms/:id/hard-delete` | Permanently delete room | Super Admin   |

---

## Public Endpoints

### 1. List Rooms

**Endpoint:** `GET /api/v1/rooms`

**Description:** Mendapatkan daftar semua room dengan pagination, search, dan filter.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter            | Type    | Default     | Description                                                    |
| -------------------- | ------- | ----------- | -------------------------------------------------------------- |
| `page`               | integer | 1           | Halaman                                                        |
| `page_size`          | integer | 10          | Jumlah data per halaman (max: 100)                             |
| `search`             | string  | -           | Cari berdasarkan room_number                                   |
| `room_type`          | string  | -           | Filter by room type                                            |
| `department_id`      | integer | -           | Filter by department                                           |
| `is_active`          | boolean | -           | Filter by active status                                        |
| `has_available_beds` | boolean | -           | Filter rooms with available beds                               |
| `sort_by`            | string  | room_number | Sort field (room_number, room_type, price_per_day, created_at) |
| `sort_dir`           | string  | asc         | Sort direction (asc, desc)                                     |

**Example Request:**

```
GET /api/v1/rooms?page=1&page_size=10&room_type=vip&has_available_beds=true&sort_by=room_number&sort_dir=asc
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Rooms retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "room_number": "301-A",
        "room_type": "vip",
        "department_id": 1,
        "department": {
          "id": 1,
          "name": "Kardiologi",
          "code": "KARDIO"
        },
        "bed_capacity": 1,
        "available_beds": 0,
        "price_per_day": 1500000.0,
        "is_active": true,
        "occupancy_rate": 100.0,
        "created_at": "2024-01-19T10:00:00Z"
      },
      {
        "id": 2,
        "room_number": "302-A",
        "room_type": "vip",
        "department_id": 1,
        "department": {
          "id": 1,
          "name": "Kardiologi",
          "code": "KARDIO"
        },
        "bed_capacity": 1,
        "available_beds": 1,
        "price_per_day": 1500000.0,
        "is_active": true,
        "occupancy_rate": 0.0,
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
curl -X GET "http://localhost:8080/api/v1/rooms?page=1&page_size=10&room_type=vip" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 2. List Available Rooms

**Endpoint:** `GET /api/v1/rooms/available`

**Description:** Mendapatkan daftar room yang masih tersedia (ada bed kosong).

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter       | Type    | Description                     |
| --------------- | ------- | ------------------------------- |
| `room_type`     | string  | Filter by room type             |
| `department_id` | integer | Filter by department            |
| `min_beds`      | integer | Minimum available beds required |
| `max_price`     | decimal | Maximum price per day           |
| `page`          | integer | Page number                     |
| `page_size`     | integer | Items per page                  |

**Example Request:**

```
GET /api/v1/rooms/available?room_type=class_1&min_beds=1&max_price=500000
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Available rooms retrieved successfully",
  "data": {
    "total_available_rooms": 15,
    "total_available_beds": 28,
    "data": [
      {
        "id": 5,
        "room_number": "201-A",
        "room_type": "class_1",
        "department": {
          "id": 2,
          "name": "Penyakit Dalam"
        },
        "bed_capacity": 2,
        "available_beds": 2,
        "price_per_day": 450000.0,
        "is_active": true
      },
      {
        "id": 6,
        "room_number": "201-B",
        "room_type": "class_1",
        "department": {
          "id": 2,
          "name": "Penyakit Dalam"
        },
        "bed_capacity": 2,
        "available_beds": 1,
        "price_per_day": 450000.0,
        "is_active": true
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
curl -X GET "http://localhost:8080/api/v1/rooms/available?room_type=class_1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Use Case:**

- Patient registration untuk hospitalisasi
- Emergency room allocation
- Transfer patient planning

---

### 3. List Occupied Rooms

**Endpoint:** `GET /api/v1/rooms/occupied`

**Description:** Mendapatkan daftar room yang terisi penuh (tidak ada bed kosong).

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**
Same as List Available Rooms.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Occupied rooms retrieved successfully",
  "data": {
    "total_occupied_rooms": 8,
    "total_occupied_beds": 12,
    "data": [
      {
        "id": 1,
        "room_number": "301-A",
        "room_type": "vip",
        "department": {
          "id": 1,
          "name": "Kardiologi"
        },
        "bed_capacity": 1,
        "available_beds": 0,
        "occupancy_rate": 100.0,
        "current_patients": [
          {
            "patient_code": "P-2024-001",
            "patient_name": "John Doe",
            "admission_date": "2024-01-15T10:00:00Z"
          }
        ]
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
curl -X GET "http://localhost:8080/api/v1/rooms/occupied" \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN"
```

---

### 4. List Inactive Rooms

**Endpoint:** `GET /api/v1/rooms/inactive`

**Description:** Mendapatkan daftar room yang tidak aktif (maintenance, renovation, dll).

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**
Same as List Rooms.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Inactive rooms retrieved successfully",
  "data": {
    "data": [
      {
        "id": 20,
        "room_number": "401-A",
        "room_type": "icu",
        "department": {
          "id": 5,
          "name": "ICU"
        },
        "bed_capacity": 1,
        "is_active": false,
        "created_at": "2024-01-10T10:00:00Z",
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
curl -X GET "http://localhost:8080/api/v1/rooms/inactive" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 5. List Deleted Rooms

**Endpoint:** `GET /api/v1/rooms/deleted`

**Description:** Mendapatkan daftar room yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**
Same as List Rooms.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Deleted rooms retrieved successfully",
  "data": {
    "data": [
      {
        "id": 50,
        "room_number": "OLD-101",
        "room_type": "class_3",
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
curl -X GET "http://localhost:8080/api/v1/rooms/deleted" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 6. Get Room by ID

**Endpoint:** `GET /api/v1/rooms/:id`

**Description:** Mendapatkan detail room berdasarkan ID, termasuk current occupancy.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Room ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Room retrieved successfully",
  "data": {
    "id": 1,
    "room_number": "301-A",
    "room_type": "vip",
    "department_id": 1,
    "department": {
      "id": 1,
      "name": "Kardiologi",
      "code": "KARDIO",
      "floor_location": "Lantai 3"
    },
    "bed_capacity": 1,
    "available_beds": 0,
    "price_per_day": 1500000.0,
    "is_active": true,
    "occupancy_rate": 100.0,
    "is_available": false,
    "is_full": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z",
    "current_patients": [
      {
        "hospitalization_id": 1,
        "patient_code": "P-2024-001",
        "patient_name": "John Doe",
        "admission_date": "2024-01-15T10:00:00Z",
        "attending_doctor": "Dr. Jane Smith, Sp.JP"
      }
    ],
    "amenities": [
      "AC",
      "TV",
      "Private Bathroom",
      "WiFi",
      "Sofa Bed for Guardian"
    ]
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Room not found",
  "error": "room not found"
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/rooms/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 7. Get Rooms by Type

**Endpoint:** `GET /api/v1/rooms/type/:room_type`

**Description:** Mendapatkan daftar room berdasarkan tipe.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `room_type`: Room Type (string), values: vip, class_1, class_2, class_3, icu, emergency

**Query Parameters:**

| Parameter            | Type    | Description           |
| -------------------- | ------- | --------------------- |
| `is_active`          | boolean | Filter active only    |
| `has_available_beds` | boolean | Filter available only |
| `page`               | integer | Page number           |
| `page_size`          | integer | Items per page        |

**Example Request:**

```
GET /api/v1/rooms/type/vip?is_active=true&has_available_beds=true
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Rooms by type retrieved successfully",
  "data": {
    "room_type": "vip",
    "total_rooms": 10,
    "total_beds": 10,
    "available_beds": 3,
    "occupancy_rate": 70.0,
    "price_range": {
      "min": 1500000.0,
      "max": 2000000.0
    },
    "data": [
      {
        "id": 2,
        "room_number": "302-A",
        "department": {
          "id": 1,
          "name": "Kardiologi"
        },
        "bed_capacity": 1,
        "available_beds": 1,
        "price_per_day": 1500000.0,
        "is_active": true
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
curl -X GET "http://localhost:8080/api/v1/rooms/type/vip?has_available_beds=true" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 8. Get Rooms by Department

**Endpoint:** `GET /api/v1/rooms/department/:dept_id`

**Description:** Mendapatkan daftar room berdasarkan department.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `dept_id`: Department ID (integer)

**Query Parameters:**
Same as Get Rooms by Type.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Rooms by department retrieved successfully",
  "data": {
    "department": {
      "id": 1,
      "name": "Kardiologi",
      "code": "KARDIO",
      "floor_location": "Lantai 3"
    },
    "total_rooms": 15,
    "total_beds": 25,
    "available_beds": 8,
    "occupancy_rate": 68.0,
    "data": [
      {
        "id": 1,
        "room_number": "301-A",
        "room_type": "vip",
        "bed_capacity": 1,
        "available_beds": 0,
        "price_per_day": 1500000.0,
        "is_active": true
      },
      {
        "id": 5,
        "room_number": "201-A",
        "room_type": "class_1",
        "bed_capacity": 2,
        "available_beds": 1,
        "price_per_day": 450000.0,
        "is_active": true
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
curl -X GET "http://localhost:8080/api/v1/rooms/department/1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Admin Endpoints

### 9. Create Room

**Endpoint:** `POST /api/v1/rooms`

**Description:** Admin membuat room baru.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "room_number": "301-A",
  "room_type": "vip",
  "department_id": 1,
  "bed_capacity": 1,
  "available_beds": 1,
  "price_per_day": 1500000.0,
  "is_active": true
}
```

**Field Rules:**

- `room_number`: required, unique, max 20 characters, indexed
- `room_type`: required, enum (vip, class_1, class_2, class_3, icu, emergency)
- `department_id`: optional, FK to departments table
- `bed_capacity`: required, min 1
- `available_beds`: optional, defaults to bed_capacity
- `price_per_day`: optional, decimal(10,2)
- `is_active`: optional, boolean (default: true)

**Response Success (201 Created):**

```json
{
  "success": true,
  "message": "Room created successfully",
  "data": {
    "id": 1,
    "room_number": "301-A",
    "room_type": "vip",
    "department_id": 1,
    "bed_capacity": 1,
    "available_beds": 1,
    "price_per_day": 1500000.0,
    "is_active": true,
    "occupancy_rate": 0.0,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z"
  }
}
```

**Response Error (400 Bad Request):**

```json
{
  "success": false,
  "message": "Failed to create room",
  "error": "room_number already exists"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/rooms \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "room_number": "301-A",
    "room_type": "vip",
    "department_id": 1,
    "bed_capacity": 1,
    "price_per_day": 1500000
  }'
```

**Notes:**

- Available beds automatically set to bed_capacity if not specified
- Room number must follow hospital naming convention
- Price per day optional (untuk internal use/free rooms)

---

### 10. Update Room

**Endpoint:** `PUT /api/v1/rooms/:id`

**Description:** Admin mengupdate data room berdasarkan ID.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Room ID (integer)

**Request Body:**

```json
{
  "room_number": "301-A-Updated",
  "room_type": "vip",
  "department_id": 2,
  "bed_capacity": 2,
  "price_per_day": 1800000.0,
  "is_active": true
}
```

**Field Rules:**

- All fields optional
- Validation same as create
- Cannot reduce bed_capacity below currently occupied
- Cannot change room_number if room is occupied

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Room updated successfully",
  "data": {
    "id": 1,
    "room_number": "301-A-Updated",
    "room_type": "vip",
    "department_id": 2,
    "bed_capacity": 2,
    "available_beds": 2,
    "price_per_day": 1800000.0,
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T15:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/rooms/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "price_per_day": 1800000
  }'
```

---

### 11. Activate Room

**Endpoint:** `PATCH /api/v1/rooms/:id/activate`

**Description:** Admin mengaktifkan room yang sedang inactive.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Room ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Room activated successfully",
  "data": null
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/rooms/1/activate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

**Use Case:**

- Room selesai maintenance
- Room selesai renovation
- Room ready untuk digunakan kembali

---

### 12. Deactivate Room

**Endpoint:** `PATCH /api/v1/rooms/:id/deactivate`

**Description:** Admin menonaktifkan room (maintenance, renovation).

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Room ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Room deactivated successfully",
  "data": null
}
```

**Notes:**

- Room yang deactivated tidak bisa digunakan untuk hospitalisasi baru
- Patient yang sudah ada tetap bisa stay sampai discharge
- Available beds set to 0

**⚠️ Business Rules:**

- Harus transfer patients terlebih dahulu sebelum deactivate
- Atau tunggu sampai semua patients discharge

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/rooms/1/deactivate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 13. Occupy Bed

**Endpoint:** `PATCH /api/v1/rooms/:id/occupy`

**Description:** Mengurangi available beds (saat patient admission).

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Room ID (integer)

**Request Body:**

```json
{
  "beds_count": 1
}
```

**Field Rules:**

- `beds_count`: optional, default 1, min 1

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Bed occupied successfully",
  "data": {
    "id": 5,
    "room_number": "201-A",
    "bed_capacity": 2,
    "available_beds": 1,
    "occupancy_rate": 50.0,
    "is_available": true,
    "is_full": false
  }
}
```

**Response Error (400 Bad Request):**

```json
{
  "success": false,
  "message": "Failed to occupy bed",
  "error": "no available beds in this room"
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/rooms/5/occupy \
  -H "Authorization: Bearer RECEPTIONIST_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "beds_count": 1
  }'
```

**Notes:**

- Automatically called saat patient admission
- Checks if room has available beds
- Updates occupancy rate

---

### 14. Release Bed

**Endpoint:** `PATCH /api/v1/rooms/:id/release`

**Description:** Menambah available beds (saat patient discharge).

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Room ID (integer)

**Request Body:**

```json
{
  "beds_count": 1
}
```

**Field Rules:**

- `beds_count`: optional, default 1, min 1

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Bed released successfully",
  "data": {
    "id": 5,
    "room_number": "201-A",
    "bed_capacity": 2,
    "available_beds": 2,
    "occupancy_rate": 0.0,
    "is_available": true,
    "is_full": false
  }
}
```

**Response Error (400 Bad Request):**

```json
{
  "success": false,
  "message": "Failed to release bed",
  "error": "available beds cannot exceed bed capacity"
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/rooms/5/release \
  -H "Authorization: Bearer RECEPTIONIST_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "beds_count": 1
  }'
```

**Notes:**

- Automatically called saat patient discharge
- Cannot exceed bed capacity
- Updates occupancy rate

---

### 15. Soft Delete Room

**Endpoint:** `DELETE /api/v1/rooms/:id`

**Description:** Admin menghapus room (soft delete). Data tetap ada dan bisa di-restore.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Room ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Room deleted successfully",
  "data": null
}
```

**Notes:**

- Room yang dihapus tidak muncul di list normal
- Hospitalization records tetap tersimpan
- Bisa di-restore dengan endpoint restore

**⚠️ Business Rules:**

- Tidak bisa delete room yang masih occupied
- Harus tunggu semua patients discharge
- Atau transfer patients ke room lain

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/rooms/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 16. Restore Room

**Endpoint:** `PATCH /api/v1/rooms/:id/restore`

**Description:** Admin me-restore room yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Room ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Room restored successfully",
  "data": null
}
```

**Notes:**

- Room di-restore dengan status inactive
- Perlu activate manual jika ingin langsung active

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/rooms/1/restore \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

## Super Admin Endpoints

### 17. Hard Delete Room

**Endpoint:** `DELETE /api/v1/rooms/:id/hard-delete`

**Description:** Super Admin menghapus room secara permanen dari database.

**Authentication:** Required (Super Admin Only)

**Request Headers:**

```
Authorization: Bearer <super-admin-token>
```

**URL Parameters:**

- `id`: Room ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Room permanently deleted",
  "data": null
}
```

**⚠️ WARNING:**

- Data dihapus permanen dari database
- Tidak bisa di-restore
- Gunakan dengan sangat hati-hati

**⚠️ Business Rules:**

- Tidak bisa hard delete jika masih ada hospitalization records
- Must archive all historical data first
- Requires special approval

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/rooms/1/hard-delete \
  -H "Authorization: Bearer SUPER_ADMIN_JWT_TOKEN"
```

---

## Database Model

### Table: rooms

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier untuk room |
| room_number | VARCHAR(20) | UNIQUE, NOT NULL, INDEX | Nomor ruangan (e.g., 101, 202-A) |
| room_type | VARCHAR(20) | NOT NULL, INDEX | Tipe ruangan (Standard, VIP, ICU, dll) |
| department_id | BIGINT | FOREIGN KEY (departments.id), INDEX | Department tempat ruangan berada |
| bed_capacity | INT | NOT NULL, DEFAULT 1 | Jumlah maksimal tempat tidur |
| available_beds | INT | NOT NULL, DEFAULT 1 | Jumlah tempat tidur yang tersedia |
| price_per_day | DECIMAL(10,2) | DEFAULT 0 | Harga per hari (dalam IDR) |
| is_active | BOOLEAN | DEFAULT true, INDEX | Status aktif ruangan |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan record |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update terakhir |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

**Indexes:**
- Primary Key: id
- Unique Index: room_number
- Foreign Key: department_id
- Regular Index: room_type, is_active, deleted_at

**Relationships:**
- Belongs To Department (many-to-one)
- Has Many Hospitalizations (one-to-many)

**Notes:**
- room_type bisa hardcoded atau reference ke master data
- available_beds otomatis update saat ada admission/discharge
- price_per_day bisa berbeda untuk setiap room type dan department
- is_active untuk maintenance atau permanent closure
- bed_capacity harus >= 1, available_beds <= bed_capacity

---

## Error Responses

### 400 Bad Request

```json
{
  "success": false,
  "message": "Validation error",
  "error": "room_number already exists"
}
```

### 404 Not Found

```json
{
  "success": false,
  "message": "Room not found",
  "error": "room not found"
}
```

### 409 Conflict

```json
{
  "success": false,
  "message": "Cannot delete room",
  "error": "room is currently occupied"
}
```

### 422 Unprocessable Entity

```json
{
  "success": false,
  "message": "Cannot occupy bed",
  "error": "no available beds in this room"
}
```

---

## Data Models

### Room Object (Full)

```json
{
  "id": 1,
  "room_number": "301-A",
  "room_type": "vip",
  "department_id": 1,
  "department": {
    "id": 1,
    "name": "Kardiologi",
    "code": "KARDIO",
    "floor_location": "Lantai 3"
  },
  "bed_capacity": 1,
  "available_beds": 0,
  "price_per_day": 1500000.0,
  "is_active": true,
  "occupancy_rate": 100.0,
  "is_available": false,
  "is_full": true,
  "created_at": "2024-01-19T10:00:00Z",
  "updated_at": "2024-01-19T10:00:00Z",
  "deleted_at": null
}
```

### Room Types & Pricing

```json
{
  "room_types": [
    {
      "type": "vip",
      "name": "VIP Room",
      "typical_capacity": 1,
      "price_range": "1,500,000 - 2,500,000",
      "amenities": ["AC", "TV", "Private Bathroom", "WiFi", "Sofa Bed"]
    },
    {
      "type": "class_1",
      "name": "Class 1",
      "typical_capacity": 2,
      "price_range": "400,000 - 600,000",
      "amenities": ["AC", "TV", "Shared Bathroom"]
    },
    {
      "type": "class_2",
      "name": "Class 2",
      "typical_capacity": 3,
      "price_range": "200,000 - 350,000",
      "amenities": ["Fan", "Shared Bathroom"]
    },
    {
      "type": "class_3",
      "name": "Class 3",
      "typical_capacity": 4-6,
      "price_range": "100,000 - 180,000",
      "amenities": ["Basic", "Shared Bathroom"]
    },
    {
      "type": "icu",
      "name": "ICU",
      "typical_capacity": 1,
      "price_range": "2,000,000 - 5,000,000",
      "amenities": ["Full ICU Equipment", "24/7 Monitoring"]
    },
    {
      "type": "emergency",
      "name": "Emergency Room",
      "typical_capacity": 1,
      "price_range": "Varies",
      "amenities": ["Emergency Equipment"]
    }
  ]
}
```

---

## Business Rules

1. **Room Number Uniqueness**: Room number harus unik
2. **Bed Management**: Available beds tidak boleh negative atau exceed capacity
3. **Occupancy Calculation**: (occupied_beds / total_beds) \* 100
4. **Active Status**: Hanya room active yang bisa untuk admission
5. **Price Management**: Price per day optional (free rooms)
6. **Department Assignment**: Room sebaiknya assign ke department
7. **Soft Delete Protection**: Cannot delete occupied rooms
8. **Historical Data**: Hospitalization records tetap accessible
9. **Bed Allocation**: First-come-first-served atau by priority
10. **Transfer Protocol**: Must follow transfer procedures

---

## Common Use Cases

### Use Case 1: Patient Admission

```bash
# 1. Find available rooms
GET /api/v1/rooms/available?room_type=class_1

# 2. Check room details
GET /api/v1/rooms/5

# 3. Occupy bed (automatic via hospitalization creation)
PATCH /api/v1/rooms/5/occupy

# 4. Create hospitalization record
POST /api/v1/hospitalizations
```

### Use Case 2: Patient Discharge

```bash
# 1. Update hospitalization record
PATCH /api/v1/hospitalizations/1/discharge

# 2. Release bed (automatic)
PATCH /api/v1/rooms/5/release
```

### Use Case 3: Room Maintenance

```bash
# 1. Check if room is empty
GET /api/v1/rooms/10

# 2. Transfer patients if needed
POST /api/v1/hospitalizations/bulk-transfer

# 3. Deactivate room
PATCH /api/v1/rooms/10/deactivate

# 4. After maintenance, activate
PATCH /api/v1/rooms/10/activate
```

---

## Testing Examples

### Test 1: Complete Room Management Flow

```bash
# 1. Create Room
curl -X POST http://localhost:8080/api/v1/rooms \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "room_number": "TEST-101",
    "room_type": "class_1",
    "bed_capacity": 2,
    "price_per_day": 450000
  }'

# 2. List Available Rooms
curl -X GET "http://localhost:8080/api/v1/rooms/available" \
  -H "Authorization: Bearer PATIENT_TOKEN"

# 3. Occupy Bed
curl -X PATCH http://localhost:8080/api/v1/rooms/1/occupy \
  -H "Authorization: Bearer RECEPTIONIST_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"beds_count": 1}'

# 4. Release Bed
curl -X PATCH http://localhost:8080/api/v1/rooms/1/release \
  -H "Authorization: Bearer RECEPTIONIST_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"beds_count": 1}'

# 5. Deactivate Room
curl -X PATCH http://localhost:8080/api/v1/rooms/1/deactivate \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

---

## Notes

- Room numbering: [Floor][Wing]-[Number] (e.g., 301-A, 201-B)
- Occupancy rate calculated real-time
- Price per day in IDR (Indonesian Rupiah)
- ICU and Emergency rooms usually single bed
- Class 3 untuk BPJS/free treatment
- VIP rooms untuk private patients
- Bed management automatic via hospitalization
- Support for bed transfer between rooms

---

**Last Updated:** 2024-01-19  
**API Version:** 1.0.0
