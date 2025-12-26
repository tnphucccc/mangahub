# MangaHub Web Frontend Documentation

---

## 1. Overview

The MangaHub Web Frontend is a modern **Next.js 16** web application that provides a user-friendly interface for manga tracking and management. Built with React 19 and Tailwind CSS 4, it connects to the MangaHub backend API to deliver a seamless manga browsing and reading experience.

**Technology Stack:**

- **Framework**: Next.js 16.1.0 (App Router)
- **UI Library**: React 19.2.3
- **Language**: TypeScript 5
- **Styling**: Tailwind CSS 4
- **HTTP Client**: Axios 1.13.2
- **State Management**: React Context API
- **Package Manager**: Yarn 4.0.0

**Location**: `apps/web/`

---

## 2. Project Structure

```
apps/web/
├── src/
│   ├── app/                   # Next.js App Router pages
│   │   ├── page.tsx           # Landing page (login/register)
│   │   ├── layout.tsx         # Root layout with AuthProvider
│   │   ├── main/              # Main manga browse page
│   │   │   └── page.tsx
│   │   ├── library/           # User's manga library
│   │   │   └── page.tsx
│   │   ├── contexts/          # React Context providers
│   │   │   └── AuthContext.tsx
│   │   ├── hoc/               # Higher-Order Components
│   │   │   └── withAuth.tsx
│   │   └── helpers/           # Utility functions
│   │       ├── upperCaseFirstLetter.ts
│   │       └── hook.ts
│   │
│   ├── component/             # React components
│   │   ├── header.tsx         # Navigation header
│   │   ├── login.tsx          # Login form
│   │   ├── register.tsx       # Registration form
│   │   ├── searchAndFilter.tsx# Search and filter controls
│   │   ├── mangaCard.tsx      # Manga display card
│   │   ├── libraryMangaCard.tsx # Library-specific card
│   │   └── mangaModal.tsx     # Manga detail modal
│   │
│   └── lib/                   # Shared libraries
│       └── apiClient.ts       # HTTP API client
│
├── public/                    # Static assets
├── package.json               # Dependencies
├── next.config.js             # Next.js configuration
├── tailwind.config.ts         # Tailwind configuration
└── tsconfig.json              # TypeScript configuration
```

---

## 3. Key Features

### Implemented Features

✅ **User Authentication**

- Login with username/password
- User registration
- JWT token-based auth
- Auto-logout on token expiry
- Persistent login (localStorage)

✅ **Manga Browse**

- Search by title, author, genre
- Filter by manga status (ongoing, completed, etc.)
- Grid view of manga cards
- Cover images, titles, authors, genres

✅ **Manga Detail Modal**

- Full manga information
- Chapter progress tracking
- Add to library functionality
- "Add New Chapter" button (triggers notifications)

✅ **User Library**

- View all manga in library
- Update reading progress
- Change reading status (reading, completed, plan to read, etc.)
- Rate manga (1-10 stars)
- Remove from library

✅ **Real-Time Chat** (WebSocket integration ready)

- Chat interface scaffolded
- WebSocket connection support
- Message sending/receiving

### Planned Features (Future)

- 📖 Reading statistics dashboard
- 👥 User profiles
- ⭐ Manga reviews and ratings
- 🔔 Real-time notifications (WebSocket)
- 📱 Responsive mobile design
- 🌙 Dark mode toggle
- 🔍 Advanced search filters
- 📊 Reading analytics

---

## 4. Pages and Routes

### 4.1 Landing Page (`/`)

**File**: `src/app/page.tsx`

**Purpose**: Login and registration entry point

**Features**:

- Toggle between Login and Register forms
- Form validation
- Error handling
- Auto-redirect after successful login

**Components Used**:

- `<Login />` - Login form
- `<Register />` - Registration form

**Authentication Flow**:

```
User fills form → Submit → apiClient.login() →
Response with token → localStorage.setItem('token') →
AuthContext.login(token) → Redirect to /main
```

---

### 4.2 Main Page (`/main`)

**File**: `src/app/main/page.tsx`

**Purpose**: Browse and search manga catalog

**Features**:

- Search by title, author, genre
- Filter by status (ongoing, completed, hiatus, cancelled)
- Grid layout of manga cards
- Click card to open detail modal
- Loading states and error handling

**Components Used**:

- `<Header />` - Navigation with username, logout
- `<SearchAndFilter />` - Search and filter controls
- `<MangaCard />` - Individual manga display
- `<MangaModal />` - Manga detail popup

**State Management**:

```typescript
const [manga, setManga] = useState<Manga[]>([]);
const [loading, setLoading] = useState(false);
const [selectedManga, setSelectedManga] = useState<Manga | null>(null);
const [searchParams, setSearchParams] = useState({
  title: "",
  author: "",
  genre: "",
  status: undefined,
});
```

**Data Flow**:

```
User searches → Update searchParams →
apiClient.getManga(searchParams) → setManga(results) →
Render <MangaCard /> grid
```

---

### 4.3 Library Page (`/library`)

**File**: `src/app/library/page.tsx`

**Purpose**: Manage user's manga library

**Features**:

- View all added manga
- Update reading progress (chapter number)
- Change reading status (reading, completed, etc.)
- Rate manga (1-10 stars)
- Remove from library
- Real-time progress updates

**Components Used**:

- `<Header />` - Navigation
- `<LibraryMangaCard />` - Enhanced card with progress controls

**State Management**:

```typescript
const [library, setLibrary] = useState<UserProgressWithManga[]>([]);
const [loading, setLoading] = useState(true);
```

**Progress Update Flow**:

```
User updates chapter/status →
apiClient.updateMangaProgress(mangaId, data) →
Re-fetch library → Update UI
```

---

## 5. Key Components

### 5.1 `<Header />`

**File**: `src/component/header.tsx`

**Purpose**: Navigation bar with user info and actions

**Features**:

- Display username
- Navigation links (Browse, Library)
- Logout button
- Active route highlighting

**Props**: None (uses AuthContext)

**Usage**:

```tsx
import Header from "@/component/header";

export default function Page() {
  return (
    <div>
      <Header />
      {/* Page content */}
    </div>
  );
}
```

---

### 5.2 `<MangaCard />`

**File**: `src/component/mangaCard.tsx`

**Purpose**: Display manga in grid view

**Features**:

- Cover image with fallback
- Title, author, genres
- Click to open modal

**Props**:

```typescript
interface MangaCardProps {
  manga: Manga;
  onClick: (manga: Manga) => void;
}
```

**Usage**:

```tsx
<MangaCard manga={manga} onClick={(m) => setSelectedManga(m)} />
```

---

### 5.3 `<MangaModal />`

**File**: `src/component/mangaModal.tsx`

**Purpose**: Detailed manga view in modal

**Features**:

- Full manga information
- Cover image
- Genres, author, description
- Add to library button
- "Add New Chapter" button (sends UDP notification)
- Close button

**Props**:

```typescript
interface MangaModalProps {
  manga: Manga | null;
  onClose: () => void;
  isInLibrary: boolean;
  onAddToLibrary: (mangaId: string) => void;
}
```

**Notable Feature - UDP Notification**:

```typescript
// Send new chapter notification
const handleAddChapter = async () => {
  await fetch("http://localhost:8080/api/v1/admin/notifications", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      manga_id: manga.id,
      manga_title: manga.title,
      chapter_number: newChapterNumber,
      message: `New chapter of ${manga.title} released!`,
    }),
  });
};
```

This triggers the UDP notification server to broadcast to all registered clients!

---

### 5.4 `<LibraryMangaCard />`

**File**: `src/component/libraryMangaCard.tsx`

**Purpose**: Enhanced manga card for library with progress controls

**Features**:

- All features of MangaCard
- Current chapter input
- Status dropdown (reading, completed, etc.)
- Rating input (1-10 stars)
- Update button
- Remove button

**Props**:

```typescript
interface LibraryMangaCardProps {
  item: UserProgressWithManga;
  onUpdate: (mangaId: string, data: ProgressUpdateRequest) => void;
  onRemove: (mangaId: string) => void;
}
```

---

### 5.5 `<SearchAndFilter />`

**File**: `src/component/searchAndFilter.tsx`

**Purpose**: Search and filter controls

**Features**:

- Title search input
- Author search input
- Genre search input
- Status filter dropdown
- Search button

**Props**:

```typescript
interface SearchAndFilterProps {
  onSearch: (params: GetMangaParams) => void;
}
```

---

## 6. Authentication System

### 6.1 AuthContext

**File**: `src/app/contexts/AuthContext.tsx`

**Purpose**: Global authentication state management

**State**:

```typescript
interface AuthContextType {
  user: User | null; // Current logged-in user
  token: string | null; // JWT token
  login: (token: string) => Promise<void>;
  logout: () => void;
  loading: boolean; // Auth check in progress
}
```

**Features**:

- Auto-fetch user on mount (if token in localStorage)
- Auto-redirect to login if unauthenticated
- Token persistence in localStorage
- Axios header management

**Usage in Components**:

```tsx
import { useAuth } from "@/app/contexts/AuthContext";

function MyComponent() {
  const { user, logout, loading } = useAuth();

  if (loading) return <div>Loading...</div>;
  if (!user) return <div>Not logged in</div>;

  return <div>Welcome, {user.username}!</div>;
}
```

---

### 6.2 `withAuth` HOC

**File**: `src/app/hoc/withAuth.tsx`

**Purpose**: Higher-Order Component for protecting routes

**Usage**:

```tsx
import withAuth from "@/app/hoc/withAuth";

function ProtectedPage() {
  return <div>Only logged-in users see this</div>;
}

export default withAuth(ProtectedPage);
```

**Behavior**:

- Redirects to `/` if not authenticated
- Shows loading spinner while checking auth
- Passes user and token as props

---

## 7. API Client

### 7.1 API Client Overview

**File**: `src/lib/apiClient.ts`

**Purpose**: Centralized HTTP client for backend API

**Base Configuration**:

```typescript
const instance = axios.create({
  baseURL: "http://localhost:8080/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
});
```

**Axios Interceptors**:

- **Response Interceptor**: Unwraps `data` object, handles errors
- **Error Interceptor**: Converts error responses to `APIError` type

---

### 7.2 Available Methods

**Authentication**:

```typescript
apiClient.login(credentials: UserLoginRequest)
  → Promise<{ data: { user, token }, success }>

apiClient.register(details: UserRegisterRequest)
  → Promise<{ data: { user, token }, success }>

apiClient.getCurrentUser()
  → Promise<{ data: { user }, success }>
```

**Manga Operations**:

```typescript
apiClient.getManga(params: GetMangaParams)
  → Promise<{ data: { items: Manga[] }, meta, success }>
// params: { title?, author?, genre?, status?, limit?, offset? }
```

**Library Operations**:

```typescript
apiClient.addToLibrary(mangaId: string)
  → Promise<APIResponse>

apiClient.getLibrary()
  → Promise<{ data: { items: UserProgressWithManga[] }, meta, success }>

apiClient.updateMangaProgress(mangaId: string, data: ProgressUpdateRequest)
  → Promise<APIResponse>
// data: { current_chapter?, status?, rating? }
```

**Setting Headers** (for JWT):

```typescript
apiClient.setDefaultHeader("Authorization", `Bearer ${token}`);
```

---

### 7.3 Type Safety

All API calls use TypeScript types from `@mangahub/types`:

```typescript
import type {
  User,
  Manga,
  UserProgressWithManga,
  ProgressUpdateRequest,
  APIError,
} from "@mangahub/types";
```

**Benefits**:

- Compile-time type checking
- IDE autocomplete
- Runtime error prevention
- Consistent data structures

---

## 8. Type System

### 8.1 Shared Types (`@mangahub/types`)

**Package**: `packages/types/`

**Generated From**: OpenAPI specification (`packages/spec/openapi.yaml`)

**Key Types**:

```typescript
interface User {
  id: string;
  username: string;
  email: string;
  created_at: string;
  updated_at: string;
}

interface Manga {
  id: string;
  title: string;
  author: string;
  genres: string[];
  status: MangaStatus;
  total_chapters: number;
  description: string;
  cover_image_url: string;
}

type MangaStatus = "ongoing" | "completed" | "hiatus" | "cancelled";
type ReadingStatus =
  | "reading"
  | "completed"
  | "plan_to_read"
  | "on_hold"
  | "dropped";

interface UserProgressWithManga {
  manga: Manga;
  current_chapter: number;
  status: ReadingStatus;
  rating?: number;
  updated_at: string;
}
```

**Import in Components**:

```typescript
import type { Manga, User } from "@mangahub/types";
```

---

### 8.2 Local Types

**Component Props**:

```typescript
// Define props inline or as interface
interface MyComponentProps {
  title: string;
  onClose: () => void;
  manga?: Manga;
}

export default function MyComponent({
  title,
  onClose,
  manga,
}: MyComponentProps) {
  // ...
}
```

---

## 9. Styling with Tailwind CSS 4

### 9.1 Configuration

**File**: `tailwind.config.ts`

**Key Settings**:

- **Content**: `src/**/*.{ts,tsx}`
- **Theme**: Custom colors, fonts, etc.
- **Plugins**: None (vanilla Tailwind)

---

### 9.2 Common Patterns

**Card Component**:

```tsx
<div className="bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow p-4">
  {/* Card content */}
</div>
```

**Button**:

```tsx
<button className="bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded-lg transition-colors">
  Click Me
</button>
```

**Form Input**:

```tsx
<input
  type="text"
  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
  placeholder="Search manga..."
/>
```

**Grid Layout**:

```tsx
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
  {manga.map((m) => (
    <MangaCard key={m.id} manga={m} />
  ))}
</div>
```

---

## 10. State Management

### 10.1 useState for Component State

```tsx
function SearchPage() {
  const [manga, setManga] = useState<Manga[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchManga = async () => {
    setLoading(true);
    try {
      const response = await apiClient.getManga({});
      setManga(response.data.items);
    } catch (err) {
      setError("Failed to fetch manga");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchManga();
  }, []);

  return <div>{/* Render manga */}</div>;
}
```

---

### 10.2 Context API for Global State

**AuthContext** (already covered):

- User authentication state
- Available across all components

**Future Contexts** (potential):

- **ThemeContext**: Dark mode toggle
- **NotificationContext**: Toast notifications
- **WebSocketContext**: Real-time chat

---

### 10.3 Local Storage for Persistence

```tsx
// Save
localStorage.setItem("token", token);

// Read
const token = localStorage.getItem("token");

// Remove
localStorage.removeItem("token");
```

**Used For**:

- JWT token persistence
- User preferences (future)
- Reading history cache (future)

---

## 11. Development Workflow

### 11.1 Running the Dev Server

```bash
# From project root
make js-dev

# OR from apps/web directory
yarn dev
```

**Access**: http://localhost:3000

**Features**:

- Hot Module Replacement (HMR)
- Fast Refresh for React
- TypeScript type checking
- Auto-reload on file changes

---

### 11.2 Building for Production

```bash
# Build
yarn workspace @mangahub/web build

# Start production server
yarn workspace @mangahub/web start
```

**Output**: `.next/` directory with optimized bundles

---

### 11.3 Type Checking

```bash
# Check types without building
yarn workspace @mangahub/web typecheck
```

**Runs**: `tsc --noEmit` (type check only, no output)

---

### 11.4 Linting

```bash
# Lint code
yarn workspace @mangahub/web lint
```

**Uses**: ESLint with Next.js config

---

## 12. Environment Variables

### 12.1 API Base URL

**Current** (hardcoded):

```typescript
// src/lib/apiClient.ts
baseURL: "http://localhost:8080/api/v1";
```

**Production** (recommended):

```typescript
baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
```

Create `.env.local`:

```bash
NEXT_PUBLIC_API_URL=https://api.mangahub.com/api/v1
```

---

### 12.2 WebSocket URL

**Future**:

```typescript
const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:9093";
```

`.env.local`:

```bash
NEXT_PUBLIC_WS_URL=wss://chat.mangahub.com
```

---

## 13. Integration with Backend

### 13.1 HTTP API Integration

**Endpoints Used**:

- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `GET /users/me` - Get current user
- `GET /manga` - Search manga
- `GET /users/library` - Get user's library
- `POST /users/library` - Add manga to library
- `PUT /users/progress/:manga_id` - Update reading progress
- `POST /admin/notifications` - Trigger chapter notification

**Authentication**:

- JWT token in `Authorization: Bearer <token>` header
- Set via `apiClient.setDefaultHeader()`

---

### 13.2 WebSocket Integration (Ready)

**Scaffolded** but not fully implemented.

**How to Connect**:

```typescript
const ws = new WebSocket(
  "ws://localhost:9093/ws?username=" + user.username + "&room=general"
);

ws.onopen = () => console.log("Connected to chat");

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  // Handle chat message, system notification, etc.
};

ws.send(
  JSON.stringify({
    type: "chat",
    username: user.username,
    room: "general",
    content: "Hello!",
    timestamp: new Date().toISOString(),
  })
);
```

**Future Implementation**:

- Create `<ChatWidget />` component
- Create `useWebSocket()` custom hook
- Real-time chapter notifications via WebSocket
- Live chat with other users

---

### 13.3 UDP Notifications (Indirect)

**Flow**:

1. User clicks "Add New Chapter" in MangaModal
2. Frontend sends HTTP POST to `/admin/notifications`
3. Backend API server sends message to UDP server
4. UDP server broadcasts to all registered clients
5. API server (as UDP client) receives notification
6. API server bridges to WebSocket "general" room
7. WebSocket clients (browsers) receive notification

**This demonstrates inter-protocol integration!**

---

## 14. Performance Optimization

### 14.1 Next.js Optimizations

**Automatic**:

- Code splitting per route
- Image optimization with `next/image`
- Font optimization with `next/font`
- Bundle size optimization

**Manual**:

```tsx
// Lazy load components
const MangaModal = dynamic(() => import("@/component/mangaModal"), {
  ssr: false,
  loading: () => <p>Loading...</p>,
});
```

---

### 14.2 React Optimizations

**useMemo** for expensive computations:

```tsx
const filteredManga = useMemo(() => {
  return manga.filter((m) => m.status === "ongoing");
}, [manga]);
```

**useCallback** for event handlers:

```tsx
const handleSearch = useCallback((params: GetMangaParams) => {
  // Search logic
}, []);
```

**React.memo** for component memoization:

```tsx
export default React.memo(MangaCard);
```

---

### 14.3 Image Optimization

**Use Next.js Image Component**:

```tsx
import Image from "next/image";

<Image
  src={manga.cover_image_url}
  alt={manga.title}
  width={200}
  height={300}
  quality={75}
  loading="lazy"
/>;
```

**Benefits**:

- Automatic WebP conversion
- Lazy loading
- Responsive images
- Optimized delivery

---

## 15. Deployment

### 15.1 Vercel (Recommended)

**Steps**:

1. Push code to GitHub
2. Import project in Vercel
3. Set environment variables
4. Deploy

**Environment Variables** (Vercel dashboard):

```
NEXT_PUBLIC_API_URL=https://api.mangahub.com/api/v1
NEXT_PUBLIC_WS_URL=wss://chat.mangahub.com
```

---

### 15.2 Docker (Self-Hosted)

**Dockerfile** (create in `apps/web/`):

```dockerfile
FROM node:18-alpine AS builder

WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install
COPY . .
RUN yarn build

FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./
RUN yarn install --production
EXPOSE 3000
CMD ["yarn", "start"]
```

**Build and run**:

```bash
docker build -t mangahub-web .
docker run -p 3000:3000 -e NEXT_PUBLIC_API_URL=http://api:8080/api/v1 mangahub-web
```

---

### 15.3 Static Export (Future)

**For fully static site**:

```javascript
// next.config.js
module.exports = {
  output: "export",
};
```

**Build**:

```bash
yarn build
```

**Output**: `out/` directory with static HTML/CSS/JS

**Serve**:

```bash
npx serve out
```

---

## 16. Troubleshooting

### Common Issues

#### CORS Errors

**Symptom**:

```
Access to fetch at 'http://localhost:8080/api/v1/manga' from origin 'http://localhost:3000' has been blocked by CORS policy
```

**Solution**:
Backend must include CORS middleware allowing `http://localhost:3000`

---

#### API Not Found (404)

**Check**:

1. Is API server running? `curl http://localhost:8080/health`
2. Is baseURL correct in apiClient.ts?
3. Is endpoint path correct?

---

#### Token Expiry

**Symptom**: User randomly logged out

**Solution**:

- Implement token refresh mechanism
- Or increase JWT expiry time in backend config

---

#### Image Not Loading

**Symptom**: Manga cover images don't show

**Solution**:

1. Check `next.config.js` `remotePatterns`
2. Add image hostname to allowed list
3. Verify image URLs are valid

---

## 17. Future Improvements

### Planned Enhancements

- **React Query**: Better server state management
- **Zustand/Redux**: Client state management (if needed)
- **React Hook Form**: Better form handling
- **Zod**: Runtime validation
- **SWR**: Data fetching with cache
- **Framer Motion**: Animations
- **PWA**: Progressive Web App features
- **Service Worker**: Offline support
- **IndexedDB**: Local caching

---

## 18. References

### Internal Documentation

- [Architecture](./architecture.md) - System architecture
- [API Documentation](./api-documentation.md) - HTTP API reference
- [WebSocket Documentation](./websocket-documentation.md) - Chat protocol
- [MONOREPO.md](../MONOREPO.md) - Monorepo structure

### External Resources

- [Next.js Documentation](https://nextjs.org/docs)
- [React Documentation](https://react.dev)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
- [Axios Documentation](https://axios-http.com/docs/intro)

---

**Last Updated**: 2025-12-26
**Version**: 1.0.0
**Framework**: Next.js 16 + React 19
**Status**: ✅ Fully Functional
