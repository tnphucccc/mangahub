/**
 * MangaHub API Client - Fetch Implementation
 *
 * A type-safe API client using native Fetch API
 * Copy this file to your frontend project and customize as needed
 */

import type {
  UserRegisterRequest,
  UserLoginRequest,
  LibraryAddRequest,
  ProgressUpdateRequest,
  AuthResponse,
  UserProfileResponse,
  MangaListResponse,
  MangaDetailResponse,
  LibraryResponse,
  ProgressResponse,
  APIResponse,
} from '@mangahub/types';

export interface MangaHubClientConfig {
  baseUrl?: string;
  onTokenExpired?: () => void;
  onError?: (error: APIError) => void;
}

export interface APIError {
  code: string;
  message: string;
  status?: number;
}

export class MangaHubClient {
  private baseUrl: string;
  private token: string | null = null;
  private config: MangaHubClientConfig;

  constructor(config: MangaHubClientConfig = {}) {
    this.baseUrl = config.baseUrl || 'http://localhost:8080/api/v1';
    this.config = config;
  }

  // ==========================================
  // Authentication Methods
  // ==========================================

  /**
   * Register a new user
   */
  async register(data: UserRegisterRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });

    // Auto-save token
    if (response.success && response.data) {
      this.setToken(response.data.token);
    }

    return response;
  }

  /**
   * Login user
   */
  async login(data: UserLoginRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });

    // Auto-save token
    if (response.success && response.data) {
      this.setToken(response.data.token);
    }

    return response;
  }

  /**
   * Get current user profile
   */
  async getProfile(): Promise<UserProfileResponse> {
    return this.request<UserProfileResponse>('/users/me', {
      method: 'GET',
      authenticated: true,
    });
  }

  /**
   * Logout (client-side only)
   */
  logout(): void {
    this.token = null;
    localStorage.removeItem('mangahub_token');
  }

  // ==========================================
  // Manga Methods
  // ==========================================

  /**
   * Search manga with filters
   */
  async searchManga(params: {
    title?: string;
    author?: string;
    genre?: string;
    status?: string;
    limit?: number;
    offset?: number;
  }): Promise<MangaListResponse> {
    const queryString = new URLSearchParams(
      Object.entries(params)
        .filter(([_, value]) => value !== undefined)
        .map(([key, value]) => [key, String(value)])
    ).toString();

    return this.request<MangaListResponse>(`/manga?${queryString}`, {
      method: 'GET',
    });
  }

  /**
   * Get all manga with pagination
   */
  async getAllManga(limit = 20, offset = 0): Promise<MangaListResponse> {
    return this.request<MangaListResponse>(
      `/manga/all?limit=${limit}&offset=${offset}`,
      { method: 'GET' }
    );
  }

  /**
   * Get manga by ID
   */
  async getManga(id: string): Promise<MangaDetailResponse> {
    return this.request<MangaDetailResponse>(`/manga/${id}`, {
      method: 'GET',
    });
  }

  // ==========================================
  // Library Methods
  // ==========================================

  /**
   * Get user's manga library
   */
  async getLibrary(): Promise<LibraryResponse> {
    return this.request<LibraryResponse>('/users/library', {
      method: 'GET',
      authenticated: true,
    });
  }

  /**
   * Add manga to library
   */
  async addToLibrary(data: LibraryAddRequest): Promise<APIResponse> {
    return this.request<APIResponse>('/users/library', {
      method: 'POST',
      body: JSON.stringify(data),
      authenticated: true,
    });
  }

  // ==========================================
  // Progress Methods
  // ==========================================

  /**
   * Get reading progress for a specific manga
   */
  async getProgress(mangaId: string): Promise<ProgressResponse> {
    return this.request<ProgressResponse>(`/users/progress/${mangaId}`, {
      method: 'GET',
      authenticated: true,
    });
  }

  /**
   * Update reading progress
   */
  async updateProgress(
    mangaId: string,
    data: ProgressUpdateRequest
  ): Promise<APIResponse> {
    return this.request<APIResponse>(`/users/progress/${mangaId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
      authenticated: true,
    });
  }

  // ==========================================
  // Token Management
  // ==========================================

  /**
   * Set authentication token
   */
  setToken(token: string): void {
    this.token = token;
    localStorage.setItem('mangahub_token', token);
  }

  /**
   * Get current token
   */
  getToken(): string | null {
    if (!this.token) {
      this.token = localStorage.getItem('mangahub_token');
    }
    return this.token;
  }

  /**
   * Check if user is authenticated
   */
  isAuthenticated(): boolean {
    return this.getToken() !== null;
  }

  // ==========================================
  // Core Request Method
  // ==========================================

  private async request<T>(
    endpoint: string,
    options: {
      method: string;
      body?: string;
      authenticated?: boolean;
      headers?: Record<string, string>;
    }
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    // Add authentication token if required
    if (options.authenticated) {
      const token = this.getToken();
      if (!token) {
        throw new Error('Authentication required but no token found');
      }
      headers['Authorization'] = `Bearer ${token}`;
    }

    try {
      const response = await fetch(url, {
        method: options.method,
        headers,
        body: options.body,
      });

      // Handle HTTP errors
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({
          success: false,
          error: {
            code: 'UNKNOWN_ERROR',
            message: `HTTP ${response.status}: ${response.statusText}`,
          },
        }));

        // Handle token expiration
        if (response.status === 401 && this.config.onTokenExpired) {
          this.logout();
          this.config.onTokenExpired();
        }

        const error: APIError = {
          code: errorData.error?.code || 'HTTP_ERROR',
          message: errorData.error?.message || response.statusText,
          status: response.status,
        };

        if (this.config.onError) {
          this.config.onError(error);
        }

        throw error;
      }

      return response.json();
    } catch (error) {
      if (error instanceof Error && !(error as any).code) {
        // Network error
        const apiError: APIError = {
          code: 'NETWORK_ERROR',
          message: error.message,
        };

        if (this.config.onError) {
          this.config.onError(apiError);
        }

        throw apiError;
      }

      throw error;
    }
  }
}

// ==========================================
// Export Singleton Instance
// ==========================================

export const mangaHubApi = new MangaHubClient();

// ==========================================
// Usage Example
// ==========================================

/*
import { mangaHubApi } from './api-client';

// Login
const { data } = await mangaHubApi.login({
  username: 'testuser',
  password: 'password123'
});

console.log('Logged in as:', data.user.username);
console.log('Token:', data.token);

// Search manga
const searchResults = await mangaHubApi.searchManga({
  title: 'naruto',
  genre: 'action',
  limit: 10
});

console.log('Found:', searchResults.meta.total, 'manga');

// Get library
const library = await mangaHubApi.getLibrary();
console.log('Your library:', library.data.items);

// Update progress
await mangaHubApi.updateProgress('manga-001', {
  current_chapter: 75,
  status: 'reading',
  rating: 9
});
*/
