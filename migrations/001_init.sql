-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE users (
    id UUID PRIMARY KEY,
    avatar_url VARCHAR DEFAULT NULL,
    name VARCHAR NOT NULL UNIQUE,
    bio TEXT,
    telegram VARCHAR
);

CREATE TABLE groups (
    id UUID PRIMARY KEY,
    avatar_url VARCHAR DEFAULT NULL,
    name VARCHAR NOT NULL
);

CREATE TABLE roles (
    group_id UUID REFERENCES groups(id),
    user_id UUID REFERENCES users(id),
    role VARCHAR NOT NULL,
    PRIMARY KEY (group_id, user_id)
);

CREATE TABLE help_requests (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    pending_mentor_ids UUID[],
    active_mentor_ids UUID[],
    max_mentors_count INT NOT NULL,
    goal VARCHAR NOT NULL,
    status TEXT NOT NULL
);

CREATE TABLE pairs (
    user_id UUID REFERENCES users(id),
    mentor_id UUID REFERENCES users(id),
    goal VARCHAR NOT NULL,
    PRIMARY KEY (user_id, mentor_id)
);

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back

-- Add DROP statements for the UP migration here
