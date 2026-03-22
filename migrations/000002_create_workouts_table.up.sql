CREATE TABLE IF NOT EXISTS workouts (
    id               serial PRIMARY KEY,
    created_at       timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    user_id          int NOT NULL REFERENCES users(id),
    type             text NOT NULL,
    duration_minutes int NOT NULL,
    calories_burned  int NOT NULL
);