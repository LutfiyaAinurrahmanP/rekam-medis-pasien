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
| GET /medical-history/allergies/patient/:pid | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| GET /medical-history/allergies/:id          | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| POST /medical-history/allergies             | ❌            | ✅     | ❌           | ✅    | ✅          |
| PUT /medical-history/allergies/:id          | ❌            | ✅     | ❌           | ✅    | ✅          |
| DELETE /medical-history/allergies/:id       | ❌            | ❌     | ❌           | ✅    | ✅          |

### Medical Conditions

| Endpoint                                     | Patient (Own) | Doctor | Receptionist | Admin | Super Admin |
| -------------------------------------------- | ------------- | ------ | ------------ | ----- | ----------- |
| GET /medical-history/conditions/patient/:pid | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| GET /medical-history/conditions/:id          | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| POST /medical-history/conditions             | ❌            | ✅     | ❌           | ✅    | ✅          |
| PUT /medical-history/conditions/:id          | ❌            | ✅     | ❌           | ✅    | ✅          |
| DELETE /medical-history/conditions/:id       | ❌            | ❌     | ❌           | ✅    | ✅          |

### Surgical History

| Endpoint                                    | Patient (Own) | Doctor | Receptionist | Admin | Super Admin |
| ------------------------------------------- | ------------- | ------ | ------------ | ----- | ----------- |
| GET /medical-history/surgeries/patient/:pid | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| GET /medical-history/surgeries/:id          | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| POST /medical-history/surgeries             | ❌            | ✅     | ❌           | ✅    | ✅          |
| PUT /medical-history/surgeries/:id          | ❌            | ✅     | ❌           | ✅    | ✅          |
| DELETE /medical-history/surgeries/:id       | ❌            | ❌     | ❌           | ✅    | ✅          |

### Family History

| Endpoint                                         | Patient (Own) | Doctor | Receptionist | Admin | Super Admin |
| ------------------------------------------------ | ------------- | ------ | ------------ | ----- | ----------- |
| GET /medical-history/family-history/patient/:pid | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| GET /medical-history/family-history/:id          | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |
| POST /medical-history/family-history             | ❌            | ✅     | ✅           | ✅    | ✅          |
| PUT /medical-history/family-history/:id          | ❌            | ✅     | ✅           | ✅    | ✅          |
| DELETE /medical-history/family-history/:id       | ❌            | ❌     | ❌           | ✅    | ✅          |

### Summary

| Endpoint                                         | Patient (Own) | Doctor | Receptionist | Admin | Super Admin |
| ------------------------------------------------ | ------------- | ------ | ------------ | ----- | ----------- |
| GET /medical-history/patient/:patient_id/summary | ✅ (Own)      | ✅     | ✅           | ✅    | ✅          |

---

## Endpoints Summary

| Method                 | Endpoint                                          | Description                      | Role Required                            |
| ---------------------- | ------------------------------------------------- | -------------------------------- | ---------------------------------------- |
| GET                    | `/medical-history/patient/:patient_id/summary`    | Get full medical history summary | All Authenticated (own for patient)      |
| **Allergies**          |
| GET                    | `/medical-history/allergies/patient/:patient_id`  | List patient allergies           | All Authenticated (own for patient)      |
| GET                    | `/medical-history/allergies/:id`                  | Get allergy detail by ID         | All Authenticated (own for patient)      |
| POST                   | `/medical-history/allergies`                      | Add allergy                      | Doctor, Admin, Super Admin               |
| PUT                    | `/medical-history/allergies/:id`                  | Update allergy                   | Doctor, Admin, Super Admin               |
| DELETE                 | `/medical-history/allergies/:id`                  | Delete allergy                   | Admin, Super Admin                       |
| **Medical Conditions** |
| GET                    | `/medical-history/conditions/patient/:patient_id` | List patient medical conditions  | All Authenticated (own for patient)      |
| GET                    | `/medical-history/conditions/:id`                 | Get condition detail by ID       | All Authenticated (own for patient)      |
| POST                   | `/medical-history/conditions`                     | Add medical condition            | Doctor, Admin, Super Admin               |
| PUT                    | `/medical-history/conditions/:id`                 | Update medical condition         | Doctor, Admin, Super Admin               |
| DELETE                 | `/medical-history/conditions/:id`                 | Delete medical condition         | Admin, Super Admin                       |
| **Surgical History**   |
| GET                    | `/medical-history/surgeries/patient/:patient_id`  | List patient surgical history    | All Authenticated (own for patient)      |
| GET                    | `/medical-history/surgeries/:id`                  | Get surgery record by ID         | All Authenticated (own for patient)      |
| POST                   | `/medical-history/surgeries`                      | Add surgery record               | Doctor, Admin, Super Admin               |
| PUT                    | `/medical-history/surgeries/:id`                  | Update surgery record            | Doctor, Admin, Super Admin               |
| DELETE                 | `/medical-history/surgeries/:id`                  | Delete surgery record            | Admin, Super Admin                       |
| **Family History**     |
| GET                    | `/medical-history/family-history/patient/:pid`    | List patient family history      | All Authenticated (own for patient)      |
| GET                    | `/medical-history/family-history/:id`             | Get family history entry by ID   | All Authenticated (own for patient)      |
| POST                   | `/medical-history/family-history`                 | Add family history entry         | Doctor, Receptionist, Admin, Super Admin |
| PUT                    | `/medical-history/family-history/:id`             | Update family history entry      | Doctor, Receptionist, Admin, Super Admin |
| DELETE                 | `/medical-history/family-history/:id`             | Delete family history entry      | Admin, Super Admin                       |

---

## Summary Endpoint

### Get Full Medical History Summary

**Endpoint:** `GET /api/v1/medical-history/patient/:patient_id/summary`

**Description:** Mendapatkan ringkasan lengkap riwayat medis pasien dalam satu respons, mencakup alergi, kondisi kronis, riwayat operasi, dan riwayat penyakit keluarga.

**Authentication:** Required (All Authenticated — Patient hanya bisa melihat miliknya sendiri)

**Path Parameters:**

| Parameter  | Type    | Required | Description |
| ---------- | ------- | -------- | ----------- |
| patient_id | integer | Yes      | ID pasien   |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Medical history summary retrieved successfully",
  "data": {
    "patient_id": 10,
    "patient_name": "Budi Santoso",
    "allergies": [
      {
        "id": 1,
        "allergen": "Penisilin",
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
        "hospital_name": "RS Dr. Soetomo",
        "surgeon_name": "dr. Bima Sutaryo, SpB",
        "notes": "Tanpa komplikasi"
      }
    ],
    "family_history": [
      {
        "id": 1,
        "relation": "father",
        "condition_name": "Hipertensi",
        "age_at_diagnosis": 52,
        "notes": "Meninggal karena stroke di usia 68 tahun"
      }
    ]
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/medical-history/patient/10/summary" \
  -H "Authorization: Bearer <token>"
```

---

## Allergy Endpoints

### 1. List Patient Allergies

**Endpoint:** `GET /api/v1/medical-history/allergies/patient/:patient_id`

**Description:** Mendapatkan daftar semua alergi yang dimiliki pasien.

**Authentication:** Required (All Authenticated — Patient hanya bisa melihat miliknya sendiri)

**Path Parameters:**

| Parameter  | Type    | Required | Description |
| ---------- | ------- | -------- | ----------- |
| patient_id | integer | Yes      | ID pasien   |

**Query Parameters:**

| Parameter     | Type   | Required | Default | Description                                           |
| ------------- | ------ | -------- | ------- | ----------------------------------------------------- |
| allergen_type | string | No       | -       | Filter tipe: `drug`, `food`, `environmental`, `other` |
| severity      | string | No       | -       | Filter keparahan: `mild`, `moderate`, `severe`        |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Allergies retrieved successfully",
  "data": [
    {
      "id": 1,
      "patient_id": 10,
      "allergen": "Penisilin",
      "allergen_type": "drug",
      "reaction": "Ruam kulit dan sesak napas",
      "severity": "severe",
      "notes": "Hindari semua golongan Beta-laktam",
      "created_at": "2024-01-15T09:00:00Z",
      "updated_at": "2024-01-15T09:00:00Z"
    },
    {
      "id": 2,
      "patient_id": 10,
      "allergen": "Udang",
      "allergen_type": "food",
      "reaction": "Gatal-gatal dan biduran",
      "severity": "moderate",
      "notes": null,
      "created_at": "2024-01-15T09:05:00Z",
      "updated_at": "2024-01-15T09:05:00Z"
    }
  ]
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/medical-history/allergies/patient/10?allergen_type=drug" \
  -H "Authorization: Bearer <token>"
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
    "allergen": "Penisilin",
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
| allergen      | string  | Yes      | Nama zat/bahan yang menjadi alergen (max 255 karakter)         |
| allergen_type | string  | Yes      | Tipe alergen: `drug`, `food`, `environmental`, `other`         |
| reaction      | string  | Yes      | Deskripsi reaksi alergi yang timbul                            |
| severity      | string  | Yes      | Tingkat keparahan: `mild`, `moderate`, `severe`                |
| notes         | string  | No       | Catatan tambahan, mis. obat pengganti atau tindakan pencegahan |

**Example Request Body:**

```json
{
  "patient_id": 10,
  "allergen": "Penisilin",
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
    "allergen": "Penisilin",
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
    "allergen": "Penisilin",
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
    "allergen": "Penisilin",
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

### 1. List Patient Medical Conditions

**Endpoint:** `GET /api/v1/medical-history/conditions/patient/:patient_id`

**Description:** Mendapatkan daftar kondisi medis/penyakit kronis yang dimiliki pasien.

**Authentication:** Required (All Authenticated — Patient hanya bisa melihat miliknya sendiri)

**Path Parameters:**

| Parameter  | Type    | Required | Description |
| ---------- | ------- | -------- | ----------- |
| patient_id | integer | Yes      | ID pasien   |

**Query Parameters:**

| Parameter | Type   | Required | Default | Description                                     |
| --------- | ------ | -------- | ------- | ----------------------------------------------- |
| status    | string | No       | -       | Filter status: `ongoing`, `resolved`, `managed` |

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
    },
    {
      "id": 2,
      "patient_id": 10,
      "condition_name": "Hipertensi",
      "icd_code": "I10",
      "diagnosed_date": "2020-06-01",
      "status": "managed",
      "notes": "Terkontrol dengan Amlodipine 5mg",
      "created_at": "2024-01-15T09:10:00Z",
      "updated_at": "2024-01-15T09:10:00Z"
    }
  ]
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/medical-history/conditions/patient/10?status=ongoing" \
  -H "Authorization: Bearer <token>"
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

### 1. List Patient Surgical History

**Endpoint:** `GET /api/v1/medical-history/surgeries/patient/:patient_id`

**Description:** Mendapatkan daftar riwayat operasi/tindakan medis besar yang pernah dijalani pasien.

**Authentication:** Required (All Authenticated — Patient hanya bisa melihat miliknya sendiri)

**Path Parameters:**

| Parameter  | Type    | Required | Description |
| ---------- | ------- | -------- | ----------- |
| patient_id | integer | Yes      | ID pasien   |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Surgical history retrieved successfully",
  "data": [
    {
      "id": 1,
      "patient_id": 10,
      "procedure_name": "Appendektomi",
      "surgery_date": "2015-07-20",
      "hospital_name": "RS Dr. Soetomo Surabaya",
      "surgeon_name": "dr. Bima Sutaryo, SpB",
      "complications": null,
      "notes": "Operasi berjalan lancar, tanpa komplikasi",
      "created_at": "2024-01-15T09:00:00Z",
      "updated_at": "2024-01-15T09:00:00Z"
    }
  ]
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/medical-history/surgeries/patient/10" \
  -H "Authorization: Bearer <token>"
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
| hospital_name  | string  | No       | Nama fasilitas kesehatan tempat operasi     |
| surgeon_name   | string  | No       | Nama dokter/operator yang melakukan operasi |
| complications  | string  | No       | Komplikasi yang dialami (jika ada)          |
| notes          | string  | No       | Catatan tambahan                            |

**Example Request Body:**

```json
{
  "patient_id": 10,
  "procedure_name": "Appendektomi Laparoskopi",
  "surgery_date": "2015-07-20",
  "hospital_name": "RS Dr. Soetomo Surabaya",
  "surgeon_name": "dr. Bima Sutaryo, SpB",
  "complications": null,
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
    "hospital_name": "RS Dr. Soetomo Surabaya",
    "surgeon_name": "dr. Bima Sutaryo, SpB",
    "complications": null,
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
    "hospital_name": "RS Dr. Soetomo Surabaya",
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
  "complications": "Infeksi luka pasca operasi, dirawat selama 3 hari tambahan"
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Surgical history updated successfully",
  "data": {
    "id": 1,
    "complications": "Infeksi luka pasca operasi, dirawat selama 3 hari tambahan",
    "updated_at": "2024-01-20T10:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT "http://localhost:8080/api/v1/medical-history/surgeries/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"complications": "Infeksi luka pasca operasi"}'
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

### 1. List Patient Family History

**Endpoint:** `GET /api/v1/medical-history/family-history/patient/:patient_id`

**Description:** Mendapatkan daftar riwayat penyakit keluarga pasien yang berpotensi bersifat herediter.

**Authentication:** Required (All Authenticated — Patient hanya bisa melihat miliknya sendiri)

**Path Parameters:**

| Parameter  | Type    | Required | Description |
| ---------- | ------- | -------- | ----------- |
| patient_id | integer | Yes      | ID pasien   |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Family history retrieved successfully",
  "data": [
    {
      "id": 1,
      "patient_id": 10,
      "relation": "father",
      "condition_name": "Hipertensi",
      "age_at_diagnosis": 52,
      "notes": "Meninggal karena stroke di usia 68 tahun",
      "created_at": "2024-01-15T09:00:00Z",
      "updated_at": "2024-01-15T09:00:00Z"
    },
    {
      "id": 2,
      "patient_id": 10,
      "relation": "mother",
      "condition_name": "Diabetes Mellitus",
      "age_at_diagnosis": 48,
      "notes": "Masih hidup, terkontrol dengan insulin",
      "created_at": "2024-01-15T09:05:00Z",
      "updated_at": "2024-01-15T09:05:00Z"
    }
  ]
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/medical-history/family-history/patient/10" \
  -H "Authorization: Bearer <token>"
```

---

### 2. Add Family History Entry

**Endpoint:** `POST /api/v1/medical-history/family-history`

**Description:** Menambahkan riwayat penyakit pada anggota keluarga pasien.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Body:**

| Field            | Type    | Required | Description                                                                       |
| ---------------- | ------- | -------- | --------------------------------------------------------------------------------- |
| patient_id       | integer | Yes      | ID pasien                                                                         |
| relation         | string  | Yes      | Hubungan keluarga: `father`, `mother`, `sibling`, `grandparent`, `child`, `other` |
| condition_name   | string  | Yes      | Nama penyakit anggota keluarga (max 255 karakter)                                 |
| age_at_diagnosis | integer | No       | Usia anggota keluarga saat didiagnosis (tahun)                                    |
| notes            | string  | No       | Catatan tambahan (kondisi saat ini, meninggal, dll.)                              |

**Example Request Body:**

```json
{
  "patient_id": 10,
  "relation": "father",
  "condition_name": "Hipertensi",
  "age_at_diagnosis": 52,
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
    "relation": "father",
    "condition_name": "Hipertensi",
    "age_at_diagnosis": 52,
    "notes": "Meninggal karena komplikasi stroke di usia 68 tahun",
    "created_at": "2024-01-15T09:00:00Z",
    "updated_at": "2024-01-15T09:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/medical-history/family-history" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 10,
    "relation": "father",
    "condition_name": "Hipertensi",
    "age_at_diagnosis": 52,
    "notes": "Meninggal karena komplikasi stroke di usia 68 tahun"
  }'
```

---

### 3. Update Family History Entry

**Endpoint:** `PUT /api/v1/medical-history/family-history/:id`

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
curl -X PUT "http://localhost:8080/api/v1/medical-history/family-history/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"notes": "Meninggal karena stroke usia 68 tahun."}'
```

---

### 4. Delete Family History Entry

**Endpoint:** `DELETE /api/v1/medical-history/family-history/:id`

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
    "allergen": "allergen is required",
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

### Family Relation Values

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
