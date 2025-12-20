# MangaHub Monorepo Structure

This project uses a **Yarn Workspaces + Turborepo** monorepo setup to manage both the Go backend and TypeScript/JavaScript packages in a single repository.

## 📁 Project Structure

```
mangahub/
├── packages/                    # Shared TypeScript packages
│   ├── spec/                   # @mangahub/spec - OpenAPI specification
│   ├── types/                  # @mangahub/types - Generated TypeScript types
│   ├── api/                    # @mangahub/api - HTTP client SDK
│   └── hooks/                  # @mangahub/hooks - React hooks
│
├── apps/                       # Applications
│   └── web/                    # @mangahub/web - Next.js web application
│
├── cmd/                        # Go server binaries
│   ├── api-server/            # HTTP REST API server
│   ├── tcp-server/            # TCP sync server
│   ├── udp-server/            # UDP notification server
│   ├── grpc-server/           # gRPC service
│   └── cli/                   # CLI tool
│
├── internal/                   # Go internal packages
│   ├── auth/                  # Authentication & JWT
│   ├── manga/                 # Manga business logic
│   ├── user/                  # User management
│   ├── tcp/                   # TCP server implementation
│   ├── udp/                   # UDP server implementation
│   └── grpc/                  # gRPC server implementation
│
├── pkg/                        # Go shared libraries
│   ├── config/                # Configuration
│   ├── database/              # Database utilities
│   ├── models/                # Data models
│   └── utils/                 # Helper functions
│
├── docs/                       # Documentation
├── data/                       # Manga data & database
├── migrations/                 # SQL migrations
├── proto/                      # Protocol Buffer definitions
└── test/                       # Tests
```

## 📦 Packages

### Shared Packages (`packages/`)

#### `@mangahub/spec`
- **Description:** OpenAPI 3.0 specification and API documentation
- **Location:** `packages/spec/`
- **Main Files:** `openapi.yaml`, `README.md`, `QUICKSTART.md`
- **Commands:**
  ```bash
  yarn workspace @mangahub/spec preview    # Preview API docs
  yarn workspace @mangahub/spec validate   # Validate spec
  ```

#### `@mangahub/types`
- **Description:** Auto-generated TypeScript types from OpenAPI spec
- **Location:** `packages/types/`
- **Generated From:** `@mangahub/spec/openapi.yaml`
- **Commands:**
  ```bash
  yarn workspace @mangahub/types generate  # Generate types
  yarn workspace @mangahub/types build     # Compile types
  ```

#### `@mangahub/api`
- **Description:** Type-safe HTTP API client SDK
- **Location:** `packages/api/`
- **Dependencies:** `@mangahub/types`
- **Exports:** `MangaHubClient`, `mangaHubApi`
- **Commands:**
  ```bash
  yarn workspace @mangahub/api build       # Build client
  yarn workspace @mangahub/api typecheck   # Type-check
  ```

#### `@mangahub/hooks`
- **Description:** React hooks for API integration
- **Location:** `packages/hooks/`
- **Dependencies:** `@mangahub/api`, `@mangahub/types`
- **Exports:** `useAuth`, `useMangaSearch`, `useLibrary`, etc.
- **Peer Dependencies:** React 18+ or 19+
- **Commands:**
  ```bash
  yarn workspace @mangahub/hooks build     # Build hooks
  ```

### Applications (`apps/`)

#### `@mangahub/web`
- **Description:** Next.js 16 web application
- **Location:** `apps/web/`
- **Framework:** Next.js 16.1.0 + React 19 + Tailwind CSS 4
- **Dependencies:** `@mangahub/api`, `@mangahub/hooks`, `@mangahub/types`
- **Commands:**
  ```bash
  yarn workspace @mangahub/web dev         # Development server
  yarn workspace @mangahub/web build       # Production build
  yarn workspace @mangahub/web start       # Start production server
  ```

## 🚀 Getting Started

### Prerequisites

- **Node.js** >= 18.0.0
- **Yarn** >= 4.0.0
- **Go** >= 1.19 (for backend)

### Installation

```bash
# Install all dependencies (JavaScript + Go)
yarn install

# Or use the Makefile
make js-install
```

### Development Workflow

#### 1. Start Go Backend Servers

```bash
# Terminal 1: HTTP API Server
make run-api

# Terminal 2: TCP Server
make run-tcp

# Terminal 3: UDP Server
make run-udp

# Terminal 4: gRPC Server
make run-grpc
```

#### 2. Start Next.js Frontend

```bash
# Terminal 5: Web application
yarn workspace @mangahub/web dev
# Or
make js-dev
```

Access the web app at: http://localhost:3000

### Building Everything

```bash
# Build Go backend
make build

# Build JavaScript packages
yarn build
# Or
make js-build

# Build everything
make build && make js-build
```

## 🔧 Common Tasks

### Code Generation

```bash
# Generate TypeScript types from OpenAPI
yarn workspace @mangahub/types generate
# Or
make generate-types

# Generate gRPC code from Protocol Buffers
make generate-proto

# Generate all
make generate
```

### Database Operations

```bash
# Run migrations
make migrate-up

# Seed database
make seed

# Reset database
make db-reset
```

### Testing

```bash
# Test Go backend
make test

# Test JavaScript packages
yarn test

# Run with coverage
make test-coverage
```

### Linting & Formatting

```bash
# Lint Go code
make lint

# Format Go code
make fmt

# Lint JavaScript code
yarn lint

# Format JavaScript code
yarn format
```

## 📝 Package Dependencies

```
@mangahub/spec (OpenAPI spec)
       ↓
@mangahub/types (Generated types)
       ↓
@mangahub/api (HTTP client)
       ↓
@mangahub/hooks (React hooks)
       ↓
@mangahub/web (Next.js app)
```

## 🏗️ Build System

This monorepo uses:

- **Yarn Workspaces** - Package management and workspace linking
- **Turborepo** - Smart task orchestration and caching
- **Make** - Go build automation

### Turborepo Tasks

The `turbo.json` file defines the build pipeline:

- `build` - Compile packages
- `test` - Run tests
- `lint` - Lint code
- `typecheck` - Type checking
- `generate` - Code generation
- `dev` - Development servers

Tasks automatically respect dependencies (e.g., `@mangahub/api` builds after `@mangahub/types`).

### Caching

Turborepo caches task outputs in `.turbo/` for instant rebuilds when nothing changed.

## 📚 Adding a New Package

### 1. Create Package Directory

```bash
mkdir -p packages/new-package/src
cd packages/new-package
```

### 2. Create `package.json`

```json
{
  "name": "@mangahub/new-package",
  "version": "1.0.0",
  "private": true,
  "main": "dist/index.js",
  "types": "dist/index.d.ts",
  "scripts": {
    "build": "tsc",
    "clean": "rm -rf dist"
  },
  "dependencies": {
    "@mangahub/types": "workspace:*"
  },
  "devDependencies": {
    "typescript": "^5.9.3"
  }
}
```

### 3. Create `tsconfig.json`

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "declaration": true,
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true
  },
  "include": ["src/**/*"]
}
```

### 4. Install Dependencies

```bash
yarn install
```

The new package will automatically be part of the workspace.

## 🌐 Deployment

### Docker (Recommended)

```bash
# Build all services
docker-compose build

# Start all services
docker-compose up
```

### Manual Deployment

#### Backend (Go)

```bash
# Build binaries
make build

# Deploy binaries to server
scp bin/* user@server:/opt/mangahub/

# Start services (systemd, supervisor, etc.)
```

#### Frontend (Next.js)

```bash
# Build for production
yarn workspace @mangahub/web build

# Start production server
yarn workspace @mangahub/web start

# Or deploy to Vercel/Netlify
```

## 🧪 Testing Strategy

### Go Tests

```bash
# Unit tests
make test

# Integration tests
make test-integration

# With coverage
make test-coverage
```

### JavaScript Tests

```bash
# Run all tests
yarn test

# Test specific package
yarn workspace @mangahub/api test
```

## 🔍 Troubleshooting

### Yarn Workspace Issues

```bash
# Clean and reinstall
yarn clean
rm -rf node_modules
yarn install
```

### Type Generation Issues

```bash
# Regenerate types
yarn workspace @mangahub/types generate
```

### Build Cache Issues

```bash
# Clear Turborepo cache
rm -rf .turbo
yarn build
```

## 📖 Documentation

- **API Documentation:** See `packages/spec/README.md`
- **Frontend Integration:** See `docs/FRONTEND_INTEGRATION_GUIDE.md`
- **Architecture:** See `CLAUDE.md`
- **Project Specification:** See `docs/project_specification.md`

## 🛠️ Available Make Commands

```bash
make help                 # Show all commands

# Go Backend
make build                # Build all Go binaries
make run-api              # Run HTTP API server
make run-tcp              # Run TCP server
make run-udp              # Run UDP server
make run-grpc             # Run gRPC server
make test                 # Run Go tests

# JavaScript/TypeScript
make js-install           # Install dependencies
make js-build             # Build all packages
make js-dev               # Run Next.js dev server
make js-test              # Run JS tests
make js-lint              # Lint JS code
make js-typecheck         # Type-check JS code
make js-clean             # Clean JS build artifacts

# Code Generation
make generate-types       # Generate TypeScript types
make generate-proto       # Generate gRPC code
make generate             # Generate all

# Database
make migrate-up           # Run migrations
make seed                 # Seed database
make db-reset             # Reset and reseed

# Documentation
make docs-preview         # Preview OpenAPI docs
make docs-validate        # Validate OpenAPI spec

# Cleanup
make clean                # Clean Go artifacts
make clean-all            # Clean everything
```

## 📜 License

MIT License - See LICENSE file for details

## 👥 Contributors

MangaHub Team - Network Programming Course (IT096IU)

---

**Last Updated:** 2025-12-21
