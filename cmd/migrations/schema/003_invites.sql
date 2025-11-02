-- +migrate Up
CREATE TYPE invite_status AS ENUM (
    'sent',
    'accepted',
    'declined'
);

CREATE TABLE city_invites (
    id          UUID PRIMARY KEY,
    city_id     UUID NOT NULL REFERENCES city(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL,
    role        city_moderator_roles NOT NULL,
    status      invite_status        NOT NULL DEFAULT 'sent',

    created_at  TIMESTAMPTZ        NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ        NOT NULL DEFAULT now(),

    CONSTRAINT chk_city_invites_expires_future CHECK (expires_at > now())
);

CREATE INDEX IF NOT EXISTS idx_city_invites_city_status ON city_invites (city_id, status);
CREATE INDEX IF NOT EXISTS idx_city_invites_user_status ON city_invites (user_id, status);

-- +migrate Down
DROP INDEX IF EXISTS idx_city_invites_city_status;
DROP INDEX IF EXISTS idx_city_invites_user_status;

DROP TABLE IF EXISTS city_invites CASCADE;
DROP TYPE IF EXISTS invite_status;
DROP TYPE IF EXISTS city_invite_target;
