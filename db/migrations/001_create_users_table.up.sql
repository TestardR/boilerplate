CREATE TABLE IF NOT EXISTS users(
   id UUID PRIMARY KEY,
   username VARCHAR (50) NOT NULL,
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);