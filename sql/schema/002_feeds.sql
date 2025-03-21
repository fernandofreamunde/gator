-- +goose Up
CREATE TABLE feeds (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name VARCHAR(256),
	url VARCHAR(256) UNIQUE,
	user_id UUID,
	CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
