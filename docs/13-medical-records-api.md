# Medical Records API Documentation

## Overview

API untuk manajemen data medical records (rekam medis) dalam sistem rekam medis. Medical Records adalah catatan medis lengkap dari konsultasi pasien dengan dokter, termasuk diagnosis, treatment plan, dan medical history.

**Base URL:** `/api/v1/medical-records`

---

## Table of Contents

- [Authentication](#authentication)
- [Authorization](#authorization)
- [Endpoints Summary](#endpoints-summary)
- [Self-Owned Endpoints](#self-owned-endpoints)
- [Public Endpoints](#public-endpoints)
- [Doctor Endpoints](#doctor-endpoints)
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

| Endpoint                                 | Patient (Own) | Doctor   | Receptionist | Admin | Super Admin |
| ---------------------------------------- | ------------- | -------- | ------------ | ----- | ----------- |
| GET /medical-records/my-records          | ✅            | ❌       | ❌           | ❌    | ❌          |
| GET /medical-records                     | ❌            | ✅       | ✅           | ✅    | ✅          |
| GET /medical-records/patient/:patient_id | ✅ (Own)      | ✅       | ✅           | ✅    | ✅          |
| GET /medical-records/deleted             | ❌            | ❌       | ❌           | ✅    | ✅          |
| GET /medical-records/:id                 | ✅ (Own)      | ✅       | ✅           | ✅    | ✅          |
| POST /medical-records                    | ❌            | ✅       | ❌           | ❌    | ❌          |
| PUT /medical-records/:id                 | ❌            | ✅ (Own) | ❌           | ✅    | ✅          |
| PATCH /medical-records/:id/finalize      | ❌            | ✅ (Own) | ❌           | ❌    | ❌          |
| DELETE /medical-records/:id              | ❌            | ❌       | ❌           | ✅    | ✅          |
| PATCH /medical-records/:id/restore       | ❌            | ❌       | ❌           | ✅    | ✅          |
| DELETE /medical-records/:id/hard-delete  | ❌            | ❌       | ❌           | ❌    | ✅          |

---

## Endpoints Summary

### Self-Owned Endpoints

| Method | Endpoint                      | Description                      | Auth    |
| ------ | ----------------------------- | -------------------------------- | ------- |
| GET    | `/medical-records/my-records` | Get my medical records (patient) | Patient |

### Public Endpoints

| Method | Endpoint                               | Description                   | Role Required                                           |
| ------ | -------------------------------------- | ----------------------------- | ------------------------------------------------------- |
| GET    | `/medical-records`                     | List all medical records      | Doctor, Receptionist, Admin, Super Admin                |
| GET    | `/medical-records/patient/:patient_id` | Get patient's medical records | Patient (own), Doctor, Receptionist, Admin, Super Admin |
| GET    | `/medical-records/deleted`             | List deleted medical records  | Admin, Super Admin                                      |
| GET    | `/medical-records/:id`                 | Get medical record by ID      | Patient (own), Doctor, Receptionist, Admin, Super Admin |

### Doctor Endpoints

| Method | Endpoint                        | Description             | Role Required                    |
| ------ | ------------------------------- | ----------------------- | -------------------------------- |
| POST   | `/medical-records`              | Create medical record   | Doctor                           |
| PUT    | `/medical-records/:id`          | Update medical record   | Doctor (own), Admin, Super Admin |
| PATCH  | `/medical-records/:id/finalize` | Finalize medical record | Doctor (own)                     |

### Admin Endpoints

| Method | Endpoint                       | Description                    | Role Required      |
| ------ | ------------------------------ | ------------------------------ | ------------------ |
| DELETE | `/medical-records/:id`         | Soft delete medical record     | Admin, Super Admin |
| PATCH  | `/medical-records/:id/restore` | Restore deleted medical record | Admin, Super Admin |

### Super Admin Endpoints

| Method | Endpoint                           | Description                       | Role Required |
| ------ | ---------------------------------- | --------------------------------- | ------------- |
| DELETE | `/medical-records/:id/hard-delete` | Permanently delete medical record | Super Admin   |

---

## Self-Owned Endpoints

### 1. Get My Medical Records

**Endpoint:** `GET /api/v1/medical-records/my-records`

**Description:** Patient mendapatkan daftar rekam medis milik sendiri.

**Authentication:** Required (Patient Role)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter   | Type    | Default    | Description                         |
| ----------- | ------- | ---------- | ----------------------------------- |
| `page`      | integer | 1          | Halaman                             |
| `page_size` | integer | 10         | Jumlah data per halaman (max: 100)  |
| `date_from` | date    | -          | Filter from date (YYYY-MM-DD)       |
| `date_to`   | date    | -          | Filter to date (YYYY-MM-DD)         |
| `doctor_id` | integer | -          | Filter by doctor                    |
| `sort_by`   | string  | visit_date | Sort field (visit_date, created_at) |
| `sort_dir`  | string  | desc       | Sort direction (asc, desc)          |

**Example Request:**

```
GET /api/v1/medical-records/my-records?date_from=2024-01-01&date_to=2024-01-31&sort_by=visit_date&sort_dir=desc
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "My medical records retrieved successfully",
  "data": {
    "patient_info": {
      "id": 10,
      "patient_code": "P-2024-001",
      "full_name": "Jane Doe",
      "date_of_birth": "1978-05-15",
      "age": 45
    },
    "total_records": 15,
    "last_visit": "2024-01-20",
    "data": [
      {
        "id": 1,
        "visit_date": "2024-01-20",
        "doctor": {
          "full_name": "Dr. John Smith, Sp.JP",
          "specialization": "Cardiologist",
          "department": "Kardiologi"
        },
        "chief_complaint": "Chest pain and shortness of breath",
        "diagnosis": "Stable angina pectoris",
        "treatment_plan": "Medication adjustment and lifestyle modification",
        "is_finalized": true,
        "created_at": "2024-01-20T10:00:00Z"
      },
      {
        "id": 2,
        "visit_date": "2024-01-15",
        "doctor": {
          "full_name": "Dr. Jane Smith, Sp.PD",
          "specialization": "Internal Medicine",
          "department": "Penyakit Dalam"
        },
        "chief_complaint": "Routine checkup",
        "diagnosis": "Hypertension grade 1",
        "treatment_plan": "Continue current medication",
        "is_finalized": true,
        "created_at": "2024-01-15T14:00:00Z"
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
curl -X GET "http://localhost:8080/api/v1/medical-records/my-records?page=1&page_size=10" \
  -H "Authorization: Bearer PATIENT_JWT_TOKEN"
```

**Use Case:**

- Patient melihat riwayat medis
- Patient download medical history
- Patient share dengan dokter lain

---

## Public Endpoints

### 2. List All Medical Records

**Endpoint:** `GET /api/v1/medical-records`

**Description:** Mendapatkan daftar semua medical records (staff only).

**Authentication:** Required (Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter       | Type    | Default    | Description                        |
| --------------- | ------- | ---------- | ---------------------------------- |
| `page`          | integer | 1          | Halaman                            |
| `page_size`     | integer | 10         | Jumlah data per halaman (max: 100) |
| `patient_id`    | integer | -          | Filter by patient                  |
| `doctor_id`     | integer | -          | Filter by doctor                   |
| `department_id` | integer | -          | Filter by department               |
| `date_from`     | date    | -          | Filter from date                   |
| `date_to`       | date    | -          | Filter to date                     |
| `is_finalized`  | boolean | -          | Filter by finalized status         |
| `sort_by`       | string  | visit_date | Sort field                         |
| `sort_dir`      | string  | desc       | Sort direction                     |

**Example Request:**

```
GET /api/v1/medical-records?patient_id=10&date_from=2024-01-01&date_to=2024-01-31
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medical records retrieved successfully",
  "data": {
    "filters": {
      "patient_id": 10,
      "date_from": "2024-01-01",
      "date_to": "2024-01-31"
    },
    "summary": {
      "total_records": 50,
      "finalized": 45,
      "draft": 5
    },
    "data": [
      {
        "id": 1,
        "patient": {
          "id": 10,
          "patient_code": "P-2024-001",
          "full_name": "Jane Doe",
          "age": 45
        },
        "doctor": {
          "id": 5,
          "full_name": "Dr. John Smith, Sp.JP",
          "specialization": "Cardiologist"
        },
        "visit_date": "2024-01-20",
        "chief_complaint": "Chest pain",
        "diagnosis": "Stable angina pectoris",
        "is_finalized": true,
        "created_at": "2024-01-20T10:00:00Z"
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
curl -X GET "http://localhost:8080/api/v1/medical-records?patient_id=10" \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN"
```

---

### 3. Get Patient's Medical Records

**Endpoint:** `GET /api/v1/medical-records/patient/:patient_id`

**Description:** Mendapatkan semua medical records dari patient tertentu.

**Authentication:** Required (Patient own, Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `patient_id`: Patient ID (integer)

**Query Parameters:**

| Parameter   | Type    | Description      |
| ----------- | ------- | ---------------- |
| `date_from` | date    | Filter from date |
| `date_to`   | date    | Filter to date   |
| `doctor_id` | integer | Filter by doctor |
| `page`      | integer | Page number      |
| `page_size` | integer | Items per page   |

**Example Request:**

```
GET /api/v1/medical-records/patient/10?date_from=2024-01-01&page=1&page_size=10
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Patient medical records retrieved successfully",
  "data": {
    "patient": {
      "id": 10,
      "patient_code": "P-2024-001",
      "full_name": "Jane Doe",
      "date_of_birth": "1978-05-15",
      "age": 45,
      "gender": "female",
      "blood_type": "O+",
      "allergies": "Penicillin, Seafood"
    },
    "medical_summary": {
      "total_visits": 15,
      "first_visit": "2020-05-10",
      "last_visit": "2024-01-20",
      "chronic_conditions": ["Hypertension", "Stable angina pectoris"],
      "active_medications": [
        "Amlodipine 10mg",
        "Aspirin 100mg",
        "Atorvastatin 20mg"
      ]
    },
    "records": [
      {
        "id": 1,
        "visit_date": "2024-01-20",
        "doctor": {
          "full_name": "Dr. John Smith, Sp.JP",
          "specialization": "Cardiologist"
        },
        "chief_complaint": "Chest pain and shortness of breath",
        "history_of_present_illness": "Patient complains of chest pain on exertion for the past 2 weeks. Pain is substernal, pressure-like, radiating to left arm. Associated with shortness of breath. No fever, no cough.",
        "vital_signs": {
          "blood_pressure": "140/90 mmHg",
          "heart_rate": "88 bpm",
          "temperature": "36.5°C",
          "respiratory_rate": "18/min",
          "oxygen_saturation": "98%"
        },
        "physical_examination": "General: alert, conscious, cooperative. Cardiovascular: regular rhythm, no murmur. Lungs: clear. Abdomen: soft, no tenderness.",
        "diagnosis": "Stable angina pectoris",
        "treatment_plan": "1. Adjust medication dosage\n2. Lifestyle modification\n3. Follow-up in 2 weeks\n4. ECG and cardiac enzymes if symptoms worsen",
        "is_finalized": true,
        "created_at": "2024-01-20T10:00:00Z",
        "finalized_at": "2024-01-20T10:30:00Z"
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
curl -X GET http://localhost:8080/api/v1/medical-records/patient/10 \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN"
```

**Use Case:**

- Doctor melihat riwayat medis pasien
- Continuity of care
- Medical history review

---

### 4. List Deleted Medical Records

**Endpoint:** `GET /api/v1/medical-records/deleted`

**Description:** Mendapatkan daftar medical records yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**
Same as List All Medical Records.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Deleted medical records retrieved successfully",
  "data": {
    "data": [
      {
        "id": 100,
        "patient": {
          "patient_code": "P-2024-050",
          "full_name": "Deleted Patient"
        },
        "doctor": {
          "full_name": "Dr. Smith"
        },
        "visit_date": "2024-01-10",
        "diagnosis": "Test diagnosis",
        "is_finalized": false,
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
curl -X GET "http://localhost:8080/api/v1/medical-records/deleted" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 5. Get Medical Record by ID

**Endpoint:** `GET /api/v1/medical-records/:id`

**Description:** Mendapatkan detail medical record berdasarkan ID.

**Authentication:** Required (Patient own, Doctor, Receptionist, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Medical Record ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medical record retrieved successfully",
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
      "blood_type": "O+",
      "phone": "081234567890",
      "email": "janedoe@example.com",
      "allergies": "Penicillin, Seafood"
    },
    "doctor_id": 5,
    "doctor": {
      "id": 5,
      "employee_id": "DOC-001",
      "full_name": "Dr. John Smith, Sp.JP, FIHA",
      "specialization": "Cardiologist",
      "license_number": "LIC-12345-2024",
      "department": {
        "id": 1,
        "name": "Kardiologi",
        "code": "KARDIO"
      }
    },
    "appointment_id": 1,
    "visit_date": "2024-01-20",
    "chief_complaint": "Chest pain and shortness of breath on exertion",
    "history_of_present_illness": "Patient is a 45-year-old female with a history of hypertension who presents with complaints of chest pain on exertion for the past 2 weeks. The pain is described as substernal, pressure-like, radiating to the left arm. It is associated with shortness of breath but no diaphoresis. Pain typically lasts 5-10 minutes and is relieved by rest. No fever, cough, or other constitutional symptoms. Patient has been compliant with her antihypertensive medications.",
    "past_medical_history": "Hypertension (diagnosed 5 years ago), No diabetes, No previous cardiac events",
    "family_history": "Father had myocardial infarction at age 60, Mother has hypertension",
    "social_history": "Non-smoker, Occasional alcohol consumption, Works as an office manager (sedentary lifestyle)",
    "current_medications": "Amlodipine 10mg OD, Aspirin 100mg OD",
    "allergies": "Penicillin (rash), Seafood (allergic reaction)",
    "vital_signs": {
      "blood_pressure": "140/90 mmHg",
      "heart_rate": "88 bpm",
      "temperature": "36.5°C",
      "respiratory_rate": "18/min",
      "oxygen_saturation": "98%",
      "weight": "65 kg",
      "height": "160 cm",
      "bmi": "25.4"
    },
    "physical_examination": "General: Patient is alert, conscious, cooperative, appears comfortable at rest.\nCardiovascular: Regular rate and rhythm, no murmurs, rubs or gallops. Peripheral pulses intact bilaterally.\nRespiratory: Clear to auscultation bilaterally, no wheezes or crackles.\nAbdomen: Soft, non-tender, no organomegaly.\nExtremities: No edema, no cyanosis.",
    "diagnosis": "1. Stable angina pectoris\n2. Hypertension grade 1 (controlled)",
    "icd_codes": ["I20.8", "I10"],
    "treatment_plan": "1. Medication adjustment:\n   - Continue Amlodipine 10mg OD\n   - Continue Aspirin 100mg OD\n   - Add Atorvastatin 20mg OD\n   - Add Isosorbide dinitrate 5mg SL PRN for chest pain\n\n2. Lifestyle modification:\n   - Cardiac rehabilitation program\n   - Diet: low salt, low fat\n   - Exercise: gradual increase in physical activity\n   - Stress management\n\n3. Investigations:\n   - ECG (done: shows non-specific ST-T changes)\n   - Lipid profile\n   - Cardiac enzymes if symptoms worsen\n   - Consider stress test in 2 weeks\n\n4. Follow-up in 2 weeks or sooner if symptoms worsen\n\n5. Patient education:\n   - Warning signs of heart attack\n   - When to seek emergency care\n   - Medication compliance importance",
    "notes": "Patient understands the treatment plan and agrees to lifestyle modifications. Emphasized the importance of medication compliance and follow-up. Patient instructed to return immediately if chest pain becomes more severe, prolonged, or occurs at rest.",
    "is_finalized": true,
    "created_at": "2024-01-20T10:00:00Z",
    "updated_at": "2024-01-20T10:30:00Z",
    "finalized_at": "2024-01-20T10:30:00Z",
    "deleted_at": null,
    "related_data": {
      "prescriptions_count": 1,
      "lab_tests_count": 2,
      "follow_up_appointment": {
        "id": 5,
        "appointment_date": "2024-02-03",
        "appointment_time": "10:00:00"
      }
    }
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "Medical record not found",
  "error": "medical record not found"
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/medical-records/1 \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN"
```

---

## Doctor Endpoints

### 6. Create Medical Record

**Endpoint:** `POST /api/v1/medical-records`

**Description:** Doctor membuat medical record baru setelah konsultasi.

**Authentication:** Required (Doctor Role)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "patient_id": 10,
  "appointment_id": 1,
  "visit_date": "2024-01-20",
  "chief_complaint": "Chest pain and shortness of breath on exertion",
  "history_of_present_illness": "Patient complains of chest pain on exertion for the past 2 weeks. Pain is substernal, pressure-like, radiating to left arm.",
  "past_medical_history": "Hypertension (diagnosed 5 years ago)",
  "family_history": "Father had MI at age 60",
  "social_history": "Non-smoker, occasional alcohol",
  "current_medications": "Amlodipine 10mg OD, Aspirin 100mg OD",
  "allergies": "Penicillin (rash), Seafood",
  "vital_signs": {
    "blood_pressure": "140/90 mmHg",
    "heart_rate": "88 bpm",
    "temperature": "36.5°C",
    "respiratory_rate": "18/min",
    "oxygen_saturation": "98%"
  },
  "physical_examination": "General: alert, conscious. CV: regular rhythm. Lungs: clear.",
  "diagnosis": "Stable angina pectoris",
  "icd_codes": ["I20.8"],
  "treatment_plan": "Medication adjustment and lifestyle modification",
  "notes": "Patient understands treatment plan"
}
```

**Field Rules:**

- `patient_id`: required, FK to patients table
- `appointment_id`: optional, FK to appointments table
- `visit_date`: required, date (YYYY-MM-DD)
- `chief_complaint`: required, text
- `history_of_present_illness`: optional, text
- `past_medical_history`: optional, text
- `family_history`: optional, text
- `social_history`: optional, text
- `current_medications`: optional, text
- `allergies`: optional, text
- `vital_signs`: optional, JSON object
- `physical_examination`: optional, text
- `diagnosis`: required, text
- `icd_codes`: optional, array of strings
- `treatment_plan`: required, text
- `notes`: optional, text
- `doctor_id`: auto-set from authenticated user

**Response Success (201 Created):**

```json
{
  "success": true,
  "message": "Medical record created successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "doctor_id": 5,
    "appointment_id": 1,
    "visit_date": "2024-01-20",
    "chief_complaint": "Chest pain and shortness of breath on exertion",
    "diagnosis": "Stable angina pectoris",
    "treatment_plan": "Medication adjustment and lifestyle modification",
    "is_finalized": false,
    "created_at": "2024-01-20T10:00:00Z",
    "updated_at": "2024-01-20T10:00:00Z"
  }
}
```

**Response Error (400 Bad Request):**

```json
{
  "success": false,
  "message": "Failed to create medical record",
  "error": "patient not found"
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/medical-records \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 10,
    "visit_date": "2024-01-20",
    "chief_complaint": "Chest pain",
    "diagnosis": "Stable angina pectoris",
    "treatment_plan": "Medication and lifestyle modification"
  }'
```

**Notes:**

- Medical record created as draft (is_finalized = false)
- Doctor ID automatically set from authenticated user
- Can link to appointment if created from appointment
- Vital signs stored as JSON for flexibility

---

### 7. Update Medical Record

**Endpoint:** `PUT /api/v1/medical-records/:id`

**Description:** Doctor mengupdate medical record (hanya jika belum finalized).

**Authentication:** Required (Doctor own, Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: Medical Record ID (integer)

**Request Body:**

```json
{
  "chief_complaint": "Chest pain and shortness of breath - updated",
  "history_of_present_illness": "Updated history",
  "vital_signs": {
    "blood_pressure": "135/85 mmHg",
    "heart_rate": "85 bpm"
  },
  "physical_examination": "Updated examination findings",
  "diagnosis": "Stable angina pectoris, Hypertension grade 1",
  "treatment_plan": "Updated treatment plan",
  "notes": "Additional notes"
}
```

**Field Rules:**

- All fields optional
- Cannot update if is_finalized = true (unless admin)
- Can only update own medical records (unless admin)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medical record updated successfully",
  "data": {
    "id": 1,
    "patient_id": 10,
    "doctor_id": 5,
    "visit_date": "2024-01-20",
    "chief_complaint": "Chest pain and shortness of breath - updated",
    "diagnosis": "Stable angina pectoris, Hypertension grade 1",
    "is_finalized": false,
    "created_at": "2024-01-20T10:00:00Z",
    "updated_at": "2024-01-20T10:15:00Z"
  }
}
```

**Response Error (403 Forbidden):**

```json
{
  "success": false,
  "message": "Cannot update medical record",
  "error": "medical record is already finalized"
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/medical-records/1 \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "diagnosis": "Updated diagnosis",
    "treatment_plan": "Updated treatment"
  }'
```

---

### 8. Finalize Medical Record

**Endpoint:** `PATCH /api/v1/medical-records/:id/finalize`

**Description:** Doctor finalize medical record (lock untuk editing).

**Authentication:** Required (Doctor own)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Medical Record ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medical record finalized successfully",
  "data": {
    "id": 1,
    "is_finalized": true,
    "finalized_at": "2024-01-20T10:30:00Z",
    "finalized_by": 5
  }
}
```

**Response Error (400 Bad Request):**

```json
{
  "success": false,
  "message": "Cannot finalize medical record",
  "error": "required fields are missing: diagnosis, treatment_plan"
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/medical-records/1/finalize \
  -H "Authorization: Bearer DOCTOR_JWT_TOKEN"
```

**Notes:**

- Validates that all required fields are filled
- Once finalized, cannot be edited (except by admin)
- Digital signature dapat ditambahkan di sini
- Triggers notification to patient (optional)
- Creates audit trail

**Business Rules:**

- Required fields before finalize:
  - chief_complaint
  - diagnosis
  - treatment_plan
- Cannot finalize twice
- Only doctor who created can finalize

---

## Admin Endpoints

### 9. Soft Delete Medical Record

**Endpoint:** `DELETE /api/v1/medical-records/:id`

**Description:** Admin menghapus medical record (soft delete).

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id`: Medical Record ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medical record deleted successfully",
  "data": null
}
```

**Notes:**

- Medical record yang dihapus tidak muncul di list normal
- Data tetap tersimpan untuk audit trail
- Bisa di-restore dengan endpoint restore
- Harus ada alasan kuat untuk delete (biasanya error/duplicate)

**⚠️ Business Rules:**

- Requires special permission (audit trail)
- Deletion reason must be documented
- Cannot delete if linked to billing/insurance claim

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/medical-records/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 10. Restore Medical Record

**Endpoint:** `PATCH /api/v1/medical-records/:id/restore`

**Description:** Admin me-restore medical record yang sudah di-soft delete.

**Authentication:** Required (Admin, Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: Medical Record ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medical record restored successfully",
  "data": null
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/medical-records/1/restore \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

## Super Admin Endpoints

### 11. Hard Delete Medical Record

**Endpoint:** `DELETE /api/v1/medical-records/:id/hard-delete`

**Description:** Super Admin menghapus medical record secara permanen dari database.

**Authentication:** Required (Super Admin Only)

**Request Headers:**

```
Authorization: Bearer <super-admin-token>
```

**URL Parameters:**

- `id`: Medical Record ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Medical record permanently deleted",
  "data": null
}
```

**⚠️ WARNING:**

- Data dihapus permanen dari database
- Tidak bisa di-restore
- Gunakan dengan sangat hati-hati
- Violates medical record retention policy

**⚠️ Business Rules:**

- Should NEVER be used in production
- Violates legal requirements (medical records must be kept)
- Only for test data cleanup
- Requires special authorization and documentation

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/medical-records/1/hard-delete \
  -H "Authorization: Bearer SUPER_ADMIN_JWT_TOKEN"
```

---

## Error Responses

### 400 Bad Request

```json
{
  "success": false,
  "message": "Validation error",
  "error": "chief_complaint is required"
}
```

### 403 Forbidden

```json
{
  "success": false,
  "message": "Access denied",
  "error": "you can only access your own medical records"
}
```

### 404 Not Found

```json
{
  "success": false,
  "message": "Medical record not found",
  "error": "medical record not found"
}
```

### 409 Conflict

```json
{
  "success": false,
  "message": "Cannot update medical record",
  "error": "medical record is already finalized"
}
```

---

## Data Models

### Medical Record Object (Full)

```json
{
  "id": 1,
  "patient_id": 10,
  "doctor_id": 5,
  "appointment_id": 1,
  "visit_date": "2024-01-20",
  "chief_complaint": "Chest pain and shortness of breath",
  "history_of_present_illness": "Patient complains...",
  "past_medical_history": "Hypertension",
  "family_history": "Father had MI",
  "social_history": "Non-smoker",
  "current_medications": "Amlodipine 10mg",
  "allergies": "Penicillin",
  "vital_signs": {
    "blood_pressure": "140/90 mmHg",
    "heart_rate": "88 bpm",
    "temperature": "36.5°C"
  },
  "physical_examination": "General: alert, conscious",
  "diagnosis": "Stable angina pectoris",
  "icd_codes": ["I20.8"],
  "treatment_plan": "Medication adjustment",
  "notes": "Patient understands",
  "is_finalized": true,
  "finalized_at": "2024-01-20T10:30:00Z",
  "created_at": "2024-01-20T10:00:00Z",
  "updated_at": "2024-01-20T10:30:00Z",
  "deleted_at": null
}
```

### Vital Signs Structure

```json
{
  "vital_signs": {
    "blood_pressure": "120/80 mmHg",
    "heart_rate": "75 bpm",
    "temperature": "36.5°C",
    "respiratory_rate": "16/min",
    "oxygen_saturation": "98%",
    "weight": "70 kg",
    "height": "170 cm",
    "bmi": "24.2"
  }
}
```

---

## Business Rules

1. **HIPAA Compliance**: Medical records are highly confidential
2. **Access Control**: Strict role-based access (patient own, treating doctor, authorized staff)
3. **Audit Trail**: All access and changes logged
4. **Retention Policy**: Records must be kept for minimum 7 years (Indonesia law)
5. **Finalization**: Once finalized, cannot be edited (except admin with reason)
6. **Digital Signature**: Doctor signature required for finalized records
7. **Patient Access**: Patients have right to access own records
8. **Data Integrity**: No deletion in normal operations (soft delete only)
9. **Continuity of Care**: Complete history accessible to authorized providers
10. **Quality Assurance**: Peer review and quality checks

---

## SOAP Note Format

Medical records commonly follow SOAP format:

```
S - Subjective (Patient's complaint)
  - Chief Complaint
  - History of Present Illness
  - Review of Systems

O - Objective (Doctor's observations)
  - Vital Signs
  - Physical Examination
  - Lab/Imaging Results

A - Assessment (Diagnosis)
  - Primary Diagnosis
  - Differential Diagnosis
  - ICD-10 Codes

P - Plan (Treatment)
  - Medications
  - Procedures
  - Follow-up
  - Patient Education
```

---

## Common Use Cases

### Use Case 1: Doctor Creates Medical Record After Consultation

```bash
# 1. Complete appointment
PATCH /api/v1/appointments/1/complete

# 2. Create medical record
POST /api/v1/medical-records
{
  "patient_id": 10,
  "appointment_id": 1,
  "visit_date": "2024-01-20",
  "chief_complaint": "Chest pain",
  "diagnosis": "Stable angina",
  "treatment_plan": "Medication"
}

# 3. Create prescription
POST /api/v1/prescriptions

# 4. Order lab tests
POST /api/v1/lab-tests

# 5. Finalize record
PATCH /api/v1/medical-records/1/finalize
```

### Use Case 2: Patient Views Medical History

```bash
# 1. View all my medical records
GET /api/v1/medical-records/my-records

# 2. View specific record detail
GET /api/v1/medical-records/1

# 3. Download/Print (additional endpoint)
GET /api/v1/medical-records/1/pdf
```

### Use Case 3: Doctor Reviews Patient History Before Consultation

```bash
# 1. View patient's medical records
GET /api/v1/medical-records/patient/10

# 2. View specific previous visit
GET /api/v1/medical-records/5

# 3. Check prescriptions history
GET /api/v1/prescriptions/patient/10

# 4. Check lab results
GET /api/v1/lab-tests/patient/10
```

---

## Testing Examples

### Test 1: Complete Medical Record Lifecycle

```bash
# 1. Create Medical Record
curl -X POST http://localhost:8080/api/v1/medical-records \
  -H "Authorization: Bearer DOCTOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 10,
    "visit_date": "2024-01-20",
    "chief_complaint": "Chest pain",
    "diagnosis": "Stable angina pectoris",
    "treatment_plan": "Medication adjustment"
  }'

# 2. Update Medical Record
curl -X PUT http://localhost:8080/api/v1/medical-records/1 \
  -H "Authorization: Bearer DOCTOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "diagnosis": "Stable angina pectoris, Hypertension"
  }'

# 3. Finalize Medical Record
curl -X PATCH http://localhost:8080/api/v1/medical-records/1/finalize \
  -H "Authorization: Bearer DOCTOR_TOKEN"

# 4. View Medical Record (Patient)
curl -X GET http://localhost:8080/api/v1/medical-records/1 \
  -H "Authorization: Bearer PATIENT_TOKEN"
```

---

## Notes

- Medical records are the core of the system
- HIPAA/privacy compliance critical
- Comprehensive audit trail required
- Integration with prescriptions, lab tests, imaging
- Support for templates (SOAP, specialty-specific)
- Digital signature for finalized records
- Export to PDF for printing/sharing
- Electronic Medical Record (EMR) standards
- HL7/FHIR integration capability
- Telemedicine consultation records
- Support for multimedia (images, scans)
- Version control for amendments
- Peer review and quality assurance

---

## Database Model

### Table: medical_records

| Field | Type | Constraints | Description |
| --- | --- | --- | --- |
| id | BIGINT | PRIMARY KEY, AUTO_INCREMENT | Unique identifier |
| appointment_id | BIGINT | FOREIGN KEY (appointments.id), NULLABLE, INDEX | Reference ke appointment |
| patient_id | BIGINT | FOREIGN KEY (patients.id), NOT NULL, INDEX | Reference ke pasien |
| doctor_id | BIGINT | FOREIGN KEY (doctors.id), NOT NULL, INDEX | Reference ke dokter |
| visit_date | DATE | NOT NULL | Tanggal kunjungan |
| chief_complaint | TEXT | NOT NULL | Keluhan utama pasien |
| history_of_illness | TEXT | NULLABLE | Riwayat penyakit pasien |
| physical_examination | TEXT | NULLABLE | Hasil pemeriksaan fisik |
| diagnosis | TEXT | NOT NULL | Diagnosis dokter |
| treatment_plan | TEXT | NOT NULL | Rencana treatment |
| notes | TEXT | NULLABLE | Catatan tambahan dokter |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'draft', INDEX | Status (draft, finalized, amended) |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | Waktu pembuatan |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP | Waktu update |
| deleted_at | TIMESTAMP | INDEX, NULLABLE | Soft delete timestamp |

**Indexes:**
- Primary Key: id
- Foreign Keys: patient_id, doctor_id, appointment_id
- Regular Index: visit_date, status, deleted_at

**Relationships:**
- Belongs To Patient (many-to-one)
- Belongs To Doctor (many-to-one)
- Belongs To Appointment (many-to-one, optional)
- Has Many Prescriptions (one-to-many)
- Has Many Lab Tests (one-to-many)
- Has Many Vital Signs (one-to-many)

**Notes:**
- Status: draft (masih diedit) -> finalized (selesai) -> amended (koreksi)
- Setiap record adalah catatan kunjungan satu pasien ke satu dokter
- Diagnosis dan treatment plan wajib diisi
- Dapat direferensikan dari appointment atau standalone

---

**Last Updated:** 2024-01-19  
**API Version:** 1.0.0
