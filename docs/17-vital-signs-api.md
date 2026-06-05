# Vital Signs API Documentation

## Overview

API untuk manajemen data tanda vital pasien (vital signs) dalam sistem rekam medis. Vital Signs mencatat pengukuran tanda-tanda vital pasien saat kunjungan, meliputi tekanan darah, denyut jantung, suhu tubuh, laju pernafasan, saturasi oksigen, berat badan, tinggi badan, dan BMI yang dihitung secara otomatis.

**Base URL:** `/api/v1/vital-signs`

---

## Table of Contents

- [Authentication](#authentication)
- [Authorization](#authorization)
- [Endpoints Summary](#endpoints-summary)
- [Endpoints Detail](#endpoints-detail)
- [Database Model](#database-model)
- [Error Responses](#error-responses)

---

## Authentication

Semua endpoints memerlukan JWT token di header:

```
Authorization: Bearer <your-jwt-token>
```

---

## Authorization

| Endpoint                                   | Patient | Doctor | Receptionist | Admin | Super Admin |
| ------------------------------------------ | ------- | ------ | ------------ | ----- | ----------- |
| GET /vital-signs                           | ❌      | ✅     | ✅           | ✅    | ✅          |
| GET /vital-signs/deleted                   | ❌      | ❌     | ❌           | ✅    | ✅          |
| GET /vital-signs/:id                       | ✅      | ✅     | ✅           | ✅    | ✅          |
| POST /vital-signs                          | ❌      | ✅     | ✅           | ✅    | ✅          |
| PUT /vital-signs/:id                       | ❌      | ✅     | ✅           | ✅    | ✅          |
| DELETE /vital-signs/:id                    | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /vital-signs/:id/restore             | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /vital-signs/:id/hard-delete        | ❌      | ❌     | ❌           | ❌    | ✅          |

---

## Endpoints Summary

| Method | Endpoint                                 | Description                        | Role Required                            |
| ------ | ---------------------------------------- | ---------------------------------- | ---------------------------------------- |
| GET    | `/vital-signs`                           | List all vital signs records       | Doctor, Receptionist, Admin, Super Admin |
| GET    | `/vital-signs/deleted`                   | List deleted vital signs records   | Admin, Super Admin                       |
| GET    | `/vital-signs/:id`                       | Get vital signs by ID              | All Authenticated                        |
| GET    | `/vital-signs/medical-record/:record_id` | Get vital signs by medical record  | All Authenticated                        |
| POST   | `/vital-signs`                           | Record new vital signs             | Doctor, Receptionist, Admin, Super Admin |
| PUT    | `/vital-signs/:id`                       | Update vital signs record          | Doctor, Receptionist, Admin, Super Admin |
| DELETE | `/vital-signs/:id`                       | Soft delete vital signs record     | Admin, Super Admin                       |
| PATCH  | `/vital-signs/:id/restore`               | Restore deleted vital signs record | Admin, Super Admin                       |
| DELETE | `/vital-signs/:id/hard-delete`           | Permanently delete vital signs     | Super Admin                              |

---

## Endpoints Detail

### 1. List All Vital Signs Records

**Endpoint:** `GET /api/v1/vital-signs`

**Description:** Mendapatkan daftar semua pencatatan tanda vital dengan pagination, search, dan filter.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

| Header        | Type   | Required | Description      |
| ------------- | ------ | -------- | ---------------- |
| Authorization | string | Yes      | `Bearer <token>` |

**Query Parameters:**

| Parameter         | Type    | Required | Default     | Description                                      |
| ----------------- | ------- | -------- | ----------- | ------------------------------------------------ |
| page              | integer | No       | 1           | Halaman saat ini                                 |
| page_size         | integer | No       | 10          | Jumlah data per halaman (max: 100)               |
| sort_by           | string  | No       | recorded_at | Field sorting: `id`, `recorded_at`, `created_at` |
| sort_dir          | string  | No       | desc        | Arah sorting: `asc` atau `desc`                  |
| medical_record_id | integer | No       | -           | Filter berdasarkan ID rekam medis                |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Vital signs retrieved successfully",
  "data": [
    {
      "id": 1,
      "medical_record_id": 5,
      "blood_pressure_systolic": 120,
      "blood_pressure_diastolic": 80,
      "heart_rate": 75,
      "temperature": 36.8,
      "respiratory_rate": 18,
      "oxygen_saturation": 98.5,
      "weight_kg": 65.5,
      "height_cm": 168.0,
      "bmi": 23.21,
      "recorded_at": "2024-01-15T08:30:00Z",
      "created_at": "2024-01-15T08:30:00Z",
      "updated_at": "2024-01-15T08:30:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "page_size": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/vital-signs?page=1&page_size=10" \
  -H "Authorization: Bearer <token>"
```

---

### 2. List Deleted Vital Signs Records

**Endpoint:** `GET /api/v1/vital-signs/deleted`

**Description:** Mendapatkan daftar pencatatan tanda vital yang telah dihapus (soft deleted).

**Authentication:** Required (Admin, Super Admin)

**Query Parameters:**

| Parameter | Type    | Required | Default     | Description                        |
| --------- | ------- | -------- | ----------- | ---------------------------------- |
| page      | integer | No       | 1           | Halaman saat ini                   |
| page_size | integer | No       | 10          | Jumlah data per halaman (max: 100) |
| sort_by   | string  | No       | recorded_at | Field untuk sorting                |
| sort_dir  | string  | No       | desc        | Arah sorting: `asc` atau `desc`    |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Deleted vital signs retrieved successfully",
  "data": [
    {
      "id": 2,
      "medical_record_id": 3,
      "blood_pressure_systolic": 130,
      "blood_pressure_diastolic": 85,
      "heart_rate": 80,
      "recorded_at": "2024-01-10T09:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "page_size": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/vital-signs/deleted" \
  -H "Authorization: Bearer <token>"
```

---

### 3. Get Vital Signs by ID

**Endpoint:** `GET /api/v1/vital-signs/:id`

**Description:** Mendapatkan detail pencatatan tanda vital berdasarkan ID.

**Authentication:** Required (All Authenticated Users)

**Path Parameters:**

| Parameter | Type    | Required | Description               |
| --------- | ------- | -------- | ------------------------- |
| id        | integer | Yes      | ID pencatatan tanda vital |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Vital signs retrieved successfully",
  "data": {
    "id": 1,
    "medical_record_id": 5,
    "blood_pressure_systolic": 120,
    "blood_pressure_diastolic": 80,
    "heart_rate": 75,
    "temperature": 36.8,
    "respiratory_rate": 18,
    "oxygen_saturation": 98.5,
    "weight_kg": 65.5,
    "height_cm": 168.0,
    "bmi": 23.21,
    "recorded_at": "2024-01-15T08:30:00Z",
    "created_at": "2024-01-15T08:30:00Z",
    "updated_at": "2024-01-15T08:30:00Z",
    "medical_record": {
      "id": 5,
      "visit_date": "2024-01-15",
      "chief_complaint": "Demam dan lemas"
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/vital-signs/1" \
  -H "Authorization: Bearer <token>"
```

---

### 4. Record Vital Signs

**Endpoint:** `POST /api/v1/vital-signs`

**Description:** Mencatat tanda vital baru untuk suatu rekam medis. BMI dihitung otomatis jika `weight_kg` dan `height_cm` diisi. Setiap rekam medis hanya dapat memiliki satu data vital signs.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

| Header        | Type   | Required | Description        |
| ------------- | ------ | -------- | ------------------ |
| Authorization | string | Yes      | `Bearer <token>`   |
| Content-Type  | string | Yes      | `application/json` |

**Request Body:**

| Field                    | Type    | Required | Description                                               |
| ------------------------ | ------- | -------- | --------------------------------------------------------- |
| medical_record_id        | integer | Yes      | ID rekam medis (unik — satu rekam medis satu vital signs) |
| recorded_at              | string  | Yes      | Waktu pencatatan (format: `YYYY-MM-DDTHH:MM:SSZ`)         |
| blood_pressure_systolic  | integer | No       | Tekanan darah sistolik (mmHg)                             |
| blood_pressure_diastolic | integer | No       | Tekanan darah diastolik (mmHg)                            |
| heart_rate               | integer | No       | Denyut jantung (bpm)                                      |
| temperature              | float   | No       | Suhu tubuh (°C), 2 desimal                                |
| respiratory_rate         | integer | No       | Laju pernapasan (per menit)                               |
| oxygen_saturation        | float   | No       | Saturasi oksigen (%), 2 desimal                           |
| weight_kg                | float   | No       | Berat badan (kg), 2 desimal                               |
| height_cm                | float   | No       | Tinggi badan (cm), 2 desimal                              |

> BMI dihitung otomatis: `weight_kg / (height_cm/100)²` jika `weight_kg` dan `height_cm` tersedia.

**Example Request Body:**

```json
{
  "medical_record_id": 5,
  "recorded_at": "2024-01-15T08:30:00Z",
  "blood_pressure_systolic": 120,
  "blood_pressure_diastolic": 80,
  "heart_rate": 75,
  "temperature": 36.8,
  "respiratory_rate": 18,
  "oxygen_saturation": 98.5,
  "weight_kg": 65.5,
  "height_cm": 168.0
}
```

**Response Success (201):**

```json
{
  "status": "success",
  "message": "Vital signs recorded successfully",
  "data": {
    "id": 1,
    "medical_record_id": 5,
    "blood_pressure_systolic": 120,
    "blood_pressure_diastolic": 80,
    "heart_rate": 75,
    "temperature": 36.8,
    "respiratory_rate": 18,
    "oxygen_saturation": 98.5,
    "weight_kg": 65.5,
    "height_cm": 168.0,
    "bmi": 23.21,
    "recorded_at": "2024-01-15T08:30:00Z",
    "created_at": "2024-01-15T08:30:00Z",
    "updated_at": "2024-01-15T08:30:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/vital-signs" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "medical_record_id": 5,
    "recorded_at": "2024-01-15T08:30:00Z",
    "blood_pressure_systolic": 120,
    "blood_pressure_diastolic": 80,
    "heart_rate": 75,
    "temperature": 36.8,
    "respiratory_rate": 18,
    "oxygen_saturation": 98.5,
    "weight_kg": 65.5,
    "height_cm": 168.0
  }'
```

---

### 5. Update Vital Signs

**Endpoint:** `PUT /api/v1/vital-signs/:id`

**Description:** Memperbarui data tanda vital berdasarkan ID. BMI dihitung ulang secara otomatis jika `weight_kg` atau `height_cm` diperbarui.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description    |
| --------- | ------- | -------- | -------------- |
| id        | integer | Yes      | ID tanda vital |

**Request Body:** (Same fields as Create, all optional for update)

**Example Request Body:**

```json
{
  "blood_pressure_systolic": 125,
  "blood_pressure_diastolic": 82,
  "temperature": 37.2,
  "oxygen_saturation": 97.8
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Vital signs updated successfully",
  "data": {
    "id": 1,
    "medical_record_id": 5,
    "blood_pressure_systolic": 125,
    "blood_pressure_diastolic": 82,
    "heart_rate": 75,
    "temperature": 37.2,
    "respiratory_rate": 18,
    "oxygen_saturation": 97.8,
    "weight_kg": 65.5,
    "height_cm": 168.0,
    "bmi": 23.21,
    "recorded_at": "2024-01-15T08:30:00Z",
    "updated_at": "2024-01-15T09:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT "http://localhost:8080/api/v1/vital-signs/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "blood_pressure_systolic": 125,
    "blood_pressure_diastolic": 82,
    "temperature": 37.2,
    "oxygen_saturation": 97.8
  }'
```

---

### 6. Soft Delete Vital Signs

**Endpoint:** `DELETE /api/v1/vital-signs/:id`

**Description:** Menghapus pencatatan tanda vital secara soft delete.

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description    |
| --------- | ------- | -------- | -------------- |
| id        | integer | Yes      | ID tanda vital |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Vital signs deleted successfully"
}
```

**cURL Example:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/vital-signs/1" \
  -H "Authorization: Bearer <token>"
```

---

### 7. Restore Deleted Vital Signs

**Endpoint:** `PATCH /api/v1/vital-signs/:id/restore`

**Description:** Memulihkan pencatatan tanda vital yang telah dihapus (soft delete).

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description    |
| --------- | ------- | -------- | -------------- |
| id        | integer | Yes      | ID tanda vital |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Vital signs restored successfully",
  "data": {
    "id": 1,
    "medical_record_id": 5,
    "recorded_at": "2024-01-15T08:30:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/vital-signs/1/restore" \
  -H "Authorization: Bearer <token>"
```

---

### 8. Hard Delete Vital Signs

**Endpoint:** `DELETE /api/v1/vital-signs/:id/hard-delete`

**Description:** Menghapus pencatatan tanda vital secara permanen dari database. Operasi ini tidak dapat dibatalkan.

**Authentication:** Required (Super Admin Only)

**Path Parameters:**

| Parameter | Type    | Required | Description    |
| --------- | ------- | -------- | -------------- |
| id        | integer | Yes      | ID tanda vital |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Vital signs permanently deleted successfully"
}
```

**cURL Example:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/vital-signs/1/hard-delete" \
  -H "Authorization: Bearer <token>"
```

---

## Error Responses

### 400 Bad Request

```json
{
  "status": "error",
  "message": "Validation error",
  "errors": {
    "medical_record_id": "medical_record_id is required",
    "recorded_at": "recorded_at is required"
  }
}
```

### 401 Unauthorized

```json
{
  "status": "error",
  "message": "Unauthorized: missing or invalid token"
}
```

### 403 Forbidden

```json
{
  "status": "error",
  "message": "Forbidden: insufficient permissions"
}
```

### 404 Not Found

```json
{
  "status": "error",
  "message": "Vital signs not found"
}
```

### 409 Conflict

```json
{
  "status": "error",
  "message": "Vital signs already recorded for this medical record"
}
```

### 500 Internal Server Error

```json
{
  "status": "error",
  "message": "Internal server error"
}
```

---

## Database Model

### Table: vital_signs

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| medical_record_id | BIGINT | FOREIGN KEY (medical_records.id), NOT NULL, INDEX | Reference ke medical record |
| measurement_date | DATE | NOT NULL | Tanggal pengukuran |
| measurement_time | TIME | NOT NULL | Waktu pengukuran |
| systolic_bp | INT | NULLABLE | Tekanan darah sistolik (mmHg) |
| diastolic_bp | INT | NULLABLE | Tekanan darah diastolik (mmHg) |
| heart_rate | INT | NULLABLE | Denyut jantung (bpm) |
| body_temperature | DECIMAL(4,2) | NULLABLE | Suhu tubuh (°C) |
| respiratory_rate | INT | NULLABLE | Laju pernafasan (breaths/min) |
| oxygen_saturation | DECIMAL(5,2) | NULLABLE | Saturasi oksigen (%) |
| weight_kg | DECIMAL(5,2) | NULLABLE | Berat badan (kg) |
| height_cm | INT | NULLABLE | Tinggi badan (cm) |
| bmi | DECIMAL(5,2) | NULLABLE | BMI (dihitung otomatis) |
| notes | TEXT | NULLABLE | Catatan khusus |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

**Indexes:**
- Primary Key: id
- Foreign Key: medical_record_id
- Regular Index: measurement_date, deleted_at

**Relationships:**
- Belongs To Medical Record (many-to-one)

**Notes:**
- BMI otomatis dihitung dari weight dan height
- Semua field vital signs NULLABLE (tidak semua harus diisi)
- Multiple measurements per day diperbolehkan
- Tracking tren vital signs selama perawatan

---

## Vital Signs Reference Ranges

| Parameter                | Normal Range         | Unit      |
| ------------------------ | -------------------- | --------- |
| blood_pressure_systolic  | 90 – 120             | mmHg      |
| blood_pressure_diastolic | 60 – 80              | mmHg      |
| heart_rate               | 60 – 100             | bpm       |
| temperature              | 36.1 – 37.2          | °C        |
| respiratory_rate         | 12 – 20              | per menit |
| oxygen_saturation        | ≥ 95                 | %         |
| bmi                      | 18.5 – 24.9 (normal) | kg/m²     |

**BMI Categories:**

- `< 18.5` — Berat badan kurang (Underweight)
- `18.5 – 24.9` — Normal
- `25.0 – 29.9` — Berat badan lebih (Overweight)
- `≥ 30.0` — Obesitas

> **Notes:**
>
> - Semua field selain `medical_record_id` dan `recorded_at` bersifat opsional — catat hanya parameter yang tersedia
> - BMI dihitung otomatis saat `weight_kg` dan `height_cm` keduanya diisi
> - Setiap rekam medis (`medical_record_id`) bersifat **unik** — hanya dapat memiliki satu entri vital signs
> - Tekanan darah dicatat sebagai dua nilai terpisah: sistolik (`blood_pressure_systolic`) dan diastolik (`blood_pressure_diastolic`)
