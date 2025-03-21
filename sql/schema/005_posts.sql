-- +goose Up
CREATE TABLE posts (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	title VARCHAR(256),
	url VARCHAR(256) UNIQUE,
	description TEXT,
	published_at TIMESTAMP NOT NULL,
	feed_id UUID,
	CONSTRAINT fk_feed_id FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
