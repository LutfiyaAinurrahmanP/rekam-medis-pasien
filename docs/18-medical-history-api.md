# Medical History API Documentation

## Overview

API untuk manajemen riwayat medis pasien (medical history) dalam sistem rekam medis. Medical History mencakup data alergi, kondisi medis kronis, riwayat operasi/tindakan sebelumnya, dan riwayat penyakit keluarga. Data ini bersifat permanen dan terus diperbarui sepanjang riwayat perawatan pasien.

**Base URL:** `/api/v1/medical-history`

---

## Table of Contents

- [Authentication](#authentication)
- [Authorization](#authorization)
- [Endpoints Summary](#endpoints-summary)
- [Allergy Endpoints](#allergy-endpoints)
- [Medical Condition Endpoints](#medical-condition-endpoints)
- [Surgical History Endpoints](#surgical-history-endpoints)
- [Family History Endpoints](#family-history-endpoints)
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

### Allergies

| Endpoint                                    | Patient (Own) | Doctor | Receptionist | Admin | Super Admin |
| ------------------------------------------- | ------------- | ------ | ------------ | ----- | ----------- |
| GET /medical-history/allergies              | ❌            | ✅     | ✅           | ✅    | ✅          |
| GET /medical-history/allergies/:id          | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| POST /medical-history/allergies             | ❌            | ✅     | ❌           | ✅    | ✅          |
| PUT /medical-history/allergies/:id          | ❌            | ✅     | ❌           | ✅    | ✅          |
| DELETE /medical-history/allergies/:id       | ❌            | ❌     | ❌           | ✅    | ✅          |

### Medical Conditions

| Endpoint                                     | Patient (Own) | Doctor | Receptionist | Admin | Super Admin |
| -------------------------------------------- | ------------- | ------ | ------------ | ----- | ----------- |
| GET /medical-history/conditions              | ❌            | ✅     | ✅           | ✅    | ✅          |
| GET /medical-history/conditions/:id          | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| POST /medical-history/conditions             | ❌            | ✅     | ❌           | ✅    | ✅          |
| PUT /medical-history/conditions/:id          | ❌            | ✅     | ❌           | ✅    | ✅          |
| DELETE /medical-history/conditions/:id       | ❌            | ❌     | ❌           | ✅    | ✅          |

### Surgical History

| Endpoint                                    | Patient (Own) | Doctor | Receptionist | Admin | Super Admin |
| ------------------------------------------- | ------------- | ------ | ------------ | ----- | ----------- |
| GET /medical-history/surgeries              | ❌            | ✅     | ✅           | ✅    | ✅          |
| GET /medical-history/surgeries/:id          | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| POST /medical-history/surgeries             | ❌            | ✅     | ❌           | ✅    | ✅          |
| PUT /medical-history/surgeries/:id          | ❌            | ✅     | ❌           | ✅    | ✅          |
| DELETE /medical-history/surgeries/:id       | ❌            | ❌     | ❌           | ✅    | ✅          |

### Family History

| Endpoint                                         | Patient (Own) | Doctor | Receptionist | Admin | Super Admin |
| ------------------------------------------------ | ------------- | ------ | ------------ | ----- | ----------- |
| GET /medical-history/family-histories              | ❌            | ✅     | ✅           | ✅    | ✅          |
| GET /medical-history/family-histories/:id          | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| POST /medical-history/family-histories             | ❌            | ✅     | ✅           | ✅    | ✅          |
| PUT /medical-history/family-histories/:id          | ❌            | ✅     | ✅           | ✅    | ✅          |
| DELETE /medical-history/family-histories/:id       | ❌            | ❌     | ❌           | ✅    | ✅          |

### General Medical History

| Endpoint                                         | Patient (Own) | Doctor | Receptionist | Admin | Super Admin |
| ------------------------------------------------ | ------------- | ------ | ------------ | ----- | ----------- |
| GET /medical-history                             | ❌            | ✅     | ✅           | ✅    | ✅          |
| GET /medical-history/:id                         | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| GET /medical-history/patient/:pid                | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |

---

## Endpoints Summary

| Method                 | Endpoint                                          | Description                      | Role Required                            |
| ---------------------- | ------------------------------------------------- | -------------------------------- | ---------------------------------------- |
| **Medical History**    |                                                   |                                  |                                          |
| GET                    | `/medical-history`                                | List all medical histories       | Doctor, Receptionist, Admin, Super Admin |
| GET                    | `/medical-history/:id`                            | Get medical history by ID        | All Authenticated (own for patient)      |
| GET                    | `/medical-history/patient/:patient_id`            | Get medical history by patient   | All Authenticated (own for patient)      |
| **Allergies**          |                                                   |                                  |                                          |
| GET                    | `/medical-history/allergies`                      | List all allergies               | Doctor, Receptionist, Admin, Super Admin |
| GET                    | `/medical-history/allergies/:id`                  | Get allergy detail by ID         | All Authenticated (own for patient)      |
| POST                   | `/medical-history/allergies`                      | Add allergy                      | Doctor, Admin, Super Admin               |
| PUT                    | `/medical-history/allergies/:id`                  | Update allergy                   | Doctor, Admin, Super Admin               |
| DELETE                 | `/medical-history/allergies/:id`                  | Delete allergy                   | Admin, Super Admin                       |
| **Medical Conditions** |                                                   |                                  |                                          |
| GET                    | `/medical-history/conditions`                     | List all medical conditions      | Doctor, Receptionist, Admin, Super Admin |
| GET                    | `/medical-history/conditions/:id`                 | Get condition detail by ID       | All Authenticated (own for patient)      |
| POST                   | `/medical-history/conditions`                     | Add medical condition            | Doctor, Admin, Super Admin               |
| PUT                    | `/medical-history/conditions/:id`                 | Update medical condition         | Doctor, Admin, Super Admin               |
| DELETE                 | `/medical-history/conditions/:id`                 | Delete medical condition         | Admin, Super Admin                       |
| **Surgical History**   |                                                   |                                  |                                          |
| GET                    | `/medical-history/surgeries`                      | List all surgical histories      | Doctor, Receptionist, Admin, Super Admin |
| GET                    | `/medical-history/surgeries/:id`                  | Get surgery record by ID         | All Authenticated (own for patient)      |
| POST                   | `/medical-history/surgeries`                      | Add surgery record               | Doctor, Admin, Super Admin               |
| PUT                    | `/medical-history/surgeries/:id`                  | Update surgery record            | Doctor, Admin, Super Admin               |
| DELETE                 | `/medical-history/surgeries/:id`                  | Delete surgery record            | Admin, Super Admin                       |
| **Family History**     |                                                   |                                  |                                          |
| GET                    | `/medical-history/family-histories`                 | List all family histories        | Doctor, Receptionist, Admin, Super Admin |
| GET                    | `/medical-history/family-histories/:id`             | Get family history entry by ID   | All Authenticated (own for patient)      |
| POST                   | `/medical-history/family-histories`                 | Add family history entry         | Doctor, Receptionist, Admin, Super Admin |
| PUT                    | `/medical-history/family-histories/:id`             | Update family history entry      | Doctor, Receptionist, Admin, Super Admin |
| DELETE                 | `/medical-history/family-histories/:id`             | Delete family history entry      | Admin, Super Admin                       |

---

## General Medical History Endpoints

### 1. List Medical Histories

**Endpoint:** `GET /api/v1/medical-history`

**Description:** Mendapatkan daftar seluruh riwayat medis (semua pasien, untuk keperluan admin/dokter).

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Medical histories retrieved successfully",
  "data": [
    {
      "id": 1,
      "patient_id": 10,
      "patient_name": "Budi Santoso",
      "allergies_count": 1,
      "conditions_count": 1,
      "surgeries_count": 1,
      "family_history_count": 1,
      "last_updated": "2024-01-15T09:00:00Z"
    }
  ]
}
```

### 2. Get Medical History by ID

**Endpoint:** `GET /api/v1/medical-history/:id`

**Description:** Mendapatkan detail riwayat medis berdasarkan ID agregat.

**Authentication:** Required (All Authenticated — Patient hanya bisa melihat miliknya sendiri)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID Medical History |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Medical history retrieved successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "allergies": [],
    "medical_conditions": [],
    "surgical_history": [],
    "family_history": []
  }
}
```

### 3. Get Medical History by Patient

**Endpoint:** `GET /api/v1/medical-history/patient/:patient_id`

**Description:** Mendapatkan ringkasan riwayat medis untuk satu pasien (menggabungkan data alergi, kondisi medis, operasi, dan riwayat keluarga).

**Authentication:** Required (All Authenticated — Patient hanya bisa melihat miliknya sendiri)

**Path Parameters:**

| Parameter  | Type    | Required | Description |
| ---------- | ------- | -------- | ----------- |
| patient_id | integer | Yes      | ID pasien   |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Patient medical history retrieved successfully",
  "data": {
    "patient_id": 10,
    "patient_name": "Budi Santoso",
    "allergies": [
      {
        "id": 1,
        "allergen_name": "Penisilin",
        "allergen_type": "drug",
        "reaction": "Ruam kulit dan sesak napas",
        "severity": "severe",
        "notes": "Hindari semua golongan Beta-laktam"
      }
    ],
    "medical_conditions": [
      {
        "id": 1,
        "condition_name": "Diabetes Mellitus Tipe 2",
        "icd_code": "E11",
        "diagnosed_date": "2018-03-15",
        "status": "ongoing",
        "notes": "Terkontrol dengan Metformin 500mg"
      }
    ],
    "surgical_history": [
      {
        "id": 1,
        "procedure_name": "Appendektomi",
        "surgery_date": "2015-07-20",
        "hospital": "RS Dr. Soetomo",
        "surgeon_name": "dr. Bima Sutaryo, SpB",
        "notes": "Tanpa komplikasi"
      }
    ],
    "family_history": [
      {
        "id": 1,
        "family_member": "father",
        "condition_name": "Hipertensi",
        "notes": "Meninggal karena stroke di usia 68 tahun"
      }
    ]
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/medical-history/patient/10" \
  -H "Authorization: Bearer <token>"
```

---

## Allergy Endpoints

### 1. List Allergies

**Endpoint:** `GET /api/v1/medical-history/allergies`

**Description:** Mendapatkan daftar semua alergi yang tercatat di sistem.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| page | integer | No | 1 | Halaman yang ingin diambil |
| limit | integer | No | 10 | Jumlah data per halaman |
| patient_id | integer | No | - | Filter berdasarkan ID pasien |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Allergies retrieved successfully",
  "data": [
    {
      "id": 1,
      "patient_id": 10,
      "allergen_name": "Penisilin",
      "allergen_type": "drug",
      "reaction": "Ruam kulit dan sesak napas",
      "severity": "severe",
      "notes": "Hindari semua golongan Beta-laktam",
      "created_at": "2024-01-15T09:00:00Z",
      "updated_at": "2024-01-15T09:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

---

### 2. Get Allergy by ID

**Endpoint:** `GET /api/v1/medical-history/allergies/:id`

**Description:** Mendapatkan detail alergi berdasarkan ID.

**Authentication:** Required (All Authenticated — Patient hanya bisa melihat miliknya sendiri)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID alergi   |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Allergy retrieved successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "allergen_name": "Penisilin",
    "allergen_type": "drug",
    "reaction": "Ruam kulit dan sesak napas",
    "severity": "severe",
    "notes": "Hindari semua golongan Beta-laktam",
    "created_at": "2024-01-15T09:00:00Z",
    "updated_at": "2024-01-15T09:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/medical-history/allergies/1" \
  -H "Authorization: Bearer <token>"
```

---

### 3. Add Allergy

**Endpoint:** `POST /api/v1/medical-history/allergies`

**Description:** Menambahkan data alergi baru untuk pasien.

**Authentication:** Required (Doctor, Admin, Super Admin)

**Request Headers:**

| Header        | Type   | Required | Description        |
| ------------- | ------ | -------- | ------------------ |
| Authorization | string | Yes      | `Bearer <token>`   |
| Content-Type  | string | Yes      | `application/json` |

**Request Body:**

| Field         | Type    | Required | Description                                                    |
| ------------- | ------- | -------- | -------------------------------------------------------------- |
| patient_id    | integer | Yes      | ID pasien                                                      |
| allergen_name | string  | Yes      | Nama zat/bahan yang menjadi alergen (max 255 karakter)         |
| allergen_type | string  | Yes      | Tipe alergen: `drug`, `food`, `environmental`, `other`         |
| reaction      | string  | Yes      | Deskripsi reaksi alergi yang timbul                            |
| severity      | string  | Yes      | Tingkat keparahan: `mild`, `moderate`, `severe`                |
| notes         | string  | No       | Catatan tambahan, mis. obat pengganti atau tindakan pencegahan |

**Example Request Body:**

```json
{
  "patient_id": 10,
  "allergen_name": "Penisilin",
  "allergen_type": "drug",
  "reaction": "Ruam kulit, biduran, dan sesak napas dalam 30 menit",
  "severity": "severe",
  "notes": "Hindari semua golongan Beta-laktam. Gunakan Eritromisin sebagai alternatif."
}
```

**Response Success (201):**

```json
{
  "status": "success",
  "message": "Allergy added successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "allergen_name": "Penisilin",
    "allergen_type": "drug",
    "reaction": "Ruam kulit, biduran, dan sesak napas dalam 30 menit",
    "severity": "severe",
    "notes": "Hindari semua golongan Beta-laktam. Gunakan Eritromisin sebagai alternatif.",
    "created_at": "2024-01-15T09:00:00Z",
    "updated_at": "2024-01-15T09:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/medical-history/allergies" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 10,
    "allergen_name": "Penisilin",
    "allergen_type": "drug",
    "reaction": "Ruam kulit, biduran, dan sesak napas dalam 30 menit",
    "severity": "severe",
    "notes": "Hindari semua golongan Beta-laktam."
  }'
```

---

### 4. Update Allergy

**Endpoint:** `PUT /api/v1/medical-history/allergies/:id`

**Description:** Memperbarui data alergi berdasarkan ID.

**Authentication:** Required (Doctor, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID alergi   |

**Request Body:** (Same fields as Add Allergy, all optional for update)

**Example Request Body:**

```json
{
  "severity": "severe",
  "notes": "Diperbarui: Hindari semua Beta-laktam dan Karbapenem."
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Allergy updated successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "allergen_name": "Penisilin",
    "severity": "severe",
    "notes": "Diperbarui: Hindari semua Beta-laktam dan Karbapenem.",
    "updated_at": "2024-01-20T10:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT "http://localhost:8080/api/v1/medical-history/allergies/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "severity": "severe",
    "notes": "Diperbarui: Hindari semua Beta-laktam dan Karbapenem."
  }'
```

---

### 5. Delete Allergy

**Endpoint:** `DELETE /api/v1/medical-history/allergies/:id`

**Description:** Menghapus data alergi.

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID alergi   |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Allergy deleted successfully"
}
```

**cURL Example:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/medical-history/allergies/1" \
  -H "Authorization: Bearer <token>"
```

---

## Medical Condition Endpoints

### 1. List Medical Conditions

**Endpoint:** `GET /api/v1/medical-history/conditions`

**Description:** Mendapatkan daftar semua kondisi medis yang tercatat di sistem.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| page | integer | No | 1 | Halaman yang ingin diambil |
| limit | integer | No | 10 | Jumlah data per halaman |
| patient_id | integer | No | - | Filter berdasarkan ID pasien |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Medical conditions retrieved successfully",
  "data": [
    {
      "id": 1,
      "patient_id": 10,
      "condition_name": "Diabetes Mellitus Tipe 2",
      "icd_code": "E11",
      "diagnosed_date": "2018-03-15",
      "status": "ongoing",
      "notes": "Terkontrol dengan Metformin 500mg",
      "created_at": "2024-01-15T09:00:00Z",
      "updated_at": "2024-01-15T09:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

---

### 2. Add Medical Condition

**Endpoint:** `POST /api/v1/medical-history/conditions`

**Description:** Menambahkan kondisi medis/penyakit kronis baru untuk pasien.

**Authentication:** Required (Doctor, Admin, Super Admin)

**Request Body:**

| Field          | Type    | Required | Description                                                    |
| -------------- | ------- | -------- | -------------------------------------------------------------- |
| patient_id     | integer | Yes      | ID pasien                                                      |
| condition_name | string  | Yes      | Nama kondisi/penyakit (max 255 karakter)                       |
| icd_code       | string  | No       | Kode ICD-10 (mis. `E11`, `I10`)                                |
| diagnosed_date | string  | No       | Tanggal diagnosis (format: `YYYY-MM-DD`)                       |
| status         | string  | Yes      | Status: `ongoing`, `resolved`, `managed`                       |
| notes          | string  | No       | Catatan tambahan (obat yang digunakan, kondisi terakhir, dll.) |

**Example Request Body:**

```json
{
  "patient_id": 10,
  "condition_name": "Diabetes Mellitus Tipe 2",
  "icd_code": "E11",
  "diagnosed_date": "2018-03-15",
  "status": "ongoing",
  "notes": "Terkontrol dengan Metformin 500mg 2x sehari"
}
```

**Response Success (201):**

```json
{
  "status": "success",
  "message": "Medical condition added successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "condition_name": "Diabetes Mellitus Tipe 2",
    "icd_code": "E11",
    "diagnosed_date": "2018-03-15",
    "status": "ongoing",
    "notes": "Terkontrol dengan Metformin 500mg 2x sehari",
    "created_at": "2024-01-15T09:00:00Z",
    "updated_at": "2024-01-15T09:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/medical-history/conditions" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 10,
    "condition_name": "Diabetes Mellitus Tipe 2",
    "icd_code": "E11",
    "diagnosed_date": "2018-03-15",
    "status": "ongoing",
    "notes": "Terkontrol dengan Metformin 500mg 2x sehari"
  }'
```

---

### 3. Update Medical Condition

**Endpoint:** `PUT /api/v1/medical-history/conditions/:id`

**Description:** Memperbarui data kondisi medis berdasarkan ID.

**Authentication:** Required (Doctor, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description      |
| --------- | ------- | -------- | ---------------- |
| id        | integer | Yes      | ID kondisi medis |

**Example Request Body:**

```json
{
  "status": "managed",
  "notes": "Terkontrol baik dengan Metformin 500mg + Glibenklamid 5mg"
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Medical condition updated successfully",
  "data": {
    "id": 1,
    "condition_name": "Diabetes Mellitus Tipe 2",
    "status": "managed",
    "notes": "Terkontrol baik dengan Metformin 500mg + Glibenklamid 5mg",
    "updated_at": "2024-06-01T08:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT "http://localhost:8080/api/v1/medical-history/conditions/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "managed",
    "notes": "Terkontrol baik dengan Metformin 500mg + Glibenklamid 5mg"
  }'
```

---

### 4. Delete Medical Condition

**Endpoint:** `DELETE /api/v1/medical-history/conditions/:id`

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description      |
| --------- | ------- | -------- | ---------------- |
| id        | integer | Yes      | ID kondisi medis |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Medical condition deleted successfully"
}
```

---

## Surgical History Endpoints

### 1. List Surgical Histories

**Endpoint:** `GET /api/v1/medical-history/surgeries`

**Description:** Mendapatkan daftar semua riwayat operasi yang tercatat di sistem.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| page | integer | No | 1 | Halaman yang ingin diambil |
| limit | integer | No | 10 | Jumlah data per halaman |
| patient_id | integer | No | - | Filter berdasarkan ID pasien |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Surgical histories retrieved successfully",
  "data": [
    {
      "id": 1,
      "patient_id": 10,
      "procedure_name": "Appendektomi",
      "surgery_date": "2015-07-20",
      "hospital": "RS Dr. Soetomo Surabaya",
      "surgeon_name": "dr. Bima Sutaryo, SpB",
      "notes": "Operasi berjalan lancar, tanpa komplikasi",
      "created_at": "2024-01-15T09:00:00Z",
      "updated_at": "2024-01-15T09:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

---

### 2. Add Surgical History

**Endpoint:** `POST /api/v1/medical-history/surgeries`

**Description:** Menambahkan riwayat operasi/tindakan medis besar.

**Authentication:** Required (Doctor, Admin, Super Admin)

**Request Body:**

| Field          | Type    | Required | Description                                 |
| -------------- | ------- | -------- | ------------------------------------------- |
| patient_id     | integer | Yes      | ID pasien                                   |
| procedure_name | string  | Yes      | Nama prosedur/operasi (max 255 karakter)    |
| surgery_date   | string  | No       | Tanggal operasi (format: `YYYY-MM-DD`)      |
| hospital       | string  | No       | Nama fasilitas kesehatan tempat operasi     |
| surgeon_name   | string  | No       | Nama dokter/operator yang melakukan operasi |
| notes          | string  | No       | Catatan tambahan                            |

**Example Request Body:**

```json
{
  "patient_id": 10,
  "procedure_name": "Appendektomi Laparoskopi",
  "surgery_date": "2015-07-20",
  "hospital": "RS Dr. Soetomo Surabaya",
  "surgeon_name": "dr. Bima Sutaryo, SpB",
  "notes": "Operasi laparoskopi berjalan lancar, tanpa komplikasi"
}
```

**Response Success (201):**

```json
{
  "status": "success",
  "message": "Surgical history added successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "procedure_name": "Appendektomi Laparoskopi",
    "surgery_date": "2015-07-20",
    "hospital": "RS Dr. Soetomo Surabaya",
    "surgeon_name": "dr. Bima Sutaryo, SpB",
    "notes": "Operasi laparoskopi berjalan lancar, tanpa komplikasi",
    "created_at": "2024-01-15T09:00:00Z",
    "updated_at": "2024-01-15T09:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/medical-history/surgeries" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 10,
    "procedure_name": "Appendektomi Laparoskopi",
    "surgery_date": "2015-07-20",
    "hospital": "RS Dr. Soetomo Surabaya",
    "surgeon_name": "dr. Bima Sutaryo, SpB",
    "notes": "Operasi berjalan lancar"
  }'
```

---

### 3. Update Surgical History

**Endpoint:** `PUT /api/v1/medical-history/surgeries/:id`

**Authentication:** Required (Doctor, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description        |
| --------- | ------- | -------- | ------------------ |
| id        | integer | Yes      | ID riwayat operasi |

**Example Request Body:**

```json
{
  "notes": "Infeksi luka pasca operasi, dirawat selama 3 hari tambahan"
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Surgical history updated successfully",
  "data": {
    "id": 1,
    "notes": "Infeksi luka pasca operasi, dirawat selama 3 hari tambahan",
    "updated_at": "2024-01-20T10:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT "http://localhost:8080/api/v1/medical-history/surgeries/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"notes": "Infeksi luka pasca operasi"}'
```

---

### 4. Delete Surgical History

**Endpoint:** `DELETE /api/v1/medical-history/surgeries/:id`

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description        |
| --------- | ------- | -------- | ------------------ |
| id        | integer | Yes      | ID riwayat operasi |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Surgical history deleted successfully"
}
```

---

## Family History Endpoints

### 1. List Family Histories

**Endpoint:** `GET /api/v1/medical-history/family-histories`

**Description:** Mendapatkan daftar semua riwayat penyakit keluarga yang tercatat di sistem.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
| --- | --- | --- | --- | --- |
| page | integer | No | 1 | Halaman yang ingin diambil |
| limit | integer | No | 10 | Jumlah data per halaman |
| patient_id | integer | No | - | Filter berdasarkan ID pasien |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Family histories retrieved successfully",
  "data": [
    {
      "id": 1,
      "patient_id": 10,
      "family_member": "father",
      "condition_name": "Hipertensi",
      "notes": "Meninggal karena stroke di usia 68 tahun",
      "created_at": "2024-01-15T09:00:00Z",
      "updated_at": "2024-01-15T09:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

---

### 2. Add Family History Entry

**Endpoint:** `POST /api/v1/medical-history/family-histories`

**Description:** Menambahkan riwayat penyakit pada anggota keluarga pasien.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Body:**

| Field            | Type    | Required | Description                                                                       |
| ---------------- | ------- | -------- | --------------------------------------------------------------------------------- |
| patient_id       | integer | Yes      | ID pasien                                                                         |
| family_member    | string  | Yes      | Hubungan keluarga: `father`, `mother`, `sibling`, `grandparent`, `child`, `other` |
| condition_name   | string  | Yes      | Nama penyakit anggota keluarga (max 255 karakter)                                 |
| notes            | string  | No       | Catatan tambahan (kondisi saat ini, meninggal, dll.)                              |

**Example Request Body:**

```json
{
  "patient_id": 10,
  "family_member": "father",
  "condition_name": "Hipertensi",
  "notes": "Meninggal karena komplikasi stroke di usia 68 tahun"
}
```

**Response Success (201):**

```json
{
  "status": "success",
  "message": "Family history added successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "family_member": "father",
    "condition_name": "Hipertensi",
    "notes": "Meninggal karena komplikasi stroke di usia 68 tahun",
    "created_at": "2024-01-15T09:00:00Z",
    "updated_at": "2024-01-15T09:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/medical-history/family-histories" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 10,
    "family_member": "father",
    "condition_name": "Hipertensi",
    "notes": "Meninggal karena komplikasi stroke di usia 68 tahun"
  }'
```

---

### 3. Update Family History Entry

**Endpoint:** `PUT /api/v1/medical-history/family-histories/:id`

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description         |
| --------- | ------- | -------- | ------------------- |
| id        | integer | Yes      | ID riwayat keluarga |

**Example Request Body:**

```json
{
  "notes": "Meninggal karena stroke usia 68 tahun. Hipertensi tidak terkontrol."
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Family history updated successfully",
  "data": {
    "id": 1,
    "notes": "Meninggal karena stroke usia 68 tahun. Hipertensi tidak terkontrol.",
    "updated_at": "2024-02-01T08:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT "http://localhost:8080/api/v1/medical-history/family-histories/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"notes": "Meninggal karena stroke usia 68 tahun."}'
```

---

### 4. Delete Family History Entry

**Endpoint:** `DELETE /api/v1/medical-history/family-histories/:id`

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description         |
| --------- | ------- | -------- | ------------------- |
| id        | integer | Yes      | ID riwayat keluarga |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Family history entry deleted successfully"
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
    "allergen_name": "allergen is required",
    "allergen_type": "allergen_type must be one of: drug, food, environmental, other",
    "severity": "severity must be one of: mild, moderate, severe"
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
  "message": "Medical history record not found"
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

### Table: allergies

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| patient_id | BIGINT | FOREIGN KEY (patients.id), NOT NULL, INDEX | Reference ke pasien |
| allergen_type | VARCHAR(50) | NOT NULL | Tipe alergen (Drug, Food, Environmental, Latex, dll) |
| allergen_name | VARCHAR(100) | NOT NULL | Nama alergen spesifik |
| reaction | TEXT | NOT NULL | Reaksi alergi |
| severity | VARCHAR(20) | NOT NULL | Tingkat keparahan (Mild, Moderate, Severe) |
| notes | TEXT | NULLABLE | Catatan tambahan |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update |

### Table: medical_conditions

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| patient_id | BIGINT | FOREIGN KEY (patients.id), NOT NULL, INDEX | Reference ke pasien |
| condition_name | VARCHAR(200) | NOT NULL | Nama kondisi medis (e.g., Diabetes, Hypertension) |
| icd_code | VARCHAR(20) | NULLABLE | Kode ICD-10 |
| diagnosed_date | DATE | NULLABLE | Tanggal diagnosis |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'ongoing', INDEX | Status (ongoing, resolved, inactive) |
| notes | TEXT | NULLABLE | Catatan detail |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update |

### Table: surgical_histories

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| patient_id | BIGINT | FOREIGN KEY (patients.id), NOT NULL, INDEX | Reference ke pasien |
| procedure_name | VARCHAR(200) | NOT NULL | Nama prosedur/operasi |
| surgery_date | DATE | NOT NULL | Tanggal operasi |
| surgeon_name | VARCHAR(100) | NULLABLE | Nama dokter bedah |
| hospital | VARCHAR(200) | NULLABLE | Nama rumah sakit |
| notes | TEXT | NULLABLE | Catatan hasil operasi |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update |

### Table: family_histories

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| patient_id | BIGINT | FOREIGN KEY (patients.id), NOT NULL, INDEX | Reference ke pasien |
| family_member | VARCHAR(50) | NOT NULL | Hubungan keluarga (Parent, Sibling, Grand Parent, dll) |
| condition_name | VARCHAR(200) | NOT NULL | Nama penyakit dalam keluarga |
| notes | TEXT | NULLABLE | Catatan detail |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update |

**Indexes:**
- Primary Keys: id (each table)
- Foreign Keys: patient_id (all tables)
- Regular Index: allergen_type, condition_name, surgery_date, family_member

**Relationships:**
- All Belong To Patient (many-to-one)

**Notes:**
- Riwayat medis bersifat permanen dan terus diakumulasi
- Multiple records per patient diperbolehkan
- Important untuk clinical decision making
- ICD-10 codes untuk standardisasi diagnosis

---

## Reference Values

### Allergen Type Values

- `drug` — Alergi obat-obatan
- `food` — Alergi makanan
- `environmental` — Alergi lingkungan (debu, serbuk sari, bulu hewan, dll.)
- `other` — Lainnya

### Allergy Severity Values

- `mild` — Ringan (mis. ruam lokal, bersin)
- `moderate` — Sedang (mis. urtikaria, angioedema ringan)
- `severe` — Berat/Mengancam jiwa (mis. anafilaksis, bronkospasme)

### Medical Condition Status Values

- `ongoing` — Masih aktif/berlangsung
- `managed` — Terkontrol dengan pengobatan
- `resolved` — Sudah sembuh/tidak aktif

### Family Member Values

- `father` — Ayah
- `mother` — Ibu
- `sibling` — Saudara kandung
- `grandparent` — Kakek/nenek
- `child` — Anak
- `other` — Hubungan keluarga lainnya

> **Notes:**
>
> - Data riwayat medis bersifat **kumulatif** — tidak dihapus permanen kecuali oleh Super Admin
> - Informasi alergi sangat kritis karena digunakan untuk memvalidasi penulisan resep agar menghindari obat yang menyebabkan alergi
> - Data riwayat keluarga digunakan untuk menilai faktor risiko penyakit herediter seperti diabetes, hipertensi, jantung koroner, dan kanker
