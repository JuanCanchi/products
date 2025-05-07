CREATE TABLE products (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    price NUMERIC(10, 2),
    image_url TEXT,
    user_id TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

