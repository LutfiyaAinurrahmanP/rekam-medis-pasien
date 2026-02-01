# Users API Documentation

## Overview

API untuk manajemen data users (pengguna) dalam sistem rekam medis. Mendukung 5 role: `patient`, `doctor`, `receptionist`, `admin`, dan `super_admin`.

**Base URL:** `/api/v1/users`

---

## Table of Contents

- [Authentication](#authentication)
- [User Roles](#user-roles)
- [Self-Owned Endpoints](#self-owned-endpoints)
- [Admin Endpoints](#admin-endpoints)
- [Super Admin Endpoints](#super-admin-endpoints)
- [Request & Response Examples](#request--response-examples)
- [Error Responses](#error-responses)

---

## Authentication

Semua endpoints (kecuali auth) memerlukan JWT token di header:

```
Authorization: Bearer <your-jwt-token>
```

---

## User Roles

| Role           | Description         | Permissions                                      |
| -------------- | ------------------- | ------------------------------------------------ |
| `patient`      | Pasien              | View & update own profile                        |
| `doctor`       | Dokter              | View & update own profile + medical functions    |
| `receptionist` | Resepsionis         | View & update own profile + patient registration |
| `admin`        | Administrator       | Full access except hard delete                   |
| `super_admin`  | Super Administrator | Full access including hard delete                |

---

## Endpoints Summary

### Public Endpoints

| Method | Endpoint         | Description       |
| ------ | ---------------- | ----------------- |
| POST   | `/auth/register` | Register new user |
| POST   | `/auth/login`    | Login user        |

### Self-Owned Endpoints (`/me`)

| Method | Endpoint                    | Description           | Auth |
| ------ | --------------------------- | --------------------- | ---- |
| GET    | `/users/me`                 | Get my profile        | ✅   |
| PUT    | `/users/me`                 | Update my profile     | ✅   |
| PATCH  | `/users/me/change-password` | Change my password    | ✅   |
| DELETE | `/users/me`                 | Delete my account     | ✅   |
| PATCH  | `/users/me/deactivate`      | Deactivate my account | ✅   |

### Admin Endpoints

| Method | Endpoint                    | Description          | Role Required      |
| ------ | --------------------------- | -------------------- | ------------------ |
| POST   | `/users`                    | Create user          | Admin, Super Admin |
| GET    | `/users`                    | List active users    | Admin, Super Admin |
| GET    | `/users/deleted`            | List deleted users   | Admin, Super Admin |
| GET    | `/users/:id`                | Get user by ID       | Admin, Super Admin |
| PUT    | `/users/:id`                | Update user          | Admin, Super Admin |
| DELETE | `/users/:id`                | Soft delete user     | Admin, Super Admin |
| PATCH  | `/users/:id/restore`        | Restore deleted user | Admin, Super Admin |
| PATCH  | `/users/:id/reset-password` | Reset user password  | Admin, Super Admin |
| PATCH  | `/users/:id/activate`       | Activate user        | Admin, Super Admin |
| PATCH  | `/users/:id/deactivate`     | Deactivate user      | Admin, Super Admin |

### Super Admin Endpoints

| Method | Endpoint                 | Description             | Role Required |
| ------ | ------------------------ | ----------------------- | ------------- |
| DELETE | `/users/:id/hard-delete` | Permanently delete user | Super Admin   |

---

## Self-Owned Endpoints

### 1. Get My Profile

**Endpoint:** `GET /api/v1/users/me`

**Description:** Mendapatkan profil user yang sedang login.

**Authentication:** Required

**Request Headers:**

```
Authorization: Bearer <token>
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "username": "johndoe",
    "email": "johndoe@example.com",
    "phone": "081234567890",
    "role": "patient",
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 2. Update My Profile

**Endpoint:** `PUT /api/v1/users/me`

**Description:** Update profil user yang sedang login.

**Authentication:** Required

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "username": "johndoe_updated",
  "email": "johndoe_new@example.com",
  "phone": "081234567899",
  "password": "newpassword123"
}
```

**Field Rules:**

- `username`: optional, min 3, max 50 characters, unique
- `email`: optional, valid email format, unique
- `phone`: optional, min 10, max 15 characters, unique
- `password`: optional, min 8 characters
- ❌ **Cannot update:** `role`, `is_active`

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "id": 1,
    "username": "johndoe_updated",
    "email": "johndoe_new@example.com",
    "phone": "081234567899",
    "role": "patient",
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T11:30:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe_updated",
    "email": "johndoe_new@example.com"
  }'
```

---

### 3. Change My Password

**Endpoint:** `PATCH /api/v1/users/me/change-password`

**Description:** Mengubah password user yang sedang login.

**Authentication:** Required

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "old_password": "password123",
  "new_password": "newpassword456"
}
```

**Field Rules:**

- `old_password`: required, must match current password
- `new_password`: required, min 8 characters

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Password changed successfully",
  "data": null
}
```

**Response Error (400 Bad Request):**

```json
{
  "success": false,
  "message": "Failed to change password",
  "error": "old password is incorrect"
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/users/me/change-password \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "password123",
    "new_password": "newpassword456"
  }'
```

---

### 4. Delete My Account

**Endpoint:** `DELETE /api/v1/users/me`

**Description:** Menghapus akun sendiri (soft delete). Akun dapat di-restore oleh admin.

**Authentication:** Required

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "password": "password123",
  "reason": "No longer need the account"
}
```

**Field Rules:**

- `password`: required, must match current password
- `reason`: optional, alasan penghapusan akun

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Account deleted successfully",
  "data": null
}
```

**Notes:**

- ⚠️ Setelah dihapus, user tidak bisa login
- ✅ Data tidak dihapus permanen
- ✅ Admin dapat restore akun

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "password": "password123",
    "reason": "Testing account deletion"
  }'
```

---

### 5. Deactivate My Account

**Endpoint:** `PATCH /api/v1/users/me/deactivate`

**Description:** Menonaktifkan akun sendiri sementara.

**Authentication:** Required

**Request Headers:**

```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "password": "password123",
  "reason": "Taking a break"
}
```

**Field Rules:**

- `password`: required, must match current password
- `reason`: optional, alasan deaktivasi

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Account deactivated successfully. You can reactivate by contacting admin.",
  "data": null
}
```

**Notes:**

- ⚠️ Akun tidak bisa login setelah dinonaktifkan
- ✅ Data tetap tersimpan
- ✅ Admin dapat mengaktifkan kembali

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/users/me/deactivate \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "password": "password123",
    "reason": "Taking a break"
  }'
```

---

## Admin Endpoints

### 6. Create User

**Endpoint:** `POST /api/v1/users`

**Description:** Admin membuat user baru.

**Authentication:** Required (Admin/Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**Request Body:**

```json
{
  "username": "drsmith",
  "email": "drsmith@example.com",
  "phone": "081234567891",
  "password": "password123",
  "role": "doctor",
  "is_active": true
}
```

**Field Rules:**

- `username`: required, min 3, max 50, unique
- `email`: required, valid email, unique
- `phone`: required, min 10, max 15, unique
- `password`: required, min 8 characters
- `role`: required, one of: patient, doctor, receptionist, admin, super_admin
- `is_active`: optional, boolean (default: true)

**Response Success (201 Created):**

```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 2,
    "username": "drsmith",
    "email": "drsmith@example.com",
    "phone": "081234567891",
    "role": "doctor",
    "is_active": true,
    "created_at": "2024-01-19T12:00:00Z",
    "updated_at": "2024-01-19T12:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "drsmith",
    "email": "drsmith@example.com",
    "phone": "081234567891",
    "password": "password123",
    "role": "doctor",
    "is_active": true
  }'
```

---

### 7. List Users

**Endpoint:** `GET /api/v1/users`

**Description:** Mendapatkan daftar user aktif dengan pagination, search, dan filter.

**Authentication:** Required (Admin/Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**Query Parameters:**

| Parameter   | Type    | Default    | Description                                                        |
| ----------- | ------- | ---------- | ------------------------------------------------------------------ |
| `page`      | integer | 1          | Halaman                                                            |
| `page_size` | integer | 10         | Jumlah data per halaman (max: 100)                                 |
| `search`    | string  | -          | Cari berdasarkan username, email, atau phone                       |
| `role`      | string  | -          | Filter by role (patient, doctor, receptionist, admin, super_admin) |
| `is_active` | boolean | -          | Filter by status aktif                                             |
| `sort_by`   | string  | created_at | Field untuk sorting (created_at, username, email)                  |
| `sort_dir`  | string  | desc       | Arah sorting (asc, desc)                                           |

**Example Request:**

```
GET /api/v1/users?page=1&page_size=10&search=john&role=patient&is_active=true&sort_by=created_at&sort_dir=desc
```

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": {
    "data": [
      {
        "id": 1,
        "username": "johndoe",
        "email": "johndoe@example.com",
        "phone": "081234567890",
        "role": "patient",
        "is_active": true,
        "created_at": "2024-01-19T10:00:00Z",
        "updated_at": "2024-01-19T10:00:00Z"
      },
      {
        "id": 3,
        "username": "johnsmith",
        "email": "johnsmith@example.com",
        "phone": "081234567892",
        "role": "patient",
        "is_active": true,
        "created_at": "2024-01-19T11:00:00Z",
        "updated_at": "2024-01-19T11:00:00Z"
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
curl -X GET "http://localhost:8080/api/v1/users?page=1&page_size=10&search=john&role=patient" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 8. List Deleted Users

**Endpoint:** `GET /api/v1/users/deleted`

**Description:** Mendapatkan daftar user yang sudah di-soft delete.

**Authentication:** Required (Admin/Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**Query Parameters:**
Same as List Users endpoint.

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Deleted users retrieved successfully",
  "data": {
    "data": [
      {
        "id": 5,
        "username": "deleteduser",
        "email": "deleted@example.com",
        "phone": "081234567895",
        "role": "patient",
        "is_active": false,
        "created_at": "2024-01-18T10:00:00Z",
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
curl -X GET "http://localhost:8080/api/v1/users/deleted?page=1&page_size=10" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 9. Get User by ID

**Endpoint:** `GET /api/v1/users/:id`

**Description:** Mendapatkan detail user berdasarkan ID.

**Authentication:** Required (Admin/Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: User ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "username": "johndoe",
    "email": "johndoe@example.com",
    "phone": "081234567890",
    "role": "patient",
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z"
  }
}
```

**Response Error (404 Not Found):**

```json
{
  "success": false,
  "message": "User not found",
  "error": "user not found"
}
```

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 10. Update User

**Endpoint:** `PUT /api/v1/users/:id`

**Description:** Admin mengupdate data user berdasarkan ID.

**Authentication:** Required (Admin/Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: User ID (integer)

**Request Body:**

```json
{
  "username": "johndoe_updated",
  "email": "johndoe_new@example.com",
  "phone": "081234567899",
  "password": "newpassword123",
  "role": "doctor",
  "is_active": false
}
```

**Field Rules:**

- All fields optional
- Admin dapat mengubah `role` dan `is_active`
- Validasi sama seperti create

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "username": "johndoe_updated",
    "email": "johndoe_new@example.com",
    "phone": "081234567899",
    "role": "doctor",
    "is_active": false,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T13:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe_updated",
    "role": "doctor"
  }'
```

---

### 11. Soft Delete User

**Endpoint:** `DELETE /api/v1/users/:id`

**Description:** Admin menghapus user (soft delete). Data tetap ada dan bisa di-restore.

**Authentication:** Required (Admin/Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: User ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "User deleted successfully",
  "data": null
}
```

**Notes:**

- User yang dihapus tidak bisa login
- Data tidak dihapus dari database
- Bisa di-restore dengan endpoint restore

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 12. Restore User

**Endpoint:** `PATCH /api/v1/users/:id/restore`

**Description:** Admin me-restore user yang sudah di-soft delete.

**Authentication:** Required (Admin/Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: User ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "User restored successfully",
  "data": null
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/users/1/restore \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 13. Reset Password

**Endpoint:** `PATCH /api/v1/users/:id/reset-password`

**Description:** Admin mereset password user tanpa perlu old password.

**Authentication:** Required (Admin/Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
Content-Type: application/json
```

**URL Parameters:**

- `id`: User ID (integer)

**Request Body:**

```json
{
  "new_password": "newpassword123"
}
```

**Field Rules:**

- `new_password`: required, min 8 characters

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "Password reset successfully",
  "data": null
}
```

**Notes:**

- ⚠️ Tidak perlu old password
- ⚠️ User akan logout otomatis
- ✅ Gunakan untuk forgot password

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/users/1/reset-password \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "new_password": "newpassword123"
  }'
```

---

### 14. Activate User

**Endpoint:** `PATCH /api/v1/users/:id/activate`

**Description:** Admin mengaktifkan user yang dinonaktifkan.

**Authentication:** Required (Admin/Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: User ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "User activated successfully",
  "data": null
}
```

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/users/1/activate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

### 15. Deactivate User

**Endpoint:** `PATCH /api/v1/users/:id/deactivate`

**Description:** Admin menonaktifkan user.

**Authentication:** Required (Admin/Super Admin)

**Request Headers:**

```
Authorization: Bearer <admin-token>
```

**URL Parameters:**

- `id`: User ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "User deactivated successfully",
  "data": null
}
```

**Notes:**

- User tidak bisa login setelah dinonaktifkan
- Data tetap ada
- Bisa diaktifkan kembali

**cURL Example:**

```bash
curl -X PATCH http://localhost:8080/api/v1/users/1/deactivate \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

## Super Admin Endpoints

### 16. Hard Delete User

**Endpoint:** `DELETE /api/v1/users/:id/hard-delete`

**Description:** Super Admin menghapus user secara permanen dari database.

**Authentication:** Required (Super Admin Only)

**Request Headers:**

```
Authorization: Bearer <super-admin-token>
```

**URL Parameters:**

- `id`: User ID (integer)

**Response Success (200 OK):**

```json
{
  "success": true,
  "message": "User permanently deleted",
  "data": null
}
```

**⚠️ WARNING:**

- Data dihapus permanen dari database
- Tidak bisa di-restore
- Gunakan dengan hati-hati

**cURL Example:**

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1/hard-delete \
  -H "Authorization: Bearer SUPER_ADMIN_JWT_TOKEN"
```

---

## Error Responses

### 400 Bad Request

```json
{
  "success": false,
  "message": "Validation error",
  "error": "Username already exists"
}
```

### 401 Unauthorized

```json
{
  "success": false,
  "message": "Authorization header is required",
  "error": null
}
```

### 403 Forbidden

```json
{
  "success": false,
  "message": "Access denied: insufficient permissions",
  "error": null
}
```

### 404 Not Found

```json
{
  "success": false,
  "message": "User not found",
  "error": "user not found"
}
```

### 500 Internal Server Error

```json
{
  "success": false,
  "message": "Internal server error",
  "error": "database connection failed"
}
```

---

## Data Models

### User Object

```json
{
  "id": 1,
  "username": "johndoe",
  "email": "johndoe@example.com",
  "phone": "081234567890",
  "role": "patient",
  "is_active": true,
  "created_at": "2024-01-19T10:00:00Z",
  "updated_at": "2024-01-19T10:00:00Z",
  "deleted_at": null
}
```

### Pagination Meta

```json
{
  "page": 1,
  "page_size": 10,
  "total_items": 50,
  "total_pages": 5
}
```

---

## Business Rules

1. **Username Uniqueness**: Username harus unik di seluruh sistem
2. **Email Uniqueness**: Email harus unik di seluruh sistem
3. **Phone Uniqueness**: Phone harus unik di seluruh sistem
4. **Password Security**: Minimum 8 karakter, di-hash dengan bcrypt
5. **Role Restriction**: User biasa tidak bisa mengubah role sendiri
6. **Soft Delete**: User yang di-soft delete bisa di-restore
7. **Hard Delete**: Hanya Super Admin yang bisa hard delete
8. **Deactivation**: User yang dinonaktifkan tidak bisa login
9. **Activation**: Hanya Admin yang bisa mengaktifkan user

---

## Testing Examples

### Test 1: Register and Login Flow

```bash
# 1. Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","phone":"081234567890","password":"password123","role":"patient"}'

# 2. Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username_or_email":"testuser","password":"password123"}'

# 3. Get Profile
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Test 2: Admin User Management

```bash
# 1. Create User (as Admin)
curl -X POST http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"username":"newuser","email":"new@example.com","phone":"081234567891","password":"password123","role":"patient"}'

# 2. List Users
curl -X GET "http://localhost:8080/api/v1/users?page=1&page_size=10" \
  -H "Authorization: Bearer ADMIN_TOKEN"

# 3. Update User
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role":"doctor"}'

# 4. Soft Delete
curl -X DELETE http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer ADMIN_TOKEN"

# 5. Restore
curl -X PATCH http://localhost:8080/api/v1/users/1/restore \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

---

## Notes

- Semua timestamps menggunakan format ISO 8601 (UTC)
- Pagination maksimal 100 items per page
- Soft deleted users muncul di endpoint `/users/deleted`
- Hard delete hanya bisa dilakukan oleh Super Admin
- Password selalu di-hash, tidak pernah dikembalikan dalam response

---

**Last Updated:** 2024-01-19  
**API Version:** 1.0.0
