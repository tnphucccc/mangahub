'use client'

import { useState } from 'react'
import { useAuth } from './contexts/AuthContext'
import { apiClient } from '../lib/apiClient'
import Login from '../component/login'
import Register from '../component/register'

export default function AuthPage() {
  // State for both forms
  const [userName, setUserName] = useState('')
  const [passWord, setPassword] = useState('')
  const [email, setEmail] = useState('')
  const [error, setError] = useState('')

  // State to toggle between Login and Register views
  const [isRegistering, setIsRegistering] = useState(false)

  // State to show success message after registration
  const [registrationSuccess, setRegistrationSuccess] = useState(false)

  const { login } = useAuth()

  const handleLoginSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (!userName || !passWord) {
      setError('Username and password are required')
      return
    }

    try {
      const response = await apiClient.login({
        username: userName,
        password: passWord,
      })
      console.log('Login response:', response)
      if (response && response.data) {
        await login(response.data.token)
      } else {
        setError('Login failed: No token received.')
      }
    } catch (err: any) {
      console.error('Login error:', err)
      setError(err.message || 'Failed to login. Please check your credentials.')
    }
  }

  const handleRegisterSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (!userName || !passWord || !email) {
      setError('Username, email, and password are required')
      return
    }

    try {
      await apiClient.register({
        username: userName,
        email: email,
        password: passWord,
      })

      // On successful registration
      setRegistrationSuccess(true)
      setIsRegistering(false) // Switch back to login view
      // Clear form fields
      setUserName('')
      setPassword('')
      setEmail('')
    } catch (error: any) {
      console.error('Registration error:', error)
      setError(error.message || 'Failed to register. Please try again.')
    }
  }

  const clearForm = () => {
    setUserName('')
    setPassword('')
    setEmail('')
    setError('')
    setRegistrationSuccess(false)
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-100">
      <div className="w-full max-w-md">
        {isRegistering ? (
          <>
            <Register
              userName={userName}
              passWord={passWord}
              email={email}
              setUserName={setUserName}
              setPassword={setPassword}
              setEmail={setEmail}
              error={error}
              handleSubmit={handleRegisterSubmit}
            />
            <p className="mt-4 text-center text-sm text-gray-600">
              Already have an account?{' '}
              <button
                onClick={() => {
                  setIsRegistering(false)
                  clearForm()
                }}
                className="font-medium text-indigo-600 hover:text-indigo-500"
              >
                Sign in
              </button>
            </p>
          </>
        ) : (
          <>
            {registrationSuccess && (
              <p className="mb-4 text-center text-sm text-green-600">
                Registration successful! Please sign in.
              </p>
            )}
            <Login
              userName={userName}
              passWord={passWord}
              setUserName={setUserName}
              setPassword={setPassword}
              error={error}
              handleSubmit={handleLoginSubmit}
            />
            <p className="mt-4 text-center text-sm text-gray-600">
              Don't have an account?{' '}
              <button
                onClick={() => {
                  setIsRegistering(true)
                  clearForm()
                }}
                className="font-medium text-indigo-600 hover:text-indigo-500"
              >
                Register
              </button>
            </p>
          </>
        )}
      </div>
    </div>
  )
}
