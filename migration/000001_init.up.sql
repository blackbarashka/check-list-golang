CREATE TABLE IF NOT EXISTS tasksTable(
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description VARCHAR(300) NOT NULL,
    status BOOLEAN,
    addTime TIMESTAMP,
    completedTime TIMESTAMP
)