# Prescriptions API Documentation

## Overview

API untuk manajemen data resep dokter (prescriptions) dan item resep (prescription items) dalam sistem rekam medis. Prescriptions mencatat resep yang ditulis dokter berdasarkan rekam medis pasien, beserta daftar obat yang diresepkan, dosis, frekuensi, dan instruksi penggunaan.

**Base URL:** `/api/v1/prescriptions`

---

## Table of Contents

- [Authentication](#authentication)
- [Authorization](#authorization)
- [Endpoints Summary](#endpoints-summary)
- [Prescription Endpoints](#prescription-endpoints)
- [Prescription Item Endpoints](#prescription-item-endpoints)
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

### Prescriptions

| Endpoint                                     | Patient | Doctor | Receptionist | Admin | Super Admin |
| -------------------------------------------- | ------- | ------ | ------------ | ----- | ----------- |
| GET /prescriptions                           | ❌      | ✅     | ✅           | ✅    | ✅          |
| GET /prescriptions/deleted                   | ❌      | ❌     | ❌           | ✅    | ✅          |
| GET /prescriptions/:id                       | ✅      | ✅     | ✅           | ✅    | ✅          |
| POST /prescriptions                          | ❌      | ✅     | ❌           | ✅    | ✅          |
| PUT /prescriptions/:id                       | ❌      | ✅     | ❌           | ✅    | ✅          |
| PATCH /prescriptions/:id/dispense            | ❌      | ❌     | ✅           | ✅    | ✅          |
| PATCH /prescriptions/:id/cancel              | ❌      | ✅     | ✅           | ✅    | ✅          |
| DELETE /prescriptions/:id                    | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /prescriptions/:id/restore             | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /prescriptions/:id/hard-delete        | ❌      | ❌     | ❌           | ❌    | ✅          |

### Prescription Items

| Endpoint                                 | Patient | Doctor | Receptionist | Admin | Super Admin |
| ---------------------------------------- | ------- | ------ | ------------ | ----- | ----------- |
| GET /prescriptions/:id/items             | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /prescriptions/:id/items/:item_id    | ✅      | ✅     | ✅           | ✅    | ✅          |
| POST /prescriptions/:id/items            | ❌      | ✅     | ❌           | ✅    | ✅          |
| PUT /prescriptions/:id/items/:item_id    | ❌      | ✅     | ❌           | ✅    | ✅          |
| DELETE /prescriptions/:id/items/:item_id | ❌      | ✅     | ❌           | ✅    | ✅          |

---

## Endpoints Summary

### Prescription Endpoints

| Method | Endpoint                                   | Description                     | Role Required                            |
| ------ | ------------------------------------------ | ------------------------------- | ---------------------------------------- |
| GET    | `/prescriptions`                           | List all prescriptions          | Doctor, Receptionist, Admin, Super Admin |
| GET    | `/prescriptions/deleted`                   | List deleted prescriptions      | Admin, Super Admin                       |
| GET    | `/prescriptions/:id`                       | Get prescription by ID          | All Authenticated                        |
| GET    | `/prescriptions/medical-record/:record_id` | Get prescriptions by record     | All Authenticated                        |
| POST   | `/prescriptions`                           | Create new prescription         | Doctor, Admin, Super Admin               |
| PUT    | `/prescriptions/:id`                       | Update prescription             | Doctor, Admin, Super Admin               |
| PATCH  | `/prescriptions/:id/dispense`              | Mark prescription as dispensed  | Receptionist, Admin, Super Admin         |
| PATCH  | `/prescriptions/:id/cancel`                | Cancel prescription             | Doctor, Receptionist, Admin, Super Admin |
| DELETE | `/prescriptions/:id`                       | Soft delete prescription        | Admin, Super Admin                       |
| PATCH  | `/prescriptions/:id/restore`               | Restore deleted prescription    | Admin, Super Admin                       |
| DELETE | `/prescriptions/:id/hard-delete`           | Permanently delete prescription | Super Admin                              |

### Prescription Item Endpoints

| Method | Endpoint                            | Description                 | Role Required              |
| ------ | ----------------------------------- | --------------------------- | -------------------------- |
| GET    | `/prescriptions/:id/items`          | List all prescription items | All Authenticated          |
| GET    | `/prescriptions/:id/items/:item_id` | Get prescription item by ID | All Authenticated          |
| POST   | `/prescriptions/:id/items`          | Add item to prescription    | Doctor, Admin, Super Admin |
| PUT    | `/prescriptions/:id/items/:item_id` | Update prescription item    | Doctor, Admin, Super Admin |
| DELETE | `/prescriptions/:id/items/:item_id` | Delete prescription item    | Doctor, Admin, Super Admin |

---

## Prescription Endpoints

### 1. List All Prescriptions

**Endpoint:** `GET /api/v1/prescriptions`

**Description:** Mendapatkan daftar semua resep dengan pagination, search, dan filter.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

| Header        | Type   | Required | Description      |
| ------------- | ------ | -------- | ---------------- |
| Authorization | string | Yes      | `Bearer <token>` |

**Query Parameters:**

| Parameter | Type    | Required | Default    | Description                                                      |
| --------- | ------- | -------- | ---------- | ---------------------------------------------------------------- |
| page      | integer | No       | 1          | Halaman saat ini                                                 |
| page_size | integer | No       | 10         | Jumlah data per halaman (max: 100)                               |
| sort_by   | string  | No       | created_at | Field sorting: `id`, `prescription_date`, `status`, `created_at` |
| sort_dir  | string  | No       | desc       | Arah sorting: `asc` atau `desc`                                  |
| search    | string  | No       | -          | Pencarian berdasarkan catatan resep                              |
| status    | string  | No       | -          | Filter status: `pending`, `dispensed`, `cancelled`               |
| doctor_id | integer | No       | -          | Filter berdasarkan ID dokter                                     |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Prescriptions retrieved successfully",
  "data": [
    {
      "id": 1,
      "medical_record_id": 5,
      "doctor_id": 2,
      "prescription_date": "2024-01-15",
      "notes": "Minum obat setelah makan",
      "status": "pending",
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
curl -X GET "http://localhost:8080/api/v1/prescriptions?page=1&page_size=10&status=pending" \
  -H "Authorization: Bearer <token>"
```

---

### 2. List Deleted Prescriptions

**Endpoint:** `GET /api/v1/prescriptions/deleted`

**Description:** Mendapatkan daftar resep yang telah dihapus (soft deleted).

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
  "message": "Deleted prescriptions retrieved successfully",
  "data": [
    {
      "id": 2,
      "medical_record_id": 3,
      "doctor_id": 1,
      "prescription_date": "2024-01-10",
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
curl -X GET "http://localhost:8080/api/v1/prescriptions/deleted" \
  -H "Authorization: Bearer <token>"
```

---

### 3. Get Prescription by ID

**Endpoint:** `GET /api/v1/prescriptions/:id`

**Description:** Mendapatkan detail data resep berdasarkan ID, termasuk daftar obat yang diresepkan.

**Authentication:** Required (All Authenticated Users)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID resep    |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Prescription retrieved successfully",
  "data": {
    "id": 1,
    "medical_record_id": 5,
    "doctor_id": 2,
    "prescription_date": "2024-01-15",
    "notes": "Minum obat setelah makan. Habiskan antibiotik.",
    "status": "pending",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z",
    "medical_record": {
      "id": 5,
      "visit_date": "2024-01-15"
    },
    "doctor": {
      "id": 2,
      "name": "dr. Siti Rahayu",
      "specialization": "Penyakit Dalam"
    },
    "items": [
      {
        "id": 1,
        "prescription_id": 1,
        "medicine_id": 3,
        "dosage": "500mg",
        "frequency": "3x sehari",
        "duration_days": 7,
        "quantity": 21,
        "instructions": "Diminum setelah makan",
        "medicine": {
          "id": 3,
          "name": "Amoxicillin",
          "unit": "tablet"
        }
      },
      {
        "id": 2,
        "prescription_id": 1,
        "medicine_id": 7,
        "dosage": "500mg",
        "frequency": "3x sehari",
        "duration_days": 5,
        "quantity": 15,
        "instructions": "Diminum jika nyeri atau demam",
        "medicine": {
          "id": 7,
          "name": "Paracetamol",
          "unit": "tablet"
        }
      }
    ]
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/prescriptions/1" \
  -H "Authorization: Bearer <token>"
```

---

### 4. Create Prescription

**Endpoint:** `POST /api/v1/prescriptions`

**Description:** Membuat resep baru oleh dokter untuk pasien.

**Authentication:** Required (Doctor, Admin, Super Admin)

**Request Headers:**

| Header        | Type   | Required | Description        |
| ------------- | ------ | -------- | ------------------ |
| Authorization | string | Yes      | `Bearer <token>`   |
| Content-Type  | string | Yes      | `application/json` |

**Request Body:**

| Field             | Type    | Required | Description                                                           |
| ----------------- | ------- | -------- | --------------------------------------------------------------------- |
| medical_record_id | integer | Yes      | ID rekam medis terkait                                                |
| doctor_id         | integer | Yes      | ID dokter penulis resep                                               |
| prescription_date | string  | Yes      | Tanggal resep (format: `YYYY-MM-DD`)                                  |
| notes             | string  | No       | Catatan atau instruksi umum untuk pasien                              |
| status            | string  | No       | Status awal: `pending`, `dispensed`, `cancelled` (default: `pending`) |

**Example Request Body:**

```json
{
  "medical_record_id": 5,
  "doctor_id": 2,
  "prescription_date": "2024-01-15",
  "notes": "Minum obat setelah makan. Habiskan antibiotik meski sudah merasa sembuh."
}
```

**Response Success (201):**

```json
{
  "status": "success",
  "message": "Prescription created successfully",
  "data": {
    "id": 1,
    "medical_record_id": 5,
    "doctor_id": 2,
    "prescription_date": "2024-01-15",
    "notes": "Minum obat setelah makan. Habiskan antibiotik meski sudah merasa sembuh.",
    "status": "pending",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/prescriptions" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "medical_record_id": 5,
    "doctor_id": 2,
    "prescription_date": "2024-01-15",
    "notes": "Minum obat setelah makan. Habiskan antibiotik meski sudah merasa sembuh."
  }'
```

---

### 5. Update Prescription

**Endpoint:** `PUT /api/v1/prescriptions/:id`

**Description:** Memperbarui data resep berdasarkan ID.

**Authentication:** Required (Doctor, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID resep    |

**Request Body:** (Same fields as Create, all optional for update)

**Example Request Body:**

```json
{
  "notes": "Minum obat setelah makan. Jika timbul alergi segera hubungi dokter."
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Prescription updated successfully",
  "data": {
    "id": 1,
    "notes": "Minum obat setelah makan. Jika timbul alergi segera hubungi dokter.",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT "http://localhost:8080/api/v1/prescriptions/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Minum obat setelah makan. Jika timbul alergi segera hubungi dokter."
  }'
```

---

### 6. Dispense Prescription

**Endpoint:** `PATCH /api/v1/prescriptions/:id/dispense`

**Description:** Menandai resep telah disiapkan/diserahkan kepada pasien. Status berubah menjadi `dispensed`.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID resep    |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Prescription dispensed successfully",
  "data": {
    "id": 1,
    "status": "dispensed",
    "updated_at": "2024-01-15T12:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/prescriptions/1/dispense" \
  -H "Authorization: Bearer <token>"
```

---

### 7. Cancel Prescription

**Endpoint:** `PATCH /api/v1/prescriptions/:id/cancel`

**Description:** Membatalkan resep. Status berubah menjadi `cancelled`.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID resep    |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Prescription cancelled successfully",
  "data": {
    "id": 1,
    "status": "cancelled",
    "updated_at": "2024-01-15T11:30:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/prescriptions/1/cancel" \
  -H "Authorization: Bearer <token>"
```

---

### 8. Soft Delete Prescription

**Endpoint:** `DELETE /api/v1/prescriptions/:id`

**Description:** Menghapus resep secara soft delete.

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID resep    |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Prescription deleted successfully"
}
```

**cURL Example:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/prescriptions/1" \
  -H "Authorization: Bearer <token>"
```

---

### 9. Restore Deleted Prescription

**Endpoint:** `PATCH /api/v1/prescriptions/:id/restore`

**Description:** Memulihkan resep yang telah dihapus (soft delete).

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID resep    |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Prescription restored successfully",
  "data": {
    "id": 1,
    "status": "pending"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/prescriptions/1/restore" \
  -H "Authorization: Bearer <token>"
```

---

### 10. Hard Delete Prescription

**Endpoint:** `DELETE /api/v1/prescriptions/:id/hard-delete`

**Description:** Menghapus resep secara permanen dari database. Operasi ini tidak dapat dibatalkan.

**Authentication:** Required (Super Admin Only)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID resep    |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Prescription permanently deleted successfully"
}
```

**cURL Example:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/prescriptions/1/hard-delete" \
  -H "Authorization: Bearer <token>"
```

---

## Prescription Item Endpoints

### 1. List All Prescription Items

**Endpoint:** `GET /api/v1/prescriptions/:id/items`

**Description:** Mendapatkan daftar semua item obat dalam sebuah resep.

**Authentication:** Required (All Authenticated Users)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID resep    |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Prescription items retrieved successfully",
  "data": [
    {
      "id": 1,
      "prescription_id": 1,
      "medicine_id": 3,
      "dosage": "500mg",
      "frequency": "3x sehari",
      "duration_days": 7,
      "quantity": 21,
      "instructions": "Diminum setelah makan",
      "medicine": {
        "id": 3,
        "name": "Amoxicillin",
        "unit": "tablet"
      }
    }
  ]
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/prescriptions/1/items" \
  -H "Authorization: Bearer <token>"
```

---

### 2. Get Prescription Item by ID

**Endpoint:** `GET /api/v1/prescriptions/:id/items/:item_id`

**Description:** Mendapatkan detail data sebuah item resep.

**Authentication:** Required (All Authenticated Users)

**Path Parameters:**

| Parameter | Type    | Required | Description     |
| --------- | ------- | -------- | --------------- |
| id        | integer | Yes      | ID resep        |
| item_id   | integer | Yes      | ID item resep   |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Prescription item retrieved successfully",
  "data": {
    "id": 1,
    "prescription_id": 1,
    "medicine_id": 3,
    "dosage": "500mg",
    "frequency": "3x sehari",
    "duration_days": 7,
    "quantity": 21,
    "instructions": "Diminum setelah makan",
    "medicine": {
      "id": 3,
      "name": "Amoxicillin",
      "unit": "tablet"
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/prescriptions/1/items/1" \
  -H "Authorization: Bearer <token>"
```

---

### 3. Add Item to Prescription

**Endpoint:** `POST /api/v1/prescriptions/:id/items`

**Description:** Menambahkan item obat ke dalam resep.

**Authentication:** Required (Doctor, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID resep    |

**Request Headers:**

| Header        | Type   | Required | Description        |
| ------------- | ------ | -------- | ------------------ |
| Authorization | string | Yes      | `Bearer <token>`   |
| Content-Type  | string | Yes      | `application/json` |

**Request Body:**

| Field         | Type    | Required | Description                      |
| ------------- | ------- | -------- | -------------------------------- |
| medicine_id   | integer | Yes      | ID obat                          |
| dosage        | string  | Yes      | Dosis obat (misal: 500mg)        |
| frequency     | string  | Yes      | Frekuensi (misal: 3x sehari)     |
| duration_days | integer | Yes      | Durasi penggunaan dalam hari     |
| quantity      | integer | Yes      | Jumlah obat yang diberikan       |
| instructions  | string  | No       | Instruksi tambahan               |

**Example Request Body:**

```json
{
  "medicine_id": 3,
  "dosage": "500mg",
  "frequency": "3x sehari",
  "duration_days": 7,
  "quantity": 21,
  "instructions": "Diminum setelah makan"
}
```

**Response Success (201):**

```json
{
  "status": "success",
  "message": "Prescription item added successfully",
  "data": {
    "id": 1,
    "prescription_id": 1,
    "medicine_id": 3,
    "dosage": "500mg",
    "frequency": "3x sehari",
    "duration_days": 7,
    "quantity": 21,
    "instructions": "Diminum setelah makan"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/prescriptions/1/items" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "medicine_id": 3,
    "dosage": "500mg",
    "frequency": "3x sehari",
    "duration_days": 7,
    "quantity": 21,
    "instructions": "Diminum setelah makan"
  }'
```

---

### 4. Update Prescription Item

**Endpoint:** `PUT /api/v1/prescriptions/:id/items/:item_id`

**Description:** Memperbarui data item resep.

**Authentication:** Required (Doctor, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description     |
| --------- | ------- | -------- | --------------- |
| id        | integer | Yes      | ID resep        |
| item_id   | integer | Yes      | ID item resep   |

**Request Body:** (Semua field optional untuk update, `medicine_id` tidak dapat diubah)

| Field         | Type    | Required | Description                      |
| ------------- | ------- | -------- | -------------------------------- |
| dosage        | string  | No       | Dosis obat (misal: 500mg)        |
| frequency     | string  | No       | Frekuensi (misal: 3x sehari)     |
| duration_days | integer | No       | Durasi penggunaan dalam hari     |
| quantity      | integer | No       | Jumlah obat yang diberikan       |
| instructions  | string  | No       | Instruksi tambahan               |

**Example Request Body:**

```json
{
  "dosage": "500mg",
  "frequency": "3x sehari",
  "duration_days": 7,
  "quantity": 21,
  "instructions": "Diminum setelah makan"
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Prescription item updated successfully",
  "data": {
    "id": 1,
    "prescription_id": 1,
    "medicine_id": 3,
    "dosage": "250mg",
    "frequency": "3x sehari",
    "duration_days": 7,
    "quantity": 14,
    "instructions": "Diminum setelah makan"
  }
}
```

**cURL Example:**

```bash
curl -X PUT "http://localhost:8080/api/v1/prescriptions/1/items/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "dosage": "250mg",
    "quantity": 14
  }'
```

---

### 5. Delete Prescription Item

**Endpoint:** `DELETE /api/v1/prescriptions/:id/items/:item_id`

**Description:** Menghapus item dari resep.

**Authentication:** Required (Doctor, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description     |
| --------- | ------- | -------- | --------------- |
| id        | integer | Yes      | ID resep        |
| item_id   | integer | Yes      | ID item resep   |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Prescription item deleted successfully"
}
```

**cURL Example:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/prescriptions/1/items/1" \
  -H "Authorization: Bearer <token>"
```

---
