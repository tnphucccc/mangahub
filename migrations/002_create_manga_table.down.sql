-- Rollback manga table
DROP INDEX IF EXISTS idx_manga_status;
DROP INDEX IF EXISTS idx_manga_title;
DROP TABLE IF EXISTS manga;
