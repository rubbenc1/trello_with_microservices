CREATE TABLE IF NOT EXISTS tasks_values (
    id SERIAL PRIMARY KEY,
    type TEXT,
    string_value TEXT,
    number_value INTEGER,
    column_id INTEGER,
    task_id INTEGER,
    option_id INTEGER
);
