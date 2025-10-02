-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE city_admins_roles AS ENUM (
    'head',
    'member',
    'moderator'
);

CREATE TABLE city_admins (
    user_id    UUID           PRIMARY KEY,
    city_id    UUID           NOT NULL REFERENCES city(id) ON DELETE CASCADE,
    role       city_admins_roles NOT NULL,
    label      VARCHAR(255),
    created_at TIMESTAMP      NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP      NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE TYPE invite_status AS ENUM (
    'sent',
    'accepted'
);

CREATE TABLE invites (
    id          UUID           PRIMARY KEY,
    status      invite_status  NOT NULL DEFAULT 'sent',
    role        city_admins_roles NOT NULL,
    city_id     UUID           NOT NULL REFERENCES city(id) ON DELETE CASCADE,
    token       VARCHAR        NOT NULL UNIQUE,
    user_id     UUID,
    answered_at TIMESTAMP,
    expires_at  TIMESTAMP      NOT NULL,
    created_at  TIMESTAMP      NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),

    CONSTRAINT invite_status_answered_ck CHECK (
        (status = 'sent' AND answered_at IS NULL)
            OR
        (status = 'accepted' AND answered_at IS NOT NULL)
    )
);

-- +migrate Down
DROP TABLE IF EXISTS invites CASCADE;
DROP TABLE IF EXISTS city_admins CASCADE;

DROP TYPE IF EXISTS invite_status;
DROP TYPE IF EXISTS city_admins_roles;

DROP EXTENSION IF EXISTS "uuid-ossp";
