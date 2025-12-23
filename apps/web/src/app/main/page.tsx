'use client'

import { useState, useEffect, useCallback } from 'react'
import withAuth from '../hoc/withAuth'
import { useAuth } from '../contexts/AuthContext'
import { apiClient, GetMangaParams } from '../../lib/apiClient'
import type { Manga, MangaStatus } from '../../../../../packages/types/src'
import MangaCard from '../../component/mangaCard'

// A simple debounce hook
function useDebounce(value: string, delay: number) {
  const [debouncedValue, setDebouncedValue] = useState(value)
  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value)
    }, delay)
    return () => {
      clearTimeout(handler)
    }
  }, [value, delay])
  return debouncedValue
}

const SearchAndFilter = ({
  onSearch,
}: {
  onSearch: (filters: {
    searchTerm: string
    status: MangaStatus | ''
    genre: string
  }) => void
}) => {
  const [searchTerm, setSearchTerm] = useState('')
  const [status, setStatus] = useState<MangaStatus | ''>('')
  const [genre, setGenre] = useState('')

  const debouncedSearchTerm = useDebounce(searchTerm, 500)
  const debouncedGenre = useDebounce(genre, 500)

  useEffect(() => {
    onSearch({ searchTerm: debouncedSearchTerm, status, genre: debouncedGenre })
  }, [debouncedSearchTerm, status, debouncedGenre, onSearch])

  return (
    <div className="bg-white p-4 rounded-lg shadow-md mb-8">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <input
          type="text"
          placeholder="Search by title..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="p-2 border rounded-md text-black"
        />
        <select
          value={status}
          onChange={(e) => setStatus(e.target.value as MangaStatus | '')}
          className="p-2 border rounded-md text-black"
        >
          <option value="">All Statuses</option>
          <option value="ongoing">Ongoing</option>
          <option value="completed">Completed</option>
          <option value="hiatus">Hiatus</option>
          <option value="cancelled">Cancelled</option>
        </select>
        <input
          type="text"
          placeholder="Filter by genre..."
          value={genre}
          onChange={(e) => setGenre(e.target.value)}
          className="p-2 border rounded-md text-black"
        />
      </div>
    </div>
  )
}

const Header = () => {
  const { user, logout } = useAuth()
  return (
    <header className="w-full bg-white shadow-md">
      <div className="container mx-auto px-6 py-4 flex justify-between items-center">
        <h1 className="text-xl font-bold text-gray-900">MangaHub</h1>
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

function HomePage() {
  const [mangaList, setMangaList] = useState<Manga[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // State for pagination and filtering
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [filters, setFilters] = useState({
    searchTerm: '',
    status: '' as MangaStatus | '',
    genre: '',
  })
  const MANGA_PER_PAGE = 10

  const fetchManga = useCallback(async () => {
    try {
      setLoading(true)
      const params: GetMangaParams = {
        title: filters.searchTerm,
        status: filters.status || undefined,
        genre: filters.genre || undefined,
        limit: MANGA_PER_PAGE,
        offset: (currentPage - 1) * MANGA_PER_PAGE,
      }
      const response = await apiClient.getManga(params)
      console.log('Fetched manga:', response)
      setMangaList(response.data.items || [])
      setTotalPages(Math.ceil((response.meta?.total || 0) / MANGA_PER_PAGE))
    } catch (err: any) {
      setError(err.message || 'Failed to fetch manga.')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [currentPage, filters])

  useEffect(() => {
    fetchManga()
  }, [fetchManga])

  const handleSearch = useCallback((newFilters: typeof filters) => {
    setCurrentPage(1) // Reset to first page on new search
    setFilters(newFilters)
  }, [])

  return (
    <div className="min-h-screen bg-gray-100">
      <Header />
      <main className="container mx-auto px-6 py-8">
        <SearchAndFilter onSearch={handleSearch} />
        {loading && <p className="text-center">Loading manga...</p>}
        {error && <p className="text-center text-red-500">{error}</p>}
        {!loading && !error && mangaList.length > 0 ? (
          <>
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-6">
              {mangaList.map((manga) => (
                <MangaCard key={manga.id} manga={manga} />
              ))}
            </div>
            <div className="mt-8 flex justify-center items-center space-x-4">
              <button
                onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                disabled={currentPage === 1}
                className="px-4 py-2 bg-indigo-600 rounded-md disabled:bg-gray-400 text-black"
              >
                Previous
              </button>
              <span className="text-black">
                Page {currentPage} of {totalPages || 1}
              </span>
              <button
                onClick={() =>
                  setCurrentPage((p) => Math.min(totalPages, p + 1))
                }
                disabled={currentPage === totalPages || totalPages === 0}
                className="px-4 py-2 bg-indigo-600 text-black rounded-md disabled:bg-gray-400"
              >
                Next
              </button>
            </div>
          </>
        ) : (
          !loading && <p className="text-center">No manga found.</p>
        )}
      </main>
    </div>
  )
}

export default withAuth(HomePage)
