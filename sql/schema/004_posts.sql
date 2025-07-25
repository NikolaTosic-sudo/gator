-- +goose Up
CREATE TABLE posts(
  ID SERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  title TEXT NOT NULL,
  url TEXT UNIQUE NOT NULL,
  description TEXT,
  published_at TIMESTAMP,
  feed_id UUID NOT NULL,
  FOREIGN KEY (feed_id)
  REFERENCES feeds(ID)
  ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
