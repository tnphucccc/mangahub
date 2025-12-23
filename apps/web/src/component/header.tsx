'use client'

import { useAuth } from '../app/contexts/AuthContext'
import Link from 'next/link'

const Header = () => {
  const { user, logout } = useAuth()
  return (
    <header className="w-full bg-white shadow-md">
      <div className="container mx-auto px-6 py-4 flex justify-between items-center">
        <Link href="/main">
          <h1 className="text-xl font-bold text-gray-900 cursor-pointer">
            MangaHub
          </h1>
        </Link>
        <nav className="flex items-center">
          <Link href="/main">
            <span className="text-gray-600 hover:text-gray-900 mr-4">Home</span>
          </Link>
          <Link href="/library">
            <span className="text-gray-600 hover:text-gray-900 mr-4">
              My Library
            </span>
          </Link>
        </nav>
        <div className="flex items-center">
          <p className="text-gray-600 mr-4">
            Welcome, {user?.username || 'User'}
          </p>
          <button
            onClick={logout}
            className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700"
          >
            Logout
          </button>
        </div>
      </div>
    </header>
  )
}

export default Header
