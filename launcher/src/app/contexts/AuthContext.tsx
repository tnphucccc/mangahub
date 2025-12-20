'use client'

import {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { User } from '../../../../api/typescript'
import { apiClient } from '../../lib/apiClient'

interface AuthContextType {
  user: User | null
  token: string | null
  login: (token: string) => Promise<void>
  logout: () => void
  loading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const router = useRouter()
  const pathname = usePathname()

  const fetchUser = async (authToken: string) => {
    try {
      apiClient.setDefaultHeader('Authorization', `Bearer ${authToken}`)
      // The `getCurrentUser` response is wrapped in a `user` object.
      const response = await apiClient.getCurrentUser()
      setUser(response.user)
    } catch (error) {
      console.error('Failed to fetch user', error)
      logout()
    } finally {
      setLoading(false)
      router.push('/main')
    }
  }

  useEffect(() => {
    const storedToken = localStorage.getItem('token')
    if (storedToken) {
      setToken(storedToken)
      fetchUser(storedToken)
    } else {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    if (!loading && !user && pathname !== '/') {
      // Redirect to login if not authenticated and not already on login page
      router.push('/')
    }
  }, [user, loading, pathname, router])

  const login = async (newToken: string) => {
    setLoading(true)
    localStorage.setItem('token', newToken)
    setToken(newToken)
    await fetchUser(newToken)
  }

  const logout = () => {
    setLoading(true)
    localStorage.removeItem('token')
    setToken(null)
    setUser(null)
    apiClient.setDefaultHeader('Authorization', '') // Clear the header
    router.push('/')
    setLoading(false)
  }

  return (
    <AuthContext.Provider value={{ user, token, login, logout, loading }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
