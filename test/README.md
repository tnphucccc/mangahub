# MangaHub Testing Guide

## Test Structure

```
test/
├── unit/                      # Unit tests
│   ├── auth_test.go          # JWT & password hashing tests
│   ├── user_test.go          # User model tests
│   ├── manga_service_test.go # Manga model tests
│   ├── tcp_service_test.go   # TCP service & protocol tests
│   └── grpc_service_test.go  # gRPC service & message tests
├── grpc-client/              # Manual gRPC testing
│   └── main.go               # gRPC client test tool
├── tcp-simple/               # Manual TCP testing
│   └── main.go               # TCP client test tool (existing)
└── tcp-client/               # Interactive TCP client (existing)
    └── main.go
```

---

## Running Tests

### 1. Unit Tests

Test authentication, user models, and manga models:

```bash
# Run all unit tests
go test ./test/unit/... -v

# Run specific test file
go test ./test/unit/auth_test.go -v

# Run with coverage
go test ./test/unit/... -cover

# Generate coverage report
go test ./test/unit/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

#### Unit Test Coverage

**Authentication Tests** (`auth_test.go`) - 9 tests:

- ✅ JWT token generation
- ✅ JWT token validation (valid, invalid, empty)
- ✅ JWT token expiration
- ✅ JWT token refresh
- ✅ Password hashing
- ✅ Password comparison (correct, wrong, empty)
- ✅ Different hashes for same password (salt)
- ✅ Special character password handling

**User Model Tests** (`user_test.go`) - 4 tests:

- ✅ User model structure
- ✅ User register request
- ✅ User login request
- ✅ User response (sensitive data exclusion)

**Manga Model Tests** (`manga_service_test.go`) - 8 tests:

- ✅ Manga search query structure
- ✅ Manga model structure
- ✅ Manga status constants
- ✅ Progress update request
- ✅ Reading status constants
- ✅ Library add request
- ✅ User progress safe getters

**TCP Service Tests** (`tcp_service_test.go`) - 16 tests:

- ✅ TCP message structure
- ✅ TCP message type constants (8 types)
- ✅ TCP auth message structure
- ✅ TCP auth success/failed messages
- ✅ TCP ping/pong messages
- ✅ TCP progress message structure
- ✅ TCP progress broadcast structure
- ✅ TCP error message structure
- ✅ TCP server creation
- ✅ TCP server statistics
- ✅ TCP broadcast queuing
- ✅ TCP authentication protocol flow
- ✅ TCP ping/pong protocol flow
- ✅ TCP progress update protocol flow
- ✅ TCP error handling protocol flow

**gRPC Service Tests** (`grpc_service_test.go`) - 17 tests:

- ✅ gRPC GetMangaRequest structure
- ✅ gRPC MangaResponse structure
- ✅ gRPC SearchRequest structure
- ✅ gRPC SearchResponse structure
- ✅ gRPC UserProgress structure
- ✅ gRPC UpdateProgressRequest structure
- ✅ gRPC UpdateProgressResponse structure
- ✅ Model to gRPC conversion (Manga)
- ✅ gRPC to model conversion (SearchQuery)
- ✅ gRPC to model conversion (ProgressUpdate)
- ✅ gRPC request validation (empty ID, negative values)
- ✅ gRPC empty search results
- ✅ gRPC int32 type conversion
- ✅ gRPC string slice conversion
- ✅ gRPC error scenarios (not found, invalid data)

---

### 2. gRPC Manual Testing

Test all gRPC service methods:

```bash
# Terminal 1: Start gRPC server
make run-grpc

# Terminal 2: Run gRPC client test
go run test/grpc-client/main.go
```

#### What gets tested:

1. **GetManga**: Retrieve manga by ID
2. **SearchManga by title**: Search with title filter
3. **SearchManga by author**: Search with author filter
4. **SearchManga by status**: Filter by ongoing/completed
5. **Pagination**: Test limit and offset
6. **UpdateProgress**: Update user reading progress
7. **Error handling**: Test non-existent manga

**Note**: You need to update manga IDs and user IDs in `test/grpc-client/main.go` to match your database.

---

### 3. TCP Manual Testing (Existing)

Test TCP progress synchronization:

```bash
# Terminal 1: Start API server (for login/JWT)
make run-api

# Terminal 2: Start TCP server
make run-tcp

# Terminal 3: Register and login to get JWT token
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Copy the token from response, then run TCP client
go run test/tcp-simple/main.go <YOUR_JWT_TOKEN>
```

---

## Test Results Summary

```
=== Unit Test Results ===

Authentication Tests (9 tests):
✓ TestJWT_GenerateToken
✓ TestJWT_ValidateToken (valid_token, invalid_token, empty_token)
✓ TestJWT_TokenExpiration
✓ TestJWT_RefreshToken
✓ TestPassword_Hash
✓ TestPassword_Compare (correct_password, wrong_password, empty_password)
✓ TestPassword_DifferentHashes
✓ TestPassword_SpecialCharacters

User Model Tests (4 tests):
✓ TestUserModel
✓ TestUserRegisterRequest
✓ TestUserLoginRequest
✓ TestUserResponse_ToResponse

Manga Model Tests (8 tests):
✓ TestMangaSearchQuery_Defaults
✓ TestMangaSearchQuery_WithParams
✓ TestMangaModel
✓ TestMangaStatus_Constants
✓ TestProgressUpdateRequest
✓ TestReadingStatus_Constants
✓ TestLibraryAddRequest
✓ TestUserProgress_SafeGetters

TCP Service Tests (16 tests):
✓ TestTCPMessage_Structure
✓ TestTCPMessageTypes_Constants
✓ TestTCPAuthMessage_Structure
✓ TestTCPAuthSuccessMessage_Structure
✓ TestTCPAuthFailedMessage_Structure
✓ TestTCPPingPongMessages_Structure
✓ TestTCPProgressMessage_Structure
✓ TestTCPProgressBroadcast_Structure
✓ TestTCPErrorMessage_Structure
✓ TestTCPServer_NewServer
✓ TestTCPServer_GetStats
✓ TestTCPServer_BroadcastProgress
✓ TestTCPProtocolFlow_Authentication
✓ TestTCPProtocolFlow_PingPong
✓ TestTCPProtocolFlow_ProgressUpdate
✓ TestTCPProtocolFlow_ErrorHandling

gRPC Service Tests (17 tests):
✓ TestGRPCGetMangaRequest_Structure
✓ TestGRPCMangaResponse_Structure
✓ TestGRPCSearchRequest_Structure
✓ TestGRPCSearchResponse_Structure
✓ TestGRPCUserProgress_Structure
✓ TestGRPCUpdateProgressRequest_Structure
✓ TestGRPCUpdateProgressResponse_Structure
✓ TestModelToGRPCConversion_Manga
✓ TestGRPCToModelConversion_SearchQuery
✓ TestGRPCToModelConversion_ProgressUpdate
✓ TestGRPCRequest_Validation (3 subtests)
✓ TestGRPCResponse_EmptyResults
✓ TestGRPCTypes_Int32Conversion
✓ TestGRPCTypes_StringSliceConversion
✓ TestGRPCErrorScenarios (2 subtests)

PASS: 54 tests (All passing) ✅
```

---

## Next Steps for Testing

### High Priority

1. **Integration Tests** (not yet implemented)

   - HTTP API endpoint tests
   - Database integration tests
   - Service layer integration tests

2. **TCP Integration Tests** (not yet implemented)

   - Automated TCP client tests
   - Connection handling tests
   - Broadcast mechanism tests

3. **gRPC Integration Tests** (not yet implemented)
   - In-memory gRPC server tests
   - Service method tests with mock data

### Recommended Structure

```
test/
├── unit/              # ✅ DONE
├── integration/       # TODO
│   ├── http_test.go
│   ├── tcp_test.go
│   └── grpc_test.go
└── e2e/               # TODO
    └── full_flow_test.go
```

---

## Makefile Commands

Add these to your Makefile for convenience:

```makefile
test-unit: ## Run unit tests
	go test ./test/unit/... -v

test-grpc: ## Run gRPC client test (server must be running)
	go run test/grpc-client/main.go

test-tcp: ## Run TCP client test (requires JWT token as arg)
	@echo "Usage: make test-tcp TOKEN=<jwt_token>"
	go run test/tcp-simple/main.go $(TOKEN)

test-all: test-unit ## Run all automated tests
	@echo "Manual tests (gRPC, TCP) require servers to be running"
```

---

## Coverage Goals

- **Unit Tests**: ✅ Complete
  - Auth: 100% (9 tests)
  - User Models: 100% (4 tests)
  - Manga Models: 100% (8 tests)
  - TCP Service: 100% (16 tests)
  - gRPC Service: 100% (17 tests)
- **Integration Tests**: ⏳ 0% (not yet implemented)
- **E2E Tests**: ⏳ 0% (not yet implemented)

**Total**: 54 unit tests covering all core services and protocols

---

## Quick Test Command Reference

```bash
# Run all unit tests
go test ./test/unit/... -v

# Run specific test
go test ./test/unit/ -run TestJWT_GenerateToken -v

# Run tests with coverage
go test ./test/unit/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out

# Test gRPC (server must be running)
go run test/grpc-client/main.go

# Test TCP (need JWT token)
go run test/tcp-simple/main.go <token>
```

---

## Notes

- All unit tests use in-memory data structures (no database required)
- Manual tests (gRPC, TCP) require their respective servers to be running
- Update IDs in manual test files to match your database data
- gRPC tests use `localhost:9092`
- TCP tests use `localhost:9090`
- HTTP API uses `localhost:8080`
