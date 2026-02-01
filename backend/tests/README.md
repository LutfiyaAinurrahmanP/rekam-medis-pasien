# Tests Directory

Folder ini berisi semua test files untuk backend API.

## Structure

```
tests/
└── users_test.go           # User API integration tests (65 test cases)
```

## Running Tests

### Run All Tests in This Directory
```powershell
cd backend
go test ./tests -v
```

### Run Specific Test
```powershell
cd backend
go test ./tests -run "^Test_GetMyProfile_Success$" -v
```

### Run with Makefile (Jest-like Output)
```powershell
cd backend
make test-users
```

## Test Configuration

Tests menggunakan:
- **Real Database**: PostgreSQL (sirekam_test)
- **Auto-cleanup**: Yes (before each test)
- **Package**: `tests` (bukan `routes`)
- **Router Setup**: Menggunakan exported functions dari `internal/routes`

## Test Helper Functions

Located in `users_test.go`:
- `setupRealTestRouter()` - Initialize Gin router with real dependencies
- `cleanupRealTestDB()` - Clear all test data
- `createRealUser()` - Create user and return JWT token
- `performRealRequest()` - Execute HTTP request

## Important Notes

1. **Package Name**: Tests di folder ini menggunakan `package tests`
2. **Imports**: Perlu import `internal/routes` untuk mengakses setup functions
3. **Exported Functions**: Routes package meng-export `SetupAPIRouter` dan `SetupAuthRouter`
4. **Database**: Test database otomatis dibuat jika belum ada

## Adding New Tests

Ketika menambahkan test file baru:

1. Gunakan `package tests`
2. Import yang diperlukan dari `internal/`
3. Gunakan `routes.SetupAPIRouter()` untuk setup router
4. Gunakan helper functions yang sudah ada
5. Update Makefile jika perlu command baru

## Example Test Structure

```go
package tests

import (
    "testing"
    "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/routes"
    // ... other imports
)

func Test_YourFeature_Success(t *testing.T) {
    cleanupRealTestDB()
    
    // Setup
    _, token, _ := createRealUser(...)
    
    // Execute
    w := performRealRequest("GET", "/api/v1/endpoint", token, nil)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## Test Statistics

- **Total User API Tests**: 65
- **Currently Passing**: 50
- **Currently Failing**: 15
- **Success Rate**: 77%

## CI/CD Integration

Tests can be run in CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run Tests
  run: |
    cd backend
    make test-users
```

---

**Last Updated**: February 2026
