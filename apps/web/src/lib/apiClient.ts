import axios from 'axios'
import type {
  UserLoginRequest,
  UserRegisterRequest,
  User,
  Manga,
  APIError,
  MangaStatus,
  Meta,
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
    sucess: boolean
  }>
  register: (
    details: UserRegisterRequest
  ) => Promise<{ data: { user: User; token: string }; sucess: boolean }>
  getCurrentUser: () => Promise<{ data: { user: User }; sucess: boolean }>
  getManga: (
    params: GetMangaParams
  ) => Promise<{ data: { items: Manga[] }; meta: Meta; sucess: boolean }>
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
    // The endpoint for searching is /manga, not /manga/all
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
}
