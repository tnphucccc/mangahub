# MangaHub API Documentation

This folder contains everything frontend developers need to integrate with the MangaHub backend.

## 📁 Contents

```
api/
├── openapi.yaml                    # OpenAPI 3.0 specification
├── README.md                       # This file - full API documentation
├── QUICKSTART.md                   # 5-minute quick start guide
├── typescript/                     # TypeScript types
│   ├── generated.ts                # Auto-generated from OpenAPI
│   └── index.ts                    # Convenient type exports
├── examples/                       # Ready-to-use code samples
│   ├── fetch-client.ts             # HTTP API client (fetch-based)
│   └── react-hooks.ts              # React custom hooks
├── scripts/                        # Build scripts
│   └── generate-index.js           # Auto-generate index.ts
└── package.json                    # Development tools

```

---

## 🚀 Quick Start

### 1. **View Interactive API Documentation**

```bash
# From project root
make docs-preview

# OR from api/ folder
cd api
yarn install
yarn preview
```

This opens an interactive Redoc documentation at `http://localhost:8080`

### 2. **Generate TypeScript Types**

```bash
# From project root
make generate-types

# OR from api/ folder
cd api
yarn generate
```

Types are generated at `api/typescript/generated.ts`

### 3. **Start the Backend Server**

```bash
# Terminal 1: HTTP API Server
make run-api

# Terminal 2: TCP Server (for real-time progress sync)
make run-tcp

# Terminal 3: UDP Server (for notifications)
make run-udp

# Terminal 4: gRPC Server (for internal services)
make run-grpc
```

---

## 📖 API Overview

### Base URL

- **Development**: `http://localhost:8080/api/v1`
- **Production**: `https://your-domain.com/api/v1`

### Authentication

The API uses **JWT Bearer Token** authentication.

```typescript
// Example header
Authorization: Bearer <your-jwt-token>
```

---

## 🔑 Authentication Endpoints

### Register User

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "manga_lover",
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response** (201 Created):

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid-here",
      "username": "manga_lover",
      "email": "user@example.com",
      "created_at": "2024-12-20T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "manga_lover",
  "password": "securepassword123"
}
```

**Response** (200 OK):

```json
{
  "success": true,
  "data": {
    "user": { ... },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

## 📚 Manga Endpoints (Public)

### Search Manga

```http
GET /api/v1/manga?title=one+piece&genre=action&status=ongoing&limit=20&offset=0
```

**Response**:

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "manga-001",
        "title": "One Piece",
        "author": "Eiichiro Oda",
        "genres": ["Action", "Adventure", "Fantasy"],
        "status": "ongoing",
        "total_chapters": 1100,
        "description": "...",
        "cover_image_url": "https://...",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-12-20T00:00:00Z"
      }
    ]
  },
  "meta": {
    "total": 200,
    "count": 20,
    "limit": 20,
    "offset": 0,
    "has_more": true,
    "page": 1,
    "total_pages": 10
  }
}
```

### Get Manga Details

```http
GET /api/v1/manga/{manga_id}
```

### Get All Manga (Paginated)

```http
GET /api/v1/manga/all?limit=50&offset=0
```

---

## 👤 User Library Endpoints (Protected)

### Get User Library

```http
GET /api/v1/users/library
Authorization: Bearer <token>
```

**Response**:

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "user_id": "user-123",
        "manga_id": "manga-001",
        "current_chapter": 75,
        "status": "reading",
        "rating": 9,
        "started_at": "2024-01-15T00:00:00Z",
        "updated_at": "2024-12-20T10:00:00Z",
        "manga": {
          "id": "manga-001",
          "title": "One Piece",
          ...
        }
      }
    ]
  }
}
```

### Add Manga to Library

```http
POST /api/v1/users/library
Authorization: Bearer <token>
Content-Type: application/json

{
  "manga_id": "manga-001",
  "status": "reading",
  "current_chapter": 0
}
```

### Update Reading Progress

```http
PUT /api/v1/users/progress/{manga_id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "current_chapter": 75,
  "status": "reading",
  "rating": 9
}
```

### Get Progress for Specific Manga

```http
GET /api/v1/users/progress/{manga_id}
Authorization: Bearer <token>
```

---

## 🔗 Real-Time Connections

### TCP Connection (Progress Sync)

**Purpose**: Real-time progress synchronization across devices

**Port**: `9090`

**Protocol**: JSON messages over TCP with newline delimiters

**Connection Flow**:

1. Connect to `localhost:9090`
2. Send authentication message with JWT token
3. Receive auth success/failure
4. Send/receive progress updates
5. Receive broadcasts from other users

**Example messages**: See `openapi.yaml` for TCP message schemas

### UDP Connection (Notifications)

**Purpose**: Chapter release notifications

**Port**: `9091`

**Protocol**: JSON messages over UDP

**Connection Flow**:

1. Send registration message to `localhost:9091`
2. Receive registration success
3. Listen for notifications (fire-and-forget)
4. Send unregistration when done

**Example messages**: See `openapi.yaml` for UDP message schemas

**Note**: Browsers cannot use raw UDP. For frontend integration, you'll need a WebSocket bridge or Server-Sent Events (SSE) wrapper.

### WebSocket (Chat - Coming Soon)

**Purpose**: Real-time chat between users

**Port**: `9093`

**Endpoint**: `ws://localhost:9093/ws`

---

## 📦 Using in Your Frontend

### Option 1: TypeScript + Fetch

```typescript
import type {
  UserLoginRequest,
  AuthResponse,
  MangaListResponse,
} from "./api/typescript";

const API_BASE = "http://localhost:8080/api/v1";

async function login(credentials: UserLoginRequest): Promise<AuthResponse> {
  const response = await fetch(`${API_BASE}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(credentials),
  });

  if (!response.ok) throw new Error("Login failed");
  return response.json();
}

async function searchManga(query: string): Promise<MangaListResponse> {
  const response = await fetch(
    `${API_BASE}/manga?title=${encodeURIComponent(query)}`
  );

  if (!response.ok) throw new Error("Search failed");
  return response.json();
}
```

### Option 2: With API Client Class

See `api/examples/fetch-client.ts` for a complete implementation.

### Option 3: React Hooks

See `api/examples/react-hooks.ts` for React integration.

---

## 🛠️ Development Tools

### Regenerate Types

When the API changes, regenerate TypeScript types:

```bash
make generate-types
```

### Validate OpenAPI Spec

```bash
make docs-validate
```

### Preview Documentation

```bash
make docs-preview
```

---

## 📝 Error Handling

All errors follow this format:

```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid or expired token"
  }
}
```

**Common Error Codes**:

- `BAD_REQUEST` - Invalid request body or parameters (400)
- `UNAUTHORIZED` - Missing or invalid authentication (401)
- `NOT_FOUND` - Resource not found (404)
- `CONFLICT` - Resource already exists (409)
- `INTERNAL_ERROR` - Server error (500)

---

## 🧪 Testing the API

### Using cURL

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Login
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' \
  | jq -r '.data.token')

# Search manga
curl http://localhost:8080/api/v1/manga?title=naruto

# Get library (authenticated)
curl http://localhost:8080/api/v1/users/library \
  -H "Authorization: Bearer $TOKEN"
```

### Using Postman

1. Import `openapi.yaml` into Postman
2. Postman will auto-generate a collection
3. Set up environment variables for base URL and token

---

## 📚 Type Reference

All TypeScript types are available in `typescript/index.ts`:

### Entity Types

- `User` - User account
- `Manga` - Manga details
- `UserProgress` - Reading progress
- `UserProgressWithManga` - Progress with manga details

### Enum Types

- `MangaStatus` - `"ongoing" | "completed" | "hiatus" | "cancelled"`
- `ReadingStatus` - `"reading" | "completed" | "plan_to_read" | "on_hold" | "dropped"`

### Request Types

- `UserRegisterRequest`
- `UserLoginRequest`
- `LibraryAddRequest`
- `ProgressUpdateRequest`

### Response Types

- `APIResponse` - Base response structure
- `AuthResponse`
- `MangaListResponse`
- `LibraryResponse`
- `ProgressResponse`

---

## 🔐 Security Notes

### Development

- CORS is enabled for all origins (`*`)
- JWT tokens expire in 7 days (configurable)

### Production

- Set `CORS_ALLOWED_ORIGINS` environment variable
- Use HTTPS only
- Set a strong `JWT_SECRET`
- Consider shorter token expiry

---

## 📞 Support

- **Issues**: File a bug report on GitHub
- **Questions**: Check the interactive documentation (`make docs-preview`)
- **API Changes**: See `CHANGELOG.md` in project root

---

## 🔄 Version History

- **v1.0.0** (2024-12-20) - Initial release
  - Authentication endpoints
  - Manga catalog search
  - User library management
  - Progress tracking
  - TCP/UDP real-time protocols

---

## 📄 License

See main project README for license information.
