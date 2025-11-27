-- Create user_progress table for tracking reading progress
CREATE TABLE IF NOT EXISTS user_progress (
    user_id TEXT NOT NULL,
    manga_id TEXT NOT NULL,
    current_chapter INTEGER DEFAULT 0,
    status TEXT CHECK(status IN ('reading', 'completed', 'plan_to_read', 'on_hold', 'dropped')),
    rating INTEGER CHECK(rating >= 1 AND rating <= 10),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, manga_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (manga_id) REFERENCES manga(id) ON DELETE CASCADE
);

-- Index for user's library queries
CREATE INDEX IF NOT EXISTS idx_user_progress_user_id ON user_progress(user_id);

-- Index for manga popularity queries
CREATE INDEX IF NOT EXISTS idx_user_progress_manga_id ON user_progress(manga_id);

-- Index for filtering by reading status
CREATE INDEX IF NOT EXISTS idx_user_progress_status ON user_progress(user_id, status);
