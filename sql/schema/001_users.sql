-- +goose Up
CREATE TABLE users(
  ID UUID,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT UNIQUE
);

-- +goose Down
DROP TABLE users;