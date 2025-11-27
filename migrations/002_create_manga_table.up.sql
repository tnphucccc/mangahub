-- Create manga table for manga catalog
CREATE TABLE IF NOT EXISTS manga (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT,
    genres TEXT,  -- JSON array stored as TEXT
    status TEXT CHECK(status IN ('ongoing', 'completed', 'hiatus', 'cancelled')),
    total_chapters INTEGER DEFAULT 0,
    description TEXT,
    cover_image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for title searches
CREATE INDEX IF NOT EXISTS idx_manga_title ON manga(title);

-- Index for status filtering
CREATE INDEX IF NOT EXISTS idx_manga_status ON manga(status);
