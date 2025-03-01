-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE users (
    id UUID PRIMARY KEY,
    avatar_url VARCHAR DEFAULT NULL,
    name VARCHAR NOT NULL UNIQUE,
    bio TEXT DEFAULT NULL,
    telegram VARCHAR NOT NULL
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
    id UUID,
    user_id UUID  REFERENCES users(id),
    mentor_id UUID REFERENCES users(id),
    group_id UUID REFERENCES groups(id),
    goal VARCHAR NOT NULL,
    bio TEXT DEFAULT NULL,
    status VARCHAR NOT NULL,
    PRIMARY KEY (id,user_id,mentor_id)
);

CREATE TABLE pairs (
    user_id UUID REFERENCES users(id),
    mentor_id UUID REFERENCES users(id),
    group_id UUID REFERENCES users(id),
    goal VARCHAR NOT NULL,
    PRIMARY KEY (user_id, mentor_id)
);

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back

-- Add DROP statements for the UP migration here
