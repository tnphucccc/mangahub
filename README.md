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
- Node.js 18 or later
- Yarn 4.0 or later
- SQLite3
- Git

### Installation & Setup

1.  **Clone the repository**

    ```bash
    git clone https://github.com/tnphucccc/mangahub.git
    cd mangahub
    ```

2.  **Install Go dependencies**

    ```bash
    go mod download
    ```

3.  **Install Node.js dependencies**

    ```bash
    yarn install
    ```

4.  **Run database migrations**
    This will create the necessary tables in the SQLite database.

    ```bash
    make migrate-up
    ```

5.  **Seed the database**
    This will populate the database with initial manga data.

    ```bash
    make seed
    ```

### Running the Backend (Go Servers)

The backend consists of four separate Go servers. You must run each in its own terminal. The `Makefile` provides the most convenient way to run them.

```bash
# Terminal 1: API Server (HTTP REST API & WebSockets)
make run-api

# Terminal 2: TCP Server (Real-time Sync)
make run-tcp

# Terminal 3: UDP Server (Notifications)
make run-udp

# Terminal 4: gRPC Server (Internal Services)
make run-grpc
```

### Running the Frontend (Next.js Web App)

The project includes a Next.js web application for the user interface.

```bash
# Run the web application in development mode
make js-dev
```

Once started, the application will be available at [http://localhost:3000](http://localhost:3000).

### Using the CLI

The project includes a command-line interface (CLI) for interacting with the backend.

1.  **Build the CLI tool**

    ```bash
    make build-cli
    ```

2.  **Run the CLI**
    You can see the available commands by running:

    ```bash
    ./bin/cli help
    ```

    This will output:

    ```
    MangaHub CLI - Manga Tracking System

    Usage:
      mangahub <command> [options]

    Commands:
      version              Show version information
      help                 Show this help message
      init                 Initialize configuration
      server               Manage servers (start, stop, status)
      auth                 Authentication (register, login, logout)
      manga                Manga operations (search, info, list)
      library              Library management (add, remove, list)
      progress             Progress tracking (update, history)
      chat                 Chat system (join, send)

    For more information on a command:
      mangahub <command> help
    ```

    **Note:** Most CLI commands are currently not implemented and are for demonstration purposes only.

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
- Multi-room support (per manga, general chat)
- User join/leave notifications
- Message broadcasting

**📖 Full WebSocket Documentation:** [docs/websocket-documentation.md](./docs/websocket-documentation.md)

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

### Protocol Documentation

- **[API Documentation](./docs/api-documentation.md)** - Complete REST API reference
- **[TCP Documentation](./docs/tcp-documentation.md)** - TCP progress sync protocol
- **[UDP Documentation](./docs/udp-documentation.md)** - UDP notification broadcasting
- **[WebSocket Documentation](./docs/websocket-documentation.md)** - WebSocket chat protocol
- **[gRPC Documentation](./docs/grpc-documentation.md)** - gRPC service reference

### System Documentation

- **[Architecture Documentation](./docs/architecture.md)** - System design and integration
- **[Database Documentation](./docs/database.md)** - Schema, migrations, queries
- **[Deployment Guide](./docs/deployment.md)** - Local, Docker, and production deployment
- **[Web Frontend Guide](./docs/web-frontend.md)** - Next.js application documentation

### Project Resources

- [Project Specification](./mangahub_project_spec.pdf) - Official requirements
- [Use Case Specification](<./mangahub_usecase%20(reference).pdf>) - Use cases
- [CLI Manual](<./mangahub_cli_manual%20(reference).pdf>) - CLI reference
- [Monorepo Structure](./MONOREPO.md) - Workspace organization
- [AI Assistant Context](./CLAUDE.md) - Development guidelines

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
