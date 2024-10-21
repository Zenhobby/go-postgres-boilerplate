CREATE TABLE IF NOT EXISTS persons (
    id SERIAL PRIMARY KEY,
    uid VARCHAR(36) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    traits JSONB
);

CREATE INDEX IF NOT EXISTS idx_persons_name ON persons(name);
CREATE INDEX IF NOT EXISTS idx_persons_uid ON persons(uid);