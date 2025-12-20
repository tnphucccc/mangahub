"use client"

import withAuth from '../hoc/withAuth'
import { useAuth } from '../contexts/AuthContext'

function Home() {
  const { user, logout } = useAuth()

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gray-100">
      <div className="p-8 bg-white rounded-lg shadow-md text-center">
        <h1 className="text-2xl font-bold text-gray-900">
          Welcome to MangaHub!
        </h1>
        <p className="mt-2 text-gray-600">
          {user ? `You are logged in as ${user.email}` : 'Loading...'}
        </p>
        <button
          onClick={logout}
          className="mt-6 px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
        >
          Logout
        </button>
      </div>
    </div>
  )
}

export default withAuth(Home)
