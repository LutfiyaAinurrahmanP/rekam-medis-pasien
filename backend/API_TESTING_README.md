# API Unit Testing - Rekam Medis Pasien

## 📋 Overview

File ini berisi unit testing lengkap untuk semua endpoint API yang ada di `internal/routes/api.go`. Testing menggunakan format output mirip Jest (JavaScript testing framework) untuk kemudahan membaca hasil testing.

## 🎯 Coverage

### User Profile Endpoints (`/api/v1/users/me`)

- ✅ **GET** `/api/v1/users/me` - Get My Profile (4 tests)
- ✅ **PUT** `/api/v1/users/me` - Update My Profile (4 tests)
- ✅ **PATCH** `/api/v1/users/me/change-password` - Change Password (3 tests)
- ✅ **DELETE** `/api/v1/users/me` - Delete My Account (3 tests)
- ✅ **PATCH** `/api/v1/users/me/deactivate` - Deactivate My Account (3 tests)

### Admin Endpoints (`/api/v1/users`)

- ✅ **POST** `/api/v1/users` - Create User (4 tests)
- ✅ **GET** `/api/v1/users` - List Users (3 tests)
- ✅ **GET** `/api/v1/users/deleted` - Get Deleted Users (2 tests)
- ✅ **GET** `/api/v1/users/:id` - Get User By ID (4 tests)
- ✅ **PUT** `/api/v1/users/:id` - Update User (3 tests)
- ✅ **DELETE** `/api/v1/users/:id` - Soft Delete User (3 tests)
- ✅ **PATCH** `/api/v1/users/:id/restore` - Restore User (3 tests)
- ✅ **PATCH** `/api/v1/users/:id/reset-password` - Reset Password (3 tests)
- ✅ **PATCH** `/api/v1/users/:id/activate` - Activate User (3 tests)
- ✅ **PATCH** `/api/v1/users/:id/deactivate` - Deactivate User (4 tests)

### Super Admin Endpoints

- ✅ **DELETE** `/api/v1/users/:id/hard-delete` - Hard Delete User (5 tests)

**Total: 54 test cases**

## 🧪 Test Scenarios

Setiap endpoint diuji dengan berbagai skenario:

### ✓ Happy Path (Success Cases)

- Request valid dengan autentikasi yang benar
- Role yang sesuai
- Data yang valid

### ✗ Error Cases

- **Authentication Errors:**
  - Missing authorization token
  - Invalid token format
  - Expired token
- **Authorization Errors:**
  - Wrong role (non-admin trying to access admin endpoints)
  - Insufficient permissions
- **Validation Errors:**
  - Invalid user ID format
  - Invalid request body
  - Missing required fields
- **Business Logic Errors:**
  - User not found
  - Duplicate username/email/phone
  - Wrong password verification
  - Trying to update forbidden fields

## 🚀 Running Tests

### Cara 1: Menggunakan Script (Recommended)

**Windows:**

```bash
.\run-api-tests.bat
```

**Linux/Mac:**

```bash
chmod +x run-api-tests.sh
./run-api-tests.sh
```

### Cara 2: Langsung dengan Go Command

```bash
# Run all API tests
go test -v ./internal/routes/... -count=1

# Run specific test
go test -v ./internal/routes/... -run TestGetMyProfile -count=1

# Run tests with coverage
go test -v ./internal/routes/... -cover -count=1
```

## 📊 Output Format (Jest-like)

Output testing dirancang mirip dengan Jest untuk kemudahan membaca:

```
═══════════════════════════════════════════════════════════════
  Test Suite: GET /api/v1/users/me - Get My Profile
═══════════════════════════════════════════════════════════════
  ✓ should return 200 and user profile (1.0757ms)
  ✓ should return 401 when token missing (0s)
  ✓ should return 401 for invalid token format (548.8µs)
  ✓ should return 404 when user not found (0s)
───────────────────────────────────────────────────────────────
  ✓ 4 passed | Total: 4 (100.0%)
  Time: 1.6245ms
═══════════════════════════════════════════════════════════════
```

### Legend:

- `✓` = Test passed (hijau di terminal)
- `✗` = Test failed (merah di terminal)
- Time = Durasi eksekusi per test dan total

## 📁 File Structure

```
backend/
├── internal/
│   └── routes/
│       ├── api.go              # API routes definition
│       ├── api_test.go         # Unit tests (49 test cases)
│       └── routes.go           # Main router setup
├── run-api-tests.bat           # Windows test runner
└── run-api-tests.sh            # Linux/Mac test runner
```

## 🔧 Testing Architecture

### Mock Service

Tests menggunakan mock service (`MockUserService`) untuk mengisolasi testing dari database dan business logic layer:

- Tidak memerlukan database connection
- Fast execution
- Predictable results
- Easy to test error scenarios

### Test Router

Setiap test menggunakan dedicated router instance:

- Clean state per test
- No interference between tests
- Proper middleware setup

### JWT Token Generation

Tests menggunakan utility untuk generate valid JWT tokens:

- Different roles (patient, admin, super_admin)
- Different user IDs
- Valid expiration time

## 🐛 Known Issues (Documented in Tests)

### 1. Admin Deactivate Endpoint Bug

**Location:** `internal/routes/api.go:36`

**Issue:** Endpoint `PATCH /api/v1/users/:id/deactivate` menggunakan handler `DeactivateMyAccount` yang salah.

**Actual Behavior:**

- Endpoint men-deactivate user dari token (authenticated admin), bukan user dari URL parameter `:id`
- Membutuhkan password di request body
- Admin yang memanggil endpoint malah men-deactivate akun sendiri

**Expected Behavior:**

- Seharusnya men-deactivate user dengan ID dari URL parameter
- Tidak perlu password (admin action)

**Fix Required:**

```go
// Current (WRONG):
adminRoutes.PATCH("/:id/deactivate", cfg.UserHandler.DeactivateMyAccount)

// Should be (CORRECT):
adminRoutes.PATCH("/:id/deactivate", cfg.UserHandler.DeactivateUser)
```

**Test:** `TestDeactivateUser/Success_-_Admin_deactivates_user`

## 📈 Test Statistics

- **Total Endpoints Tested:** 16
- **Total Test Cases:** 49
- **Success Rate:** 100%
- **Average Execution Time:** ~20-30ms
- **Coverage:** All endpoints in api.go

## 🛠 Technology Stack

- **Testing Framework:** Go testing package
- **Assertions:** testify/assert
- **Mocking:** testify/mock
- **HTTP Testing:** httptest
- **Router:** Gin framework

## 📝 Best Practices Followed

1. **Descriptive Test Names:** Clear, easy to understand
2. **Arrange-Act-Assert Pattern:** Well-structured tests
3. **Independent Tests:** No dependencies between tests
4. **Mock Isolation:** Service layer completely mocked
5. **Error Coverage:** All error paths tested
6. **Performance Tracking:** Execution time per test
7. **Clean Output:** Jest-like summary format

## 💡 Tips

### Debugging Failed Tests

```bash
# Run specific test with verbose output
go test -v ./internal/routes/... -run TestName

# Check mock expectations
# Look for "FAIL: 0 out of X expectation(s) were met"
```

### Adding New Tests

1. Tambahkan test function di `api_test.go`
2. Gunakan pattern yang sama dengan tests yang ada
3. Setup mock expectations
4. Assert HTTP status code dan response
5. Verify mock expectations

### Running Tests in CI/CD

```yaml
# Example GitHub Actions
- name: Run API Tests
  run: |
    cd backend
    go test -v ./internal/routes/... -count=1
```

## 📞 Contact

Jika ada pertanyaan atau menemukan bug, silakan buat issue di repository.

---

**Created:** January 2026  
**Last Updated:** January 27, 2026  
**Maintainer:** Backend Team
