'use client'

import React, { useState } from 'react'
import type {
  UserProgressWithManga,
  ReadingStatus,
} from '../../../../packages/types/src'
import { upperCaseFirstLetter } from '@/app/helpers/upperCaseFirstLetter'
import defaultCover from '@/../public/assets/bookcover_cover.png'

interface LibraryMangaCardProps {
  item: UserProgressWithManga
  onUpdate: (
    mangaId: string,
    progress: {
      current_chapter: number
      status: ReadingStatus
      rating?: number
    }
  ) => Promise<void>
}

const LibraryMangaCard = ({ item, onUpdate }: LibraryMangaCardProps) => {
  const { manga } = item

  // Initialize state from the prop
  const [currentChapter, setCurrentChapter] = useState(item.current_chapter)
  const [status, setStatus] = useState<ReadingStatus>(item.status)
  const [rating, setRating] = useState(
    item.rating.Int64 ? Number(item.rating.Int64) : 1
  )

  const [isSaving, setIsSaving] = useState(false)

  if (!manga) return null

  const handleUpdateClick = async () => {
    setIsSaving(true)
    await onUpdate(manga.id, {
      current_chapter: Number(currentChapter),
      status,
      rating: Number(rating),
    })
    setIsSaving(false)
  }

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden flex flex-col">
      <img
        src={defaultCover.src}
        alt={`Cover for ${manga.title}`}
        className="w-full h-48 object-cover"
      />
      <div className="p-4 flex flex-col grow">
        <h3
          className="text-md font-bold text-gray-900 truncate"
          title={manga.title}
        >
          {manga.title}
        </h3>
        <p className="text-sm text-gray-600">
          {manga.author || 'Unknown Author'}
        </p>

        <div className="mt-4 space-y-3 grow">
          <div>
            <label className="block text-xs font-medium text-gray-600">
              Chapter (Lastest chapter {manga.total_chapters || 'N/A'})
            </label>
            <input
              type="number"
              value={currentChapter}
              onChange={(e) => setCurrentChapter(Number(e.target.value))}
              className="mt-1 p-1 w-full border rounded-md text-sm text-black"
            />
          </div>

          <div>
            <label className="block text-xs font-medium text-gray-600">
              Status
            </label>
            <select
              value={status}
              onChange={(e) => setStatus(e.target.value as ReadingStatus)}
              className="mt-1 p-1 w-full border rounded-md text-sm text-black"
            >
              <option value="reading">Reading</option>
              <option value="completed">Completed</option>
              <option value="plan_to_read">Plan to Read</option>
              <option value="on_hold">On Hold</option>
              <option value="dropped">Dropped</option>
            </select>
          </div>

          <div>
            <label className="block text-xs font-medium text-gray-600">
              Rating (1-10)
            </label>
            <input
              type="number"
              min="0"
              max="10"
              value={rating}
              onChange={(e) => setRating(Number(e.target.value))}
              className="mt-1 p-1 w-full border rounded-md text-sm text-black"
            />
          </div>
        </div>

        <button
          onClick={handleUpdateClick}
          disabled={isSaving}
          className="mt-4 w-full px-3 py-2 text-sm font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700 disabled:bg-gray-400"
        >
          {isSaving ? 'Saving...' : 'Update'}
        </button>
      </div>
    </div>
  )
}

export default LibraryMangaCard
