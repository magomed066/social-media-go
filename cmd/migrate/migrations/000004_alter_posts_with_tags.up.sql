ALTER TABLE posts
    ADD COLUMN IF NOT EXISTS tags text[],
    ADD COLUMN IF NOT EXISTS created_at timestamp(0) with time zone NOT NULL DEFAULT NOW();
