# Doctors API Documentation

## Overview

API untuk manajemen data doctors (dokter) dalam sistem rekam medis. Doctors adalah tenaga medis professional yang memberikan layanan kesehatan kepada pasien di rumah sakit.

**Base URL:** `/api/v1/doctors`

---

## Table of Contents

- [Authentication](#authentication)
- [Authorization](#authorization)
- [Endpoints Summary](#endpoints-summary)
- [Self-Owned Endpoints](#self-owned-endpoints)
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

| Endpoint                        | Doctor (Own) | Patient | Receptionist | Admin | Super Admin |
| ------------------------------- | ------------ | ------- | ------------ | ----- | ----------- |
| GET /doctors/me                 | ✅           | ❌      | ❌           | ❌    | ❌          |
| PUT /doctors/me                 | ✅           | ❌      | ❌           | ❌    | ❌          |
| GET /doctors                    | ✅           | ✅      | ✅           | ✅    | ✅          |
| GET /doctors/deleted            | ❌           | ❌      | ❌           | ✅    | ✅          |
| GET /doctors/:id                | ✅           | ✅      | ✅           | ✅    | ✅          |
| POST /doctors                   | ❌           | ❌      | ❌           | ✅    | ✅          |
| PUT /doctors/:id                | ❌           | ❌      | ❌           | ✅    | ✅          |
| PATCH /doctors/:id/activate     | ❌           | ❌      | ❌           | ✅    | ✅          |
| PATCH /doctors/:id/deactivate   | ❌           | ❌      | ❌           | ✅    | ✅          |
| DELETE /doctors/:id             | ❌           | ❌      | ❌           | ✅    | ✅          |
| PATCH /doctors/:id/restore      | ❌           | ❌      | ❌           | ✅    | ✅          |
| DELETE /doctors/:id/hard-delete | ❌           | ❌      | ❌           | ❌    | ✅          |

---

## Endpoints Summary

### Self-Owned Endpoints (`/me`)

| Method | Endpoint      | Description           | Auth   |
| ------ | ------------- | --------------------- | ------ |
| GET    | `/doctors/me` | Get my doctor data    | Doctor |
| PUT    | `/doctors/me` | Update my doctor data | Doctor |

### Public Endpoints (All Authenticated)

| Method | Endpoint       | Description         | Role Required     |
| ------ | -------------- | ------------------- | ----------------- |
| GET    | `/doctors`     | List active doctors | All Authenticated |
| GET    | `/doctors/:id` | Get doctor by ID    | All Authenticated |

### Admin Endpoints (Management)

| Method | Endpoint                  | Description            | Role Required      |
| ------ | ------------------------- | ---------------------- | ------------------ |
| GET    | `/doctors/deleted`        | List deleted doctors   | Admin, Super Admin |
| POST   | `/doctors`                | Create doctor          | Admin, Super Admin |
| PUT    | `/doctors/:id`            | Update doctor          | Admin, Super Admin |
| PATCH  | `/doctors/:id/activate`   | Activate doctor        | Admin, Super Admin |
| PATCH  | `/doctors/:id/deactivate` | Deactivate doctor      | Admin, Super Admin |
| DELETE | `/doctors/:id`            | Soft delete doctor     | Admin, Super Admin |
| PATCH  | `/doctors/:id/restore`    | Restore deleted doctor | Admin, Super Admin |

### Super Admin Endpoints

| Method | Endpoint                   | Description               | Role Required |
| ------ | -------------------------- | ------------------------- | ------------- |
| DELETE | `/doctors/:id/hard-delete` | Permanently delete doctor | Super Admin   |

---

## Self-Owned Endpoints

### 1. Get My Doctor Profile

**Endpoint:** `GET /api/v1/doctors/me`

**Description:** Doctor mendapatkan data profile diri sendiri.

**Authentication:** Required (Doctor Role)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor profile retrieved successfully",
  "data": {
    "id": 1,
    "user_id": 2,
    "employee_id": "DOC001",
    "full_name": "Dr. John Smith",
    "license_number": "LIC123456",
    "phone": "081234567890",
    "email": "drsmith@hospital.com",
    "department_id": 1,
    "department": {
      "id": 1,
      "name": "Cardiology Department",
      "code": "CARD"
    },
    "doctor_specialization_id": 1,
    "doctor_specialization": {
      "id": 1,
      "name": "Cardiology",
      "code": "CARDIO"
    },
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Doctor profile not found",
  "error": "doctor profile not created yet"
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/doctors/me \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN"
```

**Notes:**

- Doctor harus sudah memiliki doctor record (dibuat oleh admin)
- Hanya bisa melihat data diri sendiri

---

### 2. Update My Doctor Profile

**Endpoint:** `PUT /api/v1/doctors/me`

**Description:** Doctor mengupdate data profile diri sendiri (terbatas).

**Authentication:** Required (Doctor Role)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "phone": "081234567899",
  "email": "newdremail@hospital.com"
}
```

**Field Rules:**

- `phone`: optional, max 15 characters
- `email`: optional, valid email format
- ❌ **Cannot update:** `employee_id`, `full_name`, `doctor_specialization_id`, `license_number`, `department_id`, `is_active`

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor profile updated successfully",
  "data": {
    "id": 1,
    "user_id": 2,
    "employee_id": "DOC001",
    "full_name": "Dr. John Smith",
    "license_number": "LIC123456",
    "phone": "081234567899",
    "email": "newdremail@hospital.com",
    "department_id": 1,
    "doctor_specialization_id": 1,
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T14:30:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/doctors/me \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "081234567899",
    "email": "newdremail@hospital.com"
  }'
```

---

## Public Endpoints

### 3. List Doctors

**Endpoint:** `GET /api/v1/doctors`

**Description:** Mendapatkan daftar doctor aktif dengan pagination, search, dan filter.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter                  | Type    | Default    | Description                                     |
| -------------------------- | ------- | ---------- | ----------------------------------------------- |
| `page`                     | integer | 1          | Halaman                                         |
| `page_size`                | integer | 10         | Jumlah data per halaman (max: 100)              |
| `search`                   | string  | -          | Cari berdasarkan name, employee_id              |
| `doctor_specialization_id` | integer | -          | Filter by doctor specialization ID              |
| `department_id`            | integer | -          | Filter by department                            |
| `is_active`                | boolean | true       | Filter by active status                         |
| `sort_by`                  | string  | created_at | Sort field (created_at, full_name, employee_id) |
| `sort_dir`                 | string  | desc       | Sort direction (asc, desc)                      |

**Example Request:**

```
GET /api/v1/doctors?page=1&page_size=10&search=john&doctor_specialization_id=1&is_active=true&sort_by=full_name&sort_dir=asc
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctors retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "employee_id": "DOC001",
        "full_name": "Dr. John Smith",
        "license_number": "LIC123456",
        "phone": "081234567890",
        "email": "drsmith@hospital.com",
        "department": {
          "id": 1,
          "name": "Cardiology Department",
          "code": "CARD"
        },
        "doctor_specialization": {
          "id": 1,
          "name": "Cardiology",
          "code": "CARDIO"
        },
        "is_active": true,
        "created_at": "2024-01-19T10:00:00Z"
      },
      {
        "id": 2,
        "employee_id": "DOC002",
        "full_name": "Dr. John Doe",
        "license_number": "LIC789012",
        "phone": "081234567891",
        "email": "drjohn@hospital.com",
        "department": {
          "id": 1,
          "name": "Cardiology Department",
          "code": "CARD"
        },
        "doctor_specialization": {
          "id": 1,
          "name": "Cardiology",
          "code": "CARDIO"
        },
        "is_active": true,
        "created_at": "2024-01-19T11:00:00Z"
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
curl -X GET "http://localhost:8080/api/v1/doctors?page=1&page_size=10&doctor_specialization_id=1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 4. Get Doctor by ID

**Endpoint:** `GET /api/v1/doctors/:id`

**Description:** Mendapatkan detail doctor berdasarkan ID.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Doctor ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor retrieved successfully",
  "data": {
    "id": 1,
    "user_id": 2,
    "employee_id": "DOC001",
    "full_name": "Dr. John Smith",
    "license_number": "LIC123456",
    "phone": "081234567890",
    "email": "drsmith@hospital.com",
    "department_id": 1,
    "department": {
      "id": 1,
      "name": "Cardiology Department",
      "code": "CARD",
      "description": "Department for heart and cardiovascular diseases"
    },
    "doctor_specialization_id": 1,
    "doctor_specialization": {
      "id": 1,
      "name": "Cardiology",
      "code": "CARDIO",
      "description": "Heart and cardiovascular diseases"
    },
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z",
    "total_patients": 45,
    "total_appointments": 120,
    "completed_appointments": 100
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Doctor not found",
  "error": "doctor not found"
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/doctors/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Admin Endpoints

### 5. List Deleted Doctors

**Endpoint:** `GET /api/v1/doctors/deleted`

**Description:** Mendapatkan daftar doctor yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**
Same as List Doctors endpoint.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Deleted doctors retrieved successfully",
  "data": {
    "data": [
      {
        "id": 10,
        "employee_id": "DOC010",
        "full_name": "Dr. Deleted Doctor",
        "doctor_specialization_id": 5,
        "phone": "081234567899",
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
curl -X GET "http://localhost:8080/api/v1/doctors/deleted?page=1&page_size=10" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 6. Create Doctor

**Endpoint:** `POST /api/v1/doctors`

**Description:** Admin membuat doctor record baru.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "user_id": 3,
  "employee_id": "DOC002",
  "full_name": "Dr. Jane Doe",
  "doctor_specialization_id": 2,
  "license_number": "LIC789012",
  "phone": "081234567891",
  "email": "drjane@hospital.com",
  "department_id": 2,
  "is_active": true
}
```

**Field Rules:**

- `user_id`: optional, FK to users table (if doctor has user account)
- `employee_id`: required, unique, max 50 characters, indexed
- `full_name`: required, max 100 characters
- `doctor_specialization_id`: required, FK to doctor_specializations table
- `license_number`: required, unique, max 50 characters
- `phone`: optional, max 15 characters
- `email`: optional, valid email format
- `department_id`: required, FK to departments table
- `is_active`: optional, boolean, default true

**Response Success (201 Created):**

```json
{
  "success": true,
  "message": "Doctor created successfully",
  "data": {
    "id": 2,
    "user_id": 3,
    "employee_id": "DOC002",
    "full_name": "Dr. Jane Doe",
    "license_number": "LIC789012",
    "phone": "081234567891",
    "email": "drjane@hospital.com",
    "department_id": 2,
    "department": {
      "id": 2,
      "name": "Neurology Department",
      "code": "NEUR"
    },
    "doctor_specialization_id": 2,
    "doctor_specialization": {
      "id": 2,
      "name": "Neurology",
      "code": "NEURO"
    },
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
  "message": "Failed to create doctor",
  "error": "employee_id already exists"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/doctors \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "DOC002",
    "full_name": "Dr. Jane Doe",
    "doctor_specialization_id": 2,
    "license_number": "LIC789012",
    "phone": "081234567891",
    "email": "drjane@hospital.com",
    "department_id": 2
  }'
```

**Notes:**

- Employee ID auto-generated jika tidak diisi
- Format: DOC-NNNN (DOC-0001, DOC-0002, dst)
- License number harus unique dan valid

---

### 7. Update Doctor

**Endpoint:** `PUT /api/v1/doctors/:id`

**Description:** Admin mengupdate data doctor berdasarkan ID.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Doctor ID (integer)

**Request Body:**

```json
{
  "full_name": "Dr. Jane Doe Smith",
  "doctor_specialization_id": 3,
  "phone": "081234567899",
  "email": "drjane_new@hospital.com",
  "department_id": 3,
  "is_active": true
}
```

**Field Rules:**

- All fields optional
- Validation same as create
- Employee ID cannot be changed
- License number can be updated (with caution)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor updated successfully",
  "data": {
    "id": 2,
    "employee_id": "DOC002",
    "full_name": "Dr. Jane Doe Smith",
    "license_number": "LIC789012",
    "phone": "081234567899",
    "email": "drjane_new@hospital.com",
    "department_id": 3,
    "doctor_specialization_id": 3,
    "doctor_specialization": {
      "id": 3,
      "name": "Psychiatric",
      "code": "PSY"
    },
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T15:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/doctors/2 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "081234567899",
    "department_id": 3
  }'
```

---

### 8. Activate Doctor

**Endpoint:** `PATCH /api/v1/doctors/:id/activate`

**Description:** Admin mengaktifkan doctor yang tidak aktif.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Doctor ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor activated successfully",
  "data": {
    "id": 2,
    "employee_id": "DOC002",
    "full_name": "Dr. Jane Doe",
    "is_active": true,
    "updated_at": "2024-01-19T15:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/doctors/2/activate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

**Notes:**

- Doctor yang diaktifkan dapat menerima appointment baru
- Tidak mempengaruhi data lainnya

---

### 9. Deactivate Doctor

**Endpoint:** `PATCH /api/v1/doctors/:id/deactivate`

**Description:** Admin menonaktifkan doctor (tidak bisa menerima appointment baru).

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Doctor ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor deactivated successfully",
  "data": {
    "id": 2,
    "employee_id": "DOC002",
    "full_name": "Dr. Jane Doe",
    "is_active": false,
    "updated_at": "2024-01-19T15:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/doctors/2/deactivate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

**Notes:**

- Doctor yang deactivated tidak bisa menerima appointment baru
- Appointment yang sudah ada tetap valid
- Berbeda dengan soft delete

---

### 10. Soft Delete Doctor

**Endpoint:** `DELETE /api/v1/doctors/:id`

**Description:** Admin menghapus doctor (soft delete). Data tetap ada dan bisa di-restore.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Doctor ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor deleted successfully",
  "data": null
}
```

**Notes:**

- Doctor yang dihapus tidak muncul di list normal
- Medical records dan appointment history tetap tersimpan
- Bisa di-restore dengan endpoint restore

**⚠️ Business Rules:**

- Tidak bisa delete doctor yang memiliki active appointments
- Tidak bisa delete doctor yang sedang on duty

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/doctors/2 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 11. Restore Doctor

**Endpoint:** `PATCH /api/v1/doctors/:id/restore`

**Description:** Admin me-restore doctor yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Doctor ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor restored successfully",
  "data": {
    "id": 2,
    "employee_id": "DOC002",
    "full_name": "Dr. Jane Doe",
    "is_active": false,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T16:00:00Z",
    "deleted_at": null
  }
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/doctors/2/restore \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

**Notes:**

- Setelah di-restore, doctor masih dalam status `is_active: false`
- Perlu activate manual jika ingin mengaktifkan kembali

---

## Super Admin Endpoints

### 12. Hard Delete Doctor

**Endpoint:** `DELETE /api/v1/doctors/:id/hard-delete`

**Description:** Super Admin menghapus doctor secara permanen dari database.

**Authentication:** Required (Super Admin Only)

**Request Headers:**

```
Authorization: Bearer <super-admin-token>
```

**URL Parameters:**

- `id`: Doctor ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Doctor permanently deleted",
  "data": null
}
```

**⚠️ WARNING:**

- Data dihapus permanen dari database
- Tidak bisa di-restore
- Gunakan dengan hati-hati
- Pastikan tidak ada relasi dengan medical records

**⚠️ Business Rules:**

- Tidak bisa hard delete jika masih ada medical records
- Tidak bisa hard delete jika masih ada appointments
- Tidak bisa hard delete jika masih ada schedules

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/doctors/2/hard-delete \
  -H "Authorization: Bearer SUPER_ADMIN_JWT_TOKEN"
```

---

## Database Model

### Table: doctors

| Field                    | Type         | Constraints                                              | Description                                    |
| ------------------------ | ------------ | -------------------------------------------------------- | ---------------------------------------------- |
| id                       | BIGINT       | PRIMARY KEY, AUTO_INCREMENT                              | Unique identifier untuk doctor                 |
| user_id                  | BIGINT       | FOREIGN KEY (users.id), INDEX                            | Reference ke user yang login                   |
| employee_id              | VARCHAR(50)  | UNIQUE, NOT NULL, INDEX                                  | ID karyawan dari HR                            |
| full_name                | VARCHAR(100) | NOT NULL                                                 | Nama lengkap dokter                            |
| doctor_specialization_id | BIGINT       | FOREIGN KEY (doctor_specializations.id), NOT NULL, INDEX | Spesialisasi dokter (reference ke master data) |
| license_number           | VARCHAR(50)  | UNIQUE, NOT NULL                                         | Nomor lisensi STR/SIP                          |
| phone                    | VARCHAR(15)  | NULLABLE                                                 | Nomor telepon                                  |
| email                    | VARCHAR(100) | NULLABLE                                                 | Email kontak                                   |
| department_id            | BIGINT       | FOREIGN KEY (departments.id), NOT NULL, INDEX            | Department tempat bekerja                      |
| is_active                | BOOLEAN      | NOT NULL, DEFAULT true, INDEX                            | Status aktif dokter                            |
| created_at               | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP                                | Waktu pembuatan record                         |
| updated_at               | TIMESTAMP    | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP    | Waktu update terakhir                          |
| deleted_at               | TIMESTAMP    | INDEX, NULLABLE                                          | Soft delete timestamp                          |

**Indexes:**

- Primary Key: id
- Unique Index: employee_id, license_number
- Foreign Key: user_id, department_id, doctor_specialization_id
- Regular Index: is_active, deleted_at

**Relationships:**

- Belongs To User (one-to-one)
- Belongs To Department (many-to-one)
- Belongs To DoctorSpecialization (many-to-one)
- Has Many Appointments (one-to-many)
- Has Many Medical Records (one-to-many)
- Has Many Referrals (one-to-many)

**Notes:**

- License number adalah nomor STR (Surat Tanda Registrasi) atau SIP (Surat Izin Praktik)
- Department menentukan lokasi kerja utama dokter
- Doctor Specialization adalah reference ke master data untuk standardisasi
- is_active dapat diubah untuk non-permanent leave atau retirement

---

## Error Responses

### 400 Bad Request

```json
{
  "success": false,
  "message": "Validation error",
  "error": "employee_id already exists"
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
  "message": "Doctor not found",
  "error": "doctor not found"
}
```

### 409 Conflict

```json
{
  "success": false,
  "message": "Cannot delete doctor",
  "error": "doctor has active appointments or medical records"
}
```

---

## Data Models

### Doctor Object

```json
{
  "id": 1,
  "user_id": 2,
  "employee_id": "DOC001",
  "full_name": "Dr. John Smith",
  "license_number": "LIC123456",
  "phone": "081234567890",
  "email": "drsmith@hospital.com",
  "department_id": 1,
  "department": {
    "id": 1,
    "name": "Cardiology Department",
    "code": "CARD"
  },
  "doctor_specialization_id": 1,
  "doctor_specialization": {
    "id": 1,
    "name": "Cardiology",
    "code": "CARDIO"
  },
  "is_active": true,
  "created_at": "2024-01-19T10:00:00Z",
  "updated_at": "2024-01-19T10:00:00Z",
  "deleted_at": null
}
```

### Doctor Summary (for lists)

```json
{
  "id": 1,
  "employee_id": "DOC001",
  "full_name": "Dr. John Smith",
  "phone": "081234567890",
  "department": {
    "id": 1,
    "name": "Cardiology Department"
  },
  "doctor_specialization": {
    "id": 1,
    "name": "Cardiology",
    "code": "CARDIO"
  },
  "is_active": true,
  "created_at": "2024-01-19T10:00:00Z"
}
```

---

## Common Specializations

### Medical Specializations

- **Cardiology** (Kardiologi) - Jantung dan pembuluh darah
- **Neurology** (Neurologi) - Sistem saraf
- **Pediatrics** (Pediatri) - Anak
- **Orthopedics** (Ortopedi) - Tulang dan sendi
- **Dermatology** (Dermatologi) - Kulit
- **Ophthalmology** (Oftalmologi) - Mata
- **ENT** (THT) - Telinga, Hidung, Tenggorokan
- **Internal Medicine** (Penyakit Dalam)
- **General Surgery** (Bedah Umum)
- **Obstetrics & Gynecology** (Kandungan)
- **Psychiatry** (Psikiatri) - Kesehatan mental
- **Radiology** (Radiologi) - Pencitraan medis
- **Anesthesiology** (Anestesiologi) - Pembiusan
- **Oncology** (Onkologi) - Kanker
- **Urology** (Urologi) - Saluran kemih
- **Gastroenterology** (Gastroenterologi) - Pencernaan
- **Pulmonology** (Pulmonologi) - Paru-paru
- **Endocrinology** (Endokrinologi) - Hormon
- **Nephrology** (Nefrologi) - Ginjal
- **Rheumatology** (Reumatologi) - Rematik

---

## Business Rules

1. **Employee ID Uniqueness**: Employee ID harus unik
2. **Auto-Generated ID**: Format DOC-NNNN
3. **License Number Validation**: License number harus valid dan unique
4. **Active Status**: Hanya doctor aktif yang bisa menerima appointment
5. **Department Assignment**: Doctor harus assigned ke department
6. **Specialization Required**: Specialization wajib diisi
7. **Contact Information**: Phone atau email minimal salah satu harus ada
8. **Soft Delete Protection**: Cannot delete with active appointments
9. **Deactivate vs Delete**: Deactivate untuk temporary, Delete untuk permanent
10. **Professional Credentials**: License number must be verified

---

## Common Use Cases

### Use Case 1: Add New Doctor to Hospital

```bash
# Admin creates new doctor
curl -X POST http://localhost:8080/api/v1/doctors \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "DOC015",
    "full_name": "Dr. Sarah Johnson",
    "doctor_specialization_id": 4,
    "license_number": "LIC987654",
    "phone": "081234567892",
    "email": "drsarah@hospital.com",
    "department_id": 3,
    "is_active": true
  }'
```

### Use Case 2: Patient Finding Doctors by Department or Specialization

```bash
# Patient lists all available doctors
curl -X GET "http://localhost:8080/api/v1/doctors?is_active=true&page=1&page_size=20" \
  -H "Authorization: Bearer PATIENT_TOKEN"

# Patient filters doctors by specialization
curl -X GET "http://localhost:8080/api/v1/doctors?is_active=true&doctor_specialization_id=1" \
  -H "Authorization: Bearer PATIENT_TOKEN"

# Patient views doctor details
curl -X GET http://localhost:8080/api/v1/doctors/1 \
  -H "Authorization: Bearer PATIENT_TOKEN"
```

### Use Case 3: Doctor Updates Own Profile

```bash
# Doctor updates contact info
curl -X PUT http://localhost:8080/api/v1/doctors/me \
  -H "Authorization: Bearer DOCTOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "081234567899",
    "email": "newemail@hospital.com"
  }'
```

### Use Case 4: Admin Manages Doctor Status

```bash
# Deactivate doctor temporarily (vacation, leave)
curl -X PATCH http://localhost:8080/api/v1/doctors/5/deactivate \
  -H "Authorization: Bearer ADMIN_TOKEN"

# Reactivate doctor
curl -X PATCH http://localhost:8080/api/v1/doctors/5/activate \
  -H "Authorization: Bearer ADMIN_TOKEN"

# Soft delete doctor (resigned)
curl -X DELETE http://localhost:8080/api/v1/doctors/5 \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

---

## Testing Examples

### Test 1: Complete Doctor Management Flow

```bash
# 1. Create Doctor (Admin)
curl -X POST http://localhost:8080/api/v1/doctors \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Dr. Test Doctor",
    "doctor_specialization_id": 5,
    "license_number": "LIC-TEST-001",
    "phone": "081234567890",
    "email": "testdoc@hospital.com",
    "department_id": 1
  }'

# 2. List All Active Doctors
curl -X GET "http://localhost:8080/api/v1/doctors?page=1&page_size=10" \
  -H "Authorization: Bearer PATIENT_TOKEN"

# 3. Filter Doctors by Specialization
curl -X GET "http://localhost:8080/api/v1/doctors?doctor_specialization_id=5&is_active=true" \
  -H "Authorization: Bearer PATIENT_TOKEN"

# 4. View Doctor Details
curl -X GET http://localhost:8080/api/v1/doctors/1 \
  -H "Authorization: Bearer DOCTOR_TOKEN"

# 5. Update Doctor (Admin)
curl -X PUT http://localhost:8080/api/v1/doctors/1 \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"phone": "081234567899"}'

# 6. Deactivate Doctor
curl -X PATCH http://localhost:8080/api/v1/doctors/1/deactivate \
  -H "Authorization: Bearer ADMIN_TOKEN"

# 7. Reactivate Doctor
curl -X PATCH http://localhost:8080/api/v1/doctors/1/activate \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

### Test 2: Doctor Self-Management

```bash
# 1. Doctor views own profile
curl -X GET http://localhost:8080/api/v1/doctors/me \
  -H "Authorization: Bearer DOCTOR_TOKEN"

# 2. Doctor updates contact info
curl -X PUT http://localhost:8080/api/v1/doctors/me \
  -H "Authorization: Bearer DOCTOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "081234567899",
    "email": "updated@hospital.com"
  }'
```

### Test 3: Search and Filter Doctors

```bash
# 1. Search by name
curl -X GET "http://localhost:8080/api/v1/doctors?search=John&page=1" \
  -H "Authorization: Bearer PATIENT_TOKEN"

# 2. Filter by department and specialization
curl -X GET "http://localhost:8080/api/v1/doctors?department_id=1&specialization=Cardiology&is_active=true" \
  -H "Authorization: Bearer PATIENT_TOKEN"

# 3. Sort by name ascending
curl -X GET "http://localhost:8080/api/v1/doctors?sort_by=full_name&sort_dir=asc" \
  -H "Authorization: Bearer PATIENT_TOKEN"
```

---

## Notes

- Employee ID auto-generated dengan format DOC-NNNN
- License number harus unique dan divalidasi
- Doctor bisa deactivated (temporary) atau deleted (soft delete)
- Deactivated doctor tidak bisa menerima appointment baru
- Deleted doctor tidak muncul di list tetapi data tetap tersimpan
- Specialization mempengaruhi appointment routing
- Department assignment penting untuk organization
- Active status controlled oleh admin

---

**Last Updated:** 2024-01-19  
**API Version:** 1.0.0
