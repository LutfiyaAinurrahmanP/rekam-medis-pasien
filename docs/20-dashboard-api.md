# Dashboard & Statistics API Documentation

## Overview

API untuk mendapatkan data ringkasan, statistik, dan laporan operasional sistem rekam medis. Setiap role mendapatkan tampilan dashboard yang berbeda sesuai dengan tanggung jawab dan kebutuhannya. Data dashboard ditampilkan dalam format ringkasan (counts, totals, averages) tanpa pagination.

**Base URL:** `/api/v1/dashboard`

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

| Endpoint                            | Patient | Doctor | Receptionist | Admin | Super Admin |
| ----------------------------------- | ------- | ------ | ------------ | ----- | ----------- |
| GET /dashboard/overview             | ❌      | ❌     | ❌           | ✅    | ✅          |
| GET /dashboard/admin                | ❌      | ❌     | ❌           | ✅    | ✅          |
| GET /dashboard/doctor               | ❌      | ✅     | ❌           | ✅    | ✅          |
| GET /dashboard/receptionist         | ❌      | ❌     | ✅           | ✅    | ✅          |
| GET /dashboard/patient              | ✅      | ❌     | ❌           | ✅    | ✅          |
| GET /dashboard/reports/appointments | ❌      | ✅     | ✅           | ✅    | ✅          |
| GET /dashboard/reports/revenue      | ❌      | ❌     | ❌           | ✅    | ✅          |
| GET /dashboard/reports/patients     | ❌      | ❌     | ✅           | ✅    | ✅          |

---

## Endpoints Summary

| Method | Endpoint                          | Description                   | Role Required                    |
| ------ | --------------------------------- | ----------------------------- | -------------------------------- |
| GET    | `/dashboard/overview`             | Global system overview        | Admin, Super Admin               |
| GET    | `/dashboard/admin`                | Admin operational dashboard   | Admin, Super Admin               |
| GET    | `/dashboard/doctor`               | Doctor personal dashboard     | Doctor                           |
| GET    | `/dashboard/receptionist`         | Receptionist daily dashboard  | Receptionist                     |
| GET    | `/dashboard/patient`              | Patient personal dashboard    | Patient                          |
| GET    | `/dashboard/reports/appointments` | Appointment statistics report | Doctor, Receptionist, Admin, SA  |
| GET    | `/dashboard/reports/revenue`      | Revenue and billing report    | Admin, Super Admin               |
| GET    | `/dashboard/reports/patients`     | Patient registration report   | Receptionist, Admin, Super Admin |

---

## Endpoints Detail

### 1. System Overview (Admin/Super Admin)

**Endpoint:** `GET /api/v1/dashboard/overview`

**Description:** Mendapatkan ringkasan statistik sistem secara keseluruhan. Berisi jumlah total data master dan aktivitas hari ini.

**Authentication:** Required (Admin, Super Admin)

**Query Parameters:**

| Parameter | Type   | Required | Default | Description                                             |
| --------- | ------ | -------- | ------- | ------------------------------------------------------- |
| date      | string | No       | today   | Tanggal referensi untuk statistik harian (`YYYY-MM-DD`) |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Dashboard overview retrieved successfully",
  "data": {
    "summary_date": "2024-01-15",
    "master_data": {
      "total_patients": 1254,
      "total_doctors": 38,
      "total_departments": 12,
      "total_rooms": 80,
      "total_medicines": 320
    },
    "today": {
      "appointments": {
        "total": 45,
        "pending": 10,
        "confirmed": 20,
        "completed": 12,
        "cancelled": 3
      },
      "new_patients": 8,
      "new_medical_records": 30,
      "active_hospitalizations": 22
    },
    "rooms": {
      "total": 80,
      "available": 58,
      "occupied": 22,
      "out_of_service": 0,
      "occupancy_rate": 27.5
    },
    "billing": {
      "today_revenue": 12500000,
      "unpaid_count": 15,
      "unpaid_total": 8750000
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/overview?date=2024-01-15" \
  -H "Authorization: Bearer <token>"
```

---

### 2. Admin Dashboard

**Endpoint:** `GET /api/v1/dashboard/admin`

**Description:** Dashboard operasional detail untuk admin dan super admin. Mencakup statistik appointments, rawat inap, billing, dan tren registrasi dalam periode tertentu.

**Authentication:** Required (Admin, Super Admin)

**Query Parameters:**

| Parameter  | Type   | Required | Default | Description                                                                 |
| ---------- | ------ | -------- | ------- | --------------------------------------------------------------------------- |
| period     | string | No       | today   | Periode laporan: `today`, `this_week`, `this_month`, `last_month`, `custom` |
| start_date | string | No       | -       | Tanggal mulai (format `YYYY-MM-DD`, wajib jika `period=custom`)             |
| end_date   | string | No       | -       | Tanggal akhir (format `YYYY-MM-DD`, wajib jika `period=custom`)             |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Admin dashboard retrieved successfully",
  "data": {
    "period": "this_month",
    "period_range": {
      "start": "2024-01-01",
      "end": "2024-01-31"
    },
    "appointments": {
      "total": 780,
      "completed": 650,
      "cancelled": 45,
      "no_show": 20,
      "completion_rate": 83.3
    },
    "patients": {
      "total_active": 1254,
      "new_registrations": 85,
      "returning": 695
    },
    "hospitalization": {
      "total_admissions": 42,
      "currently_hospitalized": 22,
      "total_discharged": 20,
      "average_length_of_stay_days": 3.8,
      "available_beds": 58
    },
    "billing": {
      "total_revenue": 254000000,
      "paid_count": 720,
      "unpaid_count": 60,
      "unpaid_total": 42000000,
      "average_bill": 352778
    },
    "referrals": {
      "total_issued": 35,
      "internal": 25,
      "external": 10,
      "pending": 8,
      "completed": 22
    },
    "top_departments": [
      {
        "department_id": 1,
        "department_name": "Poli Umum",
        "appointment_count": 210
      },
      {
        "department_id": 4,
        "department_name": "Kardiologi",
        "appointment_count": 145
      },
      {
        "department_id": 3,
        "department_name": "Neurologi",
        "appointment_count": 120
      }
    ],
    "appointment_trend": [
      { "date": "2024-01-01", "count": 22 },
      { "date": "2024-01-02", "count": 30 },
      { "date": "2024-01-03", "count": 28 }
    ],
    "revenue_trend": [
      { "date": "2024-01-01", "amount": 7800000 },
      { "date": "2024-01-02", "amount": 10200000 },
      { "date": "2024-01-03", "amount": 9500000 }
    ]
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/admin?period=this_month" \
  -H "Authorization: Bearer <token>"
```

```bash
# Custom period
curl -X GET "http://localhost:8080/api/v1/dashboard/admin?period=custom&start_date=2024-01-01&end_date=2024-01-15" \
  -H "Authorization: Bearer <token>"
```

---

### 3. Doctor Dashboard

**Endpoint:** `GET /api/v1/dashboard/doctor`

**Description:** Dashboard personal untuk dokter. Menampilkan jadwal hari ini, pasien tertangani, dan status lab/resep yang pending.

**Authentication:** Required (Doctor)

**Query Parameters:**

| Parameter | Type   | Required | Default | Description                                          |
| --------- | ------ | -------- | ------- | ---------------------------------------------------- |
| date      | string | No       | today   | Tanggal referensi untuk jadwal harian (`YYYY-MM-DD`) |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Doctor dashboard retrieved successfully",
  "data": {
    "doctor": {
      "id": 2,
      "name": "dr. Siti Rahayu, SpPD",
      "specialization": "Penyakit Dalam"
    },
    "date": "2024-01-15",
    "today_schedule": {
      "total_appointments": 12,
      "completed": 5,
      "pending": 6,
      "cancelled": 1,
      "next_patient": {
        "appointment_id": 45,
        "appointment_time": "10:30",
        "patient_name": "Budi Santoso",
        "chief_complaint": "Jantung berdebar-debar"
      }
    },
    "upcoming_appointments": [
      {
        "appointment_id": 45,
        "appointment_time": "10:30",
        "patient_name": "Budi Santoso",
        "patient_id": 10
      },
      {
        "appointment_id": 46,
        "appointment_time": "11:00",
        "patient_name": "Dewi Lestari",
        "patient_id": 11
      }
    ],
    "statistics": {
      "total_patients_this_month": 145,
      "total_patients_today": 12,
      "pending_lab_results": 8,
      "pending_prescriptions": 3,
      "pending_referrals": 4
    },
    "recent_medical_records": [
      {
        "id": 120,
        "patient_name": "Ahmad Wijaya",
        "visit_date": "2024-01-15",
        "chief_complaint": "Sesak napas",
        "status": "final"
      }
    ]
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/doctor" \
  -H "Authorization: Bearer <token>"
```

```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/doctor?date=2024-01-15" \
  -H "Authorization: Bearer <token>"
```

---

### 4. Receptionist Dashboard

**Endpoint:** `GET /api/v1/dashboard/receptionist`

**Description:** Dashboard harian untuk resepsionis, menampilkan antrian appointment hari ini, status rawat inap, dan tagihan belum terbayar.

**Authentication:** Required (Receptionist)

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Receptionist dashboard retrieved successfully",
  "data": {
    "date": "2024-01-15",
    "appointments_today": {
      "total": 45,
      "pending": 10,
      "confirmed": 20,
      "in_progress": 8,
      "completed": 5,
      "cancelled": 2,
      "appointment_queue": [
        {
          "appointment_id": 42,
          "scheduled_time": "09:00",
          "patient_name": "Ahmad Wijaya",
          "patient_id": 5,
          "doctor_name": "dr. Siti Rahayu, SpPD",
          "status": "confirmed"
        },
        {
          "appointment_id": 43,
          "scheduled_time": "09:30",
          "patient_name": "Dewi Lestari",
          "patient_id": 11,
          "doctor_name": "dr. Budi Hartono, SpA",
          "status": "pending"
        }
      ]
    },
    "hospitalization": {
      "currently_hospitalized": 22,
      "available_beds": 58,
      "new_admissions_today": 4,
      "discharged_today": 2
    },
    "billing": {
      "unpaid_count": 15,
      "unpaid_total": 8750000
    },
    "new_patient_registrations_today": 8
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/receptionist" \
  -H "Authorization: Bearer <token>"
```

---

### 5. Patient Dashboard

**Endpoint:** `GET /api/v1/dashboard/patient`

**Description:** Dashboard personal untuk pasien. Menampilkan jadwal mendatang, tagihan, resep aktif, dan riwayat kunjungan singkat.

**Authentication:** Required (Patient)

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Patient dashboard retrieved successfully",
  "data": {
    "patient": {
      "id": 10,
      "name": "Budi Santoso",
      "medical_record_number": "MRN-2024-000010",
      "date_of_birth": "1985-03-20",
      "age": 38
    },
    "upcoming_appointments": [
      {
        "id": 50,
        "scheduled_date": "2024-01-20",
        "scheduled_time": "10:00",
        "doctor_name": "dr. Ahmad Fauzi, SpJP",
        "department_name": "Kardiologi",
        "status": "confirmed"
      }
    ],
    "billing": {
      "unpaid_count": 2,
      "unpaid_total": 350000,
      "unpaid_bills": [
        {
          "id": 88,
          "invoice_number": "INV-2024-000088",
          "total_amount": 250000,
          "due_date": "2024-01-25",
          "status": "unpaid"
        }
      ]
    },
    "active_prescriptions": [
      {
        "id": 12,
        "prescribed_date": "2024-01-10",
        "doctor_name": "dr. Siti Rahayu, SpPD",
        "medicines_count": 3,
        "status": "active"
      }
    ],
    "pending_lab_results": [
      {
        "id": 7,
        "test_name": "Elektrokardiogram (EKG)",
        "ordered_date": "2024-01-15",
        "status": "processing"
      }
    ],
    "active_referrals": [
      {
        "id": 1,
        "referral_number": "REF-2024-000001",
        "referred_to": "Kardiologi — dr. Ahmad Fauzi, SpJP",
        "referral_date": "2024-01-15",
        "status": "accepted"
      }
    ],
    "recent_visits": [
      {
        "id": 5,
        "visit_date": "2024-01-15",
        "doctor_name": "dr. Siti Rahayu, SpPD",
        "chief_complaint": "Jantung berdebar-debar",
        "status": "final"
      },
      {
        "id": 4,
        "visit_date": "2023-12-20",
        "doctor_name": "dr. Siti Rahayu, SpPD",
        "chief_complaint": "Kontrol tekanan darah",
        "status": "final"
      }
    ]
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/patient" \
  -H "Authorization: Bearer <token>"
```

---

### 6. Appointment Statistics Report

**Endpoint:** `GET /api/v1/dashboard/reports/appointments`

**Description:** Laporan statistik appointment berdasarkan periode dan filter tertentu. Berguna untuk analisis kapasitas dan perencanaan jadwal.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Query Parameters:**

| Parameter     | Type    | Required | Default    | Description                                                                     |
| ------------- | ------- | -------- | ---------- | ------------------------------------------------------------------------------- |
| period        | string  | No       | this_month | Periode: `today`, `this_week`, `this_month`, `last_month`, `custom`             |
| start_date    | string  | No       | -          | Tanggal mulai (`YYYY-MM-DD`), wajib jika `period=custom`                        |
| end_date      | string  | No       | -          | Tanggal akhir (`YYYY-MM-DD`), wajib jika `period=custom`                        |
| doctor_id     | integer | No       | -          | Filter berdasarkan dokter tertentu (doctor hanya bisa melihat miliknya sendiri) |
| department_id | integer | No       | -          | Filter berdasarkan departemen                                                   |
| group_by      | string  | No       | day        | Granularitas trend: `day`, `week`, `month`                                      |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Appointment report retrieved successfully",
  "data": {
    "period": "this_month",
    "period_range": {
      "start": "2024-01-01",
      "end": "2024-01-31"
    },
    "totals": {
      "total": 780,
      "scheduled": 50,
      "confirmed": 60,
      "completed": 620,
      "cancelled": 38,
      "no_show": 12,
      "completion_rate": 79.5,
      "cancellation_rate": 4.9
    },
    "by_department": [
      {
        "department_id": 1,
        "department_name": "Poli Umum",
        "total": 210,
        "completed": 180,
        "completion_rate": 85.7
      },
      {
        "department_id": 4,
        "department_name": "Kardiologi",
        "total": 145,
        "completed": 120,
        "completion_rate": 82.8
      }
    ],
    "by_doctor": [
      {
        "doctor_id": 2,
        "doctor_name": "dr. Siti Rahayu, SpPD",
        "total": 145,
        "completed": 128,
        "completion_rate": 88.3
      }
    ],
    "trend": [
      { "date": "2024-01-01", "total": 22, "completed": 18, "cancelled": 2 },
      { "date": "2024-01-02", "total": 30, "completed": 26, "cancelled": 2 },
      { "date": "2024-01-03", "total": 28, "completed": 25, "cancelled": 1 }
    ],
    "peak_hours": [
      { "hour": "09:00", "count": 95 },
      { "hour": "10:00", "count": 120 },
      { "hour": "11:00", "count": 100 }
    ]
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/reports/appointments?period=this_month&department_id=4" \
  -H "Authorization: Bearer <token>"
```

---

### 7. Revenue & Billing Report

**Endpoint:** `GET /api/v1/dashboard/reports/revenue`

**Description:** Laporan pendapatan dan billing dalam periode tertentu beserta tren harian/mingguan/bulanan.

**Authentication:** Required (Admin, Super Admin)

**Query Parameters:**

| Parameter  | Type   | Required | Default    | Description                                                         |
| ---------- | ------ | -------- | ---------- | ------------------------------------------------------------------- |
| period     | string | No       | this_month | Periode: `today`, `this_week`, `this_month`, `last_month`, `custom` |
| start_date | string | No       | -          | Tanggal mulai (`YYYY-MM-DD`), wajib jika `period=custom`            |
| end_date   | string | No       | -          | Tanggal akhir (`YYYY-MM-DD`), wajib jika `period=custom`            |
| group_by   | string | No       | day        | Granularitas trend: `day`, `week`, `month`                          |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Revenue report retrieved successfully",
  "data": {
    "period": "this_month",
    "period_range": {
      "start": "2024-01-01",
      "end": "2024-01-31"
    },
    "revenue": {
      "total_billed": 310000000,
      "total_paid": 254000000,
      "total_unpaid": 56000000,
      "total_bills": 780,
      "paid_bills": 720,
      "unpaid_bills": 60,
      "average_bill_amount": 397436,
      "collection_rate": 81.9
    },
    "by_category": [
      {
        "category": "Consultation",
        "total_amount": 120000000,
        "bill_count": 620
      },
      {
        "category": "Laboratory",
        "total_amount": 85000000,
        "bill_count": 340
      },
      {
        "category": "Hospitalization",
        "total_amount": 65000000,
        "bill_count": 42
      },
      {
        "category": "Pharmacy",
        "total_amount": 40000000,
        "bill_count": 510
      }
    ],
    "trend": [
      { "date": "2024-01-01", "billed": 9500000, "paid": 7800000 },
      { "date": "2024-01-02", "billed": 12000000, "paid": 10200000 },
      { "date": "2024-01-03", "billed": 11500000, "paid": 9500000 }
    ],
    "comparison": {
      "vs_previous_period": {
        "revenue_change_percent": 12.5,
        "bill_count_change_percent": 8.2
      }
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/reports/revenue?period=this_month&group_by=day" \
  -H "Authorization: Bearer <token>"
```

---

### 8. Patient Registration Report

**Endpoint:** `GET /api/v1/dashboard/reports/patients`

**Description:** Laporan statistik registrasi, kunjungan, dan demografi pasien.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Query Parameters:**

| Parameter  | Type   | Required | Default    | Description                                                         |
| ---------- | ------ | -------- | ---------- | ------------------------------------------------------------------- |
| period     | string | No       | this_month | Periode: `today`, `this_week`, `this_month`, `last_month`, `custom` |
| start_date | string | No       | -          | Tanggal mulai (`YYYY-MM-DD`), wajib jika `period=custom`            |
| end_date   | string | No       | -          | Tanggal akhir (`YYYY-MM-DD`), wajib jika `period=custom`            |

**Response Success (200):**

```json
{
  "status": "success",
  "message": "Patient report retrieved successfully",
  "data": {
    "period": "this_month",
    "period_range": {
      "start": "2024-01-01",
      "end": "2024-01-31"
    },
    "registrations": {
      "new_patients": 85,
      "returning_patients": 695,
      "total_visits": 780,
      "unique_visitors": 650
    },
    "demographics": {
      "by_gender": [
        { "gender": "male", "count": 580 },
        { "gender": "female", "count": 674 }
      ],
      "by_age_group": [
        { "group": "0-12", "label": "Anak", "count": 120 },
        { "group": "13-17", "label": "Remaja", "count": 85 },
        { "group": "18-40", "label": "Dewasa Muda", "count": 380 },
        { "group": "41-60", "label": "Dewasa", "count": 450 },
        { "group": "61+", "label": "Lansia", "count": 219 }
      ]
    },
    "registration_trend": [
      { "date": "2024-01-01", "new_patients": 4, "total_visits": 22 },
      { "date": "2024-01-02", "new_patients": 6, "total_visits": 30 },
      { "date": "2024-01-03", "new_patients": 3, "total_visits": 28 }
    ],
    "comparison": {
      "vs_previous_period": {
        "new_patients_change_percent": 18.0,
        "total_visits_change_percent": 9.5
      }
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/dashboard/reports/patients?period=this_month" \
  -H "Authorization: Bearer <token>"
```

---

## Error Responses

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
  "message": "Forbidden: insufficient permissions to access this dashboard"
}
```

### 400 Bad Request

```json
{
  "status": "error",
  "message": "Validation error",
  "errors": {
    "start_date": "start_date is required when period is 'custom'",
    "end_date": "end_date is required when period is 'custom'"
  }
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

> **Note:** Dashboard endpoints tidak memiliki model penyimpanan dedicated. Dashboard mengagregasi data dari berbagai tabel lain (Users, Patients, Appointments, Medical Records, Billing, dll) untuk menghasilkan summary dan statistics.

### Aggregated Data Sources:

| Dashboard Type | Primary Data Sources | Metrics Shown |
| --- | --- | --- |
| Admin Dashboard | Users, Patients, Appointments, Medical Records, Billing | Total users by role, new registrations, appointments, revenue |
| Doctor Dashboard | Appointments, Medical Records, Patients, Lab Tests | My appointments, patients today, pending records, lab requests |
| Receptionist Dashboard | Appointments, Patients, Billing, Hospitalizations | Today's appointments, new registrations, pending payments, occupancy |
| Patient Dashboard | Medical Records, Appointments, Billing, Prescriptions | My records, appointments, bills, medications |

### Typical Data Aggregation:

```sql
-- Example: Count appointments by status
SELECT status, COUNT(*) as count 
FROM appointments 
WHERE appointment_date >= CURDATE() 
GROUP BY status;

-- Example: Revenue summary by date
SELECT DATE(payment_date), SUM(amount) 
FROM payments 
WHERE status = 'completed' 
GROUP BY DATE(payment_date);
```

**Notes:**
- Dashboard data typically cached untuk performance
- Real-time atau periodic refresh tergantung requirements
- No direct database model, hanya queries aggregation
- Security filtered berdasarkan user role dan department access

---

## Reference Values

**Period Values:**

- `today` — Data hari ini (00:00 sampai 23:59)
- `this_week` — Minggu berjalan (Senin–Minggu)
- `this_month` — Bulan berjalan (tanggal 1 hingga hari ini)
- `last_month` — Bulan kalender sebelumnya (penuh)
- `custom` — Rentang tanggal bebas, wajib mengisi `start_date` dan `end_date`

**group_by Values (untuk trend):**

- `day` — Dikelompokkan per hari
- `week` — Dikelompokkan per minggu
- `month` — Dikelompokkan per bulan

> **Notes:**
>
> - Dashboard doctor hanya menampilkan data dokter yang sedang login.
> - Dashboard patient hanya menampilkan data pasien yang sedang login.
> - Admin dapat mengakses dashboard doctor dengan menambahkan query `doctor_id` pada endpoint reports/appointments.
> - Semua nilai mata uang (revenue, billing) dalam satuan Rupiah (IDR).
> - `completion_rate` dan `collection_rate` dihitung dalam persen (%).
> - Data trend menggunakan tanggal kalender lokal (WIB, UTC+7).
