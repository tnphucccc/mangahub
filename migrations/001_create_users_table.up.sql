-- Create users table for authentication
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster username lookups during login
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Index for email lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
