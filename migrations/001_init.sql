-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE users (
                       id UUID PRIMARY KEY,
                       avatar_url VARCHAR DEFAULT NULL,
                       name VARCHAR NOT NULL,
                       bio TEXT DEFAULT NULL,
                       telegram VARCHAR NOT NULL UNIQUE,
                       password bytea,
                       banned BOOLEAN DEFAULT FALSE
);


CREATE TABLE groups (
                        id UUID PRIMARY KEY,
                        avatar_url VARCHAR DEFAULT NULL,
                        name VARCHAR NOT NULL,
                        invite_code VARCHAR DEFAULT NULL
);

CREATE TABLE roles (
                       user_id UUID REFERENCES users(id),
                        group_id UUID REFERENCES groups(id),
                       role VARCHAR NOT NULL,
                       UNIQUE (user_id, group_id,role)
);

CREATE TABLE help_requests (
                               id UUID PRIMARY KEY ,
                               user_id UUID  REFERENCES users(id),
                               mentor_id UUID REFERENCES users(id),
                               group_id UUID REFERENCES groups(id),
                               goal VARCHAR NOT NULL,
                               bio TEXT DEFAULT NULL,
                               status VARCHAR NOT NULL
);

CREATE TABLE fast_helps (
                            id UUID PRIMARY KEY,
                            user_id UUID  REFERENCES users(id),
                            question VARCHAR NOT NULL,
                            status BOOLEAN NOT NULL
);


CREATE TABLE pairs (
                       user_id UUID REFERENCES users(id),
                       mentor_id UUID REFERENCES users(id),
                       group_id UUID REFERENCES groups(id),
                       goal VARCHAR NOT NULL,
                   UNIQUE (user_id,mentor_id,group_id)
);


-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back

-- Add DROP statements for the UP migration here
