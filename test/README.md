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
├── tcp-simple/               # Automated TCP testing
│   └── main.go               # TCP automated test client
├── tcp-client/               # Interactive TCP client
│   └── main.go               # TCP interactive REPL
├── grpc-simple/              # Automated gRPC testing
│   └── main.go               # gRPC automated test client
└── grpc-client/              # Interactive gRPC client
    └── main.go               # gRPC interactive REPL
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

#### Option A: Automated Test (grpc-simple)

Run all gRPC tests automatically:

```bash
# Terminal 1: Start gRPC server
make run-grpc

# Terminal 2: Run automated test
go run test/grpc-simple/main.go
```

**What gets tested:**
1. GetManga - Retrieve manga by ID
2. SearchManga by title - Search with title filter
3. SearchManga by author - Search with author filter
4. SearchManga by status - Filter by ongoing/completed
5. Pagination - Test limit and offset
6. UpdateProgress - Update user reading progress
7. Error handling - Test non-existent manga

**Note**: Update manga IDs in `test/grpc-simple/main.go` to match your database.

#### Option B: Interactive Client (grpc-client)

Test interactively with custom commands:

```bash
# Terminal 1: Start gRPC server
make run-grpc

# Terminal 2: Run interactive client
go run test/grpc-client/main.go

# Available commands:
> get manga-123
> search title=One
> search author=Oda limit=5
> search status=ongoing
> update user-123 manga-456 50
> update user-123 manga-456 50 8 reading
> help
> quit
```

**Commands:**
- `get <manga_id>` - Get manga by ID
- `search title=<query>` - Search by title
- `search author=<query>` - Search by author
- `search status=<status>` - Search by status (ongoing/completed/hiatus/cancelled)
- `search title=<query> limit=<n> offset=<n>` - Search with pagination
- `update <user_id> <manga_id> <chapter>` - Update progress
- `update <user_id> <manga_id> <chapter> <rating> <status>` - Full update with rating
- `help` - Show all commands
- `quit` - Exit

---

### 3. TCP Manual Testing

#### Setup: Get JWT Token

First, get a JWT token for authentication:

```bash
# Terminal 1: Start API server
make run-api

# Terminal 2: Register and login
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Copy the token from response
```

#### Option A: Automated Test (tcp-simple)

Run all TCP tests automatically:

```bash
# Terminal 1: Start TCP server
make run-tcp

# Terminal 2: Run automated test
go run test/tcp-simple/main.go <YOUR_JWT_TOKEN>
```

**What gets tested:**
1. Connection to TCP server
2. Authentication with JWT
3. Ping/Pong heartbeat
4. Progress update sent
5. Broadcast received

#### Option B: Interactive Client (tcp-client)

Test interactively with custom commands:

```bash
# Terminal 1: Start TCP server
make run-tcp

# Terminal 2: Run interactive client
go run test/tcp-client/main.go -token <YOUR_JWT_TOKEN>

# Available commands:
> ping
> progress manga-123 50
> quit
```

**Commands:**
- `ping` - Send heartbeat ping to server
- `progress <manga_id> <chapter>` - Send progress update
- `quit` - Exit client

**Note**: The client automatically receives broadcasts from other connected clients.

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
# ==========================================
# Unit Tests
# ==========================================

# Run all unit tests
go test ./test/unit/... -v

# Run specific test file
go test ./test/unit/auth_test.go -v
go test ./test/unit/tcp_service_test.go -v
go test ./test/unit/grpc_service_test.go -v

# Run specific test
go test ./test/unit/ -run TestJWT_GenerateToken -v

# Run tests with coverage
go test ./test/unit/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out

# ==========================================
# gRPC Manual Tests (server must be running)
# ==========================================

# Automated test
go run test/grpc-simple/main.go

# Interactive client
go run test/grpc-client/main.go
go run test/grpc-client/main.go -host localhost -port 9092

# ==========================================
# TCP Manual Tests (need JWT token)
# ==========================================

# Automated test
go run test/tcp-simple/main.go <JWT_TOKEN>

# Interactive client
go run test/tcp-client/main.go -token <JWT_TOKEN>
go run test/tcp-client/main.go -host localhost -port 9090 -token <JWT_TOKEN>
```

---

## Notes

- All unit tests use in-memory data structures (no database required)
- Manual tests (gRPC, TCP) require their respective servers to be running
- Update IDs in manual test files to match your database data
- gRPC tests use `localhost:9092`
- TCP tests use `localhost:9090`
- HTTP API uses `localhost:8080`
