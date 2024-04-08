DROP TABLE IF EXISTS tasks_labels, tasks, labels, users;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE labels (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    opened BIGINT NOT NULL DEFAULT extract(epoch from now()),
    closed BIGINT DEFAULT 0,
    author_id INTEGER REFERENCES users(id) DEFAULT 0,
    assigned_id INTEGER REFERENCES users(id) DEFAULT 0,
    title TEXT,
    content TEXT
);

CREATE TABLE tasks_labels (
    task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    label_id INTEGER REFERENCES labels(id)
);

-- стартовые данные
INSERT INTO users (id, name) VALUES (0, 'default');
INSERT INTO users (id, name) VALUES (1, 'user01');
INSERT INTO users (id, name) VALUES (2, 'user02');
