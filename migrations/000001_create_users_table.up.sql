CREATE TABLE IF NOT EXISTS users (
    id         serial PRIMARY KEY,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    name       text NOT NULL,
    email      text NOT NULL
);