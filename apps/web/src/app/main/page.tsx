'use client'

import { useState, useEffect, useCallback } from 'react'
import withAuth from '../hoc/withAuth'
import { useAuth } from '../contexts/AuthContext'
import { apiClient, GetMangaParams } from '../../lib/apiClient'
import type { Manga, MangaStatus } from '../../../../../packages/types/src'
import MangaCard from '../../component/mangaCard'
import Header from '../../component/header'
import MangaModal from '../../component/mangaModal'
import { ToastContainer, toast } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'
import { SearchAndFilter } from '../../component/searchAndFilter'
import { ClipLoader } from 'react-spinners'

function HomePage() {
  const [mangaList, setMangaList] = useState<Manga[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // State for modal
  const [selectedManga, setSelectedManga] = useState<Manga | null>(null)

  // State for pagination and filtering
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [filters, setFilters] = useState({
    searchTerm: '',
    status: '' as MangaStatus | '',
    genre: '',
    limit: 10,
  })

  const fetchManga = useCallback(async () => {
    try {
      setLoading(true)
      const params: GetMangaParams = {
        title: filters.searchTerm,
        status: filters.status || undefined,
        genre: filters.genre || undefined,
        limit: filters.limit,
        offset: (currentPage - 1) * filters.limit,
      }
      const response = await apiClient.getManga(params)
      setMangaList(response.data.items || [])
      setTotalPages(Math.ceil((response.meta?.total || 0) / filters.limit))
    } catch (err: any) {
      setError(err.message || 'Failed to fetch manga.')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [currentPage, filters])

  const handleSearch = useCallback((newFilters: typeof filters) => {
    setCurrentPage(1) // Reset to first page on new search
    setFilters(newFilters)
  }, [])

  const handleCardClick = (manga: Manga) => {
    setSelectedManga(manga)
  }

  const handleCloseModal = () => {
    setSelectedManga(null)
  }

  const handleAddToLibrary = async (mangaId: string) => {
    try {
      await apiClient.addToLibrary(mangaId)
      toast.success('Manga added to your library!')
      handleCloseModal() // Close modal after adding
    } catch (err: any) {
      toast.error(err.message || 'Failed to add manga to library.')
      console.error('Add to library error:', err)
    }
  }

  useEffect(() => {
    fetchManga()
  }, [fetchManga])

  return (
    <div className="min-h-screen bg-gray-100">
      <Header />
      <main className="container mx-auto px-6 py-8">
        <SearchAndFilter onSearch={handleSearch} />
        {loading && (
          <div className="flex justify-center">
            <ClipLoader size={100} />
          </div>
        )}
        {error && <p className="text-center text-red-500">{error}</p>}
        {!loading && !error && mangaList.length > 0 ? (
          <>
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-6">
              {mangaList.map((manga) => (
                <MangaCard
                  key={manga.id}
                  manga={manga}
                  onClick={() => handleCardClick(manga)}
                />
              ))}
            </div>
            <div className="mt-8 flex justify-center items-center space-x-4">
              <button
                onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                disabled={currentPage === 1}
                className="px-4 py-2 bg-indigo-600 text-white rounded-md disabled:bg-gray-400"
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
                className="px-4 py-2 bg-indigo-600 text-white rounded-md disabled:bg-gray-400"
              >
                Next
              </button>
            </div>
          </>
        ) : (
          !loading && (
            <p className="text-center text-gray-500">No manga found.</p>
          )
        )}
      </main>
      <MangaModal
        manga={selectedManga}
        onClose={handleCloseModal}
        onAddToLibrary={handleAddToLibrary}
      />
      <ToastContainer
        position="bottom-right"
        autoClose={3000}
        hideProgressBar={false}
        newestOnTop={false}
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
      />
    </div>
  )
}

export default withAuth(HomePage)
