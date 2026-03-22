CREATE TABLE IF NOT EXISTS meals (
    id         serial PRIMARY KEY,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    user_id    int NOT NULL REFERENCES users(id),
    food_name  text NOT NULL,
    calories   int NOT NULL
);