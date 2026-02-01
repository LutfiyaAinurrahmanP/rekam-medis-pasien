# Patients API Documentation

## Overview

API untuk manajemen data patients (pasien) dalam sistem rekam medis. Patients adalah pengguna yang menerima layanan medis di rumah sakit.

**Base URL:** `/api/v1/patients`

---

## Table of Contents

- [Authentication](#authentication)
- [Authorization](#authorization)
- [Endpoints Summary](#endpoints-summary)
- [Self-Owned Endpoints](#self-owned-endpoints)
- [Public Endpoints](#public-endpoints)
- [Admin Endpoints](#admin-endpoints)
- [Super Admin Endpoints](#super-admin-endpoints)
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

| Endpoint                         | Patient (Own) | Doctor | Receptionist | Admin | Super Admin |
| -------------------------------- | ------------- | ------ | ------------ | ----- | ----------- |
| GET /patients/me                 | ✅            | ❌     | ❌           | ❌    | ❌          |
| PUT /patients/me                 | ✅            | ❌     | ❌           | ❌    | ❌          |
| GET /patients                    | ❌            | ✅     | ✅           | ✅    | ✅          |
| GET /patients/deleted            | ❌            | ❌     | ✅           | ✅    | ✅          |
| GET /patients/:id                | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| POST /patients                   | ❌            | ❌     | ✅           | ✅    | ✅          |
| PUT /patients/:id                | ❌            | ❌     | ✅           | ✅    | ✅          |
| DELETE /patients/:id             | ❌            | ❌     | ✅           | ✅    | ✅          |
| PATCH /patients/:id/restore      | ❌            | ❌     | ✅           | ✅    | ✅          |
| DELETE /patients/:id/hard-delete | ❌            | ❌     | ❌           | ❌    | ✅          |

---

## Endpoints Summary

### Self-Owned Endpoints (`/me`)

| Method | Endpoint       | Description            | Auth    |
| ------ | -------------- | ---------------------- | ------- |
| GET    | `/patients/me` | Get my patient data    | Patient |
| PUT    | `/patients/me` | Update my patient data | Patient |

### Public Endpoints (Staff Access)

| Method | Endpoint               | Description           | Role Required                            |
| ------ | ---------------------- | --------------------- | ---------------------------------------- |
| GET    | `/patients`            | List active patients  | Doctor, Receptionist, Admin, Super Admin |
| GET    | `/patients/deleted`    | List deleted patients | Receptionist, Admin, Super Admin         |
| GET    | `/patients/:id`        | Get patient by ID     | All (with ownership check)               |
| GET    | `/patients/code/:code` | Get patient by code   | Doctor, Receptionist, Admin, Super Admin |
| GET    | `/patients/search`     | Advanced search       | Doctor, Receptionist, Admin, Super Admin |

### Admin Endpoints (Management)

| Method | Endpoint                | Description             | Role Required                    |
| ------ | ----------------------- | ----------------------- | -------------------------------- |
| POST   | `/patients`             | Create patient          | Receptionist, Admin, Super Admin |
| PUT    | `/patients/:id`         | Update patient          | Receptionist, Admin, Super Admin |
| DELETE | `/patients/:id`         | Soft delete patient     | Receptionist, Admin, Super Admin |
| PATCH  | `/patients/:id/restore` | Restore deleted patient | Receptionist, Admin, Super Admin |

### Super Admin Endpoints

| Method | Endpoint                    | Description                | Role Required |
| ------ | --------------------------- | -------------------------- | ------------- |
| DELETE | `/patients/:id/hard-delete` | Permanently delete patient | Super Admin   |

---

## Self-Owned Endpoints

### 1. Get My Patient Data

**Endpoint:** `GET /api/v1/patients/me`

**Description:** Patient mendapatkan data diri sendiri.

**Authentication:** Required (Patient Role)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Patient data retrieved successfully",
  "data": {
    "id": 1,
    "user_id": 5,
    "patient_code": "P-2024-001",
    "full_name": "John Doe",
    "date_of_birth": "1990-05-15",
    "gender": "male",
    "blood_type": "O+",
    "phone": "081234567890",
    "email": "johndoe@example.com",
    "address": "Jl. Merdeka No. 123, Jakarta",
    "emergency_contact_name": "Jane Doe",
    "emergency_contact_phone": "081234567891",
    "insurance_number": "INS-123456",
    "insurance_provider": "BPJS Kesehatan",
    "allergies": "Penicillin, Seafood",
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Patient data not found",
  "error": "patient profile not created yet"
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/patients/me \
  -H "Authorization: Bearer PATIENT_JWT_TOKEN"
```

**Notes:**

- Patient harus sudah memiliki patient record (dibuat saat registrasi atau oleh receptionist)
- Hanya bisa melihat data diri sendiri

---

### 2. Update My Patient Data

**Endpoint:** `PUT /api/v1/patients/me`

**Description:** Patient mengupdate data diri sendiri.

**Authentication:** Required (Patient Role)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "full_name": "John Doe Updated",
  "phone": "081234567899",
  "email": "johndoe_new@example.com",
  "address": "Jl. Merdeka No. 124, Jakarta",
  "emergency_contact_name": "Jane Doe Smith",
  "emergency_contact_phone": "081234567892",
  "insurance_number": "INS-123456-NEW",
  "insurance_provider": "BPJS Kesehatan Premium",
  "allergies": "Penicillin, Seafood, Peanuts"
}
```

**Field Rules:**

- `full_name`: optional, max 100 characters
- `phone`: optional, max 15 characters
- `email`: optional, valid email format
- `address`: optional, text
- `emergency_contact_name`: optional, max 100
- `emergency_contact_phone`: optional, max 15
- `insurance_number`: optional, max 50
- `insurance_provider`: optional, max 100
- `allergies`: optional, text
- ❌ **Cannot update:** `patient_code`, `date_of_birth`, `gender`, `blood_type`

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Patient data updated successfully",
  "data": {
    "id": 1,
    "patient_code": "P-2024-001",
    "full_name": "John Doe Updated",
    "date_of_birth": "1990-05-15",
    "gender": "male",
    "blood_type": "O+",
    "phone": "081234567899",
    "email": "johndoe_new@example.com",
    "address": "Jl. Merdeka No. 124, Jakarta",
    "emergency_contact_name": "Jane Doe Smith",
    "emergency_contact_phone": "081234567892",
    "insurance_number": "INS-123456-NEW",
    "insurance_provider": "BPJS Kesehatan Premium",
    "allergies": "Penicillin, Seafood, Peanuts",
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T14:30:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/patients/me \
  -H "Authorization: Bearer PATIENT_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "081234567899",
    "address": "Jl. Merdeka No. 124, Jakarta"
  }'
```

---

## Public Endpoints (Staff Access)

### 3. List Patients

**Endpoint:** `GET /api/v1/patients`

**Description:** Mendapatkan daftar patient aktif dengan pagination, search, dan filter.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter            | Type    | Default    | Description                                                     |
| -------------------- | ------- | ---------- | --------------------------------------------------------------- |
| `page`               | integer | 1          | Halaman                                                         |
| `page_size`          | integer | 10         | Jumlah data per halaman (max: 100)                              |
| `search`             | string  | -          | Cari berdasarkan name, patient_code, phone, atau email          |
| `gender`             | string  | -          | Filter by gender (male, female, other)                          |
| `blood_type`         | string  | -          | Filter by blood type (A+, A-, B+, B-, AB+, AB-, O+, O-)         |
| `insurance_provider` | string  | -          | Filter by insurance provider                                    |
| `min_age`            | integer | -          | Filter minimum age                                              |
| `max_age`            | integer | -          | Filter maximum age                                              |
| `sort_by`            | string  | created_at | Sort field (created_at, full_name, patient_code, date_of_birth) |
| `sort_dir`           | string  | desc       | Sort direction (asc, desc)                                      |

**Example Request:**

```
GET /api/v1/patients?page=1&page_size=10&search=john&gender=male&blood_type=O+&sort_by=full_name&sort_dir=asc
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Patients retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "patient_code": "P-2024-001",
        "full_name": "John Doe",
        "date_of_birth": "1990-05-15",
        "age": 33,
        "gender": "male",
        "blood_type": "O+",
        "phone": "081234567890",
        "email": "johndoe@example.com",
        "insurance_provider": "BPJS Kesehatan",
        "created_at": "2024-01-19T10:00:00Z"
      },
      {
        "id": 2,
        "patient_code": "P-2024-002",
        "full_name": "John Smith",
        "date_of_birth": "1985-03-20",
        "age": 38,
        "gender": "male",
        "blood_type": "O+",
        "phone": "081234567892",
        "email": "johnsmith@example.com",
        "insurance_provider": "Mandiri Inhealth",
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
curl -X GET "http://localhost:8080/api/v1/patients?page=1&page_size=10&search=john&gender=male" \
  -H "Authorization: Bearer STAFF_JWT_TOKEN"
```

---

### 4. List Deleted Patients

**Endpoint:** `GET /api/v1/patients/deleted`

**Description:** Mendapatkan daftar patient yang sudah di-soft delete.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**
Same as List Patients endpoint.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Deleted patients retrieved successfully",
  "data": {
    "data": [
      {
        "id": 10,
        "patient_code": "P-2024-010",
        "full_name": "Deleted Patient",
        "date_of_birth": "1995-01-01",
        "age": 29,
        "gender": "female",
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
curl -X GET "http://localhost:8080/api/v1/patients/deleted?page=1&page_size=10" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 5. Get Patient by ID

**Endpoint:** `GET /api/v1/patients/:id`

**Description:** Mendapatkan detail patient berdasarkan ID, termasuk medical history.

**Authentication:** Required (All authenticated, with ownership check for patients)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Patient ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Patient retrieved successfully",
  "data": {
    "id": 1,
    "user_id": 5,
    "patient_code": "P-2024-001",
    "full_name": "John Doe",
    "date_of_birth": "1990-05-15",
    "age": 33,
    "gender": "male",
    "blood_type": "O+",
    "phone": "081234567890",
    "email": "johndoe@example.com",
    "address": "Jl. Merdeka No. 123, Jakarta",
    "emergency_contact_name": "Jane Doe",
    "emergency_contact_phone": "081234567891",
    "insurance_number": "INS-123456",
    "insurance_provider": "BPJS Kesehatan",
    "allergies": "Penicillin, Seafood",
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z",
    "medical_records_count": 5,
    "appointments_count": 3,
    "last_visit": "2024-01-15T14:00:00Z"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Patient not found",
  "error": "patient not found"
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/patients/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Notes:**

- Patient hanya bisa melihat data diri sendiri
- Staff (Doctor, Receptionist, Admin) bisa melihat semua patient

---

### 6. Get Patient by Code

**Endpoint:** `GET /api/v1/patients/code/:code`

**Description:** Mendapatkan patient berdasarkan patient code (untuk quick search).

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `code`: Patient Code (string), e.g., "P-2024-001"

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Patient retrieved successfully",
  "data": {
    "id": 1,
    "patient_code": "P-2024-001",
    "full_name": "John Doe",
    "date_of_birth": "1990-05-15",
    "age": 33,
    "gender": "male",
    "blood_type": "O+",
    "phone": "081234567890",
    "email": "johndoe@example.com"
  }
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/patients/code/P-2024-001 \
  -H "Authorization: Bearer STAFF_JWT_TOKEN"
```

---

### 7. Advanced Search Patients

**Endpoint:** `GET /api/v1/patients/search`

**Description:** Pencarian patient dengan multiple criteria.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter          | Type   | Description                 |
| ------------------ | ------ | --------------------------- |
| `full_name`        | string | Exact or partial name match |
| `patient_code`     | string | Exact code match            |
| `phone`            | string | Exact phone match           |
| `email`            | string | Exact email match           |
| `date_of_birth`    | date   | Exact DOB (YYYY-MM-DD)      |
| `insurance_number` | string | Insurance number match      |

**Example Request:**

```
GET /api/v1/patients/search?full_name=John&date_of_birth=1990-05-15
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Search results retrieved successfully",
  "data": [
    {
      "id": 1,
      "patient_code": "P-2024-001",
      "full_name": "John Doe",
      "date_of_birth": "1990-05-15",
      "gender": "male",
      "phone": "081234567890"
    }
  ]
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/patients/search?full_name=John&date_of_birth=1990-05-15" \
  -H "Authorization: Bearer STAFF_JWT_TOKEN"
```

---

## Admin Endpoints (Management)

### 8. Create Patient

**Endpoint:** `POST /api/v1/patients`

**Description:** Staff membuat patient record baru.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "user_id": 5,
  "patient_code": "P-2024-001",
  "full_name": "John Doe",
  "date_of_birth": "1990-05-15",
  "gender": "male",
  "blood_type": "O+",
  "phone": "081234567890",
  "email": "johndoe@example.com",
  "address": "Jl. Merdeka No. 123, Jakarta",
  "emergency_contact_name": "Jane Doe",
  "emergency_contact_phone": "081234567891",
  "insurance_number": "INS-123456",
  "insurance_provider": "BPJS Kesehatan",
  "allergies": "Penicillin, Seafood"
}
```

**Field Rules:**

- `user_id`: optional, FK to users table (if patient has user account)
- `patient_code`: required, unique, max 20 characters, indexed
- `full_name`: required, max 100 characters
- `date_of_birth`: required, date format (YYYY-MM-DD)
- `gender`: required, enum (male, female, other)
- `blood_type`: optional, max 5 characters (A+, A-, B+, B-, AB+, AB-, O+, O-)
- `phone`: optional, max 15 characters
- `email`: optional, valid email format
- `address`: optional, text
- `emergency_contact_name`: optional, max 100
- `emergency_contact_phone`: optional, max 15
- `insurance_number`: optional, max 50
- `insurance_provider`: optional, max 100
- `allergies`: optional, text

**Response Success (201 Created):**

```json
{
  "success": true,
  "message": "Patient created successfully",
  "data": {
    "id": 1,
    "user_id": 5,
    "patient_code": "P-2024-001",
    "full_name": "John Doe",
    "date_of_birth": "1990-05-15",
    "age": 33,
    "gender": "male",
    "blood_type": "O+",
    "phone": "081234567890",
    "email": "johndoe@example.com",
    "address": "Jl. Merdeka No. 123, Jakarta",
    "emergency_contact_name": "Jane Doe",
    "emergency_contact_phone": "081234567891",
    "insurance_number": "INS-123456",
    "insurance_provider": "BPJS Kesehatan",
    "allergies": "Penicillin, Seafood",
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z"
  }
}
```

**Response Error (400 Bad Request):**

```json
{
  "success": false,
  "message": "Failed to create patient",
  "error": "patient_code already exists"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/patients \
  -H "Authorization: Bearer RECEPTIONIST_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "patient_code": "P-2024-001",
    "full_name": "John Doe",
    "date_of_birth": "1990-05-15",
    "gender": "male",
    "phone": "081234567890"
  }'
```

**Notes:**

- Patient code auto-generated jika tidak diisi
- Format: P-YYYY-NNNN (P-2024-0001, P-2024-0002, dst)

---

### 9. Update Patient

**Endpoint:** `PUT /api/v1/patients/:id`

**Description:** Staff mengupdate data patient berdasarkan ID.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Patient ID (integer)

**Request Body:**

```json
{
  "full_name": "John Doe Updated",
  "phone": "081234567899",
  "email": "johndoe_new@example.com",
  "address": "Jl. Merdeka No. 124, Jakarta",
  "blood_type": "O+",
  "emergency_contact_name": "Jane Doe Smith",
  "emergency_contact_phone": "081234567892",
  "insurance_number": "INS-123456-NEW",
  "insurance_provider": "BPJS Kesehatan Premium",
  "allergies": "Penicillin, Seafood, Peanuts"
}
```

**Field Rules:**

- All fields optional
- Validation same as create
- Patient code cannot be changed
- Date of birth can be updated (with caution)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Patient updated successfully",
  "data": {
    "id": 1,
    "patient_code": "P-2024-001",
    "full_name": "John Doe Updated",
    "date_of_birth": "1990-05-15",
    "age": 33,
    "gender": "male",
    "blood_type": "O+",
    "phone": "081234567899",
    "email": "johndoe_new@example.com",
    "address": "Jl. Merdeka No. 124, Jakarta",
    "emergency_contact_name": "Jane Doe Smith",
    "emergency_contact_phone": "081234567892",
    "insurance_number": "INS-123456-NEW",
    "insurance_provider": "BPJS Kesehatan Premium",
    "allergies": "Penicillin, Seafood, Peanuts",
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T15:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/patients/1 \
  -H "Authorization: Bearer RECEPTIONIST_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "081234567899",
    "address": "Jl. Merdeka No. 124, Jakarta"
  }'
```

---

### 10. Soft Delete Patient

**Endpoint:** `DELETE /api/v1/patients/:id`

**Description:** Staff menghapus patient (soft delete). Data tetap ada dan bisa di-restore.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Patient ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Patient deleted successfully",
  "data": null
}
```

**Notes:**

- Patient yang dihapus tidak muncul di list normal
- Medical records tetap tersimpan
- Bisa di-restore dengan endpoint restore

**⚠️ Business Rules:**

- Tidak bisa delete patient yang memiliki active appointments
- Tidak bisa delete patient yang sedang dirawat (hospitalized)

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/patients/1 \
  -H "Authorization: Bearer RECEPTIONIST_JWT_TOKEN"
```

---

### 11. Restore Patient

**Endpoint:** `PATCH /api/v1/patients/:id/restore`

**Description:** Staff me-restore patient yang sudah di-soft delete.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Patient ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Patient restored successfully",
  "data": null
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/patients/1/restore \
  -H "Authorization: Bearer RECEPTIONIST_JWT_TOKEN"
```

---

## Super Admin Endpoints

### 12. Hard Delete Patient

**Endpoint:** `DELETE /api/v1/patients/:id/hard-delete`

**Description:** Super Admin menghapus patient secara permanen dari database.

**Authentication:** Required (Super Admin Only)

**Request Headers:**

```
Authorization: Bearer <super-admin-token>
```

**URL Parameters:**

- `id`: Patient ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Patient permanently deleted",
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
- Tidak bisa hard delete jika masih ada billing records

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/patients/1/hard-delete \
  -H "Authorization: Bearer SUPER_ADMIN_JWT_TOKEN"
```

---

## Error Responses

### 400 Bad Request

```json
{
  "success": false,
  "message": "Validation error",
  "error": "patient_code already exists"
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
  "message": "Patient not found",
  "error": "patient not found"
}
```

### 409 Conflict

```json
{
  "success": false,
  "message": "Cannot delete patient",
  "error": "patient has active appointments or medical records"
}
```

---

## Data Models

### Patient Object

```json
{
  "id": 1,
  "user_id": 5,
  "patient_code": "P-2024-001",
  "full_name": "John Doe",
  "date_of_birth": "1990-05-15",
  "age": 33,
  "gender": "male",
  "blood_type": "O+",
  "phone": "081234567890",
  "email": "johndoe@example.com",
  "address": "Jl. Merdeka No. 123, Jakarta",
  "emergency_contact_name": "Jane Doe",
  "emergency_contact_phone": "081234567891",
  "insurance_number": "INS-123456",
  "insurance_provider": "BPJS Kesehatan",
  "allergies": "Penicillin, Seafood",
  "created_at": "2024-01-19T10:00:00Z",
  "updated_at": "2024-01-19T10:00:00Z",
  "deleted_at": null
}
```

### Patient Summary (for lists)

```json
{
  "id": 1,
  "patient_code": "P-2024-001",
  "full_name": "John Doe",
  "date_of_birth": "1990-05-15",
  "age": 33,
  "gender": "male",
  "blood_type": "O+",
  "phone": "081234567890",
  "insurance_provider": "BPJS Kesehatan",
  "created_at": "2024-01-19T10:00:00Z"
}
```

---

## Business Rules

1. **Patient Code Uniqueness**: Patient code harus unik
2. **Auto-Generated Code**: Format P-YYYY-NNNN
3. **Age Calculation**: Age dihitung otomatis dari date_of_birth
4. **Gender Validation**: Must be male, female, or other
5. **Blood Type Format**: A+, A-, B+, B-, AB+, AB-, O+, O-
6. **Emergency Contact**: Wajib diisi untuk keamanan
7. **Insurance**: Optional, untuk billing purposes
8. **Allergies**: Critical information, must be up-to-date
9. **Soft Delete Protection**: Cannot delete with active appointments
10. **GDPR Compliance**: Patient can update own data

---

## Common Use Cases

### Use Case 1: Patient Registration (Walk-in)

```bash
# Receptionist creates patient record
POST /api/v1/patients
{
  "patient_code": "P-2024-050",
  "full_name": "New Patient",
  "date_of_birth": "1995-01-01",
  "gender": "male",
  "phone": "081234567890",
  "emergency_contact_name": "Family Member",
  "emergency_contact_phone": "081234567891"
}
```

### Use Case 2: Patient Self-Registration (Online)

```bash
# 1. User registers account
POST /api/v1/auth/register

# 2. System auto-creates patient record
# 3. Patient completes profile
PUT /api/v1/patients/me
```

### Use Case 3: Doctor Views Patient

```bash
# Doctor searches patient by code
GET /api/v1/patients/code/P-2024-001

# Doctor views full patient details
GET /api/v1/patients/1
```

---

## Testing Examples

### Test 1: Complete Patient Registration Flow

```bash
# 1. Create Patient (Receptionist)
curl -X POST http://localhost:8080/api/v1/patients \
  -H "Authorization: Bearer RECEPTIONIST_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test Patient",
    "date_of_birth": "1990-01-01",
    "gender": "male",
    "phone": "081234567890"
  }'

# 2. Search Patient
curl -X GET "http://localhost:8080/api/v1/patients?search=Test+Patient" \
  -H "Authorization: Bearer DOCTOR_TOKEN"

# 3. View Patient Details
curl -X GET http://localhost:8080/api/v1/patients/1 \
  -H "Authorization: Bearer DOCTOR_TOKEN"

# 4. Update Patient
curl -X PUT http://localhost:8080/api/v1/patients/1 \
  -H "Authorization: Bearer RECEPTIONIST_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"phone": "081234567899"}'
```

---

## Notes

- Patient code auto-generated jika tidak diisi
- Age calculated from date_of_birth
- Allergies information critical untuk safety
- Emergency contact wajib untuk emergency situations
- Insurance info untuk billing dan claims
- Soft delete preserves medical history
- HIPAA/GDPR compliant data handling

---

**Last Updated:** 2024-01-19  
**API Version:** 1.0.0
