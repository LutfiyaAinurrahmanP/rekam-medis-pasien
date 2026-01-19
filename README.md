# Sistem Rekam Medis Pasien - RESTful API

Sistem Rekam Medis Pasien berbasis Golang dengan Gin Framework, PostgreSQL, dan Docker.

## 🚀 Features

- ✅ RESTful API Architecture
- ✅ Clean Architecture (Repository, Service, Handler Pattern)
- ✅ JWT Authentication & Authorization
- ✅ Role-Based Access Control (RBAC)
- ✅ GORM ORM with PostgreSQL
- ✅ Request Validation
- ✅ Pagination, Search & Filtering
- ✅ Soft Delete & Hard Delete
- ✅ Graceful Shutdown
- ✅ Docker & Docker Compose Support
- ✅ Environment Configuration with Viper
- ✅ Password Hashing (bcrypt)
- ✅ CORS Middleware

## 📋 User Roles

1. **Patient** - Pasien
2. **Doctor** - Dokter
3. **Receptionist** - Resepsionis
4. **Admin** - Administrator
5. **Super Admin** - Super Administrator

## 🛠️ Tech Stack

- **Language**: Go 1.21
- **Framework**: Gin
- **Database**: PostgreSQL 15
- **ORM**: GORM
- **Authentication**: JWT
- **Configuration**: Viper
- **Validation**: go-playground/validator
- **Password Hashing**: bcrypt
- **Containerization**: Docker & Docker Compose

## 📁 Project Structure

```
sirekam-medis-pasien/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go                 # Application entry point
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go               # Configuration management
│   │   ├── models/
│   │   │   └── user.go                 # Database models
│   │   ├── repository/
│   │   │   └── user_repository.go      # Data access layer
│   │   ├── service/
│   │   │   └── user_service.go         # Business logic layer
│   │   ├── handler/
│   │   │   └── user_handler.go         # HTTP handlers
│   │   ├── middleware/
│   │   │   ├── auth_middleware.go      # JWT authentication
│   │   │   └── cors_middleware.go      # CORS handling
│   │   ├── dto/
│   │   │   └── user_dto.go             # Data transfer objects
│   │   ├── utils/
│   │   │   ├── jwt. go                  # JWT utilities
│   │   │   ├── password. go             # Password hashing
│   │   │   ├── response.go             # Response helpers
│   │   │   └── validator.go            # Validation helpers
│   │   └── database/
│   │       ├── database.go             # Database connection
│   │       └── migration. go            # Database migrations
│   ├── . env. example                    # Environment variables example
│   ├── . env                            # Environment variables
│   ├── Dockerfile                      # Docker configuration
│   ├── docker-compose. yml              # Docker Compose configuration
│   ├── go.mod                          # Go modules
│   └── go.sum                          # Go dependencies checksum
└── README.md                           # Project documentation
```

## 🔧 Installation & Setup

### Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose
- PostgreSQL 15 (if running without Docker)

### 1. Clone Repository

```bash
git clone https://github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien.git
cd sirekam-medis-pasien
```

### 2. Setup Environment Variables

```bash
cp .env.example . env
```

Edit `.env` file and configure your settings:

```env
# Application
APP_ENV=development
APP_NAME=Sirekam Medis API
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=sirekam_medis
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Jakarta

# JWT
JWT_SECRET=your-super-secret-key-change-this-in-production
JWT_EXPIRED_TIME=24h

# Pagination
DEFAULT_PAGE_SIZE=10
MAX_PAGE_SIZE=100
```

### 3. Run with Docker Compose (Recommended)

```bash
docker-compose up -d
```

The API will be available at `http://localhost:8080`

### 4. Run without Docker

#### Install Dependencies

```bash
go mod download
```

#### Run PostgreSQL

```bash
# Install and start PostgreSQL
# Create database
createdb sirekam_medis
```

#### Run Application

```bash
go run cmd/api/main.go
```

## 📚 API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

---

## 🔐 Authentication Endpoints

### 1. Register User

**POST** `/api/v1/auth/register`

Register a new user (default role: patient)

**Request Body:**

```json
{
  "username": "johndoe",
  "email": "johndoe@example.com",
  "phone": "081234567890",
  "password": "password123",
  "role": "patient"
}
```

**Response:** `201 Created`

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "username": "johndoe",
    "email": "johndoe@example. com",
    "phone": "081234567890",
    "role": "patient",
    "is_active": true,
    "created_at": "2024-01-19T10:00:00Z",
    "updated_at": "2024-01-19T10:00:00Z"
  }
}
```

### 2. Login

**POST** `/api/v1/auth/login`

Authenticate user and get JWT token

**Request Body:**

```json
{
  "username_or_email": "johndoe",
  "password": "password123"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-20T10:00:00Z",
    "user": {
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
}
```

---

## 👥 User Management Endpoints

### 3. Get Profile

**GET** `/api/v1/users/profile`

Get current authenticated user's profile

**Headers:**

```
Authorization: Bearer <token>
```

**Response:** `200 OK`

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

### 4. Create User

**POST** `/api/v1/users`

Create a new user (Admin/Super Admin only)

**Headers:**

```
Authorization: Bearer <token>
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

**Response:** `201 Created`

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
    "created_at": "2024-01-19T10:05:00Z",
    "updated_at": "2024-01-19T10:05:00Z"
  }
}
```

### 5. List Users

**GET** `/api/v1/users`

Get list of users with pagination and filters (Admin/Super Admin only)

**Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

| Parameter | Type    | Default    | Description                              |
| --------- | ------- | ---------- | ---------------------------------------- |
| page      | integer | 1          | Page number                              |
| page_size | integer | 10         | Items per page (max: 100)                |
| search    | string  | -          | Search by username, email, or phone      |
| role      | string  | -          | Filter by role                           |
| is_active | boolean | -          | Filter by active status                  |
| sort_by   | string  | created_at | Sort field (created_at, username, email) |
| sort_dir  | string  | desc       | Sort direction (asc, desc)               |

**Example Request:**

```
GET /api/v1/users?page=1&page_size=10&search=john&role=patient&is_active=true&sort_by=created_at&sort_dir=desc
```

**Response:** `200 OK`

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

### 6. Get User by ID

**GET** `/api/v1/users/: id`

Get user details by ID

**Headers:**

```
Authorization: Bearer <token>
```

**Response:** `200 OK`

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

### 7. Update User

**PUT** `/api/v1/users/:id`

Update user information (Admin/Super Admin only)

**Headers:**

```
Authorization: Bearer <token>
```

**Request Body:** (all fields are optional)

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

**Response:** `200 OK`

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
    "updated_at": "2024-01-19T11:00:00Z"
  }
}
```

### 8. Change Password

**PATCH** `/api/v1/users/:id/change-password`

Change user password

**Headers:**

```
Authorization: Bearer <token>
```

**Request Body:**

```json
{
  "old_password": "password123",
  "new_password": "newpassword456"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Password changed successfully",
  "data": null
}
```

### 9. Soft Delete User

**DELETE** `/api/v1/users/:id`

Soft delete user (mark as deleted) (Admin/Super Admin only)

**Headers:**

```
Authorization: Bearer <token>
```

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "User deleted successfully",
  "data": null
}
```

### 10. Hard Delete User

**DELETE** `/api/v1/users/:id/hard-delete`

Permanently delete user from database (Super Admin only)

**Headers:**

```
Authorization: Bearer <token>
```

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "User permanently deleted",
  "data": null
}
```

### 11. Restore User

**PATCH** `/api/v1/users/:id/restore`

Restore a soft-deleted user (Admin/Super Admin only)

**Headers:**

```
Authorization: Bearer <token>
```

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "User restored successfully",
  "data": null
}
```

---

## 🔒 Authorization Matrix

| Endpoint        | Patient  | Doctor   | Receptionist | Admin | Super Admin |
| --------------- | -------- | -------- | ------------ | ----- | ----------- |
| Register        | ✅       | ✅       | ✅           | ✅    | ✅          |
| Login           | ✅       | ✅       | ✅           | ✅    | ✅          |
| Get Profile     | ✅       | ✅       | ✅           | ✅    | ✅          |
| Create User     | ❌       | ❌       | ❌           | ✅    | ✅          |
| List Users      | ❌       | ❌       | ❌           | ✅    | ✅          |
| Get User by ID  | ✅       | ✅       | ✅           | ✅    | ✅          |
| Update User     | ❌       | ❌       | ❌           | ✅    | ✅          |
| Change Password | ✅ (own) | ✅ (own) | ✅ (own)     | ✅    | ✅          |
| Soft Delete     | ❌       | ❌       | ❌           | ✅    | ✅          |
| Hard Delete     | ❌       | ❌       | ❌           | ❌    | ✅          |
| Restore User    | ❌       | ❌       | ❌           | ✅    | ✅          |

---

## 🧪 Testing with cURL

### Register User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "johndoe@example.com",
    "phone": "081234567890",
    "password": "password123",
    "role": "patient"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type:  application/json" \
  -d '{
    "username_or_email": "johndoe",
    "password": "password123"
  }'
```

### Get Profile

```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### List Users

```bash
curl -X GET "http://localhost:8080/api/v1/users? page=1&page_size=10&search=john" \
  -H "Authorization:  Bearer YOUR_JWT_TOKEN"
```

### Update User

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe_updated",
    "is_active": true
  }'
```

---

## 🐳 Docker Commands

### Build and Run

```bash
# Build and start containers
docker-compose up -d

# View logs
docker-compose logs -f

# Stop containers
docker-compose down

# Rebuild containers
docker-compose up -d --build

# Remove volumes (caution:  deletes database data)
docker-compose down -v
```

### Database Management

```bash
# Access PostgreSQL container
docker exec -it sirekam_postgres psql -U postgres -d sirekam_medis

# Backup database
docker exec sirekam_postgres pg_dump -U postgres sirekam_medis > backup. sql

# Restore database
docker exec -i sirekam_postgres psql -U postgres sirekam_medis < backup. sql
```

---

## 🔍 Error Responses

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
  "message": "Access denied:  insufficient permissions",
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
  "error": "error details"
}
```

---

## 📊 Database Schema

### Users Table

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(15) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'patient',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

---

## 🚀 Production Deployment

### Environment Variables for Production

```env
APP_ENV=production
JWT_SECRET=use-a-very-strong-secret-key-here-min-32-chars
DB_SSLMODE=require
```

### Security Recommendations

1. **Change default credentials**
2. **Use strong JWT secret** (minimum 32 characters)
3. **Enable SSL/TLS** for database connections
4. **Use HTTPS** in production
5. **Enable rate limiting**
6. **Implement request logging**
7. **Use environment-specific configs**
8. **Regular security updates**

---

## 🛠️ Development

### Install Dependencies

```bash
go mod download
```

### Run Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
golangci-lint run
```

---

## 📝 License

MIT License

---

## 👨‍💻 Author

**Lutfiya Ainurrahman P**

---

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📞 Support

For support, email: support@sirekammedis.com

---

## 🔄 Changelog

### Version 1.0.0 (2024-01-19)

- ✅ Initial release
- ✅ User management CRUD operations
- ✅ JWT authentication
- ✅ Role-based access control
- ✅ Pagination and search
- ✅ Soft delete and hard delete
- ✅ Docker support
- ✅ Graceful shutdown

---

## 🎯 Roadmap

- [ ] Patient management
- [ ] Doctor management
- [ ] Appointment scheduling
- [ ] Medical records
- [ ] Prescription management
- [ ] Lab test results
- [ ] Billing system
- [ ] Reports and analytics
- [ ] Email notifications
- [ ] File upload (medical documents)
- [ ] API rate limiting
- [ ] API documentation (Swagger)
- [ ] Unit tests
- [ ] Integration tests
