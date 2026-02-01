# Doctors API Documentation

## Base URL

```
/api/v1/doctors
```

## Authentication

All endpoints require JWT token.

---

## Endpoints Overview

| Method | Endpoint                        | Access            | Description                   |
| ------ | ------------------------------- | ----------------- | ----------------------------- |
| GET    | `/doctors/me`                   | Doctor            | Get own doctor profile        |
| PUT    | `/doctors/me`                   | Doctor            | Update own profile            |
| GET    | `/doctors`                      | All Authenticated | List active doctors           |
| GET    | `/doctors/deleted`              | Admin             | List deleted doctors          |
| GET    | `/doctors/:id`                  | All Authenticated | Get doctor by ID              |
| GET    | `/doctors/specialization/:spec` | All Authenticated | Get doctors by specialization |
| POST   | `/doctors`                      | Admin             | Create doctor                 |
| PUT    | `/doctors/:id`                  | Admin             | Update doctor                 |
| PATCH  | `/doctors/:id/activate`         | Admin             | Activate doctor               |
| PATCH  | `/doctors/:id/deactivate`       | Admin             | Deactivate doctor             |
| DELETE | `/doctors/:id`                  | Admin             | Soft delete doctor            |
| PATCH  | `/doctors/:id/restore`          | Admin             | Restore deleted doctor        |
| DELETE | `/doctors/:id/hard-delete`      | Super Admin       | Permanently delete doctor     |

---

## Self-Owned Endpoints

### 1. Get My Doctor Profile

**GET** `/api/v1/doctors/me`

**Access:** Doctor (Self)

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Doctor profile retrieved successfully",
  "data": {
    "id": 1,
    "user_id": 2,
    "employee_id": "DOC001",
    "full_name": "Dr. John Smith",
    "specialization": "Cardiology",
    "license_number": "LIC123456",
    "phone": "081234567890",
    "email": "drsmith@hospital.com",
    "department_id": 1,
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z",
    "department": {
      "id": 1,
      "name": "Cardiology Department",
      "code": "CARD"
    }
  }
}
```

### 2. Update My Profile

**PUT** `/api/v1/doctors/me`

**Access:** Doctor (Self)

**Request Body:**

```json
{
  "phone": "081234567899",
  "email": "newemail@hospital.com"
}
```

**Note:** Can only update phone and email

**Response:** `200 OK`

---

## Public Endpoints (Authenticated)

### 3. List Doctors

**GET** `/api/v1/doctors`

**Access:** All Authenticated Users

**Query Parameters:**

- `page` (int)
- `page_size` (int)
- `search` (string) - Search by name, employee_id
- `specialization` (string)
- `department_id` (int)
- `is_active` (boolean)
- `sort_by` (string)
- `sort_dir` (string)

**Response:** `200 OK`

### 4. Get Doctor by ID

**GET** `/api/v1/doctors/:id`

**Access:** All Authenticated Users

**Response:** `200 OK`

### 5. Get Doctors by Specialization

**GET** `/api/v1/doctors/specialization/:spec`

**Access:** All Authenticated Users

**Example:** `/api/v1/doctors/specialization/Cardiology`

**Response:** `200 OK`

---

## Admin Endpoints

### 6. Create Doctor

**POST** `/api/v1/doctors`

**Access:** Admin, Super Admin

**Request Body:**

```json
{
  "user_id": 3,
  "employee_id": "DOC002",
  "full_name": "Dr. Jane Doe",
  "specialization": "Neurology",
  "license_number": "LIC789012",
  "phone": "081234567891",
  "email": "drjane@hospital.com",
  "department_id": 2,
  "is_active": true
}
```

**Response:** `201 Created`

### 7. Update Doctor

**PUT** `/api/v1/doctors/:id`

**Access:** Admin, Super Admin

**Response:** `200 OK`

### 8. Activate Doctor

**PATCH** `/api/v1/doctors/:id/activate`

**Access:** Admin, Super Admin

**Response:** `200 OK`

### 9. Deactivate Doctor

**PATCH** `/api/v1/doctors/:id/deactivate`

**Access:** Admin, Super Admin

**Response:** `200 OK`

### 10. List Deleted Doctors

**GET** `/api/v1/doctors/deleted`

**Access:** Admin, Super Admin

**Response:** `200 OK`

### 11. Soft Delete Doctor

**DELETE** `/api/v1/doctors/:id`

**Access:** Admin, Super Admin

**Response:** `200 OK`

### 12. Restore Doctor

**PATCH** `/api/v1/doctors/:id/restore`

**Access:** Admin, Super Admin

**Response:** `200 OK`

---

## Super Admin Endpoints

### 13. Hard Delete Doctor

**DELETE** `/api/v1/doctors/:id/hard-delete`

**Access:** Super Admin Only

**Response:** `200 OK`

---

## Common Specializations

- Cardiology (Jantung)
- Neurology (Saraf)
- Pediatrics (Anak)
- Orthopedics (Tulang)
- Dermatology (Kulit)
- Ophthalmology (Mata)
- ENT (THT)
- Internal Medicine (Penyakit Dalam)
- Surgery (Bedah)
- Obstetrics & Gynecology (Kandungan)
