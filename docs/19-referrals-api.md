# Referrals API Documentation

## Overview

API untuk manajemen surat rujukan pasien (referrals) dalam sistem rekam medis. Referrals mencatat proses pengiriman pasien dari satu dokter/dokter umum ke spesialis atau dari satu fasilitas kesehatan ke fasilitas lain yang lebih lengkap. Fitur ini sangat penting untuk mendukung alur layanan BPJS maupun rujukan mandiri.

**Base URL:** `/api/v1/referrals`

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

| Endpoint                           | Patient (Own) | Doctor | Receptionist | Admin | Super Admin |
| ---------------------------------- | ------------- | ------ | ------------ | ----- | ----------- |
| GET /referrals                     | ❌            | ✅     | ✅           | ✅    | ✅          |
| GET /referrals/deleted             | ❌            | ❌     | ❌           | ✅    | ✅          |
| GET /referrals/my-referrals        | ✅            | ❌     | ❌           | ❌    | ❌          |
| GET /referrals/:id                 | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| GET /referrals/patient/:patient_id | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| GET /referrals/doctor/:doctor_id   | ❌            | ✅     | ✅           | ✅    | ✅          |
| POST /referrals                    | ❌            | ✅     | ❌           | ✅    | ✅          |
| PUT /referrals/:id                 | ❌            | ✅     | ❌           | ✅    | ✅          |
| PATCH /referrals/:id/accept        | ❌            | ✅     | ❌           | ✅    | ✅          |
| PATCH /referrals/:id/reject        | ❌            | ✅     | ❌           | ✅    | ✅          |
| PATCH /referrals/:id/complete      | ❌            | ✅     | ✅           | ✅    | ✅          |
| PATCH /referrals/:id/cancel        | ❌            | ✅     | ✅           | ✅    | ✅          |
| DELETE /referrals/:id              | ❌            | ❌     | ❌           | ✅    | ✅          |
| PATCH /referrals/:id/restore       | ❌            | ❌     | ❌           | ✅    | ✅          |
| DELETE /referrals/:id/hard-delete  | ❌            | ❌     | ❌           | ❌    | ✅          |

---

## Endpoints Summary

| Method | Endpoint                         | Description                     | Role Required                            |
| ------ | -------------------------------- | ------------------------------- | ---------------------------------------- |
| GET    | `/referrals`                     | List all referrals              | Doctor, Receptionist, Admin, Super Admin |
| GET    | `/referrals/deleted`             | List deleted referrals          | Admin, Super Admin                       |
| GET    | `/referrals/my-referrals`        | Get my referrals (patient)      | Patient                                  |
| GET    | `/referrals/:id`                 | Get referral by ID              | All Authenticated (own for patient)      |
| GET    | `/referrals/patient/:patient_id` | Get referrals by patient        | All Authenticated (own for patient)      |
| GET    | `/referrals/doctor/:doctor_id`   | Get referrals by doctor         | Doctor, Receptionist, Admin, Super Admin |
| POST   | `/referrals`                     | Create new referral             | Doctor, Admin, Super Admin               |
| PUT    | `/referrals/:id`                 | Update referral                 | Doctor, Admin, Super Admin               |
| PATCH  | `/referrals/:id/accept`          | Accept referral (by specialist) | Doctor, Admin, Super Admin               |
| PATCH  | `/referrals/:id/reject`          | Reject referral                 | Doctor, Admin, Super Admin               |
| PATCH  | `/referrals/:id/complete`        | Complete referral               | Doctor, Receptionist, Admin, Super Admin |
| PATCH  | `/referrals/:id/cancel`          | Cancel referral                 | Doctor, Receptionist, Admin, Super Admin |
| DELETE | `/referrals/:id`                 | Soft delete referral            | Admin, Super Admin                       |
| PATCH  | `/referrals/:id/restore`         | Restore deleted referral        | Admin, Super Admin                       |
| DELETE | `/referrals/:id/hard-delete`     | Permanently delete referral     | Super Admin                              |

---

## Endpoints Detail

### 1. List All Referrals

**Endpoint:** `GET /api/v1/referrals`

**Description:** Mendapatkan daftar semua surat rujukan dengan pagination, search, dan filter.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

| Header        | Type   | Required | Description      |
| ------------- | ------ | -------- | ---------------- |
| Authorization | string | Yes      | `Bearer <token>` |

**Query Parameters:**

| Parameter           | Type    | Required | Default    | Description                                                                |
| ------------------- | ------- | -------- | ---------- | -------------------------------------------------------------------------- |
| page                | integer | No       | 1          | Halaman saat ini                                                           |
| page_size           | integer | No       | 10         | Jumlah data per halaman (max: 100)                                         |
| sort_by             | string  | No       | created_at | Field sorting: `id`, `referral_date`, `status`, `priority`, `created_at`   |
| sort_dir            | string  | No       | desc       | Arah sorting: `asc` atau `desc`                                            |
| search              | string  | No       | -          | Pencarian berdasarkan nomor rujukan atau alasan rujukan                    |
| status              | string  | No       | -          | Filter status: `pending`, `accepted`, `rejected`, `completed`, `cancelled` |
| priority            | string  | No       | -          | Filter prioritas: `routine`, `urgent`, `emergency`                         |
| referral_type       | string  | No       | -          | Filter tipe: `internal`, `external`                                        |
| referring_doctor_id | integer | No       | -          | Filter berdasarkan ID dokter perujuk                                       |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Referrals retrieved successfully",
  "data": [
    {
      "id": 1,
      "referral_number": "REF-2024-000001",
      "patient_id": 10,
      "medical_record_id": 5,
      "referring_doctor_id": 2,
      "referral_date": "2024-01-15",
      "referral_type": "internal",
      "referred_to_department_id": 4,
      "referred_to_doctor_id": 6,
      "referred_to_facility": null,
      "reason": "Pasien memerlukan konsultasi spesialis jantung untuk evaluasi aritmia",
      "diagnosis": "Aritmia, kemungkinan atrial fibrilasi",
      "priority": "urgent",
      "status": "pending",
      "notes": null,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
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
curl -X GET "http://localhost:8080/api/v1/referrals?page=1&page_size=10&status=pending&priority=urgent" \
  -H "Authorization: Bearer <token>"
```

---

### 2. List Deleted Referrals

**Endpoint:** `GET /api/v1/referrals/deleted`

**Description:** Mendapatkan daftar surat rujukan yang telah dihapus (soft deleted).

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
  "message": "Deleted referrals retrieved successfully",
  "data": [
    {
      "id": 2,
      "referral_number": "REF-2024-000002",
      "patient_id": 7,
      "status": "cancelled",
      "referral_date": "2024-01-10"
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
curl -X GET "http://localhost:8080/api/v1/referrals/deleted" \
  -H "Authorization: Bearer <token>"
```

---

### 3. Get My Referrals (Patient)

**Endpoint:** `GET /api/v1/referrals/my-referrals`

**Description:** Mendapatkan daftar semua surat rujukan milik pasien yang sedang login.

**Authentication:** Required (Patient Only)

**Query Parameters:**

| Parameter | Type   | Required | Default | Description                                                                |
| --------- | ------ | -------- | ------- | -------------------------------------------------------------------------- |
| status    | string | No       | -       | Filter status: `pending`, `accepted`, `rejected`, `completed`, `cancelled` |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "My referrals retrieved successfully",
  "data": [
    {
      "id": 1,
      "referral_number": "REF-2024-000001",
      "referral_date": "2024-01-15",
      "referral_type": "internal",
      "reason": "Pasien memerlukan konsultasi spesialis jantung",
      "priority": "urgent",
      "status": "accepted",
      "referring_doctor": {
        "id": 2,
        "name": "dr. Siti Rahayu",
        "specialization": "Penyakit Dalam"
      },
      "referred_to_department": {
        "id": 4,
        "name": "Kardiologi"
      }
    }
  ]
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/referrals/my-referrals?status=pending" \
  -H "Authorization: Bearer <token>"
```

---

### 4. Get Referral by ID

**Endpoint:** `GET /api/v1/referrals/:id`

**Description:** Mendapatkan detail surat rujukan berdasarkan ID.

**Authentication:** Required (All Authenticated — Patient hanya bisa melihat miliknya sendiri)

**Path Parameters:**

| Parameter | Type    | Required | Description      |
| --------- | ------- | -------- | ---------------- |
| id        | integer | Yes      | ID surat rujukan |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Referral retrieved successfully",
  "data": {
    "id": 1,
    "referral_number": "REF-2024-000001",
    "patient_id": 10,
    "medical_record_id": 5,
    "referring_doctor_id": 2,
    "referral_date": "2024-01-15",
    "referral_type": "internal",
    "referred_to_department_id": 4,
    "referred_to_doctor_id": 6,
    "referred_to_facility": null,
    "reason": "Pasien memerlukan evaluasi lebih lanjut oleh spesialis jantung karena ditemukan irama jantung tidak teratur",
    "diagnosis": "Aritmia, kemungkinan atrial fibrilasi — perlu EKG dan ekokardiografi",
    "priority": "urgent",
    "status": "accepted",
    "accepted_at": "2024-01-15T14:00:00Z",
    "completed_at": null,
    "rejection_reason": null,
    "notes": "Harap segera diperiksa, pasien mengeluh jantung berdebar-debar sejak 3 hari",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T14:00:00Z",
    "patient": {
      "id": 10,
      "name": "Budi Santoso",
      "date_of_birth": "1985-03-20",
      "medical_record_number": "MRN-2024-000010"
    },
    "medical_record": {
      "id": 5,
      "visit_date": "2024-01-15",
      "chief_complaint": "Jantung berdebar-debar"
    },
    "referring_doctor": {
      "id": 2,
      "name": "dr. Siti Rahayu",
      "specialization": "Penyakit Dalam"
    },
    "referred_to_department": {
      "id": 4,
      "name": "Kardiologi"
    },
    "referred_to_doctor": {
      "id": 6,
      "name": "dr. Ahmad Fauzi, SpJP",
      "specialization": "Kardiologi"
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/referrals/1" \
  -H "Authorization: Bearer <token>"
```

---

### 5. Get Referrals by Patient

**Endpoint:** `GET /api/v1/referrals/patient/:patient_id`

**Description:** Mendapatkan riwayat semua surat rujukan berdasarkan ID pasien.

**Authentication:** Required (All Authenticated — Patient hanya bisa melihat miliknya sendiri)

**Path Parameters:**

| Parameter  | Type    | Required | Description |
| ---------- | ------- | -------- | ----------- |
| patient_id | integer | Yes      | ID pasien   |

**Query Parameters:**

| Parameter | Type    | Required | Default | Description                                                                |
| --------- | ------- | -------- | ------- | -------------------------------------------------------------------------- |
| status    | string  | No       | -       | Filter status: `pending`, `accepted`, `rejected`, `completed`, `cancelled` |
| page      | integer | No       | 1       | Halaman saat ini                                                           |
| page_size | integer | No       | 10      | Jumlah data per halaman                                                    |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Patient referrals retrieved successfully",
  "data": [
    {
      "id": 1,
      "referral_number": "REF-2024-000001",
      "referral_date": "2024-01-15",
      "referral_type": "internal",
      "reason": "Evaluasi aritmia",
      "priority": "urgent",
      "status": "completed",
      "referring_doctor": {
        "id": 2,
        "name": "dr. Siti Rahayu"
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
curl -X GET "http://localhost:8080/api/v1/referrals/patient/10?status=completed" \
  -H "Authorization: Bearer <token>"
```

---

### 6. Get Referrals by Doctor

**Endpoint:** `GET /api/v1/referrals/doctor/:doctor_id`

**Description:** Mendapatkan daftar rujukan yang diterbitkan oleh dokter tertentu.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| doctor_id | integer | Yes      | ID dokter   |

**Query Parameters:**

| Parameter | Type    | Required | Default | Description             |
| --------- | ------- | -------- | ------- | ----------------------- |
| status    | string  | No       | -       | Filter status rujukan   |
| page      | integer | No       | 1       | Halaman saat ini        |
| page_size | integer | No       | 10      | Jumlah data per halaman |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Doctor referrals retrieved successfully",
  "data": [
    {
      "id": 1,
      "referral_number": "REF-2024-000001",
      "patient_id": 10,
      "referral_date": "2024-01-15",
      "reason": "Evaluasi aritmia",
      "priority": "urgent",
      "status": "accepted"
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
curl -X GET "http://localhost:8080/api/v1/referrals/doctor/2?status=pending" \
  -H "Authorization: Bearer <token>"
```

---

### 7. Create Referral

**Endpoint:** `POST /api/v1/referrals`

**Description:** Membuat surat rujukan baru untuk pasien. Dokter dapat merujuk pasien ke spesialis lain di dalam fasilitas (internal) atau ke fasilitas kesehatan lain (external).

**Authentication:** Required (Doctor, Admin, Super Admin)

**Request Headers:**

| Header        | Type   | Required | Description        |
| ------------- | ------ | -------- | ------------------ |
| Authorization | string | Yes      | `Bearer <token>`   |
| Content-Type  | string | Yes      | `application/json` |

**Request Body:**

| Field                     | Type    | Required | Description                                                                  |
| ------------------------- | ------- | -------- | ---------------------------------------------------------------------------- |
| referral_number           | string  | No       | Nomor surat rujukan (auto-generate jika tidak diisi, max 50 karakter)        |
| patient_id                | integer | Yes      | ID pasien yang dirujuk                                                       |
| medical_record_id         | integer | Yes      | ID rekam medis yang menjadi dasar rujukan                                    |
| referring_doctor_id       | integer | Yes      | ID dokter yang merujuk                                                       |
| referral_date             | string  | Yes      | Tanggal surat rujukan (format: `YYYY-MM-DD`)                                 |
| referral_type             | string  | Yes      | Tipe rujukan: `internal` (dalam fasilitas) atau `external` (antar fasilitas) |
| referred_to_department_id | integer | No       | ID departemen/poli tujuan (wajib jika `referral_type = internal`)            |
| referred_to_doctor_id     | integer | No       | ID dokter spesialis tujuan (opsional jika internal)                          |
| referred_to_facility      | string  | No       | Nama fasilitas kesehatan tujuan (wajib jika `referral_type = external`)      |
| reason                    | string  | Yes      | Alasan/indikasi rujukan                                                      |
| diagnosis                 | string  | No       | Diagnosis sementara atau kerja yang menjadi dasar rujukan                    |
| priority                  | string  | Yes      | Prioritas: `routine`, `urgent`, `emergency`                                  |
| notes                     | string  | No       | Catatan tambahan untuk dokter penerima rujukan                               |

**Example Request Body (Internal Referral):**

```json
{
  "patient_id": 10,
  "medical_record_id": 5,
  "referring_doctor_id": 2,
  "referral_date": "2024-01-15",
  "referral_type": "internal",
  "referred_to_department_id": 4,
  "referred_to_doctor_id": 6,
  "reason": "Pasien memerlukan evaluasi dan tatalaksana lebih lanjut oleh spesialis jantung",
  "diagnosis": "Aritmia, kemungkinan atrial fibrilasi — perlu EKG 12 lead dan ekokardiografi",
  "priority": "urgent",
  "notes": "Pasien mengeluh jantung berdebar-debar dan mudah lelah sejak 3 hari. Riwayat hipertensi."
}
```

**Example Request Body (External Referral):**

```json
{
  "patient_id": 10,
  "medical_record_id": 7,
  "referring_doctor_id": 3,
  "referral_date": "2024-01-20",
  "referral_type": "external",
  "referred_to_facility": "RSUP Dr. Sardjito Yogyakarta",
  "reason": "Pasien memerlukan tindakan bedah jantung yang tidak tersedia di fasilitas ini",
  "diagnosis": "Penyakit Jantung Koroner — perlu CABG",
  "priority": "urgent",
  "notes": "Sertakan hasil pemeriksaan angiografi yang telah dilakukan"
}
```

**Response Success (201):**

```json
{
  "status": "success",
  "message": "Referral created successfully",
  "data": {
    "id": 1,
    "referral_number": "REF-2024-000001",
    "patient_id": 10,
    "medical_record_id": 5,
    "referring_doctor_id": 2,
    "referral_date": "2024-01-15",
    "referral_type": "internal",
    "referred_to_department_id": 4,
    "referred_to_doctor_id": 6,
    "reason": "Pasien memerlukan evaluasi oleh spesialis jantung",
    "diagnosis": "Aritmia, kemungkinan atrial fibrilasi",
    "priority": "urgent",
    "status": "pending",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/referrals" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 10,
    "medical_record_id": 5,
    "referring_doctor_id": 2,
    "referral_date": "2024-01-15",
    "referral_type": "internal",
    "referred_to_department_id": 4,
    "referred_to_doctor_id": 6,
    "reason": "Pasien memerlukan evaluasi oleh spesialis jantung",
    "diagnosis": "Aritmia, kemungkinan atrial fibrilasi",
    "priority": "urgent"
  }'
```

---

### 8. Update Referral

**Endpoint:** `PUT /api/v1/referrals/:id`

**Description:** Memperbarui data surat rujukan. Hanya dapat diperbarui selama status masih `pending`.

**Authentication:** Required (Doctor, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description      |
| --------- | ------- | -------- | ---------------- |
| id        | integer | Yes      | ID surat rujukan |

**Request Body:** (Same fields as Create, all optional for update)

**Example Request Body:**

```json
{
  "priority": "emergency",
  "notes": "Kondisi pasien memburuk, ditemukan tanda-tanda gagal jantung akut"
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Referral updated successfully",
  "data": {
    "id": 1,
    "referral_number": "REF-2024-000001",
    "priority": "emergency",
    "notes": "Kondisi pasien memburuk, ditemukan tanda-tanda gagal jantung akut",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT "http://localhost:8080/api/v1/referrals/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "priority": "emergency",
    "notes": "Kondisi pasien memburuk, ditemukan tanda-tanda gagal jantung akut"
  }'
```

---

### 9. Accept Referral

**Endpoint:** `PATCH /api/v1/referrals/:id/accept`

**Description:** Dokter spesialis/penerima menerima surat rujukan. Status berubah menjadi `accepted` dan `accepted_at` diisi otomatis.

**Authentication:** Required (Doctor, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description      |
| --------- | ------- | -------- | ---------------- |
| id        | integer | Yes      | ID surat rujukan |

**Request Body:**

| Field | Type   | Required | Description                                      |
| ----- | ------ | -------- | ------------------------------------------------ |
| notes | string | No       | Catatan penerimaan (mis. jadwal yang disediakan) |

**Example Request Body:**

```json
{
  "notes": "Diterima. Pasien dijadwalkan untuk EKG dan konsultasi pada 17 Januari 2024."
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Referral accepted successfully",
  "data": {
    "id": 1,
    "referral_number": "REF-2024-000001",
    "status": "accepted",
    "accepted_at": "2024-01-15T14:00:00Z",
    "notes": "Diterima. Pasien dijadwalkan untuk EKG dan konsultasi pada 17 Januari 2024.",
    "updated_at": "2024-01-15T14:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/referrals/1/accept" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Diterima. Pasien dijadwalkan untuk EKG pada 17 Januari 2024."
  }'
```

---

### 10. Reject Referral

**Endpoint:** `PATCH /api/v1/referrals/:id/reject`

**Description:** Dokter atau admin menolak surat rujukan dan mencatat alasan penolakan.

**Authentication:** Required (Doctor, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description      |
| --------- | ------- | -------- | ---------------- |
| id        | integer | Yes      | ID surat rujukan |

**Request Body:**

| Field            | Type   | Required | Description              |
| ---------------- | ------ | -------- | ------------------------ |
| rejection_reason | string | Yes      | Alasan penolakan rujukan |

**Example Request Body:**

```json
{
  "rejection_reason": "Dokter spesialis kardiologi sedang cuti. Silakan rujuk ke RS lain atau jadwalkan ulang setelah tanggal 25 Januari 2024."
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Referral rejected",
  "data": {
    "id": 1,
    "referral_number": "REF-2024-000001",
    "status": "rejected",
    "rejection_reason": "Dokter spesialis kardiologi sedang cuti.",
    "updated_at": "2024-01-15T15:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/referrals/1/reject" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "rejection_reason": "Dokter spesialis sedang cuti. Silakan jadwalkan ulang setelah 25 Januari."
  }'
```

---

### 11. Complete Referral

**Endpoint:** `PATCH /api/v1/referrals/:id/complete`

**Description:** Menandai rujukan sebagai selesai setelah pasien mendapatkan layanan di tujuan rujukan. `completed_at` diisi otomatis.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description      |
| --------- | ------- | -------- | ---------------- |
| id        | integer | Yes      | ID surat rujukan |

**Request Body:**

| Field | Type   | Required | Description                                       |
| ----- | ------ | -------- | ------------------------------------------------- |
| notes | string | No       | Catatan hasil penanganan di tempat tujuan rujukan |

**Example Request Body:**

```json
{
  "notes": "Pasien telah menjalani EKG dan didiagnosis atrial fibrilasi. Saat ini mendapat terapi antikoagulan dan beta-blocker dari dr. Ahmad Fauzi, SpJP."
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Referral completed successfully",
  "data": {
    "id": 1,
    "referral_number": "REF-2024-000001",
    "status": "completed",
    "completed_at": "2024-01-17T11:00:00Z",
    "notes": "Pasien telah didiagnosis atrial fibrilasi dan mendapat terapi.",
    "updated_at": "2024-01-17T11:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/referrals/1/complete" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Pasien telah didiagnosis atrial fibrilasi dan mendapat terapi antikoagulan."
  }'
```

---

### 12. Cancel Referral

**Endpoint:** `PATCH /api/v1/referrals/:id/cancel`

**Description:** Membatalkan surat rujukan.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description      |
| --------- | ------- | -------- | ---------------- |
| id        | integer | Yes      | ID surat rujukan |

**Request Body:**

| Field  | Type   | Required | Description       |
| ------ | ------ | -------- | ----------------- |
| reason | string | No       | Alasan pembatalan |

**Example Request Body:**

```json
{
  "reason": "Pasien meninggal dunia sebelum sempat menjalani rujukan"
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Referral cancelled successfully",
  "data": {
    "id": 1,
    "status": "cancelled",
    "updated_at": "2024-01-16T09:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/referrals/1/cancel" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"reason": "Pasien memilih berobat ke luar kota"}'
```

---

### 13. Soft Delete Referral

**Endpoint:** `DELETE /api/v1/referrals/:id`

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description      |
| --------- | ------- | -------- | ---------------- |
| id        | integer | Yes      | ID surat rujukan |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Referral deleted successfully"
}
```

**cURL Example:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/referrals/1" \
  -H "Authorization: Bearer <token>"
```

---

### 14. Restore Deleted Referral

**Endpoint:** `PATCH /api/v1/referrals/:id/restore`

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description      |
| --------- | ------- | -------- | ---------------- |
| id        | integer | Yes      | ID surat rujukan |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Referral restored successfully",
  "data": {
    "id": 1,
    "referral_number": "REF-2024-000001",
    "status": "pending"
  }
}
```

---

### 15. Hard Delete Referral

**Endpoint:** `DELETE /api/v1/referrals/:id/hard-delete`

**Authentication:** Required (Super Admin Only)

**Path Parameters:**

| Parameter | Type    | Required | Description      |
| --------- | ------- | -------- | ---------------- |
| id        | integer | Yes      | ID surat rujukan |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Referral permanently deleted successfully"
}
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
    "referring_doctor_id": "referring_doctor_id is required",
    "referral_type": "referral_type must be one of: internal, external",
    "priority": "priority must be one of: routine, urgent, emergency",
    "referred_to_department_id": "referred_to_department_id is required for internal referral"
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
  "message": "Referral not found"
}
```

### 422 Unprocessable Entity

```json
{
  "status": "error",
  "message": "Cannot update a referral that has already been accepted or completed"
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

### Table: referrals

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| patient_id | BIGINT | FOREIGN KEY (patients.id), NOT NULL, INDEX | Reference ke pasien |
| from_doctor_id | BIGINT | FOREIGN KEY (doctors.id), NOT NULL, INDEX | Dokter pengirim rujukan |
| to_doctor_id | BIGINT | FOREIGN KEY (doctors.id), NULLABLE, INDEX | Dokter penerima rujukan |
| medical_record_id | BIGINT | FOREIGN KEY (medical_records.id), NULLABLE | Reference ke medical record |
| referral_date | DATE | NOT NULL | Tanggal pembuat rujukan |
| reason | TEXT | NOT NULL | Alasan rujukan |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'pending', INDEX | Status (pending, accepted, rejected, completed, cancelled) |
| acceptance_date | DATE | NULLABLE | Tanggal dokter menerima rujukan |
| completion_date | DATE | NULLABLE | Tanggal rujukan selesai |
| notes | TEXT | NULLABLE | Catatan tambahan |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

**Indexes:**
- Primary Key: id
- Foreign Keys: patient_id, from_doctor_id, to_doctor_id, medical_record_id
- Regular Index: referral_date, status, deleted_at

**Relationships:**
- Belongs To Patient (many-to-one)
- Belongs To From Doctor (many-to-one)
- Belongs To To Doctor (many-to-one, optional)
- Belongs To Medical Record (many-to-one, optional)

**Notes:**
- Status flow: pending -> accepted -> completed (atau rejected/cancelled)
- to_doctor_id optional jika general referral
- acceptance_date dan completion_date diisi saat dokter accept/complete
- Important untuk patient routing dan specialty care

---

## Referral Status Flow

```
pending → accepted → completed
pending → rejected
pending → cancelled
accepted → completed
accepted → cancelled
```

**Status Values:**

- `pending` — Surat rujukan telah dibuat, menunggu konfirmasi penerima
- `accepted` — Rujukan diterima oleh dokter/fasilitas tujuan (`accepted_at` diisi otomatis)
- `rejected` — Rujukan ditolak disertai alasan penolakan
- `completed` — Pasien telah mendapatkan layanan di tujuan rujukan (`completed_at` diisi otomatis)
- `cancelled` — Rujukan dibatalkan

**Referral Type Values:**

- `internal` — Rujukan dalam fasilitas yang sama (mis. dari poli umum ke poli spesialis)
- `external` — Rujukan ke fasilitas kesehatan lain (mis. dari puskesmas ke rumah sakit, atau antar RS)

**Priority Values:**

- `routine` — Kondisi stabil, jadwal normal
- `urgent` — Perlu penanganan lebih cepat, dalam 1–3 hari
- `emergency` — Kondisi mengancam jiwa, segera dirujuk pada hari yang sama

> **Notes:**
>
> - Untuk rujukan `internal`, wajib mengisi `referred_to_department_id`
> - Untuk rujukan `external`, wajib mengisi `referred_to_facility`
> - `rejection_reason` wajib diisi ketika status diubah ke `rejected`
> - Nomor rujukan (`referral_number`) digenerate otomatis jika tidak diisi secara manual
> - Pada konteks BPJS, dokter umum di FKTP (Fasilitas Kesehatan Tingkat Pertama) hanya dapat merujuk ke FKRTL (Rumah Sakit) sesuai aturan yang berlaku
