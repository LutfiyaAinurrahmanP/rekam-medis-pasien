# Medicines API Documentation

## Overview

API untuk manajemen data medicines (obat-obatan) dalam sistem rekam medis. Medicines adalah master data untuk obat-obatan yang tersedia di rumah sakit, termasuk stock management.

**Base URL:** `/api/v1/medicines`

---

## Table of Contents

- [Authentication](#authentication)
- [Authorization](#authorization)
- [Endpoints Summary](#endpoints-summary)
- [Public Endpoints](#public-endpoints)
- [Pharmacist Endpoints](#pharmacist-endpoints)
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

| Endpoint                          | Patient | Doctor | Receptionist | Admin | Super Admin |
| --------------------------------- | ------- | ------ | ------------ | ----- | ----------- |
| GET /medicines                    | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /medicines/available          | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /medicines/low-stock          | ❌      | ✅     | ✅           | ✅    | ✅          |
| GET /medicines/out-of-stock       | ❌      | ✅     | ✅           | ✅    | ✅          |
| GET /medicines/inactive           | ❌      | ❌     | ✅           | ✅    | ✅          |
| GET /medicines/deleted            | ❌      | ❌     | ❌           | ✅    | ✅          |
| GET /medicines/:id                | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /medicines/name/:name         | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /medicines/type/:type         | ✅      | ✅     | ✅           | ✅    | ✅          |
| GET /medicines/search             | ✅      | ✅     | ✅           | ✅    | ✅          |
| POST /medicines                   | ❌      | ❌     | ❌           | ✅    | ✅          |
| PUT /medicines/:id                | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /medicines/:id/add-stock    | ❌      | ❌     | ✅           | ✅    | ✅          |
| PATCH /medicines/:id/reduce-stock | ❌      | ❌     | ✅           | ✅    | ✅          |
| PATCH /medicines/:id/activate     | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /medicines/:id/deactivate   | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /medicines/:id             | ❌      | ❌     | ❌           | ✅    | ✅          |
| PATCH /medicines/:id/restore      | ❌      | ❌     | ❌           | ✅    | ✅          |
| DELETE /medicines/:id/hard-delete | ❌      | ❌     | ❌           | ❌    | ✅          |

---

## Endpoints Summary

### Public Endpoints (All Authenticated)

| Method | Endpoint                  | Description                 | Role Required                            |
| ------ | ------------------------- | --------------------------- | ---------------------------------------- |
| GET    | `/medicines`              | List all medicines          | All Authenticated                        |
| GET    | `/medicines/available`    | List available medicines    | All Authenticated                        |
| GET    | `/medicines/low-stock`    | List low stock medicines    | Doctor, Receptionist, Admin, Super Admin |
| GET    | `/medicines/out-of-stock` | List out of stock medicines | Doctor, Receptionist, Admin, Super Admin |
| GET    | `/medicines/inactive`     | List inactive medicines     | Receptionist, Admin, Super Admin         |
| GET    | `/medicines/deleted`      | List deleted medicines      | Admin, Super Admin                       |
| GET    | `/medicines/:id`          | Get medicine by ID          | All Authenticated                        |
| GET    | `/medicines/name/:name`   | Get medicine by name        | All Authenticated                        |
| GET    | `/medicines/type/:type`   | Get medicines by type       | All Authenticated                        |
| GET    | `/medicines/search`       | Advanced search             | All Authenticated                        |

### Pharmacist/Stock Endpoints

| Method | Endpoint                      | Description  | Role Required                    |
| ------ | ----------------------------- | ------------ | -------------------------------- |
| PATCH  | `/medicines/:id/add-stock`    | Add stock    | Receptionist, Admin, Super Admin |
| PATCH  | `/medicines/:id/reduce-stock` | Reduce stock | Receptionist, Admin, Super Admin |

### Admin Endpoints

| Method | Endpoint                    | Description              | Role Required      |
| ------ | --------------------------- | ------------------------ | ------------------ |
| POST   | `/medicines`                | Create medicine          | Admin, Super Admin |
| PUT    | `/medicines/:id`            | Update medicine          | Admin, Super Admin |
| PATCH  | `/medicines/:id/activate`   | Activate medicine        | Admin, Super Admin |
| PATCH  | `/medicines/:id/deactivate` | Deactivate medicine      | Admin, Super Admin |
| DELETE | `/medicines/:id`            | Soft delete medicine     | Admin, Super Admin |
| PATCH  | `/medicines/:id/restore`    | Restore deleted medicine | Admin, Super Admin |

### Super Admin Endpoints

| Method | Endpoint                     | Description                 | Role Required |
| ------ | ---------------------------- | --------------------------- | ------------- |
| DELETE | `/medicines/:id/hard-delete` | Permanently delete medicine | Super Admin   |

---

## Public Endpoints

### 1. List Medicines

**Endpoint:** `GET /api/v1/medicines`

**Description:** Mendapatkan daftar semua medicines dengan pagination, search, dan filter.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter      | Type    | Default | Description                                                         |
| -------------- | ------- | ------- | ------------------------------------------------------------------- |
| `page`         | integer | 1       | Halaman                                                             |
| `page_size`    | integer | 10      | Jumlah data per halaman (max: 100)                                  |
| `search`       | string  | -       | Cari berdasarkan name, generic_name, brand_name                     |
| `type`         | string  | -       | Filter by type (tablet, capsule, syrup, injection, ointment, other) |
| `manufacturer` | string  | -       | Filter by manufacturer                                              |
| `is_active`    | boolean | -       | Filter by active status                                             |
| `has_stock`    | boolean | -       | Filter medicines with stock > 0                                     |
| `min_stock`    | integer | -       | Filter minimum stock                                                |
| `max_stock`    | integer | -       | Filter maximum stock                                                |
| `sort_by`      | string  | name    | Sort field (name, generic_name, stock_quantity, price, created_at)  |
| `sort_dir`     | string  | asc     | Sort direction (asc, desc)                                          |

**Example Request:**

```
GET /api/v1/medicines?page=1&page_size=10&type=tablet&has_stock=true&sort_by=name&sort_dir=asc
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicines retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "name": "Paracetamol 500mg",
        "generic_name": "Paracetamol",
        "brand_name": "Sanmol",
        "type": "tablet",
        "strength": "500mg",
        "manufacturer": "PT. Sanbe Farma",
        "unit": "tablet",
        "stock_quantity": 1000,
        "price": 500.0,
        "is_active": true,
        "is_low_stock": false,
        "is_out_of_stock": false,
        "created_at": "2024-01-19T10:00:00Z"
      },
      {
        "id": 2,
        "name": "Amoxicillin 500mg",
        "generic_name": "Amoxicillin",
        "brand_name": "Amoxan",
        "type": "capsule",
        "strength": "500mg",
        "manufacturer": "PT. Kimia Farma",
        "unit": "capsule",
        "stock_quantity": 500,
        "price": 2000.0,
        "is_active": true,
        "is_low_stock": false,
        "is_out_of_stock": false,
        "created_at": "2024-01-19T10:05:00Z"
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
curl -X GET "http://localhost:8080/api/v1/medicines?page=1&page_size=10&type=tablet" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 2. List Available Medicines

**Endpoint:** `GET /api/v1/medicines/available`

**Description:** Mendapatkan daftar medicines yang available (active & in stock).

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter      | Type    | Description            |
| -------------- | ------- | ---------------------- |
| `type`         | string  | Filter by type         |
| `manufacturer` | string  | Filter by manufacturer |
| `min_stock`    | integer | Minimum stock quantity |
| `page`         | integer | Page number            |
| `page_size`    | integer | Items per page         |

**Example Request:**

```
GET /api/v1/medicines/available?type=tablet&min_stock=100
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Available medicines retrieved successfully",
  "data": {
    "total_available": 150,
    "total_stock_value": 50000000.0,
    "medicine_types": [
      {
        "type": "tablet",
        "count": 60
      },
      {
        "type": "capsule",
        "count": 40
      },
      {
        "type": "syrup",
        "count": 30
      },
      {
        "type": "injection",
        "count": 15
      },
      {
        "type": "ointment",
        "count": 5
      }
    ],
    "data": [
      {
        "id": 1,
        "name": "Paracetamol 500mg",
        "generic_name": "Paracetamol",
        "type": "tablet",
        "strength": "500mg",
        "stock_quantity": 1000,
        "price": 500.0,
        "is_active": true
      }
    ],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total_items": 150,
      "total_pages": 15
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/medicines/available?type=tablet" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Use Case:**

- Doctor prescribing medicines
- Pharmacist dispensing medicines
- Patient inquiry available medicines

---

### 3. List Low Stock Medicines

**Endpoint:** `GET /api/v1/medicines/low-stock`

**Description:** Mendapatkan daftar medicines dengan stock rendah (perlu restock).

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter   | Type    | Description                   |
| ----------- | ------- | ----------------------------- |
| `threshold` | integer | Stock threshold (default: 50) |
| `page`      | integer | Page number                   |
| `page_size` | integer | Items per page                |

**Example Request:**

```
GET /api/v1/medicines/low-stock?threshold=100
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Low stock medicines retrieved successfully",
  "data": {
    "threshold": 100,
    "total_low_stock": 15,
    "critical_count": 5,
    "data": [
      {
        "id": 10,
        "name": "Insulin Glargine 100IU/ml",
        "generic_name": "Insulin Glargine",
        "type": "injection",
        "stock_quantity": 10,
        "minimum_stock": 50,
        "price": 350000.0,
        "stock_status": "critical",
        "days_until_stockout": 3,
        "is_active": true
      },
      {
        "id": 15,
        "name": "Salbutamol Inhaler",
        "generic_name": "Salbutamol",
        "type": "other",
        "stock_quantity": 25,
        "minimum_stock": 100,
        "price": 75000.0,
        "stock_status": "low",
        "days_until_stockout": 7,
        "is_active": true
      }
    ],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total_items": 15,
      "total_pages": 2
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/medicines/low-stock?threshold=100" \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN"
```

**Use Case:**

- Inventory management
- Purchase planning
- Stock alert notifications

---

### 4. List Out of Stock Medicines

**Endpoint:** `GET /api/v1/medicines/out-of-stock`

**Description:** Mendapatkan daftar medicines yang habis (stock = 0).

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**
Same as List Low Stock.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Out of stock medicines retrieved successfully",
  "data": {
    "total_out_of_stock": 8,
    "data": [
      {
        "id": 20,
        "name": "Ceftriaxone 1g Injection",
        "generic_name": "Ceftriaxone",
        "type": "injection",
        "stock_quantity": 0,
        "price": 45000.0,
        "last_stock_date": "2024-01-10T10:00:00Z",
        "days_out_of_stock": 9,
        "is_active": true
      }
    ],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total_items": 8,
      "total_pages": 1
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/medicines/out-of-stock" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 5. List Inactive Medicines

**Endpoint:** `GET /api/v1/medicines/inactive`

**Description:** Mendapatkan daftar medicines yang tidak aktif (discontinued, expired, dll).

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**
Same as List Medicines.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Inactive medicines retrieved successfully",
  "data": {
    "data": [
      {
        "id": 100,
        "name": "Old Medicine Brand",
        "generic_name": "Generic Name",
        "type": "tablet",
        "stock_quantity": 50,
        "is_active": false,
        "deactivated_reason": "Discontinued by manufacturer",
        "created_at": "2020-01-01T10:00:00Z",
        "updated_at": "2024-01-19T10:00:00Z"
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
curl -X GET "http://localhost:8080/api/v1/medicines/inactive" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 6. List Deleted Medicines

**Endpoint:** `GET /api/v1/medicines/deleted`

**Description:** Mendapatkan daftar medicines yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**
Same as List Medicines.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Deleted medicines retrieved successfully",
  "data": {
    "data": [
      {
        "id": 150,
        "name": "Deleted Medicine",
        "generic_name": "Generic",
        "type": "tablet",
        "stock_quantity": 0,
        "is_active": false,
        "created_at": "2023-01-01T10:00:00Z",
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
curl -X GET "http://localhost:8080/api/v1/medicines/deleted" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 7. Get Medicine by ID

**Endpoint:** `GET /api/v1/medicines/:id`

**Description:** Mendapatkan detail medicine berdasarkan ID, termasuk usage statistics.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Medicine ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine retrieved successfully",
  "data": {
    "id": 1,
    "name": "Paracetamol 500mg",
    "generic_name": "Paracetamol",
    "brand_name": "Sanmol",
    "type": "tablet",
    "strength": "500mg",
    "manufacturer": "PT. Sanbe Farma",
    "unit": "tablet",
    "stock_quantity": 1000,
    "price": 500.0,
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z",
    "stock_info": {
      "is_low_stock": false,
      "is_out_of_stock": false,
      "minimum_stock_level": 100,
      "reorder_level": 200,
      "last_restock_date": "2024-01-15T10:00:00Z",
      "average_daily_usage": 50
    },
    "statistics": {
      "total_prescribed": 5000,
      "prescribed_this_month": 500,
      "total_dispensed": 4800,
      "dispensed_this_month": 480,
      "revenue_this_month": 240000.0
    },
    "details": {
      "indication": "Analgesic and antipyretic",
      "dosage_form": "Tablet",
      "route": "Oral",
      "storage": "Store at room temperature",
      "expiry_tracking": true
    }
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Medicine not found",
  "error": "medicine not found"
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/medicines/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 8. Get Medicine by Name

**Endpoint:** `GET /api/v1/medicines/name/:name`

**Description:** Mendapatkan medicine berdasarkan name (untuk quick search).

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `name`: Medicine Name (string), e.g., "Paracetamol"

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine retrieved successfully",
  "data": {
    "id": 1,
    "name": "Paracetamol 500mg",
    "generic_name": "Paracetamol",
    "brand_name": "Sanmol",
    "type": "tablet",
    "strength": "500mg",
    "stock_quantity": 1000,
    "price": 500.0,
    "is_active": true,
    "is_available": true
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/medicines/name/Paracetamol" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 9. Get Medicines by Type

**Endpoint:** `GET /api/v1/medicines/type/:type`

**Description:** Mendapatkan daftar medicines berdasarkan tipe.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `type`: Medicine Type (string), values: tablet, capsule, syrup, injection, ointment, other

**Query Parameters:**

| Parameter   | Type    | Description                        |
| ----------- | ------- | ---------------------------------- |
| `is_active` | boolean | Filter active only (default: true) |
| `has_stock` | boolean | Filter with stock only             |
| `page`      | integer | Page number                        |
| `page_size` | integer | Items per page                     |

**Example Request:**

```
GET /api/v1/medicines/type/tablet?is_active=true&has_stock=true
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicines by type retrieved successfully",
  "data": {
    "type": "tablet",
    "total_medicines": 60,
    "available_medicines": 55,
    "total_stock": 15000,
    "data": [
      {
        "id": 1,
        "name": "Paracetamol 500mg",
        "generic_name": "Paracetamol",
        "strength": "500mg",
        "stock_quantity": 1000,
        "price": 500.0,
        "is_active": true
      },
      {
        "id": 2,
        "name": "Ibuprofen 400mg",
        "generic_name": "Ibuprofen",
        "strength": "400mg",
        "stock_quantity": 800,
        "price": 1500.0,
        "is_active": true
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
curl -X GET "http://localhost:8080/api/v1/medicines/type/tablet?has_stock=true" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 10. Advanced Search Medicines

**Endpoint:** `GET /api/v1/medicines/search`

**Description:** Pencarian medicines dengan multiple criteria.

**Authentication:** Required (All Authenticated Users)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter      | Type    | Description                              |
| -------------- | ------- | ---------------------------------------- |
| `keyword`      | string  | Search in name, generic_name, brand_name |
| `type`         | string  | Filter by type                           |
| `manufacturer` | string  | Filter by manufacturer                   |
| `min_price`    | decimal | Minimum price                            |
| `max_price`    | decimal | Maximum price                            |
| `has_stock`    | boolean | Filter with stock                        |
| `is_active`    | boolean | Filter by status                         |
| `page`         | integer | Page number                              |
| `page_size`    | integer | Items per page                           |

**Example Request:**

```
GET /api/v1/medicines/search?keyword=para&type=tablet&max_price=1000&has_stock=true
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Search results retrieved successfully",
  "data": {
    "search_criteria": {
      "keyword": "para",
      "type": "tablet",
      "max_price": 1000,
      "has_stock": true
    },
    "results_found": 3,
    "data": [
      {
        "id": 1,
        "name": "Paracetamol 500mg",
        "generic_name": "Paracetamol",
        "brand_name": "Sanmol",
        "type": "tablet",
        "strength": "500mg",
        "stock_quantity": 1000,
        "price": 500.0,
        "is_active": true,
        "relevance_score": 0.98
      },
      {
        "id": 5,
        "name": "Paracetamol 650mg",
        "generic_name": "Paracetamol",
        "brand_name": "Panadol",
        "type": "tablet",
        "strength": "650mg",
        "stock_quantity": 500,
        "price": 800.0,
        "is_active": true,
        "relevance_score": 0.95
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
curl -X GET "http://localhost:8080/api/v1/medicines/search?keyword=para&type=tablet" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Pharmacist/Stock Endpoints

### 11. Add Stock

**Endpoint:** `PATCH /api/v1/medicines/:id/add-stock`

**Description:** Menambah stock medicine (saat pembelian/restock).

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Medicine ID (integer)

**Request Body:**

```json
{
  "quantity": 500,
  "notes": "Pembelian dari distributor XYZ"
}
```

**Field Rules:**

- `quantity`: required, min 1
- `notes`: optional, text

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Stock added successfully",
  "data": {
    "id": 1,
    "name": "Paracetamol 500mg",
    "previous_stock": 1000,
    "added_quantity": 500,
    "current_stock": 1500,
    "stock_value": 750000.0,
    "updated_at": "2024-01-19T15:00:00Z"
  }
}
```

**Response Error (400 Bad Request):**

```json
{
  "success": false,
  "message": "Failed to add stock",
  "error": "invalid quantity"
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/medicines/1/add-stock \
  -H "Authorization: Bearer RECEPTIONIST_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 500,
    "notes": "Restock from supplier"
  }'
```

**Notes:**

- Automatically updates stock value
- Creates stock movement history
- Sends notification if reaches normal level

---

### 12. Reduce Stock

**Endpoint:** `PATCH /api/v1/medicines/:id/reduce-stock`

**Description:** Mengurangi stock medicine (saat dispensing/penjualan).

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Medicine ID (integer)

**Request Body:**

```json
{
  "quantity": 100,
  "notes": "Dispensed for prescription #12345"
}
```

**Field Rules:**

- `quantity`: required, min 1
- `notes`: optional, text

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Stock reduced successfully",
  "data": {
    "id": 1,
    "name": "Paracetamol 500mg",
    "previous_stock": 1500,
    "reduced_quantity": 100,
    "current_stock": 1400,
    "stock_value": 700000.0,
    "is_low_stock": false,
    "updated_at": "2024-01-19T16:00:00Z"
  }
}
```

**Response Error (400 Bad Request):**

```json
{
  "success": false,
  "message": "Failed to reduce stock",
  "error": "insufficient stock quantity"
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/medicines/1/reduce-stock \
  -H "Authorization: Bearer RECEPTIONIST_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 100,
    "notes": "Prescription dispensing"
  }'
```

**Notes:**

- Validates sufficient stock before reducing
- Automatically called saat prescription dispensing
- Sends low stock alert if below threshold
- Creates stock movement history

---

## Admin Endpoints

### 13. Create Medicine

**Endpoint:** `POST /api/v1/medicines`

**Description:** Admin membuat medicine record baru.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "name": "Paracetamol 500mg",
  "generic_name": "Paracetamol",
  "brand_name": "Sanmol",
  "type": "tablet",
  "strength": "500mg",
  "manufacturer": "PT. Sanbe Farma",
  "unit": "tablet",
  "stock_quantity": 1000,
  "price": 500.0,
  "is_active": true
}
```

**Field Rules:**

- `name`: required, max 200 characters, indexed
- `generic_name`: optional, max 200 characters
- `brand_name`: optional, max 200 characters
- `type`: required, enum (tablet, capsule, syrup, injection, ointment, other)
- `strength`: optional, max 50 (e.g., "500mg", "10ml")
- `manufacturer`: optional, max 100
- `unit`: optional, max 20 (tablet, capsule, ml, mg, etc)
- `stock_quantity`: required, default 0
- `price`: optional, decimal(10,2)
- `is_active`: optional, boolean (default: true)

**Response Success (201 Created):**

```json
{
  "success": true,
  "message": "Medicine created successfully",
  "data": {
    "id": 1,
    "name": "Paracetamol 500mg",
    "generic_name": "Paracetamol",
    "brand_name": "Sanmol",
    "type": "tablet",
    "strength": "500mg",
    "manufacturer": "PT. Sanbe Farma",
    "unit": "tablet",
    "stock_quantity": 1000,
    "price": 500.0,
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
  "message": "Failed to create medicine",
  "error": "medicine with same name and strength already exists"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/medicines \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Paracetamol 500mg",
    "generic_name": "Paracetamol",
    "type": "tablet",
    "strength": "500mg",
    "stock_quantity": 1000,
    "price": 500
  }'
```

**Notes:**

- Name + Strength combination should be unique
- Type validation using enum
- Stock quantity defaults to 0 if not provided
- Price per unit dalam IDR

---

### 14. Update Medicine

**Endpoint:** `PUT /api/v1/medicines/:id`

**Description:** Admin mengupdate data medicine berdasarkan ID.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Medicine ID (integer)

**Request Body:**

```json
{
  "name": "Paracetamol 500mg Updated",
  "generic_name": "Paracetamol",
  "brand_name": "Sanmol Plus",
  "type": "tablet",
  "strength": "500mg",
  "manufacturer": "PT. Sanbe Farma",
  "price": 550.0,
  "is_active": true
}
```

**Field Rules:**

- All fields optional
- Validation same as create
- Cannot update stock_quantity via this endpoint (use add/reduce stock)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine updated successfully",
  "data": {
    "id": 1,
    "name": "Paracetamol 500mg Updated",
    "generic_name": "Paracetamol",
    "brand_name": "Sanmol Plus",
    "type": "tablet",
    "strength": "500mg",
    "manufacturer": "PT. Sanbe Farma",
    "unit": "tablet",
    "stock_quantity": 1000,
    "price": 550.0,
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T15:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/medicines/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "price": 550,
    "brand_name": "Sanmol Plus"
  }'
```

---

### 15. Activate Medicine

**Endpoint:** `PATCH /api/v1/medicines/:id/activate`

**Description:** Admin mengaktifkan medicine yang sedang inactive.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Medicine ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine activated successfully",
  "data": null
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/medicines/1/activate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

**Use Case:**

- Medicine kembali tersedia dari manufacturer
- Medicine approved untuk digunakan kembali
- Regulatory clearance obtained

---

### 16. Deactivate Medicine

**Endpoint:** `PATCH /api/v1/medicines/:id/deactivate`

**Description:** Admin menonaktifkan medicine (recalled, expired batch, discontinued).

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Medicine ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine deactivated successfully",
  "data": null
}
```

**Notes:**

- Medicine yang deactivated tidak bisa di-prescribe
- Prescription items yang sudah ada tetap valid
- Stock tetap ada tapi tidak available untuk dispensing

**⚠️ Business Rules:**

- Pending prescriptions dengan medicine ini akan di-alert
- Alternative medicine suggestions provided
- Historical data tetap accessible

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/medicines/1/deactivate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 17. Soft Delete Medicine

**Endpoint:** `DELETE /api/v1/medicines/:id`

**Description:** Admin menghapus medicine (soft delete). Data tetap ada dan bisa di-restore.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Medicine ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine deleted successfully",
  "data": null
}
```

**Notes:**

- Medicine yang dihapus tidak muncul di list normal
- Prescription records tetap tersimpan
- Stock movements tetap accessible
- Bisa di-restore dengan endpoint restore

**⚠️ Business Rules:**

- Cannot delete if medicine has active prescriptions
- Stock should be 0 or transferred first
- Historical data preserved

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/medicines/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 18. Restore Medicine

**Endpoint:** `PATCH /api/v1/medicines/:id/restore`

**Description:** Admin me-restore medicine yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Medicine ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine restored successfully",
  "data": null
}
```

**Notes:**

- Medicine di-restore dengan status inactive
- Perlu activate manual jika ingin langsung available
- Stock quantity preserved

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/medicines/1/restore \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

## Super Admin Endpoints

### 19. Hard Delete Medicine

**Endpoint:** `DELETE /api/v1/medicines/:id/hard-delete`

**Description:** Super Admin menghapus medicine secara permanen dari database.

**Authentication:** Required (Super Admin Only)

**Request Headers:**

```
Authorization: Bearer <super-admin-token>
```

**URL Parameters:**

- `id`: Medicine ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medicine permanently deleted",
  "data": null
}
```

**⚠️ WARNING:**

- Data dihapus permanen dari database
- Tidak bisa di-restore
- Gunakan dengan sangat hati-hati

**⚠️ Business Rules:**

- Cannot hard delete if medicine has prescription records
- Cannot hard delete if medicine has stock movements
- Must archive all historical data first
- Requires special approval

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/medicines/1/hard-delete \
  -H "Authorization: Bearer SUPER_ADMIN_JWT_TOKEN"
```

---

## Error Responses

### 400 Bad Request

```json
{
  "success": false,
  "message": "Validation error",
  "error": "insufficient stock quantity"
}
```

### 404 Not Found

```json
{
  "success": false,
  "message": "Medicine not found",
  "error": "medicine not found"
}
```

### 409 Conflict

```json
{
  "success": false,
  "message": "Cannot delete medicine",
  "error": "medicine has associated prescription records"
}
```

### 422 Unprocessable Entity

```json
{
  "success": false,
  "message": "Cannot reduce stock",
  "error": "requested quantity exceeds available stock"
}
```

---

## Data Models

### Medicine Object (Full)

```json
{
  "id": 1,
  "name": "Paracetamol 500mg",
  "generic_name": "Paracetamol",
  "brand_name": "Sanmol",
  "type": "tablet",
  "strength": "500mg",
  "manufacturer": "PT. Sanbe Farma",
  "unit": "tablet",
  "stock_quantity": 1000,
  "price": 500.0,
  "is_active": true,
  "created_at": "2024-01-19T10:00:00Z",
  "updated_at": "2024-01-19T10:00:00Z",
  "deleted_at": null
}
```

### Medicine Types

```json
{
  "types": [
    {
      "value": "tablet",
      "label": "Tablet",
      "description": "Solid oral dosage form"
    },
    {
      "value": "capsule",
      "label": "Capsule",
      "description": "Gelatin shell containing medicine"
    },
    {
      "value": "syrup",
      "label": "Syrup",
      "description": "Liquid oral dosage form"
    },
    {
      "value": "injection",
      "label": "Injection",
      "description": "Parenteral dosage form"
    },
    {
      "value": "ointment",
      "label": "Ointment/Cream",
      "description": "Topical dosage form"
    },
    {
      "value": "other",
      "label": "Other",
      "description": "Inhaler, suppository, etc"
    }
  ]
}
```

---

## Business Rules

1. **Stock Management**: Real-time stock tracking
2. **Low Stock Alert**: Automatic notification at threshold
3. **Price Management**: Price per unit dalam IDR
4. **Type Validation**: Must be one of defined types
5. **Name Uniqueness**: Name + Strength combination unique
6. **Active Status**: Only active medicines can be prescribed
7. **Stock Movement**: All stock changes logged
8. **Expiry Tracking**: Support for batch and expiry management
9. **BPJS Integration**: Support for formularium nasional
10. **Generics Priority**: Generic names for better tracking

---

## Common Medicines (Indonesia)

### Analgesics & Antipyretics

```json
[
  { "name": "Paracetamol 500mg", "type": "tablet", "price": 500 },
  { "name": "Ibuprofen 400mg", "type": "tablet", "price": 1500 },
  { "name": "Asam Mefenamat 500mg", "type": "capsule", "price": 2000 },
  { "name": "Paracetamol Syrup 120mg/5ml", "type": "syrup", "price": 15000 }
]
```

### Antibiotics

```json
[
  { "name": "Amoxicillin 500mg", "type": "capsule", "price": 2000 },
  { "name": "Ciprofloxacin 500mg", "type": "tablet", "price": 3500 },
  { "name": "Cefadroxil 500mg", "type": "capsule", "price": 5000 },
  { "name": "Azithromycin 500mg", "type": "tablet", "price": 8000 }
]
```

### Antidiabetics

```json
[
  { "name": "Metformin 500mg", "type": "tablet", "price": 1000 },
  { "name": "Glimepiride 2mg", "type": "tablet", "price": 3000 },
  { "name": "Insulin Glargine 100IU/ml", "type": "injection", "price": 350000 }
]
```

### Antihypertensives

```json
[
  { "name": "Amlodipine 10mg", "type": "tablet", "price": 2500 },
  { "name": "Captopril 25mg", "type": "tablet", "price": 1500 },
  { "name": "Valsartan 80mg", "type": "tablet", "price": 5000 }
]
```

---

## Common Use Cases

### Use Case 1: Doctor Prescribes Medicine

```bash
# 1. Search available medicines
GET /api/v1/medicines/search?keyword=para&has_stock=true

# 2. Check medicine details and stock
GET /api/v1/medicines/1

# 3. Create prescription (different endpoint)
POST /api/v1/prescriptions
```

### Use Case 2: Pharmacist Dispenses Prescription

```bash
# 1. View prescription items
GET /api/v1/prescriptions/123/items

# 2. Check stock availability
GET /api/v1/medicines/1

# 3. Reduce stock (automatic)
PATCH /api/v1/medicines/1/reduce-stock

# 4. Dispense prescription
PATCH /api/v1/prescriptions/123/dispense
```

### Use Case 3: Stock Management

```bash
# 1. Check low stock
GET /api/v1/medicines/low-stock?threshold=100

# 2. Purchase order created
# 3. Receive stock
PATCH /api/v1/medicines/1/add-stock

# 4. View stock movement history
GET /api/v1/medicines/1/stock-movements
```

---

## Testing Examples

### Test 1: Complete Medicine Management Flow

```bash
# 1. Create Medicine
curl -X POST http://localhost:8080/api/v1/medicines \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Paracetamol 500mg",
    "generic_name": "Paracetamol",
    "type": "tablet",
    "stock_quantity": 1000,
    "price": 500
  }'

# 2. Add Stock
curl -X PATCH http://localhost:8080/api/v1/medicines/1/add-stock \
  -H "Authorization: Bearer RECEPTIONIST_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"quantity": 500}'

# 3. Reduce Stock
curl -X PATCH http://localhost:8080/api/v1/medicines/1/reduce-stock \
  -H "Authorization: Bearer RECEPTIONIST_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"quantity": 100}'

# 4. Check Low Stock
curl -X GET "http://localhost:8080/api/v1/medicines/low-stock" \
  -H "Authorization: Bearer ADMIN_TOKEN"

# 5. Deactivate Medicine
curl -X PATCH http://localhost:8080/api/v1/medicines/1/deactivate \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

---

## Notes

- Stock management real-time
- Price per unit dalam IDR
- Support batch and expiry tracking
- Generic name untuk BPJS integration
- Automatic low stock alerts
- Stock movement history tracking
- Integration with prescription system
- Support for drug interactions checking
- Formularium nasional compliance
- LASA (Look-Alike Sound-Alike) warnings

---

**Last Updated:** 2024-01-19  
**API Version:** 1.0.0
