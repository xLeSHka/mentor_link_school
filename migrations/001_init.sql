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
    name VARCHAR NOT NULL,
    invite_code VARCHAR DEFAULT NULL
);

CREATE TABLE roles (
    group_id UUID REFERENCES groups(id),
    user_id UUID REFERENCES users(id),
    role VARCHAR NOT NULL
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
    goal VARCHAR NOT NULL
);


INSERT INTO users (id, name, telegram) VALUES
('f9e0a67c-ee38-4fac-8a32-ad69d36ba661', 'Демо Пользователь 1', 't_prodano'),
('c0f998f2-63b5-49ed-b9b5-35f0815bbac3', 'Демо Пользователь 2', 't_prodano'),
('a5dcd87c-37c5-4504-9ff6-7b32c6cce139', 'Демо Пользователь 3', 't_prodano'),
('58e016b0-8956-44df-939d-5418f3f5d234', 'Демо Пользователь 4', 't_prodano'),
('7c176601-b275-4f21-9b8c-19410679c395', 'Демо Пользователь 5', 't_prodano'),
('96d39048-36d9-4349-96da-b5119be10b34', 'Демо Пользователь 6', 't_prodano'),
('01dc08e8-fc5c-469b-bdcf-e826c7b497b5', 'Демо Пользователь 7', 't_prodano'),
('1a575bec-6e2b-41bf-bb34-04a93c43e1f8', 'Демо Пользователь 8', 't_prodano'),
('b234d17c-21b2-458f-9641-3a5ff9510dd4', 'Демо Пользователь 9', 't_prodano'),
('0ed523b1-9827-433c-b969-67686e8b812b', 'Демо Пользователь 10', 't_prodano');


INSERT INTO groups (id, name, invite_code) VALUES
('74878d9b-e348-40a0-a4a5-a07515864760', 'Организация 1', 'gnjkf'),
('24406b8c-f0cc-4e31-9970-9f31f1b232e1', 'Организация 2', 'gnjkq');

INSERT INTO roles (group_id, user_id, role) VALUES
('74878d9b-e348-40a0-a4a5-a07515864760', 'f9e0a67c-ee38-4fac-8a32-ad69d36ba661', 'owner'),
('24406b8c-f0cc-4e31-9970-9f31f1b232e1', 'c0f998f2-63b5-49ed-b9b5-35f0815bbac3', 'owner'),
('74878d9b-e348-40a0-a4a5-a07515864760', 'a5dcd87c-37c5-4504-9ff6-7b32c6cce139', 'student'),
('24406b8c-f0cc-4e31-9970-9f31f1b232e1', '58e016b0-8956-44df-939d-5418f3f5d234', 'student'),
('74878d9b-e348-40a0-a4a5-a07515864760', '58e016b0-8956-44df-939d-5418f3f5d234', 'student'),
('24406b8c-f0cc-4e31-9970-9f31f1b232e1', '7c176601-b275-4f21-9b8c-19410679c395', 'mentor'),
('74878d9b-e348-40a0-a4a5-a07515864760', '96d39048-36d9-4349-96da-b5119be10b34', 'student-mentor'),
('24406b8c-f0cc-4e31-9970-9f31f1b232e1', '96d39048-36d9-4349-96da-b5119be10b34', 'mentor'),
('74878d9b-e348-40a0-a4a5-a07515864760', '01dc08e8-fc5c-469b-bdcf-e826c7b497b5', 'mentor'),
('24406b8c-f0cc-4e31-9970-9f31f1b232e1', '01dc08e8-fc5c-469b-bdcf-e826c7b497b5', 'mentor'),
('24406b8c-f0cc-4e31-9970-9f31f1b232e1', '1a575bec-6e2b-41bf-bb34-04a93c43e1f8', 'student'),
('24406b8c-f0cc-4e31-9970-9f31f1b232e1', 'b234d17c-21b2-458f-9641-3a5ff9510dd4', 'student-mentor'),
('24406b8c-f0cc-4e31-9970-9f31f1b232e1', '0ed523b1-9827-433c-b969-67686e8b812b', 'mentor');

INSERT INTO users (id, name, telegram) VALUES
('1b63ef07-f7bb-41d4-8289-ccd705dfcf59', 'Сваггер Пользователь 1', 't_prodano'),
('0dc729df-9cad-4500-81b4-a20d6b18b691', 'Сваггер Пользователь 2', 't_prodano'),
('e02cb2f1-b69f-4805-8a6c-e96c68c2a0b5', 'Сваггер Пользователь 3', 't_prodano'),
('ba24317a-14a4-433f-b441-493c1a912f1f', 'Сваггер Пользователь 4', 't_prodano'),
('dbe8ebfb-5aa8-4c80-9c48-4466014bc4ce', 'Сваггер Пользователь 5', 't_prodano');

INSERT INTO groups (id, name, invite_code) VALUES
('9aff3974-8bb7-44ec-8a88-6415df608044', 'Организация 1', 'gnnkf');

INSERT INTO roles (group_id, user_id, role) VALUES
('9aff3974-8bb7-44ec-8a88-6415df608044', '1b63ef07-f7bb-41d4-8289-ccd705dfcf59', 'owner'),
('9aff3974-8bb7-44ec-8a88-6415df608044', '0dc729df-9cad-4500-81b4-a20d6b18b691', 'student'),
('9aff3974-8bb7-44ec-8a88-6415df608044', 'e02cb2f1-b69f-4805-8a6c-e96c68c2a0b5', 'student'),
('9aff3974-8bb7-44ec-8a88-6415df608044', 'ba24317a-14a4-433f-b441-493c1a912f1f', 'mentor'),
('9aff3974-8bb7-44ec-8a88-6415df608044', 'dbe8ebfb-5aa8-4c80-9c48-4466014bc4ce', 'mentor');
-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back

-- Add DROP statements for the UP migration here
