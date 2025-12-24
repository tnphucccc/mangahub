'use client'

import { useState, useEffect, useCallback } from 'react'
import withAuth from '../hoc/withAuth'
import { apiClient } from '../../lib/apiClient'
import type {
  UserProgressWithManga,
  ProgressUpdateRequest,
} from '../../../../../packages/types/src'
import Header from '../../component/header'
import LibraryMangaCard from '../../component/libraryMangaCard'
import { ToastContainer, toast } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'
import { ClipLoader } from 'react-spinners'

function LibraryPage() {
  const [library, setLibrary] = useState<UserProgressWithManga[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchLibrary = useCallback(async () => {
    try {
      setLoading(true)
      const response = await apiClient.getLibrary()
      setLibrary(response.data.items || [])
    } catch (err: any) {
      setError(err.message || 'Failed to fetch library.')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchLibrary()
  }, [fetchLibrary])

  const handleUpdateProgress = async (
    mangaId: string,
    progress: ProgressUpdateRequest
  ) => {
    try {
      await apiClient.updateMangaProgress(mangaId, progress)
      toast.success('Progress updated successfully!')
      // Optionally, refetch or update state locally
      setLibrary((prevLibrary) =>
        prevLibrary.map((item) =>
          item.manga_id === mangaId ? { ...item, ...progress } : item
        )
      )
    } catch (err: any) {
      toast.error(err.message || 'Failed to update progress.')
      console.error(err)
    }
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <Header />
      <main className="container mx-auto px-6 py-8">
        <h2 className="text-3xl font-bold text-gray-900 mb-8">My Library</h2>
        {loading && (
          <div className="flex justify-center">
            <ClipLoader size={100} />
          </div>
        )}
        {error && <p className="text-center text-red-500">{error}</p>}
        {!loading && !error && library.length > 0 ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
            {library.map((item) => (
              <LibraryMangaCard
                key={item.manga_id}
                item={item}
                onUpdate={handleUpdateProgress}
              />
            ))}
          </div>
        ) : (
          !loading && (
            <p className="text-center text-gray-600">Your library is empty.</p>
          )
        )}
      </main>
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

export default withAuth(LibraryPage)
