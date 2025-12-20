# Frontend Integration Guide

## 📋 Summary: How the `api/` Folder Works

The `api/` folder is your **Frontend Developer SDK** - it contains everything needed to build a UI for MangaHub:

```
api/
├── openapi.yaml          # Complete API documentation (OpenAPI 3.0)
├── typescript/           # Auto-generated TypeScript types
│   ├── generated.ts      # Generated from openapi.yaml
│   └── index.ts          # Convenient exports
├── examples/             # Ready-to-use code samples
│   ├── fetch-client.ts   # Full API client implementation
│   └── react-hooks.ts    # React custom hooks
├── README.md             # Full API documentation
├── QUICKSTART.md         # 5-minute quick start
└── package.json          # Tools for type generation
```

---

## 🎯 Purpose of Each Component

### 1. **openapi.yaml** - The Source of Truth

**What it is**: OpenAPI 3.0 specification documenting all REST endpoints, request/response schemas, and real-time protocols.

**What it's used for**:
- ✅ Auto-generate TypeScript types
- ✅ Generate interactive API documentation
- ✅ Import into Postman/Insomnia
- ✅ Validate API requests/responses
- ✅ Single source of truth for API contract

**How to use**:
```bash
# View interactive docs
make docs-preview

# Validate specification
make docs-validate

# Generate types from it
make generate-types
```

### 2. **typescript/** - Type-Safe Development

**What it is**: Auto-generated TypeScript types from openapi.yaml + convenience exports.

**What it's used for**:
- ✅ Type-safe API calls in TypeScript/JavaScript projects
- ✅ Autocomplete in VSCode/other IDEs
- ✅ Catch errors at compile-time, not runtime
- ✅ Document data structures

**How to use**:
```typescript
// Import types
import type { Manga, User, AuthResponse } from './api/typescript';

// Use in your code with full type safety
const user: User = {
  id: 'abc123',
  username: 'manga_lover',
  email: 'user@example.com',
  created_at: new Date().toISOString()
};

// TypeScript will error if you forget required fields!
```

**Regenerate when API changes**:
```bash
# Manual
make generate-types

# Auto-watch mode
cd api && yarn generate:watch
```

### 3. **examples/** - Copy-Paste Ready Code

**What it is**: Production-ready code samples that frontend developers can copy directly into their projects.

**What's included**:
- `fetch-client.ts` - Complete API client class with authentication, error handling, and all endpoints
- `react-hooks.ts` - Custom React hooks for auth, search, library management, etc.

**How to use**:

#### Option A: Copy the client
```bash
# Copy to your project
cp api/examples/fetch-client.ts your-project/src/lib/api-client.ts
cp -r api/typescript your-project/src/types/mangahub/
```

```typescript
// In your app
import { mangaHubApi } from './lib/api-client';

// Login
await mangaHubApi.login({ username, password });

// Search
const results = await mangaHubApi.searchManga({ title: 'naruto' });
```

#### Option B: Use React hooks
```bash
cp api/examples/react-hooks.ts your-react-app/src/hooks/useMangaHub.ts
```

```tsx
// In React components
function App() {
  const { user, login } = useAuth();
  const { manga, search } = useMangaSearch();

  return <div>...</div>;
}
```

---

## 🚀 For Frontend Developers: Step-by-Step

### **Phase 1: Get Started (5 minutes)**

1. **Start the backend**:
   ```bash
   make run-api  # HTTP API on port 8080
   ```

2. **Test the API**:
   ```bash
   curl http://localhost:8080/health
   # Should return: {"status":"ok","service":"MangaHub HTTP API Server","version":"1.0.0"}
   ```

3. **View API docs**:
   ```bash
   make docs-preview
   # Opens http://localhost:8080 with interactive docs
   ```

### **Phase 2: Copy Types to Your Project**

**For TypeScript/React/Vue/Svelte projects**:

```bash
# Copy TypeScript types
cp -r api/typescript your-frontend/src/types/mangahub/

# Or link them during development
ln -s $(pwd)/api/typescript your-frontend/src/types/mangahub
```

**Update as API changes**:
```bash
make generate-types  # Regenerates typescript/generated.ts
```

### **Phase 3: Choose Integration Method**

#### **Method 1: Simple Fetch (Vanilla JS/TS)**

```typescript
import type { AuthResponse } from './types/mangahub';

const API = 'http://localhost:8080/api/v1';

async function login(username: string, password: string): Promise<AuthResponse> {
  const res = await fetch(`${API}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
  });
  return res.json();
}

const { data } = await login('testuser', 'password123');
console.log('Token:', data.token);
localStorage.setItem('token', data.token);
```

#### **Method 2: Use API Client Class (Recommended)**

```bash
# Copy the pre-built client
cp api/examples/fetch-client.ts your-project/src/lib/api-client.ts
```

```typescript
import { mangaHubApi } from './lib/api-client';

// Login (token auto-saved)
await mangaHubApi.login({ username: 'test', password: 'pass' });

// All methods are type-safe and authenticated automatically
const library = await mangaHubApi.getLibrary();
const results = await mangaHubApi.searchManga({ title: 'naruto' });
await mangaHubApi.updateProgress('manga-001', { current_chapter: 75 });
```

#### **Method 3: React Hooks (Best for React)**

```bash
cp api/examples/react-hooks.ts your-react-app/src/hooks/useMangaHub.ts
cp api/examples/fetch-client.ts your-react-app/src/lib/api-client.ts
```

```tsx
import { useAuth, useMangaSearch } from './hooks/useMangaHub';

function App() {
  const { user, login, logout, isAuthenticated } = useAuth();
  const { manga, search, isLoading } = useMangaSearch();

  return (
    <div>
      {isAuthenticated ? (
        <h1>Welcome, {user?.username}!</h1>
      ) : (
        <button onClick={() => login({ username: 'test', password: 'pass' })}>
          Login
        </button>
      )}
    </div>
  );
}
```

### **Phase 4: Handle Authentication**

```typescript
// The API client handles this automatically!

// Login
const { data } = await mangaHubApi.login({ username, password });
// Token is saved to localStorage automatically

// Make authenticated requests
const library = await mangaHubApi.getLibrary();
// Token is included in Authorization header automatically

// Check if authenticated
if (mangaHubApi.isAuthenticated()) {
  console.log('User is logged in');
}

// Logout
mangaHubApi.logout();
// Token is cleared from localStorage
```

### **Phase 5: Error Handling**

```typescript
try {
  await mangaHubApi.login({ username, password });
} catch (error: any) {
  // All errors follow this format
  console.error('Code:', error.code);      // e.g., "UNAUTHORIZED"
  console.error('Message:', error.message); // e.g., "Invalid credentials"
  console.error('Status:', error.status);   // e.g., 401

  // Common error codes:
  // - BAD_REQUEST (400)
  // - UNAUTHORIZED (401)
  // - NOT_FOUND (404)
  // - CONFLICT (409)
  // - INTERNAL_ERROR (500)
}
```

---

## 📊 Complete API Endpoint Reference

### **Public Endpoints (No Authentication)**

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login user |
| GET | `/api/v1/manga` | Search manga (with filters) |
| GET | `/api/v1/manga/all` | Get all manga (paginated) |
| GET | `/api/v1/manga/:id` | Get manga details |
| GET | `/health` | Health check |

### **Protected Endpoints (Require JWT Token)**

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users/me` | Get current user profile |
| GET | `/api/v1/users/library` | Get user's manga library |
| POST | `/api/v1/users/library` | Add manga to library |
| GET | `/api/v1/users/progress/:manga_id` | Get progress for manga |
| PUT | `/api/v1/users/progress/:manga_id` | Update reading progress |

### **Real-Time Protocols**

| Protocol | Port | Purpose |
|----------|------|---------|
| TCP | 9090 | Real-time progress sync across devices |
| UDP | 9091 | Chapter release notifications |
| WebSocket | 9093 | Chat (coming soon) |
| gRPC | 9092 | Internal services only |

---

## 🎨 Frontend Project Structure Recommendation

```
your-frontend-project/
├── src/
│   ├── lib/
│   │   └── api-client.ts        # Copied from api/examples/fetch-client.ts
│   ├── hooks/                   # (React only)
│   │   └── useMangaHub.ts       # Copied from api/examples/react-hooks.ts
│   ├── types/
│   │   └── mangahub/            # Copied from api/typescript/
│   │       ├── generated.ts
│   │       └── index.ts
│   ├── components/
│   │   ├── MangaCard.tsx
│   │   ├── SearchBar.tsx
│   │   └── Library.tsx
│   ├── pages/
│   │   ├── Home.tsx
│   │   ├── Search.tsx
│   │   ├── Library.tsx
│   │   └── MangaDetail.tsx
│   └── App.tsx
└── package.json
```

---

## 🔧 Development Workflow

### **1. API Changes**

When the backend API changes:

```bash
# Backend developer updates openapi.yaml
cd mangahub/api
yarn validate  # Check spec is valid
yarn generate  # Regenerate TypeScript types

# Frontend developer pulls changes
cd your-frontend
git pull
cp -r ../mangahub/api/typescript src/types/mangahub/
# TypeScript will now show errors where types changed!
```

### **2. Environment Configuration**

```typescript
// config.ts
export const config = {
  apiUrl: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
  tcpUrl: import.meta.env.VITE_TCP_URL || 'localhost:9090',
  udpUrl: import.meta.env.VITE_UDP_URL || 'localhost:9091',
};

// Use in API client
import { config } from './config';

export const mangaHubApi = new MangaHubClient({
  baseUrl: config.apiUrl,
  onTokenExpired: () => {
    // Redirect to login
    window.location.href = '/login';
  },
  onError: (error) => {
    // Global error handling
    console.error('API Error:', error);
  }
});
```

### **3. CORS Setup**

```bash
# Development: Allow your frontend origin
export CORS_ALLOWED_ORIGINS="http://localhost:3000,http://localhost:5173"
make run-api

# Production: Set in configs/prod.yaml or environment variable
```

---

## 📦 Recommended npm Packages

For an optimal developer experience:

```json
{
  "devDependencies": {
    "typescript": "^5.0.0",          // Type checking
    "@types/node": "^20.0.0"         // Node types
  },
  "dependencies": {
    // React Query (optional, for data fetching)
    "@tanstack/react-query": "^5.0.0",

    // Axios (alternative to fetch)
    "axios": "^1.6.0",

    // Zod (runtime validation)
    "zod": "^3.22.0"
  }
}
```

---

## ✅ Checklist: Ready for Frontend Development

- [ ] Backend servers are running (`make run-api`)
- [ ] API health check passes (`curl http://localhost:8080/health`)
- [ ] TypeScript types copied to frontend project
- [ ] API client copied and configured
- [ ] Test login/register works
- [ ] Test search manga works
- [ ] Test authenticated endpoints work
- [ ] Error handling implemented
- [ ] CORS configured correctly
- [ ] Environment variables set

---

## 🆘 Common Issues & Solutions

### "CORS policy: No 'Access-Control-Allow-Origin'"

**Solution**: Set CORS allowed origins
```bash
export CORS_ALLOWED_ORIGINS="http://localhost:3000"
make run-api
```

### "401 Unauthorized" on protected endpoints

**Solution**: Make sure you're logged in and token is valid
```typescript
// Check authentication
console.log('Token:', mangaHubApi.getToken());
console.log('Is authenticated:', mangaHubApi.isAuthenticated());

// Login again if needed
await mangaHubApi.login({ username, password });
```

### "Cannot find module './api/typescript'"

**Solution**: Copy the types folder
```bash
cp -r api/typescript your-project/src/types/mangahub
```

### Types are outdated

**Solution**: Regenerate types
```bash
cd mangahub
make generate-types
# Then copy updated types to your frontend
```

---

## 📞 Getting Help

1. **Check the docs**: [api/README.md](../api/README.md)
2. **Try quick start**: [api/QUICKSTART.md](../api/QUICKSTART.md)
3. **View examples**: [api/examples/](../api/examples/)
4. **Interactive API docs**: `make docs-preview`
5. **File an issue**: GitHub Issues
6. **Ask the team**: Discord/Slack

---

## 🎓 Learning Resources

- [OpenAPI Specification](https://swagger.io/specification/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [React Query Docs](https://tanstack.com/query/latest)
- [Fetch API MDN](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API)

---

**Ready to build?** Start with the [QUICKSTART.md](../api/QUICKSTART.md)!
