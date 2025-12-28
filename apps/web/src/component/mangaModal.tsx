import React, { useState } from 'react'
import type { Manga } from '../../../../packages/types/src'
import defaultCover from '@/../public/assets/bookcover_cover.png'
import { upperCaseFirstLetter } from '@/app/helpers/upperCaseFirstLetter'
import Image from 'next/image'

interface MangaModalProps {
  manga: Manga | null
  onClose: () => void
  onAddToLibrary: (mangaId: string) => void
  onSendNotification: (newChapNum: number, newChapTitle: string) => void
}

const MangaModal = ({
  manga,
  onClose,
  onAddToLibrary,
  onSendNotification,
}: MangaModalProps) => {
  const [showNewChapterForm, setShowNewChapterForm] = useState(false)
  const [newChapterNumber, setNewChapterNumber] = useState(0)
  const [newChapterTitle, setNewChapterTitle] = useState('')

  if (!manga) {
    return null
  }

  const handleSendNotification = () => {
    onSendNotification(newChapterNumber, newChapterTitle)
    setShowNewChapterForm(false)
    setNewChapterNumber(0)
    setNewChapterTitle('')
  }

  // Stop propagation to prevent clicks inside the modal from closing it
  const handleModalContentClick = (e: React.MouseEvent) => {
    e.stopPropagation()
  }

  const coverImage = manga.cover_image_url?.includes('mangadex')
    ? manga.cover_image_url
    : defaultCover.src

  return (
    // Backdrop
    <div
      className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
      onClick={onClose}
    >
      {/* Modal Content */}
      <div
        className="bg-white rounded-lg shadow-xl w-full max-w-4xl max-h-[70vh] flex flex-col"
        onClick={handleModalContentClick}
      >
        <div className="flex flex-col md:flex-row grow overflow-hidden">
          <Image
            width={400}
            height={600}
            src={coverImage!}
            alt={`Cover for ${manga.title}`}
            className="h-full object-cover"
          />
          <div className="p-6 flex flex-col flex-1 overflow-y-auto relative">
            <button
              onClick={onClose}
              className="absolute top-4 right-4 text-gray-500 hover:text-gray-800"
              aria-label="Close modal"
            >
              <svg
                className="w-6 h-6"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">
              {manga.title}
            </h2>
            <p className="text-lg text-gray-700 mb-4">
              by {manga.author || 'Unknown Author'}
            </p>

            <div className="flex items-center mb-4">
              <span
                className={`inline-block px-3 py-1 text-sm font-semibold rounded-full ${
                  manga.status === 'ongoing'
                    ? 'bg-blue-200 text-blue-800'
                    : manga.status === 'completed'
                      ? 'bg-green-200 text-green-800'
                      : manga.status === 'hiatus'
                        ? 'bg-yellow-200 text-yellow-800'
                        : 'bg-red-200 text-red-800'
                }`}
              >
                {upperCaseFirstLetter(manga.status || 'unknown')}
              </span>
              <span className="ml-4 text-sm text-gray-600">
                {manga.total_chapters || 'N/A'} Chapters
              </span>
            </div>

            <div className="mb-4">
              <h4 className="font-semibold text-gray-800">Genres</h4>
              <div className="flex flex-wrap gap-2 mt-2">
                {manga.genres && manga.genres.length > 0 ? (
                  manga.genres.map((genre) => (
                    <span
                      key={genre}
                      className="bg-gray-200 text-gray-800 px-2 py-1 text-xs rounded-full"
                    >
                      {genre}
                    </span>
                  ))
                ) : (
                  <p className="text-sm text-gray-500">No genres listed.</p>
                )}
              </div>
            </div>

            <div>
              <h4 className="font-semibold text-gray-800">Description</h4>
              <p className="text-base text-gray-600 mt-2 whitespace-pre-wrap">
                {manga.description || 'No description available.'}
              </p>
            </div>
          </div>
        </div>
        <div className="p-4 border-t flex flex-col">
          {showNewChapterForm ? (
            <div className="flex flex-col space-y-4">
              <div className="flex items-center space-x-4">
                <input
                  type="number"
                  placeholder="Chapter Number"
                  value={newChapterNumber}
                  onChange={(e) => setNewChapterNumber(Number(e.target.value))}
                  className="w-1/3 px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm text-black"
                />
                <input
                  type="text"
                  placeholder="Chapter Title (optional)"
                  value={newChapterTitle}
                  onChange={(e) => setNewChapterTitle(e.target.value)}
                  className="grow px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm text-black"
                />
              </div>
              <div className="flex justify-end space-x-2">
                <button
                  onClick={() => setShowNewChapterForm(false)}
                  className="px-4 py-2 bg-gray-300 text-black rounded-md hover:bg-gray-400"
                >
                  Cancel
                </button>
                <button
                  onClick={handleSendNotification}
                  className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700"
                >
                  Send Notification
                </button>
              </div>
            </div>
          ) : (
            <div className="flex justify-end">
              <button
                onClick={() => onAddToLibrary(manga.id)}
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
              >
                Add to Library
              </button>
              <button
                onClick={() => setShowNewChapterForm(true)}
                className="ml-4 px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700"
              >
                Add New Chapter
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default MangaModal
