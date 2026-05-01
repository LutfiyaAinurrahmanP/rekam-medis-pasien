# Make Test Commands - Quick Reference

## 🚀 Essential Commands

```bash
# Run all unit tests
make test-unit

# Run specific layer
make test-unit-handler      # 25 tests - HTTP endpoints
make test-unit-service      # 17 tests - Business logic  
make test-unit-repository   # 31 tests - Data access
```

## 📊 Coverage Reports

```bash
# Terminal coverage report
make test-unit-coverage

# HTML report (opens in browser)
make test-unit-coverage-html

# Layer-specific coverage
make test-unit-coverage-handler
make test-unit-coverage-service
make test-unit-coverage-repository
```

## 🔍 Detailed Output

```bash
# Verbose output (see all test names)
make test-unit-verbose

# Layer-specific verbose
make test-unit-handler-verbose
make test-unit-service-verbose
make test-unit-repository-verbose
```

## 🛡️ Advanced Testing

```bash
# Detect race conditions
make test-unit-race

# Stop at first failure
make test-unit-failfast

# Run with timeout (60s)
make test-unit-timeout

# Run specific test
make test-unit-run TEST=TestGetMyProfile

# List all tests
make test-unit-list
```

## 🧹 Cleanup

```bash
# Remove coverage files
make clean-coverage

# Clean all test artifacts
make test-all-cleanup
```

## 🔄 CI/CD

```bash
# Full CI suite (coverage + race detection)
make test-ci
```

## 📋 All Test Targets

| Command | Description | Tests |
|---------|-------------|-------|
| `test-unit` | All unit tests | 73 |
| `test-unit-handler` | Handler layer only | 25 |
| `test-unit-service` | Service layer only | 17 |
| `test-unit-repository` | Repository layer only | 31 |
| `test-unit-verbose` | All tests + verbose | 73 |
| `test-unit-handler-verbose` | Handler + verbose | 25 |
| `test-unit-service-verbose` | Service + verbose | 17 |
| `test-unit-repository-verbose` | Repository + verbose | 31 |
| `test-unit-race` | All tests + race detector | 73 |
| `test-unit-failfast` | Stop at first failure | 73 |
| `test-unit-timeout` | 60s timeout | 73 |
| `test-unit-coverage` | Coverage summary | 73 |
| `test-unit-coverage-html` | HTML coverage report | 73 |
| `test-unit-coverage-handler` | Handler coverage | 25 |
| `test-unit-coverage-service` | Service coverage | 17 |
| `test-unit-coverage-repository` | Repository coverage | 31 |
| `test-ci` | CI mode (coverage+race) | 73 |

## 📁 Test Files Location

```
tests/unit/
├── mocks/
│   ├── mock_user_repository.go
│   ├── mock_user_service.go
│   └── test_helpers.go
├── handler/
│   └── user_handler_test.go (25 tests)
├── service/
│   └── user_service_test.go (17 tests)
└── repository/
    └── user_repository_test.go (31 tests)
```

## ✅ Test Status

- **Total Tests**: 73
- **Status**: All ✅ PASSING
- **Coverage**: All 16 Users API endpoints
- **Execution Time**: ~0.4 seconds

## 📖 Documentation

See `TESTING.md` for:
- Complete endpoint coverage mapping
- Detailed test descriptions
- Testing patterns used
- Error scenarios covered
- Next steps for other entities
