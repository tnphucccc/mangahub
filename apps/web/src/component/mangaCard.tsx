import React from 'react'
import type { Manga } from '../../../../packages/types/src'
import defaultCover from '@/../public/assets/bookcover_cover.png'
import { upperCaseFirstLetter } from '@/app/helpers/upperCaseFirstLetter'
import Image from 'next/image'

interface MangaCardProps {
  manga: Manga
  onClick: () => void
}

const MangaCard = ({ manga, onClick }: MangaCardProps) => {
  // Use a placeholder if the cover image is missing
  const coverImage = manga.cover_image_url || defaultCover.src
  return (
    <button
      onClick={onClick}
      className="bg-white rounded-lg shadow-md overflow-hidden transform hover:-translate-y-1 transition-transform duration-300 cursor-pointer text-left w-full"
    >
      <Image
        width={400}
        height={600}
        src={coverImage}
        alt={`Cover for ${manga.title}`}
        className="w-full h-72 object-cover"
      />
      <div className="p-4">
        <h3
          className="text-lg font-bold text-gray-900 truncate"
          title={manga.title}
        >
          {manga.title}
        </h3>
        <p className="text-sm text-gray-600 mt-1">
          {manga.author || 'Unknown Author'}
        </p>
        <span
          className={`mt-2 inline-block px-2 py-1 text-xs font-semibold rounded-full ${
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
      </div>
    </button>
  )
}

export default MangaCard
