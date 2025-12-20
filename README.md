# MangaHub - Manga & Comic Tracking System

A network programming project demonstrating **5 network protocols** (TCP, UDP, HTTP, gRPC, WebSocket) in Go for tracking manga reading progress.

**Course**: Network Programming (Net Centric Programming) – IT096IU
**Team Size**: 2 students
**Language**: Go 1.19+
**Timeline**: 10-11 weeks

---

## 🎯 Project Objectives

- Implement all 5 required network protocols (TCP, UDP, HTTP, gRPC, WebSocket)
- Build a practical manga tracking system with user authentication
- Demonstrate concurrent programming with goroutines
- Create a functional CLI tool for user interaction
- Track reading progress across multiple devices in real-time

---

## 🏗️ Architecture

```
MangaHub
├── HTTP REST API Server (port 8080)   - User auth, manga CRUD, library management
├── TCP Sync Server (port 9090)        - Real-time progress synchronization
├── UDP Notification Server (port 9091) - Chapter release notifications
├── gRPC Internal Service (port 9092)   - Internal service communication
├── WebSocket Chat (port 9093)          - Real-time manga discussions
└── CLI Client                          - Command-line interface
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.19 or later
- SQLite3
- Git

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/tnphucccc/mangahub.git
   cd mangahub
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Run database migrations**

   ```bash
   make migrate-up
   ```

4. **Seed the database**
   ```bash
   make seed-db
   ```

### Running the Servers

**Option 1: Development (separate terminals)**

```bash
# Terminal 1 - HTTP API Server
go run cmd/api-server/main.go

# Terminal 2 - TCP Sync Server
go run cmd/tcp-server/main.go

# Terminal 3 - UDP Notification Server
go run cmd/udp-server/main.go

# Terminal 4 - gRPC Internal Service
go run cmd/grpc-server/main.go
```

**Option 2: Production (Docker Compose)**

```bash
docker-compose up
```

### Using the CLI

```bash
# Build CLI tool
go build -o mangahub cmd/cli/main.go

# Register a new account
./mangahub auth register --username johndoe --email john@example.com

# Login
./mangahub auth login --username johndoe

# Search for manga
./mangahub manga search "one piece"

# Add manga to library
./mangahub library add --manga-id one-piece --status reading

# Update reading progress
./mangahub progress update --manga-id one-piece --chapter 1095
```

---

## 📂 Project Structure

```
mangahub/
├── cmd/                      # Main applications (5 servers + CLI)
│   ├── api-server/          # HTTP REST API server
│   ├── tcp-server/          # TCP sync server
│   ├── udp-server/          # UDP notification server
│   ├── grpc-server/         # gRPC internal service
│   └── cli/                 # CLI client tool
│
├── internal/                 # Private application code
│   ├── auth/                # Authentication & JWT
│   ├── manga/               # Manga management
│   ├── user/                # User management
│   ├── library/             # User library
│   ├── progress/            # Progress tracking
│   ├── tcp/                 # TCP server implementation
│   ├── udp/                 # UDP server implementation
│   ├── websocket/           # WebSocket chat
│   └── grpc/                # gRPC service implementation
│
├── pkg/                      # Shared libraries
│   ├── models/              # Data models
│   ├── database/            # Database utilities
│   ├── config/              # Configuration management
│   └── utils/               # Helper functions
│
├── proto/                    # Protocol Buffer definitions
├── migrations/               # SQL database migrations
├── data/                     # Manga data (JSON files)
├── scripts/                  # Utility scripts
├── test/                     # Tests (unit, integration, e2e)
├── docs/                     # Documentation
└── configs/                  # Configuration files
```

---

## 🔌 Network Protocols

### 1. HTTP REST API (25 points)

- User registration and authentication (JWT)
- Manga search and CRUD operations
- Library management
- Reading progress tracking

**📖 Full API Documentation:** [docs/api-documentation.md](./docs/api-documentation.md)

**Quick Example:**

```bash
# Register a new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "email": "alice@example.com", "password": "alice123"}'

# Login and get JWT token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "password": "alice123"}'
```

### 2. TCP Progress Sync (20 points)

- Real-time progress synchronization across devices
- Concurrent connection handling with goroutines
- JSON-based message protocol over TCP
- JWT authentication for secure connections
- Broadcast mechanism for instant updates

**📖 Full TCP Documentation:** [docs/tcp-documentation.md](./docs/tcp-documentation.md)

**Quick Example:**

```bash
# Start TCP server
go run cmd/tcp-server/main.go

# In another terminal, test with automated client
TOKEN="your-jwt-token"
go run test/tcp-simple/main.go $TOKEN

# Or use interactive client for manual testing
go run test/tcp-client/main.go -token $TOKEN
```

### 3. UDP Notifications (15 points)

- Chapter release notifications
- Client registration mechanism
- Broadcast to multiple clients

### 4. WebSocket Chat (15 points)

- Real-time manga discussions
- User join/leave notifications
- Message broadcasting

### 5. gRPC Internal Service (10 points)

- Internal service-to-service communication
- Protocol Buffer definitions
- Unary RPC calls for manga retrieval and progress updates

**📖 Full gRPC Documentation:** [docs/grpc-documentation.md](./docs/grpc-documentation.md)

**Quick Example:**

```bash
# Install grpcurl for testing
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Get manga by ID
grpcurl -plaintext -d '{"manga_id": "manga-001"}' \
  localhost:9092 manga.MangaService/GetManga

# Search by title
grpcurl -plaintext -d '{"title": "naruto", "limit": 10}' \
  localhost:9092 manga.MangaService/SearchManga
```

---

## 🗄️ Database Schema

**SQLite3 database with 3 core tables:**

```sql
-- users table
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- manga table
CREATE TABLE manga (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT,
    genres TEXT,  -- JSON array
    status TEXT,
    total_chapters INTEGER,
    description TEXT
);

-- user_progress table
CREATE TABLE user_progress (
    user_id TEXT NOT NULL,
    manga_id TEXT NOT NULL,
    current_chapter INTEGER,
    status TEXT,
    updated_at TIMESTAMP,
    PRIMARY KEY (user_id, manga_id)
);
```

---

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run integration tests
make test-integration

# Test specific protocol
go test ./internal/tcp/...
```

---

## 📝 Development

### Available Make Commands

```bash
make help           # Show all available commands
make build          # Build all binaries
make run-api        # Run HTTP API server
make run-tcp        # Run TCP server
make run-udp        # Run UDP server
make run-grpc       # Run gRPC server
make migrate-up     # Run database migrations
make migrate-down   # Rollback migrations
make seed-db        # Seed database with manga data
make proto          # Generate gRPC code from .proto files
make test           # Run all tests
make clean          # Clean build artifacts
```

### Adding a New Feature

Follow the pattern:

1. Create model in `pkg/models/`
2. Create migration in `migrations/`
3. Create repository in `internal/<feature>/repository.go`
4. Create service in `internal/<feature>/service.go`
5. Create handler in `internal/<feature>/handler.go`
6. Add routes to server
7. Write tests

---

## 📚 Documentation

- **[API Documentation](./docs/api-documentation.md)** - Complete REST API reference
- **[TCP Documentation](./docs/tcp-documentation.md)** - TCP progress sync protocol
- **[gRPC Documentation](./docs/grpc-documentation.md)** - gRPC service reference
- [Project Specification](./mangahub_project_spec.pdf)
- [Use Case Specification](<./mangahub_usecase%20(reference).pdf>)
- [CLI Manual](<./mangahub_cli_manual%20(reference).pdf>)
- [CLAUDE.md](./CLAUDE.md) - AI assistant context

---

## 🎓 Academic Requirements

### Grading Criteria (100 points)

- **Core Protocol Implementation (40 pts)**

  - HTTP REST API: 15 pts
  - TCP Progress Sync: 13 pts
  - UDP Notifications: 18 pts
  - WebSocket Chat: 10 pts
  - gRPC Service: 7 pts

- **System Integration (20 pts)**

  - Database Integration: 8 pts
  - Service Communication: 7 pts
  - Error Handling: 3 pts
  - Code Organization: 2 pts

- **Code Quality (10 pts)**

  - Go Idioms: 5 pts
  - Testing: 3 pts
  - Documentation: 2 pts

- **Documentation & Demo (10 pts)**

  - Technical Documentation: 5 pts
  - Live Demonstration: 5 pts

- **Bonus Features (up to 20 pts)**
  - Docker Compose: 10 pts
  - Advanced Features: 5-10 pts each

---

## 🤝 Contributing

This is an academic project. For team members:

1. Create a feature branch: `git checkout -b feature/your-feature`
2. Make your changes following Go best practices
3. Write tests for new features
4. Update documentation
5. Create a pull request

---

## 📄 License

This project is for educational purposes as part of IT096IU coursework.

---

## 👥 Team

- [Your Name] - [Your Student ID]
- [Team Member] - [Student ID]

**Instructor**: Lê Thanh Sơn - Nguyễn Trung Nghĩa

---

## 🔗 Resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [gRPC Go Documentation](https://grpc.io/docs/languages/go/)
- [Gin Web Framework](https://gin-gonic.com/)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)

---

**Status**: 🚧 In Development
**Version**: 1.0.0-dev
