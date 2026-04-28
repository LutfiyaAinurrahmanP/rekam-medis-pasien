# Payment Gateway API Documentation

## Overview

API untuk integrasi payment gateway pada sistem rekam medis. Pasien dapat membayar tagihan secara online melalui berbagai metode pembayaran digital (QRIS, Virtual Account, kartu kredit, e-wallet). Sistem mendukung integrasi dengan **Midtrans** dan **Xendit** sebagai payment gateway provider.

Alur pembayaran online:

1. Pasien memilih tagihan (`billing`) yang ingin dibayar
2. Sistem membuat transaksi ke payment gateway dan mengembalikan payment URL/token
3. Pasien diarahkan ke halaman pembayaran gateway
4. Gateway mengirim notifikasi pembayaran ke webhook endpoint
5. Sistem memperbarui status billing secara otomatis

**Base URL:** `/api/v1/payments`

> **Catatan:** Endpoint webhook (`POST /payments/webhook/*`) bersifat **publik** (tidak memerlukan JWT) tetapi diverifikasi menggunakan signature dari payment gateway provider.

---

## Table of Contents

- [Authentication](#authentication)
- [Authorization](#authorization)
- [Endpoints Summary](#endpoints-summary)
- [Endpoints Detail](#endpoints-detail)
- [Webhook Events](#webhook-events)
- [Database Model](#database-model)
- [Error Responses](#error-responses)
- [Payment Method Reference](#payment-method-reference)

---

## Authentication

Semua endpoint (kecuali webhook) memerlukan JWT token di header:

```
Authorization: Bearer <your-jwt-token>
```

Endpoint webhook diverifikasi dengan signature header dari provider:

- Midtrans: `X-Midtrans-Signature-SHA512`
- Xendit: `x-callback-token`

---

## Authorization

| Endpoint                          | Patient  | Doctor | Receptionist | Admin | Super Admin |
| --------------------------------- | -------- | ------ | ------------ | ----- | ----------- |
| POST /payments/initiate           | ✅ (Own) | ❌     | ✅           | ✅    | ✅          |
| GET /payments/:payment_id         | ✅ (Own) | ❌     | ✅           | ✅    | ✅          |
| GET /payments/:payment_id/status  | ✅ (Own) | ❌     | ✅           | ✅    | ✅          |
| GET /payments/billing/:billing_id | ✅ (Own) | ❌     | ✅           | ✅    | ✅          |
| GET /payments/my-payments         | ✅       | ❌     | ❌           | ❌    | ❌          |
| GET /payments                     | ❌       | ❌     | ✅           | ✅    | ✅          |
| POST /payments/:payment_id/refund | ❌       | ❌     | ❌           | ✅    | ✅          |
| POST /payments/:payment_id/cancel | ✅ (Own) | ❌     | ✅           | ✅    | ✅          |
| POST /payments/webhook/midtrans   | Public\* | -      | -            | -     | -           |
| POST /payments/webhook/xendit     | Public\* | -      | -            | -     | -           |

\*Terverifikasi via signature header dari provider

---

## Endpoints Summary

| Method | Endpoint                        | Description                                | Role Required                             |
| ------ | ------------------------------- | ------------------------------------------ | ----------------------------------------- |
| POST   | `/payments/initiate`            | Inisiasi transaksi pembayaran baru         | Patient, Receptionist, Admin, Super Admin |
| GET    | `/payments`                     | List semua transaksi pembayaran            | Receptionist, Admin, Super Admin          |
| GET    | `/payments/my-payments`         | List transaksi pembayaran milik pasien     | Patient                                   |
| GET    | `/payments/:payment_id`         | Detail transaksi pembayaran                | Patient (Own), Receptionist, Admin, SA    |
| GET    | `/payments/:payment_id/status`  | Cek status terkini ke payment gateway      | Patient (Own), Receptionist, Admin, SA    |
| GET    | `/payments/billing/:billing_id` | List transaksi berdasarkan billing         | Patient (Own), Receptionist, Admin, SA    |
| POST   | `/payments/:payment_id/cancel`  | Batalkan transaksi yang belum dibayar      | Patient (Own), Receptionist, Admin, SA    |
| POST   | `/payments/:payment_id/refund`  | Proses refund transaksi yang sudah dibayar | Admin, Super Admin                        |
| POST   | `/payments/webhook/midtrans`    | Webhook callback dari Midtrans             | Public (verified by signature)            |
| POST   | `/payments/webhook/xendit`      | Webhook callback dari Xendit               | Public (verified by signature)            |

---

## Endpoints Detail

### 1. Initiate Payment

**Endpoint:** `POST /api/v1/payments/initiate`

**Description:** Membuat transaksi pembayaran baru ke payment gateway berdasarkan billing ID. Sistem akan mengembalikan payment URL atau token yang dapat digunakan pasien untuk menyelesaikan pembayaran.

**Authentication:** Required (Patient, Receptionist, Admin, Super Admin)

**Request Headers:**

| Header        | Type   | Required | Description        |
| ------------- | ------ | -------- | ------------------ |
| Authorization | string | Yes      | `Bearer <token>`   |
| Content-Type  | string | Yes      | `application/json` |

**Request Body:**

| Field           | Type    | Required | Description                                                                                                                 |
| --------------- | ------- | -------- | --------------------------------------------------------------------------------------------------------------------------- |
| billing_id      | integer | Yes      | ID billing yang akan dibayar                                                                                                |
| payment_channel | string  | Yes      | Channel pembayaran. Lihat [Payment Method Reference](#payment-method-reference)                                             |
| provider        | string  | Yes      | Payment gateway provider: `midtrans` atau `xendit`                                                                          |
| bank            | string  | No       | Kode bank untuk Virtual Account (wajib jika `payment_channel = virtual_account`): `bca`, `bni`, `bri`, `mandiri`, `permata` |
| ewallet_type    | string  | No       | Jenis e-wallet (wajib jika `payment_channel = ewallet`): `gopay`, `ovo`, `dana`, `shopeepay`                                |
| callback_url    | string  | No       | URL redirect setelah pembayaran selesai di sisi frontend/mobile                                                             |

**Example Request Body (QRIS):**

```json
{
  "billing_id": 88,
  "payment_channel": "qris",
  "provider": "midtrans",
  "callback_url": "https://myapp.com/payment/result"
}
```

**Example Request Body (Virtual Account BCA):**

```json
{
  "billing_id": 88,
  "payment_channel": "virtual_account",
  "provider": "xendit",
  "bank": "bca",
  "callback_url": "https://myapp.com/payment/result"
}
```

**Example Request Body (GoPay):**

```json
{
  "billing_id": 88,
  "payment_channel": "ewallet",
  "provider": "midtrans",
  "ewallet_type": "gopay",
  "callback_url": "https://myapp.com/payment/result"
}
```

**Response Success (201):**

```json
{
  "status": "success",
  "message": "Payment initiated successfully",
  "data": {
    "payment_id": "PAY-2024-000088",
    "billing_id": 88,
    "invoice_number": "INV-2024-000088",
    "amount": 350000,
    "payment_channel": "qris",
    "provider": "midtrans",
    "payment_status": "pending",
    "expired_at": "2024-01-15T11:00:00Z",
    "payment_url": "https://app.midtrans.com/snap/v2/vtweb/abc123xyz",
    "token": "abc123xyz",
    "qr_code_url": "https://api.midtrans.com/v2/qris/abc123/qr-code",
    "virtual_account_number": null,
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

**Response Success — Virtual Account (201):**

```json
{
  "status": "success",
  "message": "Payment initiated successfully",
  "data": {
    "payment_id": "PAY-2024-000089",
    "billing_id": 89,
    "invoice_number": "INV-2024-000089",
    "amount": 250000,
    "payment_channel": "virtual_account",
    "provider": "xendit",
    "bank": "bca",
    "payment_status": "pending",
    "expired_at": "2024-01-16T10:00:00Z",
    "payment_url": null,
    "token": null,
    "qr_code_url": null,
    "virtual_account_number": "8277082000088123",
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/payments/initiate" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "billing_id": 88,
    "payment_channel": "qris",
    "provider": "midtrans",
    "callback_url": "https://myapp.com/payment/result"
  }'
```

---

### 2. List All Payments

**Endpoint:** `GET /api/v1/payments`

**Description:** Mendapatkan daftar semua transaksi pembayaran dengan pagination dan filter.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Query Parameters:**

| Parameter       | Type    | Required | Default    | Description                                                                    |
| --------------- | ------- | -------- | ---------- | ------------------------------------------------------------------------------ |
| page            | integer | No       | 1          | Halaman saat ini                                                               |
| page_size       | integer | No       | 10         | Jumlah data per halaman (max: 100)                                             |
| sort_by         | string  | No       | created_at | Field sorting: `created_at`, `amount`, `payment_status`                        |
| sort_dir        | string  | No       | desc       | Arah sorting: `asc` atau `desc`                                                |
| payment_status  | string  | No       | -          | Filter status: `pending`, `paid`, `failed`, `expired`, `cancelled`, `refunded` |
| payment_channel | string  | No       | -          | Filter channel: `qris`, `virtual_account`, `ewallet`, `credit_card`            |
| provider        | string  | No       | -          | Filter provider: `midtrans`, `xendit`                                          |
| start_date      | string  | No       | -          | Filter tanggal mulai (`YYYY-MM-DD`)                                            |
| end_date        | string  | No       | -          | Filter tanggal akhir (`YYYY-MM-DD`)                                            |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Payments retrieved successfully",
  "data": [
    {
      "payment_id": "PAY-2024-000088",
      "billing_id": 88,
      "invoice_number": "INV-2024-000088",
      "patient_name": "Budi Santoso",
      "amount": 350000,
      "payment_channel": "qris",
      "provider": "midtrans",
      "payment_status": "paid",
      "paid_at": "2024-01-15T10:45:00Z",
      "created_at": "2024-01-15T10:00:00Z"
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
curl -X GET "http://localhost:8080/api/v1/payments?payment_status=pending&payment_channel=virtual_account" \
  -H "Authorization: Bearer <token>"
```

---

### 3. Get My Payments (Patient)

**Endpoint:** `GET /api/v1/payments/my-payments`

**Description:** Mendapatkan riwayat semua transaksi pembayaran milik pasien yang sedang login.

**Authentication:** Required (Patient Only)

**Query Parameters:**

| Parameter      | Type   | Required | Default | Description                                                                    |
| -------------- | ------ | -------- | ------- | ------------------------------------------------------------------------------ |
| payment_status | string | No       | -       | Filter status: `pending`, `paid`, `failed`, `expired`, `cancelled`, `refunded` |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "My payments retrieved successfully",
  "data": [
    {
      "payment_id": "PAY-2024-000088",
      "invoice_number": "INV-2024-000088",
      "amount": 350000,
      "payment_channel": "qris",
      "payment_status": "paid",
      "paid_at": "2024-01-15T10:45:00Z"
    },
    {
      "payment_id": "PAY-2024-000092",
      "invoice_number": "INV-2024-000092",
      "amount": 125000,
      "payment_channel": "virtual_account",
      "payment_status": "pending",
      "expired_at": "2024-01-20T10:00:00Z"
    }
  ]
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/payments/my-payments?payment_status=pending" \
  -H "Authorization: Bearer <token>"
```

---

### 4. Get Payment Detail

**Endpoint:** `GET /api/v1/payments/:payment_id`

**Description:** Mendapatkan detail lengkap sebuah transaksi pembayaran.

**Authentication:** Required (Patient hanya bisa melihat miliknya sendiri)

**Path Parameters:**

| Parameter  | Type   | Required | Description                            |
| ---------- | ------ | -------- | -------------------------------------- |
| payment_id | string | Yes      | Payment ID (format: `PAY-YYYY-XXXXXX`) |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Payment retrieved successfully",
  "data": {
    "payment_id": "PAY-2024-000088",
    "billing_id": 88,
    "invoice_number": "INV-2024-000088",
    "patient": {
      "id": 10,
      "name": "Budi Santoso"
    },
    "amount": 350000,
    "payment_channel": "qris",
    "provider": "midtrans",
    "provider_transaction_id": "midtrans-order-88",
    "payment_status": "paid",
    "expired_at": "2024-01-15T11:00:00Z",
    "paid_at": "2024-01-15T10:45:00Z",
    "payment_url": "https://app.midtrans.com/snap/v2/vtweb/abc123xyz",
    "qr_code_url": "https://api.midtrans.com/v2/qris/abc123/qr-code",
    "virtual_account_number": null,
    "bank": null,
    "ewallet_type": null,
    "refund_id": null,
    "refund_amount": null,
    "refunded_at": null,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:45:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/payments/PAY-2024-000088" \
  -H "Authorization: Bearer <token>"
```

---

### 5. Check Payment Status

**Endpoint:** `GET /api/v1/payments/:payment_id/status`

**Description:** Mengecek status transaksi terkini langsung ke payment gateway (real-time inquiry). Gunakan endpoint ini jika status di sistem terlihat `pending` padahal pembayaran sudah dilakukan.

**Authentication:** Required (Patient hanya bisa mengecek miliknya sendiri)

**Path Parameters:**

| Parameter  | Type   | Required | Description  |
| ---------- | ------ | -------- | ------------ |
| payment_id | string | Yes      | ID transaksi |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Payment status retrieved successfully",
  "data": {
    "payment_id": "PAY-2024-000088",
    "system_status": "paid",
    "gateway_status": "settlement",
    "is_synced": true,
    "last_checked_at": "2024-01-15T10:50:00Z",
    "paid_at": "2024-01-15T10:45:00Z"
  }
}
```

**Response — Status Belum Sinkron (200):**

```json
{
  "status": "success",
  "message": "Payment status retrieved successfully",
  "data": {
    "payment_id": "PAY-2024-000089",
    "system_status": "pending",
    "gateway_status": "settlement",
    "is_synced": false,
    "note": "Payment completed at gateway but not yet reflected in system. Please wait or contact support.",
    "last_checked_at": "2024-01-15T11:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/payments/PAY-2024-000088/status" \
  -H "Authorization: Bearer <token>"
```

---

### 6. Get Payments by Billing

**Endpoint:** `GET /api/v1/payments/billing/:billing_id`

**Description:** Mendapatkan semua transaksi pembayaran berdasarkan billing ID. Satu billing bisa memiliki lebih dari satu transaksi (jika transaksi pertama expired lalu pasien membuat transaksi baru).

**Authentication:** Required (Patient hanya bisa melihat miliknya sendiri)

**Path Parameters:**

| Parameter  | Type    | Required | Description |
| ---------- | ------- | -------- | ----------- |
| billing_id | integer | Yes      | ID billing  |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Billing payments retrieved successfully",
  "data": [
    {
      "payment_id": "PAY-2024-000091",
      "amount": 350000,
      "payment_channel": "qris",
      "payment_status": "expired",
      "expired_at": "2024-01-15T10:30:00Z"
    },
    {
      "payment_id": "PAY-2024-000088",
      "amount": 350000,
      "payment_channel": "virtual_account",
      "bank": "bca",
      "payment_status": "paid",
      "paid_at": "2024-01-15T10:45:00Z"
    }
  ]
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/payments/billing/88" \
  -H "Authorization: Bearer <token>"
```

---

### 7. Cancel Payment

**Endpoint:** `POST /api/v1/payments/:payment_id/cancel`

**Description:** Membatalkan transaksi pembayaran yang masih berstatus `pending`. Jika transaksi sudah diteruskan ke gateway, pembatalan juga akan dikirimkan ke gateway provider.

**Authentication:** Required (Patient hanya bisa membatalkan miliknya sendiri)

**Path Parameters:**

| Parameter  | Type   | Required | Description  |
| ---------- | ------ | -------- | ------------ |
| payment_id | string | Yes      | ID transaksi |

**Request Body:**

| Field  | Type   | Required | Description       |
| ------ | ------ | -------- | ----------------- |
| reason | string | No       | Alasan pembatalan |

**Example Request Body:**

```json
{
  "reason": "Ingin mengganti metode pembayaran ke transfer bank"
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Payment cancelled successfully",
  "data": {
    "payment_id": "PAY-2024-000090",
    "payment_status": "cancelled",
    "cancelled_at": "2024-01-15T10:10:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/payments/PAY-2024-000090/cancel" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"reason": "Ingin mengganti metode pembayaran"}'
```

---

### 8. Refund Payment

**Endpoint:** `POST /api/v1/payments/:payment_id/refund`

**Description:** Memproses pengembalian dana (refund) ke pasien untuk transaksi yang sudah berstatus `paid`. Refund dikirim ke payment gateway provider dan dana dikembalikan sesuai kebijakan provider (biasanya 1–7 hari kerja).

**Authentication:** Required (Admin, Super Admin Only)

**Path Parameters:**

| Parameter  | Type   | Required | Description  |
| ---------- | ------ | -------- | ------------ |
| payment_id | string | Yes      | ID transaksi |

**Request Body:**

| Field         | Type    | Required | Description                                                               |
| ------------- | ------- | -------- | ------------------------------------------------------------------------- |
| refund_amount | integer | Yes      | Jumlah refund dalam Rupiah. Tidak boleh melebihi `amount` transaksi asal. |
| reason        | string  | Yes      | Alasan refund (mis. pembatalan layanan, kelebihan pembayaran)             |

**Example Request Body:**

```json
{
  "refund_amount": 350000,
  "reason": "Layanan dibatalkan oleh fasilitas kesehatan karena dokter tidak tersedia"
}
```

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Refund processed successfully",
  "data": {
    "payment_id": "PAY-2024-000088",
    "refund_id": "REFUND-2024-000001",
    "original_amount": 350000,
    "refund_amount": 350000,
    "payment_status": "refunded",
    "refunded_at": "2024-01-15T14:00:00Z",
    "note": "Refund akan diterima dalam 1-7 hari kerja tergantung provider dan bank pasien"
  }
}
```

**cURL Example:**

```bash
curl -X POST "http://localhost:8080/api/v1/payments/PAY-2024-000088/refund" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "refund_amount": 350000,
    "reason": "Layanan dibatalkan oleh fasilitas kesehatan"
  }'
```

---

## Webhook Events

Endpoint webhook menerima notifikasi otomatis dari payment gateway ketika status transaksi berubah. **Sistem akan memperbarui status billing dan payment secara otomatis** berdasarkan notifikasi ini.

### Midtrans Webhook

**Endpoint:** `POST /api/v1/payments/webhook/midtrans`

**Authentication:** Tidak memerlukan JWT. Diverifikasi menggunakan `X-Midtrans-Signature-SHA512` header.

**Verification Logic:**

```
SHA512(order_id + status_code + gross_amount + server_key) == X-Midtrans-Signature-SHA512
```

**Request Body (dikirim oleh Midtrans):**

```json
{
  "transaction_time": "2024-01-15 10:45:00",
  "transaction_status": "settlement",
  "transaction_id": "midtrans-trx-xyz",
  "status_message": "midtrans payment notification",
  "status_code": "200",
  "signature_key": "abc123...",
  "payment_type": "qris",
  "order_id": "PAY-2024-000088",
  "merchant_id": "G12345",
  "gross_amount": "350000.00",
  "fraud_status": "accept",
  "currency": "IDR"
}
```

**Midtrans Transaction Status Mapping:**

| Midtrans `transaction_status`      | Sistem `payment_status` |
| ---------------------------------- | ----------------------- |
| `capture` + `fraud_status: accept` | `paid`                  |
| `settlement`                       | `paid`                  |
| `pending`                          | `pending`               |
| `deny`                             | `failed`                |
| `expire`                           | `expired`               |
| `cancel`                           | `cancelled`             |
| `refund`                           | `refunded`              |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Webhook processed"
}
```

**cURL Simulation:**

```bash
curl -X POST "http://localhost:8080/api/v1/payments/webhook/midtrans" \
  -H "Content-Type: application/json" \
  -H "X-Midtrans-Signature-SHA512: <signature>" \
  -d '{
    "transaction_status": "settlement",
    "order_id": "PAY-2024-000088",
    "gross_amount": "350000.00",
    "status_code": "200"
  }'
```

---

### Xendit Webhook

**Endpoint:** `POST /api/v1/payments/webhook/xendit`

**Authentication:** Tidak memerlukan JWT. Diverifikasi menggunakan `x-callback-token` header (Webhook Verification Token dari dashboard Xendit).

**Request Body — Invoice Paid (dikirim oleh Xendit):**

```json
{
  "id": "xendit-invoice-xyz",
  "external_id": "PAY-2024-000089",
  "status": "PAID",
  "amount": 250000,
  "paid_amount": 250000,
  "payment_channel": "BCA",
  "payment_destination": "8277082000088123",
  "paid_at": "2024-01-15T10:45:00.000Z",
  "payment_method": "BANK_TRANSFER",
  "currency": "IDR"
}
```

**Xendit Status Mapping:**

| Xendit `status` | Sistem `payment_status` |
| --------------- | ----------------------- |
| `PAID`          | `paid`                  |
| `SETTLED`       | `paid`                  |
| `PENDING`       | `pending`               |
| `EXPIRED`       | `expired`               |
| `REFUNDED`      | `refunded`              |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Webhook processed"
}
```

**cURL Simulation:**

```bash
curl -X POST "http://localhost:8080/api/v1/payments/webhook/xendit" \
  -H "Content-Type: application/json" \
  -H "x-callback-token: <xendit-webhook-token>" \
  -d '{
    "external_id": "PAY-2024-000089",
    "status": "PAID",
    "amount": 250000
  }'
```

---

## Error Responses

### 400 Bad Request

```json
{
  "status": "error",
  "message": "Validation error",
  "errors": {
    "billing_id": "billing_id is required",
    "payment_channel": "payment_channel must be one of: qris, virtual_account, ewallet, credit_card",
    "provider": "provider must be one of: midtrans, xendit",
    "bank": "bank is required when payment_channel is virtual_account"
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
  "message": "Payment not found"
}
```

### 409 Conflict

```json
{
  "status": "error",
  "message": "Billing already has an active payment transaction. Cancel the existing transaction before creating a new one."
}
```

### 422 Unprocessable Entity

```json
{
  "status": "error",
  "message": "Cannot refund a transaction that is not in paid status"
}
```

```json
{
  "status": "error",
  "message": "Refund amount exceeds original transaction amount"
}
```

### 502 Bad Gateway

```json
{
  "status": "error",
  "message": "Payment gateway error: failed to create transaction",
  "gateway_error": "Merchant account is not active"
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

### Table: payments

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| payment_number | VARCHAR(50) | UNIQUE, NOT NULL, INDEX | Nomor transaksi pembayaran |
| billing_id | BIGINT | FOREIGN KEY (billing.id), NOT NULL, INDEX | Reference ke billing/invoice |
| patient_id | BIGINT | FOREIGN KEY (patients.id), NOT NULL, INDEX | Reference ke pasien |
| amount | DECIMAL(12,2) | NOT NULL | Jumlah pembayaran |
| payment_method | VARCHAR(50) | NOT NULL | Metode pembayaran (QRIS, VA, CC, E-WALLET, dll) |
| payment_channel | VARCHAR(50) | NOT NULL | Channel/provider (midtrans, xendit, manual, dll) |
| payment_gateway_id | VARCHAR(100) | NULLABLE | Transaction ID dari payment gateway |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'pending', INDEX | Status (pending, processing, completed, failed, refunded) |
| payment_date | DATE | NULLABLE | Tanggal pembayaran berhasil |
| payment_time | TIME | NULLABLE | Waktu pembayaran berhasil |
| receipt_url | VARCHAR(500) | NULLABLE | URL untuk receipt/invoice |
| notes | TEXT | NULLABLE | Catatan transaksi |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

### Table: payment_details (untuk detail payment gateway response)

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| payment_id | BIGINT | FOREIGN KEY (payments.id), NOT NULL | Reference ke payment |
| gateway_response | LONGTEXT | NOT NULL | Full response dari payment gateway (JSON) |
| signature_key | VARCHAR(500) | NULLABLE | Signature untuk verifikasi dari gateway |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan |

**Indexes:**
- Primary Keys: id (each table)
- Foreign Keys: billing_id, patient_id (payments), payment_id (details)
- Regular Index: payment_method, status, payment_date, deleted_at

**Relationships:**
- Payment Belongs To Billing, Patient
- Payment Has One Payment Details
- Payment Details Belongs To Payment

**Notes:**
- Status flow: pending -> processing -> completed (atau failed/refunded)
- payment_gateway_id untuk tracking di pihak payment gateway
- Gateway response stored as JSON untuk audit trail
- Signature key untuk webhook verification
n- Integration dengan midtrans (Snap API) atau xendit
- Support retry mechanism untuk failed payments

---

## Payment Method Reference

### payment_channel Values

| Value             | Description                                 | Provider Support |
| ----------------- | ------------------------------------------- | ---------------- |
| `qris`            | QR Code pembayaran universal (QRIS)         | Midtrans, Xendit |
| `virtual_account` | Transfer bank melalui nomor Virtual Account | Midtrans, Xendit |
| `ewallet`         | E-wallet (GoPay, OVO, DANA, ShopeePay)      | Midtrans, Xendit |
| `credit_card`     | Kartu kredit / debit Visa, Mastercard       | Midtrans         |

### bank Values (Virtual Account)

| Value     | Bank                  | Provider Support |
| --------- | --------------------- | ---------------- |
| `bca`     | Bank Central Asia     | Midtrans, Xendit |
| `bni`     | Bank Negara Indonesia | Midtrans, Xendit |
| `bri`     | Bank Rakyat Indonesia | Midtrans, Xendit |
| `mandiri` | Bank Mandiri          | Midtrans, Xendit |
| `permata` | Bank Permata          | Midtrans         |
| `cimb`    | Bank CIMB Niaga       | Xendit           |

### ewallet_type Values

| Value       | E-Wallet  | Provider Support |
| ----------- | --------- | ---------------- |
| `gopay`     | GoPay     | Midtrans         |
| `ovo`       | OVO       | Xendit           |
| `dana`      | DANA      | Xendit           |
| `shopeepay` | ShopeePay | Midtrans, Xendit |

### payment_status Values

| Value       | Description                                                       |
| ----------- | ----------------------------------------------------------------- |
| `pending`   | Transaksi dibuat, menunggu pembayaran dari pasien                 |
| `paid`      | Pembayaran berhasil dikonfirmasi oleh gateway                     |
| `failed`    | Pembayaran gagal (kartu ditolak, saldo tidak cukup, dll.)         |
| `expired`   | Transaksi kedaluwarsa sebelum dibayar (default: 1 jam untuk QRIS) |
| `cancelled` | Transaksi dibatalkan secara manual                                |
| `refunded`  | Dana telah dikembalikan ke pasien                                 |

> **Catatan Integrasi:**
>
> - Konfigurasi API Key Midtrans/Xendit disimpan di environment variable server, tidak terekspos ke client
> - Untuk environment development/staging, gunakan Sandbox API Key dari masing-masing provider
> - Expired time default QRIS: **1 jam**; Virtual Account: **24 jam**; dapat dikonfigurasi per environment
> - Sistem otomatis memperbarui status `billing.payment_status` di tabel billing ketika webhook `paid` diterima
> - Satu billing yang sudah `paid` tidak dapat diinisiasi pembayaran baru
