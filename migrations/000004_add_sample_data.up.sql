INSERT INTO users (name, email) VALUES
    ('John Doe',     'john@email.com'),
    ('Jane Smith',   'jane@email.com'),
    ('Carlos Reyes', 'carlos@email.com');

INSERT INTO workouts (user_id, type, duration_minutes, calories_burned) VALUES
    (1, 'Running',  30, 300),
    (2, 'Cycling',  45, 400),
    (3, 'Swimming', 60, 500);

INSERT INTO meals (user_id, food_name, calories) VALUES
    (1, 'Grilled Chicken', 350),
    (2, 'Caesar Salad',    200),
    (3, 'Brown Rice Bowl', 450);