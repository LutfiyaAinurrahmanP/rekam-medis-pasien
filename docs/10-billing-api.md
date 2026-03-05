# Billing API Documentation

## Overview

API untuk manajemen data billing (tagihan/invoice) dan item tagihan dalam sistem rekam medis. Billing mencatat seluruh biaya layanan kesehatan yang diberikan kepada pasien beserta rincian item biaya dan status pembayarannya.

**Base URL:** `/api/v1/billing`

---

## Table of Contents

- [Authentication](#authentication)
- [Authorization](#authorization)
- [Endpoints Summary](#endpoints-summary)
- [Billing Endpoints](#billing-endpoints)
- [Billing Item Endpoints](#billing-item-endpoints)
- [Error Responses](#error-responses)

---

## Authentication

Semua endpoints memerlukan JWT token di header:

```
Authorization: Bearer <your-jwt-token>
```

---

## Authorization

### Billing

| Endpoint                             | Patient | Doctor | Receptionist | Admin | Super Admin |
| ------------------------------------ | ------- | ------ | ------------ | ----- | ----------- |
| GET /billing                         | ❌      | ✅     | ✅           | ✅    | ✅          |
| GET /billing/deleted                 | ❌      | ❌     | ❌           | ✅    | ✅          |
| GET /billing/:id                     | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /billing/patient/:patient_id     | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /billing/invoice/:invoice_number | ✅      | ✅     | ✅           | ✅    | ✅          |
| POST /billing                        | ❌      | ❌     | ✅           | ✅    | ✅          |
| PUT /billing/:id                     | ❌      | ❌     | ✅           | ✅    | ✅          |
| PATCH /billing/:id/pay               | ❌      | ❌     | ✅           | ✅    | ✅          |
| PATCH /billing/:id/cancel            | ❌      | ❌     | ✅           | ✅    | ✅          |
| DELETE /billing/:id                  | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /billing/:id/restore           | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /billing/:id/hard-delete      | ❌      | ❌     | ❌           | ❌    | ✅          |

### Billing Items

| Endpoint                           | Patient | Doctor | Receptionist | Admin | Super Admin |
| ---------------------------------- | ------- | ------ | ------------ | ----- | ----------- |
| GET /billing/:id/items             | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /billing/:id/items/:item_id    | ✅      | ✅     | ✅           | ✅    | ✅          |
| POST /billing/:id/items            | ❌      | ❌     | ✅           | ✅    | ✅          |
| PUT /billing/:id/items/:item_id    | ❌      | ❌     | ✅           | ✅    | ✅          |
| DELETE /billing/:id/items/:item_id | ❌      | ❌     | ❌           | ✅    | ✅          |

---

## Endpoints Summary

### Billing Endpoints

| Method | Endpoint                           | Description                   | Role Required                            |
| ------ | ---------------------------------- | ----------------------------- | ---------------------------------------- |
| GET    | `/billing`                         | List all billing records      | Doctor, Receptionist, Admin, Super Admin |
| GET    | `/billing/deleted`                 | List deleted billing records  | Admin, Super Admin                       |
| GET    | `/billing/:id`                     | Get billing by ID             | All Authenticated                        |
| GET    | `/billing/patient/:patient_id`     | Get billing by patient        | All Authenticated                        |
| GET    | `/billing/invoice/:invoice_number` | Get billing by invoice number | All Authenticated                        |
| POST   | `/billing`                         | Create new billing            | Receptionist, Admin, Super Admin         |
| PUT    | `/billing/:id`                     | Update billing                | Receptionist, Admin, Super Admin         |
| PATCH  | `/billing/:id/pay`                 | Record payment                | Receptionist, Admin, Super Admin         |
| PATCH  | `/billing/:id/cancel`              | Cancel billing                | Receptionist, Admin, Super Admin         |
| DELETE | `/billing/:id`                     | Soft delete billing           | Admin, Super Admin                       |
| PATCH  | `/billing/:id/restore`             | Restore deleted billing       | Admin, Super Admin                       |
| DELETE | `/billing/:id/hard-delete`         | Permanently delete billing    | Super Admin                              |

### Billing Item Endpoints

| Method | Endpoint                      | Description               | Role Required                    |
| ------ | ----------------------------- | ------------------------- | -------------------------------- |
| GET    | `/billing/:id/items`          | List all items in billing | All Authenticated                |
| GET    | `/billing/:id/items/:item_id` | Get billing item by ID    | All Authenticated                |
| POST   | `/billing/:id/items`          | Add item to billing       | Receptionist, Admin, Super Admin |
| PUT    | `/billing/:id/items/:item_id` | Update billing item       | Receptionist, Admin, Super Admin |
| DELETE | `/billing/:id/items/:item_id` | Delete billing item       | Admin, Super Admin               |

---

## Billing Endpoints

### 1. List All Billing Records

**Endpoint:** `GET /api/v1/billing`

**Description:** Mendapatkan daftar semua data billing dengan pagination, search, dan filter.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

| Header        | Type   | Required | Description      |
| ------------- | ------ | -------- | ---------------- |
| Authorization | string | Yes      | `Bearer <token>` |

**Query Parameters:**

| Parameter      | Type    | Required | Default    | Description                                                                                           |
| -------------- | ------- | -------- | ---------- | ----------------------------------------------------------------------------------------------------- |
| page           | integer | No       | 1          | Halaman saat ini                                                                                      |
| page_size      | integer | No       | 10         | Jumlah data per halaman (max: 100)                                                                    |
| sort_by        | string  | No       | created_at | Field sorting: `id`, `invoice_number`, `invoice_date`, `total_amount`, `payment_status`, `created_at` |
| sort_dir       | string  | No       | desc       | Arah sorting: `asc` atau `desc`                                                                       |
| search         | string  | No       | -          | Pencarian berdasarkan invoice number atau notes                                                       |
| payment_status | string  | No       | -          | Filter status: `unpaid`, `partial`, `paid`, `cancelled`                                               |
| payment_method | string  | No       | -          | Filter metode: `cash`, `debit_card`, `credit_card`, `insurance`, `transfer`, `other`                  |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Billing records retrieved successfully",
  "data": [
    {
      "id": 1,
      "patient_id": 10,
      "medical_record_id": 5,
      "invoice_number": "INV-2024-000001",
      "invoice_date": "2024-01-15",
      "total_amount": 500000.0,
      "paid_amount": 0.0,
      "discount_amount": 50000.0,
      "payment_status": "unpaid",
      "payment_method": null,
      "payment_date": null,
      "notes": "Pembayaran layanan rawat jalan",
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
curl -X GET "http://localhost:8080/api/v1/billing?page=1&page_size=10&payment_status=unpaid" \
  -H "Authorization: Bearer <token>"
```

---

### 2. List Deleted Billing Records

**Endpoint:** `GET /api/v1/billing/deleted`

**Description:** Mendapatkan daftar billing yang telah dihapus (soft deleted).

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
  "message": "Deleted billing records retrieved successfully",
  "data": [
    {
      "id": 2,
      "invoice_number": "INV-2024-000002",
      "invoice_date": "2024-01-10",
      "total_amount": 250000.0,
      "payment_status": "cancelled",
      "created_at": "2024-01-10T09:00:00Z",
      "updated_at": "2024-01-12T10:00:00Z"
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
curl -X GET "http://localhost:8080/api/v1/billing/deleted" \
  -H "Authorization: Bearer <token>"
```

---

### 3. Get Billing by ID

**Endpoint:** `GET /api/v1/billing/:id`

**Description:** Mendapatkan detail data billing berdasarkan ID, termasuk item-item billing.

**Authentication:** Required (All Authenticated Users)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID billing  |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Billing retrieved successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "medical_record_id": 5,
    "invoice_number": "INV-2024-000001",
    "invoice_date": "2024-01-15",
    "total_amount": 500000.0,
    "paid_amount": 500000.0,
    "discount_amount": 50000.0,
    "payment_status": "paid",
    "payment_method": "cash",
    "payment_date": "2024-01-15T10:30:00Z",
    "notes": "Pembayaran layanan rawat jalan",
    "created_at": "2024-01-15T08:00:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "patient": {
      "id": 10,
      "name": "Budi Santoso",
      "medical_record_number": "MRN-2024-000010"
    },
    "medical_record": {
      "id": 5,
      "visit_date": "2024-01-15"
    },
    "items": [
      {
        "id": 1,
        "billing_id": 1,
        "item_type": "consultation",
        "item_description": "Biaya konsultasi dokter umum",
        "quantity": 1,
        "unit_price": 150000.0,
        "total_price": 150000.0
      },
      {
        "id": 2,
        "billing_id": 1,
        "item_type": "medicine",
        "item_description": "Paracetamol 500mg",
        "quantity": 10,
        "unit_price": 5000.0,
        "total_price": 50000.0
      }
    ]
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/billing/1" \
  -H "Authorization: Bearer <token>"
```

---

### 4. Get Billing by Patient

**Endpoint:** `GET /api/v1/billing/patient/:patient_id`

**Description:** Mendapatkan semua data billing berdasarkan ID pasien.

**Authentication:** Required (All Authenticated Users)

**Path Parameters:**

| Parameter  | Type    | Required | Description |
| ---------- | ------- | -------- | ----------- |
| patient_id | integer | Yes      | ID pasien   |

**Query Parameters:**

| Parameter      | Type    | Required | Default | Description                                             |
| -------------- | ------- | -------- | ------- | ------------------------------------------------------- |
| page           | integer | No       | 1       | Halaman saat ini                                        |
| page_size      | integer | No       | 10      | Jumlah data per halaman (max: 100)                      |
| payment_status | string  | No       | -       | Filter status: `unpaid`, `partial`, `paid`, `cancelled` |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Patient billing records retrieved successfully",
  "data": [
    {
      "id": 1,
      "invoice_number": "INV-2024-000001",
      "invoice_date": "2024-01-15",
      "total_amount": 500000.0,
      "paid_amount": 500000.0,
      "payment_status": "paid",
      "payment_method": "cash",
      "payment_date": "2024-01-15T10:30:00Z"
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
curl -X GET "http://localhost:8080/api/v1/billing/patient/10?payment_status=unpaid" \
  -H "Authorization: Bearer <token>"
```

---

### 5. Get Billing by Invoice Number

**Endpoint:** `GET /api/v1/billing/invoice/:invoice_number`

**Description:** Mendapatkan detail data billing berdasarkan nomor invoice.

**Authentication:** Required (All Authenticated Users)

**Path Parameters:**

| Parameter      | Type   | Required | Description   |
| -------------- | ------ | -------- | ------------- |
| invoice_number | string | Yes      | Nomor invoice |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Billing retrieved successfully",
  "data": {
    "id": 1,
    "invoice_number": "INV-2024-000001",
    "invoice_date": "2024-01-15",
    "total_amount": 500000.0,
    "paid_amount": 0.0,
    "payment_status": "unpaid"
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/billing/invoice/INV-2024-000001" \
  -H "Authorization: Bearer <token>"
```

---

### 6. Create Billing

**Endpoint:** `POST /api/v1/billing`

**Description:** Membuat data billing baru untuk pasien.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Request Headers:**

| Header        | Type   | Required | Description        |
| ------------- | ------ | -------- | ------------------ |
| Authorization | string | Yes      | `Bearer <token>`   |
| Content-Type  | string | Yes      | `application/json` |

**Request Body:**

| Field             | Type    | Required | Description                                                                   |
| ----------------- | ------- | -------- | ----------------------------------------------------------------------------- |
| patient_id        | integer | Yes      | ID pasien                                                                     |
| medical_record_id | integer | No       | ID rekam medis terkait                                                        |
| invoice_number    | string  | Yes      | Nomor invoice (unik, max 50 karakter)                                         |
| invoice_date      | string  | Yes      | Tanggal invoice (format: `YYYY-MM-DD`)                                        |
| total_amount      | float   | Yes      | Total jumlah tagihan (>= 0)                                                   |
| paid_amount       | float   | No       | Jumlah yang sudah dibayar (default: 0)                                        |
| discount_amount   | float   | No       | Jumlah diskon                                                                 |
| payment_status    | string  | Yes      | Status: `unpaid`, `partial`, `paid`, `cancelled` (default: `unpaid`)          |
| payment_method    | string  | No       | Metode: `cash`, `debit_card`, `credit_card`, `insurance`, `transfer`, `other` |
| payment_date      | string  | No       | Tanggal pembayaran (format: `YYYY-MM-DDTHH:MM:SSZ`)                           |
| notes             | string  | No       | Catatan tambahan                                                              |

**Example Request Body:**

```json
{
  "patient_id": 10,
  "medical_record_id": 5,
  "invoice_number": "INV-2024-000001",
  "invoice_date": "2024-01-15",
  "total_amount": 500000.0,
  "discount_amount": 50000.0,
  "payment_status": "unpaid",
  "notes": "Pembayaran layanan rawat jalan"
}
```

**Response Success (201):**

```json
{
  "status": "success",
  "message": "Billing created successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "medical_record_id": 5,
    "invoice_number": "INV-2024-000001",
    "invoice_date": "2024-01-15",
    "total_amount": 500000.0,
    "paid_amount": 0.0,
    "discount_amount": 50000.0,
    "payment_status": "unpaid",
    "created_at": "2024-01-15T08:00:00Z",
    "updated_at": "2024-01-15T08:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/billing" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 10,
    "medical_record_id": 5,
    "invoice_number": "INV-2024-000001",
    "invoice_date": "2024-01-15",
    "total_amount": 500000.00,
    "discount_amount": 50000.00,
    "payment_status": "unpaid",
    "notes": "Pembayaran layanan rawat jalan"
  }'
```

---

### 7. Update Billing

**Endpoint:** `PUT /api/v1/billing/:id`

**Description:** Memperbarui data billing berdasarkan ID.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID billing  |

**Request Body:** (Same fields as Create, all optional for update)

**Example Request Body:**

```json
{
  "total_amount": 550000.0,
  "discount_amount": 25000.0,
  "notes": "Updated: Biaya tambahan tindakan"
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Billing updated successfully",
  "data": {
    "id": 1,
    "invoice_number": "INV-2024-000001",
    "total_amount": 550000.0,
    "discount_amount": 25000.0,
    "payment_status": "unpaid",
    "updated_at": "2024-01-15T09:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT "http://localhost:8080/api/v1/billing/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "total_amount": 550000.00,
    "discount_amount": 25000.00,
    "notes": "Updated: Biaya tambahan tindakan"
  }'
```

---

### 8. Record Payment

**Endpoint:** `PATCH /api/v1/billing/:id/pay`

**Description:** Mencatat pembayaran untuk billing. Status pembayaran akan diperbarui secara otomatis berdasarkan jumlah yang dibayar.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID billing  |

**Request Body:**

| Field          | Type   | Required | Description                                                                   |
| -------------- | ------ | -------- | ----------------------------------------------------------------------------- |
| paid_amount    | float  | Yes      | Jumlah pembayaran (>= 0)                                                      |
| payment_method | string | Yes      | Metode: `cash`, `debit_card`, `credit_card`, `insurance`, `transfer`, `other` |

**Example Request Body:**

```json
{
  "paid_amount": 500000.0,
  "payment_method": "cash"
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Payment recorded successfully",
  "data": {
    "id": 1,
    "invoice_number": "INV-2024-000001",
    "total_amount": 500000.0,
    "paid_amount": 500000.0,
    "payment_status": "paid",
    "payment_method": "cash",
    "payment_date": "2024-01-15T10:30:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/billing/1/pay" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "paid_amount": 500000.00,
    "payment_method": "cash"
  }'
```

> **Notes:**
>
> - Jika `paid_amount` >= `total_amount - discount_amount`, status akan berubah menjadi `paid`
> - Jika `paid_amount` > 0 dan < net amount, status akan berubah menjadi `partial`
> - `payment_date` diisi otomatis saat status berubah menjadi `paid`

---

### 9. Cancel Billing

**Endpoint:** `PATCH /api/v1/billing/:id/cancel`

**Description:** Membatalkan billing. Status pembayaran akan diubah menjadi `cancelled`.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID billing  |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Billing cancelled successfully",
  "data": {
    "id": 1,
    "invoice_number": "INV-2024-000001",
    "payment_status": "cancelled",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/billing/1/cancel" \
  -H "Authorization: Bearer <token>"
```

---

### 10. Soft Delete Billing

**Endpoint:** `DELETE /api/v1/billing/:id`

**Description:** Menghapus billing secara soft delete (data masih tersimpan di database).

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID billing  |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Billing deleted successfully"
}
```

**cURL Example:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/billing/1" \
  -H "Authorization: Bearer <token>"
```

---

### 11. Restore Deleted Billing

**Endpoint:** `PATCH /api/v1/billing/:id/restore`

**Description:** Memulihkan billing yang telah dihapus (soft delete).

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID billing  |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Billing restored successfully",
  "data": {
    "id": 1,
    "invoice_number": "INV-2024-000001",
    "payment_status": "unpaid"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH "http://localhost:8080/api/v1/billing/1/restore" \
  -H "Authorization: Bearer <token>"
```

---

### 12. Hard Delete Billing

**Endpoint:** `DELETE /api/v1/billing/:id/hard-delete`

**Description:** Menghapus billing secara permanen dari database. Operasi ini tidak dapat dibatalkan.

**Authentication:** Required (Super Admin Only)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID billing  |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Billing permanently deleted successfully"
}
```

**cURL Example:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/billing/1/hard-delete" \
  -H "Authorization: Bearer <token>"
```

---

## Billing Item Endpoints

### 1. List Billing Items

**Endpoint:** `GET /api/v1/billing/:id/items`

**Description:** Mendapatkan semua item dalam sebuah billing.

**Authentication:** Required (All Authenticated Users)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID billing  |

**Query Parameters:**

| Parameter | Type   | Required | Default | Description                                                                       |
| --------- | ------ | -------- | ------- | --------------------------------------------------------------------------------- |
| item_type | string | No       | -       | Filter tipe: `consultation`, `medicine`, `lab_test`, `procedure`, `room`, `other` |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Billing items retrieved successfully",
  "data": [
    {
      "id": 1,
      "billing_id": 1,
      "item_type": "consultation",
      "item_description": "Biaya konsultasi dokter umum",
      "quantity": 1,
      "unit_price": 150000.0,
      "total_price": 150000.0,
      "created_at": "2024-01-15T08:00:00Z",
      "updated_at": "2024-01-15T08:00:00Z"
    },
    {
      "id": 2,
      "billing_id": 1,
      "item_type": "medicine",
      "item_description": "Paracetamol 500mg",
      "quantity": 10,
      "unit_price": 5000.0,
      "total_price": 50000.0,
      "created_at": "2024-01-15T08:05:00Z",
      "updated_at": "2024-01-15T08:05:00Z"
    }
  ]
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/billing/1/items?item_type=medicine" \
  -H "Authorization: Bearer <token>"
```

---

### 2. Get Billing Item by ID

**Endpoint:** `GET /api/v1/billing/:id/items/:item_id`

**Description:** Mendapatkan detail satu item billing berdasarkan ID.

**Authentication:** Required (All Authenticated Users)

**Path Parameters:**

| Parameter | Type    | Required | Description     |
| --------- | ------- | -------- | --------------- |
| id        | integer | Yes      | ID billing      |
| item_id   | integer | Yes      | ID item billing |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Billing item retrieved successfully",
  "data": {
    "id": 1,
    "billing_id": 1,
    "item_type": "consultation",
    "item_description": "Biaya konsultasi dokter umum",
    "quantity": 1,
    "unit_price": 150000.0,
    "total_price": 150000.0,
    "created_at": "2024-01-15T08:00:00Z",
    "updated_at": "2024-01-15T08:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/billing/1/items/1" \
  -H "Authorization: Bearer <token>"
```

---

### 3. Add Billing Item

**Endpoint:** `POST /api/v1/billing/:id/items`

**Description:** Menambahkan item baru ke dalam billing. `total_price` dihitung otomatis dari `quantity * unit_price`.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description |
| --------- | ------- | -------- | ----------- |
| id        | integer | Yes      | ID billing  |

**Request Body:**

| Field            | Type    | Required | Description                                                                     |
| ---------------- | ------- | -------- | ------------------------------------------------------------------------------- |
| item_type        | string  | Yes      | Tipe item: `consultation`, `medicine`, `lab_test`, `procedure`, `room`, `other` |
| item_description | string  | Yes      | Deskripsi item (max 255 karakter)                                               |
| quantity         | integer | Yes      | Jumlah item (min: 1)                                                            |
| unit_price       | float   | Yes      | Harga per unit (>= 0)                                                           |

**Example Request Body:**

```json
{
  "item_type": "lab_test",
  "item_description": "Pemeriksaan darah lengkap",
  "quantity": 1,
  "unit_price": 200000.0
}
```

**Response Success (201):**

```json
{
  "status": "success",
  "message": "Billing item added successfully",
  "data": {
    "id": 3,
    "billing_id": 1,
    "item_type": "lab_test",
    "item_description": "Pemeriksaan darah lengkap",
    "quantity": 1,
    "unit_price": 200000.0,
    "total_price": 200000.0,
    "created_at": "2024-01-15T09:00:00Z",
    "updated_at": "2024-01-15T09:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/billing/1/items" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "item_type": "lab_test",
    "item_description": "Pemeriksaan darah lengkap",
    "quantity": 1,
    "unit_price": 200000.00
  }'
```

---

### 4. Update Billing Item

**Endpoint:** `PUT /api/v1/billing/:id/items/:item_id`

**Description:** Memperbarui data item billing. `total_price` dihitung ulang secara otomatis.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description     |
| --------- | ------- | -------- | --------------- |
| id        | integer | Yes      | ID billing      |
| item_id   | integer | Yes      | ID item billing |

**Request Body:** (Same fields as Add Billing Item, all optional for update)

**Example Request Body:**

```json
{
  "quantity": 2,
  "unit_price": 200000.0
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Billing item updated successfully",
  "data": {
    "id": 3,
    "billing_id": 1,
    "item_type": "lab_test",
    "item_description": "Pemeriksaan darah lengkap",
    "quantity": 2,
    "unit_price": 200000.0,
    "total_price": 400000.0,
    "updated_at": "2024-01-15T09:30:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT "http://localhost:8080/api/v1/billing/1/items/3" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 2,
    "unit_price": 200000.00
  }'
```

---

### 5. Delete Billing Item

**Endpoint:** `DELETE /api/v1/billing/:id/items/:item_id`

**Description:** Menghapus item dari billing (soft delete).

**Authentication:** Required (Admin, Super Admin)

**Path Parameters:**

| Parameter | Type    | Required | Description     |
| --------- | ------- | -------- | --------------- |
| id        | integer | Yes      | ID billing      |
| item_id   | integer | Yes      | ID item billing |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Billing item deleted successfully"
}
```

**cURL Example:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/billing/1/items/3" \
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
    "invoice_number": "invoice_number is required",
    "total_amount": "total_amount must be greater than or equal to 0"
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
  "message": "Billing not found"
}
```

### 409 Conflict

```json
{
  "status": "error",
  "message": "Invoice number already exists"
}
```

### 422 Unprocessable Entity

```json
{
  "status": "error",
  "message": "Cannot add item to a cancelled billing"
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

## Payment Status Flow

```
unpaid → partial → paid
unpaid → cancelled
partial → paid
partial → cancelled
```

**Payment Status Values:**

- `unpaid` — Belum ada pembayaran
- `partial` — Pembayaran sebagian (paid_amount < net amount)
- `paid` — Lunas (paid_amount >= total_amount - discount_amount)
- `cancelled` — Billing dibatalkan

**Payment Method Values:**

- `cash` — Tunai
- `debit_card` — Kartu debit
- `credit_card` — Kartu kredit
- `insurance` — Asuransi
- `transfer` — Transfer bank
- `other` — Lainnya

**Billing Item Type Values:**

- `consultation` — Biaya konsultasi
- `medicine` — Obat-obatan
- `lab_test` — Pemeriksaan laboratorium
- `procedure` — Tindakan medis
- `room` — Biaya kamar/rawat inap
- `other` — Lainnya
