-- Rollback user_progress table
DROP INDEX IF EXISTS idx_user_progress_status;
DROP INDEX IF EXISTS idx_user_progress_manga_id;
DROP INDEX IF EXISTS idx_user_progress_user_id;
DROP TABLE IF EXISTS user_progress;
