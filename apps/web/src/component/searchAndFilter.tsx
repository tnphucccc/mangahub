import { useDebounce } from '@/app/helpers/hook'
import { useEffect, useState } from 'react'
import { MangaStatus } from '../../../../packages/types/src'

export const SearchAndFilter = ({
  onSearch,
}: {
  onSearch: (filters: {
    searchTerm: string
    status: MangaStatus | ''
    genre: string
    limit: number
  }) => void
}) => {
  const [searchTerm, setSearchTerm] = useState('')
  const [status, setStatus] = useState<MangaStatus | ''>('')
  const [genre, setGenre] = useState('')
  const [limit, setLimit] = useState(10)

  const debouncedSearchTerm = useDebounce(searchTerm, 500)
  const debouncedGenre = useDebounce(genre, 500)

  useEffect(() => {
    onSearch({
      searchTerm: debouncedSearchTerm,
      status,
      genre: debouncedGenre,
      limit,
    })
  }, [debouncedSearchTerm, status, debouncedGenre, limit, onSearch])

  return (
    <div className="bg-white p-4 rounded-lg shadow-md mb-8">
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
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
        <select
          value={limit}
          onChange={(e) => setLimit(Number(e.target.value))}
          className="p-2 border rounded-md text-black"
        >
          <option value="10">10 items per page</option>
          <option value="20">20 items per page</option>
          <option value="30">30 items per page</option>
        </select>
      </div>
    </div>
  )
}
