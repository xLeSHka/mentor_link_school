-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL,
    surname VARCHAR NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    avatar_url VARCHAR DEFAULT NULL,
    age INT NOT NULL,
    password BYTEA NOT NULL
);

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back

-- Add DROP statements for the UP migration here
