# Lab Tests API Documentation

## Overview

API untuk manajemen data pemeriksaan laboratorium (lab tests) dalam sistem rekam medis. Lab Tests mencatat permintaan pemeriksaan laboratorium yang diorder oleh dokter berdasarkan rekam medis pasien, beserta status pengerjaan, hasil pemeriksaan, dan nilai referensi.

**Base URL:** `/api/v1/lab-tests`

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

| Endpoint                                 | Patient | Doctor | Receptionist | Admin | Super Admin |
| ---------------------------------------- | ------- | ------ | ------------ | ----- | ----------- |
| GET /lab-tests                           | ❌      | ✅     | ✅           | ✅    | ✅          |
| GET /lab-tests/deleted                   | ❌      | ❌     | ❌           | ✅    | ✅          |
| GET /lab-tests/:id                       | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /lab-tests/medical-record/:record_id | ✅      | ✅     | ✅           | ✅    | ✅          |
| POST /lab-tests                          | ❌      | ✅     | ❌           | ✅    | ✅          |
| PUT /lab-tests/:id                       | ❌      | ✅     | ❌           | ✅    | ✅          |
| PATCH /lab-tests/:id/collect-sample      | ❌      | ❌     | ✅           | ✅    | ✅          |
| PATCH /lab-tests/:id/start               | ❌      | ❌     | ✅           | ✅    | ✅          |
| PATCH /lab-tests/:id/complete            | ❌      | ❌     | ✅           | ✅    | ✅          |
| PATCH /lab-tests/:id/cancel              | ❌      | ✅     | ✅           | ✅    | ✅          |
| DELETE /lab-tests/:id                    | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /lab-tests/:id/restore             | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /lab-tests/:id/hard-delete        | ❌      | ❌     | ❌           | ❌    | ✅          |

---

## Endpoints Summary

| Method | Endpoint                               | Description                     | Role Required                            |
| ------ | -------------------------------------- | ------------------------------- | ---------------------------------------- |
| GET    | `/lab-tests`                           | List all lab tests              | Doctor, Receptionist, Admin, Super Admin |
| GET    | `/lab-tests/deleted`                   | List deleted lab tests          | Admin, Super Admin                       |
| GET    | `/lab-tests/:id`                       | Get lab test by ID              | All Authenticated                        |
| GET    | `/lab-tests/medical-record/:record_id` | Get lab tests by medical record | All Authenticated                        |
| POST   | `/lab-tests`                           | Create/order new lab test       | Doctor, Admin, Super Admin               |
| PUT    | `/lab-tests/:id`                       | Update lab test                 | Doctor, Admin, Super Admin               |
| PATCH  | `/lab-tests/:id/collect-sample`        | Mark sample as collected        | Receptionist, Admin, Super Admin         |
| PATCH  | `/lab-tests/:id/start`                 | Mark lab test as in-progress    | Receptionist, Admin, Super Admin         |
| PATCH  | `/lab-tests/:id/complete`              | Complete lab test with results  | Receptionist, Admin, Super Admin         |
| PATCH  | `/lab-tests/:id/cancel`                | Cancel lab test                 | Doctor, Receptionist, Admin, Super Admin |
| DELETE | `/lab-tests/:id`                       | Soft delete lab test            | Admin, Super Admin                       |
| PATCH  | `/lab-tests/:id/restore`               | Restore deleted lab test        | Admin, Super Admin                       |
| DELETE | `/lab-tests/:id/hard-delete`           | Permanently delete lab test     | Super Admin                              |

---

## Endpoints Detail

### 1. List All Lab Tests

**Endpoint:** `GET /api/v1/lab-tests`

**Description:** Mendapatkan daftar semua pemeriksaan laboratorium dengan pagination, search, dan filter.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

| Header        | Type   | Required | Description      |
| ------------- | ------ | -------- | ---------------- |
| Authorization | string | Yes      | `Bearer <token>` |

**Query Parameters:**

| Parameter            | Type    | Required | Default    | Description                                                                           |
| -------------------- | ------- | -------- | ---------- | ------------------------------------------------------------------------------------- |
| page                 | integer | No       | 1          | Halaman saat ini                                                                      |
| page_size            | integer | No       | 10         | Jumlah data per halaman (max: 100)                                                    |
| sort_by              | string  | No       | created_at | Field sorting: `id`, `order_date`, `result_date`, `status`, `created_at`              |
| sort_dir             | string  | No       | desc       | Arah sorting: `asc` atau `desc`                                                       |
| search               | string  | No       | -          | Pencarian berdasarkan catatan atau nilai hasil                                        |
| status               | string  | No       | -          | Filter status: `ordered`, `sample_collected`, `in_progress`, `completed`, `cancelled` |
| test_type_id         | integer | No       | -          | Filter berdasarkan ID jenis pemeriksaan                                               |
| ordered_by_doctor_id | integer | No       | -          | Filter berdasarkan ID dokter yang memesan                                             |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Lab tests retrieved successfully",
  "data": [
    {
      "id": 1,
      "medical_record_id": 5,
      "test_type_id": 3,
      "ordered_by_doctor_id": 2,
      "order_date": "2024-01-15",
      "sample_collection_date": null,
      "result_date": null,
      "result_value": null,
      "result_unit": null,
      "reference_range": null,
      "status": "ordered",
      "notes": "Segera",
      "created_at": "2024-01-15T09:00:00Z",
      "updated_at": "2024-01-15T09:00:00Z"
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
curl -X GET "http://localhost:8080/api/v1/lab-tests?page=1&page_size=10&status=ordered" \
  -H "Authorization: Bearer <token>"
```

---

### 2. List Deleted Lab Tests

**Endpoint:** `GET /api/v1/lab-tests/deleted`

**Description:** Mendapatkan daftar pemeriksaan laboratorium yang telah dihapus (soft deleted).

**Authentication:** Required (Admin, Super Admin)

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
  "message": "Deleted lab tests retrieved successfully",
  "data": [
    {
      "id": 2,
      "medical_record_id": 3,
      "test_type_id": 1,
      "order_date": "2024-01-10",
      "status": "cancelled"
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
curl -X GET "http://localhost:8080/api/v1/lab-tests/deleted" \
  -H "Authorization: Bearer <token>"
```

---

### 3. Get Lab Test by ID

**Endpoint:** `GET /api/v1/lab-tests/:id`

**Description:** Mendapatkan detail data pemeriksaan laboratorium berdasarkan ID.

**Authentication:** Required (All Authenticated Users)

**Path Parameters:**

| Parameter | Type    | Required | Description        |
| --------- | ------- | -------- | ------------------ |
| id        | integer | Yes      | ID pemeriksaan lab |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Lab test retrieved successfully",
  "data": {
    "id": 1,
    "medical_record_id": 5,
    "test_type_id": 3,
    "ordered_by_doctor_id": 2,
    "order_date": "2024-01-15",
    "sample_collection_date": "2024-01-15T10:00:00Z",
    "result_date": "2024-01-15T14:00:00Z",
    "result_value": "12.5",
    "result_unit": "g/dL",
    "reference_range": "12.0 - 16.0 g/dL",
    "status": "completed",
    "notes": "Segera",
    "created_at": "2024-01-15T09:00:00Z",
    "updated_at": "2024-01-15T14:00:00Z",
    "medical_record": {
      "id": 5,
      "visit_date": "2024-01-15",
      "chief_complaint": "Demam dan lemas"
    },
    "test_type": {
      "id": 3,
      "name": "Hemoglobin",
      "code": "HGB",
      "category": "hematology"
    },
    "ordered_by_doctor": {
      "id": 2,
      "name": "dr. Siti Rahayu",
      "specialization": "Penyakit Dalam"
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/lab-tests/1" \
  -H "Authorization: Bearer <token>"
```

---

### 4. Get Lab Tests by Medical Record

**Endpoint:** `GET /api/v1/lab-tests/medical-record/:record_id`

**Description:** Mendapatkan daftar pemeriksaan laboratorium berdasarkan ID rekam medis.

**Authentication:** Required (All Authenticated Users)

**Path Parameters:**

| Parameter | Type    | Required | Description    |
| --------- | ------- | -------- | -------------- |
| record_id | integer | Yes      | ID rekam medis |

**Query Parameters:**

| Parameter | Type   | Required | Default | Description                                                                           |
| --------- | ------ | -------- | ------- | ------------------------------------------------------------------------------------- |
| status    | string | No       | -       | Filter status: `ordered`, `sample_collected`, `in_progress`, `completed`, `cancelled` |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Medical record lab tests retrieved successfully",
  "data": [
    {
      "id": 1,
      "test_type_id": 3,
      "ordered_by_doctor_id": 2,
      "order_date": "2024-01-15",
      "result_value": "12.5",
      "result_unit": "g/dL",
      "reference_range": "12.0 - 16.0 g/dL",
      "status": "completed",
      "test_type": {
        "id": 3,
        "name": "Hemoglobin",
        "code": "HGB"
      }
    }
  ]
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/lab-tests/medical-record/5?status=completed" \
  -H "Authorization: Bearer <token>"
```

---

### 5. Create Lab Test (Order)

**Endpoint:** `POST /api/v1/lab-tests`

**Description:** Memesan/membuat permintaan pemeriksaan laboratorium baru untuk pasien.

**Authentication:** Required (Doctor, Admin, Super Admin)

**Request Headers:**

| Header        | Type   | Required | Description        |
| ------------- | ------ | -------- | ------------------ |
| Authorization | string | Yes      | `Bearer <token>`   |
| Content-Type  | string | Yes      | `application/json` |

**Request Body:**

| Field                | Type    | Required | Description                          |
| -------------------- | ------- | -------- | ------------------------------------ |
| medical_record_id    | integer | Yes      | ID rekam medis terkait               |
| test_type_id         | integer | Yes      | ID jenis pemeriksaan laboratorium    |
| ordered_by_doctor_id | integer | Yes      | ID dokter yang memesan               |
| order_date           | string  | Yes      | Tanggal order (format: `YYYY-MM-DD`) |
| status               | string  | No       | Status awal (default: `ordered`)     |
| notes                | string  | No       | Catatan atau instruksi khusus        |

**Example Request Body:**

```json
{
  "medical_record_id": 5,
  "test_type_id": 3,
  "ordered_by_doctor_id": 2,
  "order_date": "2024-01-15",
  "notes": "Segera, pasien diduga anemia"
}
```

**Response Success (201):**

```json
{
  "status": "success",
  "message": "Lab test ordered successfully",
  "data": {
    "id": 1,
    "medical_record_id": 5,
    "test_type_id": 3,
    "ordered_by_doctor_id": 2,
    "order_date": "2024-01-15",
    "sample_collection_date": null,
    "result_date": null,
    "result_value": null,
    "result_unit": null,
    "reference_range": null,
    "status": "ordered",
    "notes": "Segera, pasien diduga anemia",
    "created_at": "2024-01-15T09:00:00Z",
    "updated_at": "2024-01-15T09:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/lab-tests" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "medical_record_id": 5,
    "test_type_id": 3,
    "ordered_by_doctor_id": 2,
    "order_date": "2024-01-15",
    "notes": "Segera, pasien diduga anemia"
  }'
```

---

### 6. Update Lab Test

**Endpoint:** `PUT /api/v1/lab-tests/:id`

**Description:** Memperbarui data pemeriksaan laboratorium. Umumnya digunakan untuk mengupdate catatan atau informasi dasar.

**Authentication:** Required (Doctor, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description        |
| --------- | ------- | -------- | ------------------ |
| id        | integer | Yes      | ID pemeriksaan lab |

**Request Body:** (Same fields as Create, all optional for update)

**Example Request Body:**

```json
{
  "notes": "Pasien meminta hasil cepat untuk keperluan operasi",
  "reference_range": "12.0 - 16.0 g/dL"
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Lab test updated successfully",
  "data": {
    "id": 1,
    "notes": "Pasien meminta hasil cepat untuk keperluan operasi",
    "reference_range": "12.0 - 16.0 g/dL",
    "updated_at": "2024-01-15T09:30:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT "http://localhost:8080/api/v1/lab-tests/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Pasien meminta hasil cepat untuk keperluan operasi"
  }'
```

---

### 7. Mark Sample as Collected

**Endpoint:** `PATCH /api/v1/lab-tests/:id/collect-sample`

**Description:** Menandai bahwa sampel pasien telah diambil. Status berubah menjadi `sample_collected` dan `sample_collection_date` diisi otomatis.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description        |
| --------- | ------- | -------- | ------------------ |
| id        | integer | Yes      | ID pemeriksaan lab |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Sample collected successfully",
  "data": {
    "id": 1,
    "status": "sample_collected",
    "sample_collection_date": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/lab-tests/1/collect-sample" \
  -H "Authorization: Bearer <token>"
```

---

### 8. Mark Lab Test as In-Progress

**Endpoint:** `PATCH /api/v1/lab-tests/:id/start`

**Description:** Menandai bahwa pemeriksaan laboratorium sedang dikerjakan. Status berubah menjadi `in_progress`.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description        |
| --------- | ------- | -------- | ------------------ |
| id        | integer | Yes      | ID pemeriksaan lab |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Lab test started successfully",
  "data": {
    "id": 1,
    "status": "in_progress",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/lab-tests/1/start" \
  -H "Authorization: Bearer <token>"
```

---

### 9. Complete Lab Test

**Endpoint:** `PATCH /api/v1/lab-tests/:id/complete`

**Description:** Menyelesaikan pemeriksaan laboratorium dan mengisi hasil pemeriksaan. Status berubah menjadi `completed` dan `result_date` diisi otomatis.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description        |
| --------- | ------- | -------- | ------------------ |
| id        | integer | Yes      | ID pemeriksaan lab |

**Request Body:**

| Field           | Type   | Required | Description                             |
| --------------- | ------ | -------- | --------------------------------------- |
| result_value    | string | No       | Nilai hasil pemeriksaan                 |
| result_unit     | string | No       | Satuan hasil (max 50 karakter)          |
| reference_range | string | No       | Rentang nilai normal (max 100 karakter) |
| notes           | string | No       | Catatan hasil pemeriksaan               |

**Example Request Body:**

```json
{
  "result_value": "12.5",
  "result_unit": "g/dL",
  "reference_range": "12.0 - 16.0 g/dL",
  "notes": "Kadar hemoglobin dalam batas normal"
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Lab test completed successfully",
  "data": {
    "id": 1,
    "result_value": "12.5",
    "result_unit": "g/dL",
    "reference_range": "12.0 - 16.0 g/dL",
    "result_date": "2024-01-15T14:00:00Z",
    "status": "completed",
    "notes": "Kadar hemoglobin dalam batas normal",
    "updated_at": "2024-01-15T14:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/lab-tests/1/complete" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "result_value": "12.5",
    "result_unit": "g/dL",
    "reference_range": "12.0 - 16.0 g/dL",
    "notes": "Kadar hemoglobin dalam batas normal"
  }'
```

---

### 10. Cancel Lab Test

**Endpoint:** `PATCH /api/v1/lab-tests/:id/cancel`

**Description:** Membatalkan pemeriksaan laboratorium. Status berubah menjadi `cancelled`.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description        |
| --------- | ------- | -------- | ------------------ |
| id        | integer | Yes      | ID pemeriksaan lab |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Lab test cancelled successfully",
  "data": {
    "id": 1,
    "status": "cancelled",
    "updated_at": "2024-01-15T10:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/lab-tests/1/cancel" \
  -H "Authorization: Bearer <token>"
```

---

### 11. Soft Delete Lab Test

**Endpoint:** `DELETE /api/v1/lab-tests/:id`

**Description:** Menghapus pemeriksaan laboratorium secara soft delete.

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description        |
| --------- | ------- | -------- | ------------------ |
| id        | integer | Yes      | ID pemeriksaan lab |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Lab test deleted successfully"
}
```

**cURL Example:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/lab-tests/1" \
  -H "Authorization: Bearer <token>"
```

---

### 12. Restore Deleted Lab Test

**Endpoint:** `PATCH /api/v1/lab-tests/:id/restore`

**Description:** Memulihkan pemeriksaan laboratorium yang telah dihapus (soft delete).

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description        |
| --------- | ------- | -------- | ------------------ |
| id        | integer | Yes      | ID pemeriksaan lab |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Lab test restored successfully",
  "data": {
    "id": 1,
    "status": "ordered"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/lab-tests/1/restore" \
  -H "Authorization: Bearer <token>"
```

---

### 13. Hard Delete Lab Test

**Endpoint:** `DELETE /api/v1/lab-tests/:id/hard-delete`

**Description:** Menghapus pemeriksaan laboratorium secara permanen dari database. Operasi ini tidak dapat dibatalkan.

**Authentication:** Required (Super Admin Only)

**Path Parameters:**

| Parameter | Type    | Required | Description        |
| --------- | ------- | -------- | ------------------ |
| id        | integer | Yes      | ID pemeriksaan lab |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Lab test permanently deleted successfully"
}
```

**cURL Example:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/lab-tests/1/hard-delete" \
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
    "test_type_id": "test_type_id is required",
    "ordered_by_doctor_id": "ordered_by_doctor_id is required",
    "order_date": "order_date is required"
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
  "message": "Lab test not found"
}
```

### 422 Unprocessable Entity

```json
{
  "status": "error",
  "message": "Cannot collect sample for a cancelled lab test"
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

### Table: lab_tests

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| medical_record_id | BIGINT | FOREIGN KEY (medical_records.id), NOT NULL, INDEX | Reference ke medical record |
| test_type_id | BIGINT | FOREIGN KEY (type_tests.id), NOT NULL, INDEX | Reference ke jenis tes |
| sample_collection_date | DATE | NULLABLE | Tanggal pengambilan sampel |
| test_start_date | DATE | NULLABLE | Tanggal mulai pemeriksaan |
| test_result_date | DATE | NULLABLE | Tanggal hasil tes selesai |
| result_value | TEXT | NULLABLE | Nilai/hasil pemeriksaan |
| result_unit | VARCHAR(50) | NULLABLE | Unit hasil (mg/dL, mmol/L, dll) |
| reference_range | VARCHAR(200) | NULLABLE | Nilai normal referensi |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'ordered', INDEX | Status (ordered, collected, started, completed, cancelled) |
| notes | TEXT | NULLABLE | Catatan khusus/interpretasi |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

**Indexes:**
- Primary Key: id
- Foreign Keys: medical_record_id, test_type_id
- Regular Index: sample_collection_date, test_start_date, status, deleted_at

**Relationships:**
- Belongs To Medical Record (many-to-one)
- Belongs To Test Type (many-to-one)

**Notes:**
- Status flow: ordered -> collected -> started -> completed
- Result value, unit, dan reference range diisi saat tes selesai
- Tracking full lifecycle dari test order hingga hasil

---

## Lab Test Status Flow

```
ordered → sample_collected → in_progress → completed
ordered → cancelled
sample_collected → cancelled
in_progress → cancelled
```

**Status Values:**

- `ordered` — Pemeriksaan telah dipesan oleh dokter
- `sample_collected` — Sampel pasien telah diambil (`sample_collection_date` diisi otomatis)
- `in_progress` — Pemeriksaan sedang dikerjakan di laboratorium
- `completed` — Pemeriksaan selesai, hasil telah diisi (`result_date` diisi otomatis)
- `cancelled` — Pemeriksaan dibatalkan

> **Notes:**
>
> - `result_value` dapat berisi angka tunggal (mis. `"12.5"`) atau rentang (mis. `"reactive"`)
> - `result_unit` diisi jika hasil berupa nilai numerik (mis. `"g/dL"`, `"mg/dL"`, `"mmol/L"`)
> - `reference_range` menunjukkan nilai normal untuk perbandingan hasil
