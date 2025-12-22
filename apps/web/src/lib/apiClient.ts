import axios from 'axios'
import type {
  UserLoginRequest,
  UserRegisterRequest,
  User,
  APIError,
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
    // The backend wraps successful responses in a `data` object.
    return response.data.data || response.data
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

// Define the shape of our API client.
// Note: The axios interceptor unwraps responses, so we return the data directly
interface ApiClient {
  setDefaultHeader: (name: string, value: string) => void
  login: (
    credentials: UserLoginRequest
  ) => Promise<{ user: User; token: string }>
  register: (
    details: UserRegisterRequest
  ) => Promise<{ user: User; token: string }>
  getCurrentUser: () => Promise<User>
}

// Implement the ApiClient.
export const apiClient: ApiClient = {
  /**
   * Sets or clears a default header for all subsequent requests.
   * @param name - The name of the header (e.g., 'Authorization').
   * @param value - The value for the header. An empty string will clear it.
   */
  setDefaultHeader: (name, value) => {
    if (value) {
      instance.defaults.headers.common[name] = value
    } else {
      delete instance.defaults.headers.common[name]
    }
  },

  /**
   * Performs user login.
   * @param credentials - The user's login credentials (username, password).
   * @returns The authentication response containing the user and token.
   */
  login: async (credentials) => {
    // Note: The Go API expects `username`, but our form uses `email`.
    // The backend logic should be checked to confirm which one it uses.
    // For now, we assume the login request can handle the 'email' field as a username.
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

  /**
   * Fetches the profile of the currently authenticated user.
   * @returns The user's profile information.
   */
  getCurrentUser: async () => {
    return instance.get('/users/me')
  },
}
