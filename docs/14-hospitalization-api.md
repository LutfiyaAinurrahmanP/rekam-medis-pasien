# Hospitalization API Documentation

## Overview

API untuk manajemen data rawat inap (hospitalization) dalam sistem rekam medis. Hospitalization mencatat data pasien yang menjalani rawat inap, meliputi tanggal masuk, kamar yang ditempati, dokter yang merawat, alasan masuk, dan status kepulangan.

**Base URL:** `/api/v1/hospitalizations`

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

| Endpoint                                  | Patient | Doctor | Receptionist | Admin | Super Admin |
| ----------------------------------------- | ------- | ------ | ------------ | ----- | ----------- |
| GET /hospitalizations                     | ❌      | ✅     | ✅           | ✅    | ✅          |
| GET /hospitalizations/deleted             | ❌      | ❌     | ❌           | ✅    | ✅          |
| GET /hospitalizations/:id                 | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /hospitalizations/patient/:patient_id | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /hospitalizations/active              | ❌      | ✅     | ✅           | ✅    | ✅          |
| POST /hospitalizations                    | ❌      | ❌     | ✅           | ✅    | ✅          |
| PUT /hospitalizations/:id                 | ❌      | ✅     | ✅           | ✅    | ✅          |
| PATCH /hospitalizations/:id/discharge     | ❌      | ✅     | ✅           | ✅    | ✅          |
| PATCH /hospitalizations/:id/transfer      | ❌      | ✅     | ✅           | ✅    | ✅          |
| DELETE /hospitalizations/:id              | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /hospitalizations/:id/restore       | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /hospitalizations/:id/hard-delete  | ❌      | ❌     | ❌           | ❌    | ✅          |

---

## Endpoints Summary

| Method | Endpoint                                | Description                             | Role Required                            |
| ------ | --------------------------------------- | --------------------------------------- | ---------------------------------------- |
| GET    | `/hospitalizations`                     | List all hospitalizations               | Doctor, Receptionist, Admin, Super Admin |
| GET    | `/hospitalizations/deleted`             | List deleted hospitalizations           | Admin, Super Admin                       |
| GET    | `/hospitalizations/active`              | List active (admitted) hospitalizations | Doctor, Receptionist, Admin, Super Admin |
| GET    | `/hospitalizations/:id`                 | Get hospitalization by ID               | All Authenticated                        |
| GET    | `/hospitalizations/patient/:patient_id` | Get hospitalizations by patient         | All Authenticated                        |
| POST   | `/hospitalizations`                     | Create new hospitalization              | Receptionist, Admin, Super Admin         |
| PUT    | `/hospitalizations/:id`                 | Update hospitalization                  | Doctor, Receptionist, Admin, Super Admin |
| PATCH  | `/hospitalizations/:id/discharge`       | Discharge patient                       | Doctor, Receptionist, Admin, Super Admin |
| PATCH  | `/hospitalizations/:id/transfer`        | Transfer patient to another facility    | Doctor, Receptionist, Admin, Super Admin |
| DELETE | `/hospitalizations/:id`                 | Soft delete hospitalization             | Admin, Super Admin                       |
| PATCH  | `/hospitalizations/:id/restore`         | Restore deleted hospitalization         | Admin, Super Admin                       |
| DELETE | `/hospitalizations/:id/hard-delete`     | Permanently delete hospitalization      | Super Admin                              |

---

## Endpoints Detail

### 1. List All Hospitalizations

**Endpoint:** `GET /api/v1/hospitalizations`

**Description:** Mendapatkan daftar semua data rawat inap dengan pagination, search, dan filter.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

| Header        | Type   | Required | Description      |
| ------------- | ------ | -------- | ---------------- |
| Authorization | string | Yes      | `Bearer <token>` |

**Query Parameters:**

| Parameter  | Type    | Required | Default    | Description                                                                     |
| ---------- | ------- | -------- | ---------- | ------------------------------------------------------------------------------- |
| page       | integer | No       | 1          | Halaman saat ini                                                                |
| page_size  | integer | No       | 10         | Jumlah data per halaman (max: 100)                                              |
| sort_by    | string  | No       | created_at | Field sorting: `id`, `admission_date`, `discharge_date`, `status`, `created_at` |
| sort_dir   | string  | No       | desc       | Arah sorting: `asc` atau `desc`                                                 |
| search     | string  | No       | -          | Pencarian berdasarkan alasan masuk atau ringkasan keluar                        |
| status     | string  | No       | -          | Filter status: `admitted`, `discharged`, `transferred`                          |
| patient_id | integer | No       | -          | Filter berdasarkan ID pasien                                                    |
| room_id    | integer | No       | -          | Filter berdasarkan ID kamar                                                     |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Hospitalizations retrieved successfully",
  "data": [
    {
      "id": 1,
      "patient_id": 10,
      "medical_record_id": 5,
      "room_id": 3,
      "attending_doctor_id": 2,
      "admission_date": "2024-01-15T08:00:00Z",
      "discharge_date": null,
      "admission_reason": "Demam tinggi dan sesak napas",
      "discharge_summary": null,
      "status": "admitted",
      "created_at": "2024-01-15T08:00:00Z",
      "updated_at": "2024-01-15T08:00:00Z"
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
curl -X GET "http://localhost:8080/api/v1/hospitalizations?page=1&page_size=10&status=admitted" \
  -H "Authorization: Bearer <token>"
```

---

### 2. List Deleted Hospitalizations

**Endpoint:** `GET /api/v1/hospitalizations/deleted`

**Description:** Mendapatkan daftar data rawat inap yang telah dihapus (soft deleted).

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

| Header        | Type   | Required | Description      |
| ------------- | ------ | -------- | ---------------- |
| Authorization | string | Yes      | `Bearer <token>` |

**Query Parameters:**

| Parameter | Type    | Required | Default    | Description                        |
| --------- | ------- | -------- | ---------- | ---------------------------------- |
| page      | integer | No       | 1          | Halaman saat ini                   |
| page_size | integer | No       | 10         | Jumlah data per halaman (max: 100) |
| sort_by   | string  | No       | created_at | Field untuk sorting                |
| sort_dir  | string  | No       | desc       | Arah sorting: `asc` atau `desc`    |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Deleted hospitalizations retrieved successfully",
  "data": [
    {
      "id": 2,
      "patient_id": 7,
      "medical_record_id": 3,
      "status": "discharged",
      "admission_date": "2024-01-01T08:00:00Z",
      "discharge_date": "2024-01-05T10:00:00Z"
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
curl -X GET "http://localhost:8080/api/v1/hospitalizations/deleted" \
  -H "Authorization: Bearer <token>"
```

---

### 3. List Active Hospitalizations

**Endpoint:** `GET /api/v1/hospitalizations/active`

**Description:** Mendapatkan daftar pasien yang sedang menjalani rawat inap (status `admitted`).

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Query Parameters:**

| Parameter | Type    | Required | Default | Description                        |
| --------- | ------- | -------- | ------- | ---------------------------------- |
| page      | integer | No       | 1       | Halaman saat ini                   |
| page_size | integer | No       | 10      | Jumlah data per halaman (max: 100) |
| room_id   | integer | No       | -       | Filter berdasarkan ID kamar        |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Active hospitalizations retrieved successfully",
  "data": [
    {
      "id": 1,
      "patient_id": 10,
      "room_id": 3,
      "attending_doctor_id": 2,
      "admission_date": "2024-01-15T08:00:00Z",
      "admission_reason": "Demam tinggi dan sesak napas",
      "status": "admitted",
      "patient": {
        "id": 10,
        "name": "Budi Santoso",
        "medical_record_number": "MRN-2024-000010"
      },
      "room": {
        "id": 3,
        "room_number": "302",
        "room_type": "standard"
      },
      "attending_doctor": {
        "id": 2,
        "name": "dr. Siti Rahayu",
        "specialization": "Penyakit Dalam"
      }
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
curl -X GET "http://localhost:8080/api/v1/hospitalizations/active?room_id=3" \
  -H "Authorization: Bearer <token>"
```

---

### 4. Get Hospitalization by ID

**Endpoint:** `GET /api/v1/hospitalizations/:id`

**Description:** Mendapatkan detail data rawat inap berdasarkan ID.

**Authentication:** Required (All Authenticated Users)

**Path Parameters:**

| Parameter | Type    | Required | Description   |
| --------- | ------- | -------- | ------------- |
| id        | integer | Yes      | ID rawat inap |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Hospitalization retrieved successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "medical_record_id": 5,
    "room_id": 3,
    "attending_doctor_id": 2,
    "admission_date": "2024-01-15T08:00:00Z",
    "discharge_date": null,
    "admission_reason": "Demam tinggi dan sesak napas, diduga pneumonia",
    "discharge_summary": null,
    "status": "admitted",
    "created_at": "2024-01-15T08:00:00Z",
    "updated_at": "2024-01-15T08:00:00Z",
    "patient": {
      "id": 10,
      "name": "Budi Santoso",
      "date_of_birth": "1985-03-20",
      "medical_record_number": "MRN-2024-000010"
    },
    "medical_record": {
      "id": 5,
      "visit_date": "2024-01-15",
      "chief_complaint": "Demam dan sesak napas"
    },
    "room": {
      "id": 3,
      "room_number": "302",
      "room_type": "standard",
      "floor": 3
    },
    "attending_doctor": {
      "id": 2,
      "name": "dr. Siti Rahayu",
      "specialization": "Penyakit Dalam"
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/hospitalizations/1" \
  -H "Authorization: Bearer <token>"
```

---

### 5. Get Hospitalizations by Patient

**Endpoint:** `GET /api/v1/hospitalizations/patient/:patient_id`

**Description:** Mendapatkan riwayat rawat inap berdasarkan ID pasien.

**Authentication:** Required (All Authenticated Users)

**Path Parameters:**

| Parameter  | Type    | Required | Description |
| ---------- | ------- | -------- | ----------- |
| patient_id | integer | Yes      | ID pasien   |

**Query Parameters:**

| Parameter | Type    | Required | Default | Description                                            |
| --------- | ------- | -------- | ------- | ------------------------------------------------------ |
| page      | integer | No       | 1       | Halaman saat ini                                       |
| page_size | integer | No       | 10      | Jumlah data per halaman (max: 100)                     |
| status    | string  | No       | -       | Filter status: `admitted`, `discharged`, `transferred` |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Patient hospitalizations retrieved successfully",
  "data": [
    {
      "id": 1,
      "medical_record_id": 5,
      "room_id": 3,
      "attending_doctor_id": 2,
      "admission_date": "2024-01-15T08:00:00Z",
      "discharge_date": "2024-01-20T10:00:00Z",
      "admission_reason": "Demam tinggi dan sesak napas",
      "status": "discharged"
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
curl -X GET "http://localhost:8080/api/v1/hospitalizations/patient/10?status=discharged" \
  -H "Authorization: Bearer <token>"
```

---

### 6. Create Hospitalization

**Endpoint:** `POST /api/v1/hospitalizations`

**Description:** Membuat data rawat inap baru untuk pasien.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Request Headers:**

| Header        | Type   | Required | Description        |
| ------------- | ------ | -------- | ------------------ |
| Authorization | string | Yes      | `Bearer <token>`   |
| Content-Type  | string | Yes      | `application/json` |

**Request Body:**

| Field               | Type    | Required | Description                                                           |
| ------------------- | ------- | -------- | --------------------------------------------------------------------- |
| patient_id          | integer | Yes      | ID pasien                                                             |
| medical_record_id   | integer | Yes      | ID rekam medis terkait                                                |
| room_id             | integer | Yes      | ID kamar yang ditempati                                               |
| attending_doctor_id | integer | Yes      | ID dokter yang merawat                                                |
| admission_date      | string  | Yes      | Tanggal dan waktu masuk (format: `YYYY-MM-DDTHH:MM:SSZ`)              |
| admission_reason    | string  | Yes      | Alasan/indikasi rawat inap                                            |
| status              | string  | No       | Status: `admitted`, `discharged`, `transferred` (default: `admitted`) |

**Example Request Body:**

```json
{
  "patient_id": 10,
  "medical_record_id": 5,
  "room_id": 3,
  "attending_doctor_id": 2,
  "admission_date": "2024-01-15T08:00:00Z",
  "admission_reason": "Demam tinggi 39.5°C dan sesak napas, diduga pneumonia komunitas"
}
```

**Response Success (201):**

```json
{
  "status": "success",
  "message": "Hospitalization created successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "medical_record_id": 5,
    "room_id": 3,
    "attending_doctor_id": 2,
    "admission_date": "2024-01-15T08:00:00Z",
    "discharge_date": null,
    "admission_reason": "Demam tinggi 39.5°C dan sesak napas, diduga pneumonia komunitas",
    "discharge_summary": null,
    "status": "admitted",
    "created_at": "2024-01-15T08:00:00Z",
    "updated_at": "2024-01-15T08:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/hospitalizations" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 10,
    "medical_record_id": 5,
    "room_id": 3,
    "attending_doctor_id": 2,
    "admission_date": "2024-01-15T08:00:00Z",
    "admission_reason": "Demam tinggi 39.5°C dan sesak napas, diduga pneumonia komunitas"
  }'
```

---

### 7. Update Hospitalization

**Endpoint:** `PUT /api/v1/hospitalizations/:id`

**Description:** Memperbarui data rawat inap berdasarkan ID.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description   |
| --------- | ------- | -------- | ------------- |
| id        | integer | Yes      | ID rawat inap |

**Request Body:** (Same fields as Create, all optional for update)

**Example Request Body:**

```json
{
  "room_id": 4,
  "attending_doctor_id": 3,
  "admission_reason": "Pneumonia komunitas dengan komplikasi pleural effusion"
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Hospitalization updated successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "room_id": 4,
    "attending_doctor_id": 3,
    "admission_reason": "Pneumonia komunitas dengan komplikasi pleural effusion",
    "status": "admitted",
    "updated_at": "2024-01-16T09:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT "http://localhost:8080/api/v1/hospitalizations/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "room_id": 4,
    "attending_doctor_id": 3,
    "admission_reason": "Pneumonia komunitas dengan komplikasi pleural effusion"
  }'
```

---

### 8. Discharge Patient

**Endpoint:** `PATCH /api/v1/hospitalizations/:id/discharge`

**Description:** Memulangkan pasien rawat inap. Status berubah menjadi `discharged` dan `discharge_date` diisi otomatis.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description   |
| --------- | ------- | -------- | ------------- |
| id        | integer | Yes      | ID rawat inap |

**Request Body:**

| Field             | Type   | Required | Description                   |
| ----------------- | ------ | -------- | ----------------------------- |
| discharge_summary | string | No       | Ringkasan kondisi saat keluar |

**Example Request Body:**

```json
{
  "discharge_summary": "Pasien membaik setelah 5 hari perawatan. Pneumonia teratasi, sesak napas dan demam sudah tidak ada. Pasien diperbolehkan pulang dengan obat oral selama 7 hari."
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Patient discharged successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "admission_date": "2024-01-15T08:00:00Z",
    "discharge_date": "2024-01-20T10:00:00Z",
    "discharge_summary": "Pasien membaik setelah 5 hari perawatan...",
    "status": "discharged",
    "updated_at": "2024-01-20T10:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/hospitalizations/1/discharge" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "discharge_summary": "Pasien membaik setelah 5 hari perawatan. Pneumonia teratasi."
  }'
```

> **Notes:**
>
> - `discharge_date` diisi otomatis dengan waktu saat endpoint ini dipanggil
> - Durasi rawat inap dapat dihitung dari `admission_date` hingga `discharge_date`
> - Minimal durasi yang dicatat adalah 1 hari

---

### 9. Transfer Patient

**Endpoint:** `PATCH /api/v1/hospitalizations/:id/transfer`

**Description:** Mentransfer pasien ke fasilitas/kamar lain. Status berubah menjadi `transferred`.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description   |
| --------- | ------- | -------- | ------------- |
| id        | integer | Yes      | ID rawat inap |

**Request Body:**

| Field | Type   | Required | Description                    |
| ----- | ------ | -------- | ------------------------------ |
| notes | string | No       | Catatan alasan/tujuan transfer |

**Example Request Body:**

```json
{
  "notes": "Pasien ditransfer ke ICU karena kondisi memburuk"
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Patient transferred successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "status": "transferred",
    "updated_at": "2024-01-17T14:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/hospitalizations/1/transfer" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Pasien ditransfer ke ICU karena kondisi memburuk"
  }'
```

---

### 10. Soft Delete Hospitalization

**Endpoint:** `DELETE /api/v1/hospitalizations/:id`

**Description:** Menghapus data rawat inap secara soft delete.

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description   |
| --------- | ------- | -------- | ------------- |
| id        | integer | Yes      | ID rawat inap |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Hospitalization deleted successfully"
}
```

**cURL Example:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/hospitalizations/1" \
  -H "Authorization: Bearer <token>"
```

---

### 11. Restore Deleted Hospitalization

**Endpoint:** `PATCH /api/v1/hospitalizations/:id/restore`

**Description:** Memulihkan data rawat inap yang telah dihapus (soft delete).

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description   |
| --------- | ------- | -------- | ------------- |
| id        | integer | Yes      | ID rawat inap |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Hospitalization restored successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "status": "admitted"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/hospitalizations/1/restore" \
  -H "Authorization: Bearer <token>"
```

---

### 12. Hard Delete Hospitalization

**Endpoint:** `DELETE /api/v1/hospitalizations/:id/hard-delete`

**Description:** Menghapus data rawat inap secara permanen dari database. Operasi ini tidak dapat dibatalkan.

**Authentication:** Required (Super Admin Only)

**Path Parameters:**

| Parameter | Type    | Required | Description   |
| --------- | ------- | -------- | ------------- |
| id        | integer | Yes      | ID rawat inap |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Hospitalization permanently deleted successfully"
}
```

**cURL Example:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/hospitalizations/1/hard-delete" \
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
    "patient_id": "patient_id is required",
    "medical_record_id": "medical_record_id is required",
    "room_id": "room_id is required",
    "attending_doctor_id": "attending_doctor_id is required",
    "admission_reason": "admission_reason is required"
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
  "message": "Hospitalization not found"
}
```

### 422 Unprocessable Entity

```json
{
  "status": "error",
  "message": "Patient has already been discharged"
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

### Table: hospitalizations

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| patient_id | BIGINT | FOREIGN KEY (patients.id), NOT NULL, INDEX | Reference ke pasien |
| doctor_id | BIGINT | FOREIGN KEY (doctors.id), NOT NULL, INDEX | Reference ke dokter penanggung jawab |
| room_id | BIGINT | FOREIGN KEY (rooms.id), NOT NULL, INDEX | Reference ke ruangan |
| admission_date | DATE | NOT NULL | Tanggal masuk rawat inap |
| admission_time | TIME | NOT NULL | Waktu masuk |
| discharge_date | DATE | NULLABLE, INDEX | Tanggal pulang (NULL jika masih dirawat) |
| discharge_time | TIME | NULLABLE | Waktu pulang |
| reason_for_admission | TEXT | NOT NULL | Alasan/diagnosis saat masuk |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'admitted', INDEX | Status (admitted, transferred, discharged, cancelled) |
| notes | TEXT | NULLABLE | Catatan tambahan |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

**Indexes:**
- Primary Key: id
- Foreign Keys: patient_id, doctor_id, room_id
- Regular Index: admission_date, discharge_date, status, deleted_at

**Relationships:**
- Belongs To Patient (many-to-one)
- Belongs To Doctor (many-to-one)
- Belongs To Room (many-to-one)
- Has Many Medical Records (one-to-many)
- Has Many Vital Signs (one-to-many)

**Notes:**
- Status: admitted (sedang dirawat) -> transferred (pindah ruangan) -> discharged (pulang)
- discharge_date NULL berarti pasien masih dirawat
- Integrasi dengan room availability/occupancy
- Tracking perjalanan pasien dalam hospital

---

## Hospitalization Status Flow

```
admitted → discharged
admitted → transferred
```

**Status Values:**

- `admitted` — Pasien sedang menjalani rawat inap
- `discharged` — Pasien telah dipulangkan
- `transferred` — Pasien ditransfer ke fasilitas/unit lain

> **Notes:**
>
> - Satu pasien dapat memiliki lebih dari satu riwayat rawat inap
> - Durasi rawat inap dihitung otomatis dari `admission_date` hingga `discharge_date` (minimum 1 hari)
> - Kamar yang digunakan pasien akan terpengaruh ketersediaannya selama status `admitted`
