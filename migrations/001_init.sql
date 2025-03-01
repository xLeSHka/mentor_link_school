-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS groups (
    id UUID PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    avatar_url VARCHAR DEFAULT NULL
);
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    first_name VARCHAR NOT NULL,
    second_name VARCHAR NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    avatar_url VARCHAR DEFAULT NULL,
    bio TEXT DEFAULT NULL,
    password BYTEA NOT NULL
);
CREATE TABLE get_mentor_requests (
    user_id UUID REFERENCES users(id),
    group_id UUID REFERENCES groups(id),
    goal TEXT NOT NULL,
    bio TEXT DEFAULT NULL
);
CREATE TABLE create_mentor_requests (
    user_id UUID REFERENCES users(id),
    group_id UUID REFERENCES groups(id),
    goal TEXT NOT NULL,
    bio TEXT DEFAULT NULL
);
CREATE TABLE roles (
    user_id UUID REFERENCES users(id),
    group_id UUID REFERENCES groups(id),
    role VARCHAR NOT NULL
);
CREATE TABLE mentors (
    user_id UUID REFERENCES users(id),
    mentor_id UUID REFERENCES users(id)
);
-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back

-- Add DROP statements for the UP migration here
