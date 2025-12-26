# MangaHub System Architecture Documentation

---

## 1. Overview

MangaHub is a **multi-protocol network programming demonstration project** built for the IT096IU Network Programming course. The system showcases the practical implementation and integration of **5 distinct network protocols** in a unified manga tracking application.

**Project Type**: Academic Network Programming Project
**Primary Goal**: Demonstrate mastery of multiple network protocols in a real-world application
**Language**: Go 1.19+ (Backend), TypeScript (Frontend)
**Architecture Pattern**: Microservices with Protocol-based Separation

### Core Design Principles

1. **Protocol Isolation**: Each protocol runs in its own dedicated server process
2. **Shared Business Logic**: Common functionality lives in shared packages
3. **Clear Separation of Concerns**: Handler → Service → Repository pattern
4. **Academic Clarity**: Code optimized for learning and demonstration

---

## 2. High-Level Architecture

```
┌────────────────────────────────────────────────────────────────────┐
│                         MangaHub System                            │
├────────────────────────────────────────────────────────────────────┤
│                                                                    │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌──────────────┐  │
│  │ CLI Client │  │ Next.js    │  │ Mobile App │  │ Other Clients│  │
│  │            │  │ Web App    │  │  (Future)  │  │   (gRPC)     │  │
│  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘  └──────┬───────┘  │
│        │               │               │                │          │
└────────┼───────────────┼───────────────┼────────────────┼──────────┘
         │               │               │                │
         │          HTTP │          WebSocket        gRPC │
         │               │               │                │
┌────────┼───────────────┼───────────────┼────────────────┼──────────┐
│        │               │               │                │          │
│        ▼               ▼               ▼                ▼          │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │         API Server (HTTP + WebSocket)                       │   │
│  │         • HTTP REST API (Port 8080)                         │   │
│  │         • WebSocket Chat (Port 9093)                        │   │
│  │         • UDP Client (listens for notifications)            │   │
│  │         • Bridges UDP notifications → WebSocket             │   │
│  └────────┬─────────────────────────────────────────┬──────────┘   │
│           │                                         │              │
│           │ Notifies via TCP/UDP                    │ Reads        │
│           ▼                                         ▼              │
│  ┌────────────────┐        ┌────────────────┐  ┌──────────────┐    │
│  │  TCP Server    │        │  UDP Server    │  │ gRPC Server  │    │
│  │  Port 9090     │        │  Port 9091     │  │ Port 9092    │    │
│  │  • Progress    │        │  • Chapter     │  │ • Internal   │    │
│  │    Sync        │        │    Releases    │  │   Service    │    │
│  │  • Broadcast   │        │  • Broadcast   │  │ • Manga API  │    │
│  └────────┬───────┘        └────────┬───────┘  └──────┬───────┘    │
│           │                         │                 │            │
│           └─────────────┬───────────┘─────────────────┘            │
│                         │                                          │
│                         ▼                                          │
│                  ┌──────────────┐                                  │
│                  │   SQLite DB  │                                  │
│                  │  mangahub.db │                                  │
│                  └──────────────┘                                  │
│                                                                    │
│                   Server Layer (4 Go Processes)                    │
└────────────────────────────────────────────────────────────────────┘
```

---

## 3. Protocol Breakdown

### Overview Table

| Protocol  | Port | Server Binary     | Purpose                                | Points |
| --------- | ---- | ----------------- | -------------------------------------- | ------ |
| HTTP      | 8080 | `cmd/api-server`  | REST API, User Auth, Manga CRUD        | 25     |
| WebSocket | 9093 | `cmd/api-server`  | Real-time chat, Multi-room support     | 15     |
| TCP       | 9090 | `cmd/tcp-server`  | Progress synchronization, Broadcasting | 20     |
| UDP       | 9091 | `cmd/udp-server`  | Chapter release notifications          | 15     |
| gRPC      | 9092 | `cmd/grpc-server` | Internal service communication         | 10     |

**Note**: HTTP and WebSocket share the same process but run on different ports.

---

## 4. Component Architecture

### 4.1 API Server (`cmd/api-server/`)

**Responsibilities:**

- HTTP REST API endpoints for user and manga operations
- WebSocket server for real-time chat
- JWT authentication and authorization
- Acts as a bridge between UDP notifications and WebSocket clients
- Sends notifications to TCP/UDP servers when data changes

**Key Components:**

```
API Server
├── HTTP Server (Gin Framework)
│   ├── Auth Routes (/api/v1/auth)
│   │   ├── POST /register
│   │   └── POST /login
│   ├── Manga Routes (/api/v1/manga)
│   │   ├── GET /manga (search)
│   │   ├── GET /manga/:id
│   │   └── GET /manga/all
│   └── User Routes (/api/v1/users - protected)
│       ├── GET /me
│       ├── GET /library
│       ├── POST /library
│       ├── GET /progress/:manga_id
│       └── PUT /progress/:manga_id
│
├── WebSocket Server (Gorilla WebSocket)
│   ├── Hub (central message router)
│   ├── Clients (connected users)
│   └── Rooms (multi-room support)
│
├── UDP Client (listener)
│   ├── Registers with UDP server
│   ├── Receives chapter notifications
│   └── Forwards to WebSocket hub
│
└── Notification Bridge (goroutines)
    ├── TCP Progress Channel → TCP Server
    └── UDP Notification Channel → UDP Server
```

**Data Flow:**

1. User updates progress via HTTP PUT /api/v1/users/progress/:manga_id
2. Service saves to database
3. Service sends message to `TCPBroadcastChan`
4. Bridge goroutine picks up message and notifies TCP server
5. TCP server broadcasts to all connected clients

**Special Integration:**

- API server acts as a UDP **client** to receive notifications
- Bridges UDP notifications to WebSocket for real-time web updates
- This enables: UDP notification → WebSocket chat → "New chapter released!" message

---

### 4.2 TCP Server (`cmd/tcp-server/`)

**Responsibilities:**

- Accept multiple concurrent TCP connections
- Authenticate clients using JWT tokens
- Broadcast progress updates to all connected clients
- Maintain active connection pool

**Architecture:**

```
TCP Server
├── Listener (net.Listen on :9090)
├── Connection Pool (map[string]net.Conn)
├── Broadcast Channel (for updates)
└── Goroutines
    ├── Accept Loop (accepts new connections)
    ├── Read Pump per client (reads incoming messages)
    ├── Write Pump per client (sends outgoing messages)
    └── Broadcast Loop (distributes updates)
```

**Message Types:**

- `auth` - Client authentication with JWT
- `subscribe` - Subscribe to progress updates for a manga
- `progress` - Progress update notification
- `sync_request` - Request current progress
- `sync_response` - Response with progress data
- `ping/pong` - Keepalive heartbeat

**Concurrency Model:**

- 1 goroutine per connected client for reading
- 1 goroutine per connected client for writing
- Shared broadcast channel with mutex protection
- Non-blocking sends prevent slow clients from blocking others

---

### 4.3 UDP Server (`cmd/udp-server/`)

**Responsibilities:**

- Listen for client registrations
- Maintain list of registered client addresses
- Broadcast chapter release notifications to all clients
- Handle ping/pong heartbeats
- Automatic cleanup of stale clients

**Architecture:**

```
UDP Server
├── UDP Listener (net.ListenUDP on :9091)
├── Client Registry
│   ├── Clients map (ClientID → *UDPAddr)
│   ├── LastSeen timestamps
│   └── Cleanup routine (removes stale clients)
├── Notification Channel
└── Goroutines
    ├── Message Handler (reads incoming messages)
    ├── Broadcast Loop (sends notifications)
    └── Cleanup Loop (removes inactive clients)
```

**Message Types:**

- `register` - Register to receive notifications
- `unregister` - Stop receiving notifications
- `ping` - Heartbeat from client
- `pong` - Heartbeat response from server
- `notification` - Chapter release broadcast
- `register_success` - Confirmation of registration
- `register_failed` - Registration error
- `error` - Generic error message

**Key Features:**

- **Fire-and-forget**: UDP provides no delivery guarantees
- **Lightweight**: No connection state maintained
- **Scalable**: Can handle many clients with minimal overhead
- **Heartbeat**: Ping/pong keeps clients from being removed

---

### 4.4 WebSocket Server (Part of API Server)

**Responsibilities:**

- Real-time chat for manga discussions
- Multi-room support (per manga, general chat)
- User join/leave notifications
- Message broadcasting within rooms
- Integration with UDP notifications

**Architecture (Hub Pattern):**

```
WebSocket Hub
├── Rooms map[string]map[*Client]bool
│   ├── "general" room
│   ├── "one-piece" room
│   └── "naruto" room (dynamically created)
│
├── Channels
│   ├── Broadcast chan (messages to send)
│   ├── Register chan (new client joins)
│   └── Unregister chan (client leaves)
│
├── Hub.Run() (central goroutine)
│   ├── Listens on channels
│   ├── Manages room membership
│   └── Broadcasts messages
│
└── Per-Client Goroutines
    ├── Read Pump (client → hub)
    └── Write Pump (hub → client)
```

**Message Types:**

- `join_room` - Join a specific chat room
- `leave_room` - Leave current room
- `chat` - Regular chat message
- `system` - System notifications (joins, leaves, etc.)
- `error` - Error messages

**Room Lifecycle:**

- Created automatically when first user joins
- Deleted automatically when last user leaves
- "general" room always available

---

### 4.5 gRPC Server (`cmd/grpc-server/`)

**Responsibilities:**

- Internal service-to-service communication
- Efficient binary protocol for high performance
- Type-safe API contracts via Protocol Buffers
- Future integration point for mobile apps

**Services:**

```protobuf
service MangaService {
  rpc GetManga(GetMangaRequest) returns (MangaResponse);
  rpc SearchManga(SearchRequest) returns (SearchResponse);
  rpc GetUserProgress(GetProgressRequest) returns (ProgressResponse);
  rpc UpdateProgress(UpdateProgressRequest) returns (UpdateProgressResponse);
}
```

**Use Cases:**

- Internal microservice communication (if services were split further)
- Mobile app integration (Flutter, React Native)
- Third-party integrations
- Performance-critical operations

**Currently:**

- Implements Manga search and retrieval
- Uses same service layer as HTTP API
- Demonstrates Protocol Buffer usage

---

## 5. Shared Components

### 5.1 Database Layer (`pkg/database/`)

**Technology**: SQLite3
**Location**: `data/mangahub.db`

**Schema:**

```sql
-- Users
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Manga Library
CREATE TABLE manga (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT,
    genres TEXT,  -- JSON array as TEXT
    status TEXT CHECK(status IN ('ongoing', 'completed', 'hiatus')),
    total_chapters INTEGER DEFAULT 0,
    description TEXT,
    cover_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User Progress Tracking
CREATE TABLE user_progress (
    user_id TEXT NOT NULL,
    manga_id TEXT NOT NULL,
    current_chapter INTEGER DEFAULT 0,
    status TEXT CHECK(status IN ('reading', 'completed', 'plan_to_read', 'on_hold', 'dropped')),
    rating INTEGER CHECK(rating BETWEEN 1 AND 10),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, manga_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (manga_id) REFERENCES manga(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_progress_user ON user_progress(user_id);
CREATE INDEX idx_user_progress_manga ON user_progress(manga_id);
CREATE INDEX idx_manga_status ON manga(status);
```

**Connection Pooling:**

```go
db.SetMaxOpenConns(25)  // Maximum 25 concurrent connections
db.SetMaxIdleConns(5)   // Keep 5 idle connections
```

**Why SQLite?**

- Academic project requirement
- Simplifies deployment (single file database)
- Sufficient for 50-100 concurrent users (project requirement)
- No external database server needed

---

### 5.2 Business Logic Layer

**Pattern**: Handler → Service → Repository

```
Handler (HTTP/gRPC endpoint)
    ↓ Validates input, handles HTTP concerns
Service (Business Logic)
    ↓ Implements core functionality
Repository (Data Access)
    ↓ SQL queries, data persistence
Database
```

**Example: Update Progress**

```go
// 1. Handler (internal/manga/handler.go)
func (h *Handler) UpdateProgress(c *gin.Context) {
    var req UpdateProgressRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    userID := c.GetString("user_id") // From JWT middleware
    progress, err := h.service.UpdateProgress(userID, req.MangaID, req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, progress)
}

// 2. Service (internal/manga/service.go)
func (s *Service) UpdateProgress(userID, mangaID string, req UpdateProgressRequest) (*UserProgress, error) {
    // Business logic: validate, compute, trigger events
    progress, err := s.repo.UpdateProgress(userID, mangaID, req)
    if err != nil {
        return nil, err
    }

    // Trigger TCP broadcast for real-time sync
    s.TCPBroadcastChan <- TCPProgressBroadcast{
        MangaID:        mangaID,
        Username:       req.Username,
        CurrentChapter: req.CurrentChapter,
    }

    return progress, nil
}

// 3. Repository (internal/manga/repository.go)
func (r *Repository) UpdateProgress(userID, mangaID string, req UpdateProgressRequest) (*UserProgress, error) {
    // SQL query
    query := `UPDATE user_progress
              SET current_chapter = ?, status = ?, updated_at = ?
              WHERE user_id = ? AND manga_id = ?`
    _, err := r.db.Exec(query, req.CurrentChapter, req.Status, time.Now(), userID, mangaID)
    if err != nil {
        return nil, err
    }

    return r.GetProgress(userID, mangaID)
}
```

---

### 5.3 Configuration Management (`pkg/config/`)

**Configuration File**: `configs/dev.yaml`

```yaml
server:
  host: "localhost"
  http_port: "8080"
  tcp_port: "9090"
  udp_port: "9091"
  grpc_port: "9092"
  websocket_port: "9093"

database:
  path: "./data/mangahub.db"

jwt:
  secret: "your-secret-key-change-in-production"
  expiry_days: 7
```

**Environment Variable Overrides:**

```bash
export CONFIG_PATH="./configs/prod.yaml"
export TCP_HOST="tcp-server"  # For Docker networking
export UDP_HOST="udp-server"
export GRPC_HOST="grpc-server"
```

---

## 6. Communication Patterns

### 6.1 HTTP → TCP Notification Flow

```
User Updates Progress (HTTP)
        ↓
HTTP Handler receives request
        ↓
Service updates database
        ↓
Service sends to TCPBroadcastChan
        ↓
API Server bridge goroutine
        ↓
Connects to TCP Server (net.Dial)
        ↓
Sends TCPMessage{Type: "progress", Data: {...}}
        ↓
TCP Server receives on connection
        ↓
TCP Server broadcasts to all connected clients
        ↓
All clients receive progress update
```

---

### 6.2 UDP → WebSocket Bridge Flow

```
Admin triggers chapter release (HTTP POST /admin/notifications)
        ↓
API Server sends to UDPNotificationChan
        ↓
API Server bridge goroutine
        ↓
Sends UDP packet to UDP Server
        ↓
UDP Server receives notification
        ↓
UDP Server broadcasts to all registered clients
        ↓ (including API Server as UDP client)
API Server receives UDP notification
        ↓
API Server's UDP listener goroutine
        ↓
Converts to WebSocketMessage
        ↓
Sends to WebSocket Hub.Broadcast
        ↓
Hub broadcasts to all WebSocket clients in "general" room
        ↓
Web users see "New Chapter Released: One Piece Ch 1150!"
```

**Why this design?**

- UDP server handles notification broadcasting (its primary job)
- API server bridges UDP → WebSocket for web clients
- Demonstrates integration between different protocols
- Each server has clear, single responsibility

---

### 6.3 Client Authentication Flow

```
1. User Registration (HTTP)
   POST /api/v1/auth/register
        ↓
   Password hashed (bcrypt)
        ↓
   User stored in database
        ↓
   Return success

2. User Login (HTTP)
   POST /api/v1/auth/login
        ↓
   Verify password hash
        ↓
   Generate JWT token (7-day expiry)
        ↓
   Return {token, user_id, username}

3. TCP Connection with JWT
   Client connects to TCP server
        ↓
   Client sends: {type: "auth", data: {token: "..."}}
        ↓
   TCP server validates JWT
        ↓
   JWT valid → Client added to pool
   JWT invalid → Connection closed

4. Protected HTTP Endpoints
   Request with header: Authorization: Bearer <token>
        ↓
   Middleware validates JWT
        ↓
   Extract user_id from token claims
        ↓
   Inject into Gin context
        ↓
   Handler accesses via c.GetString("user_id")
```

---

## 7. Concurrency Model

### Goroutines per Server

**API Server:**

- 1 main goroutine (HTTP server)
- 1 WebSocket hub goroutine
- 1 UDP listener goroutine
- 1 UDP pinger goroutine (heartbeat)
- 1 UDP receiver goroutine
- 1 TCP/UDP notification bridge goroutine
- N request handler goroutines (Gin creates per request)

**TCP Server:**

- 1 main goroutine (accept loop)
- 2N goroutines per client (read pump + write pump)
- 1 broadcast goroutine

**UDP Server:**

- 1 main goroutine (message handler)
- 1 broadcast goroutine
- 1 cleanup goroutine (removes stale clients)

**WebSocket Server (part of API):**

- 1 hub goroutine
- 2N goroutines per WebSocket client (read pump + write pump)

**gRPC Server:**

- 1 main goroutine (gRPC serve)
- N goroutines per RPC call (gRPC manages internally)

**Synchronization:**

- Channels for communication between goroutines
- `sync.RWMutex` for shared data structures (client maps, room maps)
- Context for cancellation and graceful shutdown

---

## 8. Security Model

### Authentication

- **JWT (JSON Web Tokens)** for stateless authentication
- Tokens issued on login, valid for 7 days
- Tokens contain: `user_id`, `username`, `issued_at`, `expires_at`
- Secret key: Configured in `configs/dev.yaml` (change for production)

### Authorization

- **Middleware-based** on HTTP API
- TCP server validates JWT before accepting commands
- WebSocket uses query parameter username (simplified for academic project)

### Password Security

- **bcrypt hashing** (cost factor 10)
- Passwords never stored in plaintext
- No password transmission over unencrypted channels (academic environment)

### Recommendations for Production

- Use HTTPS/TLS for HTTP API
- Use TLS for TCP connections
- Implement rate limiting
- Add CORS restrictions
- Use environment variables for secrets
- Implement refresh tokens
- Add two-factor authentication

---

## 9. Data Flow Diagrams

### User Reading Progress Sync

```
┌─────────────┐
│ Web Browser │
└──────┬──────┘
       │ 1. Update progress (HTTP PUT /users/progress/manga-001)
       │    {current_chapter: 42, status: "reading"}
       ▼
┌──────────────────┐
│   API Server     │
│   HTTP Handler   │──────► 2. Save to database
└────────┬─────────┘
         │ 3. Send to TCPBroadcastChan
         ▼
   ┌─────────────┐
   │ TCP Server  │
   └─────┬───────┘
         │ 4. Broadcast to all TCP clients
         ▼
   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
   │ Desktop App  │  │  Mobile App  │  │  CLI Client  │
   │ (TCP Client) │  │ (TCP Client) │  │ (TCP Client) │
   └──────────────┘  └──────────────┘  └──────────────┘
   All receive: {"type": "progress", "manga_id": "manga-001",
                 "username": "alice", "current_chapter": 42}
```

---

### New Chapter Release Notification

```
┌────────────────┐
│ Admin/Scraper  │
└────────┬───────┘
         │ 1. Trigger notification (HTTP POST /admin/notifications)
         │    {manga_title: "One Piece", chapter_number: 1150}
         ▼
   ┌──────────────┐
   │  API Server  │──────► 2. Send to UDPNotificationChan
   └──────┬───────┘
          │ 3. UDP packet to UDP Server
          ▼
    ┌────────────┐
    │ UDP Server │
    └─────┬──────┘
          │ 4. Broadcast to all registered UDP clients
          ▼
    ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
    │  API Server  │  │  CLI Client  │  │ Mobile App   │
    │ (UDP Client) │  │ (UDP Client) │  │ (UDP Client) │
    └──────┬───────┘  └──────────────┘  └──────────────┘
           │ 5. Forward to WebSocket Hub
           ▼
    ┌──────────────┐
    │ WebSocket    │
    │ "general"    │
    │ room         │
    └──────┬───────┘
           │ 6. Broadcast to all WebSocket users
           ▼
    ┌──────────────┐  ┌──────────────┐
    │ Web Browser  │  │ Web Browser  │
    │   (Alice)    │  │    (Bob)     │
    └──────────────┘  └──────────────┘
    Both see: "New Chapter Release: One Piece - Chapter 1150!"
```

---

## 10. Deployment Architecture

### Development (Local)

```
Developer Machine (localhost)
├── Terminal 1: go run cmd/api-server/main.go     (ports 8080, 9093)
├── Terminal 2: go run cmd/tcp-server/main.go     (port 9090)
├── Terminal 3: go run cmd/udp-server/main.go     (port 9091)
├── Terminal 4: go run cmd/grpc-server/main.go    (port 9092)
└── Terminal 5: yarn workspace @mangahub/web dev  (port 3000)

Database: ./data/mangahub.db (local SQLite file)
```

### Production (Docker Compose)

```yaml
version: "3.8"
services:
  api-server:
    build: .
    command: /app/bin/api-server
    ports:
      - "8080:8080" # HTTP API
      - "9093:9093" # WebSocket
    environment:
      - TCP_HOST=tcp-server
      - UDP_HOST=udp-server
      - GRPC_HOST=grpc-server
    volumes:
      - ./data:/app/data
    depends_on:
      - tcp-server
      - udp-server
      - grpc-server

  tcp-server:
    build: .
    command: /app/bin/tcp-server
    ports:
      - "9090:9090"
    volumes:
      - ./data:/app/data

  udp-server:
    build: .
    command: /app/bin/udp-server
    ports:
      - "9091:9091/udp"

  grpc-server:
    build: .
    command: /app/bin/grpc-server
    ports:
      - "9092:9092"
    volumes:
      - ./data:/app/data

  web:
    build: ./apps/web
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://api-server:8080
      - NEXT_PUBLIC_WS_URL=ws://api-server:9093
    depends_on:
      - api-server
```

**Benefits:**

- All services start with `docker-compose up`
- Internal networking (services can resolve by name)
- Shared volume for database
- Easy scaling and replication

---

## 11. Technology Stack

### Backend (Go)

| Component        | Technology                   | Purpose                         |
| ---------------- | ---------------------------- | ------------------------------- |
| HTTP Framework   | `gin-gonic/gin`              | REST API routing and middleware |
| WebSocket        | `gorilla/websocket`          | WebSocket connections           |
| gRPC             | `google.golang.org/grpc`     | RPC framework                   |
| Protocol Buffers | `google.golang.org/protobuf` | Binary serialization            |
| Database Driver  | `mattn/go-sqlite3`           | SQLite3 database driver         |
| JWT              | `golang-jwt/jwt/v4`          | Token generation and validation |
| Password Hashing | `golang.org/x/crypto/bcrypt` | Secure password hashing         |
| Configuration    | `gopkg.in/yaml.v3`           | YAML config parsing             |

### Frontend (TypeScript)

| Component        | Technology           | Purpose                  |
| ---------------- | -------------------- | ------------------------ |
| Framework        | Next.js 16           | React framework with SSR |
| UI Library       | React 19             | Component library        |
| Styling          | Tailwind CSS 4       | Utility-first CSS        |
| API Client       | `@mangahub/api`      | Type-safe HTTP client    |
| State Management | React Query          | Server state management  |
| Forms            | React Hook Form      | Form validation          |
| WebSocket        | Native WebSocket API | Real-time chat           |

### DevOps

| Tool            | Purpose                         |
| --------------- | ------------------------------- |
| Make            | Build automation (Go)           |
| Yarn Workspaces | Monorepo package management     |
| Turborepo       | Build orchestration and caching |
| Docker          | Containerization                |
| Docker Compose  | Multi-container orchestration   |
| GitHub Actions  | CI/CD (planned)                 |

---

## 12. Scalability Considerations

### Current Design (50-100 concurrent users)

**Bottlenecks:**

- SQLite write concurrency (single writer)
- All services on one machine
- No load balancing

**Sufficient because:**

- Academic project scope
- Small user base
- Simple deployment
- Educational focus

### Production Scaling (1000+ users)

**Database:**

- Migrate to PostgreSQL for concurrent writes
- Implement read replicas
- Add Redis for caching and sessions

**API Server:**

- Horizontal scaling with load balancer (Nginx, HAProxy)
- Stateless design allows multiple instances
- Shared Redis for session storage

**TCP Server:**

- Multiple instances with connection routing
- Use message queue (RabbitMQ, Kafka) for broadcasts
- Sticky sessions to maintain client connections

**UDP Server:**

- Multiple instances with anycast networking
- No state to share (pure broadcaster)

**WebSocket:**

- Use Redis pub/sub for message distribution across instances
- WebSocket load balancer with sticky sessions
- Consider Socket.IO or similar for reconnection

**gRPC:**

- Horizontal scaling with load balancing
- Service mesh (Istio, Linkerd) for advanced routing

---

## 13. Monitoring and Observability

### Current (Academic)

- Console logging with timestamps
- Health check endpoint: `GET /health`
- Basic error messages

### Production Recommendations

**Logging:**

- Structured logging (JSON format)
- Log aggregation (ELK stack, Loki)
- Log levels (DEBUG, INFO, WARN, ERROR)

**Metrics:**

- Prometheus for metrics collection
- Grafana for visualization
- Key metrics:
  - Requests per second
  - Response times (p50, p95, p99)
  - Active connections (TCP, WebSocket)
  - Database query times
  - Error rates

**Tracing:**

- Distributed tracing (Jaeger, Zipkin)
- Trace requests across services
- Identify performance bottlenecks

**Alerting:**

- PagerDuty, Opsgenie for incident management
- Alert on: high error rates, slow responses, service downtime

---

## 14. Testing Strategy

### Unit Tests

- Test individual functions in isolation
- Mock dependencies (database, external services)
- Coverage target: >80%

### Integration Tests

- Test interaction between components
- Use test database (in-memory SQLite)
- Test HTTP endpoints end-to-end

### Protocol Tests

- TCP client connects, authenticates, receives broadcasts
- UDP client registers, receives notifications
- WebSocket client joins room, sends/receives messages
- gRPC client calls services, validates responses

### End-to-End Tests

- Test complete user flows
- Register → Login → Search → Add to Library → Update Progress
- Verify cross-protocol integration (HTTP update → TCP broadcast)

---

## 15. Future Enhancements

### Planned (Bonus Features)

- [ ] Docker Compose setup (10 bonus points) ✅ Completed
- [ ] Advanced search filters
- [ ] User reviews and ratings
- [ ] Friend system
- [ ] Reading statistics dashboard
- [ ] Redis caching layer
- [ ] CI/CD pipeline

### Beyond Academic Scope

- Mobile applications (Flutter, React Native)
- Social features (following users, recommendations)
- Reading lists and collections
- Manga recommendation engine (ML-based)
- Multi-language support (i18n)
- Admin dashboard
- Email notifications
- PWA support for offline reading

---

## 16. References

### Documentation

- [API Documentation](./api-documentation.md)
- [TCP Documentation](./tcp-documentation.md)
- [UDP Documentation](./udp-documentation.md)
- [WebSocket Documentation](./websocket-documentation.md)
- [gRPC Documentation](./grpc-documentation.md)
- [CLAUDE.md](../CLAUDE.md) - AI Assistant Context
- [MONOREPO.md](../MONOREPO.md) - Monorepo Structure

### External Resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [gRPC Go Docs](https://grpc.io/docs/languages/go/)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)
- [Gin Framework](https://gin-gonic.com/docs/)

### Academic

- [Project Specification PDF](../mangahub_project_spec.pdf)
- [Use Case Specification PDF](<../mangahub_usecase%20(reference).pdf>)
- [CLI Manual PDF](<../mangahub_cli_manual%20(reference).pdf>)

---

**Last Updated**: 2025-12-26
**Version**: 1.0.0
**Status**: ✅ Complete Implementation
**Course**: IT096IU - Network Programming
**Team**: MangaHub Development Team
