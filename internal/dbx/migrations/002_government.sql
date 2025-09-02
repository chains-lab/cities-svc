-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE city_gov_roles AS ENUM (
  'mayor',
  'advisor',
  'member',
  'moderator'
);

CREATE TABLE city_govs (
    user_id        UUID              NOT NULL PRIMARY KEY,
    city_id        UUID              NOT NULL REFERENCES city(id)   ON DELETE CASCADE,
    role           city_gov_roles    NOT NULL,
    label          VARCHAR(255),
    created_at     TIMESTAMP         NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at     TIMESTAMP         NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
);

CREATE TYPE invite_status AS ENUM (
    'sent',
    'accepted',
    'rejected'
);

CREATE TABLE invites (
    id           UUID PRIMARY KEY,
    status       status NOT NULL DEFAULT 'sent',
    role         city_gov_roles NOT NULL,
    city_id      UUID NOT NULL REFERENCES city(id) ON DELETE CASCADE,
    initiator_id UUID NOT NULL,
    user_id      UUID,
    answered_at  TIMESTAMP,
    expires_at   TIMESTAMP NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    CONSTRAINT invite_status_answered_ck CHECK (
        (status = 'sent'     AND answered_at IS NULL) OR
        (status IN ('accepted','rejected') AND answered_at IS NOT NULL)
    )
)

CREATE UNIQUE INDEX city_gov_unique_mayor_active
    ON city_governments(city_id)
    WHERE role = 'mayor' AND status = 'active';

-- +migrate Down
DROP INDEX IF EXISTS city_gov_unique_mayor_active;

DROP TABLE IF EXISTS city_governments CASCADE;
DROP TABLE IF EXISTS invites CASCADE;

DROP TYPE IF EXISTS city_gov_roles;
DROP TYPE IF EXISTS invite_status;

DROP EXTENSION IF EXISTS "uuid-ossp";
