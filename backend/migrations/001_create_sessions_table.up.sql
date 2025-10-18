-- Create sessions table for scs session management
CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

-- Create index for session expiry cleanup
CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);
