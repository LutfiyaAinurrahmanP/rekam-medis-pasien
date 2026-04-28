# Appointments API Documentation

## Overview

API untuk manajemen data appointments (jadwal konsultasi/pertemuan) antara patient dan doctor dalam sistem rekam medis. Appointments adalah booking jadwal untuk konsultasi medis.

**Base URL:** `/api/v1/appointments`

---

## Table of Contents

- [Authentication](#authentication)
- [Authorization](#authorization)
- [Endpoints Summary](#endpoints-summary)
- [Self-Owned Endpoints](#self-owned-endpoints)
- [Public Endpoints](#public-endpoints)
- [Staff Endpoints](#staff-endpoints)
- [Admin Endpoints](#admin-endpoints)
- [Super Admin Endpoints](#super-admin-endpoints)
- [Database Model](#database-model)
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

| Endpoint                             | Patient (Own) | Doctor (Own) | Receptionist | Admin | Super Admin |
| ------------------------------------ | ------------- | ------------ | ------------ | ----- | ----------- |
| GET /appointments/my-appointments    | ✅            | ✅           | ❌           | ❌    | ❌          |
| POST /appointments                   | ✅            | ❌           | ✅           | ✅    | ✅          |
| GET /appointments/my-schedule        | ❌            | ✅           | ❌           | ❌    | ❌          |
| GET /appointments                    | ❌            | ✅           | ✅           | ✅    | ✅          |
| GET /appointments/today              | ❌            | ✅           | ✅           | ✅    | ✅          |
| GET /appointments/upcoming           | ✅            | ✅           | ✅           | ✅    | ✅          |
| GET /appointments/past               | ✅            | ✅           | ✅           | ✅    | ✅          |
| GET /appointments/cancelled          | ❌            | ❌           | ✅           | ✅    | ✅          |
| GET /appointments/deleted            | ❌            | ❌           | ❌           | ✅    | ✅          |
| GET /appointments/:id                | ✅ (Own)      | ✅ (Own)     | ✅           | ✅    | ✅          |
| PUT /appointments/:id                | ✅ (Own)      | ❌           | ✅           | ✅    | ✅          |
| PATCH /appointments/:id/confirm      | ❌            | ✅           | ✅           | ✅    | ✅          |
| PATCH /appointments/:id/start        | ❌            | ✅           | ❌           | ❌    | ❌          |
| PATCH /appointments/:id/complete     | ❌            | ✅           | ❌           | ❌    | ❌          |
| PATCH /appointments/:id/cancel       | ✅ (Own)      | ✅ (Own)     | ✅           | ✅    | ✅          |
| PATCH /appointments/:id/reschedule   | ✅ (Own)      | ❌           | ✅           | ✅    | ✅          |
| PATCH /appointments/:id/no-show      | ❌            | ✅           | ✅           | ✅    | ✅          |
| DELETE /appointments/:id             | ❌            | ❌           | ❌           | ✅    | ✅          |
| PATCH /appointments/:id/restore      | ❌            | ❌           | ❌           | ✅    | ✅          |
| DELETE /appointments/:id/hard-delete | ❌            | ❌           | ❌           | ❌    | ✅          |

---

## Endpoints Summary

### Self-Owned Endpoints

| Method | Endpoint                        | Description                          | Auth            |
| ------ | ------------------------------- | ------------------------------------ | --------------- |
| GET    | `/appointments/my-appointments` | Get my appointments (patient/doctor) | Patient, Doctor |
| GET    | `/appointments/my-schedule`     | Get my schedule (doctor only)        | Doctor          |

### Public Endpoints

| Method | Endpoint                  | Description            | Role Required                             |
| ------ | ------------------------- | ---------------------- | ----------------------------------------- |
| POST   | `/appointments`           | Create appointment     | Patient, Receptionist, Admin, Super Admin |
| GET    | `/appointments`           | List all appointments  | Doctor, Receptionist, Admin, Super Admin  |
| GET    | `/appointments/today`     | Today's appointments   | Doctor, Receptionist, Admin, Super Admin  |
| GET    | `/appointments/upcoming`  | Upcoming appointments  | All (filtered by ownership)               |
| GET    | `/appointments/past`      | Past appointments      | All (filtered by ownership)               |
| GET    | `/appointments/cancelled` | Cancelled appointments | Receptionist, Admin, Super Admin          |
| GET    | `/appointments/deleted`   | Deleted appointments   | Admin, Super Admin                        |
| GET    | `/appointments/:id`       | Get appointment by ID  | All (with ownership check)                |

### Staff Endpoints

| Method | Endpoint                       | Description            | Role Required                                                 |
| ------ | ------------------------------ | ---------------------- | ------------------------------------------------------------- |
| PUT    | `/appointments/:id`            | Update appointment     | Patient (own), Receptionist, Admin, Super Admin               |
| PATCH  | `/appointments/:id/confirm`    | Confirm appointment    | Doctor, Receptionist, Admin, Super Admin                      |
| PATCH  | `/appointments/:id/start`      | Start appointment      | Doctor                                                        |
| PATCH  | `/appointments/:id/complete`   | Complete appointment   | Doctor                                                        |
| PATCH  | `/appointments/:id/cancel`     | Cancel appointment     | Patient (own), Doctor (own), Receptionist, Admin, Super Admin |
| PATCH  | `/appointments/:id/reschedule` | Reschedule appointment | Patient (own), Receptionist, Admin, Super Admin               |
| PATCH  | `/appointments/:id/no-show`    | Mark as no-show        | Doctor, Receptionist, Admin, Super Admin                      |

### Admin Endpoints

| Method | Endpoint                    | Description                 | Role Required      |
| ------ | --------------------------- | --------------------------- | ------------------ |
| DELETE | `/appointments/:id`         | Soft delete appointment     | Admin, Super Admin |
| PATCH  | `/appointments/:id/restore` | Restore deleted appointment | Admin, Super Admin |

### Super Admin Endpoints

| Method | Endpoint                        | Description                    | Role Required |
| ------ | ------------------------------- | ------------------------------ | ------------- |
| DELETE | `/appointments/:id/hard-delete` | Permanently delete appointment | Super Admin   |

---

## Self-Owned Endpoints

### 1. Get My Appointments

**Endpoint:** `GET /api/v1/appointments/my-appointments`

**Description:** Patient/Doctor mendapatkan daftar appointments milik sendiri.

**Authentication:** Required (Patient/Doctor Role)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter   | Type    | Default          | Description                        |
| ----------- | ------- | ---------------- | ---------------------------------- |
| `page`      | integer | 1                | Halaman                            |
| `page_size` | integer | 10               | Jumlah data per halaman (max: 100) |
| `status`    | string  | -                | Filter by status                   |
| `date_from` | date    | -                | Filter from date (YYYY-MM-DD)      |
| `date_to`   | date    | -                | Filter to date (YYYY-MM-DD)        |
| `sort_by`   | string  | appointment_date | Sort field                         |
| `sort_dir`  | string  | desc             | Sort direction (asc, desc)         |

**Example Request:**

```
GET /api/v1/appointments/my-appointments?status=scheduled&date_from=2024-01-20&date_to=2024-01-31
```

**Response Success (200 OK) - Patient:**

```json
{
  "success": true,
  "message": "My appointments retrieved successfully",
  "data": {
    "user_role": "patient",
    "total_appointments": 5,
    "upcoming_count": 2,
    "completed_count": 3,
    "data": [
      {
        "id": 1,
        "appointment_date": "2024-01-25",
        "appointment_time": "10:00:00",
        "duration_minutes": 30,
        "status": "scheduled",
        "doctor": {
          "id": 5,
          "full_name": "Dr. John Smith, Sp.JP",
          "specialization": "Cardiologist",
          "department": "Kardiologi"
        },
        "reason": "Kontrol rutin jantung",
        "notes": null,
        "created_at": "2024-01-19T10:00:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total_items": 5,
      "total_pages": 1
    }
  }
}
```

**Response Success (200 OK) - Doctor:**

```json
{
  "success": true,
  "message": "My schedule retrieved successfully",
  "data": {
    "user_role": "doctor",
    "total_appointments": 15,
    "today_count": 3,
    "upcoming_count": 12,
    "data": [
      {
        "id": 1,
        "appointment_date": "2024-01-20",
        "appointment_time": "10:00:00",
        "duration_minutes": 30,
        "status": "confirmed",
        "patient": {
          "id": 10,
          "patient_code": "P-2024-001",
          "full_name": "Jane Doe",
          "phone": "081234567890"
        },
        "reason": "Kontrol rutin jantung",
        "notes": "Patient has history of hypertension",
        "created_at": "2024-01-19T10:00:00Z"
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
curl -X GET "http://localhost:8080/api/v1/appointments/my-appointments?status=scheduled" \
  -H "Authorization: Bearer PATIENT_JWT_TOKEN"
```

**Use Case:**

- Patient melihat jadwal konsultasi
- Doctor melihat jadwal praktik
- Planning dan reminder

---

### 2. Get My Schedule (Doctor Only)

**Endpoint:** `GET /api/v1/appointments/my-schedule`

**Description:** Doctor mendapatkan jadwal praktik dengan detail lengkap.

**Authentication:** Required (Doctor Role)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter   | Type   | Description                                |
| ----------- | ------ | ------------------------------------------ |
| `date`      | date   | Specific date (YYYY-MM-DD), default: today |
| `date_from` | date   | Date range from                            |
| `date_to`   | date   | Date range to                              |
| `status`    | string | Filter by status                           |

**Example Request:**

```
GET /api/v1/appointments/my-schedule?date=2024-01-20
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "My schedule retrieved successfully",
  "data": {
    "date": "2024-01-20",
    "doctor_info": {
      "id": 5,
      "full_name": "Dr. John Smith, Sp.JP",
      "specialization": "Cardiologist",
      "department": "Kardiologi"
    },
    "summary": {
      "total_appointments": 8,
      "scheduled": 3,
      "confirmed": 2,
      "in_progress": 1,
      "completed": 2,
      "cancelled": 0,
      "no_show": 0
    },
    "schedule": [
      {
        "id": 1,
        "time": "08:00:00",
        "duration_minutes": 30,
        "status": "completed",
        "patient": {
          "patient_code": "P-2024-001",
          "full_name": "Jane Doe",
          "age": 45,
          "gender": "female"
        },
        "reason": "Kontrol rutin hipertensi",
        "is_new_patient": false
      },
      {
        "id": 2,
        "time": "08:30:00",
        "duration_minutes": 30,
        "status": "completed",
        "patient": {
          "patient_code": "P-2024-025",
          "full_name": "John Doe",
          "age": 50,
          "gender": "male"
        },
        "reason": "Chest pain evaluation",
        "is_new_patient": true
      },
      {
        "id": 3,
        "time": "09:00:00",
        "duration_minutes": 30,
        "status": "in_progress",
        "patient": {
          "patient_code": "P-2024-015",
          "full_name": "Mary Smith",
          "age": 38,
          "gender": "female"
        },
        "reason": "Follow-up after angioplasty",
        "is_new_patient": false
      },
      {
        "id": 4,
        "time": "09:30:00",
        "duration_minutes": 30,
        "status": "confirmed",
        "patient": {
          "patient_code": "P-2024-030",
          "full_name": "Robert Johnson",
          "age": 55,
          "gender": "male"
        },
        "reason": "Consultation for chest pain",
        "is_new_patient": true
      },
      {
        "time": "10:00:00",
        "status": "available",
        "slot_available": true
      },
      {
        "id": 5,
        "time": "10:30:00",
        "duration_minutes": 30,
        "status": "scheduled",
        "patient": {
          "patient_code": "P-2024-008",
          "full_name": "Sarah Williams",
          "age": 42,
          "gender": "female"
        },
        "reason": "Regular checkup",
        "is_new_patient": false
      }
    ]
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/appointments/my-schedule?date=2024-01-20" \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN"
```

**Use Case:**

- Doctor daily schedule overview
- Time management
- Patient preparation

---

## Public Endpoints

### 3. Create Appointment

**Endpoint:** `POST /api/v1/appointments`

**Description:** Membuat appointment baru (patient self-booking atau receptionist).

**Authentication:** Required (Patient, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "patient_id": 10,
  "doctor_id": 5,
  "appointment_date": "2024-01-25",
  "appointment_time": "10:00:00",
  "duration_minutes": 30,
  "reason": "Kontrol rutin jantung",
  "notes": "Patient has history of hypertension"
}
```

**Field Rules:**

- `patient_id`: required, FK to patients table
- `doctor_id`: required, FK to doctors table
- `appointment_date`: required, date (YYYY-MM-DD), cannot be in the past
- `appointment_time`: required, time (HH:MM:SS)
- `duration_minutes`: required, min 15, typically 30 or 60
- `reason`: optional, text
- `notes`: optional, text
- `status`: auto-set to 'scheduled'

**Response Success (201 Created):**

```json
{
  "success": true,
  "message": "Appointment created successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "patient": {
      "patient_code": "P-2024-001",
      "full_name": "Jane Doe",
      "phone": "081234567890"
    },
    "doctor_id": 5,
    "doctor": {
      "full_name": "Dr. John Smith, Sp.JP",
      "specialization": "Cardiologist",
      "department": "Kardiologi"
    },
    "appointment_date": "2024-01-25",
    "appointment_time": "10:00:00",
    "duration_minutes": 30,
    "reason": "Kontrol rutin jantung",
    "status": "scheduled",
    "notes": "Patient has history of hypertension",
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z"
  }
}
```

**Response Error (400 Bad Request):**

```json
{
  "success": false,
  "message": "Failed to create appointment",
  "error": "doctor not available at the selected time slot"
}
```

**Response Error (409 Conflict):**

```json
{
  "success": false,
  "message": "Time slot conflict",
  "error": "doctor already has an appointment at this time",
  "suggested_slots": ["10:30:00", "11:00:00", "11:30:00"]
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/appointments \
  -H "Authorization: Bearer PATIENT_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 10,
    "doctor_id": 5,
    "appointment_date": "2024-01-25",
    "appointment_time": "10:00:00",
    "duration_minutes": 30,
    "reason": "Kontrol rutin jantung"
  }'
```

**Business Rules:**

- Cannot book appointment in the past
- Check doctor availability
- Check time slot conflicts
- Check doctor working hours
- Maximum appointments per day limit
- Minimum booking time advance (e.g., 2 hours)

---

### 4. List All Appointments

**Endpoint:** `GET /api/v1/appointments`

**Description:** Mendapatkan daftar semua appointments dengan filter (staff only).

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter       | Type    | Default          | Description                        |
| --------------- | ------- | ---------------- | ---------------------------------- |
| `page`          | integer | 1                | Halaman                            |
| `page_size`     | integer | 10               | Jumlah data per halaman (max: 100) |
| `patient_id`    | integer | -                | Filter by patient                  |
| `doctor_id`     | integer | -                | Filter by doctor                   |
| `department_id` | integer | -                | Filter by department               |
| `status`        | string  | -                | Filter by status                   |
| `date`          | date    | -                | Filter by specific date            |
| `date_from`     | date    | -                | Filter from date                   |
| `date_to`       | date    | -                | Filter to date                     |
| `sort_by`       | string  | appointment_date | Sort field                         |
| `sort_dir`      | string  | asc              | Sort direction                     |

**Example Request:**

```
GET /api/v1/appointments?doctor_id=5&status=scheduled&date_from=2024-01-20&date_to=2024-01-31
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Appointments retrieved successfully",
  "data": {
    "filters": {
      "doctor_id": 5,
      "status": "scheduled",
      "date_from": "2024-01-20",
      "date_to": "2024-01-31"
    },
    "summary": {
      "total_appointments": 50,
      "scheduled": 20,
      "confirmed": 15,
      "in_progress": 2,
      "completed": 10,
      "cancelled": 3
    },
    "data": [
      {
        "id": 1,
        "patient": {
          "id": 10,
          "patient_code": "P-2024-001",
          "full_name": "Jane Doe"
        },
        "doctor": {
          "id": 5,
          "full_name": "Dr. John Smith, Sp.JP",
          "specialization": "Cardiologist"
        },
        "appointment_date": "2024-01-25",
        "appointment_time": "10:00:00",
        "duration_minutes": 30,
        "status": "scheduled",
        "reason": "Kontrol rutin jantung",
        "created_at": "2024-01-19T10:00:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total_items": 50,
      "total_pages": 5
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/appointments?doctor_id=5&status=scheduled" \
  -H "Authorization: Bearer RECEPTIONIST_JWT_TOKEN"
```

---

### 5. Today's Appointments

**Endpoint:** `GET /api/v1/appointments/today`

**Description:** Mendapatkan daftar appointments hari ini.

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter       | Type    | Description          |
| --------------- | ------- | -------------------- |
| `doctor_id`     | integer | Filter by doctor     |
| `department_id` | integer | Filter by department |
| `status`        | string  | Filter by status     |

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Today's appointments retrieved successfully",
  "data": {
    "date": "2024-01-20",
    "total_today": 25,
    "summary": {
      "scheduled": 8,
      "confirmed": 5,
      "in_progress": 2,
      "completed": 7,
      "cancelled": 2,
      "no_show": 1
    },
    "data": [
      {
        "id": 1,
        "time": "08:00:00",
        "patient": {
          "patient_code": "P-2024-001",
          "full_name": "Jane Doe"
        },
        "doctor": {
          "full_name": "Dr. John Smith, Sp.JP",
          "department": "Kardiologi"
        },
        "status": "completed",
        "reason": "Kontrol rutin"
      }
    ],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total_items": 25,
      "total_pages": 3
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/appointments/today?doctor_id=5" \
  -H "Authorization: Bearer RECEPTIONIST_JWT_TOKEN"
```

---

### 6. Upcoming Appointments

**Endpoint:** `GET /api/v1/appointments/upcoming`

**Description:** Mendapatkan daftar appointments yang akan datang.

**Authentication:** Required (All Authenticated, filtered by ownership)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter    | Type    | Description                     |
| ------------ | ------- | ------------------------------- |
| `days_ahead` | integer | Days to look ahead (default: 7) |
| `patient_id` | integer | Filter by patient (staff only)  |
| `doctor_id`  | integer | Filter by doctor (staff only)   |
| `page`       | integer | Page number                     |
| `page_size`  | integer | Items per page                  |

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Upcoming appointments retrieved successfully",
  "data": {
    "date_range": {
      "from": "2024-01-21",
      "to": "2024-01-28"
    },
    "total_upcoming": 12,
    "data": [
      {
        "id": 1,
        "appointment_date": "2024-01-25",
        "appointment_time": "10:00:00",
        "patient": {
          "patient_code": "P-2024-001",
          "full_name": "Jane Doe"
        },
        "doctor": {
          "full_name": "Dr. John Smith, Sp.JP",
          "specialization": "Cardiologist"
        },
        "status": "scheduled",
        "days_until": 5
      }
    ],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total_items": 12,
      "total_pages": 2
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/appointments/upcoming?days_ahead=7" \
  -H "Authorization: Bearer PATIENT_JWT_TOKEN"
```

---

### 7. Past Appointments

**Endpoint:** `GET /api/v1/appointments/past`

**Description:** Mendapatkan daftar appointments yang sudah lewat.

**Authentication:** Required (All Authenticated, filtered by ownership)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter    | Type    | Description                     |
| ------------ | ------- | ------------------------------- |
| `days_back`  | integer | Days to look back (default: 30) |
| `patient_id` | integer | Filter by patient (staff only)  |
| `doctor_id`  | integer | Filter by doctor (staff only)   |
| `page`       | integer | Page number                     |
| `page_size`  | integer | Items per page                  |

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Past appointments retrieved successfully",
  "data": {
    "date_range": {
      "from": "2023-12-21",
      "to": "2024-01-20"
    },
    "total_past": 20,
    "data": [
      {
        "id": 1,
        "appointment_date": "2024-01-15",
        "appointment_time": "10:00:00",
        "patient": {
          "patient_code": "P-2024-001",
          "full_name": "Jane Doe"
        },
        "doctor": {
          "full_name": "Dr. John Smith, Sp.JP",
          "specialization": "Cardiologist"
        },
        "status": "completed",
        "days_ago": 5,
        "has_medical_record": true
      }
    ],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total_items": 20,
      "total_pages": 2
    }
  }
}
```

**cURL Example:**

```bash
curl -X GET "http://localhost:8080/api/v1/appointments/past?days_back=30" \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN"
```

---

### 8. Cancelled Appointments

**Endpoint:** `GET /api/v1/appointments/cancelled`

**Description:** Mendapatkan daftar appointments yang dibatalkan.

**Authentication:** Required (Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**
Same as List All Appointments.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Cancelled appointments retrieved successfully",
  "data": {
    "total_cancelled": 15,
    "data": [
      {
        "id": 10,
        "appointment_date": "2024-01-22",
        "appointment_time": "14:00:00",
        "patient": {
          "patient_code": "P-2024-010",
          "full_name": "John Doe"
        },
        "doctor": {
          "full_name": "Dr. Jane Smith, Sp.PD"
        },
        "status": "cancelled",
        "cancellation_reason": "Patient emergency",
        "cancelled_by": "patient",
        "cancelled_at": "2024-01-20T10:00:00Z"
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
curl -X GET "http://localhost:8080/api/v1/appointments/cancelled" \
  -H "Authorization: Bearer RECEPTIONIST_JWT_TOKEN"
```

---

### 9. Deleted Appointments

**Endpoint:** `GET /api/v1/appointments/deleted`

**Description:** Mendapatkan daftar appointments yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**
Same as List All Appointments.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Deleted appointments retrieved successfully",
  "data": {
    "data": [
      {
        "id": 100,
        "appointment_date": "2024-01-15",
        "appointment_time": "10:00:00",
        "patient": {
          "patient_code": "P-2024-050",
          "full_name": "Deleted Patient"
        },
        "doctor": {
          "full_name": "Dr. Smith"
        },
        "status": "cancelled",
        "created_at": "2024-01-10T10:00:00Z",
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
curl -X GET "http://localhost:8080/api/v1/appointments/deleted" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 10. Get Appointment by ID

**Endpoint:** `GET /api/v1/appointments/:id`

**Description:** Mendapatkan detail appointment berdasarkan ID.

**Authentication:** Required (All Authenticated, with ownership check)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Appointment ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Appointment retrieved successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "patient": {
      "id": 10,
      "patient_code": "P-2024-001",
      "full_name": "Jane Doe",
      "date_of_birth": "1978-05-15",
      "age": 45,
      "gender": "female",
      "phone": "081234567890",
      "email": "janedoe@example.com"
    },
    "doctor_id": 5,
    "doctor": {
      "id": 5,
      "employee_id": "DOC-001",
      "full_name": "Dr. John Smith, Sp.JP",
      "specialization": "Cardiologist",
      "license_number": "LIC-12345-2024",
      "phone": "081234567890",
      "department": {
        "id": 1,
        "name": "Kardiologi",
        "floor_location": "Lantai 3"
      }
    },
    "appointment_date": "2024-01-25",
    "appointment_time": "10:00:00",
    "duration_minutes": 30,
    "reason": "Kontrol rutin jantung",
    "status": "scheduled",
    "notes": "Patient has history of hypertension",
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z",
    "timeline": [
      {
        "action": "created",
        "timestamp": "2024-01-19T10:00:00Z",
        "by": "patient"
      }
    ]
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Appointment not found",
  "error": "appointment not found"
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/appointments/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Staff Endpoints

### 11. Update Appointment

**Endpoint:** `PUT /api/v1/appointments/:id`

**Description:** Update appointment details (patient own atau staff).

**Authentication:** Required (Patient own, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Appointment ID (integer)

**Request Body:**

```json
{
  "appointment_date": "2024-01-26",
  "appointment_time": "11:00:00",
  "duration_minutes": 30,
  "reason": "Kontrol rutin jantung - updated",
  "notes": "Patient requested morning slot"
}
```

**Field Rules:**

- All fields optional
- Cannot update if status is 'completed' or 'in_progress'
- Must check doctor availability for new slot

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Appointment updated successfully",
  "data": {
    "id": 1,
    "appointment_date": "2024-01-26",
    "appointment_time": "11:00:00",
    "duration_minutes": 30,
    "reason": "Kontrol rutin jantung - updated",
    "status": "scheduled",
    "notes": "Patient requested morning slot",
    "updated_at": "2024-01-19T15:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/appointments/1 \
  -H "Authorization: Bearer PATIENT_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "appointment_date": "2024-01-26",
    "appointment_time": "11:00:00"
  }'
```

---

### 12. Confirm Appointment

**Endpoint:** `PATCH /api/v1/appointments/:id/confirm`

**Description:** Confirm appointment (doctor atau receptionist).

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Appointment ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Appointment confirmed successfully",
  "data": {
    "id": 1,
    "status": "confirmed",
    "confirmed_at": "2024-01-20T10:00:00Z",
    "confirmed_by": "doctor"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/appointments/1/confirm \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN"
```

**Notes:**

- Status berubah dari 'scheduled' ke 'confirmed'
- Sends confirmation notification to patient
- Cannot confirm if appointment is in the past

---

### 13. Start Appointment

**Endpoint:** `PATCH /api/v1/appointments/:id/start`

**Description:** Start appointment (doctor only, saat pasien datang).

**Authentication:** Required (Doctor)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Appointment ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Appointment started successfully",
  "data": {
    "id": 1,
    "status": "in_progress",
    "started_at": "2024-01-25T10:05:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/appointments/1/start \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN"
```

**Notes:**

- Status berubah ke 'in_progress'
- Starts timer for duration tracking
- Prevents other appointments from starting

---

### 14. Complete Appointment

**Endpoint:** `PATCH /api/v1/appointments/:id/complete`

**Description:** Complete appointment (doctor only, setelah konsultasi selesai).

**Authentication:** Required (Doctor)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Appointment ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Appointment completed successfully",
  "data": {
    "id": 1,
    "status": "completed",
    "started_at": "2024-01-25T10:05:00Z",
    "completed_at": "2024-01-25T10:28:00Z",
    "actual_duration": 23
  }
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/appointments/1/complete \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN"
```

**Notes:**

- Status berubah ke 'completed'
- Records actual duration
- Triggers medical record creation prompt

---

### 15. Cancel Appointment

**Endpoint:** `PATCH /api/v1/appointments/:id/cancel`

**Description:** Cancel appointment (patient own, doctor own, atau staff).

**Authentication:** Required (Patient own, Doctor own, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Appointment ID (integer)

**Request Body:**

```json
{
  "cancellation_reason": "Patient emergency",
  "notes": "Will reschedule next week"
}
```

**Field Rules:**

- `cancellation_reason`: optional, text
- `notes`: optional, text

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Appointment cancelled successfully",
  "data": {
    "id": 1,
    "status": "cancelled",
    "cancellation_reason": "Patient emergency",
    "cancelled_by": "patient",
    "cancelled_at": "2024-01-20T10:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/appointments/1/cancel \
  -H "Authorization: Bearer PATIENT_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "cancellation_reason": "Patient emergency"
  }'
```

**Notes:**

- Status berubah ke 'cancelled'
- Sends cancellation notification
- Frees up the time slot
- Cannot cancel if already in_progress or completed

---

### 16. Reschedule Appointment

**Endpoint:** `PATCH /api/v1/appointments/:id/reschedule`

**Description:** Reschedule appointment ke waktu lain.

**Authentication:** Required (Patient own, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Appointment ID (integer)

**Request Body:**

```json
{
  "new_appointment_date": "2024-01-26",
  "new_appointment_time": "14:00:00",
  "reason": "Patient conflict with previous schedule"
}
```

**Field Rules:**

- `new_appointment_date`: required, date (YYYY-MM-DD)
- `new_appointment_time`: required, time (HH:MM:SS)
- `reason`: optional, text

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Appointment rescheduled successfully",
  "data": {
    "id": 1,
    "previous_date": "2024-01-25",
    "previous_time": "10:00:00",
    "appointment_date": "2024-01-26",
    "appointment_time": "14:00:00",
    "status": "scheduled",
    "rescheduled_at": "2024-01-20T10:00:00Z",
    "reschedule_reason": "Patient conflict with previous schedule"
  }
}
```

**Response Error (409 Conflict):**

```json
{
  "success": false,
  "message": "Time slot conflict",
  "error": "doctor already has an appointment at the new time",
  "suggested_slots": ["14:30:00", "15:00:00", "15:30:00"]
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/appointments/1/reschedule \
  -H "Authorization: Bearer PATIENT_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "new_appointment_date": "2024-01-26",
    "new_appointment_time": "14:00:00",
    "reason": "Patient conflict"
  }'
```

**Notes:**

- Validates new time slot availability
- Keeps original appointment record
- Sends reschedule notification to both patient and doctor
- Cannot reschedule if in_progress or completed

---

### 17. Mark as No-Show

**Endpoint:** `PATCH /api/v1/appointments/:id/no-show`

**Description:** Mark appointment sebagai no-show (patient tidak datang).

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Appointment ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Appointment marked as no-show",
  "data": {
    "id": 1,
    "status": "no_show",
    "marked_at": "2024-01-25T10:15:00Z",
    "marked_by": "receptionist"
  }
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/appointments/1/no-show \
  -H "Authorization: Bearer RECEPTIONIST_JWT_TOKEN"
```

**Notes:**

- Status berubah ke 'no_show'
- Can only mark if appointment time has passed
- May trigger penalty or notification policy
- Frees up the time slot

---

## Admin Endpoints

### 18. Soft Delete Appointment

**Endpoint:** `DELETE /api/v1/appointments/:id`

**Description:** Admin menghapus appointment (soft delete).

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Appointment ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Appointment deleted successfully",
  "data": null
}
```

**Notes:**

- Appointment yang dihapus tidak muncul di list normal
- Historical data tetap preserved
- Bisa di-restore dengan endpoint restore

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/appointments/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 19. Restore Appointment

**Endpoint:** `PATCH /api/v1/appointments/:id/restore`

**Description:** Admin me-restore appointment yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Appointment ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Appointment restored successfully",
  "data": null
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/appointments/1/restore \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

## Super Admin Endpoints

### 20. Hard Delete Appointment

**Endpoint:** `DELETE /api/v1/appointments/:id/hard-delete`

**Description:** Super Admin menghapus appointment secara permanen dari database.

**Authentication:** Required (Super Admin Only)

**Request Headers:**

```
Authorization: Bearer <super-admin-token>
```

**URL Parameters:**

- `id`: Appointment ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Appointment permanently deleted",
  "data": null
}
```

**⚠️ WARNING:**

- Data dihapus permanen dari database
- Tidak bisa di-restore
- Gunakan dengan sangat hati-hati

**⚠️ Business Rules:**

- Cannot hard delete if appointment has medical record
- Must archive historical data first

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/appointments/1/hard-delete \
  -H "Authorization: Bearer SUPER_ADMIN_JWT_TOKEN"
```

---

## Error Responses

### 400 Bad Request

```json
{
  "success": false,
  "message": "Validation error",
  "error": "appointment date cannot be in the past"
}
```

### 403 Forbidden

```json
{
  "success": false,
  "message": "Access denied",
  "error": "you can only access your own appointments"
}
```

### 404 Not Found

```json
{
  "success": false,
  "message": "Appointment not found",
  "error": "appointment not found"
}
```

### 409 Conflict

```json
{
  "success": false,
  "message": "Time slot conflict",
  "error": "doctor already has an appointment at this time"
}
```

### 422 Unprocessable Entity

```json
{
  "success": false,
  "message": "Cannot complete action",
  "error": "appointment cannot be cancelled as it is already in progress"
}
```

---

## Data Models

### Appointment Object (Full)

```json
{
  "id": 1,
  "patient_id": 10,
  "doctor_id": 5,
  "appointment_date": "2024-01-25",
  "appointment_time": "10:00:00",
  "duration_minutes": 30,
  "reason": "Kontrol rutin jantung",
  "status": "scheduled",
  "notes": "Patient has history of hypertension",
  "created_at": "2024-01-19T10:00:00Z",
  "updated_at": "2024-01-19T10:00:00Z",
  "deleted_at": null
}
```

### Appointment Status Values

```json
{
  "statuses": [
    {
      "value": "scheduled",
      "label": "Scheduled",
      "description": "Appointment booked, waiting for confirmation",
      "color": "blue"
    },
    {
      "value": "confirmed",
      "label": "Confirmed",
      "description": "Appointment confirmed by doctor/staff",
      "color": "green"
    },
    {
      "value": "in_progress",
      "label": "In Progress",
      "description": "Consultation is ongoing",
      "color": "yellow"
    },
    {
      "value": "completed",
      "label": "Completed",
      "description": "Consultation finished",
      "color": "success"
    },
    {
      "value": "cancelled",
      "label": "Cancelled",
      "description": "Appointment cancelled by patient/doctor/staff",
      "color": "red"
    },
    {
      "value": "no_show",
      "label": "No Show",
      "description": "Patient did not show up",
      "color": "orange"
    }
  ]
}
```

---

## Business Rules

1. **Time Slot Validation**: Check doctor availability before booking
2. **Past Date Prevention**: Cannot book appointments in the past
3. **Minimum Advance**: Minimum 2 hours advance booking
4. **Maximum Appointments**: Limit appointments per doctor per day
5. **Working Hours**: Validate against doctor's working hours
6. **Duration Standard**: Typically 30 or 60 minutes
7. **Status Workflow**: scheduled → confirmed → in_progress → completed
8. **Cancellation Policy**: Can cancel up to 2 hours before appointment
9. **No-Show Tracking**: Track patient no-show history
10. **Reschedule Limit**: Maximum 2 reschedules per appointment

---

## Appointment Workflow

```
1. Patient/Receptionist books appointment
   ↓
2. Status: SCHEDULED
   ↓
3. Doctor/Staff confirms appointment
   ↓
4. Status: CONFIRMED
   ↓
5. Patient arrives, doctor starts
   ↓
6. Status: IN_PROGRESS
   ↓
7. Consultation finished
   ↓
8. Status: COMPLETED
   ↓
9. Medical record created (optional)
```

Alternative flows:

- Patient cancels → Status: CANCELLED
- Patient doesn't show → Status: NO_SHOW
- Reschedule → New date/time, Status: SCHEDULED

---

## Common Use Cases

### Use Case 1: Patient Books Appointment

```bash
# 1. Search available doctors
GET /api/v1/doctors/active?specialization=Cardiologist

# 2. Check doctor availability (different endpoint)
GET /api/v1/doctors/5/availability?date=2024-01-25

# 3. Book appointment
POST /api/v1/appointments
{
  "doctor_id": 5,
  "appointment_date": "2024-01-25",
  "appointment_time": "10:00:00"
}

# 4. Receive confirmation notification
```

### Use Case 2: Doctor Daily Schedule

```bash
# 1. View today's schedule
GET /api/v1/appointments/my-schedule?date=2024-01-20

# 2. Confirm appointments
PATCH /api/v1/appointments/1/confirm

# 3. Start consultation
PATCH /api/v1/appointments/1/start

# 4. Complete consultation
PATCH /api/v1/appointments/1/complete
```

### Use Case 3: Receptionist Manages Appointments

```bash
# 1. View today's appointments
GET /api/v1/appointments/today

# 2. Book walk-in appointment
POST /api/v1/appointments

# 3. Reschedule on patient request
PATCH /api/v1/appointments/1/reschedule

# 4. Mark no-show
PATCH /api/v1/appointments/2/no-show
```

---

## Testing Examples

### Test 1: Complete Appointment Lifecycle

```bash
# 1. Book Appointment
curl -X POST http://localhost:8080/api/v1/appointments \
  -H "Authorization: Bearer PATIENT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 10,
    "doctor_id": 5,
    "appointment_date": "2024-01-25",
    "appointment_time": "10:00:00",
    "duration_minutes": 30,
    "reason": "Konsultasi"
  }'

# 2. Confirm Appointment
curl -X PATCH http://localhost:8080/api/v1/appointments/1/confirm \
  -H "Authorization: Bearer DOCTOR_TOKEN"

# 3. Start Appointment
curl -X PATCH http://localhost:8080/api/v1/appointments/1/start \
  -H "Authorization: Bearer DOCTOR_TOKEN"

# 4. Complete Appointment
curl -X PATCH http://localhost:8080/api/v1/appointments/1/complete \
  -H "Authorization: Bearer DOCTOR_TOKEN"
```

---

## Notes

- Appointment status workflow strictly enforced
- Automatic notifications for all status changes
- Integration with doctor schedule management
- Support for recurring appointments (future feature)
- Waiting list management for fully booked slots
- SMS/Email reminders before appointment
- QR code check-in system
- Real-time availability calendar
- Telemedicine appointment support
- Video consultation integration

---

## Database Model

### Table: appointments

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| patient_id | BIGINT | FOREIGN KEY (patients.id), INDEX | Reference ke pasien |
| doctor_id | BIGINT | FOREIGN KEY (doctors.id), INDEX | Reference ke dokter |
| appointment_date | DATE | NOT NULL, INDEX | Tanggal appointment |
| appointment_time | TIME | NOT NULL | Waktu appointment |
| duration_minutes | INT | DEFAULT 30 | Durasi appointment dalam menit |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'scheduled', INDEX | Status (scheduled, confirmed, started, completed, cancelled) |
| notes | TEXT | NULLABLE | Catatan/keluhan pasien |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

**Indexes:**
- Primary Key: id
- Foreign Keys: patient_id, doctor_id
- Regular Index: appointment_date, status, deleted_at

**Relationships:**
- Belongs To Patient (many-to-one)
- Belongs To Doctor (many-to-one)
- Has One Medical Record (one-to-one, after appointment)

**Notes:**
- Status flow: scheduled -> confirmed -> started -> completed
- No double booking untuk doctor pada waktu yang sama
- Appointment harus scheduled dalam doctor's working hours

---

**Last Updated:** 2024-01-19  
**API Version:** 1.0.0
