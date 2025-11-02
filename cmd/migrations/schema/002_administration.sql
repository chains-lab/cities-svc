-- +migrate Up
CREATE TYPE city_moderator_roles AS ENUM (
    'chief'
    'government',
    'moderator',
);

CREATE TABLE city_administration (
    user_id    UUID      PRIMARY KEY,
    city_id    UUID      NOT NULL REFERENCES city(id) ON DELETE CASCADE,

    position   VARCHAR(255),
    label      VARCHAR(255),
    role       city_administration_roles NOT NULL,

    updated_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

-- +migrate Down
DROP TABLE IF EXISTS city_administration CASCADE;
DROP TYPE IF EXISTS city_administration_roles;

DROP EXTENSION IF EXISTS "uuid-ossp";
