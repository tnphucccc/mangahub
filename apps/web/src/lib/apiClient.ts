import axios from 'axios'
import type {
  UserLoginRequest,
  UserRegisterRequest,
  User,
  Manga,
  APIError,
  MangaStatus,
  Meta,
  LibraryAddRequest,
  APIResponse,
  UserProgressWithManga,
  ProgressUpdateRequest,
} from '../../../../packages/types/src'

// Create an axios instance with a base URL.
// The Go API server runs on port 8080.
const instance = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add a response interceptor to handle successful responses and errors globally.
instance.interceptors.response.use(
  // For successful responses (2xx), just return the data part of the response.
  (response) => {
    // The backend wraps successful responses in `data` and `meta` objects.
    // We return the whole response.data here to include both.
    return response.data
  },
  // For error responses, reject with a structured error.
  (error) => {
    if (error.response && error.response.data) {
      // The backend wraps error responses in an `error` object.
      const apiError: APIError = error.response.data.error || {
        code: 'UNKNOWN_ERROR',
        message: 'An unexpected error occurred.',
      }
      return Promise.reject(apiError)
    }
    return Promise.reject(error)
  }
)

export interface GetMangaParams {
  title?: string
  author?: string
  genre?: string
  status?: MangaStatus
  limit?: number
  offset?: number
}

// Define the shape of our API client.
interface ApiClient {
  setDefaultHeader: (name: string, value: string) => void
  login: (credentials: UserLoginRequest) => Promise<{
    data: {
      user: User
      token: string
    }
    success: boolean
  }>
  register: (
    details: UserRegisterRequest
  ) => Promise<{ data: { user: User; token: string }; success: boolean }>
  getCurrentUser: () => Promise<{ data: { user: User }; success: boolean }>
  getManga: (
    params: GetMangaParams
  ) => Promise<{ data: { items: Manga[] }; meta: Meta; success: boolean }>
  addToLibrary: (mangaId: string) => Promise<APIResponse>
  getLibrary: () => Promise<{
    data: { items: UserProgressWithManga[] }
    meta: Meta
    success: boolean
  }>
  updateMangaProgress: (
    mangaId: string,
    data: ProgressUpdateRequest
  ) => Promise<APIResponse>
  sendNotification: (
    manga_id: string,
    manga_title: string,
    chapter_number: number,
    chapter_title: string,
    release_date: string,
    message: string
  ) => Promise<APIResponse>
}

// Implement the ApiClient.
export const apiClient: ApiClient = {
  setDefaultHeader: (name, value) => {
    if (value) {
      instance.defaults.headers.common[name] = value
    } else {
      delete instance.defaults.headers.common[name]
    }
  },

  login: async (credentials) => {
    const request: UserLoginRequest = {
      username: credentials.username,
      password: credentials.password,
    }
    return instance.post('/auth/login', request)
  },

  register: async (details) => {
    const request: UserRegisterRequest = {
      username: details.username,
      email: details.email,
      password: details.password,
    }
    return instance.post('/auth/register', request)
  },

  getCurrentUser: async () => {
    return instance.get('/users/me')
  },

  getManga: async (params) => {
    const endpoint = '/manga'

    // Create a query string from the params object
    const queryParams = new URLSearchParams()

    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== '') {
        queryParams.append(key, String(value))
      }
    })

    const queryString = queryParams.toString()

    return instance.get(`${endpoint}${queryString ? `?${queryString}` : ''}`)
  },

  addToLibrary: async (mangaId: string) => {
    const request: LibraryAddRequest = {
      manga_id: mangaId,
      status: 'plan_to_read', // Default status when adding to library
      current_chapter: 0,
    }
    return instance.post('/users/library', request)
  },

  getLibrary: async () => {
    return instance.get('/users/library')
  },

  updateMangaProgress: async (mangaId, data) => {
    return instance.put(`/users/progress/${mangaId}`, data)
  },

  sendNotification: async (
    manga_id,
    manga_title,
    chapter_number,
    chapter_title,
    release_date,
    message
  ) => {
    const notification = {
      manga_id,
      manga_title,
      chapter_number,
      chapter_title,
      release_date,
      message,
    }
    return instance.post(`/admin/notifications`, notification)
  },
}
