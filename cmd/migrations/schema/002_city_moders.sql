-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE city_moders_roles AS ENUM (
  'mayor',
  'advisor',
  'member',
  'moderator'
);

CREATE TABLE city_moders (
    user_id    UUID           PRIMARY KEY,
    city_id    UUID           NOT NULL REFERENCES city(id) ON DELETE CASCADE,
    role       city_moders_roles NOT NULL,
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
    role        city_moders_roles NOT NULL,
    city_id     UUID           NOT NULL REFERENCES city(id) ON DELETE CASCADE,
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

CREATE UNIQUE INDEX city_gov_unique_mayor
    ON city_moders (city_id)
    WHERE role = 'mayor';

-- +migrate Down
DROP INDEX IF EXISTS city_moders_unique_mayor;

DROP TABLE IF EXISTS invites CASCADE;
DROP TABLE IF EXISTS city_moders CASCADE;

DROP TYPE IF EXISTS invite_status;
DROP TYPE IF EXISTS city_moders_roles;

DROP EXTENSION IF EXISTS "uuid-ossp";
