-- +goose Up

CREATE TABLE feeds(
  ID UUID NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  url TEXT UNIQUE NOT NULL,
  user_id UUID NOT NULL,
  FOREIGN KEY (user_id)
  REFERENCES users(ID)
  ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
