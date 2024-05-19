CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY,
    date TEXT UNIQUE NOT NULL,
    explanation TEXT,
    media_type VARCHAR(50),
    title VARCHAR(255),
    data BYTEA NOT NULL
);  