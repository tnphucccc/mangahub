/**
 * MangaHub React Hooks
 *
 * Custom React hooks for integrating with MangaHub API
 * Provides authentication, data fetching, and real-time updates
 *
 * Usage:
 *   Copy this file to your React project
 *   Install: npm install @tanstack/react-query
 */

import { useState, useEffect, useCallback } from 'react';
import { mangaHubApi } from '@mangahub/api';
import type {
  User,
  Manga,
  UserProgressWithManga,
  UserLoginRequest,
  UserRegisterRequest,
  LibraryAddRequest,
  ProgressUpdateRequest,
} from '@mangahub/types';

// ==========================================
// Authentication Hook
// ==========================================

interface UseAuthReturn {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  login: (credentials: UserLoginRequest) => Promise<void>;
  register: (data: UserRegisterRequest) => Promise<void>;
  logout: () => void;
}

export function useAuth(): UseAuthReturn {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Check for existing session on mount
  useEffect(() => {
    const checkAuth = async () => {
      const existingToken = mangaHubApi.getToken();
      if (existingToken) {
        try {
          const response = await mangaHubApi.getProfile();
          if (response.success && response.data) {
            setUser(response.data.user);
            setToken(existingToken);
          }
        } catch (err) {
          // Token expired or invalid
          mangaHubApi.logout();
        }
      }
      setIsLoading(false);
    };

    checkAuth();
  }, []);

  const login = useCallback(async (credentials: UserLoginRequest) => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await mangaHubApi.login(credentials);
      if (response.success && response.data) {
        setUser(response.data.user);
        setToken(response.data.token);
      }
    } catch (err: any) {
      setError(err.message || 'Login failed');
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const register = useCallback(async (data: UserRegisterRequest) => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await mangaHubApi.register(data);
      if (response.success && response.data) {
        setUser(response.data.user);
        setToken(response.data.token);
      }
    } catch (err: any) {
      setError(err.message || 'Registration failed');
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const logout = useCallback(() => {
    mangaHubApi.logout();
    setUser(null);
    setToken(null);
  }, []);

  return {
    user,
    token,
    isAuthenticated: !!user,
    isLoading,
    error,
    login,
    register,
    logout,
  };
}

// ==========================================
// Manga Search Hook
// ==========================================

interface UseMangaSearchParams {
  title?: string;
  author?: string;
  genre?: string;
  status?: string;
  limit?: number;
  offset?: number;
}

interface UseMangaSearchReturn {
  manga: Manga[];
  total: number;
  isLoading: boolean;
  error: string | null;
  hasMore: boolean;
  search: (params: UseMangaSearchParams) => Promise<void>;
}

export function useMangaSearch(): UseMangaSearchReturn {
  const [manga, setManga] = useState<Manga[]>([]);
  const [total, setTotal] = useState(0);
  const [hasMore, setHasMore] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const search = useCallback(async (params: UseMangaSearchParams) => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await mangaHubApi.searchManga(params);
      if (response.success && response.data) {
        setManga(response.data.items);
        setTotal(response.meta?.total || 0);
        setHasMore(response.meta?.has_more || false);
      }
    } catch (err: any) {
      setError(err.message || 'Search failed');
    } finally {
      setIsLoading(false);
    }
  }, []);

  return {
    manga,
    total,
    isLoading,
    error,
    hasMore,
    search,
  };
}

// ==========================================
// Manga Detail Hook
// ==========================================

interface UseMangaReturn {
  manga: Manga | null;
  isLoading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

export function useManga(mangaId: string | null): UseMangaReturn {
  const [manga, setManga] = useState<Manga | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchManga = useCallback(async () => {
    if (!mangaId) return;

    setIsLoading(true);
    setError(null);

    try {
      const response = await mangaHubApi.getManga(mangaId);
      if (response.success && response.data) {
        setManga(response.data.manga);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to fetch manga');
    } finally {
      setIsLoading(false);
    }
  }, [mangaId]);

  useEffect(() => {
    fetchManga();
  }, [fetchManga]);

  return {
    manga,
    isLoading,
    error,
    refetch: fetchManga,
  };
}

// ==========================================
// Library Hook
// ==========================================

interface UseLibraryReturn {
  library: UserProgressWithManga[];
  isLoading: boolean;
  error: string | null;
  addToLibrary: (data: LibraryAddRequest) => Promise<void>;
  updateProgress: (mangaId: string, data: ProgressUpdateRequest) => Promise<void>;
  refetch: () => Promise<void>;
}

export function useLibrary(): UseLibraryReturn {
  const [library, setLibrary] = useState<UserProgressWithManga[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchLibrary = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await mangaHubApi.getLibrary();
      if (response.success && response.data) {
        setLibrary(response.data.items);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to fetch library');
    } finally {
      setIsLoading(false);
    }
  }, []);

  const addToLibrary = useCallback(
    async (data: LibraryAddRequest) => {
      try {
        await mangaHubApi.addToLibrary(data);
        await fetchLibrary(); // Refresh library
      } catch (err: any) {
        throw new Error(err.message || 'Failed to add to library');
      }
    },
    [fetchLibrary]
  );

  const updateProgress = useCallback(
    async (mangaId: string, data: ProgressUpdateRequest) => {
      try {
        await mangaHubApi.updateProgress(mangaId, data);
        await fetchLibrary(); // Refresh library
      } catch (err: any) {
        throw new Error(err.message || 'Failed to update progress');
      }
    },
    [fetchLibrary]
  );

  useEffect(() => {
    fetchLibrary();
  }, [fetchLibrary]);

  return {
    library,
    isLoading,
    error,
    addToLibrary,
    updateProgress,
    refetch: fetchLibrary,
  };
}

// ==========================================
// Real-Time Updates Hook (TCP)
// ==========================================

interface TCPProgressUpdate {
  userId: string;
  username: string;
  mangaId: string;
  mangaTitle: string;
  currentChapter: number;
  status: string;
  timestamp: string;
}

interface UseTCPConnectionReturn {
  isConnected: boolean;
  latestUpdate: TCPProgressUpdate | null;
  connect: () => void;
  disconnect: () => void;
}

export function useTCPConnection(token: string | null): UseTCPConnectionReturn {
  const [isConnected, setIsConnected] = useState(false);
  const [latestUpdate, setLatestUpdate] = useState<TCPProgressUpdate | null>(null);
  const [socket, setSocket] = useState<WebSocket | null>(null);

  const connect = useCallback(() => {
    if (!token || isConnected) return;

    // Note: This would need a WebSocket wrapper for TCP
    // For now, this is a placeholder showing the pattern
    console.log('TCP connection would be established here');
    setIsConnected(true);
  }, [token, isConnected]);

  const disconnect = useCallback(() => {
    if (socket) {
      socket.close();
      setSocket(null);
    }
    setIsConnected(false);
  }, [socket]);

  useEffect(() => {
    return () => {
      disconnect();
    };
  }, [disconnect]);

  return {
    isConnected,
    latestUpdate,
    connect,
    disconnect,
  };
}

// ==========================================
// Usage Examples
// ==========================================

/*
// Example 1: Authentication
function LoginPage() {
  const { login, isLoading, error } = useAuth();

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      await login({
        username: 'testuser',
        password: 'password123'
      });
      // Redirect to home
    } catch (err) {
      // Error is already set in hook
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {error && <div className="error">{error}</div>}
      <input type="text" name="username" />
      <input type="password" name="password" />
      <button disabled={isLoading}>
        {isLoading ? 'Logging in...' : 'Login'}
      </button>
    </form>
  );
}

// Example 2: Search Manga
function SearchPage() {
  const { manga, total, isLoading, search } = useMangaSearch();
  const [query, setQuery] = useState('');

  const handleSearch = () => {
    search({ title: query, limit: 20 });
  };

  return (
    <div>
      <input
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
      />
      <button onClick={handleSearch} disabled={isLoading}>
        Search
      </button>

      {isLoading && <p>Loading...</p>}
      {manga.length > 0 && (
        <div>
          <p>Found {total} manga</p>
          {manga.map(m => (
            <MangaCard key={m.id} manga={m} />
          ))}
        </div>
      )}
    </div>
  );
}

// Example 3: User Library
function LibraryPage() {
  const { library, isLoading, updateProgress } = useLibrary();

  const handleUpdateProgress = async (mangaId: string, chapter: number) => {
    try {
      await updateProgress(mangaId, {
        current_chapter: chapter,
        status: 'reading'
      });
    } catch (err) {
      console.error(err);
    }
  };

  if (isLoading) return <p>Loading library...</p>;

  return (
    <div>
      <h1>My Library</h1>
      {library.map(item => (
        <div key={item.manga_id}>
          <h3>{item.manga.title}</h3>
          <p>Chapter: {item.current_chapter}</p>
          <button onClick={() => handleUpdateProgress(item.manga_id, item.current_chapter + 1)}>
            Mark Next Chapter Read
          </button>
        </div>
      ))}
    </div>
  );
}

// Example 4: Protected Route
function ProtectedRoute({ children }: { children: ReactNode }) {
  const { isAuthenticated, isLoading } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate('/login');
    }
  }, [isAuthenticated, isLoading, navigate]);

  if (isLoading) return <p>Loading...</p>;
  if (!isAuthenticated) return null;

  return <>{children}</>;
}
*/
