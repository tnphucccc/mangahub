# Quick Start Guide for Frontend Developers

Get started with the MangaHub API in 5 minutes!

## 🚀 Step 1: Start the Backend

```bash
# Clone the repository (if you haven't)
git clone https://github.com/your-org/mangahub.git
cd mangahub

# Start all servers (in separate terminals or use a process manager)
make run-api   # Terminal 1: HTTP API on port 8080
make run-tcp   # Terminal 2: TCP server on port 9090
make run-udp   # Terminal 3: UDP server on port 9091
make run-grpc  # Terminal 4: gRPC server on port 9092
```

**Or use Docker Compose (recommended)**:
```bash
docker-compose up
```

## 📚 Step 2: Explore the API

Open interactive documentation:
```bash
make docs-preview
# Opens at http://localhost:8080
```

Or test with cURL:
```bash
# Health check
curl http://localhost:8080/health

# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Login and save token
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' \
  | jq -r '.data.token')

# Search manga
curl "http://localhost:8080/api/v1/manga?title=naruto"

# Get your library (authenticated)
curl http://localhost:8080/api/v1/users/library \
  -H "Authorization: Bearer $TOKEN"
```

## 🎨 Step 3: Use in Your Frontend

### Option A: Copy TypeScript Types

```bash
# From mangahub/api directory, copy the typescript folder to your project
cp -r api/typescript your-frontend-project/src/api-types/
```

In your frontend:
```typescript
import type { Manga, User, AuthResponse } from './api-types';

const API_BASE = 'http://localhost:8080/api/v1';

async function login(username: string, password: string): Promise<AuthResponse> {
  const response = await fetch(`${API_BASE}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
  return response.json();
}
```

### Option B: Use the Pre-built API Client

```bash
# Copy the client to your project
cp api/examples/fetch-client.ts your-frontend-project/src/lib/api-client.ts
cp -r api/typescript your-frontend-project/src/api-types/
```

In your frontend:
```typescript
import { mangaHubApi } from './lib/api-client';

// Login
const { data } = await mangaHubApi.login({
  username: 'testuser',
  password: 'password123'
});

console.log('Token:', data.token);

// Search manga
const results = await mangaHubApi.searchManga({ title: 'naruto' });
console.log('Found:', results.meta.total, 'manga');
```

### Option C: React Hooks (Best for React apps)

```bash
# Copy the hooks and client
cp api/examples/react-hooks.ts your-react-app/src/hooks/useMangaHub.ts
cp api/examples/fetch-client.ts your-react-app/src/lib/api-client.ts
cp -r api/typescript your-react-app/src/api-types/
```

In your React components:
```tsx
import { useAuth, useMangaSearch, useLibrary } from './hooks/useMangaHub';

function App() {
  const { user, login, logout } = useAuth();
  const { manga, search } = useMangaSearch();
  const { library } = useLibrary();

  return (
    <div>
      {user ? (
        <>
          <h1>Welcome, {user.username}!</h1>
          <button onClick={logout}>Logout</button>
        </>
      ) : (
        <LoginForm onLogin={login} />
      )}
    </div>
  );
}
```

## 🔐 Step 4: Authentication Flow

```typescript
// 1. Login and get token
const loginResponse = await mangaHubApi.login({
  username: 'testuser',
  password: 'password123'
});

// Token is automatically saved to localStorage
console.log('Logged in!', loginResponse.data.user);

// 2. Make authenticated requests
const library = await mangaHubApi.getLibrary();
console.log('My library:', library.data.items);

// 3. Token is automatically included in requests
// No need to manually add Authorization header!

// 4. Logout when done
mangaHubApi.logout();
```

## 📖 Step 5: Common Use Cases

### Search for Manga
```typescript
const results = await mangaHubApi.searchManga({
  title: 'naruto',
  genre: 'action',
  status: 'ongoing',
  limit: 20,
  offset: 0
});

results.data.items.forEach(manga => {
  console.log(`${manga.title} by ${manga.author}`);
});
```

### Add to Library
```typescript
await mangaHubApi.addToLibrary({
  manga_id: 'manga-001',
  status: 'reading',
  current_chapter: 0
});
```

### Update Progress
```typescript
await mangaHubApi.updateProgress('manga-001', {
  current_chapter: 75,
  status: 'reading',
  rating: 9
});
```

### Get Manga Details
```typescript
const manga = await mangaHubApi.getManga('manga-001');
console.log('Manga:', manga.data.manga);
```

## 🎯 Step 6: Real-Time Features (Optional)

### TCP Progress Sync
For real-time progress updates across devices, connect to the TCP server:

```typescript
// Note: You'll need a WebSocket-to-TCP bridge or use raw TCP in Node.js
const socket = new WebSocket('ws://localhost:9090'); // Example

socket.onopen = () => {
  // Send auth message
  socket.send(JSON.stringify({
    type: 'auth',
    timestamp: new Date().toISOString(),
    data: { token: yourJWTToken }
  }));
};

socket.onmessage = (event) => {
  const message = JSON.parse(event.data);

  if (message.type === 'broadcast') {
    console.log('Progress update from:', message.data.username);
    console.log('Reading:', message.data.manga_title, 'chapter', message.data.current_chapter);
  }
};
```

### UDP Notifications
For chapter release notifications, register with the UDP server:

```typescript
// Note: UDP requires native implementation or server-sent events bridge
// See openapi.yaml for UDP message schemas
```

## 🛠️ Development Tips

### Enable CORS for Your Frontend

```bash
# Set environment variable
export CORS_ALLOWED_ORIGINS="http://localhost:3000,http://localhost:5173"

# Then start the API server
make run-api
```

### Auto-Regenerate Types When API Changes

```bash
cd api
yarn generate:watch  # Watches openapi.yaml for changes
```

### Check API Health

```typescript
const health = await fetch('http://localhost:8080/health');
const status = await health.json();

if (status.status === 'ok') {
  console.log('API is ready!');
}
```

## 📱 Sample Projects

Check out these example integrations:

- **React**: `examples/react-app/` (Coming soon)
- **Vue**: `examples/vue-app/` (Coming soon)
- **Svelte**: `examples/svelte-app/` (Coming soon)
- **Vanilla JS**: `examples/vanilla-js/` (Coming soon)

## 🐛 Troubleshooting

### CORS Errors
```bash
# Make sure API server has CORS enabled
# Check configs/dev.yaml or set CORS_ALLOWED_ORIGINS env variable
```

### 401 Unauthorized
```typescript
// Token might be expired, login again
mangaHubApi.logout();
await mangaHubApi.login({ username, password });
```

### Connection Refused
```bash
# Make sure the API server is running
curl http://localhost:8080/health

# If not running:
make run-api
```

## 📚 Next Steps

- Read the full [API Documentation](./README.md)
- Check [OpenAPI Specification](./openapi.yaml)
- Browse [Code Examples](./examples/)
- Join the team Discord for help

## 💡 Pro Tips

1. **Use TypeScript** - You get autocomplete and type safety!
2. **Check the interactive docs** - `make docs-preview` is your friend
3. **Store tokens securely** - The client handles localStorage automatically
4. **Handle errors** - All API errors follow a consistent format
5. **Test with Postman** - Import `openapi.yaml` directly

---

Happy coding! 🚀

For questions, check the [README](./README.md) or file an issue on GitHub.
